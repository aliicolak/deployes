package handlers

import (
	"encoding/json"
	"net/http"

	"deployes/internal/config"
)

type EncryptionHandler struct {
	cfg *config.Config
}

func NewEncryptionHandler(cfg *config.Config) *EncryptionHandler {
	return &EncryptionHandler{cfg: cfg}
}

type EncryptionStatusResponse struct {
	Active    bool   `json:"active"`
	Algorithm string `json:"algorithm"`
	KeyLength int    `json:"keyLength"`
}

// GET /api/encryption/status
func (h *EncryptionHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if encryption key is properly configured
	keyLength := len(h.cfg.EncryptionKey)
	isActive := keyLength == 32 // AES-256 requires 32-byte key

	response := EncryptionStatusResponse{
		Active:    isActive,
		Algorithm: "AES-256-GCM",
		KeyLength: keyLength,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
