package mcpapp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	mcperr "nfxnews/errors/src/mcp"
	"nfxnews/modules/mcp/domain/toolcall"
	"nfxnews/modules/notify/domain/channel"
	"nfxnews/modules/report/application/frequency"
	crawlpb "nfxnews/protos/gen/crawl"
	newspb "nfxnews/protos/gen/news"
	reportpb "nfxnews/protos/gen/report"
	sourcepb "nfxnews/protos/gen/source"

	"github.com/google/uuid"
)

var ToolNames = []string{
	"get_latest_news",
	"get_trending_topics",
	"get_news_by_date",
	"analyze_topic_trend",
	"analyze_data_insights",
	"analyze_sentiment",
	"find_similar_news",
	"generate_summary_report",
	"search_news",
	"search_related_news_history",
	"get_current_config",
	"get_system_status",
	"trigger_crawl",
	"list_sources",
	"list_keywords",
}

type Service struct {
	toolCallRepo *toolcall.Repo
	news         newspb.NewsServiceClient
	report       reportpb.ReportServiceClient
	source       sourcepb.SourceServiceClient
	crawl        crawlpb.CrawlServiceClient
}

func NewService(
	toolCallRepo *toolcall.Repo,
	news newspb.NewsServiceClient,
	report reportpb.ReportServiceClient,
	source sourcepb.SourceServiceClient,
	crawl crawlpb.CrawlServiceClient,
) *Service {
	return &Service{toolCallRepo: toolCallRepo, news: news, report: report, source: source, crawl: crawl}
}

func (s *Service) Tools() []string { return append([]string{}, ToolNames...) }

func (s *Service) RunTool(ctx context.Context, name, argsJSON string) (map[string]any, error) {
	args := map[string]any{}
	if argsJSON != "" {
		_ = json.Unmarshal([]byte(argsJSON), &args)
	}
	raw, _ := json.Marshal(args)
	ok := true
	var errMsg *string
	result, runErr := s.dispatch(ctx, name, args)
	if runErr != nil {
		msg := runErr.Error()
		ok = false
		errMsg = &msg
	}
	_ = s.toolCallRepo.Create.New(ctx, toolcall.NewFromState(toolcall.State{
		ID: uuid.Must(uuid.NewV7()), ToolName: name, Arguments: raw, OK: ok, ErrorMessage: errMsg,
	}))
	if runErr != nil {
		return nil, runErr
	}
	if result == nil {
		result = map[string]any{}
	}
	return result, nil
}

func (s *Service) dispatch(ctx context.Context, name string, args map[string]any) (map[string]any, error) {
	switch name {
	case "get_latest_news":
		return s.getLatestNews(ctx, args)
	case "get_trending_topics":
		return s.getTrendingTopics(ctx, args)
	case "get_news_by_date":
		return s.getNewsByDate(ctx, args)
	case "analyze_topic_trend":
		return s.analyzeTopicTrend(ctx, args)
	case "analyze_data_insights":
		return s.analyzeDataInsights(ctx, args)
	case "analyze_sentiment":
		return s.analyzeSentiment(ctx, args)
	case "find_similar_news":
		return s.findSimilarNews(ctx, args)
	case "generate_summary_report":
		return s.generateSummaryReport(ctx, args)
	case "search_news":
		return s.searchNews(ctx, args)
	case "search_related_news_history":
		return s.searchRelatedHistory(ctx, args)
	case "get_current_config":
		return s.getCurrentConfig(ctx, args)
	case "get_system_status":
		return s.getSystemStatus(ctx)
	case "trigger_crawl":
		return s.triggerCrawl(ctx, args)
	case "list_sources":
		return s.listSources(ctx)
	case "list_keywords":
		return s.listKeywords(ctx)
	default:
		return nil, mcperr.ErrMCPToolUnknown
	}
}

func (s *Service) requireNews() error {
	if s.news == nil {
		return mcperr.ErrMCPClientMissing
	}
	return nil
}

func argInt(args map[string]any, key string, fallback int) int {
	v, ok := args[key]
	if !ok {
		return fallback
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		i, err := strconv.Atoi(n)
		if err == nil {
			return i
		}
	}
	return fallback
}

