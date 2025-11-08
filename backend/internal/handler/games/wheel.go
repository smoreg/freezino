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

// WheelHandler handles wheel game HTTP requests
type WheelHandler struct{}

// NewWheelHandler creates a new wheel handler instance
func NewWheelHandler() *WheelHandler {
	return &WheelHandler{}
}

// WheelSpinRequest represents a wheel spin request
type WheelSpinRequest struct {
	UserID    uint    `json:"user_id"`
	BetAmount float64 `json:"bet_amount"`
}

// WheelSegment represents a wheel segment with multiplier and probability
type WheelSegment struct {
	Multiplier float64 `json:"multiplier"`
	Color      string  `json:"color"`
	Weight     int     `json:"weight"` // Used for probability (higher = more common)
}

// WheelSpinResponse represents a wheel spin response
type WheelSpinResponse struct {
	Success    bool    `json:"success"`
	Segment    int     `json:"segment"`    // Index of winning segment (0-9)
	Multiplier float64 `json:"multiplier"` // Winning multiplier
	Color      string  `json:"color"`      // Winning color
	BetAmount  float64 `json:"bet_amount"`
	WinAmount  float64 `json:"win_amount"`
	NewBalance float64 `json:"new_balance"`
}

// Define wheel segments with multipliers and weights
var wheelSegments = []WheelSegment{
	{Multiplier: 1.2, Color: "blue", Weight: 25},    // Common
	{Multiplier: 1.5, Color: "green", Weight: 20},   // Common
	{Multiplier: 2.0, Color: "yellow", Weight: 15},  // Medium
	{Multiplier: 3.0, Color: "orange", Weight: 12},  // Medium
	{Multiplier: 5.0, Color: "red", Weight: 10},     // Rare
	{Multiplier: 10.0, Color: "purple", Weight: 8},  // Very rare
	{Multiplier: 20.0, Color: "gold", Weight: 5},    // Super rare
	{Multiplier: 50.0, Color: "rainbow", Weight: 3}, // Ultra rare
	{Multiplier: 0.0, Color: "black", Weight: 1},    // Lose all (extremely rare)
	{Multiplier: 100.0, Color: "diamond", Weight: 1}, // Jackpot (extremely rare)
}

// Spin handles POST /api/games/wheel/spin
// @Summary Spin the wheel
// @Description Spin the wheel of fortune to win multipliers
// @Tags games
// @Accept json
// @Produce json
// @Param request body WheelSpinRequest true "Spin request"
// @Success 200 {object} WheelSpinResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/wheel/spin [post]
func (h *WheelHandler) Spin(c *fiber.Ctx) error {
	var req WheelSpinRequest
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

	// Spin the wheel (weighted random)
	winningSegmentIndex := spinWheel()
	winningSegment := wheelSegments[winningSegmentIndex]

	// Calculate winnings
	winAmount := req.BetAmount * winningSegment.Multiplier
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
		GameType: model.GameTypeWheel,
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
	description := "Wheel of Fortune - "

	if winningSegment.Multiplier == 0 {
		description += "Lost all"
	} else if netChange > 0 {
		transactionType = model.TransactionTypeGameWin
		description += strconv.FormatFloat(winningSegment.Multiplier, 'f', 1, 64) + "x win"
	} else if netChange == 0 {
		transactionType = model.TransactionType("game_push")
		description += "Break even"
	} else {
		description += strconv.FormatFloat(winningSegment.Multiplier, 'f', 1, 64) + "x (loss)"
	}

	transaction := model.Transaction{
		UserID:       req.UserID,
		Type:         transactionType,
		Amount:       math.Abs(netChange),
		BalanceAfter: newBalance,
		Description:  description,
	}

	if err := db.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to create transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(WheelSpinResponse{
		Success:    true,
		Segment:    winningSegmentIndex,
		Multiplier: winningSegment.Multiplier,
		Color:      winningSegment.Color,
		BetAmount:  req.BetAmount,
		WinAmount:  winAmount,
		NewBalance: newBalance,
	})
}

// spinWheel performs weighted random selection of wheel segment
func spinWheel() int {
	rand.Seed(time.Now().UnixNano())

	// Calculate total weight
	totalWeight := 0
	for _, segment := range wheelSegments {
		totalWeight += segment.Weight
	}

	// Generate random number
	randomValue := rand.Intn(totalWeight)

	// Find winning segment
	currentWeight := 0
	for i, segment := range wheelSegments {
		currentWeight += segment.Weight
		if randomValue < currentWeight {
			return i
		}
	}

	// Fallback (should never reach here)
	return 0
}
