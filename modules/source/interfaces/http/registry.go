package http

import (
	sourceapp "nfxnews/modules/source/application/source"
	"nfxnews/modules/source/interfaces/http/handler"
)

type Registry struct {
	Source *handler.SourceHandler
	I18n   *handler.I18nHandler
}

func NewRegistry(svc *sourceapp.Service, langs string) *Registry {
	return &Registry{Source: handler.NewSourceHandler(svc), I18n: handler.NewI18nHandler(langs)}
}
