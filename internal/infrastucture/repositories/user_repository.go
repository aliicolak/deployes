package repositories

import (
	"database/sql"

	domain "deployes/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`,
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email=$1
	`, email)

	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id=$1
	`, id)

	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
