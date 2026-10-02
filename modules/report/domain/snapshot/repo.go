package snapshot

import "context"

type Repo struct{ Create Create }
type Create interface {
	New(ctx context.Context, s *Snapshot) error
}
