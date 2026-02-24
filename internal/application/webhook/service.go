package webhook

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	domain "deployes/internal/domain/webhook"
)

type Service struct {
	repo    domain.Repository
	baseURL string
}

func NewService(repo domain.Repository, baseURL string) *Service {
	return &Service{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (s *Service) Create(userID string, req CreateWebhookRequest) (*WebhookResponse, error) {
	if req.ProjectID == "" || len(req.ServerIDs) == 0 {
		return nil, errors.New("projectId and at least one serverId are required")
	}

	// Generate random secret
	secret, err := generateSecret(32)
	if err != nil {
		return nil, err
	}

	webhook := &domain.Webhook{
		ID:        uuid.NewString(),
		UserID:    userID,
		ProjectID: req.ProjectID,
		ServerIDs: req.ServerIDs,
		Secret:    secret,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(webhook); err != nil {
		return nil, err
	}

	return &WebhookResponse{
		ID:         webhook.ID,
		ProjectID:  webhook.ProjectID,
		ServerIDs:  webhook.ServerIDs,
		Secret:     webhook.Secret, // Only shown on creation
		WebhookURL: fmt.Sprintf("%s/api/webhooks/github/%s", s.baseURL, webhook.ID),
		IsActive:   webhook.IsActive,
		CreatedAt:  webhook.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) GetByID(id string) (*WebhookResponse, error) {
	webhook, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &WebhookResponse{
		ID:         webhook.ID,
		ProjectID:  webhook.ProjectID,
		ServerIDs:  webhook.ServerIDs,
		WebhookURL: fmt.Sprintf("%s/api/webhooks/github/%s", s.baseURL, webhook.ID),
		IsActive:   webhook.IsActive,
		CreatedAt:  webhook.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) List(userID string) ([]*WebhookResponse, error) {
	webhooks, err := s.repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []*WebhookResponse
	for _, wh := range webhooks {
		result = append(result, &WebhookResponse{
			ID:         wh.ID,
			ProjectID:  wh.ProjectID,
			ServerIDs:  wh.ServerIDs,
			WebhookURL: fmt.Sprintf("%s/api/webhooks/github/%s", s.baseURL, wh.ID),
			IsActive:   wh.IsActive,
			CreatedAt:  wh.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *Service) Update(id, userID string, req UpdateWebhookRequest) (*WebhookResponse, error) {
	webhook, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if webhook.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.IsActive != nil {
		webhook.IsActive = *req.IsActive
	}
	if len(req.ServerIDs) > 0 {
		webhook.ServerIDs = req.ServerIDs
	}
	webhook.UpdatedAt = time.Now()

	if err := s.repo.Update(webhook); err != nil {
		return nil, err
	}

	return &WebhookResponse{
		ID:         webhook.ID,
		ProjectID:  webhook.ProjectID,
		ServerIDs:  webhook.ServerIDs,
		WebhookURL: fmt.Sprintf("%s/api/webhooks/github/%s", s.baseURL, webhook.ID),
		IsActive:   webhook.IsActive,
		CreatedAt:  webhook.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) Delete(id, userID string) error {
	webhook, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if webhook.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.Delete(id)
}

func (s *Service) FindByID(id string) (*domain.Webhook, error) {
	return s.repo.FindByID(id)
}

func generateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
