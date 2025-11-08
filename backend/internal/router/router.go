package router

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/auth"
	"github.com/smoreg/freezino/backend/internal/config"
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/handler"
	games "github.com/smoreg/freezino/backend/internal/handler/games"
	"github.com/smoreg/freezino/backend/internal/middleware"
	"github.com/smoreg/freezino/backend/internal/service"
)

// Setup configures all application routes
func Setup(app *fiber.App, cfg *config.Config) {
	// API group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", handler.HealthCheck)

	// Auth routes
	authHandler := auth.NewHandler(cfg)

	// Local auth (username/password)
	db := database.GetDB()
	authService := service.NewAuthService(db)
	localAuthHandler := handler.NewAuthHandler(authService, cfg)

	authGroup := api.Group("/auth")
	{
		// Local auth routes (username/password)
		authGroup.Post("/register", localAuthHandler.Register)
		authGroup.Post("/login", localAuthHandler.Login)

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

	// Shop routes
	shopHandler := handler.NewShopHandler()
	shop := api.Group("/shop")
	shop.Get("/items", shopHandler.GetItems)
	shop.Post("/buy/:itemId", shopHandler.BuyItem)
	shop.Post("/sell/:userItemId", shopHandler.SellItem)
	shop.Get("/my-items", shopHandler.GetMyItems)
	shop.Post("/equip/:userItemId", shopHandler.EquipItem)

	// Game routes
	gamesGroup := api.Group("/games")

	// Roulette routes
	rouletteHandler := handler.NewRouletteHandler()
	roulette := gamesGroup.Group("/roulette")
	roulette.Post("/bet", rouletteHandler.PlaceBet)
	roulette.Get("/history", rouletteHandler.GetHistory)
	roulette.Get("/recent", rouletteHandler.GetRecentNumbers)

	// Slots routes
	slotsHandler := handler.NewSlotsHandler()
	slots := gamesGroup.Group("/slots")
	slots.Post("/spin", slotsHandler.Spin)
	slots.Get("/payouts", slotsHandler.GetPayoutTable)

	// Crash game
	crashHandler := games.NewCrashHandler()
	crash := gamesGroup.Group("/crash")
	crash.Post("/bet", crashHandler.PlaceBet)

	// Hi-Lo game
	hiloHandler := games.NewHiLoHandler()
	hilo := gamesGroup.Group("/hilo")
	hilo.Post("/bet", hiloHandler.PlaceBet)

	// Wheel game
	wheelHandler := games.NewWheelHandler()
	wheel := gamesGroup.Group("/wheel")
	wheel.Post("/spin", wheelHandler.Spin)

	// Game history routes
	gameHistoryHandler := handler.NewGameHistoryHandler()
	gamesGroup.Get("/history", gameHistoryHandler.GetHistory)
	gamesGroup.Get("/stats", gameHistoryHandler.GetStats)

	// Game WebSocket routes
	gameHandler := handler.NewGameHandler(db)

	// WebSocket upgrade middleware and routes
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/blackjack", websocket.New(gameHandler.BlackjackWebSocket))

	// Future routes will be added here
}
