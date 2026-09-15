package systemstate

import (
	"github.com/google/uuid"
	"time"
)

type State struct{ state Inner }
type Inner struct {
	ID                    uuid.UUID
	Initialized           bool
	InitializedAt         *time.Time
	InitializationVersion *string
	ResetCount            int
	CreatedAt, UpdatedAt  time.Time
}

func NewFromState(st Inner) *State { return &State{state: st} }
func (s *State) Inner() Inner      { return s.state }
