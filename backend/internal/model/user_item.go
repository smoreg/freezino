package model

import (
	"time"
)

// UserItem represents a user's purchased item
type UserItem struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	ItemID      uint      `gorm:"not null;index" json:"item_id"`
	PurchasedAt time.Time `gorm:"not null" json:"purchased_at"`
	IsEquipped  bool      `gorm:"default:false" json:"is_equipped"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Item Item `gorm:"foreignKey:ItemID" json:"item,omitempty"`
}

// TableName specifies the table name for UserItem model
func (UserItem) TableName() string {
	return "user_items"
}
