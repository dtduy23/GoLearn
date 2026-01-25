package main

import (
	"fmt"
	"log"
	"net/http"

	"spotify-clone/internal/config"
	"spotify-clone/internal/database"
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

	// Initialize repositories
	userRepo := user.NewUserRepository(db)

	// Initialize handlers
	userHandler := user.NewHandler(userRepo)

	// Setup routes
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetByID)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed:", err)
	}
}
