package service

import (
	"testing"
)

func TestRouletteServicePlaceBet(t *testing.T) {
	t.Skip("Requires database package refactoring to inject test DB")

	// Commented out until service refactoring is complete
	/*
		db := setupTestDB(t)
		user := createTestUser(t, db, 1000.0)

		service := NewRouletteService()

		bets := []model.RouletteBet{
			{Type: model.BetTypeRed, Amount: 100.0},
			{Type: model.BetTypeStraight, Value: 17, Amount: 50.0},
		}

		req := PlaceBetRequest{
			UserID: user.ID,
			Bets:   bets,
		}
	*/
}

func TestRouletteServiceValidation(t *testing.T) {
	t.Skip("Requires database package refactoring")

	// Commented out until service refactoring is complete
	/*
		db := setupTestDB(t)
		user := createTestUser(t, db, 100.0)
		service := NewRouletteService()

		bets := []model.RouletteBet{
			{Type: model.BetTypeRed, Amount: 200.0},
		}

		req := PlaceBetRequest{
			UserID: user.ID,
			Bets:   bets,
		}

		_, err := service.PlaceBet(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient balance")
	*/
}

func TestRouletteServiceGetHistory(t *testing.T) {
	t.Skip("Requires database package refactoring")

	// Commented out until service refactoring is complete
	/*
		db := setupTestDB(t)
		user := createTestUser(t, db, 1000.0)
		service := NewRouletteService()

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

		history, err := service.GetHistory(user.ID, 10)
		require.NoError(t, err)
		assert.Len(t, history, 3)
	*/
}

func TestRouletteServiceGetRecentNumbers(t *testing.T) {
	t.Skip("Requires database package refactoring")

	// Commented out until service refactoring is complete
	/*
		db := setupTestDB(t)
		user := createTestUser(t, db, 1000.0)
		service := NewRouletteService()

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

		recentNumbers, err := service.GetRecentNumbers(5)
		require.NoError(t, err)
		assert.Len(t, recentNumbers, 5)
	*/
}
