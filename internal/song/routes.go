package song

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all song routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	songGroup := rg.Group("/songs")
	{
		songGroup.GET("/:id", h.GetSong)
		songGroup.GET("/:id/stream", h.StreamSong)
		songGroup.POST("/upload", h.UploadSong)
	}
}
