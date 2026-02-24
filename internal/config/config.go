package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppPort        string
	DBUrl          string
	JWTSecret      string
	EncryptionKey  string
	BaseURL        string
	AllowedOrigins []string
}

// Load loads configuration from environment variables.
// Critical security settings (DATABASE_URL, JWT_SECRET, ENCRYPTION_KEY) are mandatory.
func Load() *Config {
	appPort := getEnv("APP_PORT", "8080")

	// Mandatory environment variables - no defaults for security
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("❌ FATAL: DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ FATAL: JWT_SECRET environment variable is required")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("❌ FATAL: JWT_SECRET must be at least 32 characters long")
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		log.Fatal("❌ FATAL: ENCRYPTION_KEY environment variable is required")
	}
	if len(encryptionKey) != 32 {
		log.Fatal("❌ FATAL: ENCRYPTION_KEY must be exactly 32 characters long")
	}

	// Allowed origins for CORS - defaults to localhost for development
	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:3000")
	allowedOrigins := parseCSV(allowedOriginsStr)

	cfg := &Config{
		AppPort:        appPort,
		DBUrl:          dbUrl,
		JWTSecret:      jwtSecret,
		EncryptionKey:  encryptionKey,
		BaseURL:        getEnv("BASE_URL", fmt.Sprintf("http://localhost:%s", appPort)),
		AllowedOrigins: allowedOrigins,
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// parseCSV parses a comma-separated string into a slice
func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	for _, v := range splitAndTrim(s, ",") {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func splitAndTrim(s string, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			trimmed := trim(s[start:i])
			if trimmed != "" {
				result = append(result, trimmed)
			}
			start = i + 1
		}
	}
	trimmed := trim(s[start:])
	if trimmed != "" {
		result = append(result, trimmed)
	}
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
