package webhook

type Repository interface {
	Create(webhook *Webhook) error
	FindByID(id string) (*Webhook, error)
	FindBySecret(secret string) (*Webhook, error)
	FindByProjectID(projectID string) (*Webhook, error)
	ListByUserID(userID string) ([]*Webhook, error)
	Update(webhook *Webhook) error
	Delete(id string) error
}
