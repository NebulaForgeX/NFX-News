package preference

import "context"

type Repo struct{ Create Create }
type Create interface {
	Upsert(ctx context.Context, p *Preference) error
}
