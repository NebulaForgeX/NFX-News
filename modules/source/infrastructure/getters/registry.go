package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

type Item struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	URL       string         `json:"url"`
	MobileURL string         `json:"mobileUrl,omitempty"`
	PubDate   int64          `json:"pubDate,omitempty"`
	Extra     map[string]any `json:"extra,omitempty"`
}

type Meta struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Column     string `json:"column"`
	Home       string `json:"home"`
	Color      string `json:"color"`
	IntervalMS int    `json:"interval"`
	Type       string `json:"type"`
	Redirect   string `json:"redirect"`
}

type Getter func(ctx context.Context) ([]Item, error)

type Registry struct {
	getters map[string]Getter
	catalog map[string]Meta
	http    *http.Client
}

func NewRegistry() *Registry {
	r := &Registry{
		getters: map[string]Getter{},
		catalog: map[string]Meta{},
		http:    &http.Client{Timeout: 20 * time.Second},
	}
	r.registerBuiltin()
	return r
}

func (r *Registry) LoadCatalog(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	raw := map[string]map[string]any{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for id, v := range raw {
		meta := Meta{ID: id}
		if s, ok := v["name"].(string); ok {
			meta.Name = s
		}
		if s, ok := v["title"].(string); ok {
			meta.Title = s
		}
		if s, ok := v["column"].(string); ok {
			meta.Column = s
		}
		if s, ok := v["home"].(string); ok {
			meta.Home = s
		}
		if s, ok := v["color"].(string); ok {
			meta.Color = s
		}
		if s, ok := v["type"].(string); ok {
			meta.Type = s
		}
		if s, ok := v["redirect"].(string); ok {
			meta.Redirect = s
		}
		if n, ok := v["interval"].(float64); ok {
			meta.IntervalMS = int(n)
		}
		r.catalog[id] = meta
	}
	return nil
}

func (r *Registry) List() []Meta {
	out := make([]Meta, 0, len(r.catalog))
	for _, m := range r.catalog {
		out = append(out, m)
	}
	return out
}

func (r *Registry) Meta(id string) (Meta, bool) {
	m, ok := r.catalog[id]
	if ok && m.Redirect != "" {
		if t, ok2 := r.catalog[m.Redirect]; ok2 {
			return t, true
		}
	}
	return m, ok
}

func (r *Registry) Fetch(ctx context.Context, id string) ([]Item, error) {
	meta, ok := r.Meta(id)
	if !ok {
		return nil, fmt.Errorf("unknown source %s", id)
	}
	fn, ok := r.getters[id]
	if !ok {
		fn, ok = r.getters[meta.ID]
	}
	if !ok {
		return nil, fmt.Errorf("no getter for %s", id)
	}
	return fn(ctx)
}

func (r *Registry) register(id string, fn Getter) { r.getters[id] = fn }

func (r *Registry) registerBuiltin() {
	r.register("zhihu", r.zhihu)
	r.register("v2ex", r.v2ex)
	r.register("v2ex-share", r.v2ex)
	r.register("hackernews", r.hackerNews)
	r.register("github", r.githubTrending)
	r.register("github-trending-today", r.githubTrending)
	r.register("weibo", r.weibo)
	r.register("bilibili", r.bilibiliHot)
	r.register("bilibili-hot-search", r.bilibiliHot)
	r.register("bilibili-hot-video", r.bilibiliVideo)
	r.register("bilibili-ranking", r.bilibiliRank)
	r.register("ithome", r.rss("https://www.ithome.com/rss/"))
	r.register("solidot", r.rss("https://www.solidot.org/index.rss"))
	r.register("sspai", r.rss("https://sspai.com/feed"))
	r.register("hupu", r.rssHub("/bbs/hupu/all-gambia"))
	r.register("douban", r.rssHub("/douban/movie/playing"))
	r.register("producthunt", r.rss("https://www.producthunt.com/feed"))
	r.register("linuxdo", r.rss("https://linux.do/latest.rss"))
	r.register("zaobao", r.rssHub("/zaobao/realtime/china"))
	r.register("thepaper", r.rssHub("/thepaper/featured"))
	r.register("kaopu", r.rss("https://kaopu.news/rss"))
	r.register("gelonghui", r.rssHub("/gelonghui/hot-article"))
	r.register("fastbull", r.rssHub("/fastbull/express"))
	r.register("jin10", r.rssHub("/jin10"))
	r.register("wallstreetcn", r.rssHub("/wallstreetcn/news"))
	r.register("xueqiu", r.rssHub("/xueqiu/hot"))
	r.register("cls", r.rssHub("/cls/telegraph"))
	r.register("mktnews", r.rssHub("/mktnews/flash"))
	r.register("juejin", r.rssHub("/juejin/hot"))
	r.register("nowcoder", r.rssHub("/nowcoder/recommend"))
	r.register("pcbeta", r.rssHub("/pcbeta/topic"))
	r.register("chongbuluo", r.rssHub("/chongbuluo/hot"))
	r.register("sputniknewscn", r.rssHub("/sputniknewscn"))
	r.register("cankaoxiaoxi", r.rssHub("/cankaoxiaoxi"))
	r.register("ifeng", r.rssHub("/ifeng/hot"))
	r.register("toutiao", r.rssHub("/toutiao/hot"))
	r.register("tieba", r.rssHub("/tieba/topic"))
	r.register("kuaishou", r.rssHub("/kuaishou/hot"))
	r.register("douyin", r.rssHub("/douyin/hot"))
	r.register("baidu", r.baidu)
	r.register("36kr", r.rssHub("/36kr/hot-list"))
	r.register("smzdm", r.rssHub("/smzdm/ranking/pinlei/11/3"))
	r.register("steam", r.rssHub("/steam/search/hot"))
	r.register("coolapk", r.rssHub("/coolapk/hot"))
	r.register("ghxi", r.rssHub("/ghxi"))
	r.register("tencent-news", r.rssHub("/tencent/news/hot"))
}

func (r *Registry) get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; NFX-News/1.0)")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func (r *Registry) rss(url string) Getter {
	return func(ctx context.Context) ([]Item, error) {
		fp := gofeed.NewParser()
		fp.Client = r.http
		feed, err := fp.ParseURLWithContext(url, ctx)
		if err != nil {
			return nil, err
		}
		items := make([]Item, 0, len(feed.Items))
		for _, it := range feed.Items {
			link := it.Link
			if link == "" {
				link = it.GUID
			}
			item := Item{ID: link, Title: it.Title, URL: link}
			if it.PublishedParsed != nil {
				item.PubDate = it.PublishedParsed.Unix()
			}
			items = append(items, item)
		}
		return items, nil
	}
}

