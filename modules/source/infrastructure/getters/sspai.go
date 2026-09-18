package getters

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func (r *Registry) sspai(ctx context.Context) ([]Item, error) {
	url := fmt.Sprintf("https://sspai.com/api/v1/article/tag/page/get?limit=30&offset=0&created_at=%d&tag=热门文章&released=false", time.Now().UnixMilli())
	var res struct {
		Data []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, url, nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		id := strconv.Itoa(k.ID)
		out = append(out, Item{ID: id, Title: k.Title, URL: "https://sspai.com/post/" + id})
	}
	return out, nil
}
