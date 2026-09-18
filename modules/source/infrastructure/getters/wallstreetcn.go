package getters

import "context"

func (r *Registry) wallstreetcnQuick(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			Items []struct {
				URI          string `json:"uri"`
				ID           int    `json:"id"`
				Title        string `json:"title"`
				ContentText  string `json:"content_text"`
				DisplayTime  int64  `json:"display_time"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://api-one.wallstcn.com/apiv1/content/lives?channel=global-channel&limit=30", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.Items))
	for _, k := range res.Data.Items {
		title := k.Title
		if title == "" {
			title = k.ContentText
		}
		out = append(out, Item{ID: itoa(k.ID), Title: title, URL: k.URI, Extra: map[string]any{"date": k.DisplayTime * 1000}})
	}
	return out, nil
}

func (r *Registry) wallstreetcnNews(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			Items []struct {
				ResourceType string `json:"resource_type"`
				Resource     struct {
					URI          string `json:"uri"`
					ID           int    `json:"id"`
					Title        string `json:"title"`
					ContentShort string `json:"content_short"`
					DisplayTime  int64  `json:"display_time"`
					Type         string `json:"type"`
				} `json:"resource"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://api-one.wallstcn.com/apiv1/content/information-flow?channel=global-channel&accept=article&limit=30", nil, &res); err != nil {
		return nil, err
	}
	var out []Item
	for _, k := range res.Data.Items {
		h := k.Resource
		if k.ResourceType == "theme" || k.ResourceType == "ad" || h.Type == "live" || h.URI == "" {
			continue
		}
		title := h.Title
		if title == "" {
			title = h.ContentShort
		}
		out = append(out, Item{ID: itoa(h.ID), Title: title, URL: h.URI, Extra: map[string]any{"date": h.DisplayTime * 1000}})
	}
	return out, nil
}

func (r *Registry) wallstreetcnHot(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			DayItems []struct {
				URI   string `json:"uri"`
				ID    int    `json:"id"`
				Title string `json:"title"`
			} `json:"day_items"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://api-one.wallstcn.com/apiv1/content/articles/hot?period=all", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.DayItems))
	for _, h := range res.Data.DayItems {
		out = append(out, Item{ID: itoa(h.ID), Title: h.Title, URL: h.URI})
	}
	return out, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
