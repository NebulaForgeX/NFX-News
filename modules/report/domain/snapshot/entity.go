package snapshot

import (
	"github.com/google/uuid"
	"time"
)

type Snapshot struct{ state State }
type State struct {
	ID                   uuid.UUID
	AccountID, ProfileID *string
	Mode, Title          string
	Payload              []byte
	ItemCount            int
	CreatedAt            time.Time
}

func NewFromState(st State) *Snapshot { return &Snapshot{state: st} }
func (s *Snapshot) State() State      { return s.state }
