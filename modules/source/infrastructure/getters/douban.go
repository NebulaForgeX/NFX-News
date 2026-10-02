package getters

import "context"

func (r *Registry) douban(ctx context.Context) ([]Item, error) {
	var res struct {
		Items []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"items"`
	}
	if err := r.getJSON(ctx, "https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie", map[string]string{
		"Referer": "https://movie.douban.com/",
		"Accept":  "application/json, text/plain, */*",
	}, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Items))
	for _, m := range res.Items {
		out = append(out, Item{ID: m.ID, Title: m.Title, URL: "https://movie.douban.com/subject/" + m.ID})
	}
	return out, nil
}
