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

func (r *Registry) tencentHot(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			Tabs []struct {
				ArticleList []struct {
					ID       string `json:"id"`
					Title    string `json:"title"`
					Desc     string `json:"desc"`
					LinkInfo struct {
						URL string `json:"url"`
					} `json:"link_info"`
				} `json:"articleList"`
			} `json:"tabs"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://i.news.qq.com/web_backend/v2/getTagInfo?tagId=aEWqxLtdgmQ%3D", map[string]string{
		"Referer": "https://news.qq.com/",
	}, &res); err != nil {
		return nil, err
	}
	if len(res.Data.Tabs) == 0 {
		return nil, nil
	}
	var out []Item
	for _, news := range res.Data.Tabs[0].ArticleList {
		out = append(out, Item{ID: news.ID, Title: news.Title, URL: news.LinkInfo.URL, Extra: map[string]any{"hover": news.Desc}})
	}
	return out, nil
}
