package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
)

func (r *Registry) ifeng(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.ifeng.com/", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`var\s+allData\s*=\s*(\{[\s\S]*?\});`)
	m := re.FindSubmatch(body)
	if len(m) < 2 {
		return nil, fmt.Errorf("ifeng allData not found")
	}
	var real struct {
		HotNews1 []struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"hotNews1"`
	}
	if err := json.Unmarshal(m[1], &real); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(real.HotNews1))
	for _, n := range real.HotNews1 {
		out = append(out, Item{ID: n.URL, Title: n.Title, URL: n.URL})
	}
	return out, nil
}
