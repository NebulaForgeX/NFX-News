package http

import (
	"context"
	"encoding/json"

	mcpapp "nfxnews/modules/mcp/application/mcp"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func mountMCPProtocol(app *fiber.App, svc *mcpapp.Service) {
	s := server.NewMCPServer("nfx-news", "1.0.0", server.WithToolCapabilities(true))
	addTool(s, svc, "get_latest_news", "Return the latest persisted news items", map[string]string{"limit": "max items", "platforms": "comma-separated source ids"})
	addTool(s, svc, "get_trending_topics", "Count frequency-word matches in latest or daily news", map[string]string{"top_n": "top N groups", "mode": "daily|current"})
	addTool(s, svc, "get_news_by_date", "News for a calendar day", map[string]string{"date_query": "today/yesterday/YYYY-MM-DD"})
	addTool(s, svc, "analyze_topic_trend", "Count a topic by day", map[string]string{"topic": "keyword", "analysis_type": "trend|lifecycle|viral|predict"})
	addTool(s, svc, "analyze_data_insights", "Platform compare or keyword co-occurrence", map[string]string{"insight_type": "platform_compare|platform_activity|keyword_cooccur", "topic": "optional keyword"})
	addTool(s, svc, "analyze_sentiment", "Simple title sentiment split", map[string]string{"topic": "optional keyword"})
	addTool(s, svc, "find_similar_news", "Find titles similar to a reference", map[string]string{"reference_title": "title", "threshold": "0-1"})
	addTool(s, svc, "generate_summary_report", "Generate a keyword report snapshot", map[string]string{"report_type": "daily|current|incremental"})
	addTool(s, svc, "search_news", "Search persisted news by title", map[string]string{"query": "keyword", "search_mode": "keyword|fuzzy|entity"})
	addTool(s, svc, "search_related_news_history", "Related titles in a historical window", map[string]string{"reference_text": "seed title", "time_preset": "yesterday|last_week|last_month"})
	addTool(s, svc, "get_current_config", "Sources, keywords, crawl and notify settings", map[string]string{"section": "all|crawler|push|keywords"})
	addTool(s, svc, "get_system_status", "Health snapshot of news/source counts", nil)
	addTool(s, svc, "trigger_crawl", "Trigger crawl for all sources or selected platforms", map[string]string{"platforms": "comma-separated source ids"})
	addTool(s, svc, "list_sources", "List configured news sources", nil)
	addTool(s, svc, "list_keywords", "List stored keyword DSL rows", nil)
	httpServer := server.NewStreamableHTTPServer(s, server.WithEndpointPath("/mcp/protocol"))
	app.All("/mcp/protocol", adaptor.HTTPHandler(httpServer))
	app.All("/mcp/protocol/*", adaptor.HTTPHandler(httpServer))
}

func addTool(s *server.MCPServer, svc *mcpapp.Service, name, desc string, strArgs map[string]string) {
	opts := []mcp.ToolOption{mcp.WithDescription(desc)}
	for k, d := range strArgs {
		opts = append(opts, mcp.WithString(k, mcp.Description(d)))
	}
	s.AddTool(mcp.NewTool(name, opts...), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, _ := json.Marshal(req.GetArguments())
		result, err := svc.RunTool(ctx, name, string(raw))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		out, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(out)), nil
	})
}
