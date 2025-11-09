package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// WORK_DURATION is the required work duration in seconds (3 minutes)
	WORK_DURATION = 180

	// WORK_REWARD is the amount earned for completing work
	WORK_REWARD = 500.0
)

// WorkSession represents an active work session
type ActiveWorkSession struct {
	UserID    uint
	StartedAt time.Time
}

// WorkService provides business logic for work operations
type WorkService struct {
	db             *gorm.DB
	activeSessions map[uint]time.Time
	mu             sync.RWMutex
}

// NewWorkService creates a new work service instance
func NewWorkService() *WorkService {
	return &WorkService{
		db:             database.GetDB(),
		activeSessions: make(map[uint]time.Time),
	}
}

// StartWorkResponse represents the response for starting work
type StartWorkResponse struct {
	UserID      uint      `json:"user_id"`
	StartedAt   time.Time `json:"started_at"`
	DurationSec int       `json:"duration_seconds"`
	Reward      float64   `json:"reward"`
	CompletesAt time.Time `json:"completes_at"`
}

// WorkStatusResponse represents the current work status
type WorkStatusResponse struct {
	IsWorking    bool       `json:"is_working"`
	UserID       uint       `json:"user_id"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	DurationSec  int        `json:"duration_seconds"`
	ElapsedSec   int        `json:"elapsed_seconds,omitempty"`
	RemainingSec int        `json:"remaining_seconds,omitempty"`
	Progress     float64    `json:"progress,omitempty"` // 0.0 to 1.0
	CanComplete  bool       `json:"can_complete"`
	Reward       float64    `json:"reward"`
	CompletesAt  *time.Time `json:"completes_at,omitempty"`
}

// CompleteWorkResponse represents the response for completing work
type CompleteWorkResponse struct {
	UserID         uint      `json:"user_id"`
	Earned         float64   `json:"earned"`
	BaseReward     float64   `json:"base_reward"`
	NewBalance     float64   `json:"new_balance"`
	DurationSec    int       `json:"duration_seconds"`
	CompletedAt    time.Time `json:"completed_at"`
	TransactionID  uint      `json:"transaction_id"`
	WorkSessionID  uint      `json:"work_session_id"`
	HasClothing    bool      `json:"has_clothing"`
	HasCar         bool      `json:"has_car"`
	ClothingBonus  float64   `json:"clothing_bonus"`  // 0 or -250
	CarBonus       float64   `json:"car_bonus"`       // 0 or +250
}

// WorkHistoryItem represents a work session in history
type WorkHistoryItem struct {
	ID              uint      `json:"id"`
	DurationSeconds int       `json:"duration_seconds"`
	Earned          float64   `json:"earned"`
	CompletedAt     time.Time `json:"completed_at"`
}

// WorkHistoryResponse represents work history with pagination
type WorkHistoryResponse struct {
	Sessions []WorkHistoryItem `json:"sessions"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// StartWork initiates a new work session for the user
func (s *WorkService) StartWork(userID uint) (*StartWorkResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is already working
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.activeSessions[userID]; exists {
		return nil, fmt.Errorf("work session already in progress")
	}

	// Start new work session
	now := time.Now()
	s.activeSessions[userID] = now

	return &StartWorkResponse{
		UserID:      userID,
		StartedAt:   now,
		DurationSec: WORK_DURATION,
		Reward:      WORK_REWARD,
		CompletesAt: now.Add(time.Duration(WORK_DURATION) * time.Second),
	}, nil
}

// GetStatus returns the current work status for the user
func (s *WorkService) GetStatus(userID uint) (*WorkStatusResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	s.mu.RLock()
	startTime, isWorking := s.activeSessions[userID]
	s.mu.RUnlock()

	if !isWorking {
		return &WorkStatusResponse{
			IsWorking:   false,
			UserID:      userID,
			DurationSec: WORK_DURATION,
			Reward:      WORK_REWARD,
			CanComplete: false,
		}, nil
	}

	// Calculate progress
	now := time.Now()
	elapsed := int(now.Sub(startTime).Seconds())
	remaining := WORK_DURATION - elapsed
	if remaining < 0 {
		remaining = 0
	}

	progress := float64(elapsed) / float64(WORK_DURATION)
	if progress > 1.0 {
		progress = 1.0
	}

	canComplete := elapsed >= WORK_DURATION
	completesAt := startTime.Add(time.Duration(WORK_DURATION) * time.Second)

	return &WorkStatusResponse{
		IsWorking:    true,
		UserID:       userID,
		StartedAt:    &startTime,
		DurationSec:  WORK_DURATION,
		ElapsedSec:   elapsed,
		RemainingSec: remaining,
		Progress:     progress,
		CanComplete:  canComplete,
		Reward:       WORK_REWARD,
		CompletesAt:  &completesAt,
	}, nil
}

// CompleteWork completes the work session and awards the user
func (s *WorkService) CompleteWork(userID uint) (*CompleteWorkResponse, error) {
	// Get and validate active session
	s.mu.Lock()
	startTime, exists := s.activeSessions[userID]
	if !exists {
		s.mu.Unlock()
		return nil, fmt.Errorf("no active work session")
	}

	// Check if enough time has passed
	now := time.Now()
	elapsed := int(now.Sub(startTime).Seconds())
	if elapsed < WORK_DURATION {
		s.mu.Unlock()
		remaining := WORK_DURATION - elapsed
		return nil, fmt.Errorf("work not completed yet, %d seconds remaining", remaining)
	}

	// Remove from active sessions
	delete(s.activeSessions, userID)
	s.mu.Unlock()

	// Start database transaction
	var response *CompleteWorkResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get user with lock
		var user model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("failed to find user: %w", err)
		}

		// Check for clothing and car bonuses
		var clothingCount int64
		var carCount int64

		// Count equipped clothing items
		if err := tx.Model(&model.UserItem{}).
			Joins("JOIN items ON items.id = user_items.item_id").
			Where("user_items.user_id = ? AND user_items.equipped = ? AND items.type = ?",
				userID, true, model.ItemTypeClothing).
			Count(&clothingCount).Error; err != nil {
			return fmt.Errorf("failed to check clothing: %w", err)
		}

		// Count equipped car items
		if err := tx.Model(&model.UserItem{}).
			Joins("JOIN items ON items.id = user_items.item_id").
			Where("user_items.user_id = ? AND user_items.equipped = ? AND items.type = ?",
				userID, true, model.ItemTypeCar).
			Count(&carCount).Error; err != nil {
			return fmt.Errorf("failed to check car: %w", err)
		}

		// Calculate earnings with modifiers
		earnedAmount := WORK_REWARD
		clothingBonus := 0.0
		carBonus := 0.0

		// No clothing penalty: -250
		if clothingCount == 0 {
			clothingBonus = -250.0
		}

		// Car bonus: +250
		if carCount > 0 {
			carBonus = 250.0
		}

		earnedAmount += clothingBonus + carBonus

		// Ensure minimum earning of 0
		if earnedAmount < 0 {
			earnedAmount = 0
		}

		// Update user balance
		newBalance := user.Balance + earnedAmount
		if err := tx.Model(&user).Update("balance", newBalance).Error; err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		// Create transaction record
		description := fmt.Sprintf("Work reward for %d seconds", WORK_DURATION)
		if clothingBonus < 0 {
			description += " (no clothing penalty)"
		}
		if carBonus > 0 {
			description += " (car bonus)"
		}

		transaction := model.Transaction{
			UserID:      userID,
			Type:        model.TransactionTypeWork,
			Amount:      earnedAmount,
			Description: description,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Create work session record
		workSession := model.WorkSession{
			UserID:          userID,
			DurationSeconds: WORK_DURATION,
			Earned:          earnedAmount,
			CompletedAt:     now,
		}
		if err := tx.Create(&workSession).Error; err != nil {
			return fmt.Errorf("failed to create work session: %w", err)
		}

		response = &CompleteWorkResponse{
			UserID:        userID,
			Earned:        earnedAmount,
			BaseReward:    WORK_REWARD,
			NewBalance:    newBalance,
			DurationSec:   WORK_DURATION,
			CompletedAt:   now,
			TransactionID: transaction.ID,
			WorkSessionID: workSession.ID,
			HasClothing:   clothingCount > 0,
			HasCar:        carCount > 0,
			ClothingBonus: clothingBonus,
			CarBonus:      carBonus,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetHistory retrieves work session history for the user
func (s *WorkService) GetHistory(userID uint, limit int, offset int) (*WorkHistoryResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var total int64
	var sessions []model.WorkSession

	// Count total work sessions
	if err := s.db.Model(&model.WorkSession{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count work sessions: %w", err)
	}

	// Get paginated work sessions
	query := s.db.Where("user_id = ?", userID).Order("completed_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get work sessions: %w", err)
	}

	// Convert to response format
	items := make([]WorkHistoryItem, len(sessions))
	for i, session := range sessions {
		items[i] = WorkHistoryItem{
			ID:              session.ID,
			DurationSeconds: session.DurationSeconds,
			Earned:          session.Earned,
			CompletedAt:     session.CompletedAt,
		}
	}

	return &WorkHistoryResponse{
		Sessions: items,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}
