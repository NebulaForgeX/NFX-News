package getters

import "context"

func (r *Registry) juejin(ctx context.Context) ([]Item, error) {
	var res struct {
		Data []struct {
			Content struct {
				Title     string `json:"title"`
				ContentID string `json:"content_id"`
			} `json:"content"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://api.juejin.cn/content_api/v1/content/article_rank?category_id=1&type=hot&spider=0", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		id := k.Content.ContentID
		if id == "" {
			continue
		}
		out = append(out, Item{ID: id, Title: k.Content.Title, URL: "https://juejin.cn/post/" + id})
	}
	return out, nil
}
