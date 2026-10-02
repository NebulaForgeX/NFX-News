package getters

import "context"

func (r *Registry) pcbetaWindows(ctx context.Context) ([]Item, error) {
	return r.rss("https://bbs.pcbeta.com/forum.php?mod=rss&fid=521&auth=0")(ctx)
}

func (r *Registry) pcbetaWindows11(ctx context.Context) ([]Item, error) {
	return r.rss("https://bbs.pcbeta.com/forum.php?mod=rss&fid=563&auth=0")(ctx)
}
