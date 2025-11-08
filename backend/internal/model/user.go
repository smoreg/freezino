package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	GoogleID     *string        `gorm:"uniqueIndex;size:255" json:"google_id,omitempty"`
	Username     string         `gorm:"uniqueIndex;size:255" json:"username,omitempty"`
	Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string         `gorm:"size:255" json:"-"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Avatar       string         `gorm:"size:512" json:"avatar"`
	Balance      float64        `gorm:"type:decimal(15,2);default:1000.00;not null" json:"balance"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
	UserItems    []UserItem    `gorm:"foreignKey:UserID" json:"user_items,omitempty"`
	WorkSessions []WorkSession `gorm:"foreignKey:UserID" json:"work_sessions,omitempty"`
	GameSessions []GameSession `gorm:"foreignKey:UserID" json:"game_sessions,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
