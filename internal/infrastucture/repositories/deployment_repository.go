package repositories

import (
	"database/sql"
	"time"

	domain "deployes/internal/domain/deployment"
)

type DeploymentRepository struct {
	db *sql.DB
}

func NewDeploymentRepository(db *sql.DB) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

func (r *DeploymentRepository) Create(d *domain.Deployment) error {
	_, err := r.db.Exec(`
		INSERT INTO deployments (id, user_id, project_id, server_id, status, logs, created_at, commit_hash, rollback_from_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		d.ID, d.UserID, d.ProjectID, d.ServerID, d.Status, d.Logs, d.CreatedAt, d.CommitHash, d.RollbackFromID,
	)
	return err
}

func (r *DeploymentRepository) AppendLog(id string, logLine string) error {
	_, err := r.db.Exec(`
		UPDATE deployments
		SET logs = logs || $2
		WHERE id=$1
	`, id, logLine)
	return err
}

func (r *DeploymentRepository) FindByID(id string) (*domain.Deployment, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, project_id, server_id, status, logs, created_at, started_at, finished_at, COALESCE(commit_hash,''), COALESCE(rollback_from_id,'')
		FROM deployments
		WHERE id=$1
	`, id)

	d := &domain.Deployment{}
	err := row.Scan(
		&d.ID, &d.UserID, &d.ProjectID, &d.ServerID, &d.Status, &d.Logs,
		&d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.CommitHash, &d.RollbackFromID,
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// ✅ Worker için kritik method: aynı anda birden fazla worker olsa bile çakışmaz
func (r *DeploymentRepository) GetNextQueued() (*domain.Deployment, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// FOR UPDATE SKIP LOCKED:
	// Bir worker bir job'u aldıysa diğer worker o job’u görmez
	row := tx.QueryRow(`
		SELECT id, user_id, project_id, server_id, status, logs, created_at, started_at, finished_at, COALESCE(commit_hash,''), COALESCE(rollback_from_id,'')
		FROM deployments
		WHERE status='queued'
		ORDER BY created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`)

	d := &domain.Deployment{}
	err = row.Scan(
		&d.ID, &d.UserID, &d.ProjectID, &d.ServerID, &d.Status, &d.Logs,
		&d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.CommitHash, &d.RollbackFromID,
	)

	if err != nil {
		// hiç queued job yoksa
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// bu satır job’u seçtiğimiz anda "running" yapar, böylece queue’dan düşer.
	_, err = tx.Exec(`
		UPDATE deployments
		SET status='running', started_at=$2
		WHERE id=$1
	`, d.ID, time.Now())
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	d.Status = domain.StatusRunning
	return d, nil
}

func (r *DeploymentRepository) MarkRunning(id string) error {
	_, err := r.db.Exec(`
		UPDATE deployments
		SET status='running', started_at=$2
		WHERE id=$1
	`, id, time.Now())
	return err
}

func (r *DeploymentRepository) MarkFinished(id string, status string) error {
	_, err := r.db.Exec(`
		UPDATE deployments
		SET status=$2, finished_at=$3
		WHERE id=$1
	`, id, status, time.Now())
	return err
}

func (r *DeploymentRepository) ListByUserID(userID string) ([]*domain.Deployment, error) {

	rows, err := r.db.Query(`
		SELECT id, user_id, project_id, server_id, status, logs, created_at, started_at, finished_at, COALESCE(commit_hash,''), COALESCE(rollback_from_id,'')
		FROM deployments
		WHERE user_id=$1
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []*domain.Deployment

	for rows.Next() {
		d := &domain.Deployment{}
		err := rows.Scan(
			&d.ID, &d.UserID, &d.ProjectID, &d.ServerID, &d.Status, &d.Logs,
			&d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.CommitHash, &d.RollbackFromID,
		)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, d)
	}

	return deployments, nil
}

func (r *DeploymentRepository) UpdateCommitHash(id string, commitHash string) error {
	_, err := r.db.Exec(`UPDATE deployments SET commit_hash = $2 WHERE id = $1`, id, commitHash)
	return err
}

func (r *DeploymentRepository) GetStats(userID string) (*domain.Stats, error) {
	stats := &domain.Stats{}

	// 1. Aggregates
	row := r.db.QueryRow(`
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN status='success' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status='failed' THEN 1 ELSE 0 END), 0),
			COALESCE(AVG(EXTRACT(EPOCH FROM (finished_at - started_at))), 0)
		FROM deployments
		WHERE user_id=$1
	`, userID)

	if err := row.Scan(&stats.Total, &stats.Successful, &stats.Failed, &stats.AverageDurationSeconds); err != nil {
		return nil, err
	}

	// 2. Last 7 Days
	rows, err := r.db.Query(`
		SELECT to_char(d, 'YYYY-MM-DD'), count(deployments.id)
		FROM generate_series(CURRENT_DATE - INTERVAL '6 days', CURRENT_DATE, '1 day') d
		LEFT JOIN deployments ON to_char(deployments.created_at, 'YYYY-MM-DD') = to_char(d, 'YYYY-MM-DD') AND user_id=$1
		GROUP BY d
		ORDER BY d ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.Last7DaysCounts = []int{}
	stats.Last7DaysDates = []string{}

	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}
		stats.Last7DaysDates = append(stats.Last7DaysDates, date)
		stats.Last7DaysCounts = append(stats.Last7DaysCounts, count)
	}

	return stats, nil
}
