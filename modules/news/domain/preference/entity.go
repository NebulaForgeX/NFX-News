package preference

import (
	"encoding/json"
	"time"
)

type Preference struct{ state State }
type State struct {
	AccountID, ProfileID string
	ColumnOrder, Payload json.RawMessage
	UpdatedAt            time.Time
}

func NewFromState(st State) *Preference { return &Preference{state: st} }
func (p *Preference) State() State      { return p.state }
