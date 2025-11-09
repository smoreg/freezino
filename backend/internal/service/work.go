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
	JobType   model.JobType
}

// WorkService provides business logic for work operations
type WorkService struct {
	db             *gorm.DB
	activeSessions map[uint]ActiveWorkSession
	mu             sync.RWMutex
}

// NewWorkService creates a new work service instance
func NewWorkService() *WorkService {
	return &WorkService{
		db:             database.GetDB(),
		activeSessions: make(map[uint]ActiveWorkSession),
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
func (s *WorkService) StartWork(userID uint, jobType model.JobType) (*StartWorkResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check specific job requirements
	if err := s.checkJobRequirements(userID, jobType); err != nil {
		return nil, err
	}

	// Check if user is already working
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.activeSessions[userID]; exists {
		return nil, fmt.Errorf("work session already in progress")
	}

	// Start new work session
	now := time.Now()
	s.activeSessions[userID] = ActiveWorkSession{
		UserID:    userID,
		StartedAt: now,
		JobType:   jobType,
	}

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
	session, isWorking := s.activeSessions[userID]
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
	elapsed := int(now.Sub(session.StartedAt).Seconds())
	remaining := WORK_DURATION - elapsed
	if remaining < 0 {
		remaining = 0
	}

	progress := float64(elapsed) / float64(WORK_DURATION)
	if progress > 1.0 {
		progress = 1.0
	}

	canComplete := elapsed >= WORK_DURATION
	completesAt := session.StartedAt.Add(time.Duration(WORK_DURATION) * time.Second)

	return &WorkStatusResponse{
		IsWorking:    true,
		UserID:       userID,
		StartedAt:    &session.StartedAt,
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
	session, exists := s.activeSessions[userID]
	if !exists {
		s.mu.Unlock()
		return nil, fmt.Errorf("no active work session")
	}

	// Check if enough time has passed
	now := time.Now()
	elapsed := int(now.Sub(session.StartedAt).Seconds())
	if elapsed < WORK_DURATION {
		s.mu.Unlock()
		remaining := WORK_DURATION - elapsed
		return nil, fmt.Errorf("work not completed yet, %d seconds remaining", remaining)
	}

	jobType := session.JobType

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

		// Apply job-specific effects and get earnings
		earnedAmount, description, err := s.applyJobEffects(tx, userID, jobType)
		if err != nil {
			return err
		}

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
			JobType:         jobType,
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
			HasClothing:   false,
			HasCar:        false,
			ClothingBonus: 0,
			CarBonus:      0,
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

// checkJobRequirements checks if user meets requirements for a specific job
func (s *WorkService) checkJobRequirements(userID uint, jobType model.JobType) error {
	switch jobType {
	case model.JobTypeOffice:
		// Office requires at least one clothing item
		var clothingCount int64
		if err := s.db.Model(&model.UserItem{}).
			Joins("JOIN items ON items.id = user_items.item_id").
			Where("user_items.user_id = ? AND user_items.equipped = ? AND items.type = ?",
				userID, true, model.ItemTypeClothing).
			Count(&clothingCount).Error; err != nil {
			return fmt.Errorf("failed to check clothing: %w", err)
		}
		if clothingCount == 0 {
			return fmt.Errorf("office_no_clothes")
		}

	case model.JobTypeCourier:
		// Courier requires uniform
		var hasUniform bool
		err := s.db.Model(&model.UserItem{}).
			Joins("JOIN items ON items.id = user_items.item_id").
			Where("user_items.user_id = ? AND user_items.equipped = ? AND items.name = ?",
				userID, true, "Courier Uniform").
			Select("COUNT(*) > 0").
			Scan(&hasUniform).Error
		if err != nil {
			return fmt.Errorf("failed to check uniform: %w", err)
		}
		if !hasUniform {
			return fmt.Errorf("courier_no_uniform")
		}

	case model.JobTypeStuntDriver:
		// Stunt driver requires a car
		var carCount int64
		if err := s.db.Model(&model.UserItem{}).
			Joins("JOIN items ON items.id = user_items.item_id").
			Where("user_items.user_id = ? AND user_items.equipped = ? AND items.type = ?",
				userID, true, model.ItemTypeCar).
			Count(&carCount).Error; err != nil {
			return fmt.Errorf("failed to check car: %w", err)
		}
		if carCount == 0 {
			return fmt.Errorf("stunt_driver_no_car")
		}
	}

	return nil
}

// applyJobEffects applies job-specific effects and returns earned amount and description
func (s *WorkService) applyJobEffects(tx *gorm.DB, userID uint, jobType model.JobType) (float64, string, error) {
	switch jobType {
	case model.JobTypeOffice:
		return s.applyOfficeJob(tx, userID)
	case model.JobTypeCourier:
		return s.applyCourierJob(tx, userID)
	case model.JobTypeLabRat:
		return s.applyLabRatJob(tx, userID)
	case model.JobTypeStuntDriver:
		return s.applyStuntDriverJob(tx, userID)
	case model.JobTypeDrugDealer:
		return s.applyDrugDealerJob(tx, userID)
	case model.JobTypeStreamer:
		return s.applyStreamerJob(tx, userID)
	case model.JobTypeBottleCollector:
		return s.applyBottleCollectorJob(tx, userID)
	default:
		return WORK_REWARD, "Work completed", nil
	}
}

func (s *WorkService) applyOfficeJob(tx *gorm.DB, userID uint) (float64, string, error) {
	// Office job: standard 500, no bonuses
	return WORK_REWARD, "Office work completed", nil
}

func (s *WorkService) applyCourierJob(tx *gorm.DB, userID uint) (float64, string, error) {
	// Courier: base 500 + car bonus 250
	var carCount int64
	if err := tx.Model(&model.UserItem{}).
		Joins("JOIN items ON items.id = user_items.item_id").
		Where("user_items.user_id = ? AND user_items.equipped = ? AND items.type = ?",
			userID, true, model.ItemTypeCar).
		Count(&carCount).Error; err != nil {
		return 0, "", fmt.Errorf("failed to check car: %w", err)
	}

	earned := WORK_REWARD
	desc := "Courier work completed"
	if carCount > 0 {
		earned += 250.0
		desc += " (own car bonus: +$250)"
	}

	return earned, desc, nil
}

func (s *WorkService) applyLabRatJob(tx *gorm.DB, userID uint) (float64, string, error) {
	// Lab rat: 500 + random mutation
	// Get all mutations
	var mutations []model.Item
	if err := tx.Where("type = ?", model.ItemTypeMutation).Find(&mutations).Error; err != nil {
		return 0, "", fmt.Errorf("failed to get mutations: %w", err)
	}

	if len(mutations) > 0 {
		// Give random mutation
		randomMutation := mutations[time.Now().UnixNano()%int64(len(mutations))]

		// Check if user already has this mutation
		var existingCount int64
		tx.Model(&model.UserItem{}).
			Where("user_id = ? AND item_id = ?", userID, randomMutation.ID).
			Count(&existingCount)

		if existingCount == 0 {
			userItem := model.UserItem{
				UserID:   userID,
				ItemID:   randomMutation.ID,
				Equipped: true,
			}
			if err := tx.Create(&userItem).Error; err != nil {
				return 0, "", fmt.Errorf("failed to give mutation: %w", err)
			}
		}

		return WORK_REWARD, fmt.Sprintf("Lab experiment completed. You received: %s!", randomMutation.Name), nil
	}

	return WORK_REWARD, "Lab experiment completed", nil
}

func (s *WorkService) applyStuntDriverJob(tx *gorm.DB, userID uint) (float64, string, error) {
	// Stunt driver: 1500 but car gets broken (unequipped)
	// Unequip all cars
	if err := tx.Model(&model.UserItem{}).
		Where("user_id = ? AND item_id IN (SELECT id FROM items WHERE type = ?)", userID, model.ItemTypeCar).
		Update("equipped", false).Error; err != nil {
		return 0, "", fmt.Errorf("failed to break car: %w", err)
	}

	return 1500.0, "Stunt driving completed! Earned $1500 but your car is broken (unequipped)", nil
}

func (s *WorkService) applyDrugDealerJob(tx *gorm.DB, userID uint) (float64, string, error) {
	// Drug dealer: 2000 but go to jail (8 years = future timestamp)
	// Create user status "in_jail" that expires in 8 years
	jailTime := time.Now().Add(8 * 365 * 24 * time.Hour)
	status := model.UserStatus{
		UserID:    userID,
		Status:    "in_jail",
		ExpiresAt: jailTime,
	}
	if err := tx.Create(&status).Error; err != nil {
		return 0, "", fmt.Errorf("failed to create jail status: %w", err)
	}

	return 2000.0, "You got caught! Earned $2000 but you're in jail for 8 years (you can skip time)", nil
}

func (s *WorkService) applyStreamerJob(tx *gorm.DB, userID uint) (float64, string, error) {
	// Check if already popular
	var popularStatus int64
	tx.Model(&model.UserStatus{}).
		Where("user_id = ? AND status = ?", userID, "popular_streamer").
		Count(&popularStatus)

	if popularStatus > 0 {
		// Already popular, always get 10000
		return 10000.0, "Streaming as popular streamer! Earned $10,000", nil
	}

	// Not popular yet: 70% = 0, 29% = 1, 1% = 10000 + become popular
	roll := time.Now().UnixNano() % 100

	if roll < 70 {
		// 70%: nothing
		return 0.0, "Streaming session completed but no one watched...", nil
	} else if roll < 99 {
		// 29%: $1
		return 1.0, "Streaming session completed. Someone donated $1!", nil
	} else {
		// 1%: $10000 + become popular
		status := model.UserStatus{
			UserID:    userID,
			Status:    "popular_streamer",
			ExpiresAt: time.Now().Add(100 * 365 * 24 * time.Hour), // Forever
		}
		if err := tx.Create(&status).Error; err != nil {
			return 0, "", fmt.Errorf("failed to create popular status: %w", err)
		}

		return 10000.0, "YOU WENT VIRAL! Earned $10,000 and became a popular streamer! All future streams will earn $10,000!", nil
	}
}

func (s *WorkService) applyBottleCollectorJob(tx *gorm.DB, userID uint) (float64, string, error) {
	// Bottle collector: always 100, available to everyone
	return 100.0, "Collected bottles and cans. Earned $100", nil
}

// SkipJailTime removes the jail status from user
func (s *WorkService) SkipJailTime(userID uint) error {
	// Check if user is in jail
	var jailStatus model.UserStatus
	err := s.db.Where("user_id = ? AND status = ?", userID, "in_jail").First(&jailStatus).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("not in jail")
		}
		return fmt.Errorf("failed to check jail status: %w", err)
	}

	// Delete jail status
	if err := s.db.Delete(&jailStatus).Error; err != nil {
		return fmt.Errorf("failed to remove jail status: %w", err)
	}

	return nil
}
