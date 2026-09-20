package http

import (
	"nfxnews/pkgs/fiberx/middleware"
	"nfxnews/pkgs/security/token"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	app           fiber.Router
	tokenVerifier token.Verifier
	handlers      *Registry
}

func NewRouter(app fiber.Router, v token.Verifier, h *Registry) *Router {
	return &Router{app: app, tokenVerifier: v, handlers: h}
}

func (r *Router) RegisterRoutes() {
	notify := r.app.Group("/notify")
	r.RegisterLocalesGroup(notify)
	r.RegisterMeGroup(notify)
}

func (r *Router) RegisterLocalesGroup(notify fiber.Router) {
	locales := notify.Group("/locales")
	locales.Get("/:lang", r.handlers.I18n.GetErrorTranslations)
	messages := notify.Group("/messages")
	messages.Get("/:lang", r.handlers.I18n.GetMessageTranslations)
}

func (r *Router) RegisterMeGroup(notify fiber.Router) {
	me := notify.Group("", middleware.TokenAuth(r.tokenVerifier))
	me.Get("/kinds", r.handlers.App.ListKinds)
	me.Get("/channels", r.handlers.App.ListChannels)
	me.Get("/deliveries", r.handlers.App.ListDeliveries)
	me.Post("/channels", r.handlers.App.UpsertChannel)
	me.Post("/dispatch", r.handlers.App.Dispatch)
}
