// internal/song/handler.go
package song

import (
	"fmt"
	"net/http"
	"os"
	"spotify-clone/pkg/audioduration"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for songs
type Handler struct {
	repo *Repository
}

// NewHandler creates a new song handler
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// StreamSong streams audio file for a song
func (h *Handler) StreamSong(c *gin.Context) {
	songID := c.Param("id")

	// Lấy song từ DB
	song, err := h.repo.GetByID(c.Request.Context(), songID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	// Mở file audio
	file, err := os.Open(song.FileURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	// Lấy thông tin file
	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read file info"})
		return
	}

	// Set headers cho streaming
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Accept-Ranges", "bytes")

	// Gin's http.ServeContent tự động xử lý Range requests cho seek/skip
	http.ServeContent(c.Writer, c.Request, song.Title, stat.ModTime(), file)
}

// GetSong returns song details as JSON
func (h *Handler) GetSong(c *gin.Context) {
	songID := c.Param("id")

	song, err := h.repo.GetByID(c.Request.Context(), songID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// Allowed audio MIME types
var allowedAudioTypes = map[string]bool{
	"audio/mpeg": true, // MP3
	"audio/mp3":  true, // MP3 (alternative)
	// "audio/wav":   true, // WAV
	// "audio/wave":  true, // WAV (alternative)
	"audio/ogg":  true, // OGG
	"audio/flac": true, // FLAC
	// "audio/aac":   true, // AAC
	// "audio/x-m4a": true, // M4A
	// "audio/mp4":   true, // M4A (alternative)
}

// UploadSong handles audio file upload with validation and database save
func (h *Handler) UploadSong(c *gin.Context) {
	// 1. Bind form data to request struct
	var req SongUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data: " + err.Error()})
		return
	}

	// 2. Get audio file from form
	fileHeader, err := c.FormFile("audio_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Audio file is required"})
		return
	}

	// 3. Validate file type by reading magic bytes (more reliable than Content-Type header)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read file"})
		return
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read file content"})
		return
	}

	// Detect audio format using magic bytes
	audioType := -1

	// Check for MP3 (ID3 tag or MP3 frame sync)
	if (len(buffer) >= 3 && buffer[0] == 0x49 && buffer[1] == 0x44 && buffer[2] == 0x33) || // ID3v2 tag
		(len(buffer) >= 2 && buffer[0] == 0xFF && (buffer[1]&0xE0) == 0xE0) { // MP3 frame sync
		audioType = audioduration.TypeMp3
	}

	// Check for OGG (magic: OggS)
	if len(buffer) >= 4 && buffer[0] == 0x4F && buffer[1] == 0x67 && buffer[2] == 0x67 && buffer[3] == 0x53 {
		audioType = audioduration.TypeOgg
	}

	// Check for FLAC (magic: fLaC)
	if len(buffer) >= 4 && buffer[0] == 0x66 && buffer[1] == 0x4C && buffer[2] == 0x61 && buffer[3] == 0x43 {
		audioType = audioduration.TypeFlac
	}

	if audioType == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Allowed: MP3, OGG, FLAC",
		})
		return
	}

	// Reset file position for later use
	file.Seek(0, 0)

	// 4. Generate unique filename with UUID
	fileExt := getFileExtension(fileHeader.Filename)
	newFileName := fmt.Sprintf("%s%s", generateUUID(), fileExt)
	filePath := fmt.Sprintf("./assets/audio/%s", newFileName)

	// 5. Save file to disk
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 6. Get audio duration
	audioFile, err := os.Open(filePath)
	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read audio file"})
		return
	}
	defer audioFile.Close()

	duration, err := audioduration.Duration(audioFile, audioType)

	// 7. Create song in database
	songID := generateUUID()
	var albumIDPtr *string
	if req.AlbumID != "" {
		albumIDPtr = &req.AlbumID
	}

	input := CreateSongInput{
		Song: Song{
			ID:        songID,
			Title:     req.Title,
			Duration:  int(duration),
			FileURL:   filePath,
			PlayCount: 0,
			CreatedAt: getCurrentTime(),
		},
		AlbumID:   albumIDPtr,
		ArtistIDs: req.ArtistIDs,
		GenreIDs:  req.GenreIDs,
	}

	if err := h.repo.CreateSong(c.Request.Context(), input); err != nil {
		// Rollback: delete uploaded file
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database: " + err.Error()})
		return
	}

	// 8. Return success response
	c.JSON(http.StatusCreated, SongUploadResponse{
		ID:       songID,
		Title:    req.Title,
		Duration: int(duration),
		Message:  "Song uploaded successfully",
	})
}
