package middleware

import (
	"context"
	"net/http"
	"strings"

	"spotify-clone/internal/auth"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userID"
	// EmailKey is the context key for user email
	EmailKey ContextKey = "email"
	// ClaimsKey is the context key for full claims
	ClaimsKey ContextKey = "claims"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "missing authorization header")
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeUnauthorized(w, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				if err == auth.ErrTokenExpired {
					writeUnauthorized(w, "token has expired")
					return
				}
				writeUnauthorized(w, "invalid token")
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			ctx = context.WithValue(ctx, ClaimsKey, claims)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware validates token if present, but doesn't require it
func OptionalAuthMiddleware(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			// If no auth header, continue without user context
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Try to parse token
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString := parts[1]
				claims, err := jwtService.ValidateAccessToken(tokenString)
				if err == nil {
					// Valid token, add to context
					ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
					ctx = context.WithValue(ctx, EmailKey, claims.Email)
					ctx = context.WithValue(ctx, ClaimsKey, claims)
					r = r.WithContext(ctx)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetEmail extracts email from context
func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

// GetClaims extracts full claims from context
func GetClaims(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*auth.Claims)
	return claims, ok
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
