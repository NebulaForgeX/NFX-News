package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID    *string   `gorm:"column:account_id;type:uuid"`
	ProfileID    *string   `gorm:"column:profile_id;type:uuid"`
	SourceID     *string
	Status       string
	ItemCount    int
	ErrorMessage *string
	StartedAt    time.Time `gorm:"autoCreateTime"`
	FinishedAt   *time.Time
}

func (Session) TableName() string { return "crawl.sessions" }
