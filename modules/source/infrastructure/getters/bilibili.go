package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

func formatNumber(num int) string {
	if num >= 10000 {
		return fmt.Sprintf("%dw+", num/10000)
	}
	return fmt.Sprintf("%d", num)
}

func (r *Registry) bilibiliHot(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://s.search.bilibili.com/main/hotword?limit=30", nil)
	if err != nil {
		return nil, err
	}
	var res struct {
		List []struct {
			Keyword  string `json:"keyword"`
			ShowName string `json:"show_name"`
			Icon     string `json:"icon"`
		} `json:"list"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.List))
	for _, k := range res.List {
		out = append(out, Item{
			ID: k.Keyword, Title: k.ShowName,
			URL:   "https://search.bilibili.com/all?keyword=" + url.QueryEscape(k.Keyword),
			Extra: map[string]any{"icon": k.Icon},
		})
	}
	return out, nil
}

func (r *Registry) bilibiliVideo(ctx context.Context) ([]Item, error) {
	return r.bilibiliList(ctx, "https://api.bilibili.com/x/web-interface/popular")
}

func (r *Registry) bilibiliRank(ctx context.Context) ([]Item, error) {
	return r.bilibiliList(ctx, "https://api.bilibili.com/x/web-interface/ranking/v2")
}

func (r *Registry) bilibiliList(ctx context.Context, api string) ([]Item, error) {
	body, err := r.get(ctx, api, nil)
	if err != nil {
		return nil, err
	}
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				BVID  string `json:"bvid"`
				Title string `json:"title"`
				Pic   string `json:"pic"`
				Desc  string `json:"desc"`
				Owner struct {
					Name string `json:"name"`
				} `json:"owner"`
				Stat struct {
					View int `json:"view"`
					Like int `json:"like"`
				} `json:"stat"`
				Pubdate int64 `json:"pubdate"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("bilibili API error: %d", res.Code)
	}
	out := make([]Item, 0, len(res.Data.List))
	for _, v := range res.Data.List {
		extra := map[string]any{
			"info":  fmt.Sprintf("%s · %s观看 · %s点赞", v.Owner.Name, formatNumber(v.Stat.View), formatNumber(v.Stat.Like)),
			"hover": v.Desc,
		}
		if v.Pic != "" {
			extra["icon"] = map[string]any{"url": v.Pic, "scale": 1}
		}
		out = append(out, Item{ID: v.BVID, Title: v.Title, URL: "https://www.bilibili.com/video/" + v.BVID, PubDate: v.Pubdate * 1000, Extra: extra})
	}
	return out, nil
}
