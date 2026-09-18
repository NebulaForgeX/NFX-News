package getters

import (
	"context"
	"strings"

	"nfxnews/modules/source/infrastructure/getters/coolapk"
)

func (r *Registry) coolapk(ctx context.Context) ([]Item, error) {
	endpoint := "https://api.coolapk.com/v6/page/dataList?url=%2Ffeed%2FstatList%3FcacheExpires%3D300%26statType%3Dday%26sortField%3Ddetailnum%26title%3D%E4%BB%8A%E6%97%A5%E7%83%AD%E9%97%A8&title=%E4%BB%8A%E6%97%A5%E7%83%AD%E9%97%A8&subTitle=&page=1"
	var res struct {
		Data []struct {
			ID          string `json:"id"`
			Message     string `json:"message"`
			EditorTitle string `json:"editor_title"`
			URL         string `json:"url"`
			TargetRow   struct {
				SubTitle string `json:"subTitle"`
			} `json:"targetRow"`
		} `json:"data"`
	}
	headers := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"X-App-Id":         "com.coolapk.market",
		"X-App-Token":      coolapk.Token(),
		"X-Sdk-Int":        "29",
		"X-Sdk-Locale":     "zh-CN",
		"X-App-Version":    "11.0",
		"X-Api-Version":    "11",
		"X-App-Code":       "2101202",
		"User-Agent":       "Dalvik/2.1.0 (Linux; U; Android 10; Redmi K30 5G MIUI/V12.0.3.0.QGICMXM) (#Build; Redmi; Redmi K30 5G; QKQ1.191222.002 test-keys; 10) +CoolMarket/11.0-2101202",
	}
	if err := r.getJSON(ctx, endpoint, headers, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		if k.ID == "" {
			continue
		}
		title := k.EditorTitle
		if title == "" && k.Message != "" {
			title = strings.TrimSpace(strings.Split(k.Message, "\n")[0])
		}
		if title == "" {
			title = k.ID
		}
		out = append(out, Item{ID: k.ID, Title: title, URL: "https://www.coolapk.com" + k.URL, Extra: map[string]any{"info": k.TargetRow.SubTitle}})
	}
	return out, nil
}
