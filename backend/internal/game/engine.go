package game

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"

	"freezino/backend/internal/model"

	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidBet          = errors.New("invalid bet amount")
	ErrUserNotFound        = errors.New("user not found")
	ErrGameNotFound        = errors.New("game not found")
	ErrInvalidGameResult   = errors.New("invalid game result")
)

// Engine represents the game engine that handles all game operations
type Engine struct {
	db     *gorm.DB
	config *GameConfig
}

// NewEngine creates a new game engine instance
func NewEngine(db *gorm.DB, config *GameConfig) *Engine {
	if config == nil {
		config = DefaultGameConfig()
	}
	return &Engine{
		db:     db,
		config: config,
	}
}

// CheckBalance verifies if user has sufficient balance for a bet
func (e *Engine) CheckBalance(userID uint, amount float64) (bool, error) {
	var user model.User
	if err := e.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrUserNotFound
		}
		return false, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user.Balance >= amount, nil
}

// GetUserBalance retrieves user's current balance
func (e *Engine) GetUserBalance(userID uint) (float64, error) {
	var user model.User
	if err := e.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user.Balance, nil
}

// ValidateBet checks if the bet amount is within allowed limits
func (e *Engine) ValidateBet(bet float64) error {
	if bet <= 0 {
		return ErrInvalidBet
	}
	if bet < e.config.MinBet {
		return fmt.Errorf("%w: minimum bet is %.2f", ErrInvalidBet, e.config.MinBet)
	}
	if bet > e.config.MaxBet {
		return fmt.Errorf("%w: maximum bet is %.2f", ErrInvalidBet, e.config.MaxBet)
	}
	return nil
}

// CreateGameSession creates a new game session record
func (e *Engine) CreateGameSession(userID uint, gameType model.GameType, bet float64) (*model.GameSession, error) {
	// Validate bet
	if err := e.ValidateBet(bet); err != nil {
		return nil, err
	}

	// Check balance
	hasBalance, err := e.CheckBalance(userID, bet)
	if err != nil {
		return nil, err
	}
	if !hasBalance {
		return nil, ErrInsufficientBalance
	}

	// Create game session
	session := &model.GameSession{
		UserID:   userID,
		GameType: gameType,
		Bet:      bet,
		Win:      0,
	}

	if err := e.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create game session: %w", err)
	}

	return session, nil
}

// UpdateBalance updates user's balance (can be positive or negative delta)
func (e *Engine) UpdateBalance(userID uint, delta float64) error {
	result := e.db.Model(&model.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", delta))

	if result.Error != nil {
		return fmt.Errorf("failed to update balance: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// CreateTransaction creates a transaction record
func (e *Engine) CreateTransaction(userID uint, txType model.TransactionType, amount float64, description string) error {
	transaction := &model.Transaction{
		UserID:      userID,
		Type:        txType,
		Amount:      amount,
		Description: description,
	}

	if err := e.db.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// ProcessGameResult processes the game result, updates balance and creates transactions
func (e *Engine) ProcessGameResult(session *model.GameSession, winAmount float64) error {
	// Start a database transaction to ensure atomicity
	return e.db.Transaction(func(tx *gorm.DB) error {
		// Deduct bet from balance
		if err := tx.Model(&model.User{}).Where("id = ?", session.UserID).
			Update("balance", gorm.Expr("balance - ?", session.Bet)).Error; err != nil {
			return fmt.Errorf("failed to deduct bet: %w", err)
		}

		// Create bet transaction (loss)
		betTx := &model.Transaction{
			UserID:      session.UserID,
			Type:        model.TransactionTypeGameLoss,
			Amount:      -session.Bet,
			Description: fmt.Sprintf("Bet on %s", session.GameType),
		}
		if err := tx.Create(betTx).Error; err != nil {
			return fmt.Errorf("failed to create bet transaction: %w", err)
		}

		// If there's a win, add it to balance
		if winAmount > 0 {
			if err := tx.Model(&model.User{}).Where("id = ?", session.UserID).
				Update("balance", gorm.Expr("balance + ?", winAmount)).Error; err != nil {
				return fmt.Errorf("failed to add winnings: %w", err)
			}

			// Create win transaction
			winTx := &model.Transaction{
				UserID:      session.UserID,
				Type:        model.TransactionTypeGameWin,
				Amount:      winAmount,
				Description: fmt.Sprintf("Win from %s", session.GameType),
			}
			if err := tx.Create(winTx).Error; err != nil {
				return fmt.Errorf("failed to create win transaction: %w", err)
			}
		}

		// Update game session with win amount
		if err := tx.Model(session).Update("win", winAmount).Error; err != nil {
			return fmt.Errorf("failed to update game session: %w", err)
		}

		return nil
	})
}

// SecureRandomInt generates a cryptographically secure random integer in range [0, max)
func SecureRandomInt(max int64) (int64, error) {
	if max <= 0 {
		return 0, errors.New("max must be positive")
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %w", err)
	}

	return nBig.Int64(), nil
}

// SecureRandomFloat generates a cryptographically secure random float64 in range [0.0, 1.0)
func SecureRandomFloat() (float64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert bytes to uint64
	randomUint := binary.BigEndian.Uint64(b[:])

	// Convert to float in range [0, 1)
	// Divide by max uint64 value to get a float
	randomFloat := float64(randomUint) / float64(math.MaxUint64)

	return randomFloat, nil
}

// SecureRandomFloatRange generates a cryptographically secure random float64 in range [min, max)
func SecureRandomFloatRange(min, max float64) (float64, error) {
	if min >= max {
		return 0, errors.New("min must be less than max")
	}

	randomFloat, err := SecureRandomFloat()
	if err != nil {
		return 0, err
	}

	// Scale to desired range
	return min + randomFloat*(max-min), nil
}

// ApplyHouseEdge applies the house edge to a theoretical payout
// This ensures the casino has a long-term advantage
func ApplyHouseEdge(theoreticalPayout float64, houseEdge float64) float64 {
	// Reduce payout by house edge percentage
	return theoreticalPayout * (1.0 - houseEdge)
}

// CalculateExpectedReturn calculates expected return for a bet with given probability and payout
func CalculateExpectedReturn(bet float64, winProbability float64, payoutMultiplier float64, houseEdge float64) float64 {
	// Theoretical return without house edge
	theoreticalReturn := bet * winProbability * payoutMultiplier

	// Apply house edge
	actualReturn := ApplyHouseEdge(theoreticalReturn, houseEdge)

	return actualReturn
}

// ShouldWin determines if the player should win based on probability and house edge
func ShouldWin(winProbability float64, houseEdge float64) (bool, error) {
	// Adjust win probability based on house edge
	adjustedProbability := winProbability * (1.0 - houseEdge)

	// Generate random float [0, 1)
	randomValue, err := SecureRandomFloat()
	if err != nil {
		return false, err
	}

	return randomValue < adjustedProbability, nil
}

// GetConfig returns the engine configuration
func (e *Engine) GetConfig() *GameConfig {
	return e.config
}

// GetDB returns the database connection
func (e *Engine) GetDB() *gorm.DB {
	return e.db
}
