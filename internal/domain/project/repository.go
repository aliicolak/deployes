package project

type Repository interface {
	Create(project *Project) error
	Update(project *Project) error
	ListByUserID(userID string) ([]*Project, error)
	FindByID(id string) (*Project, error)
	Delete(id string) error
}
