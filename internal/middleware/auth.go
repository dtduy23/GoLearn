package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"spotify-clone/internal/auth"
)

const (
	// UserIDKey is the context key for user ID
	UserIDKey = "userID"
	// EmailKey is the context key for user email
	EmailKey = "email"
	// RoleKey is the context key for user role
	RoleKey = "role"
	// ClaimsKey is the context key for full claims
	ClaimsKey = "claims"
)

// AuthMiddleware creates a Gin middleware that validates JWT tokens
func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			if err == auth.ErrTokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Add claims to context
		c.Set(UserIDKey, claims.UserID)
		c.Set(EmailKey, claims.Email)
		c.Set(RoleKey, claims.Role)
		c.Set(ClaimsKey, claims)

		// Call next handler
		c.Next()
	}
}

// OptionalAuthMiddleware validates token if present, but doesn't require it
func OptionalAuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// If no auth header, continue without user context
		if authHeader == "" {
			c.Next()
			return
		}

		// Try to parse token
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString := parts[1]
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err == nil {
				// Valid token, add to context
				c.Set(UserIDKey, claims.UserID)
				c.Set(EmailKey, claims.Email)
				c.Set(RoleKey, claims.Role)
				c.Set(ClaimsKey, claims)
			}
		}

		c.Next()
	}
}

// GetUserID extracts user ID from Gin context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetEmail extracts email from Gin context
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(EmailKey)
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// GetClaims extracts full claims from Gin context
func GetClaims(c *gin.Context) (*auth.Claims, bool) {
	claims, exists := c.Get(ClaimsKey)
	if !exists {
		return nil, false
	}
	cl, ok := claims.(*auth.Claims)
	return cl, ok
}

// PremiumMiddleware requires user to have premium or admin role
func PremiumMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := GetClaims(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		if claims.Role != "premium" && claims.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "premium subscription required"})
			return
		}
		c.Next()
	}
}

// AdminMiddleware requires user to have admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := GetClaims(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		if claims.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

// GetRole extracts role from Gin context
func GetRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(RoleKey)
	if !exists {
		return "", false
	}
	r, ok := role.(string)
	return r, ok
}
