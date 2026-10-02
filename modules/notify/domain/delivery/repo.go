package delivery

import "context"

type Repo struct {
	Create Create
	Update Update
}
type Create interface {
	New(ctx context.Context, d *Delivery) error
}
type Update interface {
	Generic(ctx context.Context, d *Delivery) error
}
