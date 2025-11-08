package model

import (
	"time"

	"gorm.io/gorm"
)

// ItemType represents the category of an item
type ItemType string

const (
	ItemTypeClothing    ItemType = "clothing"
	ItemTypeCar         ItemType = "car"
	ItemTypeHouse       ItemType = "house"
	ItemTypeAccessories ItemType = "accessories"
)

// ItemRarity represents the rarity level of an item
type ItemRarity string

const (
	ItemRarityCommon    ItemRarity = "common"
	ItemRarityRare      ItemRarity = "rare"
	ItemRarityEpic      ItemRarity = "epic"
	ItemRarityLegendary ItemRarity = "legendary"
)

// Item represents a purchasable item in the shop
type Item struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Type        ItemType       `gorm:"size:50;not null;index" json:"type"`
	Rarity      ItemRarity     `gorm:"size:50;not null;index;default:'common'" json:"rarity"`
	Price       float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	ImageURL    string         `gorm:"size:512" json:"image_url"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	UserItems []UserItem `gorm:"foreignKey:ItemID" json:"user_items,omitempty"`
}

// TableName specifies the table name for Item model
func (Item) TableName() string {
	return "items"
}
