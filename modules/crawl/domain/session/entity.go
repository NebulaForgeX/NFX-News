package session

import (
	"github.com/google/uuid"
	"time"
)

type Session struct{ state State }
type State struct {
	ID           uuid.UUID
	SourceID     *string
	Status       string
	ItemCount    int
	ErrorMessage *string
	StartedAt    time.Time
	FinishedAt   *time.Time
}

func NewFromState(st State) *Session { return &Session{state: st} }
func (s *Session) State() State      { return s.state }
func (s *Session) Finish(status string, count int, errMsg *string, at time.Time) {
	s.state.Status = status
	s.state.ItemCount = count
	s.state.ErrorMessage = errMsg
	s.state.FinishedAt = &at
}
