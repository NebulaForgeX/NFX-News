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
	mcp := r.app.Group("/mcp")
	r.RegisterLocalesGroup(mcp)
	r.RegisterToolsGroup(mcp)
}

func (r *Router) RegisterLocalesGroup(mcp fiber.Router) {
	locales := mcp.Group("/locales")
	locales.Get("/:lang", r.handlers.I18n.GetErrorTranslations)
	messages := mcp.Group("/messages")
	messages.Get("/:lang", r.handlers.I18n.GetMessageTranslations)
}

func (r *Router) RegisterToolsGroup(mcp fiber.Router) {
	me := mcp.Group("", middleware.TokenAuth(r.tokenVerifier))
	me.Get("/tools", r.handlers.App.Tools)
	me.Post("/tools/:name", r.handlers.App.RunTool)
	me.Post("/run", r.handlers.App.RunJSON)
}
