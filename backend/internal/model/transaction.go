package model

import (
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeWork     TransactionType = "work"
	TransactionTypeGame     TransactionType = "game"
	TransactionTypeGameWin  TransactionType = "game_win"
	TransactionTypeGameLoss TransactionType = "game_loss"
	TransactionTypePurchase TransactionType = "purchase"
	TransactionTypeSale     TransactionType = "sale"
	TransactionTypeInitial  TransactionType = "initial"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID           uint            `gorm:"primarykey" json:"id"`
	UserID       uint            `gorm:"not null;index:idx_user_type;index:idx_user_created" json:"user_id"`
	Type         TransactionType `gorm:"size:50;not null;index;index:idx_user_type" json:"type"`
	Amount       float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	BalanceAfter float64         `gorm:"type:decimal(15,2);default:0.00" json:"balance_after"`
	Description  string          `gorm:"size:512" json:"description"`
	CreatedAt    time.Time       `gorm:"index:idx_user_created" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}
