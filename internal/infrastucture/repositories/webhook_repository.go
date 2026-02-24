package repositories

import (
	"database/sql"

	"github.com/google/uuid"

	domain "deployes/internal/domain/webhook"
)

type webhookRepository struct {
	db *sql.DB
}

func NewWebhookRepository(db *sql.DB) domain.Repository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) Create(webhook *domain.Webhook) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert webhook (without server_id - it's now in junction table)
	query := `
		INSERT INTO webhooks (id, user_id, project_id, secret, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(query,
		webhook.ID,
		webhook.UserID,
		webhook.ProjectID,
		webhook.Secret,
		webhook.IsActive,
		webhook.CreatedAt,
		webhook.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert server associations into junction table
	for _, serverID := range webhook.ServerIDs {
		serverQuery := `
			INSERT INTO webhook_servers (id, webhook_id, server_id, created_at)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(serverQuery,
			uuid.NewString(),
			webhook.ID,
			serverID,
			webhook.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *webhookRepository) FindByID(id string) (*domain.Webhook, error) {
	// Get webhook basic info
	query := `
		SELECT id, user_id, project_id, secret, is_active, created_at, updated_at
		FROM webhooks WHERE id = $1
	`
	webhook := &domain.Webhook{}
	err := r.db.QueryRow(query, id).Scan(
		&webhook.ID,
		&webhook.UserID,
		&webhook.ProjectID,
		&webhook.Secret,
		&webhook.IsActive,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get server IDs from junction table
	serverIDs, err := r.getServerIDsForWebhook(webhook.ID)
	if err != nil {
		return nil, err
	}
	webhook.ServerIDs = serverIDs

	return webhook, nil
}

func (r *webhookRepository) FindBySecret(secret string) (*domain.Webhook, error) {
	query := `
		SELECT id, user_id, project_id, secret, is_active, created_at, updated_at
		FROM webhooks WHERE secret = $1
	`
	webhook := &domain.Webhook{}
	err := r.db.QueryRow(query, secret).Scan(
		&webhook.ID,
		&webhook.UserID,
		&webhook.ProjectID,
		&webhook.Secret,
		&webhook.IsActive,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get server IDs from junction table
	serverIDs, err := r.getServerIDsForWebhook(webhook.ID)
	if err != nil {
		return nil, err
	}
	webhook.ServerIDs = serverIDs

	return webhook, nil
}

func (r *webhookRepository) FindByProjectID(projectID string) (*domain.Webhook, error) {
	query := `
		SELECT id, user_id, project_id, secret, is_active, created_at, updated_at
		FROM webhooks WHERE project_id = $1 AND is_active = true
		LIMIT 1
	`
	webhook := &domain.Webhook{}
	err := r.db.QueryRow(query, projectID).Scan(
		&webhook.ID,
		&webhook.UserID,
		&webhook.ProjectID,
		&webhook.Secret,
		&webhook.IsActive,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get server IDs from junction table
	serverIDs, err := r.getServerIDsForWebhook(webhook.ID)
	if err != nil {
		return nil, err
	}
	webhook.ServerIDs = serverIDs

	return webhook, nil
}

func (r *webhookRepository) ListByUserID(userID string) ([]*domain.Webhook, error) {
	query := `
		SELECT id, user_id, project_id, secret, is_active, created_at, updated_at
		FROM webhooks WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []*domain.Webhook
	for rows.Next() {
		webhook := &domain.Webhook{}
		err := rows.Scan(
			&webhook.ID,
			&webhook.UserID,
			&webhook.ProjectID,
			&webhook.Secret,
			&webhook.IsActive,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get server IDs for each webhook
		serverIDs, err := r.getServerIDsForWebhook(webhook.ID)
		if err != nil {
			return nil, err
		}
		webhook.ServerIDs = serverIDs

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

func (r *webhookRepository) Update(webhook *domain.Webhook) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update webhook basic info
	query := `
		UPDATE webhooks SET is_active = $1, updated_at = $2 WHERE id = $3
	`
	_, err = tx.Exec(query, webhook.IsActive, webhook.UpdatedAt, webhook.ID)
	if err != nil {
		return err
	}

	// If ServerIDs are provided, update junction table
	if len(webhook.ServerIDs) > 0 {
		// Delete existing server associations
		deleteQuery := `DELETE FROM webhook_servers WHERE webhook_id = $1`
		_, err = tx.Exec(deleteQuery, webhook.ID)
		if err != nil {
			return err
		}

		// Insert new server associations
		for _, serverID := range webhook.ServerIDs {
			insertQuery := `
				INSERT INTO webhook_servers (id, webhook_id, server_id, created_at)
				VALUES ($1, $2, $3, $4)
			`
			_, err = tx.Exec(insertQuery,
				uuid.NewString(),
				webhook.ID,
				serverID,
				webhook.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *webhookRepository) Delete(id string) error {
	// Junction table entries will be deleted via CASCADE
	query := `DELETE FROM webhooks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Helper function to get server IDs for a webhook from junction table
func (r *webhookRepository) getServerIDsForWebhook(webhookID string) ([]string, error) {
	query := `SELECT server_id FROM webhook_servers WHERE webhook_id = $1`
	rows, err := r.db.Query(query, webhookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serverIDs []string
	for rows.Next() {
		var serverID string
		if err := rows.Scan(&serverID); err != nil {
			return nil, err
		}
		serverIDs = append(serverIDs, serverID)
	}

	return serverIDs, nil
}
