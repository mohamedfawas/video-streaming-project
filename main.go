package main

import (
	"log"
	"net/http"

	"github.com/mohamedfawas/learn-video-streaming-project/handlers"
)

const (
	PORT      = ":8080"
	VIDEO_DIR = "./videos"
)

func main() {
	// Initialize video handler
	videoHandler := handlers.NewVideoHandler(VIDEO_DIR)

	// Setup routes
	http.HandleFunc("/upload", videoHandler.HandleUpload)
	http.HandleFunc("/videos", videoHandler.HandleListVideos)
	http.HandleFunc("/stream/", videoHandler.HandleStreamVideo)

	// Start the server
	log.Printf("Server starting on port %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
