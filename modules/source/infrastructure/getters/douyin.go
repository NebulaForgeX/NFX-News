package getters

import (
	"context"
	"net/http"
	"strings"
)

func (r *Registry) douyin(ctx context.Context) ([]Item, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.douyin.com/passport/general/login_guiding_strategy/?aid=6383", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}
	cookie := strings.Join(resp.Header.Values("Set-Cookie"), "; ")
	resp.Body.Close()
	var res struct {
		Data struct {
			WordList []struct {
				SentenceID string `json:"sentence_id"`
				Word       string `json:"word"`
			} `json:"word_list"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://www.douyin.com/aweme/v1/web/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web&detail_list=1", map[string]string{"Cookie": cookie}, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.WordList))
	for _, k := range res.Data.WordList {
		out = append(out, Item{ID: k.SentenceID, Title: k.Word, URL: "https://www.douyin.com/hot/" + k.SentenceID})
	}
	return out, nil
}
