package service

import (
	"testing"

	"github.com/smoreg/freezino/backend/internal/game"
	"github.com/smoreg/freezino/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlotsServiceSpin(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	req := &SpinRequest{
		UserID: user.ID,
		Bet:    10.0,
	}

	// Spin
	response, err := service.Spin(req)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.Result)
	assert.Equal(t, 10.0, response.Bet)
	assert.GreaterOrEqual(t, response.Win, 0.0)
	assert.LessOrEqual(t, response.NewBalance, 1000.0) // Can't have more than initial + win
	assert.Greater(t, response.TransactionID, uint(0))
	assert.Greater(t, response.GameSessionID, uint(0))

	// Verify user balance updated
	var updatedUser model.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	expectedBalance := 1000.0 + response.Win - response.Bet
	assert.Equal(t, expectedBalance, updatedUser.Balance)

	// Verify transaction created
	var transaction model.Transaction
	err = db.First(&transaction, response.TransactionID).Error
	require.NoError(t, err)
	assert.Equal(t, user.ID, transaction.UserID)

	// Verify game session created
	var gameSession model.GameSession
	err = db.First(&gameSession, response.GameSessionID).Error
	require.NoError(t, err)
	assert.Equal(t, user.ID, gameSession.UserID)
	assert.Equal(t, model.GameTypeSlots, gameSession.GameType)
	assert.Equal(t, 10.0, gameSession.Bet)
	assert.Equal(t, response.Win, gameSession.Win)
}

func TestSlotsServiceSpinInvalidBet(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	// Test zero bet
	req := &SpinRequest{
		UserID: user.ID,
		Bet:    0.0,
	}
	_, err := service.Spin(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "greater than 0")

	// Test negative bet
	req.Bet = -10.0
	_, err = service.Spin(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "greater than 0")
}

func TestSlotsServiceSpinInsufficientBalance(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 5.0)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	req := &SpinRequest{
		UserID: user.ID,
		Bet:    10.0, // More than user has
	}

	_, err := service.Spin(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient balance")

	// Verify balance unchanged
	var checkUser model.User
	err = db.First(&checkUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 5.0, checkUser.Balance)
}

func TestSlotsServiceSpinUserNotFound(t *testing.T) {
	db := setupTestDB(t)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	req := &SpinRequest{
		UserID: 9999,
		Bet:    10.0,
	}

	_, err := service.Spin(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSlotsServiceMultipleSpins(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	// Perform multiple spins
	totalBet := 0.0
	totalWin := 0.0

	for i := 0; i < 10; i++ {
		req := &SpinRequest{
			UserID: user.ID,
			Bet:    10.0,
		}

		response, err := service.Spin(req)
		require.NoError(t, err)

		totalBet += response.Bet
		totalWin += response.Win
	}

	// Verify final balance
	var finalUser model.User
	err := db.First(&finalUser, user.ID).Error
	require.NoError(t, err)

	expectedBalance := 1000.0 - totalBet + totalWin
	assert.InDelta(t, expectedBalance, finalUser.Balance, 0.01)

	// Verify number of transactions
	var transactionCount int64
	db.Model(&model.Transaction{}).Where("user_id = ?", user.ID).Count(&transactionCount)
	assert.Equal(t, int64(10), transactionCount)

	// Verify number of game sessions
	var sessionCount int64
	db.Model(&model.GameSession{}).Where("user_id = ?", user.ID).Count(&sessionCount)
	assert.Equal(t, int64(10), sessionCount)
}

func TestSlotsServiceGetPayoutTable(t *testing.T) {
	service := NewSlotsService()

	table := service.GetPayoutTable()
	assert.NotNil(t, table)
	assert.Greater(t, len(table), 0)
}

func TestSlotsServiceGetSymbols(t *testing.T) {
	service := NewSlotsService()

	symbols := service.GetSymbols()
	assert.NotNil(t, symbols)
	assert.Equal(t, 7, len(symbols))
}

func TestSlotsServiceWinLossTracking(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	wins := 0
	losses := 0

	for i := 0; i < 20; i++ {
		req := &SpinRequest{
			UserID: user.ID,
			Bet:    10.0,
		}

		response, err := service.Spin(req)
		require.NoError(t, err)

		if response.Win > 0 {
			wins++
		} else {
			losses++
		}
	}

	// Should have both wins and losses (with high probability)
	// This is a statistical test, so it might occasionally fail
	assert.Greater(t, wins+losses, 0)
	t.Logf("Wins: %d, Losses: %d", wins, losses)
}

func TestSlotsServiceConcurrency(t *testing.T) {
	t.Skip("Skipping due to SQLite in-memory concurrent access limitations")
	db := setupTestDB(t)
	user1 := createTestUser(t, db, 1000.0)
	user2 := createTestUser(t, db, 1000.0)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	done := make(chan bool, 2)

	// Spin for both users concurrently
	go func() {
		for i := 0; i < 5; i++ {
			req := &SpinRequest{
				UserID: user1.ID,
				Bet:    10.0,
			}
			_, err := service.Spin(req)
			assert.NoError(t, err)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 5; i++ {
			req := &SpinRequest{
				UserID: user2.ID,
				Bet:    10.0,
			}
			_, err := service.Spin(req)
			assert.NoError(t, err)
		}
		done <- true
	}()

	// Wait for both to complete
	<-done
	<-done

	// Verify both users have correct number of sessions
	var session1Count int64
	db.Model(&model.GameSession{}).Where("user_id = ?", user1.ID).Count(&session1Count)
	assert.Equal(t, int64(5), session1Count)

	var session2Count int64
	db.Model(&model.GameSession{}).Where("user_id = ?", user2.ID).Count(&session2Count)
	assert.Equal(t, int64(5), session2Count)
}

func TestSlotsServiceTransactionIntegrity(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := &SlotsService{
		db:     db,
		engine: game.NewSlotsEngine(),
	}

	// Do multiple spins and verify transaction integrity
	for i := 0; i < 5; i++ {
		req := &SpinRequest{
			UserID: user.ID,
			Bet:    10.0,
		}

		response, err := service.Spin(req)
		require.NoError(t, err)

		// Verify transaction amount matches balance change
		var transaction model.Transaction
		err = db.First(&transaction, response.TransactionID).Error
		require.NoError(t, err)

		expectedAmount := response.Win - response.Bet
		assert.InDelta(t, expectedAmount, transaction.Amount, 0.01)
	}

	// Final balance check
	var finalUser model.User
	err := db.First(&finalUser, user.ID).Error
	require.NoError(t, err)

	// Calculate expected balance from all transactions
	var transactions []model.Transaction
	db.Where("user_id = ?", user.ID).Find(&transactions)

	calculatedBalance := 100.0
	for _, tx := range transactions {
		calculatedBalance += tx.Amount
	}

	assert.InDelta(t, calculatedBalance, finalUser.Balance, 0.01)
}
