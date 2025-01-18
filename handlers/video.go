package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// VideoHandler handles all video-related operations
type VideoHandler struct {
	VideoDir string
}

type VideoInfo struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type UploadResponse struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Message  string `json:"message"`
}

type ListResponse struct {
	Videos []VideoInfo `json:"videos"`
}

// NewVideoHandler creates a new VideoHandler instance
func NewVideoHandler(videoDir string) *VideoHandler {
	// Create videos directory if it doesn't exist
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		panic(err)
	}
	return &VideoHandler{
		VideoDir: videoDir,
	}
}

// HandleUpload handles video file uploads
func (h *VideoHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 32MB max memory
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to get video file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create file in videos directory
	filename := filepath.Join(h.VideoDir, header.Filename)
	dst, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	response := UploadResponse{
		Filename: header.Filename,
		Size:     size,
		Message:  "Upload successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleListVideos returns a list of all available videos
func (h *VideoHandler) HandleListVideos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := os.ReadDir(h.VideoDir)
	if err != nil {
		http.Error(w, "Failed to read videos directory", http.StatusInternalServerError)
		return
	}

	videos := []VideoInfo{}
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}
			videos = append(videos, VideoInfo{
				Filename: info.Name(),
				Size:     info.Size(),
			})
		}
	}

	response := ListResponse{Videos: videos}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleStreamVideo streams a video file with support for range requests
func (h *VideoHandler) HandleStreamVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/stream/")
	filepath := filepath.Join(h.VideoDir, filename)

	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	// Handle range requests
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// Parse range header
		var start, end int64
		fmt.Sscanf(strings.TrimPrefix(rangeHeader, "bytes="), "%d-%d", &start, &end)

		if end == 0 {
			end = fileInfo.Size() - 1
		}

		// Seek to start position
		file.Seek(start, 0)

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	}

	w.Header().Set("Content-Type", "video/mp4")
	io.Copy(w, file)
}
