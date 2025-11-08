package model

import (
	"time"

	"gorm.io/gorm"
)

// ContactMessage represents a message from the contact form
type ContactMessage struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;not null" json:"email"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	IsRead    bool           `gorm:"default:false;not null" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for ContactMessage model
func (ContactMessage) TableName() string {
	return "contact_messages"
}
