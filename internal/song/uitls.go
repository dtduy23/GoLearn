package song

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// getFileExtension returns the file extension with dot (e.g., ".mp3")
func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}

// generateUUID generates a new UUID v7 string
func generateUUID() string {
	return uuid.Must(uuid.NewV7()).String()
}

// getCurrentTime returns current time for CreatedAt fields
func getCurrentTime() time.Time {
	return time.Now()
}
