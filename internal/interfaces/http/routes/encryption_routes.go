package routes

import (
	"net/http"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"deployes/internal/interfaces/http/middleware"
)

func RegisterEncryptionRoutes(cfg *config.Config, encryptionHandler *handlers.EncryptionHandler) {
	corsConfig := middleware.NewCORSConfig(cfg.AllowedOrigins)

	http.HandleFunc("/api/encryption/status", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				middleware.AuthMiddleware(cfg.JWTSecret, encryptionHandler.Status)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))
}
