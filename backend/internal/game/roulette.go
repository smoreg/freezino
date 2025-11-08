package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/smoreg/freezino/backend/internal/model"
)

// RouletteGame handles European Roulette game logic
type RouletteGame struct {
	// European roulette: 0-36
	redNumbers   []int
	blackNumbers []int
}

// NewRouletteGame creates a new roulette game instance
func NewRouletteGame() *RouletteGame {
	return &RouletteGame{
		redNumbers:   []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36},
		blackNumbers: []int{2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35},
	}
}

// Spin generates a random number from 0 to 36
func (r *RouletteGame) Spin() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(37) // 0-36
}

// IsRed checks if a number is red
func (r *RouletteGame) IsRed(number int) bool {
	for _, n := range r.redNumbers {
		if n == number {
			return true
		}
	}
	return false
}

// IsBlack checks if a number is black
func (r *RouletteGame) IsBlack(number int) bool {
	for _, n := range r.blackNumbers {
		if n == number {
			return true
		}
	}
	return false
}

// CalculatePayout calculates the payout for a bet based on the winning number
func (r *RouletteGame) CalculatePayout(bet model.RouletteBet, winningNumber int) float64 {
	switch bet.Type {
	case model.BetTypeStraight:
		// Straight bet: 35:1
		if bet.Value == winningNumber {
			return bet.Amount * 36 // bet + 35x payout
		}

	case model.BetTypeRed:
		// Red: 1:1
		if r.IsRed(winningNumber) {
			return bet.Amount * 2
		}

	case model.BetTypeBlack:
		// Black: 1:1
		if r.IsBlack(winningNumber) {
			return bet.Amount * 2
		}

	case model.BetTypeOdd:
		// Odd: 1:1 (0 is not odd)
		if winningNumber > 0 && winningNumber%2 == 1 {
			return bet.Amount * 2
		}

	case model.BetTypeEven:
		// Even: 1:1 (0 is not even)
		if winningNumber > 0 && winningNumber%2 == 0 {
			return bet.Amount * 2
		}

	case model.BetTypeDozen1:
		// First dozen (1-12): 2:1
		if winningNumber >= 1 && winningNumber <= 12 {
			return bet.Amount * 3
		}

	case model.BetTypeDozen2:
		// Second dozen (13-24): 2:1
		if winningNumber >= 13 && winningNumber <= 24 {
			return bet.Amount * 3
		}

	case model.BetTypeDozen3:
		// Third dozen (25-36): 2:1
		if winningNumber >= 25 && winningNumber <= 36 {
			return bet.Amount * 3
		}

	case model.BetTypeLow:
		// Low (1-18): 1:1
		if winningNumber >= 1 && winningNumber <= 18 {
			return bet.Amount * 2
		}

	case model.BetTypeHigh:
		// High (19-36): 1:1
		if winningNumber >= 19 && winningNumber <= 36 {
			return bet.Amount * 2
		}

	case model.BetTypeColumn1:
		// First column (1, 4, 7, ..., 34): 2:1
		if winningNumber > 0 && (winningNumber-1)%3 == 0 {
			return bet.Amount * 3
		}

	case model.BetTypeColumn2:
		// Second column (2, 5, 8, ..., 35): 2:1
		if winningNumber > 0 && (winningNumber-2)%3 == 0 {
			return bet.Amount * 3
		}

	case model.BetTypeColumn3:
		// Third column (3, 6, 9, ..., 36): 2:1
		if winningNumber > 0 && winningNumber%3 == 0 {
			return bet.Amount * 3
		}
	}

	return 0 // No win
}

// CalculateResult processes all bets and calculates total win
func (r *RouletteGame) CalculateResult(bets []model.RouletteBet) (int, float64, float64, error) {
	if len(bets) == 0 {
		return 0, 0, 0, fmt.Errorf("no bets placed")
	}

	// Validate bets
	totalBet := 0.0
	for _, bet := range bets {
		if bet.Amount <= 0 {
			return 0, 0, 0, fmt.Errorf("invalid bet amount")
		}
		totalBet += bet.Amount

		// Validate straight bet value
		if bet.Type == model.BetTypeStraight {
			if bet.Value < 0 || bet.Value > 36 {
				return 0, 0, 0, fmt.Errorf("invalid number for straight bet: %d", bet.Value)
			}
		}
	}

	// Spin the wheel
	winningNumber := r.Spin()

	// Calculate total win
	totalWin := 0.0
	for _, bet := range bets {
		totalWin += r.CalculatePayout(bet, winningNumber)
	}

	return winningNumber, totalBet, totalWin, nil
}

// GetColor returns the color of a number
func (r *RouletteGame) GetColor(number int) string {
	if number == 0 {
		return "green"
	}
	if r.IsRed(number) {
		return "red"
	}
	return "black"
}

// EncodeBets encodes bets to JSON string
func EncodeBets(bets []model.RouletteBet) (string, error) {
	data, err := json.Marshal(bets)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DecodeBets decodes bets from JSON string
func DecodeBets(data string) ([]model.RouletteBet, error) {
	var bets []model.RouletteBet
	err := json.Unmarshal([]byte(data), &bets)
	if err != nil {
		return nil, err
	}
	return bets, nil
}
