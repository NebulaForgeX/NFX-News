package views

import (
	"github.com/google/uuid"
	"time"
)

type KeywordsActiveView struct {
	ID         uuid.UUID `gorm:"column:id"`
	AccountID  *string   `gorm:"column:account_id"`
	ProfileID  *string   `gorm:"column:profile_id"`
	GroupName  string    `gorm:"column:group_name"`
	Word       string    `gorm:"column:word"`
	Kind       string    `gorm:"column:kind"`
	CountLimit int       `gorm:"column:count_limit"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (KeywordsActiveView) TableName() string { return `report."KeywordsActiveView"` }

type SnapshotsActiveView struct {
	ID        uuid.UUID `gorm:"column:id"`
	AccountID *string   `gorm:"column:account_id"`
	ProfileID *string   `gorm:"column:profile_id"`
	Mode      string    `gorm:"column:mode"`
	Title     string    `gorm:"column:title"`
	Payload   []byte    `gorm:"column:payload"`
	ItemCount int       `gorm:"column:item_count"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (SnapshotsActiveView) TableName() string { return `report."SnapshotsActiveView"` }
