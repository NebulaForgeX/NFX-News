package toolcall

import "github.com/google/uuid"

type ToolCall struct{ state State }
type State struct {
	ID           uuid.UUID
	ToolName     string
	Arguments    []byte
	OK           bool
	ErrorMessage *string
}

func NewFromState(st State) *ToolCall { return &ToolCall{state: st} }
func (t *ToolCall) State() State      { return t.state }
