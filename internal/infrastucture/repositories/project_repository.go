package repositories

import (
	"database/sql"

	domain "deployes/internal/domain/project"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *domain.Project) error {
	_, err := r.db.Exec(`
		INSERT INTO projects (
			id, user_id, name, type, repo_url, branch, local_path, deploy_script, 
			include_patterns, exclude_patterns, preserve_patterns, 
			scm_private_key, scm_public_key,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`,
		project.ID,
		project.UserID,
		project.Name,
		project.Type,
		project.RepoURL,
		project.Branch,
		project.LocalPath,
		project.DeployScript,
		project.IncludePatterns,
		project.ExcludePatterns,
		project.PreservePatterns,
		project.SCMPrivateKeyEncrypted,
		project.SCMPublicKey,
		project.CreatedAt,
		project.UpdatedAt,
	)

	return err
}

func (r *ProjectRepository) Update(project *domain.Project) error {
	// Not updating Keys usually? Or yes?
	// If user regenerates keys, yes.
	_, err := r.db.Exec(`
		UPDATE projects 
		SET name=$1, type=$2, repo_url=$3, branch=$4, local_path=$5, deploy_script=$6, 
			include_patterns=$7, exclude_patterns=$8, preserve_patterns=$9,
			scm_private_key=$10, scm_public_key=$11,
			updated_at=$12
		WHERE id=$13 AND user_id=$14
	`,
		project.Name,
		project.Type,
		project.RepoURL,
		project.Branch,
		project.LocalPath,
		project.DeployScript,
		project.IncludePatterns,
		project.ExcludePatterns,
		project.PreservePatterns,
		project.SCMPrivateKeyEncrypted,
		project.SCMPublicKey,
		project.UpdatedAt,
		project.ID,
		project.UserID,
	)

	return err
}

func (r *ProjectRepository) ListByUserID(userID string) ([]*domain.Project, error) {
	rows, err := r.db.Query(`
		SELECT 
			id, user_id, name, COALESCE(type, 'github'), repo_url, branch, COALESCE(local_path, ''), deploy_script, 
			COALESCE(include_patterns,''), COALESCE(exclude_patterns,''), COALESCE(preserve_patterns,''), 
			COALESCE(scm_private_key,''), COALESCE(scm_public_key,''),
			created_at, updated_at
		FROM projects
		WHERE user_id=$1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project

	for rows.Next() {
		p := &domain.Project{}
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Type,
			&p.RepoURL,
			&p.Branch,
			&p.LocalPath,
			&p.DeployScript,
			&p.IncludePatterns,
			&p.ExcludePatterns,
			&p.PreservePatterns,
			&p.SCMPrivateKeyEncrypted,
			&p.SCMPublicKey,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (r *ProjectRepository) FindByID(id string) (*domain.Project, error) {
	row := r.db.QueryRow(`
		SELECT 
			id, user_id, name, COALESCE(type, 'github'), repo_url, branch, COALESCE(local_path, ''), deploy_script, 
			COALESCE(include_patterns,''), COALESCE(exclude_patterns,''), COALESCE(preserve_patterns,''), 
			COALESCE(scm_private_key,''), COALESCE(scm_public_key,''),
			created_at, updated_at
		FROM projects
		WHERE id=$1
	`, id)

	p := &domain.Project{}
	err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.Type,
		&p.RepoURL,
		&p.Branch,
		&p.LocalPath,
		&p.DeployScript,
		&p.IncludePatterns,
		&p.ExcludePatterns,
		&p.PreservePatterns,
		&p.SCMPrivateKeyEncrypted,
		&p.SCMPublicKey,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *ProjectRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id=$1`, id)
	return err
}
