package snapshot

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func New(accountID, profileID *string, mode, title string, payload []byte, itemCount int) (*Snapshot, error) {
	mode, err := NormalizeMode(mode)
	if err != nil {
		return nil, err
	}
	if title == "" {
		title = strings.ToUpper(mode) + " report"
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return NewFromState(State{
		ID: id, AccountID: accountID, ProfileID: profileID, Mode: mode, Title: title, Payload: payload, ItemCount: itemCount, CreatedAt: time.Now().UTC(),
	}), nil
}
