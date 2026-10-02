package getters

import (
	"context"
	"fmt"
	"time"
)

func (r *Registry) nowcoder(ctx context.Context) ([]Item, error) {
	url := fmt.Sprintf("https://gw-c.nowcoder.com/api/sparta/hot-search/top-hot-pc?size=20&_=%d&t=", time.Now().UnixMilli())
	var res struct {
		Data struct {
			Result []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				Type  int    `json:"type"`
				UUID  string `json:"uuid"`
			} `json:"result"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, url, nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.Result))
	for _, k := range res.Data.Result {
		id, link := k.ID, "https://www.nowcoder.com/discuss/"+k.ID
		if k.Type == 74 && k.UUID != "" {
			id, link = k.UUID, "https://www.nowcoder.com/feed/main/detail/"+k.UUID
		}
		out = append(out, Item{ID: id, Title: k.Title, URL: link})
	}
	return out, nil
}
