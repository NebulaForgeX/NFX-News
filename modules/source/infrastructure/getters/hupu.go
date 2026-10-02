package getters

import (
	"context"
	"regexp"
	"strings"
)

func (r *Registry) hupu(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://bbs.hupu.com/topic-daily-hot", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<a href="(/[^"]+?\.html)"[^>]*class="p-title"[^>]*>([^<]+)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	for _, m := range matches {
		path, title := m[1], strings.TrimSpace(m[2])
		out = append(out, Item{ID: path, Title: title, URL: "https://bbs.hupu.com" + path})
	}
	return out, nil
}
