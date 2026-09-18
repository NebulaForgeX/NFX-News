package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (r *Registry) juejin(ctx context.Context) ([]Item, error) {
	var res struct {
		Data []struct {
			Content struct {
				Title     string `json:"title"`
				ContentID string `json:"content_id"`
			} `json:"content"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://api.juejin.cn/content_api/v1/content/article_rank?category_id=1&type=hot&spider=0", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		id := k.Content.ContentID
		if id == "" {
			continue
		}
		out = append(out, Item{ID: id, Title: k.Content.Title, URL: "https://juejin.cn/post/" + id})
	}
	return out, nil
}

func (r *Registry) sspai(ctx context.Context) ([]Item, error) {
	url := fmt.Sprintf("https://sspai.com/api/v1/article/tag/page/get?limit=30&offset=0&created_at=%d&tag=热门文章&released=false", time.Now().UnixMilli())
	var res struct {
		Data []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, url, nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		id := strconv.Itoa(k.ID)
		out = append(out, Item{ID: id, Title: k.Title, URL: "https://sspai.com/post/" + id})
	}
	return out, nil
}

func (r *Registry) nowcoder(ctx context.Context) ([]Item, error) {
	url := fmt.Sprintf("https://gw-c.nowcoder.com/api/sparta/hot-search/top-hot-pc?size=20&_=%d&t=", time.Now().UnixMilli())
	var res struct {
		Data struct {
			Result []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				Type  int    `json:"type"`
				UUID  string `json:"uuid"`
			} `json:"result"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, url, nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.Result))
	for _, k := range res.Data.Result {
		id, link := k.ID, "https://www.nowcoder.com/discuss/"+k.ID
		if k.Type == 74 && k.UUID != "" {
			id, link = k.UUID, "https://www.nowcoder.com/feed/main/detail/"+k.UUID
		}
		out = append(out, Item{ID: id, Title: k.Title, URL: link})
	}
	return out, nil
}

func (r *Registry) douban(ctx context.Context) ([]Item, error) {
	var res struct {
		Items []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"items"`
	}
	if err := r.getJSON(ctx, "https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie", map[string]string{
		"Referer": "https://movie.douban.com/",
		"Accept":  "application/json, text/plain, */*",
	}, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Items))
	for _, m := range res.Items {
		out = append(out, Item{ID: m.ID, Title: m.Title, URL: "https://movie.douban.com/subject/" + m.ID})
	}
	return out, nil
}

func (r *Registry) hupu(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://bbs.hupu.com/topic-daily-hot", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<a href="(/[^"]+?\.html)"[^>]*class="p-title"[^>]*>([^<]+)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	for _, m := range matches {
		path, title := m[1], strings.TrimSpace(m[2])
		out = append(out, Item{ID: path, Title: title, URL: "https://bbs.hupu.com" + path})
	}
	return out, nil
}

func (r *Registry) toutiao(ctx context.Context) ([]Item, error) {
	var res struct {
		Data []struct {
			ClusterIDStr string `json:"ClusterIdStr"`
			Title        string `json:"Title"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		out = append(out, Item{ID: k.ClusterIDStr, Title: k.Title, URL: "https://www.toutiao.com/trending/" + k.ClusterIDStr + "/"})
	}
	return out, nil
}

func (r *Registry) thepaper(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			HotNews []struct {
				ContID string `json:"contId"`
				Name   string `json:"name"`
			} `json:"hotNews"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://cache.thepaper.cn/contentapi/wwwIndex/rightSidebar", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.HotNews))
	for _, k := range res.Data.HotNews {
		out = append(out, Item{
			ID: k.ContID, Title: k.Name,
			URL: "https://www.thepaper.cn/newsDetail_forward_" + k.ContID, MobileURL: "https://m.thepaper.cn/newsDetail_forward_" + k.ContID,
		})
	}
	return out, nil
}

func (r *Registry) ifeng(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.ifeng.com/", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`var\s+allData\s*=\s*(\{[\s\S]*?\});`)
	m := re.FindSubmatch(body)
	if len(m) < 2 {
		return nil, fmt.Errorf("ifeng allData not found")
	}
	var real struct {
		HotNews1 []struct {
			URL   string `json:"url"`
			Title string `json:"title"`
		} `json:"hotNews1"`
	}
	if err := json.Unmarshal(m[1], &real); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(real.HotNews1))
	for _, n := range real.HotNews1 {
		out = append(out, Item{ID: n.URL, Title: n.Title, URL: n.URL})
	}
	return out, nil
}

func (r *Registry) tieba(ctx context.Context) ([]Item, error) {
	var res struct {
		Data struct {
			BangTopic struct {
				TopicList []struct {
					TopicID   string `json:"topic_id"`
					TopicName string `json:"topic_name"`
					TopicURL  string `json:"topic_url"`
				} `json:"topic_list"`
			} `json:"bang_topic"`
		} `json:"data"`
	}
	if err := r.getJSON(ctx, "https://tieba.baidu.com/hottopic/browse/topicList", nil, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.BangTopic.TopicList))
	for _, k := range res.Data.BangTopic.TopicList {
		out = append(out, Item{ID: k.TopicID, Title: k.TopicName, URL: k.TopicURL})
	}
	return out, nil
}

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

func (r *Registry) jin10(ctx context.Context) ([]Item, error) {
	url := fmt.Sprintf("https://www.jin10.com/flash_newest.js?t=%d", time.Now().UnixMilli())
	body, err := r.get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(string(body))
	raw = strings.TrimPrefix(raw, "var newest =")
	raw = strings.TrimPrefix(raw, "var newest=")
	raw = strings.TrimRight(raw, ";")
	var rows []struct {
		ID   string `json:"id"`
		Data struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(rows))
	for _, k := range rows {
		title := k.Data.Title
		if title == "" {
			title = k.Data.Content
		}
		if title == "" {
			continue
		}
		out = append(out, Item{ID: k.ID, Title: title, URL: "https://www.jin10.com/flash/" + k.ID})
	}
	return out, nil
}

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

func (r *Registry) smzdm(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://post.smzdm.com/hot_1/", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<a[^>]*href="(https://post\.smzdm.com/p/[^"]+)"[^>]*>([^<]+)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	seen := map[string]struct{}{}
	for _, m := range matches {
		u, title := m[1], strings.TrimSpace(m[2])
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, Item{ID: u, Title: title, URL: u})
	}
	return out, nil
}
