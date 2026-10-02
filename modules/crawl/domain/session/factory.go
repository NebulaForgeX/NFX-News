package session

import (
	"time"

	"github.com/google/uuid"
)

func New(accountID, profileID, sourceID *string) *Session {
	id := uuid.Must(uuid.NewV7())
	return NewFromState(State{
		ID: id, AccountID: accountID, ProfileID: profileID, SourceID: sourceID, Status: "running", StartedAt: time.Now().UTC(),
	})
}
