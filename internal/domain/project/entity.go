package project

import "time"

type ProjectType string

const (
	ProjectTypeGitHub ProjectType = "github"
	ProjectTypeLocal  ProjectType = "local"
)

type Project struct {
	ID                     string
	UserID                 string
	Name                   string
	Type                   ProjectType // "github" or "local"
	RepoURL                string // For github projects
	Branch                 string // For github projects
	LocalPath              string // For local projects - path to uploaded files
	DeployScript           string
	IncludePatterns        string // Comma-separated patterns: "src/*,package.json,*.go"
	ExcludePatterns        string // Patterns to exclude: "node_modules,*.log,.env.local"
	PreservePatterns       string // Files not to delete on server: "uploads/*,logs/*,.env"
	SCMPrivateKeyEncrypted string
	SCMPublicKey           string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
