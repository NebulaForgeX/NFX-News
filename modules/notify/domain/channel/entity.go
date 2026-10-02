package channel

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct{ state State }
type State struct {
	ID                   uuid.UUID
	AccountID, ProfileID *string
	Kind, Name           string
	Enabled              bool
	Config               []byte
	CreatedAt, UpdatedAt time.Time
}

func NewFromState(st State) *Channel { return &Channel{state: st} }
func (c *Channel) State() State      { return c.state }
