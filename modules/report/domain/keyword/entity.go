package keyword

import (
	"github.com/google/uuid"
	"time"
)

type Keyword struct{ state State }
type State struct {
	ID                    uuid.UUID
	AccountID, ProfileID  *string
	GroupName, Word, Kind string
	CountLimit            int
	CreatedAt             time.Time
}

func NewFromState(st State) *Keyword { return &Keyword{state: st} }
func (k *Keyword) State() State      { return k.state }
