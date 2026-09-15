package views

import (
	"encoding/json"
	"time"
)

type ItemsActiveView struct {
	ID         string     `gorm:"column:id"`
	SourceID   string     `gorm:"column:source_id"`
	OriginalID string     `gorm:"column:original_id"`
	Title      string     `gorm:"column:title"`
	URL        string     `gorm:"column:url"`
	MobileURL  *string    `gorm:"column:mobile_url"`
	PubDate    *time.Time `gorm:"column:pub_date"`
	Extra      []byte     `gorm:"column:extra"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}

func (ItemsActiveView) TableName() string { return `news."ItemsActiveView"` }

type ProfilePreferencesActiveView struct {
	AccountID   string          `gorm:"column:account_id"`
	ProfileID   string          `gorm:"column:profile_id"`
	ColumnOrder json.RawMessage `gorm:"column:column_order"`
	Payload     json.RawMessage `gorm:"column:payload"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (ProfilePreferencesActiveView) TableName() string { return `news."ProfilePreferencesActiveView"` }
