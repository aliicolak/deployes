package project

type CreateProjectRequest struct {
	Name             string `json:"name"`
	Type             string `json:"type"` // "github" or "local"
	RepoURL          string `json:"repoUrl,omitempty"`
	Branch           string `json:"branch,omitempty"`
	LocalPath        string `json:"localPath,omitempty"` // For local projects
	DeployScript     string `json:"deployScript"`
	IncludePatterns  string `json:"includePatterns,omitempty"`
	ExcludePatterns  string `json:"excludePatterns,omitempty"`
	PreservePatterns string `json:"preservePatterns,omitempty"`
}

type ProjectResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	RepoURL          string `json:"repoUrl,omitempty"`
	Branch           string `json:"branch,omitempty"`
	LocalPath        string `json:"localPath,omitempty"`
	DeployScript     string `json:"deployScript"`
	IncludePatterns  string `json:"includePatterns,omitempty"`
	ExcludePatterns  string `json:"excludePatterns,omitempty"`
	PreservePatterns string `json:"preservePatterns,omitempty"`
	SCMPublicKey     string `json:"scmPublicKey,omitempty"`
}

type TestRepoAccessRequest struct {
	RepoURL string `json:"repoUrl"`
	Branch  string `json:"branch,omitempty"`
}

type RepoAccessResult struct {
	Accessible   bool   `json:"accessible"`
	IsPrivate    bool   `json:"isPrivate"`
	Message      string `json:"message"`
	Guidance     string `json:"guidance,omitempty"`
	RepoType     string `json:"repoType"` // github, gitlab, bitbucket, unknown
	BranchExists bool   `json:"branchExists"`
}
