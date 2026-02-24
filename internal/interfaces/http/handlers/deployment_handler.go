package handlers

import (
	"encoding/json"
	"net/http"

	app "deployes/internal/application/deployment"
	"deployes/internal/interfaces/http/middleware"
)

type DeploymentHandler struct {
	service *app.Service
}

func NewDeploymentHandler(service *app.Service) *DeploymentHandler {
	return &DeploymentHandler{service: service}
}

// POST /api/deployments
func (h *DeploymentHandler) Create(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(middleware.UserIDKey).(string)

	var req app.CreateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(userId, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GET /api/deployments?id=xxx
func (h *DeploymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	res, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "deployment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GET /api/deployments/list
func (h *DeploymentHandler) List(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(middleware.UserIDKey).(string)

	res, err := h.service.List(userId)
	if err != nil {
		http.Error(w, "failed to fetch deployments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// POST /api/deployments/rollback?id={id}
func (h *DeploymentHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	res, err := h.service.Rollback(userId, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GET /api/dashboard/stats
func (h *DeploymentHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)

	res, err := h.service.GetStats(userId)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
