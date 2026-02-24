package workers

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	deploymentDomain "deployes/internal/domain/deployment"
	projectDomain "deployes/internal/domain/project"
	secretDomain "deployes/internal/domain/secret"
	serverDomain "deployes/internal/domain/server"
	"deployes/pkg/utils"

	"archive/tar"
	"io"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type DeploymentRepository interface {
	GetNextQueued() (*deploymentDomain.Deployment, error)
	AppendLog(id string, logLine string) error
	MarkRunning(id string) error
	MarkFinished(id string, status string) error
	UpdateCommitHash(id string, commitHash string) error
}

type ProjectRepository interface {
	FindByID(id string) (*projectDomain.Project, error)
}

type ServerRepository interface {
	FindByID(id string) (*serverDomain.Server, error)
}

type SecretRepository interface {
	ListByProjectID(projectID string) ([]*secretDomain.Secret, error)
}

type DeployWorker struct {
	deploymentRepo DeploymentRepository
	projectRepo    ProjectRepository
	serverRepo     ServerRepository
	secretRepo     SecretRepository
	encryptionKey  string
	broadcaster    *LogBroadcaster
}

func NewDeployWorker(
	deploymentRepo DeploymentRepository,
	projectRepo ProjectRepository,
	serverRepo ServerRepository,
	secretRepo SecretRepository,
	encryptionKey string,
) *DeployWorker {
	return &DeployWorker{
		deploymentRepo: deploymentRepo,
		projectRepo:    projectRepo,
		serverRepo:     serverRepo,
		secretRepo:     secretRepo,
		encryptionKey:  encryptionKey,
		broadcaster:    GetBroadcaster(),
	}
}

