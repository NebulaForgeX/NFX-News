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
	news := r.app.Group("/news")
	r.RegisterLocalesGroup(news)
	r.RegisterPublicGroup(news)
	r.RegisterMeGroup(news)
}

func (r *Router) RegisterLocalesGroup(news fiber.Router) {
	locales := news.Group("/locales")
	locales.Get("/:lang", r.handlers.I18n.GetErrorTranslations)
	messages := news.Group("/messages")
	messages.Get("/:lang", r.handlers.I18n.GetMessageTranslations)
}

func (r *Router) RegisterPublicGroup(news fiber.Router) {
	news.Get("/items", r.handlers.App.List)
	news.Get("/search", r.handlers.App.Search)
}

func (r *Router) RegisterMeGroup(news fiber.Router) {
	me := news.Group("", middleware.TokenAuth(r.tokenVerifier))
	me.Get("/preferences", r.handlers.App.GetPreferences)
	me.Put("/preferences", r.handlers.App.SetPreferences)
}
