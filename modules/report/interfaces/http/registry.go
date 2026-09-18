package http

import (
	authconn "nfxnews/connections/auth"
	reportapp "nfxnews/modules/report/application/report"
	"nfxnews/modules/report/interfaces/http/handler"
)

type Registry struct {
	App  *handler.ReportHandler
	I18n *handler.I18nHandler
}

func NewRegistry(svc *reportapp.Service, langs string, identity *authconn.Client) *Registry {
	return &Registry{App: handler.NewReportHandler(svc, identity), I18n: handler.NewI18nHandler(langs)}
}
