package getters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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
	IntervalMS int    `json:"interval_ms"`
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

func (r *Registry) Interval(id string) time.Duration {
	m, ok := r.Meta(id)
	if !ok || m.IntervalMS <= 0 {
		return 10 * time.Minute
	}
	return time.Duration(m.IntervalMS) * time.Millisecond
}

func (r *Registry) IDs() []string {
	ids := make([]string, 0, len(r.catalog))
	for id, m := range r.catalog {
		if m.Redirect != "" {
			continue
		}
		ids = append(ids, id)
	}
	return ids
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

func (r *Registry) nativeOrHub(native Getter, hubRoute string) Getter {
	return func(ctx context.Context) ([]Item, error) {
		items, err := native(ctx)
		if err == nil && len(items) > 0 {
			return items, nil
		}
		if hubRoute == "" {
			return items, err
		}
		return r.rssHub(hubRoute)(ctx)
	}
}

func (r *Registry) nativeOrRSS(native Getter, rssURL string) Getter {
	return func(ctx context.Context) ([]Item, error) {
		items, err := native(ctx)
		if err == nil && len(items) > 0 {
			return items, nil
		}
		return r.rss(rssURL)(ctx)
	}
}

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
	r.register("baidu", r.baidu)
	r.register("ithome", r.ithome)
	r.register("solidot", r.solidot)
	r.register("sspai", r.nativeOrRSS(r.sspai, "https://sspai.com/feed"))
	r.register("hupu", r.nativeOrHub(r.hupu, "/bbs/hupu/all-gambia"))
	r.register("douban", r.nativeOrHub(r.douban, "/douban/movie/playing"))
	r.register("producthunt", r.producthunt)
	r.register("zaobao", r.nativeOrHub(r.zaobao, "/zaobao/realtime/china"))
	r.register("thepaper", r.nativeOrHub(r.thepaper, "/thepaper/featured"))
	r.register("kaopu", r.nativeOrHub(r.kaopu, ""))
	r.register("gelonghui", r.nativeOrHub(r.gelonghui, "/gelonghui/hot-article"))
	r.register("jin10", r.nativeOrHub(r.jin10, "/jin10"))
	r.register("juejin", r.nativeOrHub(r.juejin, "/juejin/hot"))
	r.register("nowcoder", r.nativeOrHub(r.nowcoder, "/nowcoder/recommend"))
	r.register("sputniknewscn", r.nativeOrHub(r.sputnik, "/sputniknewscn"))
	r.register("cankaoxiaoxi", r.nativeOrHub(r.cankaoxiaoxi, "/cankaoxiaoxi"))
	r.register("ifeng", r.nativeOrHub(r.ifeng, "/ifeng/hot"))
	r.register("toutiao", r.nativeOrHub(r.toutiao, "/toutiao/hot"))
	r.register("tieba", r.nativeOrHub(r.tieba, "/tieba/topic"))
	r.register("kuaishou", r.nativeOrHub(r.kuaishou, "/kuaishou/hot"))
	r.register("douyin", r.nativeOrHub(r.douyin, "/douyin/hot"))
	r.register("smzdm", r.nativeOrHub(r.smzdm, "/smzdm/ranking/pinlei/11/3"))
	r.register("steam", r.nativeOrHub(r.steam, "/steam/search/hot"))
	r.register("coolapk", r.nativeOrHub(r.coolapk, "/coolapk/hot"))
	r.register("ghxi", r.nativeOrHub(r.ghxi, "/ghxi"))
	r.register("36kr-quick", r.nativeOrHub(r.kr36Quick, "/36kr/hot-list"))
	r.register("36kr", r.nativeOrHub(r.kr36Quick, "/36kr/hot-list"))
	r.register("wallstreetcn-quick", r.nativeOrHub(r.wallstreetcnQuick, "/wallstreetcn/live"))
	r.register("wallstreetcn-news", r.nativeOrHub(r.wallstreetcnNews, "/wallstreetcn/news"))
	r.register("wallstreetcn-hot", r.nativeOrHub(r.wallstreetcnHot, "/wallstreetcn/hot"))
	r.register("wallstreetcn", r.nativeOrHub(r.wallstreetcnQuick, "/wallstreetcn/live"))
	r.register("cls-telegraph", r.nativeOrHub(r.clsTelegraph, "/cls/telegraph"))
	r.register("cls-depth", r.nativeOrHub(r.clsDepth, "/cls/depth"))
	r.register("cls-hot", r.nativeOrHub(r.clsHot, "/cls/hot"))
	r.register("cls", r.nativeOrHub(r.clsTelegraph, "/cls/telegraph"))
	r.register("linuxdo-hot", r.nativeOrHub(r.linuxdoHot, ""))
	r.register("linuxdo-latest", r.nativeOrHub(r.linuxdoLatest, ""))
	r.register("linuxdo", r.nativeOrHub(r.linuxdoLatest, ""))
	r.register("mktnews-flash", r.nativeOrHub(r.mktnewsFlash, "/mktnews/flash"))
	r.register("mktnews", r.nativeOrHub(r.mktnewsFlash, "/mktnews/flash"))
	r.register("pcbeta-windows", r.pcbetaWindows)
	r.register("pcbeta-windows11", r.pcbetaWindows11)
	r.register("pcbeta", r.pcbetaWindows11)
	r.register("tencent-hot", r.nativeOrHub(r.tencentHot, "/tencent/news/hot"))
	r.register("tencent-news", r.nativeOrHub(r.tencentHot, "/tencent/news/hot"))
	r.register("xueqiu-hotstock", r.nativeOrHub(r.xueqiuHotstock, "/xueqiu/hot"))
	r.register("xueqiu", r.nativeOrHub(r.xueqiuHotstock, "/xueqiu/hot"))
	r.register("fastbull-express", r.nativeOrHub(r.fastbullExpress, "/fastbull/express"))
	r.register("fastbull-news", r.nativeOrHub(r.fastbullNews, "/fastbull/news"))
	r.register("fastbull", r.nativeOrHub(r.fastbullExpress, "/fastbull/express"))
	r.register("chongbuluo-hot", r.nativeOrHub(r.chongbuluoHot, "/chongbuluo/hot"))
	r.register("chongbuluo-latest", r.rss("https://www.chongbuluo.com/forum.php?mod=rss&view=newthread"))
	r.register("chongbuluo", r.rss("https://www.chongbuluo.com/forum.php?mod=rss&view=newthread"))
}