func (r *Registry) rssHub(route string) Getter {
	return r.rss("https://rsshub.rssforever.com" + route + "?format=xml")
}

func (r *Registry) zhihu(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.zhihu.com/api/v3/feed/topstory/hot-list-web?limit=20&desktop=true", nil)
	if err != nil {
		return nil, err
	}
	var res struct {
		Data []struct {
			Target struct {
				TitleArea struct {
					Text string `json:"text"`
				} `json:"title_area"`
				ExcerptArea struct {
					Text string `json:"text"`
				} `json:"excerpt_area"`
				MetricsArea struct {
					Text string `json:"text"`
				} `json:"metrics_area"`
				Link struct {
					URL string `json:"url"`
				} `json:"link"`
			} `json:"target"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(\d+)$`)
	out := make([]Item, 0, len(res.Data))
	for _, k := range res.Data {
		id := k.Target.Link.URL
		if m := re.FindString(k.Target.Link.URL); m != "" {
			id = m
		}
		out = append(out, Item{
			ID:    id,
			Title: k.Target.TitleArea.Text,
			URL:   k.Target.Link.URL,
			Extra: map[string]any{"info": k.Target.MetricsArea.Text, "hover": k.Target.ExcerptArea.Text},
		})
	}
	return out, nil
}

func (r *Registry) v2ex(ctx context.Context) ([]Item, error) {
	var all []Item
	for _, k := range []string{"create", "ideas", "programmer", "share"} {
		body, err := r.get(ctx, "https://www.v2ex.com/feed/"+k+".json", nil)
		if err != nil {
			continue
		}
		var res struct {
			Items []struct {
				ID            string `json:"id"`
				Title         string `json:"title"`
				URL           string `json:"url"`
				DatePublished string `json:"date_published"`
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

func (r *Registry) hackerNews(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://news.ycombinator.com", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find(".athing").Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		a := s.Find(".titleline a").First()
		title := a.Text()
		if id != "" && title != "" {
			out = append(out, Item{ID: id, Title: title, URL: "https://news.ycombinator.com/item?id=" + id})
		}
	})
	return out, nil
}

func (r *Registry) githubTrending(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://github.com/trending?spoken_language_code=", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find("main .Box article").Each(func(_ int, s *goquery.Selection) {
		a := s.Find("h2 a")
		href, _ := a.Attr("href")
		title := strings.TrimSpace(strings.ReplaceAll(a.Text(), "\n", ""))
		if href != "" && title != "" {
			out = append(out, Item{ID: href, Title: title, URL: "https://github.com" + href})
		}
	})
	return out, nil
}

func (r *Registry) weibo(ctx context.Context) ([]Item, error) {
	url := "https://s.weibo.com/top/summary?cate=realtimehot"
	body, err := r.get(ctx, url, map[string]string{"Referer": url})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find("#pl_top_realtimehot table tbody tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		a := s.Find("td.td-02 a").First()
		title := strings.TrimSpace(a.Text())
		href, _ := a.Attr("href")
		if title != "" && href != "" && !strings.Contains(href, "javascript") {
			out = append(out, Item{ID: title, Title: title, URL: "https://s.weibo.com" + href})
		}
	})
	return out, nil
}

