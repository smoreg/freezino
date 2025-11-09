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

	// User routes (protected)
	userHandler := handler.NewUserHandler()
	user := api.Group("/user", middleware.AuthMiddleware(cfg))
	user.Get("/profile", userHandler.GetProfile)
	user.Patch("/profile", userHandler.UpdateProfile)
	user.Get("/balance", userHandler.GetBalance)
	user.Get("/stats", userHandler.GetStats)
	user.Get("/transactions", userHandler.GetTransactions)
	user.Get("/items", userHandler.GetUserItems)

	// Work routes (protected)
	workHandler := handler.NewWorkHandler()
	work := api.Group("/work", middleware.AuthMiddleware(cfg))
	work.Post("/start", workHandler.StartWork)
	work.Get("/status", workHandler.GetStatus)
	work.Post("/complete", workHandler.CompleteWork)
	work.Get("/history", workHandler.GetHistory)
	work.Get("/jobs", workHandler.GetAvailableJobs)
	work.Post("/skip-jail", workHandler.SkipJailTime)

	// Stats routes
	statsHandler, err := handler.NewStatsHandler()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize stats handler: %v", err))
	}
	stats := api.Group("/stats")
	stats.Get("/countries", statsHandler.GetCountries)
	stats.Get("/countries/:code", statsHandler.GetCountryByCode)

	// Casino statistics routes
	casinoStatsHandler := handler.NewCasinoStatsHandler()
	casino := api.Group("/casino")
	casino.Get("/stats", casinoStatsHandler.GetCasinoStats) // Public - anyone can view casino stats

	// Contact routes
	contactHandler := handler.NewContactHandler()
	api.Post("/contact", contactHandler.SubmitMessage)

	// Dev/testing routes
	devHandler := handler.NewDevHandler()
	dev := api.Group("/dev")
	// Public dev routes (no auth needed)
	dev.Post("/seed", devHandler.SeedDatabase) // Public - needs to create test users
	// Protected dev routes (require auth)
	devProtected := dev.Group("", middleware.AuthMiddleware(cfg))
	devProtected.Post("/add-money", devHandler.AddMoney)
	devProtected.Post("/reset-balance", devHandler.ResetBalance)

	// Shop routes
	shopHandler := handler.NewShopHandler()
	shop := api.Group("/shop")
	shop.Get("/items", shopHandler.GetItems) // Public - can browse items without auth
	// Protected shop routes (require auth)
	shopProtected := shop.Group("", middleware.AuthMiddleware(cfg))
	shopProtected.Post("/buy/:itemId", shopHandler.BuyItem)
	shopProtected.Post("/sell/:userItemId", shopHandler.SellItem)
	shopProtected.Get("/my-items", shopHandler.GetMyItems)
	shopProtected.Post("/equip/:userItemId", shopHandler.EquipItem)

	// Game routes (protected)
	gamesGroup := api.Group("/games", middleware.AuthMiddleware(cfg))

	// Roulette routes
	rouletteHandler := handler.NewRouletteHandler()
	roulette := gamesGroup.Group("/roulette")
	roulette.Post("/bet", rouletteHandler.PlaceBet)
	roulette.Get("/history", rouletteHandler.GetHistory)
	roulette.Get("/recent", rouletteHandler.GetRecentNumbers) // Public - can view recent numbers

	// Slots routes
	slotsHandler := handler.NewSlotsHandler()
	slots := gamesGroup.Group("/slots")
	slots.Post("/spin", slotsHandler.Spin)
	slots.Get("/payouts", slotsHandler.GetPayoutTable) // Public - can view payout table

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

	// Loan routes (protected)
	loanHandler := handler.NewLoanHandler()
	loans := api.Group("/loans", middleware.AuthMiddleware(cfg))
	loans.Get("/summary", loanHandler.GetLoanSummary)
	loans.Get("", loanHandler.GetUserLoans)
	loans.Post("/take", loanHandler.TakeLoan)
	loans.Post("/repay/:loanId", loanHandler.RepayLoan)
	loans.Get("/bankruptcy-check", loanHandler.CheckBankruptcy)

	// Future routes will be added here
}
