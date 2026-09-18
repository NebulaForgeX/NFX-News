package models

import (
	"github.com/google/uuid"
	"time"
)

type Keyword struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID  *string   `gorm:"column:account_id;type:uuid"`
	ProfileID  *string   `gorm:"column:profile_id;type:uuid"`
	GroupName  string
	Word       string
	Kind       string
	CountLimit int
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (Keyword) TableName() string { return "report.keywords" }

type Snapshot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID *string   `gorm:"column:account_id;type:uuid"`
	ProfileID *string   `gorm:"column:profile_id;type:uuid"`
	Mode      string
	Title     string
	Payload   []byte `gorm:"type:jsonb"`
	ItemCount int
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Snapshot) TableName() string { return "report.snapshots" }
