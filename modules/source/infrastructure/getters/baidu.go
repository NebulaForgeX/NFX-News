package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
)

func (r *Registry) baidu(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://top.baidu.com/board?tab=realtime", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(?s)<!--s-data:(.*?)-->`)
	m := re.FindSubmatch(body)
	if len(m) < 2 {
		return nil, fmt.Errorf("baidu payload not found")
	}
	var res struct {
		Data struct {
			Cards []struct {
				Content []struct {
					IsTop  bool   `json:"isTop"`
					Word   string `json:"word"`
					RawURL string `json:"rawUrl"`
					Desc   string `json:"desc"`
				} `json:"content"`
			} `json:"cards"`
		} `json:"data"`
	}
	if err := json.Unmarshal(m[1], &res); err != nil {
		return nil, err
	}
	var out []Item
	if len(res.Data.Cards) == 0 {
		return out, nil
	}
	for _, k := range res.Data.Cards[0].Content {
		if k.IsTop {
			continue
		}
		out = append(out, Item{ID: k.RawURL, Title: k.Word, URL: k.RawURL, Extra: map[string]any{"hover": k.Desc}})
	}
	return out, nil
}
