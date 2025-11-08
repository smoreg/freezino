package model

import (
	"time"
)

// GameType represents the type of game
type GameType string

const (
	GameTypeRoulette  GameType = "roulette"
	GameTypeSlots     GameType = "slots"
	GameTypeBlackjack GameType = "blackjack"
	GameTypeCraps     GameType = "craps"
	GameTypeBaccara   GameType = "baccara"
	GameTypeWheel     GameType = "wheel"
	GameTypeKeno      GameType = "keno"
	GameTypePoker     GameType = "poker"
	GameTypeHiLo      GameType = "hilo"
	GameTypeCrash     GameType = "crash"
	GameTypeBingo     GameType = "bingo"
	GameTypePlinko    GameType = "plinko"
)

// GameSession represents a game session played by a user
type GameSession struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_game_type;index:idx_game_sessions_user_created" json:"user_id"`
	GameType  GameType  `gorm:"size:50;not null;index;index:idx_user_game_type" json:"game_type"`
	Bet       float64   `gorm:"type:decimal(15,2);not null" json:"bet"`
	Win       float64   `gorm:"type:decimal(15,2);default:0.00;index:idx_user_win" json:"win"`
	CreatedAt time.Time `gorm:"index:idx_game_sessions_user_created" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for GameSession model
func (GameSession) TableName() string {
	return "game_sessions"
}
