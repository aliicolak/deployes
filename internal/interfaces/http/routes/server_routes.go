package routes

import (
	"net/http"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"deployes/internal/interfaces/http/middleware"
)

func RegisterServerRoutes(cfg *config.Config, serverHandler *handlers.ServerHandler) {
	corsConfig := middleware.NewCORSConfig(cfg.AllowedOrigins)

	http.HandleFunc("/api/servers", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.AuthMiddleware(cfg.JWTSecret, serverHandler.Create)(w, r)
			case http.MethodPut:
				middleware.AuthMiddleware(cfg.JWTSecret, serverHandler.Update)(w, r)
			case http.MethodGet:
				middleware.AuthMiddleware(cfg.JWTSecret, serverHandler.List)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	// Test connection endpoint
	http.HandleFunc("/api/servers/test-connection", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.AuthMiddleware(cfg.JWTSecret, serverHandler.TestConnection)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))
}
