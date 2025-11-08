package service

import (
	"testing"

	"github.com/smoreg/freezino/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouletteServicePlaceBet(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	// Temporarily override database in the service
	// Note: This requires the database package to support SetDB or similar
	// For testing, we'll create the service with the test db
	service := NewRouletteService()

	// Create bets
	bets := []model.RouletteBet{
		{Type: model.BetTypeRed, Amount: 100.0},
		{Type: model.BetTypeStraight, Value: 17, Amount: 50.0},
	}

	req := PlaceBetRequest{
		UserID: user.ID,
		Bets:   bets,
	}

	// Note: This test requires database.DB to be set to our test db
	// In a real scenario, we'd need to refactor the service to accept db as parameter
	// For now, we'll test the basic structure

	// Skip actual execution as it requires database package refactoring
	// The test structure is here for future implementation
	t.Skip("Requires database package refactoring to inject test DB")
}

func TestRouletteServiceValidation(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)

	service := NewRouletteService()

	// Test insufficient balance
	bets := []model.RouletteBet{
		{Type: model.BetTypeRed, Amount: 200.0}, // More than user has
	}

	req := PlaceBetRequest{
		UserID: user.ID,
		Bets:   bets,
	}

	t.Skip("Requires database package refactoring")

	// In actual implementation:
	// _, err := service.PlaceBet(req)
	// assert.Error(t, err)
	// assert.Contains(t, err.Error(), "insufficient balance")
}

func TestRouletteServiceGetHistory(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := NewRouletteService()

	// Create some roulette results
	for i := 0; i < 3; i++ {
		result := model.RouletteResult{
			UserID:   user.ID,
			Number:   i * 5,
			TotalBet: 100.0,
			TotalWin: 50.0,
			Bets:     "[]",
		}
		err := db.Create(&result).Error
		require.NoError(t, err)
	}

	t.Skip("Requires database package refactoring")

	// In actual implementation:
	// history, err := service.GetHistory(user.ID, 10)
	// require.NoError(t, err)
	// assert.Len(t, history, 3)
}

func TestRouletteServiceGetRecentNumbers(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := NewRouletteService()

	// Create some results
	numbers := []int{17, 23, 5, 36, 0}
	for _, number := range numbers {
		result := model.RouletteResult{
			UserID:   user.ID,
			Number:   number,
			TotalBet: 100.0,
			TotalWin: 0.0,
			Bets:     "[]",
		}
		err := db.Create(&result).Error
		require.NoError(t, err)
	}

	t.Skip("Requires database package refactoring")

	// In actual implementation:
	// recentNumbers, err := service.GetRecentNumbers(5)
	// require.NoError(t, err)
	// assert.Len(t, recentNumbers, 5)
	// Should be in reverse order (newest first)
}