func (w *DeployWorker) Start() {
	log.Println("✅ Deploy worker started")

	for {
		job, err := w.deploymentRepo.GetNextQueued()
		if err != nil {
			log.Println("worker error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if job == nil {
			time.Sleep(2 * time.Second)
			continue
		}

		w.processJob(job)
	}
}

func (w *DeployWorker) processJob(job *deploymentDomain.Deployment) {
	// Mark as running
	_ = w.deploymentRepo.MarkRunning(job.ID)
	w.appendLog(job.ID, "🚀 Starting deployment...\n")

	// Get project details
	project, err := w.projectRepo.FindByID(job.ProjectID)
	if err != nil {
		w.failDeployment(job.ID, fmt.Sprintf("❌ Failed to get project: %v\n", err))
		return
	}
	w.appendLog(job.ID, fmt.Sprintf("📦 Project: %s\n", project.Name))
	w.appendLog(job.ID, fmt.Sprintf("📁 Type: %s\n", project.Type))

	// Log based on project type
	if project.Type == projectDomain.ProjectTypeGitHub {
		w.appendLog(job.ID, fmt.Sprintf("🔗 Repository: %s\n", project.RepoURL))
		w.appendLog(job.ID, fmt.Sprintf("🌿 Branch: %s\n", project.Branch))
	} else if project.Type == projectDomain.ProjectTypeLocal {
		w.appendLog(job.ID, fmt.Sprintf("📂 Local Path: %s\n", project.LocalPath))
	}

	// Get server details
	server, err := w.serverRepo.FindByID(job.ServerID)
	if err != nil {
		w.failDeployment(job.ID, fmt.Sprintf("❌ Failed to get server: %v\n", err))
		return
	}
	w.appendLog(job.ID, fmt.Sprintf("🖥️  Server: %s (%s:%d)\n", server.Name, server.Host, server.Port))

	// Decrypt SSH key
	w.appendLog(job.ID, "🔐 Decrypting SSH key...\n")
	sshKey, err := utils.Decrypt(server.SSHKeyEncrypted, w.encryptionKey)
	if err != nil {
		w.failDeployment(job.ID, fmt.Sprintf("❌ Failed to decrypt SSH key: %v\n", err))
		return
	}

	// Connect to server via SSH
	w.appendLog(job.ID, fmt.Sprintf("🔌 Connecting to %s@%s:%d...\n", server.Username, server.Host, server.Port))
	client, err := utils.CreateSSHClient(server.Host, server.Port, server.Username, sshKey)
	if err != nil {
		w.failDeployment(job.ID, fmt.Sprintf("❌ SSH connection failed: %v\n", err))
		return
	}
	defer client.Close()
	w.appendLog(job.ID, "✅ Connected to server\n")

	// --- Private Repo Handling (SSH Key Injection) - GitHub only ---
	var gitEnvPrefix string
	if project.Type == projectDomain.ProjectTypeGitHub && project.SCMPrivateKeyEncrypted != "" {
		w.appendLog(job.ID, "🔑 Configuring SSH Deploy Key...\n")
		scmKey, err := utils.Decrypt(project.SCMPrivateKeyEncrypted, w.encryptionKey)
		if err != nil {
			w.failDeployment(job.ID, fmt.Sprintf("❌ Failed to decrypt SCM Key: %v\n", err))
			return
		}

		b64Key := base64.StdEncoding.EncodeToString([]byte(scmKey))
		keyPath := fmt.Sprintf("/home/%s/.deploy_keys/%s", server.Username, project.ID)

		// Create dir, write key, chmod
		setupCmd := fmt.Sprintf("mkdir -p /home/%s/.deploy_keys && echo '%s' | base64 -d > %s && chmod 600 %s",
			server.Username, b64Key, keyPath, keyPath)

		if output, err := utils.RunSSHCommand(client, setupCmd); err != nil {
			w.failDeployment(job.ID, fmt.Sprintf("❌ Failed to setup Deploy Key: %v\nOutput: %s\n", err, output))
			return
		}
		// Set GIT_SSH_COMMAND to use this key
		gitEnvPrefix = fmt.Sprintf("GIT_SSH_COMMAND='ssh -i %s -o StrictHostKeyChecking=no' ", keyPath)
	}
	// ---------------------------------------------

	// Create project directory name
	var projectDir string
	if project.Type == projectDomain.ProjectTypeGitHub {
		repoName := extractRepoName(project.RepoURL)
		projectDir = fmt.Sprintf("/home/%s/%s", server.Username, repoName)
	} else {
		// For local projects, use project name as directory name
		projectDir = fmt.Sprintf("/home/%s/%s", server.Username, project.Name)
	}

	// Check if directory exists
	w.appendLog(job.ID, fmt.Sprintf("📂 Checking project directory: %s\n", projectDir))

	checkCmd := fmt.Sprintf("test -d %s && echo 'exists' || echo 'not_exists'", projectDir)
	checkOutput, _ := utils.RunSSHCommand(client, checkCmd)

	// Deployment Operations based on project type
	if project.Type == projectDomain.ProjectTypeLocal {
		// LOCAL PROJECT DEPLOYMENT
		w.appendLog(job.ID, "📤 Deploying local project files...\n")

		// Check if local path exists on deployes server
		if _, err := os.Stat(project.LocalPath); os.IsNotExist(err) {
			w.failDeployment(job.ID, fmt.Sprintf("❌ Local path does not exist: %s\n", project.LocalPath))
			return
		}

		// Create directory on remote server if it doesn't exist
		if strings.TrimSpace(checkOutput) != "exists" {
			w.appendLog(job.ID, fmt.Sprintf("📁 Creating directory: %s\n", projectDir))
			mkdirCmd := fmt.Sprintf("mkdir -p %s", projectDir)
			if output, err := utils.RunSSHCommand(client, mkdirCmd); err != nil {
				w.failDeployment(job.ID, fmt.Sprintf("❌ Failed to create directory: %v\nOutput: %s\n", err, output))
				return
			}
		}

		// Transfer files using TAR stream over SSH
		w.appendLog(job.ID, "📦 Transferring files to server (Native)...\n")

		if err := w.copyFilesOverSSH(client, project.LocalPath, projectDir); err != nil {
			w.failDeployment(job.ID, fmt.Sprintf("❌ File transfer failed: %v\n", err))
			return
		}
		w.appendLog(job.ID, "✅ Files transferred successfully\n")

	} else {
		// GITHUB PROJECT DEPLOYMENT
		// Git Operations
		if job.CommitHash != "" {
			// ROLLBACK / SPECIFIC COMMIT
			w.appendLog(job.ID, fmt.Sprintf("🔙 Rolling back to commit: %s\n", job.CommitHash))

			var gitCmd string
			if strings.TrimSpace(checkOutput) == "exists" {
				// Fetch all and checkout hash
				gitCmd = fmt.Sprintf("cd %s && %sgit fetch origin && %sgit checkout %s", projectDir, gitEnvPrefix, gitEnvPrefix, job.CommitHash)
			} else {
				// Clone and checkout
				gitCmd = fmt.Sprintf("%sgit clone %s %s && cd %s && %sgit checkout %s", gitEnvPrefix, project.RepoURL, projectDir, projectDir, gitEnvPrefix, job.CommitHash)
			}

			output, err := utils.RunSSHCommand(client, gitCmd)
			if err != nil {
				w.failDeployment(job.ID, fmt.Sprintf("❌ Rollback checkout failed: %v\nOutput: %s\n", err, output))
				return
			}
			w.appendLog(job.ID, fmt.Sprintf("📝 Git output:\n%s\n", output))

		} else {
			// NORMAL DEPLOYMENT
			if strings.TrimSpace(checkOutput) == "exists" {
				// Directory exists, do git pull
				w.appendLog(job.ID, "📥 Repository exists, pulling latest changes...\n")
				// Note: Pulling with specific key
				pullCmd := fmt.Sprintf("cd %s && %sgit fetch origin && %sgit checkout %s && %sgit pull origin %s",
					projectDir, gitEnvPrefix, gitEnvPrefix, project.Branch, gitEnvPrefix, project.Branch)
				output, err := utils.RunSSHCommand(client, pullCmd)
				if err != nil {
					w.failDeployment(job.ID, fmt.Sprintf("❌ Git pull failed: %v\nOutput: %s\n", err, output))
					return
				}
				w.appendLog(job.ID, fmt.Sprintf("📝 Git output:\n%s\n", output))
			} else {
				// Directory doesn't exist, clone repository
				w.appendLog(job.ID, "📥 Cloning repository check...\n")
				cloneCmd := fmt.Sprintf("%sgit clone -b %s %s %s", gitEnvPrefix, project.Branch, project.RepoURL, projectDir)
				output, err := utils.RunSSHCommand(client, cloneCmd)
				if err != nil {
					w.failDeployment(job.ID, fmt.Sprintf("❌ Git clone failed: %v\nOutput: %s\n", err, output))
					return
				}
				w.appendLog(job.ID, fmt.Sprintf("📝 Git output:\n%s\n", output))
			}

			// Capture Commit Hash
			hashCmd := fmt.Sprintf("cd %s && git rev-parse HEAD", projectDir)
			hashOutput, err := utils.RunSSHCommand(client, hashCmd)
			if err == nil {
				commitHash := strings.TrimSpace(hashOutput)
				w.appendLog(job.ID, fmt.Sprintf("📌 Deployed Commit: %s\n", commitHash))
				_ = w.deploymentRepo.UpdateCommitHash(job.ID, commitHash)
			}
		}
		w.appendLog(job.ID, "✅ Repository synchronized\n")
	}

	// Execute deploy script
	w.appendLog(job.ID, "⚙️  Executing deploy script...\n")
	w.appendLog(job.ID, fmt.Sprintf("📜 Script:\n%s\n", project.DeployScript))
	w.appendLog(job.ID, "------- Script Output -------\n")

	// Validate deploy script for dangerous patterns
	if err := utils.ValidateDeployScript(project.DeployScript); err != nil {
		w.failDeployment(job.ID, fmt.Sprintf("❌ Deploy script validation failed: %v\n", err))
		return
	}

	// 6. Injection: fetch secrets and construct env vars
	w.appendLog(job.ID, "🔐 Preparing environment variables...\n")
	secrets, err := w.secretRepo.ListByProjectID(job.ProjectID)
	var envPrefix string
	if err == nil {
		var exports []string
		for _, s := range secrets {
			decrypted, err := utils.Decrypt(s.Value, w.encryptionKey)
			if err != nil {
				w.appendLog(job.ID, fmt.Sprintf("⚠️ Failed to decrypt secret %s, skipping\n", s.Key))
				continue
			}
			// Simple escaping
			safeVal := strings.ReplaceAll(decrypted, "'", "'\\''")
			exports = append(exports, fmt.Sprintf("export %s='%s'", s.Key, safeVal))
		}
		if len(exports) > 0 {
			envPrefix = strings.Join(exports, " && ") + " && "
			w.appendLog(job.ID, fmt.Sprintf("✅ Loaded %d secrets\n", len(exports)))
		}
	}

	deployCmd := fmt.Sprintf("cd %s && %s%s", projectDir, envPrefix, project.DeployScript)
	output, err := utils.RunSSHCommand(client, deployCmd)

	// Log output line by line
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			w.appendLog(job.ID, line+"\n")
		}
	}
	w.appendLog(job.ID, "------- End Output -------\n")

	if err != nil {
		w.failDeployment(job.ID, fmt.Sprintf("❌ Deploy script failed: %v\n", err))
		return
	}

	// Success!
	w.appendLog(job.ID, "✅ Deployment completed successfully!\n")
	_ = w.deploymentRepo.MarkFinished(job.ID, deploymentDomain.StatusSuccess)
}

