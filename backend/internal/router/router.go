package router

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/auth"
	"github.com/smoreg/freezino/backend/internal/config"
	"github.com/smoreg/freezino/backend/internal/handler"
	"github.com/smoreg/freezino/backend/internal/middleware"
)

// Setup configures all application routes
func Setup(app *fiber.App, cfg *config.Config) {
	// API group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", handler.HealthCheck)

	// Auth routes
	authHandler := auth.NewHandler(cfg)
	authGroup := api.Group("/auth")
	{
		// OAuth routes
		authGroup.Get("/google", authHandler.GoogleLogin)
		authGroup.Get("/google/callback", authHandler.GoogleCallback)

		// Token refresh
		authGroup.Post("/refresh", authHandler.RefreshToken)

		// Protected routes (require authentication)
		authGroup.Get("/me", middleware.AuthMiddleware(cfg), authHandler.GetMe)
		authGroup.Post("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)
	}

	// User routes
	userHandler := handler.NewUserHandler()
	user := api.Group("/user")
	user.Get("/profile", userHandler.GetProfile)
	user.Patch("/profile", userHandler.UpdateProfile)
	user.Get("/balance", userHandler.GetBalance)
	user.Get("/stats", userHandler.GetStats)
	user.Get("/transactions", userHandler.GetTransactions)
	user.Get("/items", userHandler.GetUserItems)

	// Work routes
	workHandler := handler.NewWorkHandler()
	work := api.Group("/work")
	work.Post("/start", workHandler.StartWork)
	work.Get("/status", workHandler.GetStatus)
	work.Post("/complete", workHandler.CompleteWork)
	work.Get("/history", workHandler.GetHistory)

	// Stats routes
	statsHandler, err := handler.NewStatsHandler()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize stats handler: %v", err))
	}
	stats := api.Group("/stats")
	stats.Get("/countries", statsHandler.GetCountries)
	stats.Get("/countries/:code", statsHandler.GetCountryByCode)

	// Contact routes
	contactHandler := handler.NewContactHandler()
	api.Post("/contact", contactHandler.SubmitMessage)

	// Game routes
	games := api.Group("/games")

	// Roulette routes
	rouletteHandler := handler.NewRouletteHandler()
	roulette := games.Group("/roulette")
	roulette.Post("/bet", rouletteHandler.PlaceBet)
	roulette.Get("/history", rouletteHandler.GetHistory)
	roulette.Get("/recent", rouletteHandler.GetRecentNumbers)

	// Future routes will be added here
}
