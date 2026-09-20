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
	crawl := r.app.Group("/crawl")
	r.RegisterLocalesGroup(crawl)
	r.RegisterSessionsGroup(crawl)
}

func (r *Router) RegisterLocalesGroup(crawl fiber.Router) {
	locales := crawl.Group("/locales")
	locales.Get("/:lang", r.handlers.I18n.GetErrorTranslations)
	messages := crawl.Group("/messages")
	messages.Get("/:lang", r.handlers.I18n.GetMessageTranslations)
}

func (r *Router) RegisterSessionsGroup(crawl fiber.Router) {
	me := crawl.Group("", middleware.TokenAuth(r.tokenVerifier))
	me.Get("/sessions", r.handlers.App.List)
	me.Get("/sessions/:id", r.handlers.App.Get)
	me.Post("/sessions", r.handlers.App.Trigger)
}
