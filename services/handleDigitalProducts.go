package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	uploadPath = "./uploads/"
	// baseURL    = "http://localhost:8082/download/"
)

// Handler for file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate unique file name using timestamp
	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
	fileExt := filepath.Ext(fileHeader.Filename)
	uniqueFileName := fmt.Sprintf("%s%s", timestamp, fileExt)
	outFile, err := os.Create(filepath.Join(uploadPath, uniqueFileName))
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Generate URLs
	// fileURL := fmt.Sprintf("%s%s", baseURL, uniqueFileName)

	// Return the file URL and thumbnail URL in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{"message": "File uploaded successfully", "fileName": "%s"}`, uniqueFileName)
	w.Write([]byte(response))
}

// Handler for file downloads
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, filename)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set headers and serve the file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}
}
