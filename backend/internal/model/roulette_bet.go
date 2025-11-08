package model

import (
	"time"
)

// RouletteBetType represents the type of roulette bet
type RouletteBetType string

const (
	BetTypeStraight RouletteBetType = "straight" // Single number
	BetTypeRed      RouletteBetType = "red"      // Red color
	BetTypeBlack    RouletteBetType = "black"    // Black color
	BetTypeOdd      RouletteBetType = "odd"      // Odd numbers
	BetTypeEven     RouletteBetType = "even"     // Even numbers
	BetTypeDozen1   RouletteBetType = "dozen1"   // 1-12
	BetTypeDozen2   RouletteBetType = "dozen2"   // 13-24
	BetTypeDozen3   RouletteBetType = "dozen3"   // 25-36
	BetTypeLow      RouletteBetType = "low"      // 1-18
	BetTypeHigh     RouletteBetType = "high"     // 19-36
	BetTypeColumn1  RouletteBetType = "column1"  // First column
	BetTypeColumn2  RouletteBetType = "column2"  // Second column
	BetTypeColumn3  RouletteBetType = "column3"  // Third column
)

// RouletteBet represents a single bet in roulette
type RouletteBet struct {
	Type   RouletteBetType `json:"type"`
	Amount float64         `json:"amount"`
	Value  int             `json:"value,omitempty"` // For straight bets (0-36)
}

// RouletteResult represents the result of a roulette spin
type RouletteResult struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	Number       int       `gorm:"not null" json:"number"`        // Winning number (0-36)
	TotalBet     float64   `gorm:"type:decimal(15,2)" json:"total_bet"`
	TotalWin     float64   `gorm:"type:decimal(15,2)" json:"total_win"`
	Bets         string    `gorm:"type:text" json:"bets"` // JSON encoded bets
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for RouletteResult model
func (RouletteResult) TableName() string {
	return "roulette_results"
}
