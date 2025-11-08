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

	// User routes
	userHandler := handler.NewUserHandler()
	user := api.Group("/user")
	user.Get("/profile", userHandler.GetProfile)
	user.Patch("/profile", userHandler.UpdateProfile)
	user.Get("/balance", userHandler.GetBalance)
	user.Get("/stats", userHandler.GetStats)
	user.Get("/transactions", userHandler.GetTransactions)
	user.Get("/items", userHandler.GetUserItems)

	// Future routes will be added here
	// Example:
	// auth := api.Group("/auth")
	// auth.Get("/google", handler.GoogleAuth)
	// auth.Get("/google/callback", handler.GoogleAuthCallback)
}
