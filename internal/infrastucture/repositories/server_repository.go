package repositories

import (
	"database/sql"

	domain "deployes/internal/domain/server"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(db *sql.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) Create(server *domain.Server) error {
	_, err := r.db.Exec(`
		INSERT INTO servers (id, user_id, name, host, port, username, ssh_key_encrypted, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		server.ID,
		server.UserID,
		server.Name,
		server.Host,
		server.Port,
		server.Username,
		server.SSHKeyEncrypted,
		server.CreatedAt,
		server.UpdatedAt,
	)

	return err
}

func (r *ServerRepository) Update(server *domain.Server) error {
	_, err := r.db.Exec(`
		UPDATE servers 
		SET name=$1, host=$2, port=$3, username=$4, ssh_key_encrypted=$5, updated_at=$6
		WHERE id=$7 AND user_id=$8
	`,
		server.Name,
		server.Host,
		server.Port,
		server.Username,
		server.SSHKeyEncrypted,
		server.UpdatedAt,
		server.ID,
		server.UserID,
	)

	return err
}

func (r *ServerRepository) ListByUserID(userID string) ([]*domain.Server, error) {

	rows, err := r.db.Query(`
		SELECT id, user_id, name, host, port, username, ssh_key_encrypted, created_at, updated_at
		FROM servers
		WHERE user_id=$1
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []*domain.Server

	for rows.Next() {
		s := &domain.Server{}
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Name,
			&s.Host,
			&s.Port,
			&s.Username,
			&s.SSHKeyEncrypted,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		servers = append(servers, s)
	}

	return servers, nil
}

func (r *ServerRepository) FindByID(id string) (*domain.Server, error) {

	row := r.db.QueryRow(`
		SELECT id, user_id, name, host, port, username, ssh_key_encrypted, created_at, updated_at
		FROM servers
		WHERE id=$1
	`, id)

	s := &domain.Server{}
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.Name,
		&s.Host,
		&s.Port,
		&s.Username,
		&s.SSHKeyEncrypted,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return s, nil
}
