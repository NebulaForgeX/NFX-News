package getters

import "context"

func (r *Registry) solidot(ctx context.Context) ([]Item, error) {
	return r.rss("https://www.solidot.org/index.rss")(ctx)
}
