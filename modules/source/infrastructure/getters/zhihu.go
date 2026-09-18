package getters

import (
	"context"
	"encoding/json"
	"regexp"
)

func (r *Registry) zhihu(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.zhihu.com/api/v3/feed/topstory/hot-list-web?limit=20&desktop=true", nil)
	if err != nil {
		return nil, err
	}
	var res struct {
		Data []struct {
			Target struct {
				TitleArea   struct{ Text string `json:"text"` } `json:"title_area"`
				ExcerptArea struct{ Text string `json:"text"` } `json:"excerpt_area"`
				MetricsArea struct{ Text string `json:"text"` } `json:"metrics_area"`
				Link        struct{ URL string `json:"url"` }  `json:"link"`
			} `json:"target"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(\d+)$`)
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		id := k.Target.Link.URL
		if m := re.FindString(k.Target.Link.URL); m != "" {
			id = m
		}
		out = append(out, Item{
			ID: id, Title: k.Target.TitleArea.Text, URL: k.Target.Link.URL,
			Extra: map[string]any{"info": k.Target.MetricsArea.Text, "hover": k.Target.ExcerptArea.Text},
		})
	}
	return out, nil
}
