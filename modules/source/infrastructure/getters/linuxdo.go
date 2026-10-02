package getters

import (
	"context"
	"time"
)

func (r *Registry) linuxdoHot(ctx context.Context) ([]Item, error) {
	return r.linuxdo(ctx, "https://linux.do/top/daily.json", false)
}

func (r *Registry) linuxdoLatest(ctx context.Context) ([]Item, error) {
	return r.linuxdo(ctx, "https://linux.do/latest.json?order=created", true)
}

func (r *Registry) linuxdo(ctx context.Context, api string, withDate bool) ([]Item, error) {
	var res struct {
		TopicList struct {
			Topics []struct {
				ID        int    `json:"id"`
				Title     string `json:"title"`
				Visible   bool   `json:"visible"`
				Archived  bool   `json:"archived"`
				Pinned    bool   `json:"pinned"`
				CreatedAt string `json:"created_at"`
			} `json:"topics"`
		} `json:"topic_list"`
	}
	if err := r.getJSON(ctx, api, nil, &res); err != nil {
		return nil, err
	}
	var out []Item
	for _, k := range res.TopicList.Topics {
		if !k.Visible || k.Archived || k.Pinned {
			continue
		}
		it := Item{ID: itoa(k.ID), Title: k.Title, URL: "https://linux.do/t/topic/" + itoa(k.ID)}
		if withDate {
			if t, err := time.Parse(time.RFC3339, k.CreatedAt); err == nil {
				it.PubDate = t.UnixMilli()
			}
		}
		out = append(out, it)
	}
	return out, nil
}
