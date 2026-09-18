package http

import (
	authconn "nfxnews/connections/auth"
	crawlapp "nfxnews/modules/crawl/application/crawl"
	"nfxnews/modules/crawl/interfaces/http/handler"
)

type Registry struct {
	App  *handler.CrawlHandler
	I18n *handler.I18nHandler
}

func NewRegistry(svc *crawlapp.Service, langs string, identity *authconn.Client) *Registry {
	return &Registry{App: handler.NewCrawlHandler(svc, identity), I18n: handler.NewI18nHandler(langs)}
}
