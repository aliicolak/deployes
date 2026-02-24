package deployment

import (
	"testing"

	domain "deployes/internal/domain/deployment"
	projectDomain "deployes/internal/domain/project"
	serverDomain "deployes/internal/domain/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDeploymentRepo struct {
	mock.Mock
}

func (m *MockDeploymentRepo) Create(d *domain.Deployment) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *MockDeploymentRepo) FindByID(id string) (*domain.Deployment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Deployment), args.Error(1)
}

func (m *MockDeploymentRepo) GetNextQueued() (*domain.Deployment, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Deployment), args.Error(1)
}

func (m *MockDeploymentRepo) AppendLog(id string, logLine string) error {
	args := m.Called(id, logLine)
	return args.Error(0)
}

func (m *MockDeploymentRepo) MarkRunning(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDeploymentRepo) MarkFinished(id string, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockDeploymentRepo) ListByUserID(userID string) ([]*domain.Deployment, error) {
	args := m.Called(userID)
	return args.Get(0).([]*domain.Deployment), args.Error(1)
}

func (m *MockDeploymentRepo) GetStats(userID string) (*domain.Stats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Stats), args.Error(1)
}

func (m *MockDeploymentRepo) UpdateCommitHash(id string, commitHash string) error {
	args := m.Called(id, commitHash)
	return args.Error(0)
}

// MockProjectRepo
type MockProjectRepo struct {
	mock.Mock
}

func (m *MockProjectRepo) Create(p *projectDomain.Project) error { return m.Called(p).Error(0) }
func (m *MockProjectRepo) Update(p *projectDomain.Project) error { return m.Called(p).Error(0) }
func (m *MockProjectRepo) FindByID(id string) (*projectDomain.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*projectDomain.Project), args.Error(1)
}
func (m *MockProjectRepo) ListByUserID(userID string) ([]*projectDomain.Project, error) {
	args := m.Called(userID)
	return args.Get(0).([]*projectDomain.Project), args.Error(1)
}
func (m *MockProjectRepo) Delete(id string) error { return m.Called(id).Error(0) }

// MockServerRepo
type MockServerRepo struct {
	mock.Mock
}

func (m *MockServerRepo) Create(s *serverDomain.Server) error { return m.Called(s).Error(0) }
func (m *MockServerRepo) Update(s *serverDomain.Server) error { return m.Called(s).Error(0) }
func (m *MockServerRepo) ListByUserID(userID string) ([]*serverDomain.Server, error) {
	args := m.Called(userID)
	return args.Get(0).([]*serverDomain.Server), args.Error(1)
}
func (m *MockServerRepo) FindByID(id string) (*serverDomain.Server, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*serverDomain.Server), args.Error(1)
}
func (m *MockServerRepo) Delete(id string) error { return m.Called(id).Error(0) } // Added Delete just in case interface has it now or later

func TestRollback_Success(t *testing.T) {
	mockRepo := new(MockDeploymentRepo)
	mockProjectRepo := new(MockProjectRepo)
	mockServerRepo := new(MockServerRepo)

	service := NewService(mockRepo, mockProjectRepo, mockServerRepo)

	originalID := "dep-123"
	userID := "user-1"
	commitHash := "abc123456789"

	original := &domain.Deployment{
		ID:         originalID,
		UserID:     userID,
		ProjectID:  "proj-1",
		ServerID:   "srv-1",
		Status:     "success",
		CommitHash: commitHash,
	}

	mockRepo.On("FindByID", originalID).Return(original, nil)

	// Expect a new deployment to be created with CommitHash set
	mockRepo.On("Create", mock.MatchedBy(func(d *domain.Deployment) bool {
		return d.CommitHash == commitHash &&
			d.RollbackFromID == originalID &&
			d.Status == "queued" &&
			d.ProjectID == "proj-1"
	})).Return(nil)

	res, err := service.Rollback(userID, originalID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, commitHash, res.CommitHash)
	assert.NotEmpty(t, res.ID)
	assert.NotEqual(t, originalID, res.ID)

	mockRepo.AssertExpectations(t)
}

func TestRollback_Fail_NoHash(t *testing.T) {
	mockRepo := new(MockDeploymentRepo)
	mockProjectRepo := new(MockProjectRepo)
	mockServerRepo := new(MockServerRepo)
	service := NewService(mockRepo, mockProjectRepo, mockServerRepo)

	originalID := "dep-no-hash"
	userID := "user-1"

	original := &domain.Deployment{
		ID:         originalID,
		UserID:     userID,
		CommitHash: "", // Empty hash
	}

	mockRepo.On("FindByID", originalID).Return(original, nil)

	res, err := service.Rollback(userID, originalID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no commit hash")
	assert.Nil(t, res)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockDeploymentRepo)
	mockProjectRepo := new(MockProjectRepo)
	mockServerRepo := new(MockServerRepo)
	service := NewService(mockRepo, mockProjectRepo, mockServerRepo)

	req := CreateDeploymentRequest{
		ProjectID: "p-1",
		ServerID:  "s-1",
	}

	// Mocks for validation
	mockProjectRepo.On("FindByID", "p-1").Return(&projectDomain.Project{UserID: "u-1"}, nil)
	mockServerRepo.On("FindByID", "s-1").Return(&serverDomain.Server{UserID: "u-1"}, nil)

	mockRepo.On("Create", mock.AnythingOfType("*deployment.Deployment")).Return(nil)

	res, err := service.Create("u-1", req)

	assert.NoError(t, err)
	assert.Equal(t, "queued", res.Status)
}
