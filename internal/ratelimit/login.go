package ratelimit

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginAttempt tracks failed login attempts for a username+IP combination
type LoginAttempt struct {
	FailedCount  int
	BlockedUntil time.Time
}

// LoginRateLimiter prevents brute force password attacks
type LoginRateLimiter struct {
	attempts      map[string]*LoginAttempt
	mu            sync.RWMutex
	maxAttempts   int
	blockDuration time.Duration
}

// NewLoginRateLimiter creates a new rate limiter
// maxAttempts: number of failed attempts before blocking
// blockDuration: how long to block after max attempts reached
func NewLoginRateLimiter(maxAttempts int, blockDuration time.Duration) *LoginRateLimiter {
	rl := &LoginRateLimiter{
		attempts:      make(map[string]*LoginAttempt),
		maxAttempts:   maxAttempts,
		blockDuration: blockDuration,
	}

	// Cleanup goroutine to remove old entries
	go rl.cleanup()

	return rl
}

// makeKey creates a unique key from username and IP
func makeKey(username, ip string) string {
	return fmt.Sprintf("%s:%s", username, ip)
}

// cleanup periodically removes expired blocks
func (rl *LoginRateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, attempt := range rl.attempts {
			// Remove if not blocked and no failed attempts
			if now.After(attempt.BlockedUntil) && attempt.FailedCount == 0 {
				delete(rl.attempts, key)
			}
			// Reset failed count after block expires
			if now.After(attempt.BlockedUntil) && attempt.FailedCount > 0 {
				attempt.FailedCount = 0
			}
		}
		rl.mu.Unlock()
	}
}

// IsBlocked checks if a username+IP combination is currently blocked
func (rl *LoginRateLimiter) IsBlocked(username, ip string) (bool, time.Duration) {
	key := makeKey(username, ip)

	rl.mu.RLock()
	defer rl.mu.RUnlock()

	attempt, exists := rl.attempts[key]
	if !exists {
		return false, 0
	}

	if time.Now().Before(attempt.BlockedUntil) {
		return true, time.Until(attempt.BlockedUntil)
	}

	return false, 0
}

// Middleware returns a Gin middleware that rate limits login attempts
// Note: This middleware only blocks by IP. For username+IP blocking,
// use IsBlocked() in the handler after parsing the request body.
func (rl *LoginRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Just pass through - the actual check will be done in handler
		// after we can parse the username from request body
		c.Next()
	}
}

// CheckAndBlock checks if blocked and returns appropriate response
// Returns true if request should be blocked
func (rl *LoginRateLimiter) CheckAndBlock(c *gin.Context, username string) bool {
	ip := c.ClientIP()
	blocked, remaining := rl.IsBlocked(username, ip)

	if blocked {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error":       "too many failed login attempts",
			"retry_after": remaining.Round(time.Second).String(),
			"message":     "Please wait before trying again",
		})
		return true
	}

	return false
}

// RecordFailedAttempt should be called when login fails
func (rl *LoginRateLimiter) RecordFailedAttempt(username, ip string) {
	key := makeKey(username, ip)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	attempt, exists := rl.attempts[key]
	if !exists {
		attempt = &LoginAttempt{}
		rl.attempts[key] = attempt
	}

	// Reset if previous block has expired
	if time.Now().After(attempt.BlockedUntil) {
		attempt.FailedCount = 0
	}

	attempt.FailedCount++

	if attempt.FailedCount >= rl.maxAttempts {
		attempt.BlockedUntil = time.Now().Add(rl.blockDuration)
	}
}

// RecordSuccessfulLogin resets the counter on successful login
func (rl *LoginRateLimiter) RecordSuccessfulLogin(username, ip string) {
	key := makeKey(username, ip)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.attempts, key)
}

// GetRemainingAttempts returns how many attempts are left
func (rl *LoginRateLimiter) GetRemainingAttempts(username, ip string) int {
	key := makeKey(username, ip)

	rl.mu.RLock()
	defer rl.mu.RUnlock()

	attempt, exists := rl.attempts[key]
	if !exists {
		return rl.maxAttempts
	}

	// If block expired, return full attempts
	if time.Now().After(attempt.BlockedUntil) {
		return rl.maxAttempts
	}

	remaining := rl.maxAttempts - attempt.FailedCount
	if remaining < 0 {
		return 0
	}
	return remaining
}
