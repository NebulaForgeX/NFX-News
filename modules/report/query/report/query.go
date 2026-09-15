package report

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type KeywordVO struct {
	ID         uuid.UUID `json:"id"`
	GroupName  string    `json:"group_name"`
	Word       string    `json:"word"`
	Kind       string    `json:"kind"`
	CountLimit int       `json:"count_limit"`
	CreatedAt  time.Time `json:"created_at"`
}
type SnapshotVO struct {
	ID         uuid.UUID `json:"id"`
	Mode       string    `json:"mode"`
	Title      string    `json:"title"`
	Payload    []byte    `json:"-"`
	PayloadObj any       `json:"payload"`
	ItemCount  int       `json:"item_count"`
	CreatedAt  time.Time `json:"created_at"`
}
type Query struct {
	Keywords  Keywords
	Snapshots Snapshots
}
type Keywords interface {
	All(ctx context.Context) ([]KeywordVO, error)
}
type Snapshots interface {
	Recent(ctx context.Context, limit int) ([]SnapshotVO, error)
	ByID(ctx context.Context, id uuid.UUID) (*SnapshotVO, error)
}
