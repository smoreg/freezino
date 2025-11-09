package model

import (
	"time"
)

// JobType represents the type of job/work
type JobType string

const (
	JobTypeOffice           JobType = "office"
	JobTypeCourier          JobType = "courier"
	JobTypeLabRat           JobType = "lab_rat"
	JobTypeStuntDriver      JobType = "stunt_driver"
	JobTypeDrugDealer       JobType = "drug_dealer"
	JobTypeStreamer         JobType = "streamer"
	JobTypeBottleCollector  JobType = "bottle_collector"
)

// WorkSession represents a work session completed by a user
type WorkSession struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"not null;index:idx_user_completed" json:"user_id"`
	JobType         JobType   `gorm:"size:50;not null;default:'office'" json:"job_type"`
	DurationSeconds int       `gorm:"not null" json:"duration_seconds"`
	Earned          float64   `gorm:"type:decimal(15,2);not null" json:"earned"`
	CompletedAt     time.Time `gorm:"not null;index:idx_user_completed" json:"completed_at"`
	CreatedAt       time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for WorkSession model
func (WorkSession) TableName() string {
	return "work_sessions"
}
