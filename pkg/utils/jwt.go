package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// GenerateToken creates a new JWT token for the given user ID
func GenerateToken(userID string, secret string) (string, error) {
	return GenerateTokenWithExpiry(userID, secret, time.Hour*24, AccessToken)
}

// GenerateRefreshToken creates a new refresh token with 7-day expiry
func GenerateRefreshToken(userID string, secret string) (string, error) {
	return GenerateTokenWithExpiry(userID, secret, time.Hour*24*7, RefreshToken)
}

// GenerateTokenWithExpiry creates a JWT token with custom expiry duration
func GenerateTokenWithExpiry(userID string, secret string, expiry time.Duration, tokenType TokenType) (string, error) {
	if len(secret) < 32 {
		return "", errors.New("JWT secret must be at least 32 characters")
	}

	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(expiry).Unix(),
		"iat":    time.Now().Unix(),
		"type":   string(tokenType),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateTokenPair creates both access and refresh tokens
func GenerateTokenPair(userID string, secret string) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateToken(userID, secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateRefreshToken(userID, secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns the parsed token
// It verifies the signing method to prevent "alg: none" attacks
func ValidateToken(tokenStr string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method to prevent "alg: none" attack
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})
}

// ValidateRefreshToken validates a refresh token and ensures it's the correct type
func ValidateRefreshToken(tokenStr string, secret string) (string, error) {
	token, err := ValidateToken(tokenStr, secret)
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != string(RefreshToken) {
		return "", errors.New("token is not a refresh token")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("userId not found in token")
	}

	return userId, nil
}
