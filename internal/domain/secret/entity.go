package secret

import "time"

type Secret struct {
	ID        string
	ProjectID string
	Key       string
	Value     string // Stores encrypted value
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repository interface
type Repository interface {
	Save(secret *Secret) error
	ListByProjectID(projectID string) ([]*Secret, error)
	FindByID(id string) (*Secret, error)
	Delete(id string) error
	DeleteByProjectID(projectID string) error
}
