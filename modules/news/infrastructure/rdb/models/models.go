package models

import "time"

type Item struct {
	ID         string    `gorm:"type:varchar(255);primaryKey"`
	SourceID   string    `gorm:"type:varchar(64)"`
	OriginalID string    `gorm:"type:varchar(255)"`
	Title      string    `gorm:"type:text"`
	URL        string    `gorm:"type:text"`
	MobileURL  *string   `gorm:"type:text"`
	PubDate    *time.Time `gorm:"type:timestamp"`
	Extra      []byte    `gorm:"type:jsonb"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (Item) TableName() string { return "news.items" }
