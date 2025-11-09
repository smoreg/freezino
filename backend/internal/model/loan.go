package model

import (
	"time"

	"gorm.io/gorm"
)

// LoanType represents the type of loan
type LoanType string

const (
	LoanTypeFriends     LoanType = "friends"
	LoanTypeBank        LoanType = "bank"
	LoanTypeMicrocredit LoanType = "microcredit"
)

// Loan represents a user's loan/credit
type Loan struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	UserID            uint           `gorm:"not null;index" json:"user_id"`
	Type              LoanType       `gorm:"size:50;not null" json:"type"`
	PrincipalAmount   float64        `gorm:"type:decimal(15,2);not null" json:"principal_amount"`    // Original borrowed amount
	RemainingAmount   float64        `gorm:"type:decimal(15,2);not null" json:"remaining_amount"`    // Current debt with interest
	InterestRate      float64        `gorm:"type:decimal(10,6);not null" json:"interest_rate"`       // Interest rate (e.g., 0.05 for 5%)
	InterestPerSecond float64        `gorm:"type:decimal(15,8);not null" json:"interest_per_second"` // How much interest accrues per second
	CollateralItemID  *uint          `gorm:"index" json:"collateral_item_id,omitempty"`              // For bank loans - item held as collateral
	LastInterestAt    time.Time      `gorm:"not null" json:"last_interest_at"`                       // Last time interest was calculated
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CollateralItem *UserItem `gorm:"foreignKey:CollateralItemID" json:"collateral_item,omitempty"`
}

// TableName specifies the table name for Loan model
func (Loan) TableName() string {
	return "loans"
}

// LoanSummary represents aggregate loan information for a user
type LoanSummary struct {
	TotalDebt          float64 `json:"total_debt"`
	InterestPerSecond  float64 `json:"interest_per_second"`
	FriendsLoanCount   int     `json:"friends_loan_count"`
	TotalFriendsLoaned float64 `json:"total_friends_loaned"` // Total ever borrowed from friends
	ActiveLoans        int     `json:"active_loans"`
}
