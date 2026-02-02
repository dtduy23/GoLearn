package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"spotify-clone/internal/auth"
	"spotify-clone/internal/config"
	"spotify-clone/internal/database"
	"spotify-clone/internal/middleware"
	"spotify-clone/internal/ratelimit"
	"spotify-clone/internal/song"
	"spotify-clone/internal/user"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Initialize JWT service
	jwtService := auth.NewJWTService(auth.JWTConfig{
		SecretKey:          cfg.JWT.Secret,
		AccessTokenExpiry:  cfg.JWT.Expiry,
		RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
	})

	// Initialize repositories
	userRepo := user.NewUserRepository(db)
	songRepo := song.NewRepository(db)

	// Initialize services
	authService := auth.NewAuthService(userRepo, jwtService)

	// Initialize rate limiter: 5 failed attempts = block for 5 minutes
	loginRateLimiter := ratelimit.NewLoginRateLimiter(5, 5*time.Minute)

	// Initialize handlers
	authHandler := auth.NewHandler(authService, userRepo, loginRateLimiter)
	songHandler := song.NewHandler(songRepo)

	// Create auth middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// Setup Gin router
	r := gin.Default()

	// Serve static files (for audio streaming)
	r.Static("/static", "./web/static")

	// API routes
	api := r.Group("/api")
	{
		// Auth routes: /api/auth/...
		auth.RegisterRoutes(api, authHandler, authMiddleware, loginRateLimiter)

		// Song routes: /api/songs/...
		song.RegisterRoutes(api, songHandler)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Print all routes for reference
	log.Println("=== Available Routes ===")
	log.Println("POST   /api/auth/register    - Register new user")
	log.Println("POST   /api/auth/login       - Login")
	log.Println("POST   /api/auth/refresh     - Refresh token")
	log.Println("GET    /api/auth/me          - Get current user (protected)")
	log.Println("GET    /api/songs/:id        - Get song details")
	log.Println("GET    /api/songs/:id/stream - Stream song audio")
	log.Println("POST   /api/songs/upload     - Upload new song")
	log.Println("GET    /health               - Health check")
	log.Println("========================")

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server failed:", err)
	}
}
