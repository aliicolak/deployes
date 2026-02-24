package routes

import (
	"net/http"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"deployes/internal/interfaces/http/middleware"
)

func RegisterSecretRoutes(cfg *config.Config, secretHandler *handlers.SecretHandler) {
	corsConfig := middleware.NewCORSConfig(cfg.AllowedOrigins)

	http.HandleFunc("/api/secrets", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.AuthMiddleware(cfg.JWTSecret, secretHandler.Create)(w, r)
			case http.MethodGet:
				middleware.AuthMiddleware(cfg.JWTSecret, secretHandler.List)(w, r)
			case http.MethodDelete:
				middleware.AuthMiddleware(cfg.JWTSecret, secretHandler.Delete)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))
}
