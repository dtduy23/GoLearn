package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Env      string
	Database DatabaseConfig
	JWT      JWTConfig
	Static   StaticConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret             string
	Expiry             time.Duration
	RefreshTokenExpiry time.Duration
}

type StaticConfig struct {
	Path      string
	MusicPath string
}

func Load() (*Config, error) {
	// Load .env file
	godotenv.Load()

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))

	return &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "spotify_user"),
			Password: getEnv("DB_PASSWORD", "spotify_clone"),
			DBName:   getEnv("DB_NAME", "spotify_clone"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "secret"),
			Expiry:             jwtExpiry,
			RefreshTokenExpiry: refreshExpiry,
		},
		Static: StaticConfig{
			Path:      getEnv("STATIC_PATH", "./web/static"),
			MusicPath: getEnv("MUSIC_PATH", "./web/static/music"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
