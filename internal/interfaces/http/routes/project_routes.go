package routes

import (
	"net/http"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"deployes/internal/interfaces/http/middleware"
)

func RegisterProjectRoutes(cfg *config.Config, projectHandler *handlers.ProjectHandler) {
	corsConfig := middleware.NewCORSConfig(cfg.AllowedOrigins)

	http.HandleFunc("/api/projects", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.AuthMiddleware(cfg.JWTSecret, projectHandler.Create)(w, r)
			case http.MethodPut:
				middleware.AuthMiddleware(cfg.JWTSecret, projectHandler.Update)(w, r)
			case http.MethodGet:
				middleware.AuthMiddleware(cfg.JWTSecret, projectHandler.List)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	// Test repository access endpoint
	http.HandleFunc("/api/projects/test-access", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.AuthMiddleware(cfg.JWTSecret, projectHandler.TestRepoAccess)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	// Upload local project endpoint
	http.HandleFunc("/api/projects/upload", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.AuthMiddleware(cfg.JWTSecret, projectHandler.UploadLocalProject)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	// Delete project endpoint
	http.HandleFunc("/api/projects/delete", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				middleware.AuthMiddleware(cfg.JWTSecret, projectHandler.DeleteProject)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))
}
