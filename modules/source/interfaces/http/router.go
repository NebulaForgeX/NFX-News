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
	g := r.app.Group("/source")
	g.Get("/sources", r.handlers.Source.List)
	g.Get("/sources/:id", r.handlers.Source.Get)
	g.Post("/sources/:id/fetch", r.handlers.Source.Fetch)
	g.Get("/locales/:lang", r.handlers.I18n.GetErrorTranslations)
	g.Get("/messages/:lang", r.handlers.I18n.GetMessageTranslations)
}
