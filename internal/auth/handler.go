package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"spotify-clone/internal/ratelimit"
	"spotify-clone/internal/user"
)

type Handler struct {
	authService AuthService
	userRepo    user.UserRepository
	rateLimiter *ratelimit.LoginRateLimiter
}

func NewHandler(authService AuthService, userRepo user.UserRepository, rateLimiter *ratelimit.LoginRateLimiter) *Handler {
	return &Handler{
		authService: authService,
		userRepo:    userRepo,
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

// GET /auth/me - Get current authenticated user
func (h *Handler) Me(c *gin.Context) {
	// Get userID from context (set by AuthMiddleware)
	// Using same key as middleware.UserIDKey = "userID"
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Fetch user from database
	foundUser, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        foundUser.ID.String(),
		Email:     foundUser.Email,
		Username:  foundUser.Username,
		Role:      foundUser.Role,
		CreatedAt: foundUser.CreatedAt,
	})
}
