package getters

import (
	"context"
	"regexp"
	"strings"
)

var chongbuluoLinkRe = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*class="[^"]*xst[^"]*"[^>]*>([^<]+)</a>`)

func (r *Registry) chongbuluoHot(ctx context.Context) ([]Item, error) {
	base := "https://www.chongbuluo.com/"
	body, err := r.get(ctx, base+"forum.php?mod=guide&view=hot", nil)
	if err != nil {
		return nil, err
	}
	html := string(body)
	matches := chongbuluoLinkRe.FindAllStringSubmatch(html, -1)
	out := make([]Item, 0, len(matches))
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		u := base + m[1]
		title := strings.TrimSpace(m[2])
		out = append(out, Item{ID: u, Title: title, URL: u, Extra: map[string]any{"hover": title}})
	}
	return out, nil
}
