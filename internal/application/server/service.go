package server

import (
	"errors"
	"time"

	"github.com/google/uuid"

	domain "deployes/internal/domain/server"
	"deployes/pkg/utils"
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

func (s *Service) Create(userID string, req CreateServerRequest) (*ServerResponse, error) {

	if req.Name == "" || req.Host == "" || req.Username == "" || req.SSHKey == "" {
		return nil, errors.New("name, host, username and sshKey are required")
	}

	if req.Port == 0 {
		req.Port = 22
	}

	encryptedKey, err := utils.Encrypt(req.SSHKey, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	server := &domain.Server{
		ID:              uuid.NewString(),
		UserID:          userID,
		Name:            req.Name,
		Host:            req.Host,
		Port:            req.Port,
		Username:        req.Username,
		SSHKeyEncrypted: encryptedKey,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.repo.Create(server)
	if err != nil {
		return nil, err
	}

	return &ServerResponse{
		ID:       server.ID,
		Name:     server.Name,
		Host:     server.Host,
		Port:     server.Port,
		Username: server.Username,
	}, nil
}

func (s *Service) List(userID string) ([]*ServerResponse, error) {
	servers, err := s.repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []*ServerResponse
	for _, srv := range servers {
		result = append(result, &ServerResponse{
			ID:       srv.ID,
			Name:     srv.Name,
			Host:     srv.Host,
			Port:     srv.Port,
			Username: srv.Username,
		})
	}

	return result, nil
}

func (s *Service) Update(userID string, id string, req CreateServerRequest) (*ServerResponse, error) {
	// Find server
	server, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("server not found")
	}

	// Verify ownership
	if server.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields
	server.Name = req.Name
	server.Host = req.Host
	server.Port = req.Port
	server.Username = req.Username
	server.UpdatedAt = time.Now()

	// Only update SSH key if provided
	if req.SSHKey != "" {
		encryptedKey, err := utils.Encrypt(req.SSHKey, s.encryptionKey)
		if err != nil {
			return nil, err
		}
		server.SSHKeyEncrypted = encryptedKey
	}

	err = s.repo.Update(server)
	if err != nil {
		return nil, err
	}

	return &ServerResponse{
		ID:       server.ID,
		Name:     server.Name,
		Host:     server.Host,
		Port:     server.Port,
		Username: server.Username,
	}, nil
}

// TestConnection tests SSH connectivity to a server
func (s *Service) TestConnection(userID string, req TestConnectionRequest) (*ConnectionTestResult, error) {
	var host, username, sshKey string
	var port int

	// If ServerID is provided, use existing server credentials
	if req.ServerID != "" {
		server, err := s.repo.FindByID(req.ServerID)
		if err != nil {
			return &ConnectionTestResult{
				Success: false,
				Message: "Server not found",
			}, nil
		}

		// Verify ownership
		if server.UserID != userID {
			return &ConnectionTestResult{
				Success: false,
				Message: "Unauthorized",
			}, nil
		}

		host = server.Host
		port = server.Port
		username = server.Username

		// Decrypt SSH key
		decryptedKey, err := utils.Decrypt(server.SSHKeyEncrypted, s.encryptionKey)
		if err != nil {
			return &ConnectionTestResult{
				Success: false,
				Message: "Failed to decrypt SSH key",
			}, nil
		}
		sshKey = decryptedKey
	} else {
		// Use provided credentials
		if req.Host == "" || req.Username == "" || req.SSHKey == "" {
			return &ConnectionTestResult{
				Success: false,
				Message: "Host, username, and SSH key are required",
			}, nil
		}
		host = req.Host
		port = req.Port
		username = req.Username
		sshKey = req.SSHKey

		if port == 0 {
			port = 22
		}
	}

	// Measure connection time
	startTime := time.Now()

	// Try to establish SSH connection
	client, err := utils.CreateSSHClient(host, port, username, sshKey)
	if err != nil {
		return &ConnectionTestResult{
			Success: false,
			Message: "Connection failed: " + err.Error(),
			Latency: time.Since(startTime).Milliseconds(),
		}, nil
	}
	defer client.Close()

	latency := time.Since(startTime).Milliseconds()

	// Run a simple command to verify connection works
	output, err := utils.RunSSHCommand(client, "echo 'Connection test successful'")
	if err != nil {
		return &ConnectionTestResult{
			Success: false,
			Message: "Connection established but command execution failed: " + err.Error(),
			Latency: latency,
		}, nil
	}

	return &ConnectionTestResult{
		Success: true,
		Message: "Connection successful: " + output,
		Latency: latency,
	}, nil
}
