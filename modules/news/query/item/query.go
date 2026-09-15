package item

import (
	"context"
	"encoding/json"
	"time"
)

type ItemVO struct {
	ID         string         `json:"id"`
	SourceID   string         `json:"source_id"`
	OriginalID string         `json:"original_id"`
	Title      string         `json:"title"`
	URL        string         `json:"url"`
	MobileURL  string         `json:"mobile_url,omitempty"`
	PubDate    int64          `json:"pub_date,omitempty"`
	Extra      map[string]any `json:"extra,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type PreferenceVO struct {
	ColumnOrder json.RawMessage `json:"column_order"`
	Payload     json.RawMessage `json:"payload"`
}

type Query struct {
	List List
	Pref Pref
}

type List interface {
	BySource(ctx context.Context, sourceID string, limit int) ([]ItemVO, error)
	Search(ctx context.Context, q string, limit int) ([]ItemVO, error)
}

type Pref interface {
	ByAccountProfile(ctx context.Context, accountID, profileID string) (*PreferenceVO, error)
}
