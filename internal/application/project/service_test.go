package project

import (
	"errors"
	"testing"

	domain "deployes/internal/domain/project"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepo adheres to domain.Repository
type MockProjectRepo struct {
	mock.Mock
}

func (m *MockProjectRepo) Create(p *domain.Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockProjectRepo) Update(p *domain.Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockProjectRepo) FindByID(id string) (*domain.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepo) ListByUserID(userID string) ([]*domain.Project, error) {
	args := m.Called(userID)
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateProject(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := NewService(mockRepo, "12345678901234567890123456789012")

	req := CreateProjectRequest{
		Name:         "Test Project",
		RepoURL:      "https://github.com/user/repo",
		Branch:       "main",
		DeployScript: "echo hello",
	}

	// Expect Create to be called once with any *domain.Project
	mockRepo.On("Create", mock.AnythingOfType("*project.Project")).Return(nil)

	res, err := service.Create("user-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.RepoURL, res.RepoURL)
	assert.NotEmpty(t, res.ID)

	mockRepo.AssertExpectations(t)
}

func TestFindByID_Success(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := NewService(mockRepo, "12345678901234567890123456789012")

	expected := &domain.Project{
		ID:      "p-1",
		Name:    "Existing Project",
		RepoURL: "http://git",
		Branch:  "dev",
	}

	mockRepo.On("FindByID", "p-1").Return(expected, nil)

	res, err := service.FindByID("p-1")

	assert.NoError(t, err)
	assert.Equal(t, "Existing Project", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestFindByID_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepo)
	service := NewService(mockRepo, "12345678901234567890123456789012")

	mockRepo.On("FindByID", "p-99").Return(nil, errors.New("not found"))

	res, err := service.FindByID("p-99")

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}
