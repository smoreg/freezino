package service

import (
	"errors"
	"fmt"

	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/game"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
)

// RouletteService handles roulette game business logic
type RouletteService struct {
	rouletteGame *game.RouletteGame
}

// NewRouletteService creates a new roulette service instance
func NewRouletteService() *RouletteService {
	return &RouletteService{
		rouletteGame: game.NewRouletteGame(),
	}
}

// PlaceBetRequest represents a request to place a bet
type PlaceBetRequest struct {
	UserID uint                `json:"user_id"`
	Bets   []model.RouletteBet `json:"bets"`
}

// PlaceBetResponse represents the response after placing a bet
type PlaceBetResponse struct {
	Number     int                 `json:"number"`
	Color      string              `json:"color"`
	TotalBet   float64             `json:"total_bet"`
	TotalWin   float64             `json:"total_win"`
	Profit     float64             `json:"profit"`
	NewBalance float64             `json:"new_balance"`
	Bets       []model.RouletteBet `json:"bets"`
}

// PlaceBet processes a roulette bet
func (s *RouletteService) PlaceBet(req PlaceBetRequest) (*PlaceBetResponse, error) {
	db := database.DB

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user
	var user model.User
	if err := tx.First(&user, req.UserID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Calculate total bet
	totalBet := 0.0
	for _, bet := range req.Bets {
		totalBet += bet.Amount
	}

	// Check if user has enough balance
	if user.Balance < totalBet {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance")
	}

	// Calculate result
	winningNumber, _, totalWin, err := s.rouletteGame.CalculateResult(req.Bets)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to calculate result: %w", err)
	}

	// Update user balance
	profit := totalWin - totalBet
	user.Balance += profit

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Encode bets
	betsJSON, err := game.EncodeBets(req.Bets)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to encode bets: %w", err)
	}

	// Save roulette result
	result := model.RouletteResult{
		UserID:   req.UserID,
		Number:   winningNumber,
		TotalBet: totalBet,
		TotalWin: totalWin,
		Bets:     betsJSON,
	}

	if err := tx.Create(&result).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save result: %w", err)
	}

	// Save game session
	gameSession := model.GameSession{
		UserID:   req.UserID,
		GameType: model.GameTypeRoulette,
		Bet:      totalBet,
		Win:      totalWin,
	}

	if err := tx.Create(&gameSession).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save game session: %w", err)
	}

	// Create transaction record
	transactionType := model.TransactionTypeGameLoss
	transactionAmount := -totalBet
	description := fmt.Sprintf("Roulette bet - number %d", winningNumber)

	if totalWin > 0 {
		transactionType = model.TransactionTypeGameWin
		transactionAmount = profit
		description = fmt.Sprintf("Roulette win - number %d (won $%.2f)", winningNumber, totalWin)
	}

	transaction := model.Transaction{
		UserID:      req.UserID,
		Type:        transactionType,
		Amount:      transactionAmount,
		Description: description,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return response
	return &PlaceBetResponse{
		Number:     winningNumber,
		Color:      s.rouletteGame.GetColor(winningNumber),
		TotalBet:   totalBet,
		TotalWin:   totalWin,
		Profit:     profit,
		NewBalance: user.Balance,
		Bets:       req.Bets,
	}, nil
}

// GetHistory retrieves recent roulette game history for a user
func (s *RouletteService) GetHistory(userID uint, limit int) ([]model.RouletteResult, error) {
	db := database.DB

	var results []model.RouletteResult
	if err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	return results, nil
}

// GetRecentNumbers retrieves recent winning numbers across all users
func (s *RouletteService) GetRecentNumbers(limit int) ([]int, error) {
	db := database.DB

	var results []model.RouletteResult
	if err := db.Select("number").
		Order("created_at DESC").
		Limit(limit).
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent numbers: %w", err)
	}

	numbers := make([]int, len(results))
	for i, result := range results {
		numbers[i] = result.Number
	}

	return numbers, nil
}
