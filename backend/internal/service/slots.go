package service

import (
	"errors"
	"fmt"

	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/game"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SlotsService provides business logic for slots game
type SlotsService struct {
	db     *gorm.DB
	engine *game.SlotsEngine
}

// NewSlotsService creates a new slots service instance
func NewSlotsService() *SlotsService {
	return &SlotsService{
		db:     database.GetDB(),
		engine: game.NewSlotsEngine(),
	}
}

// SpinRequest represents a request to spin the slots
type SpinRequest struct {
	UserID uint    `json:"user_id" validate:"required"`
	Bet    float64 `json:"bet" validate:"required,gt=0"`
}

// SpinResponse represents the response from a slot spin
type SpinResponse struct {
	Result        *game.SlotResult `json:"result"`
	Bet           float64          `json:"bet"`
	Win           float64          `json:"win"`
	NewBalance    float64          `json:"new_balance"`
	TransactionID uint             `json:"transaction_id"`
	GameSessionID uint             `json:"game_session_id"`
}

// Spin performs a slot machine spin
func (s *SlotsService) Spin(req *SpinRequest) (*SpinResponse, error) {
	// Validate bet amount
	if req.Bet <= 0 {
		return nil, fmt.Errorf("bet must be greater than 0")
	}

	// Start database transaction
	var response *SpinResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get user with lock
		var user model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, req.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("failed to find user: %w", err)
		}

		// Check if user has enough balance
		if user.Balance < req.Bet {
			return fmt.Errorf("insufficient balance: have %.2f, need %.2f", user.Balance, req.Bet)
		}

		// Perform the spin
		result := s.engine.Spin(req.Bet)

		// Calculate new balance
		balanceChange := result.TotalWin - req.Bet
		newBalance := user.Balance + balanceChange

		// Update user balance
		if err := tx.Model(&user).Update("balance", newBalance).Error; err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create transaction record
		var transactionType model.TransactionType
		var transactionDesc string

		if result.TotalWin > 0 {
			transactionType = model.TransactionTypeGameWin
			transactionDesc = fmt.Sprintf("Slots win: bet %.2f, won %.2f (%.2fx)", req.Bet, result.TotalWin, result.Multiplier)
		} else {
			transactionType = model.TransactionTypeGameLoss
			transactionDesc = fmt.Sprintf("Slots loss: bet %.2f", req.Bet)
		}

		transaction := model.Transaction{
			UserID:      req.UserID,
			Type:        transactionType,
			Amount:      balanceChange,
			Description: transactionDesc,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Create game session record
		gameSession := model.GameSession{
			UserID:   req.UserID,
			GameType: model.GameTypeSlots,
			Bet:      req.Bet,
			Win:      result.TotalWin,
		}
		if err := tx.Create(&gameSession).Error; err != nil {
			return fmt.Errorf("failed to create game session: %w", err)
		}

		response = &SpinResponse{
			Result:        result,
			Bet:           req.Bet,
			Win:           result.TotalWin,
			NewBalance:    newBalance,
			TransactionID: transaction.ID,
			GameSessionID: gameSession.ID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetPayoutTable returns the payout table for display
func (s *SlotsService) GetPayoutTable() map[game.SlotSymbol]map[int]float64 {
	return game.GetPayoutTable()
}

// GetSymbols returns all available symbols
func (s *SlotsService) GetSymbols() []game.SlotSymbol {
	return game.GetAllSymbols()
}

// GetPaytableForAPI returns paytable in API-friendly format
func (s *SlotsService) GetPaytableForAPI() []game.PaytableEntry {
	return game.GetPaytableForAPI()
}
