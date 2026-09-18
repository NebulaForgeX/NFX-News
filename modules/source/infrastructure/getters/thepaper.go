package getters

import "context"

func (r *Registry) thepaper(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			HotNews []struct {
				ContID string `json:"contId"`
				Name   string `json:"name"`
			} `json:"hotNews"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://cache.thepaper.cn/contentapi/wwwIndex/rightSidebar", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.HotNews))
	for _, k := range res.Data.HotNews {
		out = append(out, Item{
			ID: k.ContID, Title: k.Name,
			URL: "https://www.thepaper.cn/newsDetail_forward_" + k.ContID, MobileURL: "https://m.thepaper.cn/newsDetail_forward_" + k.ContID,
		})
	}
	return out, nil
}
