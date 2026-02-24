package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	app "deployes/internal/application/server"
	"deployes/internal/interfaces/http/middleware"
)

type ServerHandler struct {
	service *app.Service
}

func NewServerHandler(service *app.Service) *ServerHandler {
	return &ServerHandler{service: service}
}

// POST /api/servers
func (h *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(middleware.UserIDKey).(string)

	var req app.CreateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(userId, req)
	if err != nil {
		// Don't leak internal error details
		http.Error(w, "failed to create server", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// PUT /api/servers?id={id}
func (h *ServerHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var req app.CreateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Update(userId, id, req)
	if err != nil {
		// Don't leak internal error details
		http.Error(w, "failed to update server", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GET /api/servers
func (h *ServerHandler) List(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(middleware.UserIDKey).(string)

	res, err := h.service.List(userId)
	if err != nil {
		log.Printf("Failed to fetch servers: %v", err)
		http.Error(w, "failed to fetch servers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// POST /api/servers/test-connection
func (h *ServerHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)

	var req app.TestConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.TestConnection(userId, req)
	if err != nil {
		http.Error(w, "failed to test connection", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
