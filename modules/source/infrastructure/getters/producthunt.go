package getters

import "context"

func (r *Registry) producthunt(ctx context.Context) ([]Item, error) {
	return r.rss("https://www.producthunt.com/feed")(ctx)
}
