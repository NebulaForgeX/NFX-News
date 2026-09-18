package getters

import (
	"context"
	"regexp"
	"strings"
)

func (r *Registry) smzdm(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://post.smzdm.com/hot_1/", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<a[^>]*href="(https://post\.smzdm.com/p/[^"]+)"[^>]*>([^<]+)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	seen := map[string]struct{}{}
	for _, m := range matches {
		u, title := m[1], strings.TrimSpace(m[2])
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, Item{ID: u, Title: title, URL: u})
	}
	return out, nil
}
