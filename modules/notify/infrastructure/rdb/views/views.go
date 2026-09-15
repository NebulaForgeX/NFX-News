package views

import (
	"github.com/google/uuid"
	"time"
)

type ChannelsActiveView struct {
	ID        uuid.UUID `gorm:"column:id"`
	Kind      string    `gorm:"column:kind"`
	Name      string    `gorm:"column:name"`
	Enabled   bool      `gorm:"column:enabled"`
	Config    []byte    `gorm:"column:config"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ChannelsActiveView) TableName() string { return `notify."ChannelsActiveView"` }

type DeliveriesActiveView struct {
	ID           uuid.UUID  `gorm:"column:id"`
	ChannelID    uuid.UUID  `gorm:"column:channel_id"`
	ReportID     *uuid.UUID `gorm:"column:report_id"`
	Status       string     `gorm:"column:status"`
	ErrorMessage *string    `gorm:"column:error_message"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	SentAt       *time.Time `gorm:"column:sent_at"`
}

func (DeliveriesActiveView) TableName() string { return `notify."DeliveriesActiveView"` }
