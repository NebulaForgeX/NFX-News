package getters

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *Registry) gelonghui(ctx context.Context) ([]Item, error) {
	base := "https://www.gelonghui.com"
	body, err := r.get(ctx, base+"/news/", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find(".article-content").Each(func(_ int, s *goquery.Selection) {
		a := s.Find(".detail-right>a")
		href, _ := a.Attr("href")
		title := strings.TrimSpace(a.Find("h2").Text())
		if href == "" || title == "" {
			return
		}
		out = append(out, Item{ID: href, Title: title, URL: base + href})
	})
	return out, nil
}