func argFloat(args map[string]any, key string, fallback float64) float64 {
	v, ok := args[key]
	if !ok {
		return fallback
	}
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	}
	return fallback
}

func argString(args map[string]any, key, fallback string) string {
	if v, ok := args[key].(string); ok && v != "" {
		return v
	}
	return fallback
}

func argStrings(args map[string]any, key string) []string {
	raw, ok := args[key]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case []any:
		out := make([]string, 0, len(v))
		for _, it := range v {
			if s, ok := it.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return v
	case string:
		if v == "" {
			return nil
		}
		return strings.Split(v, ",")
	}
	return nil
}

func argDateRange(args map[string]any) (start, end time.Time, ok bool) {
	raw, exists := args["date_range"]
	if !exists {
		return time.Time{}, time.Time{}, false
	}
	m, okm := raw.(map[string]any)
	if !okm {
		return time.Time{}, time.Time{}, false
	}
	ss, _ := m["start"].(string)
	es, _ := m["end"].(string)
	st, err1 := time.Parse("2006-01-02", ss)
	et, err2 := time.Parse("2006-01-02", es)
	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, false
	}
	return st, et.Add(24*time.Hour - time.Second), true
}

func pubTime(it *newspb.NewsItem) time.Time {
	v := it.GetPubDateUnix()
	if v > 1_000_000_000_000 {
		v = v / 1000
	}
	if v <= 0 {
		return time.Time{}
	}
	return time.Unix(v, 0).UTC()
}

func filterPlatforms(items []*newspb.NewsItem, platforms []string) []*newspb.NewsItem {
	if len(platforms) == 0 {
		return items
	}
	allow := map[string]struct{}{}
	for _, p := range platforms {
		allow[strings.TrimSpace(p)] = struct{}{}
	}
	out := items[:0]
	for _, it := range items {
		if _, ok := allow[it.GetSourceId()]; ok {
			out = append(out, it)
		}
	}
	return out
}

func filterRange(items []*newspb.NewsItem, start, end time.Time) []*newspb.NewsItem {
	if start.IsZero() && end.IsZero() {
		return items
	}
	out := make([]*newspb.NewsItem, 0, len(items))
	for _, it := range items {
		t := pubTime(it)
		if t.IsZero() {
			continue
		}
		if !start.IsZero() && t.Before(start) {
			continue
		}
		if !end.IsZero() && t.After(end) {
			continue
		}
		out = append(out, it)
	}
	return out
}

func parseDateQuery(q string) (start, end time.Time) {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	q = strings.TrimSpace(strings.ToLower(q))
	switch q {
	case "", "today", "今天":
		return now, now.Add(24*time.Hour - time.Second)
	case "yesterday", "昨天":
		return now.Add(-24 * time.Hour), now.Add(-time.Second)
	case "前天":
		return now.Add(-48 * time.Hour), now.Add(-24*time.Hour - time.Second)
	}
	if t, err := time.Parse("2006-01-02", q); err == nil {
		return t, t.Add(24*time.Hour - time.Second)
	}
	return now, now.Add(24*time.Hour - time.Second)
}

func itemMaps(items []*newspb.NewsItem, includeURL bool, limit int) []map[string]any {
	if limit <= 0 || limit > len(items) {
		limit = len(items)
	}
	out := make([]map[string]any, 0, limit)
	for i := 0; i < limit; i++ {
		it := items[i]
		row := map[string]any{"id": it.GetId(), "title": it.GetTitle(), "source_id": it.GetSourceId(), "pub_date_unix": it.GetPubDateUnix()}
		if includeURL {
			row["url"] = it.GetUrl()
		}
		out = append(out, row)
	}
	return out
}

