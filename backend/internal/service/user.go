package service

import (
	"fmt"

	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
)

// UserService provides business logic for user operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service instance
func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

// ProfileResponse represents user profile data
type ProfileResponse struct {
	ID        uint      `json:"id"`
	GoogleID  string    `json:"google_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Balance   float64   `json:"balance"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name   *string `json:"name,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}

// BalanceResponse represents user balance data
type BalanceResponse struct {
	UserID  uint    `json:"user_id"`
	Balance float64 `json:"balance"`
}

// StatsResponse represents user statistics
type StatsResponse struct {
	UserID          uint    `json:"user_id"`
	TotalWorkTime   int     `json:"total_work_time"`   // in seconds
	TotalEarned     float64 `json:"total_earned"`
	TotalGameSessions int   `json:"total_game_sessions"`
	TotalBet        float64 `json:"total_bet"`
	TotalWon        float64 `json:"total_won"`
	TotalLost       float64 `json:"total_lost"`
	NetProfit       float64 `json:"net_profit"` // total won - total bet
	FavoriteGame    string  `json:"favorite_game,omitempty"`
	BiggestWin      float64 `json:"biggest_win"`
	BiggestLoss     float64 `json:"biggest_loss"`
}

// GetProfile retrieves user profile by ID
func (s *UserService) GetProfile(userID uint) (*ProfileResponse, error) {
	var user model.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &ProfileResponse{
		ID:        user.ID,
		GoogleID:  user.GoogleID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uint, req UpdateProfileRequest) (*ProfileResponse, error) {
	var user model.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}

	if len(updates) > 0 {
		if err := s.db.Model(&user).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update user profile: %w", err)
		}
	}

	// Reload user to get updated data
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	return &ProfileResponse{
		ID:        user.ID,
		GoogleID:  user.GoogleID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetBalance retrieves user balance
func (s *UserService) GetBalance(userID uint) (*BalanceResponse, error) {
	var user model.User

	if err := s.db.Select("id", "balance").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	return &BalanceResponse{
		UserID:  user.ID,
		Balance: user.Balance,
	}, nil
}

// GetStats retrieves user statistics
func (s *UserService) GetStats(userID uint) (*StatsResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	stats := &StatsResponse{
		UserID: userID,
	}

	// Calculate work statistics
	var workStats struct {
		TotalTime   int
		TotalEarned float64
	}
	s.db.Model(&model.WorkSession{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(duration_seconds), 0) as total_time, COALESCE(SUM(earned), 0) as total_earned").
		Scan(&workStats)

	stats.TotalWorkTime = workStats.TotalTime
	stats.TotalEarned = workStats.TotalEarned

	// Calculate game statistics
	var gameStats struct {
		TotalSessions int
		TotalBet      float64
		TotalWon      float64
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_sessions, COALESCE(SUM(bet), 0) as total_bet, COALESCE(SUM(win), 0) as total_won").
		Scan(&gameStats)

	stats.TotalGameSessions = gameStats.TotalSessions
	stats.TotalBet = gameStats.TotalBet
	stats.TotalWon = gameStats.TotalWon
	stats.TotalLost = gameStats.TotalBet - gameStats.TotalWon
	stats.NetProfit = gameStats.TotalWon - gameStats.TotalBet

	// Find favorite game (most played)
	var favoriteGame struct {
		GameType string
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ?", userID).
		Select("game_type").
		Group("game_type").
		Order("COUNT(*) DESC").
		Limit(1).
		Scan(&favoriteGame)

	stats.FavoriteGame = favoriteGame.GameType

	// Find biggest win
	var biggestWin struct {
		Win float64
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ?", userID).
		Select("COALESCE(MAX(win), 0) as win").
		Scan(&biggestWin)

	stats.BiggestWin = biggestWin.Win

	// Find biggest loss (highest bet with zero win)
	var biggestLoss struct {
		Bet float64
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ? AND win = 0", userID).
		Select("COALESCE(MAX(bet), 0) as bet").
		Scan(&biggestLoss)

	stats.BiggestLoss = biggestLoss.Bet

	return stats, nil
}

// GetTransactions retrieves user transaction history
func (s *UserService) GetTransactions(userID uint, limit int, offset int) ([]model.Transaction, int64, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("user not found")
		}
		return nil, 0, fmt.Errorf("failed to find user: %w", err)
	}

	var transactions []model.Transaction
	var total int64

	// Count total transactions
	if err := s.db.Model(&model.Transaction{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Get paginated transactions
	query := s.db.Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, total, nil
}

// GetUserItems retrieves user's purchased items
func (s *UserService) GetUserItems(userID uint) ([]model.UserItem, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var userItems []model.UserItem

	// Get user items with item details
	if err := s.db.Preload("Item").Where("user_id = ?", userID).Order("purchased_at DESC").Find(&userItems).Error; err != nil {
		return nil, fmt.Errorf("failed to get user items: %w", err)
	}

	return userItems, nil
}
