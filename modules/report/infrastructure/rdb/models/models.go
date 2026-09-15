package models

import (
	"github.com/google/uuid"
	"time"
)

type Keyword struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupName  string
	Word       string
	Kind       string
	CountLimit int
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (Keyword) TableName() string { return "report.keywords" }

type Snapshot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Mode      string
	Title     string
	Payload   []byte `gorm:"type:jsonb"`
	ItemCount int
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Snapshot) TableName() string { return "report.snapshots" }
