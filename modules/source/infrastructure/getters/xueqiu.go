package getters

import (
	"context"
	"fmt"
)

func (r *Registry) xueqiuHotstock(ctx context.Context) ([]Item, error) {
	reqBody, err := r.get(ctx, "https://xueqiu.com/hq", nil)
	_ = reqBody
	if err != nil {
		return nil, err
	}
	// Cookie jar is not shared; fetch with Referer which is enough for many responses.
	var res struct {
		Data struct {
			Items []struct {
				Code     string  `json:"code"`
				Name     string  `json:"name"`
				Percent  float64 `json:"percent"`
				Exchange string  `json:"exchange"`
				Ad       int     `json:"ad"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://stock.xueqiu.com/v5/stock/hot_stock/list.json?size=30&_type=10&type=10", map[string]string{
		"Referer": "https://xueqiu.com/hq",
	}, &res); err != nil {
		return nil, err
	}
	var out []Item
	for _, k := range res.Data.Items {
		if k.Ad != 0 {
			continue
		}
		out = append(out, Item{
			ID: k.Code, Title: k.Name, URL: "https://xueqiu.com/s/" + k.Code,
			Extra: map[string]any{"info": fmt.Sprintf("%v%% %s", k.Percent, k.Exchange)},
		})
	}
	return out, nil
}