func (s *Service) searchAll(ctx context.Context, query string, limit int) ([]*newspb.NewsItem, error) {
	if err := s.requireNews(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	resp, err := s.news.SearchNews(ctx, &newspb.SearchNewsRequest{Query: query, Limit: int32(limit)})
	if err != nil {
		return nil, mcperr.ErrMCPToolFailed.WithCause(err)
	}
	return resp.GetItems(), nil
}

func (s *Service) getLatestNews(ctx context.Context, args map[string]any) (map[string]any, error) {
	limit := argInt(args, "limit", 50)
	items, err := s.searchAll(ctx, "", limit)
	if err != nil {
		return nil, err
	}
	items = filterPlatforms(items, argStrings(args, "platforms"))
	includeURL, _ := args["include_url"].(bool)
	return map[string]any{"count": len(items), "items": itemMaps(items, includeURL, limit)}, nil
}

func (s *Service) getNewsByDate(ctx context.Context, args map[string]any) (map[string]any, error) {
	limit := argInt(args, "limit", 50)
	items, err := s.searchAll(ctx, "", 500)
	if err != nil {
		return nil, err
	}
	start, end := parseDateQuery(argString(args, "date_query", "今天"))
	if rs, re, ok := argDateRange(args); ok {
		start, end = rs, re
	}
	items = filterRange(filterPlatforms(items, argStrings(args, "platforms")), start, end)
	includeURL, _ := args["include_url"].(bool)
	return map[string]any{"count": len(items), "start": start.Format("2006-01-02"), "end": end.Format("2006-01-02"), "items": itemMaps(items, includeURL, limit)}, nil
}

func (s *Service) searchNews(ctx context.Context, args map[string]any) (map[string]any, error) {
	q := argString(args, "query", "")
	limit := argInt(args, "limit", 50)
	items, err := s.searchAll(ctx, q, 500)
	if err != nil {
		return nil, err
	}
	if rs, re, ok := argDateRange(args); ok {
		items = filterRange(items, rs, re)
	}
	items = filterPlatforms(items, argStrings(args, "platforms"))
	mode := argString(args, "search_mode", "keyword")
	if mode == "fuzzy" || mode == "entity" {
		th := argFloat(args, "threshold", 0.6)
		filtered := items[:0]
		for _, it := range items {
			if similarity(q, it.GetTitle()) >= th {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}
	includeURL, _ := args["include_url"].(bool)
	return map[string]any{"count": len(items), "items": itemMaps(items, includeURL, limit)}, nil
}

func (s *Service) getTrendingTopics(ctx context.Context, args map[string]any) (map[string]any, error) {
	mode := argString(args, "mode", "current")
	topN := argInt(args, "top_n", 10)
	limit := 80
	if mode == "daily" {
		limit = 500
	}
	items, err := s.searchAll(ctx, "", limit)
	if err != nil {
		return nil, err
	}
	if mode == "daily" {
		start := time.Now().UTC().Truncate(24 * time.Hour)
		items = filterRange(items, start, start.Add(24*time.Hour))
	}
	dict, _ := frequency.LoadFile(frequencyPath())
	counts := map[string]int{}
	for _, it := range items {
		if g, ok := frequency.MatchTitle(it.GetTitle(), dict); ok {
			counts[g.GroupKey]++
		}
	}
	type pair struct {
		Key   string `json:"topic"`
		Count int    `json:"count"`
	}
	pairs := make([]pair, 0, len(counts))
	for k, c := range counts {
		pairs = append(pairs, pair{Key: k, Count: c})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Count > pairs[j].Count })
	if topN > 0 && len(pairs) > topN {
		pairs = pairs[:topN]
	}
	return map[string]any{"mode": mode, "topics": pairs}, nil
}

func (s *Service) analyzeTopicTrend(ctx context.Context, args map[string]any) (map[string]any, error) {
	topic := argString(args, "topic", "")
	if topic == "" {
		return nil, mcperr.ErrMCPArgsInvalid
	}
	items, err := s.searchAll(ctx, topic, 500)
	if err != nil {
		return nil, err
	}
	start, end := time.Now().UTC().AddDate(0, 0, -7), time.Now().UTC()
	if rs, re, ok := argDateRange(args); ok {
		start, end = rs, re
	}
	items = filterRange(items, start, end)
	byDay := map[string]int{}
	for _, it := range items {
		t := pubTime(it)
		if t.IsZero() {
			continue
		}
		byDay[t.Format("2006-01-02")]++
	}
	kind := argString(args, "analysis_type", "trend")
	return map[string]any{"topic": topic, "analysis_type": kind, "total": len(items), "by_day": byDay}, nil
}

func (s *Service) analyzeDataInsights(ctx context.Context, args map[string]any) (map[string]any, error) {
	kind := argString(args, "insight_type", "platform_compare")
	topic := argString(args, "topic", "")
	items, err := s.searchAll(ctx, topic, 500)
	if err != nil {
		return nil, err
	}
	if rs, re, ok := argDateRange(args); ok {
		items = filterRange(items, rs, re)
	}
	bySource := map[string]int{}
	for _, it := range items {
		bySource[it.GetSourceId()]++
	}
	result := map[string]any{"insight_type": kind, "topic": topic, "by_source": bySource, "total": len(items)}
	if kind == "keyword_cooccur" {
		minF := argInt(args, "min_frequency", 3)
		topN := argInt(args, "top_n", 20)
		result["pairs"] = cooccur(items, minF, topN)
	}
	return result, nil
}

func (s *Service) analyzeSentiment(ctx context.Context, args map[string]any) (map[string]any, error) {
	topic := argString(args, "topic", "")
	items, err := s.searchAll(ctx, topic, argInt(args, "limit", 50))
	if err != nil {
		return nil, err
	}
	items = filterPlatforms(items, argStrings(args, "platforms"))
	if rs, re, ok := argDateRange(args); ok {
		items = filterRange(items, rs, re)
	}
	pos, neg, neu := 0, 0, 0
	for _, it := range items {
		switch sentiment(it.GetTitle()) {
		case "positive":
			pos++
		case "negative":
			neg++
		default:
			neu++
		}
	}
	includeURL, _ := args["include_url"].(bool)
	return map[string]any{
		"positive": pos, "negative": neg, "neutral": neu, "total": len(items),
		"items": itemMaps(items, includeURL, argInt(args, "limit", 50)),
	}, nil
}

func (s *Service) findSimilarNews(ctx context.Context, args map[string]any) (map[string]any, error) {
	title := argString(args, "reference_title", "")
	if title == "" {
		return nil, mcperr.ErrMCPArgsInvalid
	}
	th := argFloat(args, "threshold", 0.6)
	limit := argInt(args, "limit", 50)
	items, err := s.searchAll(ctx, title, 500)
	if err != nil {
		return nil, err
	}
	type hit struct {
		Title      string  `json:"title"`
		SourceID   string  `json:"source_id"`
		URL        string  `json:"url,omitempty"`
		Similarity float64 `json:"similarity"`
	}
	includeURL, _ := args["include_url"].(bool)
	hits := []hit{}
	for _, it := range items {
		if it.GetTitle() == title {
			continue
		}
		sim := similarity(title, it.GetTitle())
		if sim < th {
			continue
		}
		h := hit{Title: it.GetTitle(), SourceID: it.GetSourceId(), Similarity: sim}
		if includeURL {
			h.URL = it.GetUrl()
		}
		hits = append(hits, h)
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].Similarity > hits[j].Similarity })
	if len(hits) > limit {
		hits = hits[:limit]
	}
	return map[string]any{"count": len(hits), "items": hits}, nil
}

