package getters

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
)

func (r *Registry) kuaishou(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.kuaishou.com/?isHome=1", nil)
	if err != nil {
		return nil, err
	}
	raw, err := extractJSONObject(string(body), "window.__APOLLO_STATE__")
	if err != nil {
		return nil, err
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, err
	}
	dc, _ := root["defaultClient"].(map[string]any)
	rq, _ := dc["ROOT_QUERY"].(map[string]any)
	hotQ, _ := rq[`visionHotRank({"page":"home"})`].(map[string]any)
	hotRankID, _ := hotQ["id"].(string)
	hotRank, _ := dc[hotRankID].(map[string]any)
	rows, _ := hotRank["items"].([]any)
	out := make([]Item, 0, len(rows))
	for _, row := range rows {
		item, _ := row.(map[string]any)
		id, _ := item["id"].(string)
		hotItem, _ := dc[id].(map[string]any)
		if hotItem == nil {
			continue
		}
		if tag, _ := hotItem["tagType"].(string); tag == "置顶" {
			continue
		}
		name, _ := hotItem["name"].(string)
		if name == "" {
			continue
		}
		word := strings.TrimPrefix(id, "VisionHotRankItem:")
		out = append(out, Item{
			ID: word, Title: name,
			URL: "https://www.kuaishou.com/search/video?searchKey=" + url.QueryEscape(name),
		})
	}
	return out, nil
}
