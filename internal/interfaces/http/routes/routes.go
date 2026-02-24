package routes

import (
	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RegisterAllRoutes(
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	projectHandler *handlers.ProjectHandler,
	serverHandler *handlers.ServerHandler,
	deploymentHandler *handlers.DeploymentHandler,
	webhookHandler *handlers.WebhookHandler,
	websocketHandler *handlers.WebSocketHandler,
	encryptionHandler *handlers.EncryptionHandler,
	secretHandler *handlers.SecretHandler,
) {
	RegisterAuthRoutes(cfg, userHandler)
	RegisterProjectRoutes(cfg, projectHandler)
	RegisterServerRoutes(cfg, serverHandler)
	RegisterDeploymentRoutes(cfg, deploymentHandler)
	RegisterWebhookRoutes(cfg, webhookHandler)
	RegisterWebSocketRoutes(websocketHandler)
	RegisterEncryptionRoutes(cfg, encryptionHandler)
	RegisterSecretRoutes(cfg, secretHandler)

	// Static file serving for Frontend
	staticPath := "./static"
	fs := http.FileServer(http.Dir(staticPath))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the path starts with /api/ or /ws, it's handled by other routes
		// Note: net/http matches the most specific pattern, but "/" is a catch-all
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/ws" || r.URL.Path == "/health" || r.URL.Path == "/terminal" {
			// These are handled by specific route registrations
			return
		}

		// Check if the file exists
		path := filepath.Join(staticPath, r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) || r.URL.Path == "/" {
			// Fallback to index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
			return
		}

		// Serve the static file
		fs.ServeHTTP(w, r)
	})
}
