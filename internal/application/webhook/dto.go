package webhook

type CreateWebhookRequest struct {
	ProjectID string   `json:"projectId"`
	ServerIDs []string `json:"serverIds"` // Changed from ServerID to support multiple servers
}

type WebhookResponse struct {
	ID         string   `json:"id"`
	ProjectID  string   `json:"projectId"`
	ServerIDs  []string `json:"serverIds"` // Changed from ServerID to support multiple servers
	Secret     string   `json:"secret,omitempty"`
	WebhookURL string   `json:"webhookUrl"`
	IsActive   bool     `json:"isActive"`
	CreatedAt  string   `json:"createdAt"`
}

type UpdateWebhookRequest struct {
	IsActive  *bool    `json:"isActive,omitempty"`
	ServerIDs []string `json:"serverIds,omitempty"` // Added to allow updating servers
}
