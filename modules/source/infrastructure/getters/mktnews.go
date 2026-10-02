package getters

import (
	"context"
	"regexp"
	"sort"
)

func (r *Registry) mktnewsFlash(ctx context.Context) ([]Item, error) {
	var res struct {
		Data []struct {
			Name  string `json:"name"`
			Child []struct {
				FlashList []struct {
					ID   string `json:"id"`
					Time string `json:"time"`
					Data struct {
						Content string `json:"content"`
						Title   string `json:"title"`
					} `json:"data"`
				} `json:"flash_list"`
			} `json:"child"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://api.mktnews.net/api/flash/host", nil, &res); err != nil {
		return nil, err
	}
	typeMap := map[string]string{"policy": "Policy", "AI": "AI", "financial": "Financial"}
	re := regexp.MustCompile(`^【([^】]*)】(.*)$`)
	type row struct {
		item Item
		time string
	}
	var rows []row
	for _, cat := range res.Data {
		label, ok := typeMap[cat.Name]
		if !ok {
			continue
		}
		for _, child := range cat.Child {
			for _, item := range child.FlashList {
				title := item.Data.Title
				if title == "" {
					if m := re.FindStringSubmatch(item.Data.Content); len(m) > 1 {
						title = m[1]
					} else {
						title = item.Data.Content
					}
				}
				rows = append(rows, row{
					time: item.Time,
					item: Item{
						ID: item.ID, Title: title, URL: "https://mktnews.net/flashDetail.html?id=" + item.ID,
						Extra: map[string]any{"info": label, "hover": item.Data.Content},
					},
				})
			}
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].time > rows[j].time })
	out := make([]Item, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.item)
	}
	return out, nil
}