func (r *Registry) baidu(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://top.baidu.com/board?tab=realtime", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(?s)<!--s-data:(.*?)-->`)
	m := re.FindSubmatch(body)
	if len(m) < 2 {
		return nil, fmt.Errorf("baidu payload not found")
	}
	var res struct {
		Data struct {
			Cards []struct {
				Content []struct {
					IsTop  bool   `json:"isTop"`
					Word   string `json:"word"`
					RawURL string `json:"rawUrl"`
					Desc   string `json:"desc"`
				} `json:"content"`
			} `json:"cards"`
		} `json:"data"`
	}
	if err := json.Unmarshal(m[1], &res); err != nil {
		return nil, err
	}
	var out []Item
	if len(res.Data.Cards) == 0 {
		return out, nil
	}
	for _, k := range res.Data.Cards[0].Content {
		if k.IsTop {
			continue
		}
		out = append(out, Item{ID: k.RawURL, Title: k.Word, URL: k.RawURL, Extra: map[string]any{"hover": k.Desc}})
	}
	return out, nil
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
		} `json:"list"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.List))
	for _, k := range res.List {
		out = append(out, Item{
			ID:    k.Keyword,
			Title: k.ShowName,
			URL:   "https://search.bilibili.com/all?keyword=" + k.Keyword,
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

func (r *Registry) bilibiliList(ctx context.Context, url string) ([]Item, error) {
	body, err := r.get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	var res struct {
		Data struct {
			List []struct {
				BVID  string `json:"bvid"`
				Title string `json:"title"`
				Desc  string `json:"desc"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(res.Data.List))
	for _, v := range res.Data.List {
		out = append(out, Item{ID: v.BVID, Title: v.Title, URL: "https://www.bilibili.com/video/" + v.BVID, Extra: map[string]any{"hover": v.Desc}})
	}
	return out, nil
}
