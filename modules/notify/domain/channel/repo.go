package channel

import "context"

type Repo struct{ Create Create }
type Create interface {
	New(ctx context.Context, c *Channel) error
}
