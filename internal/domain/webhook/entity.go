package webhook

import "time"

type Webhook struct {
	ID        string
	UserID    string
	ProjectID string
	ServerIDs []string // Changed from ServerID string to support multiple servers
	Secret    string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

