package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"emcoded/emcoded.calendar.be/internal/server"
)

func GetImages(w http.ResponseWriter, r *http.Request) {
	server.EnableCors(w)

	// Allow overriding the base path via APP_BASE_PATH. In the container we set
	// APP_BASE_PATH=/root so files baked into the image are found. Default to
	// the repo-relative ./root for local dev.
	base := os.Getenv("APP_BASE_PATH")
	if base == "" {
		base = filepath.Join("root")
	}

	dir := filepath.Join(base, "media", "img")
	entries, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "failed to read images directory", http.StatusInternalServerError)
		return
	}

	var images []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		// Ensure there's an extension and only accept jpg/jpeg/png for now
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}

		// Build response path relative to /root for consistency with earlier
		// behavior when APP_BASE_PATH is not set to an absolute /root.
		// If base is /root, this will produce paths like /root/media/img/file.jpg
		respPath := filepath.ToSlash(filepath.Join(base, "media", "img", name))
		images = append(images, respPath)
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make([]ImageResponse, len(images))
	for i, img := range images {
		resp[i] = ImageResponse{ImagePath: img}
	}

	status, body := returnJsonResponse(resp, http.StatusOK)
	w.WriteHeader(status)
	w.Write(body)
}
