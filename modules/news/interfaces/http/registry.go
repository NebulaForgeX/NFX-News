package http

import (
	authconn "nfxnews/connections/auth"
	newsapp "nfxnews/modules/news/application/news"
	"nfxnews/modules/news/interfaces/http/handler"
)

type Registry struct {
	App  *handler.NewsHandler
	I18n *handler.I18nHandler
}

func NewRegistry(svc *newsapp.Service, langs string, identity *authconn.Client) *Registry {
	return &Registry{App: handler.NewNewsHandler(svc, identity), I18n: handler.NewI18nHandler(langs)}
}
