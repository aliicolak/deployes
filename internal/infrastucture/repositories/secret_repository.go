package repositories

import (
	"database/sql"
	domain "deployes/internal/domain/secret"
)

type secretRepository struct {
	db *sql.DB
}

func NewSecretRepository(db *sql.DB) domain.Repository {
	return &secretRepository{db: db}
}

func (r *secretRepository) Save(s *domain.Secret) error {
	query := `
		INSERT INTO secrets (id, project_id, key, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (project_id, key) DO UPDATE 
		SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(query, s.ID, s.ProjectID, s.Key, s.Value, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *secretRepository) ListByProjectID(projectID string) ([]*domain.Secret, error) {
	query := `SELECT id, project_id, key, value, created_at, updated_at FROM secrets WHERE project_id = $1`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*domain.Secret
	for rows.Next() {
		s := &domain.Secret{}
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		secrets = append(secrets, s)
	}
	return secrets, nil
}

func (r *secretRepository) FindByID(id string) (*domain.Secret, error) {
	query := `SELECT id, project_id, key, value, created_at, updated_at FROM secrets WHERE id = $1`
	s := &domain.Secret{}
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.ProjectID, &s.Key, &s.Value, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *secretRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM secrets WHERE id = $1", id)
	return err
}

func (r *secretRepository) DeleteByProjectID(projectID string) error {
	_, err := r.db.Exec("DELETE FROM secrets WHERE project_id = $1", projectID)
	return err
}
