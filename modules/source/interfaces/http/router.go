package http

import (
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
	source := r.app.Group("/source")
	r.RegisterLocalesGroup(source)
	r.RegisterSourcesGroup(source)
}

func (r *Router) RegisterLocalesGroup(source fiber.Router) {
	locales := source.Group("/locales")
	locales.Get("/:lang", r.handlers.I18n.GetErrorTranslations)
	messages := source.Group("/messages")
	messages.Get("/:lang", r.handlers.I18n.GetMessageTranslations)
}

func (r *Router) RegisterSourcesGroup(source fiber.Router) {
	source.Get("/sources", r.handlers.Source.List)
	source.Get("/sources/:id", r.handlers.Source.Get)
	source.Post("/sources/:id/fetch", r.handlers.Source.Fetch)
}
