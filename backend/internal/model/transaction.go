package model

import (
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeWork     TransactionType = "work"
	TransactionTypeGameWin  TransactionType = "game_win"
	TransactionTypeGameLoss TransactionType = "game_loss"
	TransactionTypePurchase TransactionType = "purchase"
	TransactionTypeSale     TransactionType = "sale"
	TransactionTypeInitial  TransactionType = "initial"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          uint            `gorm:"primarykey" json:"id"`
	UserID      uint            `gorm:"not null;index" json:"user_id"`
	Type        TransactionType `gorm:"size:50;not null;index" json:"type"`
	Amount      float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description string          `gorm:"size:512" json:"description"`
	CreatedAt   time.Time       `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}
