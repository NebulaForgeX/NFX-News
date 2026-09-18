package getters

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *Registry) zaobao(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.zaochenbao.com/realtime/", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(decodeGBK(body)))
	if err != nil {
		return nil, err
	}
	base := "https://www.zaochenbao.com"
	var out []Item
	doc.Find("div.list-block>a.item").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := s.Find(".eps").Text()
		if href == "" || title == "" {
			return
		}
		out = append(out, Item{ID: href, Title: title, URL: base + href})
	})
	return out, nil
}
