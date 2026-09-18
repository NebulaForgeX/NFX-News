package getters

import "context"

func (r *Registry) cankaoxiaoxi(ctx context.Context) ([]Item, error) {
	channels := []string{"zhongguo", "guandian", "gj"}
	type row struct {
		Data struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			URL         string `json:"url"`
			PublishTime string `json:"publishTime"`
		} `json:"data"`
	}
	var out []Item
	for _, ch := range channels {
		var res struct {
			List []row `json:"list"`
		}
		if err := r.getJSON(ctx, "https://china.cankaoxiaoxi.com/json/channel/"+ch+"/list.json", nil, &res); err != nil {
			return nil, err
		}
		for _, k := range res.List {
			out = append(out, Item{ID: k.Data.ID, Title: k.Data.Title, URL: k.Data.URL, Extra: map[string]any{"date": k.Data.PublishTime}})
		}
	}
	return out, nil
}
