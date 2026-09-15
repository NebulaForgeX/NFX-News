package http

import (
	notifyapp "nfxnews/modules/notify/application/notify"
	"nfxnews/modules/notify/interfaces/http/handler"
)

type Registry struct {
	App  *handler.NotifyHandler
	I18n *handler.I18nHandler
}

func NewRegistry(svc *notifyapp.Service, langs string) *Registry {
	return &Registry{App: handler.NewNotifyHandler(svc), I18n: handler.NewI18nHandler(langs)}
}
