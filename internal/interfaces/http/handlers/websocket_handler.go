package handlers

import (
	"log"
	"net/http"
	"strings"

	"deployes/internal/infrastucture/workers"
	"deployes/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	jwtSecret      string
	broadcaster    *workers.LogBroadcaster
	allowedOrigins []string
}

func NewWebSocketHandler(jwtSecret string, allowedOrigins []string) *WebSocketHandler {
	return &WebSocketHandler{
		jwtSecret:      jwtSecret,
		broadcaster:    workers.GetBroadcaster(),
		allowedOrigins: allowedOrigins,
	}
}

// createUpgrader creates a WebSocket upgrader with proper origin checking
func (h *WebSocketHandler) createUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
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
			log.Printf("⚠️  WebSocket connection rejected: origin %s not allowed", origin)
			return false
		},
	}
}

// StreamDeploymentLogs handles WebSocket connections for live log streaming
// GET /api/deployments/{id}/logs/stream?token=xxx
func (h *WebSocketHandler) StreamDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	// Extract deployment ID from URL path
	// Expected: /api/deployments/{id}/logs/stream
	pathParts := strings.Split(r.URL.Path, "/")
	var deploymentID string
	for i, part := range pathParts {
		if part == "deployments" && i+1 < len(pathParts) {
			deploymentID = pathParts[i+1]
			break
		}
	}

	if deploymentID == "" {
		http.Error(w, "deployment ID is required", http.StatusBadRequest)
		return
	}

	// Validate JWT token from query parameter (WebSocket doesn't support headers)
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "token query parameter is required", http.StatusUnauthorized)
		return
	}

	token, err := utils.ValidateToken(tokenStr, h.jwtSecret)
	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket with proper origin checking
	upgrader := h.createUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket client connected for deployment: %s", deploymentID)

	// Subscribe to log broadcasts
	logChan := h.broadcaster.Subscribe(deploymentID)
	defer h.broadcaster.Unsubscribe(deploymentID, logChan)

	// Handle incoming messages (for connection management)
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Client disconnected
				return
			}
		}
	}()

	// Stream logs to client
	for msg := range logChan {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}
