package deployment

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	domain "deployes/internal/domain/deployment"
	projectDomain "deployes/internal/domain/project"
	serverDomain "deployes/internal/domain/server"
)

type Service struct {
	repo        domain.Repository
	projectRepo projectDomain.Repository
	serverRepo  serverDomain.Repository
}

func NewService(repo domain.Repository, projectRepo projectDomain.Repository, serverRepo serverDomain.Repository) *Service {
	return &Service{
		repo:        repo,
		projectRepo: projectRepo,
		serverRepo:  serverRepo,
	}
}

func (s *Service) Create(userID string, req CreateDeploymentRequest) (*DeploymentResponse, error) {

	if req.ProjectID == "" || req.ServerID == "" {
		return nil, errors.New("projectId and serverId are required")
	}

	// VALIDALATION: Check Project Ownership
	project, err := s.projectRepo.FindByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if project.UserID != userID {
		return nil, errors.New("unauthorized: project access denied")
	}

	// VALIDATION: Check Server Ownership
	server, err := s.serverRepo.FindByID(req.ServerID)
	if err != nil {
		return nil, errors.New("server not found")
	}
	if server.UserID != userID {
		return nil, errors.New("unauthorized: server access denied")
	}

	d := &domain.Deployment{
		ID:        uuid.NewString(),
		UserID:    userID,
		ProjectID: req.ProjectID,
		ServerID:  req.ServerID,
		Status:    domain.StatusQueued,
		Logs:      "queued...\n",
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(d); err != nil {
		return nil, err
	}

	return &DeploymentResponse{
		ID:         d.ID,
		Status:     d.Status,
		Logs:       d.Logs,
		ProjectID:  d.ProjectID,
		ServerID:   d.ServerID,
		CommitHash: d.CommitHash,
	}, nil
}

func (s *Service) GetByID(id string) (*DeploymentResponse, error) {
	d, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &DeploymentResponse{
		ID:         d.ID,
		Status:     d.Status,
		Logs:       d.Logs,
		ProjectID:  d.ProjectID,
		ServerID:   d.ServerID,
		CommitHash: d.CommitHash,
	}, nil
}

func (s *Service) List(userID string) ([]*DeploymentResponse, error) {

	list, err := s.repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []*DeploymentResponse
	for _, d := range list {
		result = append(result, &DeploymentResponse{
			ID:         d.ID,
			Status:     d.Status,
			Logs:       d.Logs,
			ProjectID:  d.ProjectID,
			ServerID:   d.ServerID,
			CommitHash: d.CommitHash,
		})
	}

	return result, nil
}

func (s *Service) Rollback(userID string, originalID string) (*DeploymentResponse, error) {
	original, err := s.repo.FindByID(originalID)
	if err != nil {
		return nil, err
	}

	if original.CommitHash == "" {
		return nil, errors.New("cannot rollback: original deployment has no commit hash")
	}

	// Check ownership
	// Note: strict check might be needed, assuming userID matches
	if original.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	d := &domain.Deployment{
		ID:             uuid.NewString(),
		UserID:         userID,
		ProjectID:      original.ProjectID,
		ServerID:       original.ServerID,
		Status:         domain.StatusQueued,
		Logs:           fmt.Sprintf("queued rollback to %s...\n", original.CommitHash),
		CommitHash:     original.CommitHash,
		RollbackFromID: original.ID,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.Create(d); err != nil {
		return nil, err
	}

	return &DeploymentResponse{
		ID:         d.ID,
		Status:     d.Status,
		Logs:       d.Logs,
		ProjectID:  d.ProjectID,
		ServerID:   d.ServerID,
		CommitHash: d.CommitHash,
	}, nil
}

func (s *Service) GetStats(userID string) (*StatsResponse, error) {
	stats, err := s.repo.GetStats(userID)
	if err != nil {
		return nil, err
	}

	res := &StatsResponse{
		Total:           stats.Total,
		Successful:      stats.Successful,
		Failed:          stats.Failed,
		AverageDuration: stats.AverageDurationSeconds,
	}
	res.Last7Days.Dates = stats.Last7DaysDates
	res.Last7Days.Counts = stats.Last7DaysCounts

	return res, nil
}
