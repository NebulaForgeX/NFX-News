package toolcall

import "context"

type Repo struct{ Create Create }
type Create interface {
	New(ctx context.Context, t *ToolCall) error
}
