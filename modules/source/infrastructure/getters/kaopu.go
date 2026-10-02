package getters

import "context"

func (r *Registry) kaopu(ctx context.Context) ([]Item, error) {
	var rows []struct {
		Description string `json:"description"`
		Link        string `json:"link"`
		Publisher   string `json:"publisher"`
		Title       string `json:"title"`
	}
	if err := r.getJSON(ctx, "https://kaopustorage.blob.core.windows.net/news-prod/news_list_hans_0.json", nil, &rows); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(rows))
	for _, k := range rows {
		if k.Publisher == "财新" || k.Publisher == "公视" {
			continue
		}
		out = append(out, Item{ID: k.Link, Title: k.Title, URL: k.Link, Extra: map[string]any{"hover": k.Description, "info": k.Publisher}})
	}
	return out, nil
}
