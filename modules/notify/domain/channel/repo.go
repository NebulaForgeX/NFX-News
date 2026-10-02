package channel

import "context"

type Repo struct {
	Create Create
	Update Update
}
type Create interface {
	New(ctx context.Context, c *Channel) error
}
type Update interface {
	Generic(ctx context.Context, c *Channel) error
}
