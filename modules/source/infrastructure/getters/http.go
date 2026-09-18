package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mmcdole/gofeed"
)

func (r *Registry) get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func (r *Registry) getJSON(ctx context.Context, url string, headers map[string]string, dest any) error {
	body, err := r.get(ctx, url, headers)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

func (r *Registry) rss(url string) Getter {
	return func(ctx context.Context) ([]Item, error) {
		fp := gofeed.NewParser()
		fp.Client = r.http
		feed, err := fp.ParseURLWithContext(url, ctx)
		if err != nil {
			return nil, err
		}
		items := make([]Item, 0, len(feed.Items))
		for _, it := range feed.Items {
			link := it.Link
			if link == "" {
				link = it.GUID
			}
			item := Item{ID: link, Title: it.Title, URL: link}
			if it.PublishedParsed != nil {
				item.PubDate = it.PublishedParsed.UnixMilli()
			}
			items = append(items, item)
		}
		return items, nil
	}
}

func (r *Registry) rssHub(route string) Getter {
	return r.rss("https://rsshub.rssforever.com" + route + "?format=xml")
}
