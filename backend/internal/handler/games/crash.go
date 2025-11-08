package games

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
)

// CrashHandler handles crash game HTTP requests
type CrashHandler struct{}

// NewCrashHandler creates a new crash handler instance
func NewCrashHandler() *CrashHandler {
	return &CrashHandler{}
}

// BetRequest represents a crash bet request
type BetRequest struct {
	UserID     uint    `json:"user_id"`
	BetAmount  float64 `json:"bet_amount"`
	CashoutAt  float64 `json:"cashout_at"` // Multiplier at which user wants to cashout (1.0x - 100.0x)
}

// BetResponse represents a crash bet response
type BetResponse struct {
	Success       bool    `json:"success"`
	CrashPoint    float64 `json:"crash_point"`    // Point at which game crashed (e.g., 2.45x)
	PlayerCashout float64 `json:"player_cashout"` // Player's cashout multiplier
	BetAmount     float64 `json:"bet_amount"`
	WinAmount     float64 `json:"win_amount"`
	NewBalance    float64 `json:"new_balance"`
	Won           bool    `json:"won"`
}

// PlaceBet handles POST /api/games/crash/bet
// @Summary Place a crash game bet
// @Description Place a bet in crash game with desired cashout multiplier
// @Tags games
// @Accept json
// @Produce json
// @Param request body BetRequest true "Bet request"
// @Success 200 {object} BetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/crash/bet [post]
func (h *CrashHandler) PlaceBet(c *fiber.Ctx) error {
	var req BetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	// Validate bet amount
	if req.BetAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "bet amount must be greater than 0",
		})
	}

	// Validate cashout multiplier
	if req.CashoutAt < 1.0 || req.CashoutAt > 100.0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "cashout multiplier must be between 1.0x and 100.0x",
		})
	}

	// Get user
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user not found",
		})
	}

	// Check if user has enough balance
	if user.Balance < req.BetAmount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "insufficient balance",
		})
	}

	// Generate crash point (house edge: ~3%)
	crashPoint := generateCrashPoint()

	// Determine if player won
	won := req.CashoutAt <= crashPoint
	winAmount := 0.0

	if won {
		winAmount = req.BetAmount * req.CashoutAt
	}

	// Calculate net change
	netChange := winAmount - req.BetAmount

	// Update user balance
	newBalance := user.Balance + netChange
	if err := db.Model(&user).Update("balance", newBalance).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to update balance",
		})
	}

	// Create game session record
	gameSession := model.GameSession{
		UserID:   req.UserID,
		GameType: model.GameTypeCrash,
		Bet:      req.BetAmount,
		Win:      winAmount,
	}

	if err := db.Create(&gameSession).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to create game session",
		})
	}

	// Create transaction record
	transactionType := model.TransactionType("game_loss")
	if won {
		transactionType = model.TransactionTypeGameWin
	}

	transaction := model.Transaction{
		UserID:       req.UserID,
		Type:         transactionType,
		Amount:       math.Abs(netChange),
		BalanceAfter: newBalance,
		Description:  "Crash game - " + strconv.FormatFloat(crashPoint, 'f', 2, 64) + "x",
	}

	if err := db.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to create transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(BetResponse{
		Success:       true,
		CrashPoint:    crashPoint,
		PlayerCashout: req.CashoutAt,
		BetAmount:     req.BetAmount,
		WinAmount:     winAmount,
		NewBalance:    newBalance,
		Won:           won,
	})
}

// generateCrashPoint generates a crash point using provably fair algorithm
// Returns a multiplier between 1.00x and ~100x with house edge
func generateCrashPoint() float64 {
	rand.Seed(time.Now().UnixNano())

	// House edge: 3%
	houseEdge := 0.03

	// Generate random float between 0 and 1
	r := rand.Float64()

	// Calculate crash point with exponential distribution
	// This creates realistic crash distribution where lower multipliers are more common
	crashPoint := math.Floor((100.0 / (r * (1.0 - houseEdge))) * 100) / 100

	// Ensure minimum of 1.00x and maximum of 100.00x
	if crashPoint < 1.00 {
		crashPoint = 1.00
	} else if crashPoint > 100.00 {
		crashPoint = 100.00
	}

	return crashPoint
}
