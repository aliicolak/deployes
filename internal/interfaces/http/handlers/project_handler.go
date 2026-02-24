package handlers

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	app "deployes/internal/application/project"
	"deployes/internal/interfaces/http/middleware"
)

type ProjectHandler struct {
	service *app.Service
}

func NewProjectHandler(service *app.Service) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// POST /api/projects
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(middleware.UserIDKey).(string)

	var req app.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(userId, req)
	if err != nil {
		// Don't leak internal error details
		http.Error(w, "failed to create project", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// PUT /api/projects?id={id}
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var req app.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Update(userId, id, req)
	if err != nil {
		// Don't leak internal error details
		http.Error(w, "failed to update project", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GET /api/projects
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(middleware.UserIDKey).(string)

	res, err := h.service.List(userId)
	if err != nil {
		log.Printf("Failed to fetch projects: %v", err)
		http.Error(w, "failed to fetch projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// POST /api/projects/test-access
func (h *ProjectHandler) TestRepoAccess(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in TestRepoAccess: %v", r)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}()

	var req app.TestRepoAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("TestRepoAccess: invalid body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.TestRepoAccess(req)
	if err != nil {
		log.Printf("TestRepoAccess: service failed: %v", err)
		http.Error(w, "failed to test repository access", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// POST /api/projects/upload
func (h *ProjectHandler) UploadLocalProject(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 500MB for folders)
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		log.Printf("UploadLocalProject: failed to parse form: %v", err)
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// Check for folder upload first (multiple files with "files" key)
	if r.MultipartForm != nil && len(r.MultipartForm.File["files"]) > 0 {
		h.handleFolderUpload(w, r)
		return
	}

	// Single file upload (legacy)
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("UploadLocalProject: no file provided: %v", err)
		http.Error(w, "no file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get project name from form
	projectName := r.FormValue("projectName")
	if projectName == "" {
		projectName = strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	}

	// Sanitize project name
	projectName = strings.ReplaceAll(projectName, " ", "_")
	projectName = strings.ReplaceAll(projectName, "/", "_")
	projectName = strings.ReplaceAll(projectName, "\\", "_")

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads/projects"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("UploadLocalProject: failed to create uploads dir: %v", err)
		http.Error(w, "failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	// Create project directory
	projectPath := filepath.Join(uploadsDir, projectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		log.Printf("UploadLocalProject: failed to create project dir: %v", err)
		http.Error(w, "failed to create project directory", http.StatusInternalServerError)
		return
	}

	// Handle zip file or single file
	if strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		// For zip files, extract to project directory
		zipPath := filepath.Join(projectPath, header.Filename)
		dst, err := os.Create(zipPath)
		if err != nil {
			log.Printf("UploadLocalProject: failed to create zip file: %v", err)
			http.Error(w, "failed to save zip file", http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			dst.Close()
			log.Printf("UploadLocalProject: failed to write zip file: %v", err)
			http.Error(w, "failed to write zip file", http.StatusInternalServerError)
			return
		}
		dst.Close()

		// Extract zip
		if err := h.extractZip(zipPath, projectPath); err != nil {
			log.Printf("UploadLocalProject: failed to extract zip: %v", err)
			http.Error(w, "failed to extract zip file", http.StatusInternalServerError)
			return
		}

		// Remove zip file after extraction
		os.Remove(zipPath)

	} else {
		// For single files, save directly
		filePath := filepath.Join(projectPath, header.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("UploadLocalProject: failed to create file: %v", err)
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("UploadLocalProject: failed to write file: %v", err)
			http.Error(w, "failed to write file", http.StatusInternalServerError)
			return
		}
	}

	// Return the local path
	response := map[string]string{
		"localPath": projectPath,
		"message":   "File uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleFolderUpload handles multiple file uploads (folder structure preserved)
func (h *ProjectHandler) handleFolderUpload(w http.ResponseWriter, r *http.Request) {
	// Get project name from form
	projectName := r.FormValue("projectName")
	if projectName == "" {
		projectName = "project_" + fmt.Sprintf("%d", time.Now().Unix())
	}

	// Sanitize project name
	projectName = strings.ReplaceAll(projectName, " ", "_")
	projectName = strings.ReplaceAll(projectName, "/", "_")
	projectName = strings.ReplaceAll(projectName, "\\", "_")

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads/projects"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("handleFolderUpload: failed to create uploads dir: %v", err)
		http.Error(w, "failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	// Create project directory
	projectPath := filepath.Join(uploadsDir, projectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		log.Printf("handleFolderUpload: failed to create project dir: %v", err)
		http.Error(w, "failed to create project directory", http.StatusInternalServerError)
		return
	}

	// Get all files from the form
	files := r.MultipartForm.File["files"]
	uploadedCount := 0

	for i, fileHeader := range files {
		// Get the relative path from form data (sent from frontend)
		relativePath := r.FormValue(fmt.Sprintf("file_path_%d", i))
		if relativePath == "" {
			// Fallback to filename if no path provided
			relativePath = fileHeader.Filename
		}

		// Open uploaded file
		src, err := fileHeader.Open()
		if err != nil {
			log.Printf("handleFolderUpload: failed to open file %s: %v", relativePath, err)
			continue
		}

		// Use the relative path to preserve folder structure
		// The path from webkitRelativePath contains the folder name as first part
		// We need to strip the root folder name to avoid nesting
		pathParts := strings.Split(relativePath, "/")
		if len(pathParts) > 1 {
			// Remove the root folder name, keep subdirectories
			relativePath = strings.Join(pathParts[1:], "/")
		} else {
			// Single file in root
			relativePath = pathParts[0]
		}

		dstPath := filepath.Join(projectPath, relativePath)

		// Create parent directories if needed
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			log.Printf("handleFolderUpload: failed to create directory for %s: %v", dstPath, err)
			src.Close()
			continue
		}

		// Create destination file
		dst, err := os.Create(dstPath)
		if err != nil {
			log.Printf("handleFolderUpload: failed to create file %s: %v", dstPath, err)
			src.Close()
			continue
		}

		// Copy file content
		if _, err := io.Copy(dst, src); err != nil {
			log.Printf("handleFolderUpload: failed to copy file %s: %v", dstPath, err)
			dst.Close()
			src.Close()
			continue
		}

		dst.Close()
		src.Close()
		uploadedCount++
	}

	if uploadedCount == 0 {
		http.Error(w, "no files were successfully uploaded", http.StatusBadRequest)
		return
	}

	// Return the local path
	response := map[string]string{
		"localPath": projectPath,
		"message":   fmt.Sprintf("Successfully uploaded %d files", uploadedCount),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// extractZip extracts a zip file to the specified directory
func (h *ProjectHandler) extractZip(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Extract each file
	for _, file := range reader.File {
		// Construct the full file path
		path := filepath.Join(destDir, file.Name)

		// Check for Zip Slip vulnerability
		if !strings.HasPrefix(path, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			// Create directory
			os.MkdirAll(path, file.Mode())
			continue
		}

		// Create parent directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Extract file
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		destFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return fmt.Errorf("failed to create destination file: %w", err)
		}

		if _, err := io.Copy(destFile, fileReader); err != nil {
			fileReader.Close()
			destFile.Close()
			return fmt.Errorf("failed to write file: %w", err)
		}

		fileReader.Close()
		destFile.Close()
	}

	return nil
}

// DELETE /api/projects
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteProject: Deleting project %s for user %s", id, userID)

	if err := h.service.Delete(userID, id); err != nil {
		log.Printf("DeleteProject: failed to delete project: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteProject: Successfully deleted project %s", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Project deleted successfully"})
}
