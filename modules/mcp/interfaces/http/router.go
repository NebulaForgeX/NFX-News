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
	g := r.app.Group("/mcp")
	g.Get("/i18n/errors/:lang", r.handlers.I18n.GetErrorTranslations)
	g.Get("/tools", r.handlers.App.Tools)
	g.Post("/tools/:name", r.handlers.App.RunTool)
	g.Post("/run", r.handlers.App.RunJSON)
}
