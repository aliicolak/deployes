package server

type Repository interface {
	Create(server *Server) error
	Update(server *Server) error
	ListByUserID(userID string) ([]*Server, error)
	FindByID(id string) (*Server, error)
}
