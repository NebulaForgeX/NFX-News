package getters

import (
	"context"
	"regexp"
	"strings"
)

func (r *Registry) ghxi(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.ghxi.com/category/all", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<li[^>]*>[\s\S]*?<a[^>]*href="([^"]*)"[^>]*class="[^"]*item-title[^"]*"[^>]*>([^<]+)</a>[\s\S]*?<div[^>]*class="[^"]*item-excerpt[^"]*"[^>]*>([^<]+)</div>[\s\S]*?<div[^>]*class="[^"]*date[^"]*"[^>]*>([^<]+)</div>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	for _, m := range matches {
		link, title := m[1], strings.TrimSpace(m[2])
		out = append(out, Item{ID: link, Title: title, URL: link, Extra: map[string]any{"hover": strings.TrimSpace(m[3]), "date": strings.TrimSpace(m[4])}})
	}
	return out, nil
}
