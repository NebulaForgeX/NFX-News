package keyword

import "context"

type Repo struct{ Create Create }
type Create interface {
	New(ctx context.Context, k *Keyword) error
}
