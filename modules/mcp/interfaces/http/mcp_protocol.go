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
	addTool(s, svc, "get_latest_news", "Return the latest persisted news items", map[string]string{"limit": "max items"})
	addTool(s, svc, "search_news", "Search persisted news by title", map[string]string{"query": "keyword"})
	addTool(s, svc, "list_sources", "List configured news sources", nil)
	addTool(s, svc, "generate_report", "Generate a keyword report snapshot", map[string]string{"mode": "daily|current|incremental"})
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
