package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FileRequest represents the file operation request
// @Description File operation request payload
type FileRequest struct {
	// File path
	// required: true
	Path string `json:"path" example:"/var/www/html/index.html"`
	// File content (for write operations)
	// required: false
	Content string `json:"content,omitempty" example:"<html>...</html>"`
}

// FileResponse represents the file operation response
// @Description File operation response payload
type FileResponse struct {
	// File content
	// required: true
	Content string `json:"content" example:"<html>...</html>"`
	// File size in bytes
	// required: true
	Size int64 `json:"size" example:"1234"`
	// Error message if any
	// required: false
	Error string `json:"error,omitempty" example:"file not found"`
}

// HandleFileRead implements a vulnerable file read endpoint
// @Summary Read file content
// @Description Read file content (intentionally vulnerable to path traversal)
// @Tags files
// @Accept json
// @Produce json
// @Param request body FileRequest true "File path"
// @Security BearerAuth
// @Success 200 {object} FileResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /files/read [post]
func HandleFileRead(w http.ResponseWriter, r *http.Request) {
	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Path traversal through direct path usage
	// Vulnerable: No path sanitization
	// Vulnerable: No access control
	// Vulnerable: No path validation
	content, err := os.ReadFile(req.Path)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	fileInfo, err := os.Stat(req.Path)
	if err != nil {
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	response := FileResponse{
		Content: string(content),
		Size:    fileInfo.Size(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleFileDownload implements a vulnerable file download endpoint
// @Summary Download file
// @Description Download file (intentionally vulnerable to path traversal)
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Param request body FileRequest true "File path"
// @Security BearerAuth
// @Success 200 {file} binary
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /files/download [post]
func HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Path traversal through direct path usage
	// Vulnerable: No path sanitization
	// Vulnerable: No access control
	// Vulnerable: No path validation
	file, err := os.Open(req.Path)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file info for Content-Disposition
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(req.Path)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Stream file to response
	io.Copy(w, file)
}

// HandleFileWrite implements a vulnerable file write endpoint
// @Summary Write file content
// @Description Write file content (intentionally vulnerable to path traversal)
// @Tags files
// @Accept json
// @Produce json
// @Param request body FileRequest true "File path and content"
// @Security BearerAuth
// @Success 200 {object} FileResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /files/write [post]
func HandleFileWrite(w http.ResponseWriter, r *http.Request) {
	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Path traversal through direct path usage
	// Vulnerable: No path sanitization
	// Vulnerable: No access control
	// Vulnerable: No path validation
	// Vulnerable: No content validation
	err := os.WriteFile(req.Path, []byte(req.Content), 0644)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	fileInfo, err := os.Stat(req.Path)
	if err != nil {
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	response := FileResponse{
		Content: req.Content,
		Size:    fileInfo.Size(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleFileList implements a vulnerable file listing endpoint
// @Summary List directory contents
// @Description List directory contents (intentionally vulnerable to path traversal)
// @Tags files
// @Accept json
// @Produce json
// @Param request body FileRequest true "Directory path"
// @Security BearerAuth
// @Success 200 {array} string
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /files/list [post]
func HandleFileList(w http.ResponseWriter, r *http.Request) {
	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Path traversal through direct path usage
	// Vulnerable: No path sanitization
	// Vulnerable: No access control
	// Vulnerable: No path validation
	entries, err := os.ReadDir(req.Path)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
