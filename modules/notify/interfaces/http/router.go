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
	g := r.app.Group("/notify")
	g.Get("/locales/:lang", r.handlers.I18n.GetErrorTranslations)
	g.Get("/messages/:lang", r.handlers.I18n.GetMessageTranslations)
	protected := g.Group("", middleware.TokenAuth(r.tokenVerifier))
	protected.Get("/kinds", r.handlers.App.ListKinds)
	protected.Get("/channels", r.handlers.App.ListChannels)
	protected.Get("/deliveries", r.handlers.App.ListDeliveries)
	protected.Post("/channels", r.handlers.App.UpsertChannel)
	protected.Post("/dispatch", r.handlers.App.Dispatch)
}
