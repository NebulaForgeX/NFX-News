package getters

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *Registry) steam(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://store.steampowered.com/stats/stats/", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find("#detailStats tr.player_count_row").Each(func(_ int, s *goquery.Selection) {
		a := s.Find("a.gameLink")
		href, _ := a.Attr("href")
		name := strings.TrimSpace(a.Text())
		players := strings.TrimSpace(s.Find("td:first-child .currentServers").Text())
		if href == "" || name == "" {
			return
		}
		out = append(out, Item{ID: href, Title: name, URL: href, Extra: map[string]any{"info": players}})
	})
	return out, nil
}
