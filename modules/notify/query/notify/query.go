package notify

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type ChannelVO struct {
	ID        uuid.UUID      `json:"id"`
	Kind      string         `json:"kind"`
	Name      string         `json:"name"`
	Enabled   bool           `json:"enabled"`
	Config    []byte         `json:"-"`
	ConfigObj map[string]any `json:"config"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
type DeliveryVO struct {
	ID           uuid.UUID  `json:"id"`
	ChannelID    uuid.UUID  `json:"channel_id"`
	ReportID     *uuid.UUID `json:"report_id"`
	Status       string     `json:"status"`
	ErrorMessage *string    `json:"error_message"`
	CreatedAt    time.Time  `json:"created_at"`
	SentAt       *time.Time `json:"sent_at"`
}
type Query struct {
	Channels   Channels
	Deliveries Deliveries
}
type Channels interface {
	All(ctx context.Context, accountID string) ([]ChannelVO, error)
}
type Deliveries interface {
	Recent(ctx context.Context, accountID string, limit int) ([]DeliveryVO, error)
}
