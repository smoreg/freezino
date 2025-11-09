package model

import (
	"time"
)

// UserItem represents a user's purchased item
type UserItem struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"not null;index:idx_user_purchased;index:idx_user_equipped" json:"user_id"`
	ItemID       uint      `gorm:"not null;index" json:"item_id"`
	PurchasedAt  time.Time `gorm:"not null;index:idx_user_purchased" json:"purchased_at"`
	IsEquipped   bool      `gorm:"default:false;index:idx_user_equipped" json:"is_equipped"`
	IsCollateral bool      `gorm:"default:false;index" json:"is_collateral"` // Item is used as loan collateral
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Item Item `gorm:"foreignKey:ItemID" json:"item,omitempty"`
}

// TableName specifies the table name for UserItem model
func (UserItem) TableName() string {
	return "user_items"
}
