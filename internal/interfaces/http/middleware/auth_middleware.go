package middleware

import (
	"context"
	"net/http"
	"strings"

	"deployes/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userId"

func AuthMiddleware(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// "Bearer token"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		token, err := utils.ValidateToken(tokenStr, jwtSecret)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// ✅ claims artık jwt.MapClaims olmalı (v5)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userId, ok := claims["userId"].(string)
		if !ok {
			http.Error(w, "userId not found in token", http.StatusUnauthorized)
			return
		}

		// userId'yi request context içine koyuyoruz
		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		next(w, r.WithContext(ctx))
	}
}
