package getters

import (
	"context"
	"encoding/json"
)

func (r *Registry) v2ex(ctx context.Context) ([]Item, error) {
	var all []Item
	for _, k := range []string{"create", "ideas", "programmer", "share"} {
		body, err := r.get(ctx, "https://www.v2ex.com/feed/"+k+".json", nil)
		if err != nil {
			continue
		}
		var res struct {
			Items []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"items"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			continue
		}
		for _, it := range res.Items {
			all = append(all, Item{ID: it.ID, Title: it.Title, URL: it.URL})
		}
	}
	return all, nil
}
