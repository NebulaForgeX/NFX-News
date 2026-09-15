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
	g := r.app.Group("/report")
	g.Get("/i18n/errors/:lang", r.handlers.I18n.GetErrorTranslations)
	g.Get("/keywords", r.handlers.App.ListKeywords)
	g.Get("/snapshots", r.handlers.App.ListSnapshots)
	g.Get("/snapshots/:id", r.handlers.App.Get)
	protected := g.Group("", middleware.TokenAuth(r.tokenVerifier))
	protected.Post("/keywords", r.handlers.App.AddKeyword)
	protected.Post("/snapshots", r.handlers.App.Generate)
}
