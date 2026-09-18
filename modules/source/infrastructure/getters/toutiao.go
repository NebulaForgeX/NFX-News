package getters

import "context"

func (r *Registry) toutiao(ctx context.Context) ([]Item, error) {
	var res struct {
		Data []struct {
			ClusterIDStr string `json:"ClusterIdStr"`
			Title        string `json:"Title"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		out = append(out, Item{ID: k.ClusterIDStr, Title: k.Title, URL: "https://www.toutiao.com/trending/" + k.ClusterIDStr + "/"})
	}
	return out, nil
}
