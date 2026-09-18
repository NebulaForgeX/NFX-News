package getters

import (
	"context"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var chongbuluoLinkRe = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*class="[^"]*xst[^"]*"[^>]*>([^<]+)</a>`)

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
