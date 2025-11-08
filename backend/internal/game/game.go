package game

import (
	"github.com/smoreg/freezino/backend/internal/model"
)

// Game interface defines methods that all games must implement
type Game interface {
	// PlaceBet validates and places a bet for the user
	PlaceBet(userID uint, bet float64) (*model.GameSession, error)

	// Play executes the game logic and returns the result
	Play(session *model.GameSession) (interface{}, error)

	// CalculateWin calculates the winnings based on the game result
	CalculateWin(result interface{}, bet float64) float64

	// GetGameType returns the type of the game
	GetGameType() model.GameType

	// GetHouseEdge returns the house edge percentage (e.g., 0.027 for 2.7%)
	GetHouseEdge() float64
}

// GameResult represents a generic game result
type GameResult struct {
	GameType model.GameType `json:"game_type"`
	Bet      float64        `json:"bet"`
	Win      float64        `json:"win"`
	Data     interface{}    `json:"data"` // Game-specific data
}

// GameConfig holds configuration for game engine
type GameConfig struct {
	MinBet         float64 // Minimum bet amount
	MaxBet         float64 // Maximum bet amount
	DefaultHouseEdge float64 // Default house edge (2.7% for European roulette)
}

// DefaultGameConfig returns default configuration for games
func DefaultGameConfig() *GameConfig {
	return &GameConfig{
		MinBet:         1.0,
		MaxBet:         10000.0,
		DefaultHouseEdge: 0.027, // 2.7% house edge
	}
}

// Predefined house edges for different games
const (
	HouseEdgeRoulette  = 0.027  // 2.7% - European Roulette
	HouseEdgeSlots     = 0.05   // 5% - Slot machines
	HouseEdgeBlackjack = 0.005  // 0.5% - Blackjack (with optimal strategy)
	HouseEdgeCraps     = 0.014  // 1.4% - Craps (pass line)
	HouseEdgeBaccara   = 0.0106 // 1.06% - Baccarat (banker bet)
	HouseEdgeWheel     = 0.056  // 5.6% - Wheel of Fortune
	HouseEdgeKeno      = 0.25   // 25% - Keno (very high)
	HouseEdgePoker     = 0.02   // 2% - Video Poker
	HouseEdgeHiLo      = 0.02   // 2% - Hi-Lo
	HouseEdgeCrash     = 0.03   // 3% - Crash game
	HouseEdgeBingo     = 0.10   // 10% - Bingo
	HouseEdgePlinko    = 0.04   // 4% - Plinko
)

// GetHouseEdgeForGame returns the configured house edge for a specific game type
func GetHouseEdgeForGame(gameType model.GameType) float64 {
	switch gameType {
	case model.GameTypeRoulette:
		return HouseEdgeRoulette
	case model.GameTypeSlots:
		return HouseEdgeSlots
	case model.GameTypeBlackjack:
		return HouseEdgeBlackjack
	case model.GameTypeCraps:
		return HouseEdgeCraps
	case model.GameTypeBaccara:
		return HouseEdgeBaccara
	case model.GameTypeWheel:
		return HouseEdgeWheel
	case model.GameTypeKeno:
		return HouseEdgeKeno
	case model.GameTypePoker:
		return HouseEdgePoker
	case model.GameTypeHiLo:
		return HouseEdgeHiLo
	case model.GameTypeCrash:
		return HouseEdgeCrash
	case model.GameTypeBingo:
		return HouseEdgeBingo
	case model.GameTypePlinko:
		return HouseEdgePlinko
	default:
		return DefaultGameConfig().DefaultHouseEdge
	}
}
