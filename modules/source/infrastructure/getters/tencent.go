package getters

import "context"

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
