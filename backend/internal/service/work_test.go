package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/smoreg/freezino/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use unique in-memory database per test for isolation
	// For concurrency tests, use shared memory mode
	dbName := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "failed to connect to test database")

	// Enable WAL mode for better concurrency
	db.Exec("PRAGMA journal_mode=WAL")
	// Set busy timeout to 5 seconds for locked tables
	db.Exec("PRAGMA busy_timeout=5000")

	// Auto migrate models
	err = db.AutoMigrate(
		&model.User{},
		&model.Transaction{},
		&model.WorkSession{},
		&model.GameSession{},
		&model.Item{},
		&model.UserItem{},
	)
	require.NoError(t, err, "failed to migrate test database")

	return db
}

func createTestUser(t *testing.T, db *gorm.DB, balance float64) *model.User {
	timestamp := time.Now().UnixNano()
	googleID := fmt.Sprintf("test-google-id-%d", timestamp)
	username := fmt.Sprintf("testuser%d", timestamp)
	email := fmt.Sprintf("test%d@example.com", timestamp)

	user := &model.User{
		GoogleID: &googleID,
		Username: username,
		Email:    email,
		Name:     "Test User",
		Balance:  balance,
	}

	err := db.Create(user).Error
	require.NoError(t, err, "failed to create test user")
	return user
}

func TestWorkServiceStartWork(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 0.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Start work
	response, err := service.StartWork(user.ID, model.JobTypeBottleCollector)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, user.ID, response.UserID)
	assert.Equal(t, WORK_DURATION, response.DurationSec)
	assert.Equal(t, WORK_REWARD, response.Reward)
	assert.WithinDuration(t, time.Now(), response.StartedAt, 1*time.Second)
	assert.WithinDuration(t, time.Now().Add(WORK_DURATION*time.Second), response.CompletesAt, 1*time.Second)

	// Verify session is active
	_, exists := service.activeSessions[user.ID]
	assert.True(t, exists, "work session should be active")
}

func TestWorkServiceStartWorkAlreadyWorking(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 0.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Start first work session
	_, err := service.StartWork(user.ID, model.JobTypeBottleCollector)
	require.NoError(t, err)

	// Try to start another session (should fail)
	_, err = service.StartWork(user.ID, model.JobTypeBottleCollector)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already in progress")
}

func TestWorkServiceStartWorkUserNotFound(t *testing.T) {
	db := setupTestDB(t)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Try to start work for non-existent user
	_, err := service.StartWork(9999, model.JobTypeBottleCollector)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestWorkServiceGetStatusNotWorking(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	status, err := service.GetStatus(user.ID)
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.False(t, status.IsWorking)
	assert.Equal(t, user.ID, status.UserID)
	assert.Equal(t, WORK_DURATION, status.DurationSec)
	assert.Equal(t, WORK_REWARD, status.Reward)
	assert.False(t, status.CanComplete)
	assert.Nil(t, status.StartedAt)
	assert.Nil(t, status.CompletesAt)
}

func TestWorkServiceGetStatusWorking(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 0.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Start work
	_, err := service.StartWork(user.ID, model.JobTypeBottleCollector)
	require.NoError(t, err)

	// Small delay to have elapsed time (at least 1 second for int conversion)
	time.Sleep(1100 * time.Millisecond)

	// Get status
	status, err := service.GetStatus(user.ID)
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.True(t, status.IsWorking)
	assert.Equal(t, user.ID, status.UserID)
	assert.NotNil(t, status.StartedAt)
	assert.NotNil(t, status.CompletesAt)
	assert.Greater(t, status.ElapsedSec, 0)
	assert.Less(t, status.ElapsedSec, WORK_DURATION)
	assert.Greater(t, status.RemainingSec, 0)
	assert.Greater(t, status.Progress, 0.0)
	assert.Less(t, status.Progress, 1.0)
	assert.False(t, status.CanComplete, "should not be able to complete immediately")
}

func TestWorkServiceCompleteWorkSuccess(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Start work and simulate it being completed
	startTime := time.Now().Add(-WORK_DURATION * time.Second)
	service.activeSessions[user.ID] = ActiveWorkSession{
		UserID:    user.ID,
		StartedAt: startTime,
		JobType:   model.JobTypeBottleCollector,
	}

	// Complete work (bottle collector earns 100)
	expectedEarned := 100.0 // Bottle collector specific amount
	response, err := service.CompleteWork(user.ID)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, user.ID, response.UserID)
	assert.Equal(t, expectedEarned, response.Earned)
	assert.Equal(t, 200.0, response.NewBalance) // 100 + 100
	assert.Equal(t, WORK_DURATION, response.DurationSec)
	assert.Greater(t, response.TransactionID, uint(0))
	assert.Greater(t, response.WorkSessionID, uint(0))

	// Verify session is removed
	_, exists := service.activeSessions[user.ID]
	assert.False(t, exists, "work session should be removed after completion")

	// Verify user balance updated
	var updatedUser model.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 200.0, updatedUser.Balance)

	// Verify transaction created
	var transaction model.Transaction
	err = db.Where("user_id = ? AND type = ?", user.ID, model.TransactionTypeWork).First(&transaction).Error
	require.NoError(t, err)
	assert.Equal(t, expectedEarned, transaction.Amount)

	// Verify work session created
	var workSession model.WorkSession
	err = db.Where("user_id = ?", user.ID).First(&workSession).Error
	require.NoError(t, err)
	assert.Equal(t, WORK_DURATION, workSession.DurationSeconds)
	assert.Equal(t, expectedEarned, workSession.Earned)
}

