package routes

import (
	"net/http"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"deployes/internal/interfaces/http/middleware"
)

func RegisterDeploymentRoutes(cfg *config.Config, deploymentHandler *handlers.DeploymentHandler) {
	corsConfig := middleware.NewCORSConfig(cfg.AllowedOrigins)

	http.HandleFunc("/api/deployments", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.AuthMiddleware(cfg.JWTSecret, deploymentHandler.Create)(w, r)
			case http.MethodGet:
				middleware.AuthMiddleware(cfg.JWTSecret, deploymentHandler.GetByID)(w, r)
			case http.MethodOptions:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	http.HandleFunc("/api/deployments/list", middleware.SecurityHeaders(
		middleware.CORS(corsConfig,
			middleware.AuthMiddleware(cfg.JWTSecret, deploymentHandler.List))))
	http.HandleFunc("/api/deployments/rollback", middleware.SecurityHeaders(
		middleware.CORS(corsConfig, func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost {
				middleware.AuthMiddleware(cfg.JWTSecret, deploymentHandler.Rollback)(w, r)
			} else if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
			} else {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	http.HandleFunc("/api/dashboard/stats", middleware.SecurityHeaders(
		middleware.CORS(corsConfig,
			middleware.AuthMiddleware(cfg.JWTSecret, deploymentHandler.GetStats))))
}
