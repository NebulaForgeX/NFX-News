package getters

import (
	"context"
	"regexp"
	"strconv"
	"strings"
)

func (r *Registry) sputnik(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://sputniknews.cn/services/widget/lenta/", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<div[^>]*class="[^"]*lenta__item[^"]*"[^>]*>[\s\S]*?<a[^>]*href="([^"]*)"[^>]*>[\s\S]*?<div[^>]*class="[^"]*lenta__item-text[^"]*"[^>]*>([^<]+)</div>[\s\S]*?<div[^>]*class="[^"]*lenta__item-date[^"]*"[^>]*data-unixtime="(\d+)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	for _, m := range matches {
		path, title := m[1], strings.TrimSpace(m[2])
		unix, _ := strconv.ParseInt(m[3], 10, 64)
		out = append(out, Item{ID: path, Title: title, URL: "https://sputniknews.cn" + path, PubDate: unix * 1000})
	}
	return out, nil
}