func (s *Service) generateSummaryReport(ctx context.Context, args map[string]any) (map[string]any, error) {
	if s.report == nil {
		return nil, mcperr.ErrMCPClientMissing
	}
	mode := argString(args, "report_type", argString(args, "mode", "daily"))
	if mode == "weekly" {
		mode = "daily"
	}
	resp, err := s.report.GenerateReport(ctx, &reportpb.GenerateReportRequest{Mode: mode})
	if err != nil {
		return nil, mcperr.ErrMCPToolFailed.WithCause(err)
	}
	out := map[string]any{"report_id": resp.GetReportId(), "item_count": resp.GetItemCount(), "mode": mode}
	if got, err := s.report.GetReport(ctx, &reportpb.GetReportRequest{ReportId: resp.GetReportId()}); err == nil {
		out["title"] = got.GetTitle()
		out["payload_json"] = got.GetPayloadJson()
	}
	return out, nil
}

func (s *Service) searchRelatedHistory(ctx context.Context, args map[string]any) (map[string]any, error) {
	ref := argString(args, "reference_text", "")
	if ref == "" {
		return nil, mcperr.ErrMCPArgsInvalid
	}
	preset := argString(args, "time_preset", "yesterday")
	now := time.Now().UTC().Truncate(24 * time.Hour)
	var start, end time.Time
	switch preset {
	case "last_week":
		start, end = now.AddDate(0, 0, -7), now.Add(-time.Second)
	case "last_month":
		start, end = now.AddDate(0, 0, -30), now.Add(-time.Second)
	default:
		start, end = now.Add(-24*time.Hour), now.Add(-time.Second)
	}
	items, err := s.searchAll(ctx, ref, 500)
	if err != nil {
		return nil, err
	}
	items = filterRange(items, start, end)
	th := argFloat(args, "threshold", 0.4)
	limit := argInt(args, "limit", 50)
	includeURL, _ := args["include_url"].(bool)
	type hit struct {
		Title    string  `json:"title"`
		SourceID string  `json:"source_id"`
		URL      string  `json:"url,omitempty"`
		Score    float64 `json:"score"`
	}
	hits := []hit{}
	for _, it := range items {
		score := 0.7*overlap(ref, it.GetTitle()) + 0.3*similarity(ref, it.GetTitle())
		if score < th {
			continue
		}
		h := hit{Title: it.GetTitle(), SourceID: it.GetSourceId(), Score: score}
		if includeURL {
			h.URL = it.GetUrl()
		}
		hits = append(hits, h)
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	if len(hits) > limit {
		hits = hits[:limit]
	}
	return map[string]any{"count": len(hits), "items": hits}, nil
}

func (s *Service) getCurrentConfig(ctx context.Context, args map[string]any) (map[string]any, error) {
	section := argString(args, "section", "all")
	out := map[string]any{"section": section}
	if section == "all" || section == "crawler" {
		if s.source != nil {
			resp, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
			if err != nil {
				return nil, mcperr.ErrMCPToolFailed.WithCause(err)
			}
			out["sources"] = resp.GetSources()
		}
		out["crawl_schedule_seconds"] = os.Getenv("CRAWL_SCHEDULE_SECONDS")
	}
	if section == "all" || section == "keywords" {
		kw, err := s.listKeywords(ctx)
		if err != nil {
			return nil, err
		}
		out["keywords"] = kw["keywords"]
		out["frequency_words_path"] = frequencyPath()
	}
	if section == "all" || section == "push" {
		out["notify_kinds"] = channel.Kinds()
	}
	return out, nil
}

func (s *Service) getSystemStatus(ctx context.Context) (map[string]any, error) {
	status := map[string]any{"ok": true, "time": time.Now().UTC().Format(time.RFC3339)}
	if s.source != nil {
		resp, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
		if err == nil {
			status["source_count"] = len(resp.GetSources())
		}
	}
	if s.news != nil {
		resp, err := s.news.SearchNews(ctx, &newspb.SearchNewsRequest{Limit: 20})
		if err == nil {
			status["latest_news_sample"] = len(resp.GetItems())
		}
	}
	status["tools"] = ToolNames
	return status, nil
}

func (s *Service) triggerCrawl(ctx context.Context, args map[string]any) (map[string]any, error) {
	if s.crawl == nil {
		return nil, mcperr.ErrMCPClientMissing
	}
	platforms := argStrings(args, "platforms")
	if len(platforms) == 0 {
		resp, err := s.crawl.TriggerCrawl(ctx, &crawlpb.TriggerCrawlRequest{})
		if err != nil {
			return nil, mcperr.ErrMCPToolFailed.WithCause(err)
		}
		return map[string]any{"session_id": resp.GetSessionId(), "platforms": []string{"*"}}, nil
	}
	sessions := []string{}
	failed := []string{}
	for _, p := range platforms {
		resp, err := s.crawl.TriggerCrawl(ctx, &crawlpb.TriggerCrawlRequest{SourceId: strings.TrimSpace(p)})
		if err != nil {
			failed = append(failed, p)
			continue
		}
		sessions = append(sessions, resp.GetSessionId())
	}
	return map[string]any{"platforms": platforms, "sessions": sessions, "failed_platforms": failed}, nil
}

func (s *Service) listSources(ctx context.Context) (map[string]any, error) {
	if s.source == nil {
		return nil, mcperr.ErrMCPClientMissing
	}
	resp, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
	if err != nil {
		return nil, mcperr.ErrMCPToolFailed.WithCause(err)
	}
	return map[string]any{"sources": resp.GetSources()}, nil
}

func (s *Service) listKeywords(ctx context.Context) (map[string]any, error) {
	if s.report == nil {
		return nil, mcperr.ErrMCPClientMissing
	}
	resp, err := s.report.ListKeywords(ctx, &reportpb.ListKeywordsRequest{})
	if err != nil {
		return nil, mcperr.ErrMCPToolFailed.WithCause(err)
	}
	return map[string]any{"keywords": resp.GetKeywords()}, nil
}

func frequencyPath() string {
	if p := os.Getenv("FREQUENCY_WORDS_PATH"); p != "" {
		return p
	}
	return filepath.Join("config", "frequency_words.txt")
}

func fold(s string) string {
	return strings.Map(func(r rune) rune { return unicode.ToLower(r) }, s)
}

func tokens(s string) map[string]struct{} {
	s = fold(s)
	out := map[string]struct{}{}
	for _, p := range strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		if p != "" {
			out[p] = struct{}{}
		}
	}
	runes := []rune(s)
	for i := 0; i+1 < len(runes); i++ {
		out[string(runes[i:i+2])] = struct{}{}
	}
	return out
}