func (w *DeployWorker) appendLog(deploymentID, message string) {
	// Save to database
	_ = w.deploymentRepo.AppendLog(deploymentID, message)

	// Broadcast to WebSocket subscribers
	w.broadcaster.Broadcast(LogMessage{
		DeploymentID: deploymentID,
		Message:      message,
		Timestamp:    time.Now().UnixMilli(),
	})
}

func (w *DeployWorker) failDeployment(deploymentID, errorMessage string) {
	w.appendLog(deploymentID, errorMessage)
	_ = w.deploymentRepo.MarkFinished(deploymentID, deploymentDomain.StatusFailed)
}

func (w *DeployWorker) copyFilesOverSSH(client *ssh.Client, localPath, remotePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Start remote tar command
	// -x: extract
	// -C: change directory
	// -: read from stdin
	// m: touch modification time
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		tw := tar.NewWriter(w)
		defer tw.Close()

		// Walk through local path
		filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Get relative path
			relPath, err := filepath.Rel(localPath, path)
			if err != nil {
				return err
			}

			// Skip root directory entry
			if relPath == "." {
				return nil
			}

			// Create compatible header
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(relPath) // Ensure forward slashes

			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				if _, err := io.Copy(tw, file); err != nil {
					return err
				}
			}
			return nil
		})
	}()

	// Wait for remote command to finish
	if err := session.Run(fmt.Sprintf("tar -xm -C %s", remotePath)); err != nil {
		return fmt.Errorf("remote tar failed: %w", err)
	}

	return nil
}

// extractRepoName extracts the repository name from a GitHub URL
func extractRepoName(repoURL string) string {
	// Handle both HTTPS and SSH URLs
	// https://github.com/user/repo.git -> repo
	// git@github.com:user/repo.git -> repo

	url := strings.TrimSuffix(repoURL, ".git")
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "project"
}
