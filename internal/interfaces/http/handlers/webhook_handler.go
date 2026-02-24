package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	appDeployment "deployes/internal/application/deployment"
	appProject "deployes/internal/application/project"
	app "deployes/internal/application/webhook"
	"deployes/internal/interfaces/http/middleware"

	"github.com/google/uuid"
)

type WebhookHandler struct {
	service           *app.Service
	deploymentService *appDeployment.Service
	projectService    *appProject.Service
}

func NewWebhookHandler(
	service *app.Service,
	deploymentService *appDeployment.Service,
	projectService *appProject.Service,
) *WebhookHandler {
	return &WebhookHandler{
		service:           service,
		deploymentService: deploymentService,
		projectService:    projectService,
	}
}

// POST /api/webhooks
func (h *WebhookHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req app.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GET /api/webhooks
func (h *WebhookHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	res, err := h.service.List(userID)
	if err != nil {
		http.Error(w, "failed to fetch webhooks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// PUT /api/webhooks?id=xxx
func (h *WebhookHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var req app.UpdateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Update(id, userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// DELETE /api/webhooks?id=xxx
func (h *WebhookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GitHubPushPayload represents the GitHub webhook push event payload
type GitHubPushPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		FullName string `json:"full_name"`
		CloneURL string `json:"clone_url"`
		SSHURL   string `json:"ssh_url"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	HeadCommit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
}

// POST /api/webhooks/github/{webhookId}
func (h *WebhookHandler) HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	// Extract webhook ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	var webhookID string
	for i, part := range pathParts {
		if part == "github" && i+1 < len(pathParts) {
			webhookID = pathParts[i+1]
			break
		}
	}

	if webhookID == "" {
		http.Error(w, "webhook ID is required", http.StatusBadRequest)
		return
	}

	// Get webhook from database
	webhook, err := h.service.FindByID(webhookID)
	if err != nil {
		http.Error(w, "webhook not found", http.StatusNotFound)
		return
	}

	if !webhook.IsActive {
		http.Error(w, "webhook is inactive", http.StatusForbidden)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify GitHub signature - MANDATORY for security
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		log.Printf("⚠️  Webhook request rejected: missing signature for webhook %s", webhookID)
		http.Error(w, "missing signature header", http.StatusUnauthorized)
		return
	}
	if !verifyGitHubSignature(body, signature, webhook.Secret) {
		log.Printf("⚠️  Webhook request rejected: invalid signature for webhook %s", webhookID)
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse payload
	var payload GitHubPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Check event type
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "push" {
		// We only handle push events, but acknowledge other events
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "event type not handled",
			"event":   eventType,
		})
		return
	}

	log.Printf("📥 GitHub webhook received for project: %s, ref: %s",
		payload.Repository.FullName, payload.Ref)

	// Get project to check branch
	project, err := h.projectService.FindByID(webhook.ProjectID)
	if err != nil {
		log.Printf("❌ Failed to get project for webhook: %v", err)
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	// Extract branch name from ref (refs/heads/main -> main)
	pushedBranch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	// Check if pushed branch matches project's configured branch
	if pushedBranch != project.Branch {
		log.Printf("⏭️  Branch mismatch: pushed=%s, configured=%s. Skipping deployment.",
			pushedBranch, project.Branch)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":          "branch mismatch, deployment skipped",
			"pushedBranch":     pushedBranch,
			"configuredBranch": project.Branch,
		})
		return
	}

	log.Printf("✅ Branch matched: %s. Triggering deployment...", pushedBranch)

	// Create deployments for all associated servers
	var deploymentIDs []string
	for _, serverID := range webhook.ServerIDs {
		deployReq := appDeployment.CreateDeploymentRequest{
			ProjectID: webhook.ProjectID,
			ServerID:  serverID,
		}

		deployment, err := h.deploymentService.Create(webhook.UserID, deployReq)
		if err != nil {
			log.Printf("❌ Failed to create deployment for server %s: %v", serverID, err)
			continue // Continue with other servers even if one fails
		}

		deploymentIDs = append(deploymentIDs, deployment.ID)
		log.Printf("🚀 Deployment triggered for server %s: %s", serverID, deployment.ID)
	}

	if len(deploymentIDs) == 0 {
		log.Printf("❌ No deployments were created from webhook")
		http.Error(w, "failed to create any deployments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "deployments triggered",
		"deploymentIds": deploymentIDs,
		"serverCount":   len(deploymentIDs),
		"branch":        pushedBranch,
		"triggeredAt":   time.Now().Format(time.RFC3339),
	})
}

// verifyGitHubSignature verifies the HMAC-SHA256 signature from GitHub
func verifyGitHubSignature(payload []byte, signature, secret string) bool {
	// Signature format: sha256=xxxxx
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	expectedSig := strings.TrimPrefix(signature, "sha256=")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	actualSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedSig), []byte(actualSig))
}

// Generate unique request ID for tracking
func generateRequestID() string {
	return uuid.NewString()[:8]
}
