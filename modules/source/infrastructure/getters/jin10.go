package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (r *Registry) jin10(ctx context.Context) ([]Item, error) {
	url := fmt.Sprintf("https://www.jin10.com/flash_newest.js?t=%d", time.Now().UnixMilli())
	body, err := r.get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(string(body))
	raw = strings.TrimPrefix(raw, "var newest =")
	raw = strings.TrimPrefix(raw, "var newest=")
	raw = strings.TrimRight(raw, ";")
	var rows []struct {
		ID   string `json:"id"`
		Data struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(rows))
	for _, k := range rows {
		title := k.Data.Title
		if title == "" {
			title = k.Data.Content
		}
		if title == "" {
			continue
		}
		out = append(out, Item{ID: k.ID, Title: title, URL: "https://www.jin10.com/flash/" + k.ID})
	}
	return out, nil
}
