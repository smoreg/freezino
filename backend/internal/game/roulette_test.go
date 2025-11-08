package game

import (
	"testing"

	"github.com/smoreg/freezino/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouletteGame(t *testing.T) {
	game := NewRouletteGame()

	assert.NotNil(t, game)
	assert.Len(t, game.redNumbers, 18, "should have 18 red numbers")
	assert.Len(t, game.blackNumbers, 18, "should have 18 black numbers")
}

func TestRouletteSpin(t *testing.T) {
	game := NewRouletteGame()

	// Test multiple spins to ensure randomness
	results := make(map[int]bool)
	for i := 0; i < 100; i++ {
		number := game.Spin()
		assert.GreaterOrEqual(t, number, 0, "spin result should be >= 0")
		assert.LessOrEqual(t, number, 36, "spin result should be <= 36")
		results[number] = true
	}

	// Should have some variety in 100 spins
	assert.Greater(t, len(results), 10, "should have variety in spin results")
}

func TestRouletteIsRed(t *testing.T) {
	game := NewRouletteGame()

	tests := []struct {
		number   int
		expected bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{3, true},
		{18, true},
		{19, true},
		{36, true},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.number)), func(t *testing.T) {
			result := game.IsRed(tt.number)
			assert.Equal(t, tt.expected, result, "number %d should be red=%v", tt.number, tt.expected)
		})
	}
}

func TestRouletteIsBlack(t *testing.T) {
	game := NewRouletteGame()

	tests := []struct {
		number   int
		expected bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{4, true},
		{20, true},
		{35, true},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.number)), func(t *testing.T) {
			result := game.IsBlack(tt.number)
			assert.Equal(t, tt.expected, result, "number %d should be black=%v", tt.number, tt.expected)
		})
	}
}

func TestRouletteGetColor(t *testing.T) {
	game := NewRouletteGame()

	tests := []struct {
		number   int
		expected string
	}{
		{0, "green"},
		{1, "red"},
		{2, "black"},
		{36, "red"},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.number)), func(t *testing.T) {
			color := game.GetColor(tt.number)
			assert.Equal(t, tt.expected, color)
		})
	}
}

func TestRouletteCalculatePayoutStraight(t *testing.T) {
	game := NewRouletteGame()

	bet := model.RouletteBet{
		Type:   model.BetTypeStraight,
		Value:  17,
		Amount: 10.0,
	}

	// Winning bet
	payout := game.CalculatePayout(bet, 17)
	assert.Equal(t, 360.0, payout, "straight bet should pay 36x")

	// Losing bet
	payout = game.CalculatePayout(bet, 18)
	assert.Equal(t, 0.0, payout, "losing bet should return 0")
}

func TestRouletteCalculatePayoutRed(t *testing.T) {
	game := NewRouletteGame()

	bet := model.RouletteBet{
		Type:   model.BetTypeRed,
		Amount: 100.0,
	}

	// Win on red number
	payout := game.CalculatePayout(bet, 1)
	assert.Equal(t, 200.0, payout, "red bet should pay 2x")

	// Lose on black
	payout = game.CalculatePayout(bet, 2)
	assert.Equal(t, 0.0, payout)

	// Lose on green (0)
	payout = game.CalculatePayout(bet, 0)
	assert.Equal(t, 0.0, payout)
}

func TestRouletteCalculatePayoutBlack(t *testing.T) {
	game := NewRouletteGame()

	bet := model.RouletteBet{
		Type:   model.BetTypeBlack,
		Amount: 50.0,
	}

	// Win on black number
	payout := game.CalculatePayout(bet, 2)
	assert.Equal(t, 100.0, payout)

	// Lose on red
	payout = game.CalculatePayout(bet, 1)
	assert.Equal(t, 0.0, payout)
}

func TestRouletteCalculatePayoutOddEven(t *testing.T) {
	game := NewRouletteGame()

	oddBet := model.RouletteBet{
		Type:   model.BetTypeOdd,
		Amount: 25.0,
	}

	evenBet := model.RouletteBet{
		Type:   model.BetTypeEven,
		Amount: 25.0,
	}

	// Odd wins
	payout := game.CalculatePayout(oddBet, 17)
	assert.Equal(t, 50.0, payout)

	// Even wins
	payout = game.CalculatePayout(evenBet, 18)
	assert.Equal(t, 50.0, payout)

	// 0 is neither odd nor even
	payout = game.CalculatePayout(oddBet, 0)
	assert.Equal(t, 0.0, payout)

	payout = game.CalculatePayout(evenBet, 0)
	assert.Equal(t, 0.0, payout)
}

