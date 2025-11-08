package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Port        string
	Environment string

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// JWT
	JWTSecret            string
	JWTAccessExpiration  string
	JWTRefreshExpiration string

	// Frontend URL
	FrontendURL string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Try to load .env file (optional in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		Port:        getEnv("PORT", "3000"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Google OAuth
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/api/auth/google/callback"),

		// JWT
		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTAccessExpiration:  getEnv("JWT_ACCESS_EXPIRATION", "15m"),
		JWTRefreshExpiration: getEnv("JWT_REFRESH_EXPIRATION", "7d"),

		// Frontend
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}

	return cfg
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
