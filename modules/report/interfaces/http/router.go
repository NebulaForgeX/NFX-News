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
	report := r.app.Group("/report")
	r.RegisterLocalesGroup(report)
	r.RegisterKeywordsGroup(report)
	r.RegisterSnapshotsGroup(report)
}

func (r *Router) RegisterLocalesGroup(report fiber.Router) {
	locales := report.Group("/locales")
	locales.Get("/:lang", r.handlers.I18n.GetErrorTranslations)
	messages := report.Group("/messages")
	messages.Get("/:lang", r.handlers.I18n.GetMessageTranslations)
}

func (r *Router) RegisterKeywordsGroup(report fiber.Router) {
	me := report.Group("", middleware.TokenAuth(r.tokenVerifier))
	me.Get("/keywords", r.handlers.App.ListKeywords)
	me.Post("/keywords", r.handlers.App.AddKeyword)
}

func (r *Router) RegisterSnapshotsGroup(report fiber.Router) {
	me := report.Group("", middleware.TokenAuth(r.tokenVerifier))
	me.Get("/snapshots", r.handlers.App.ListSnapshots)
	me.Get("/snapshots/:id/html", r.handlers.App.HTML)
	me.Get("/snapshots/:id", r.handlers.App.Get)
	me.Post("/snapshots", r.handlers.App.Generate)
}
