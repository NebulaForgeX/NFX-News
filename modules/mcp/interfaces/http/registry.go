package http

import (
	mcpapp "nfxnews/modules/mcp/application/mcp"
	"nfxnews/modules/mcp/interfaces/http/handler"
)

type Registry struct {
	App  *handler.MCPHandler
	I18n *handler.I18nHandler
}

func NewRegistry(svc *mcpapp.Service, langs string) *Registry {
	return &Registry{App: handler.NewMCPHandler(svc), I18n: handler.NewI18nHandler(langs)}
}
