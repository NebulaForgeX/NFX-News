package views

import (
	"github.com/google/uuid"
	"time"
)

type SessionsActiveView struct {
	ID           uuid.UUID  `gorm:"column:id"`
	SourceID     *string    `gorm:"column:source_id"`
	Status       string     `gorm:"column:status"`
	ItemCount    int        `gorm:"column:item_count"`
	ErrorMessage *string    `gorm:"column:error_message"`
	StartedAt    time.Time  `gorm:"column:started_at"`
	FinishedAt   *time.Time `gorm:"column:finished_at"`
}

func (SessionsActiveView) TableName() string { return `crawl."SessionsActiveView"` }
