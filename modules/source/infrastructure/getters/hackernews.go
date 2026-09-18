package getters

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *Registry) hackerNews(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://news.ycombinator.com", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find(".athing").Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		a := s.Find(".titleline a").First()
		title := a.Text()
		if id != "" && title != "" {
			out = append(out, Item{ID: id, Title: title, URL: "https://news.ycombinator.com/item?id=" + id})
		}
	})
	return out, nil
}
