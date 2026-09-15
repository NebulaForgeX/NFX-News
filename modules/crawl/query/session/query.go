package session

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type SessionVO struct {
	ID           uuid.UUID  `json:"id"`
	SourceID     *string    `json:"source_id"`
	Status       string     `json:"status"`
	ItemCount    int        `json:"item_count"`
	ErrorMessage *string    `json:"error_message"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
}
type Query struct{ List List }
type List interface {
	Recent(ctx context.Context, limit int) ([]SessionVO, error)
	ByID(ctx context.Context, id uuid.UUID) (*SessionVO, error)
}
