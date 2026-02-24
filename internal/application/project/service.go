package project

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "deployes/internal/domain/project"
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

func (s *Service) Create(userID string, req CreateProjectRequest) (*ProjectResponse, error) {
	if req.Name == "" || req.DeployScript == "" {
		return nil, errors.New("name and deployScript are required")
	}

	// Validate project type
	projectType := domain.ProjectType(req.Type)
	if projectType != domain.ProjectTypeGitHub && projectType != domain.ProjectTypeLocal {
		return nil, errors.New("invalid project type: must be 'github' or 'local'")
	}

	// For GitHub projects, validate repository URL and branch
	if projectType == domain.ProjectTypeGitHub {
		if req.RepoURL == "" || req.Branch == "" {
			return nil, errors.New("repoUrl and branch are required for GitHub projects")
		}
		if err := utils.ValidateRepoURL(req.RepoURL); err != nil {
			return nil, errors.New("invalid repository URL format (supported: GitHub, GitLab, Bitbucket)")
		}
	}

	// For Local projects, validate local path
	if projectType == domain.ProjectTypeLocal && req.LocalPath == "" {
		return nil, errors.New("localPath is required for local projects")
	}

	p := &domain.Project{
		ID:               uuid.NewString(),
		UserID:           userID,
		Name:             req.Name,
		Type:             projectType,
		RepoURL:          req.RepoURL,
		Branch:           req.Branch,
		LocalPath:        req.LocalPath,
		DeployScript:     req.DeployScript,
		IncludePatterns:  req.IncludePatterns,
		ExcludePatterns:  req.ExcludePatterns,
		PreservePatterns: req.PreservePatterns,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Generate SSH Key Pair only for GitHub projects
	if projectType == domain.ProjectTypeGitHub {
		priv, pub, err := utils.GenerateSSHKeyPair()
		if err == nil {
			if encPriv, err := utils.Encrypt(priv, s.encryptionKey); err == nil {
				p.SCMPrivateKeyEncrypted = encPriv
				p.SCMPublicKey = pub
			}
		}
	}

	if err := s.repo.Create(p); err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:               p.ID,
		Name:             p.Name,
		Type:             string(p.Type),
		RepoURL:          p.RepoURL,
		Branch:           p.Branch,
		LocalPath:        p.LocalPath,
		DeployScript:     p.DeployScript,
		IncludePatterns:  p.IncludePatterns,
		ExcludePatterns:  p.ExcludePatterns,
		PreservePatterns: p.PreservePatterns,
		SCMPublicKey:     p.SCMPublicKey,
	}, nil
}

func (s *Service) Update(userID string, id string, req CreateProjectRequest) (*ProjectResponse, error) {
	// Find project
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("project not found")
	}

	// Verify ownership
	if p.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Validate based on project type
	projectType := domain.ProjectType(req.Type)
	if projectType != domain.ProjectTypeGitHub && projectType != domain.ProjectTypeLocal {
		return nil, errors.New("invalid project type: must be 'github' or 'local'")
	}

	// For GitHub projects, validate repository URL
	if projectType == domain.ProjectTypeGitHub && req.RepoURL != "" {
		if err := utils.ValidateRepoURL(req.RepoURL); err != nil {
			return nil, errors.New("invalid repository URL format (supported: GitHub, GitLab, Bitbucket)")
		}
	}

	// Update fields
	p.Name = req.Name
	p.Type = projectType
	p.RepoURL = req.RepoURL
	p.Branch = req.Branch
	p.LocalPath = req.LocalPath
	p.DeployScript = req.DeployScript
	p.IncludePatterns = req.IncludePatterns
	p.ExcludePatterns = req.ExcludePatterns
	p.PreservePatterns = req.PreservePatterns
	p.UpdatedAt = time.Now()

	// If no key exists for GitHub project (legacy project), generate one
	if p.Type == domain.ProjectTypeGitHub && p.SCMPrivateKeyEncrypted == "" {
		priv, pub, err := utils.GenerateSSHKeyPair()
		if err == nil {
			if encPriv, err := utils.Encrypt(priv, s.encryptionKey); err == nil {
				p.SCMPrivateKeyEncrypted = encPriv
				p.SCMPublicKey = pub
			}
		}
	}

	if err := s.repo.Update(p); err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:               p.ID,
		Name:             p.Name,
		Type:             string(p.Type),
		RepoURL:          p.RepoURL,
		Branch:           p.Branch,
		LocalPath:        p.LocalPath,
		DeployScript:     p.DeployScript,
		IncludePatterns:  p.IncludePatterns,
		ExcludePatterns:  p.ExcludePatterns,
		PreservePatterns: p.PreservePatterns,
		SCMPublicKey:     p.SCMPublicKey,
	}, nil
}