func TestWorkServiceCompleteWorkNoActiveSession(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Try to complete without starting
	_, err := service.CompleteWork(user.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active work session")
}

func TestWorkServiceCompleteWorkTooEarly(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Start work
	_, err := service.StartWork(user.ID, model.JobTypeBottleCollector)
	require.NoError(t, err)

	// Try to complete immediately (should fail)
	_, err = service.CompleteWork(user.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not completed yet")
	assert.Contains(t, err.Error(), "remaining")
}

func TestWorkServiceGetHistory(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Create some work sessions
	for i := 0; i < 3; i++ {
		workSession := model.WorkSession{
			UserID:          user.ID,
			DurationSeconds: WORK_DURATION,
			Earned:          WORK_REWARD,
			CompletedAt:     time.Now().Add(-time.Duration(i) * time.Hour),
		}
		err := db.Create(&workSession).Error
		require.NoError(t, err)
	}

	// Get history
	history, err := service.GetHistory(user.ID, 10, 0)
	require.NoError(t, err)
	assert.NotNil(t, history)
	assert.Len(t, history.Sessions, 3)
	assert.Equal(t, int64(3), history.Total)
	assert.Equal(t, 10, history.Limit)
	assert.Equal(t, 0, history.Offset)

	// Verify sessions are ordered by completed_at DESC (newest first)
	for i := 0; i < len(history.Sessions)-1; i++ {
		assert.True(t, history.Sessions[i].CompletedAt.After(history.Sessions[i+1].CompletedAt))
	}
}

func TestWorkServiceGetHistoryWithPagination(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Create 5 work sessions
	for i := 0; i < 5; i++ {
		workSession := model.WorkSession{
			UserID:          user.ID,
			DurationSeconds: WORK_DURATION,
			Earned:          WORK_REWARD,
			CompletedAt:     time.Now().Add(-time.Duration(i) * time.Hour),
		}
		err := db.Create(&workSession).Error
		require.NoError(t, err)
	}

	// Get first page (limit 2)
	history, err := service.GetHistory(user.ID, 2, 0)
	require.NoError(t, err)
	assert.Len(t, history.Sessions, 2)
	assert.Equal(t, int64(5), history.Total)

	// Get second page (limit 2, offset 2)
	history, err = service.GetHistory(user.ID, 2, 2)
	require.NoError(t, err)
	assert.Len(t, history.Sessions, 2)
	assert.Equal(t, int64(5), history.Total)

	// Get third page (limit 2, offset 4)
	history, err = service.GetHistory(user.ID, 2, 4)
	require.NoError(t, err)
	assert.Len(t, history.Sessions, 1)
	assert.Equal(t, int64(5), history.Total)
}

func TestWorkServiceGetHistoryUserNotFound(t *testing.T) {
	db := setupTestDB(t)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	_, err := service.GetHistory(9999, 10, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestWorkServiceConcurrency(t *testing.T) {
	db := setupTestDB(t)
	user1 := createTestUser(t, db, 0.0)
	user2 := createTestUser(t, db, 0.0)

	service := &WorkService{
		db:             db,
		activeSessions: make(map[uint]ActiveWorkSession),
	}

	// Start work for both users concurrently
	done := make(chan bool, 2)

	go func() {
		_, err := service.StartWork(user1.ID, model.JobTypeBottleCollector)
		assert.NoError(t, err)
		done <- true
	}()

	go func() {
		_, err := service.StartWork(user2.ID, model.JobTypeBottleCollector)
		assert.NoError(t, err)
		done <- true
	}()

	// Wait for both to complete
	<-done
	<-done

	// Both should have active sessions
	_, exists1 := service.activeSessions[user1.ID]
	_, exists2 := service.activeSessions[user2.ID]
	assert.True(t, exists1)
	assert.True(t, exists2)
}
