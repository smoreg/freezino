package model

import (
	"time"
)

// UserStatus represents special status effects for a user
type UserStatus struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_status" json:"user_id"`
	Status    string    `gorm:"size:50;not null" json:"status"` // "in_jail", "popular_streamer", "car_broken"
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for UserStatus model
func (UserStatus) TableName() string {
	return "user_statuses"
}

// IsExpired checks if the status has expired
func (us *UserStatus) IsExpired() bool {
	return time.Now().After(us.ExpiresAt)
}
