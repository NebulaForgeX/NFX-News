package models

import (
	"encoding/json"
	"time"
)

type Preference struct {
	AccountID   string          `gorm:"column:account_id;type:uuid;primaryKey"`
	ProfileID   string          `gorm:"column:profile_id;type:uuid;primaryKey"`
	ColumnOrder json.RawMessage `gorm:"column:column_order;type:jsonb"`
	Payload     json.RawMessage `gorm:"column:payload;type:jsonb"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (Preference) TableName() string { return "news.profile_preferences" }
