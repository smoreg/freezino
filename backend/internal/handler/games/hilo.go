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

// HiLoHandler handles hi-lo game HTTP requests
type HiLoHandler struct{}

// NewHiLoHandler creates a new hi-lo handler instance
func NewHiLoHandler() *HiLoHandler {
	return &HiLoHandler{}
}

// HiLoBetRequest represents a hi-lo bet request
type HiLoBetRequest struct {
	UserID    uint    `json:"user_id"`
	BetAmount float64 `json:"bet_amount"`
	Guess     string  `json:"guess"` // "higher" or "lower"
}

// HiLoBetResponse represents a hi-lo bet response
type HiLoBetResponse struct {
	Success      bool    `json:"success"`
	CurrentCard  int     `json:"current_card"`  // The card shown (1-13)
	NextCard     int     `json:"next_card"`     // The next card drawn
	BetAmount    float64 `json:"bet_amount"`
	WinAmount    float64 `json:"win_amount"`
	NewBalance   float64 `json:"new_balance"`
	Won          bool    `json:"won"`
	CurrentSuit  string  `json:"current_suit"`  // Hearts, Diamonds, Clubs, Spades
	NextSuit     string  `json:"next_suit"`
}

// PlaceBet handles POST /api/games/hilo/bet
// @Summary Place a hi-lo game bet
// @Description Guess if next card is higher or lower
// @Tags games
// @Accept json
// @Produce json
// @Param request body HiLoBetRequest true "Bet request"
// @Success 200 {object} HiLoBetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/hilo/bet [post]
func (h *HiLoHandler) PlaceBet(c *fiber.Ctx) error {
	var req HiLoBetRequest
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

	// Validate guess
	if req.Guess != "higher" && req.Guess != "lower" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "guess must be 'higher' or 'lower'",
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

	// Generate cards
	rand.Seed(time.Now().UnixNano())
	currentCard := rand.Intn(13) + 1 // 1-13 (Ace to King)
	nextCard := rand.Intn(13) + 1

	currentSuit := getRandomSuit()
	nextSuit := getRandomSuit()

	// Determine if player won
	won := false
	if req.Guess == "higher" && nextCard > currentCard {
		won = true
	} else if req.Guess == "lower" && nextCard < currentCard {
		won = true
	}

	// If cards are equal, it's a push (return bet)
	isPush := currentCard == nextCard

	// Calculate winnings (2x multiplier for correct guess)
	winAmount := 0.0
	netChange := -req.BetAmount // Default to loss

	if isPush {
		// Push - return the bet
		winAmount = req.BetAmount
		netChange = 0
	} else if won {
		// Win - 2x payout
		winAmount = req.BetAmount * 2.0
		netChange = winAmount - req.BetAmount
	}

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
		GameType: model.GameTypeHiLo,
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
	description := "Hi-Lo game - Loss"

	if isPush {
		transactionType = model.TransactionType("game_push")
		description = "Hi-Lo game - Push (tie)"
	} else if won {
		transactionType = model.TransactionTypeGameWin
		description = "Hi-Lo game - Win"
	}

	transaction := model.Transaction{
		UserID:       req.UserID,
		Type:         transactionType,
		Amount:       math.Abs(netChange),
		BalanceAfter: newBalance,
		Description:  description + " (" + strconv.Itoa(currentCard) + " vs " + strconv.Itoa(nextCard) + ")",
	}

	if err := db.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to create transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(HiLoBetResponse{
		Success:     true,
		CurrentCard: currentCard,
		NextCard:    nextCard,
		CurrentSuit: currentSuit,
		NextSuit:    nextSuit,
		BetAmount:   req.BetAmount,
		WinAmount:   winAmount,
		NewBalance:  newBalance,
		Won:         won || isPush,
	})
}

// getRandomSuit returns a random card suit
func getRandomSuit() string {
	suits := []string{"hearts", "diamonds", "clubs", "spades"}
	return suits[rand.Intn(len(suits))]
}
