package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"spotify-clone/internal/ratelimit"
	"spotify-clone/internal/user"
)

type Handler struct {
	authService AuthService
	rateLimiter *ratelimit.LoginRateLimiter
}

func NewHandler(authService AuthService, rateLimiter *ratelimit.LoginRateLimiter) *Handler {
	return &Handler{
		authService: authService,
		rateLimiter: rateLimiter,
	}
}

// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate
	if req.Email == "" || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email, username and password are required"})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, user.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		if errors.Is(err, user.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	ip := c.ClientIP()

	// Check if this username+IP is blocked
	if h.rateLimiter.CheckAndBlock(c, req.Username) {
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			// Record failed attempt with username + IP
			h.rateLimiter.RecordFailedAttempt(req.Username, ip)
			remaining := h.rateLimiter.GetRemainingAttempts(req.Username, ip)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":              "invalid username or password",
				"attempts_remaining": remaining,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	// Reset on successful login
	h.rateLimiter.RecordSuccessfulLogin(req.Username, ip)

	c.JSON(http.StatusOK, resp)
}

// POST /auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
