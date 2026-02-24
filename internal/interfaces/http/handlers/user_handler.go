package handlers

import (
	app "deployes/internal/application/user"
	"deployes/internal/interfaces/http/middleware"
	"deployes/pkg/utils"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service *app.Service
}

func NewUserHandler(service *app.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Register -> POST /api/auth/register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	// Sadece POST isteğine izin verelim
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req app.RegisterRequest

	// JSON body parse ediyoruz
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Email validation
	if err := utils.ValidateEmail(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Password validation
	if err := utils.ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Service çağırıyoruz
	user, err := h.service.Register(req)
	if err != nil {
		// Don't leak internal error details
		http.Error(w, "registration failed", http.StatusBadRequest)
		return
	}

	// Response dön
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}

// Login -> POST /api/auth/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req app.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	res, err := h.service.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"protected route ✅", "userId":"` + userId + `"}`))
}

// Refresh -> POST /api/auth/refresh
func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req app.RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	res, err := h.service.RefreshToken(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
