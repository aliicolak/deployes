package handlers

import (
	"encoding/json"
	"net/http"

	server "deployes/internal/domain/server"
	"deployes/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type TerminalHandler struct {
	serverRepo     server.Repository
	encKey         string
	jwtSecret      string
	allowedOrigins []string
}

func NewTerminalHandler(serverRepo server.Repository, encKey, jwtSecret string, allowedOrigins []string) *TerminalHandler {
	return &TerminalHandler{
		serverRepo:     serverRepo,
		encKey:         encKey,
		jwtSecret:      jwtSecret,
		allowedOrigins: allowedOrigins,
	}
}

type TerminalMessage struct {
	Type string `json:"type"` // "input", "resize"
	Data string `json:"data"`
	Rows int    `json:"rows,omitempty"`
	Cols int    `json:"cols,omitempty"`
}

func (h *TerminalHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Auth via Query Param
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	token, err := utils.ValidateToken(tokenStr, h.jwtSecret)
	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// 2. Server ID
	serverID := r.URL.Query().Get("serverId")
	if serverID == "" {
		http.Error(w, "server id required", http.StatusBadRequest)
		return
	}

	// 3. Fetch Server
	srv, err := h.serverRepo.FindByID(serverID)
	if err != nil {
		http.Error(w, "server not found", http.StatusNotFound)
		return
	}

	// Authorization check
	if srv.UserID != userID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	// 4. Upgrade
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return false
			}
			for _, allowed := range h.allowedOrigins {
				if allowed == origin {
					return true
				}
			}
			return false
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // upgrade error
	}
	defer ws.Close()

	// 5. Decrypt Key
	key, err := utils.Decrypt(srv.SSHKeyEncrypted, h.encKey)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Decryption failed"))
		return
	}

	// 6. Connect SSH
	client, err := utils.CreateSSHClient(srv.Host, srv.Port, srv.Username, key)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: SSH Connection failed: "+err.Error()))
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: New Session failed"))
		return
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 24, 80, modes); err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: PTY Request failed"))
		return
	}

	if err := session.Shell(); err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Shell failed"))
		return
	}

	// Pump stdout
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				return
			}
			ws.WriteMessage(websocket.BinaryMessage, buf[:n])
		}
	}()

	// Pump stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				return
			}
			ws.WriteMessage(websocket.BinaryMessage, buf[:n])
		}
	}()

	// Pump inputs
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}

		// Try JSON
		var tm TerminalMessage
		if err := json.Unmarshal(msg, &tm); err == nil && tm.Type != "" {
			if tm.Type == "resize" {
				session.WindowChange(tm.Rows, tm.Cols)
			} else if tm.Type == "input" {
				stdin.Write([]byte(tm.Data))
			}
		} else {
			// Raw input
			stdin.Write(msg)
		}
	}
}
