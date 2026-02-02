package auth

import (
	"github.com/gin-gonic/gin"

	"spotify-clone/internal/ratelimit"
)

// RegisterRoutes registers all auth routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc, rateLimiter *ratelimit.LoginRateLimiter) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", rateLimiter.Middleware(), h.Login)
		authGroup.POST("/refresh", h.RefreshToken)
		// Protected route - requires valid JWT
		authGroup.GET("/me", authMiddleware, h.Me)
	}
}