func TestRouletteCalculatePayoutDozen(t *testing.T) {
	game := NewRouletteGame()

	dozen1 := model.RouletteBet{
		Type:   model.BetTypeDozen1,
		Amount: 30.0,
	}

	dozen2 := model.RouletteBet{
		Type:   model.BetTypeDozen2,
		Amount: 30.0,
	}

	dozen3 := model.RouletteBet{
		Type:   model.BetTypeDozen3,
		Amount: 30.0,
	}

	// Dozen 1 (1-12)
	payout := game.CalculatePayout(dozen1, 5)
	assert.Equal(t, 90.0, payout, "dozen bet should pay 3x")

	// Dozen 2 (13-24)
	payout = game.CalculatePayout(dozen2, 20)
	assert.Equal(t, 90.0, payout)

	// Dozen 3 (25-36)
	payout = game.CalculatePayout(dozen3, 30)
	assert.Equal(t, 90.0, payout)

	// Loss
	payout = game.CalculatePayout(dozen1, 25)
	assert.Equal(t, 0.0, payout)
}

func TestRouletteCalculatePayoutLowHigh(t *testing.T) {
	game := NewRouletteGame()

	lowBet := model.RouletteBet{
		Type:   model.BetTypeLow,
		Amount: 40.0,
	}

	highBet := model.RouletteBet{
		Type:   model.BetTypeHigh,
		Amount: 40.0,
	}

	// Low (1-18) wins
	payout := game.CalculatePayout(lowBet, 10)
	assert.Equal(t, 80.0, payout)

	// High (19-36) wins
	payout = game.CalculatePayout(highBet, 25)
	assert.Equal(t, 80.0, payout)

	// 0 loses
	payout = game.CalculatePayout(lowBet, 0)
	assert.Equal(t, 0.0, payout)
}

func TestRouletteCalculatePayoutColumn(t *testing.T) {
	game := NewRouletteGame()

	col1 := model.RouletteBet{
		Type:   model.BetTypeColumn1,
		Amount: 20.0,
	}

	col2 := model.RouletteBet{
		Type:   model.BetTypeColumn2,
		Amount: 20.0,
	}

	col3 := model.RouletteBet{
		Type:   model.BetTypeColumn3,
		Amount: 20.0,
	}

	// Column 1: 1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34
	payout := game.CalculatePayout(col1, 1)
	assert.Equal(t, 60.0, payout, "column bet should pay 3x")

	payout = game.CalculatePayout(col1, 34)
	assert.Equal(t, 60.0, payout)

	// Column 2: 2, 5, 8, ...
	payout = game.CalculatePayout(col2, 2)
	assert.Equal(t, 60.0, payout)

	// Column 3: 3, 6, 9, ...
	payout = game.CalculatePayout(col3, 3)
	assert.Equal(t, 60.0, payout)

	// Loss
	payout = game.CalculatePayout(col1, 2)
	assert.Equal(t, 0.0, payout)
}

func TestRouletteCalculateResultWithMultipleBets(t *testing.T) {
	game := NewRouletteGame()

	bets := []model.RouletteBet{
		{Type: model.BetTypeStraight, Value: 17, Amount: 10.0},
		{Type: model.BetTypeRed, Amount: 50.0},
		{Type: model.BetTypeOdd, Amount: 25.0},
	}

	// Mock spin by replacing it temporarily
	// In real scenario, we'd test the CalculatePayout logic separately
	winningNumber, totalBet, totalWin, err := game.CalculateResult(bets)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, winningNumber, 0)
	assert.LessOrEqual(t, winningNumber, 36)
	assert.Equal(t, 85.0, totalBet)
	assert.GreaterOrEqual(t, totalWin, 0.0, "win should be non-negative")
}

func TestRouletteCalculateResultValidation(t *testing.T) {
	game := NewRouletteGame()

	// Empty bets
	_, _, _, err := game.CalculateResult([]model.RouletteBet{})
	assert.Error(t, err, "should error on empty bets")

	// Negative bet amount
	_, _, _, err = game.CalculateResult([]model.RouletteBet{
		{Type: model.BetTypeRed, Amount: -10.0},
	})
	assert.Error(t, err, "should error on negative bet")

	// Invalid straight bet number
	_, _, _, err = game.CalculateResult([]model.RouletteBet{
		{Type: model.BetTypeStraight, Value: 37, Amount: 10.0},
	})
	assert.Error(t, err, "should error on invalid number")

	_, _, _, err = game.CalculateResult([]model.RouletteBet{
		{Type: model.BetTypeStraight, Value: -1, Amount: 10.0},
	})
	assert.Error(t, err, "should error on negative number")
}

func TestRouletteEncodeDecode(t *testing.T) {
	bets := []model.RouletteBet{
		{Type: model.BetTypeStraight, Value: 17, Amount: 10.0},
		{Type: model.BetTypeRed, Amount: 50.0},
	}

	// Encode
	encoded, err := EncodeBets(bets)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	// Decode
	decoded, err := DecodeBets(encoded)
	require.NoError(t, err)
	assert.Equal(t, len(bets), len(decoded))
	assert.Equal(t, bets[0].Type, decoded[0].Type)
	assert.Equal(t, bets[0].Value, decoded[0].Value)
	assert.Equal(t, bets[0].Amount, decoded[0].Amount)
}

func TestRouletteDecodeInvalidJSON(t *testing.T) {
	_, err := DecodeBets("invalid json")
	assert.Error(t, err, "should error on invalid JSON")
}
