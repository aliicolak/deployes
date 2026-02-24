package deployment

type Repository interface {
	Create(d *Deployment) error
	FindByID(id string) (*Deployment, error)

	GetNextQueued() (*Deployment, error)
	AppendLog(id string, logLine string) error
	MarkRunning(id string) error
	MarkFinished(id string, status string) error
	ListByUserID(userID string) ([]*Deployment, error)
	GetStats(userID string) (*Stats, error)
}
