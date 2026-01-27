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

	// Initialize services
	authService := auth.NewAuthService(userRepo, jwtService)

	// Initialize handlers
	userHandler := user.NewHandler(userRepo)

	// Initialize rate limiter: 5 failed attempts = block for 5 minutes
	loginRateLimiter := ratelimit.NewLoginRateLimiter(5, 5*time.Minute)
	authHandler := auth.NewHandler(authService, userRepo, loginRateLimiter)

	// Setup Gin router
	r := gin.Default()

	// Auth routes (public)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", loginRateLimiter.Middleware(), authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		// Protected route - requires valid JWT
		authGroup.GET("/me", middleware.AuthMiddleware(jwtService), authHandler.Me)
	}

	// API routes
	api := r.Group("/api")
	{
		// User routes
		api.GET("/users/:id", userHandler.GetByID)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server failed:", err)
	}
}
