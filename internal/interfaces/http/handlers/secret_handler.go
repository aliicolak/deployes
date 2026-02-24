package handlers

import (
	"encoding/json"
	"net/http"

	app "deployes/internal/application/secret"
)

type SecretHandler struct {
	service *app.Service
}

func NewSecretHandler(service *app.Service) *SecretHandler {
	return &SecretHandler{service: service}
}

// POST /api/secrets
// Body: { "projectId": "...", "key": "...", "value": "..." }
// Note: CreatedRequest in service only has Key, Value. I need a wrapper or pass ProjectID separately.
// Service.Create takes projectID separately.
// Let's expect body: { "projectId": "...", "key": "...", "value": "..." }
type CreateSecretRequestBody struct {
	ProjectID string `json:"projectId"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (h *SecretHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body CreateSecretRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.ProjectID == "" {
		http.Error(w, "projectId is required", http.StatusBadRequest)
		return
	}

	req := app.CreateSecretRequest{
		Key:   body.Key,
		Value: body.Value,
	}

	res, err := h.service.Create(body.ProjectID, req)
	if err != nil {
		http.Error(w, "failed to create secret: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GET /api/secrets?projectId={projectId}
func (h *SecretHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		http.Error(w, "projectId query param is required", http.StatusBadRequest)
		return
	}

	res, err := h.service.List(projectID)
	if err != nil {
		http.Error(w, "failed to list secrets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// DELETE /api/secrets?id={id}
func (h *SecretHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query param is required", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "failed to delete secret", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