func similarity(a, b string) float64 {
	ta, tb := tokens(a), tokens(b)
	if len(ta) == 0 || len(tb) == 0 {
		return 0
	}
	inter := 0
	for k := range ta {
		if _, ok := tb[k]; ok {
			inter++
		}
	}
	union := len(ta) + len(tb) - inter
	if union == 0 {
		return 0
	}
	return float64(inter) / float64(union)
}

func overlap(a, b string) float64 {
	ta, tb := tokens(a), tokens(b)
	if len(ta) == 0 {
		return 0
	}
	hit := 0
	for k := range ta {
		if _, ok := tb[k]; ok {
			hit++
		}
	}
	return float64(hit) / float64(len(ta))
}

func sentiment(title string) string {
	t := fold(title)
	pos := []string{"涨", "利好", "突破", "增长", "胜利", "up", "gain", "surge", "record"}
	neg := []string{"跌", "暴跌", "下滑", "危机", "事故", "down", "crash", "loss", "risk"}
	p, n := 0, 0
	for _, w := range pos {
		if strings.Contains(t, w) {
			p++
		}
	}
	for _, w := range neg {
		if strings.Contains(t, w) {
			n++
		}
	}
	if p > n {
		return "positive"
	}
	if n > p {
		return "negative"
	}
	return "neutral"
}

func cooccur(items []*newspb.NewsItem, minF, topN int) []map[string]any {
	counts := map[string]int{}
	for _, it := range items {
		toks := tokens(it.GetTitle())
		list := make([]string, 0, len(toks))
		for t := range toks {
			if len([]rune(t)) < 2 {
				continue
			}
			list = append(list, t)
		}
		sort.Strings(list)
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				counts[list[i]+"|"+list[j]]++
			}
		}
	}
	type pair struct {
		A, B  string
		Count int
	}
	pairs := []pair{}
	for k, c := range counts {
		if c < minF {
			continue
		}
		ab := strings.SplitN(k, "|", 2)
		if len(ab) != 2 {
			continue
		}
		pairs = append(pairs, pair{A: ab[0], B: ab[1], Count: c})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Count > pairs[j].Count })
	if topN > 0 && len(pairs) > topN {
		pairs = pairs[:topN]
	}
	out := make([]map[string]any, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, map[string]any{"a": p.A, "b": p.B, "count": p.Count})
	}
	return out
}
