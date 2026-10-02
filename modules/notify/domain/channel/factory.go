package channel

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

func New(accountID, profileID *string, kind, name string, enabled bool, cfg []byte) (*Channel, error) {
	kind, err := NormalizeKind(kind)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(name) == "" {
		name = kind
	}
	if len(cfg) == 0 {
		cfg = []byte("{}")
	}
	if !json.Valid(cfg) {
		cfg = []byte("{}")
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return NewFromState(State{
		ID: id, AccountID: accountID, ProfileID: profileID, Kind: kind, Name: name, Enabled: enabled, Config: cfg, CreatedAt: now, UpdatedAt: now,
	}), nil
}
