package http

import (
	"encoding/json"
	"time"

	authconn "nfxnews/connections/auth"
	notifyapp "nfxnews/modules/notify/application/notify"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/fiberx/middleware"
	"nfxnews/pkgs/httpx"
	"nfxnews/pkgs/security/token"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type httpDeps interface {
	AppSvc() *notifyapp.Service
	UserTokenVerifier() token.Verifier
	ErrorsLangsPath() string
	AuthClient() *authconn.Client
}

func NewHTTPServer(d httpDeps, accessLog httpx.AccessLogConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal, ErrorHandler: fiberx.ErrorHandler,
		ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 120 * time.Second,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-Api-Key", "X-Request-ID"},
		AllowCredentials: false, ExposeHeaders: []string{"Content-Length", "Content-Type"}, MaxAge: 3600,
	}))
	app.Use(middleware.Logger(), middleware.AccessLog(accessLog), middleware.Recover())
	NewRouter(app, d.UserTokenVerifier(), NewRegistry(d.AppSvc(), d.ErrorsLangsPath(), d.AuthClient())).RegisterRoutes()
	return app
}
