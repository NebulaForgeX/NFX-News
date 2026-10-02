package getters

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *Registry) fastbullExpress(ctx context.Context) ([]Item, error) {
	base := "https://www.fastbull.com"
	body, err := r.get(ctx, base+"/cn/express-news", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find(".news-list").Each(func(_ int, s *goquery.Selection) {
		a := s.Find(".title_name")
		href, _ := a.Attr("href")
		titleText := a.Text()
		title := titleText
		if i := strings.Index(titleText, "【"); i >= 0 {
			if j := strings.Index(titleText, "】"); j > i {
				inner := titleText[i+len("【") : j]
				if len(inner) >= 4 {
					title = inner
				}
			}
		}
		date, _ := s.Attr("data-date")
		if href == "" || title == "" || date == "" {
			return
		}
		out = append(out, Item{ID: href, Title: title, URL: base + href})
	})
	return out, nil
}

func (r *Registry) fastbullNews(ctx context.Context) ([]Item, error) {
	base := "https://www.fastbull.com"
	body, err := r.get(ctx, base+"/cn/news", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find(".trending_type").Each(func(_ int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		title := a.Find(".title").Text()
		if href == "" || title == "" {
			return
		}
		out = append(out, Item{ID: href, Title: title, URL: base + href})
	})
	return out, nil
}
