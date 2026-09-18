package getters

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
)

func clsSign(extra map[string]string) string {
	params := map[string]string{"appName": "CailianpressWeb", "os": "web", "sv": "7.7.5"}
	for k, v := range extra {
		params[k] = v
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	vals := url.Values{}
	for _, k := range keys {
		vals.Set(k, params[k])
	}
	sum := sha1.Sum([]byte(vals.Encode()))
	md := md5.Sum([]byte(hex.EncodeToString(sum[:])))
	vals.Set("sign", hex.EncodeToString(md[:]))
	return vals.Encode()
}

type clsItem struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Brief    string `json:"brief"`
	ShareURL string `json:"shareurl"`
	CTime    int64  `json:"ctime"`
	IsAd     int    `json:"is_ad"`
}

func (k clsItem) toItem() Item {
	title := k.Title
	if title == "" {
		title = k.Brief
	}
	return Item{
		ID: itoa(k.ID), Title: title, URL: fmt.Sprintf("https://www.cls.cn/detail/%d", k.ID),
		MobileURL: k.ShareURL, PubDate: k.CTime * 1000,
	}
}

func (r *Registry) clsTelegraph(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			RollData []clsItem `json:"roll_data"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://www.cls.cn/nodeapi/updateTelegraphList?"+clsSign(nil), nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.RollData))
	for _, k := range res.Data.RollData {
		if k.IsAd != 0 {
			continue
		}
		out = append(out, k.toItem())
	}
	return out, nil
}

func (r *Registry) clsDepth(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			DepthList []clsItem `json:"depth_list"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://www.cls.cn/v3/depth/home/assembled/1000?"+clsSign(nil), nil, &res); err != nil {
		return nil, err
	}
	sort.Slice(res.Data.DepthList, func(i, j int) bool { return res.Data.DepthList[i].CTime > res.Data.DepthList[j].CTime })
	out := make([]Item, 0, len(res.Data.DepthList))
	for _, k := range res.Data.DepthList {
		out = append(out, k.toItem())
	}
	return out, nil
}

func (r *Registry) clsHot(ctx context.Context) ([]Item, error) {
	var res struct {
		Data []clsItem `json:"data"`
	}
	if err := r.getJSON(ctx, "https://www.cls.cn/v2/article/hot/list?"+clsSign(nil), nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		out = append(out, k.toItem())
	}
	return out, nil
}
