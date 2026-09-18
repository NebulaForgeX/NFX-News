package getters

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (r *Registry) weibo(ctx context.Context) ([]Item, error) {
	url := "https://s.weibo.com/top/summary?cate=realtimehot"
	body, err := r.get(ctx, url, map[string]string{
		"Referer": url,
		"Cookie":  "SUB=_2AkMWIuNSf8NxqwJRmP8dy2rhaoV2ygrEieKgfhKJJRMxHRl-yT9jqk86tRB6PaLNvQZR6zYUcYVT1zSjoSreQHidcUq7",
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find("#pl_top_realtimehot table tbody tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		a := s.Find("td.td-02 a").FilterFunction(func(_ int, sel *goquery.Selection) bool {
			href, _ := sel.Attr("href")
			return href != "" && !strings.Contains(href, "javascript:void(0);")
		}).First()
		title := strings.TrimSpace(a.Text())
		href, _ := a.Attr("href")
		if title == "" || href == "" {
			return
		}
		flag := strings.TrimSpace(s.Find("td.td-03").Text())
		extra := map[string]any{}
		switch flag {
		case "新":
			extra["icon"] = map[string]any{"url": "https://simg.s.weibo.com/moter/flags/1_0.png", "scale": 1.5}
		case "热":
			extra["icon"] = map[string]any{"url": "https://simg.s.weibo.com/moter/flags/2_0.png", "scale": 1.5}
		}
		out = append(out, Item{ID: title, Title: title, URL: "https://s.weibo.com" + href, MobileURL: "https://s.weibo.com" + href, Extra: extra})
	})
	return out, nil
}
