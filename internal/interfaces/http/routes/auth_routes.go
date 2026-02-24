package routes

import (
	"net/http"
	"time"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
	"deployes/internal/interfaces/http/middleware"
)

// Rate limiters for auth endpoints
var (
	loginRateLimiter    = middleware.NewRateLimiter(10, time.Minute) // 10 requests per minute
	registerRateLimiter = middleware.NewRateLimiter(5, time.Minute)  // 5 requests per minute
	refreshRateLimiter  = middleware.NewRateLimiter(10, time.Minute) // 10 requests per minute
)

func RegisterAuthRoutes(cfg *config.Config, userHandler *handlers.UserHandler) {
	corsConfig := middleware.NewCORSConfig(cfg.AllowedOrigins)

	// Register with rate limiting and security headers
	http.HandleFunc("/api/auth/register",
		middleware.SecurityHeaders(
			middleware.CORS(corsConfig,
				middleware.RateLimitMiddleware(registerRateLimiter, userHandler.Register))))

	// Login with rate limiting and security headers
	http.HandleFunc("/api/auth/login",
		middleware.SecurityHeaders(
			middleware.CORS(corsConfig,
				middleware.RateLimitMiddleware(loginRateLimiter, userHandler.Login))))

	// Refresh token endpoint
	http.HandleFunc("/api/auth/refresh",
		middleware.SecurityHeaders(
			middleware.CORS(corsConfig,
				middleware.RateLimitMiddleware(refreshRateLimiter, userHandler.Refresh))))

	// protected profile endpoint
	http.HandleFunc("/api/profile",
		middleware.SecurityHeaders(
			middleware.CORS(corsConfig,
				middleware.AuthMiddleware(cfg.JWTSecret, userHandler.Profile))))
}
