package getters

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func decodeGBK(raw []byte) string {
	out, err := io.ReadAll(transform.NewReader(bytes.NewReader(raw), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return string(raw)
	}
	return string(out)
}

func extractJSONObject(s, marker string) ([]byte, error) {
	i := strings.Index(s, marker)
	if i < 0 {
		return nil, fmt.Errorf("%s not found", marker)
	}
	i += len(marker)
	for i < len(s) && (s[i] == ' ' || s[i] == '\n' || s[i] == '\r' || s[i] == '\t' || s[i] == '=') {
		i++
	}
	if i >= len(s) || s[i] != '{' {
		return nil, fmt.Errorf("%s is not a JSON object", marker)
	}
	depth, inStr, esc := 0, false, false
	for j := i; j < len(s); j++ {
		c := s[j]
		if inStr {
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return []byte(s[i : j+1]), nil
			}
		}
	}
	return nil, fmt.Errorf("%s JSON object unclosed", marker)
}

func (r *Registry) zaobao(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.zaochenbao.com/realtime/", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(decodeGBK(body)))
	if err != nil {
		return nil, err
	}
	base := "https://www.zaochenbao.com"
	var out []Item
	doc.Find("div.list-block>a.item").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := s.Find(".eps").Text()
		if href == "" || title == "" {
			return
		}
		out = append(out, Item{ID: href, Title: title, URL: base + href})
	})
	return out, nil
}

func (r *Registry) gelonghui(ctx context.Context) ([]Item, error) {
	base := "https://www.gelonghui.com"
	body, err := r.get(ctx, base+"/news/", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find(".article-content").Each(func(_ int, s *goquery.Selection) {
		a := s.Find(".detail-right>a")
		href, _ := a.Attr("href")
		title := strings.TrimSpace(a.Find("h2").Text())
		if href == "" || title == "" {
			return
		}
		out = append(out, Item{ID: href, Title: title, URL: base + href})
	})
	return out, nil
}

func (r *Registry) sputnik(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://sputniknews.cn/services/widget/lenta/", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<div[^>]*class="[^"]*lenta__item[^"]*"[^>]*>[\s\S]*?<a[^>]*href="([^"]*)"[^>]*>[\s\S]*?<div[^>]*class="[^"]*lenta__item-text[^"]*"[^>]*>([^<]+)</div>[\s\S]*?<div[^>]*class="[^"]*lenta__item-date[^"]*"[^>]*data-unixtime="(\d+)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	for _, m := range matches {
		path, title := m[1], strings.TrimSpace(m[2])
		unix, _ := strconv.ParseInt(m[3], 10, 64)
		out = append(out, Item{ID: path, Title: title, URL: "https://sputniknews.cn" + path, PubDate: unix * 1000})
	}
	return out, nil
}

func (r *Registry) cankaoxiaoxi(ctx context.Context) ([]Item, error) {
	channels := []string{"zhongguo", "guandian", "gj"}
	type row struct {
		Data struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			URL         string `json:"url"`
			PublishTime string `json:"publishTime"`
		} `json:"data"`
	}
	var out []Item
	for _, ch := range channels {
		var res struct {
			List []row `json:"list"`
		}
		if err := r.getJSON(ctx, "https://china.cankaoxiaoxi.com/json/channel/"+ch+"/list.json", nil, &res); err != nil {
			return nil, err
		}
		for _, k := range res.List {
			out = append(out, Item{ID: k.Data.ID, Title: k.Data.Title, URL: k.Data.URL, Extra: map[string]any{"date": k.Data.PublishTime}})
		}
	}
	return out, nil
}

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

func (r *Registry) steam(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://store.steampowered.com/stats/stats/", nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	var out []Item
	doc.Find("#detailStats tr.player_count_row").Each(func(_ int, s *goquery.Selection) {
		a := s.Find("a.gameLink")
		href, _ := a.Attr("href")
		name := strings.TrimSpace(a.Text())
		players := strings.TrimSpace(s.Find("td:first-child .currentServers").Text())
		if href == "" || name == "" {
			return
		}
		out = append(out, Item{ID: href, Title: name, URL: href, Extra: map[string]any{"info": players}})
	})
	return out, nil
}

func coolapkDeviceID() string {
	lens := []int{8, 4, 4, 4, 12}
	parts := make([]string, len(lens))
	for i, n := range lens {
		buf := make([]byte, n)
		_, _ = rand.Read(buf)
		const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
		for j := range buf {
			buf[j] = alphabet[int(buf[j])%len(alphabet)]
		}
		parts[i] = string(buf)
	}
	return strings.Join(parts, "-")
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func coolapkToken() string {
	deviceID := coolapkDeviceID()
	now := time.Now().Unix()
	hexNow := fmt.Sprintf("0x%x", now)
	md5Now := md5Hex(strconv.FormatInt(now, 10))
	s := "token://com.coolapk.market/c67ef5943784d09750dcfbb31020f0ab?" + md5Now + "$" + deviceID + "&com.coolapk.market"
	return md5Hex(base64.StdEncoding.EncodeToString([]byte(s))) + deviceID + hexNow
}

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
		"X-App-Token":      coolapkToken(),
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

func (r *Registry) ghxi(ctx context.Context) ([]Item, error) {
	body, err := r.get(ctx, "https://www.ghxi.com/category/all", nil)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`<li[^>]*>[\s\S]*?<a[^>]*href="([^"]*)"[^>]*class="[^"]*item-title[^"]*"[^>]*>([^<]+)</a>[\s\S]*?<div[^>]*class="[^"]*item-excerpt[^"]*"[^>]*>([^<]+)</div>[\s\S]*?<div[^>]*class="[^"]*date[^"]*"[^>]*>([^<]+)</div>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	out := make([]Item, 0, len(matches))
	for _, m := range matches {
		link, title := m[1], strings.TrimSpace(m[2])
		out = append(out, Item{ID: link, Title: title, URL: link, Extra: map[string]any{"hover": strings.TrimSpace(m[3]), "date": strings.TrimSpace(m[4])}})
	}
	return out, nil
}
