package models

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID            *string   `gorm:"column:account_id;type:uuid"`
	ProfileID            *string   `gorm:"column:profile_id;type:uuid"`
	Kind, Name           string
	Enabled              bool
	Config               []byte `gorm:"type:jsonb"`
	CreatedAt, UpdatedAt time.Time
}

func (Channel) TableName() string { return "notify.channels" }

type Delivery struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID    *string   `gorm:"column:account_id;type:uuid"`
	ProfileID    *string   `gorm:"column:profile_id;type:uuid"`
	ChannelID    uuid.UUID
	ReportID     *uuid.UUID
	Status       string
	ErrorMessage *string
	CreatedAt    time.Time
	SentAt       *time.Time
}

func (Delivery) TableName() string { return "notify.deliveries" }
