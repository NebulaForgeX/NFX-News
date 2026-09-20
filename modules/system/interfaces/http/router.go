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
	g := r.app.Group("/system")
	g.Get("/locales/:lang", r.handlers.I18n.GetErrorTranslations)
	g.Get("/messages/:lang", r.handlers.I18n.GetMessageTranslations)
	protected := g.Group("", middleware.TokenAuth(r.tokenVerifier))
	protected.Get("/system-state/latest", r.handlers.App.Latest)
	protected.Post("/system-state/initialize", r.handlers.App.Initialize)
}
