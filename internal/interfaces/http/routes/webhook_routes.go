package routes

import (
	"net/http"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"deployes/internal/interfaces/http/middleware"
)

func RegisterWebhookRoutes(cfg *config.Config, handler *handlers.WebhookHandler) {
	corsConfig := middleware.NewCORSConfig(cfg.AllowedOrigins)

	// Protected routes (require authentication)
	http.HandleFunc("/api/webhooks", middleware.SecurityHeaders(
		middleware.CORS(corsConfig,
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					middleware.AuthMiddleware(cfg.JWTSecret, handler.List)(w, r)
				case http.MethodPost:
					middleware.AuthMiddleware(cfg.JWTSecret, handler.Create)(w, r)
				case http.MethodPut:
					middleware.AuthMiddleware(cfg.JWTSecret, handler.Update)(w, r)
				case http.MethodDelete:
					middleware.AuthMiddleware(cfg.JWTSecret, handler.Delete)(w, r)
				case http.MethodOptions:
					w.WriteHeader(http.StatusNoContent)
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}
			})))

	// Public route - GitHub webhook endpoint (no auth, uses secret verification)
	http.HandleFunc("/api/webhooks/github/", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, handler.HandleGitHubWebhook)))
}
