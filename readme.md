# Video Streaming Project

## Overview
This project is a simple HTTP server built with Go that provides RESTful API endpoints for uploading, listing, and streaming video files. The server uses the standard `net/http` package and stores videos in the local file system.

## Getting Started

### Prerequisites
- Go 1.23.4 or later
- A working Go environment

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/mohamedfawas/learn-video-streaming-project.git
   cd learn-video-streaming-project
   ```

2. Install dependencies (if any):
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

- `POST /upload`: Upload a video file.
- `GET /list`: List all uploaded videos.
- `GET /stream/:filename`: Stream a video file.
