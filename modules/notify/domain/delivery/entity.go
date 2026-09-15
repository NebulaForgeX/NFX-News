package delivery

import (
	"github.com/google/uuid"
	"time"
)

type Delivery struct{ state State }
type State struct {
	ID, ChannelID uuid.UUID
	ReportID      *uuid.UUID
	Status        string
	ErrorMessage  *string
	CreatedAt     time.Time
	SentAt        *time.Time
}

func NewFromState(st State) *Delivery { return &Delivery{state: st} }
func (d *Delivery) State() State      { return d.state }
func (d *Delivery) Mark(status string, errMsg *string, sent *time.Time) {
	d.state.Status = status
	d.state.ErrorMessage = errMsg
	d.state.SentAt = sent
}
