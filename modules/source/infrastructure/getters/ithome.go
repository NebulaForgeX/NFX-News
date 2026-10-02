package getters

import "context"

func (r *Registry) ithome(ctx context.Context) ([]Item, error) {
	return r.rss("https://www.ithome.com/rss/")(ctx)
}