func (s *Service) List(userID string) ([]*ProjectResponse, error) {

	projects, err := s.repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []*ProjectResponse
	for _, p := range projects {
		result = append(result, &ProjectResponse{
			ID:               p.ID,
			Name:             p.Name,
			Type:             string(p.Type),
			RepoURL:          p.RepoURL,
			Branch:           p.Branch,
			LocalPath:        p.LocalPath,
			DeployScript:     p.DeployScript,
			IncludePatterns:  p.IncludePatterns,
			ExcludePatterns:  p.ExcludePatterns,
			PreservePatterns: p.PreservePatterns,
			SCMPublicKey:     p.SCMPublicKey,
		})
	}

	return result, nil
}

// FindByID returns the project domain object by ID (used by webhook handler)
func (s *Service) FindByID(id string) (*domain.Project, error) {
	return s.repo.FindByID(id)
}

// TestRepoAccess tests if a repository can be accessed
func (s *Service) TestRepoAccess(req TestRepoAccessRequest) (*RepoAccessResult, error) {
	// Validate URL format first
	if err := utils.ValidateRepoURL(req.RepoURL); err != nil {
		return &RepoAccessResult{
			Accessible: false,
			IsPrivate:  false,
			Message:    "Geçersiz repository URL formatı",
			RepoType:   "unknown",
		}, nil
	}

	// Determine repo type
	repoType := "unknown"
	if strings.Contains(req.RepoURL, "github.com") {
		repoType = "github"
	} else if strings.Contains(req.RepoURL, "gitlab.com") {
		repoType = "gitlab"
	} else if strings.Contains(req.RepoURL, "bitbucket.org") {
		repoType = "bitbucket"
	}

	// Try to access the repository using git ls-remote
	branch := req.Branch
	if branch == "" {
		branch = "HEAD"
	}

	cmd := exec.Command("git", "ls-remote", "--exit-code", req.RepoURL, branch)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0") // Disable password prompt

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Printf("[TestRepoAccess] Executing: %s\n", cmd.String())
	err := cmd.Run()

	if err != nil {
		stderrStr := stderr.String()
		fmt.Printf("[TestRepoAccess] Failed: %v\nStderr: %s\n", err, stderrStr)

		// Check for common error patterns
		if strings.Contains(stderrStr, "Repository not found") ||
			strings.Contains(stderrStr, "not found") ||
			strings.Contains(stderrStr, "authentication") ||
			strings.Contains(stderrStr, "Permission denied") ||
			strings.Contains(stderrStr, "could not read Username") {

			// Likely private repo
			guidance := s.getPrivateRepoGuidance(repoType, req.RepoURL)
			return &RepoAccessResult{
				Accessible:   false,
				IsPrivate:    true,
				Message:      "Bu repository'ye erişilemiyor. Muhtemelen private bir repository.",
				Guidance:     guidance,
				RepoType:     repoType,
				BranchExists: false,
			}, nil
		}

		// Branch doesn't exist
		if strings.Contains(stderrStr, "Could not find remote branch") ||
			cmd.ProcessState.ExitCode() == 2 {
			return &RepoAccessResult{
				Accessible:   true,
				IsPrivate:    false,
				Message:      fmt.Sprintf("Repository erişilebilir ancak '%s' branch'i bulunamadı.", branch),
				RepoType:     repoType,
				BranchExists: false,
			}, nil
		}

		return &RepoAccessResult{
			Accessible: false,
			IsPrivate:  false,
			Message:    "Repository'ye erişim başarısız: " + stderrStr,
			RepoType:   repoType,
		}, nil
	}

	// Success - repo is accessible
	return &RepoAccessResult{
		Accessible:   true,
		IsPrivate:    false,
		Message:      "Repository erişilebilir ve branch mevcut! ✓",
		RepoType:     repoType,
		BranchExists: true,
	}, nil
}

func (s *Service) getPrivateRepoGuidance(repoType, repoURL string) string {
	return "deployes bu proje için otomatik bir SSH Deploy Key oluşturacak.\n" +
		"Projeyi oluşturduktan sonra, 'Public Key'i kopyalayıp GitHub/GitLab 'Deploy Keys' bölümüne eklemeniz yeterlidir."
}

// Delete deletes a project by ID
func (s *Service) Delete(userID, id string) error {
	// Find project
	p, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("project not found")
	}

	// Verify ownership
	if p.UserID != userID {
		return errors.New("unauthorized")
	}

	// If it's a local project, delete the uploaded files FIRST
	if p.Type == domain.ProjectTypeLocal && p.LocalPath != "" {
		fmt.Printf("Delete: Removing local files at %s\n", p.LocalPath)
		if err := os.RemoveAll(p.LocalPath); err != nil {
			// Return error if file deletion fails - don't delete the project record
			fmt.Printf("Delete: Failed to delete local files: %v\n", err)
			return fmt.Errorf("failed to delete local files: %w", err)
		}
		fmt.Printf("Delete: Successfully deleted local files\n")
	}

	// Delete the project from database
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
