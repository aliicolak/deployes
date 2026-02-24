package secret

import (
	domain "deployes/internal/domain/secret"
	"deployes/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo          domain.Repository
	encryptionKey string
}

func NewService(repo domain.Repository, encryptionKey string) *Service {
	return &Service{
		repo:          repo,
		encryptionKey: encryptionKey,
	}
}

type CreateSecretRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SecretResponse struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt"`
}

func (s *Service) Create(projectID string, req CreateSecretRequest) (*SecretResponse, error) {
	if req.Key == "" || req.Value == "" {
		return nil, errors.New("key and value are required")
	}

	encryptedVal, err := utils.Encrypt(req.Value, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	secret := &domain.Secret{
		ID:        uuid.NewString(),
		ProjectID: projectID,
		Key:       req.Key,
		Value:     encryptedVal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Save(secret); err != nil {
		return nil, err
	}

	return &SecretResponse{
		ID:        secret.ID,
		Key:       secret.Key,
		CreatedAt: secret.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) List(projectID string) ([]*SecretResponse, error) {
	secrets, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	var res []*SecretResponse
	for _, sc := range secrets {
		res = append(res, &SecretResponse{
			ID:        sc.ID,
			Key:       sc.Key,
			CreatedAt: sc.CreatedAt.Format(time.RFC3339),
		})
	}
	return res, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

// Internal use only (for deployment worker)
func (s *Service) GetDecryptedSecrets(projectID string) (map[string]string, error) {
	secrets, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, sc := range secrets {
		decrypted, err := utils.Decrypt(sc.Value, s.encryptionKey)
		if err != nil {
			continue // Skip if decryption fails? Or error out? Better skip or log.
		}
		result[sc.Key] = decrypted
	}
	return result, nil
}
