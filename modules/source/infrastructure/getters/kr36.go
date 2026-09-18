package getters

// kr36.go is the Go name for copy _36kr.ts (Go ignores files that start with _).

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *Registry) kr36Quick(ctx context.Context) ([]Item, error) {
	base := "https://www.36kr.com"
	body, err := r.get(ctx, base+"/newsflashes", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find(".newsflash-item").Each(func(_ int, s *goquery.Selection) {
		a := s.Find("a.item-title")
		href, _ := a.Attr("href")
		title := a.Text()
		if href == "" || title == "" {
			return
		}
		out = append(out, Item{ID: href, Title: title, URL: base + href})
	})
	return out, nil
}
