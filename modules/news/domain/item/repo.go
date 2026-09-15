package item

import "context"

type Repo struct{ Create Create }
type Create interface {
	Upsert(ctx context.Context, i *Item) error
}
