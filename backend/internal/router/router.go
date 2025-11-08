package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/handler"
)

// Setup configures all application routes
func Setup(app *fiber.App) {
	// API group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", handler.HealthCheck)

	// Future routes will be added here
	// Example:
	// auth := api.Group("/auth")
	// auth.Get("/google", handler.GoogleAuth)
	// auth.Get("/google/callback", handler.GoogleAuthCallback)
}
