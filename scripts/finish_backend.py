#!/usr/bin/env python3
"""Generate remaining NFX-News HTTP/gRPC/wiring files and fix {{ }} leftovers."""
from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def write(rel: str, content: str) -> None:
    path = ROOT / rel
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content.lstrip("\n") if content.startswith("\n") else content)
    if not content.endswith("\n"):
        path.write_text(path.read_text() + "\n")


def fix_double_braces() -> None:
    for path in ROOT.glob("modules/**/*.go"):
        text = path.read_text()
        if "{{" not in text:
            continue
        path.write_text(text.replace("{{", "{").replace("}}", "}"))


HTTP_SERVER = '''package http

import (
	"encoding/json"
	"time"

	{appsvc} "{modpath}"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/fiberx/middleware"
	"nfxnews/pkgs/httpx"
	"nfxnews/pkgs/security/token"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type httpDeps interface {{
	AppSvc() *{appsvc}.Service
	UserTokenVerifier() token.Verifier
	ErrorsLangsPath() string
}}

func NewHTTPServer(d httpDeps, accessLog httpx.AccessLogConfig) *fiber.App {{
	app := fiber.New(fiber.Config{{
		JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal, ErrorHandler: fiberx.ErrorHandler,
		ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 120 * time.Second,
	}})
	app.Use(cors.New(cors.Config{{
		AllowOrigins: []string{{"*"}},
		AllowMethods: []string{{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}},
		AllowHeaders: []string{{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-Api-Key", "X-Request-ID"}},
		AllowCredentials: false, ExposeHeaders: []string{{"Content-Length", "Content-Type"}}, MaxAge: 3600,
	}}))
	app.Use(middleware.Logger(), middleware.AccessLog(accessLog), middleware.Recover())
	NewRouter(app, d.UserTokenVerifier(), NewRegistry(d.AppSvc(), d.ErrorsLangsPath())).RegisterRoutes()
	return app
}}
'''

REGISTRY = '''package http

import (
	{appsvc} "{modpath}"
	"nfxnews/modules/{module}/interfaces/http/handler"
)

type Registry struct {{
	App  *handler.{Handler}Handler
	I18n *handler.I18nHandler
}}

func NewRegistry(svc *{appsvc}.Service, langs string) *Registry {{
	return &Registry{{App: handler.New{Handler}Handler(svc), I18n: handler.NewI18nHandler(langs)}}
}}
'''

GRPC_SERVER = '''package grpc

import (
	"nfxnews/modules/{module}/application/resource"
	{appsvc} "{modpath}"
	grpcHandler "nfxnews/modules/{module}/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	{pbalias} "{pbimport}"

	"google.golang.org/grpc"
)

type Deps interface {{
	AppSvc() *{appsvc}.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}}

func NewServer(d Deps) *grpc.Server {{
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	{pbalias}.Register{Svc}Server(s, grpcHandler.New{Handler}Handler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "{module}"))
	return s
}}
'''

WIRING = '''package server

import (
	"context"
	"fmt"
	"time"

	{appsvc} "{modpath}"
	resourceApp "nfxnews/modules/{module}/application/resource"
	"nfxnews/modules/{module}/config"
	"nfxnews/pkgs/cachex"
	"nfxnews/pkgs/health"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/postgresqlx"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	"nfxnews/pkgs/tokenx"
	{extra_imports}
)

type Dependencies struct {{
	healthMgr           *health.Manager
	cache               *cachex.Connection
	postgres            *postgresqlx.Connection
	kafkaConfig         *kafkax.Config
	busPublisher        *eventbus.BusPublisher
	appSvc              *{appsvc}.Service
	resourceSvc         *resourceApp.Service
	userTokenVerifier   token.Verifier
	serverTokenVerifier token.Verifier
	errorsLangsPath     string
	conns               []*grpc.ClientConn
}}

func NewDeps(ctx context.Context, cfg *config.Config) (*Dependencies, error) {{
	postgres, err := postgresqlx.Init(ctx, cfg.PostgreSQL)
	if err != nil {{
		return nil, fmt.Errorf("init PostgreSQL: %w", err)
	}}
	cacheConn, err := cachex.InitConn(ctx, cfg.Cache)
	if err != nil {{
		return nil, fmt.Errorf("init Redis: %w", err)
	}}
	healthMgr := health.NewManager(ctx, 30*time.Second)
	healthMgr.Register(postgres)
	healthMgr.Register(cacheConn)
	kafkaConfig := cfg.KafkaConfig
	busPublisher, err := kafkax.NewPublisher(&kafkaConfig)
	if err != nil {{
		return nil, fmt.Errorf("kafka publisher: %w", err)
	}}
	tokenxInstance := tokenx.New(cfg.Token)
	userTokenVerifier := &tokenxVerifierAdapter{{tokenx: tokenxInstance}}
	serverTokenVerifier := servertoken.NewVerifier(
		&servertoken.HMACSigner{{Key: []byte(cfg.Token.SecretKey)}},
		cfg.Token.Issuer,
		servertoken.WithAllowedSkew(5*time.Second),
	)
	provider := servertoken.NewProvider(
		&servertoken.HMACSigner{{Key: []byte(cfg.Token.SecretKey)}},
		cfg.Token.Issuer,
		"{module}",
	)
	errorsLangsPath := cfg.I18n.ErrorsLangsPath
	if errorsLangsPath == "" {{
		errorsLangsPath = "./errors/langs"
	}}
	d := &Dependencies{{
		healthMgr: healthMgr, postgres: postgres, cache: cacheConn, kafkaConfig: &kafkaConfig,
		busPublisher: busPublisher,
		resourceSvc: resourceApp.NewService(postgres, cacheConn, &kafkaConfig),
		userTokenVerifier: userTokenVerifier, serverTokenVerifier: serverTokenVerifier, errorsLangsPath: errorsLangsPath,
	}}
	{dial_and_svc}
	return d, nil
}}

func (d *Dependencies) Cleanup() {{
	d.healthMgr.Stop()
	d.postgres.Close()
	d.cache.Close()
	for _, c := range d.conns {{
		_ = c.Close()
	}}
}}

func (d *Dependencies) AppSvc() *{appsvc}.Service              {{ return d.appSvc }}
func (d *Dependencies) ResourceSvc() *resourceApp.Service      {{ return d.resourceSvc }}
func (d *Dependencies) UserTokenVerifier() token.Verifier     {{ return d.userTokenVerifier }}
func (d *Dependencies) ServerTokenVerifier() token.Verifier   {{ return d.serverTokenVerifier }}
func (d *Dependencies) KafkaConfig() *kafkax.Config           {{ return d.kafkaConfig }}
func (d *Dependencies) BusPublisher() *eventbus.BusPublisher {{ return d.busPublisher }}
func (d *Dependencies) ErrorsLangsPath() string                {{ return d.errorsLangsPath }}

type tokenxVerifierAdapter struct{{ tokenx *tokenx.Tokenx }}

func (a *tokenxVerifierAdapter) Verify(ctx context.Context, tokenStr string) (*token.Claims, error) {{
	claims, err := a.tokenx.VerifyAccessToken(tokenStr)
	if err != nil {{
		return nil, err
	}}
	return &token.Claims{{Registered: claims.RegisteredClaims, Raw: map[string]any{{"user_id": claims.UserID}}}}, nil
}}
'''


def main() -> None:
    fix_double_braces()

    write(
        "pkgs/grpcx/dial.go",
        """
package grpcx

import (
	"fmt"

	"nfxnews/pkgs/security/token/servertoken"
	"nfxnews/pkgs/tokenx"

	"google.golang.org/grpc"
)

func Dial(addr, serviceID string, tok tokenx.Config) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, fmt.Errorf("empty grpc address for %s", serviceID)
	}
	provider := servertoken.NewProvider(
		&servertoken.HMACSigner{Key: []byte(tok.SecretKey)},
		tok.Issuer,
		serviceID,
	)
	return grpc.NewClient(addr, DefaultClientOptions(provider)...)
}
""",
    )

    write(
        "connections/source/client.go",
        """
package source

import (
	sourcepb "nfxnews/protos/gen/source"

	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) sourcepb.SourceServiceClient {
	return sourcepb.NewSourceServiceClient(conn)
}
""",
    )
    write(
        "connections/news/client.go",
        """
package news

import (
	newspb "nfxnews/protos/gen/news"

	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) newspb.NewsServiceClient {
	return newspb.NewNewsServiceClient(conn)
}
""",
    )
    write(
        "connections/report/client.go",
        """
package report

import (
	reportpb "nfxnews/protos/gen/report"

	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) reportpb.ReportServiceClient {
	return reportpb.NewReportServiceClient(conn)
}
""",
    )

    # --- news HTTP/gRPC ---
    write(
        "modules/news/interfaces/http/handler/news.go",
        """
package handler

import (
	"strconv"

	newsapp "nfxnews/modules/news/application/news"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type NewsHandler struct{ svc *newsapp.Service }

func NewNewsHandler(svc *newsapp.Service) *NewsHandler { return &NewsHandler{svc: svc} }

func (h *NewsHandler) List(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, err := h.svc.ListBySource(c.Context(), c.Query("source_id"), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: items})
}

func (h *NewsHandler) Search(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, err := h.svc.Search(c.Context(), c.Query("q"), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: items})
}
""",
    )
    write(
        "modules/news/interfaces/http/router.go",
        """
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
	g := r.app.Group("/news")
	g.Get("/items", r.handlers.App.List)
	g.Get("/search", r.handlers.App.Search)
	g.Get("/i18n/errors/:lang", r.handlers.I18n.GetErrorTranslations)
}
""",
    )
    write(
        "modules/news/interfaces/http/registry.go",
        REGISTRY.format(appsvc="newsapp", modpath="nfxnews/modules/news/application/news", module="news", Handler="News"),
    )
    write(
        "modules/news/interfaces/http/server.go",
        HTTP_SERVER.format(appsvc="newsapp", modpath="nfxnews/modules/news/application/news"),
    )
    write(
        "modules/news/interfaces/grpc/handler/news.go",
        """
package handler

import (
	"context"
	"encoding/json"

	"nfxnews/events"
	newsapp "nfxnews/modules/news/application/news"
	newspb "nfxnews/protos/gen/news"
)

type NewsHandler struct {
	newspb.UnimplementedNewsServiceServer
	svc *newsapp.Service
}

func NewNewsHandler(svc *newsapp.Service) *NewsHandler { return &NewsHandler{svc: svc} }

func toPB(items []newsapp.ItemView) []*newspb.NewsItem {
	out := make([]*newspb.NewsItem, 0, len(items))
	for _, it := range items {
		extra, _ := json.Marshal(it.Extra)
		out = append(out, &newspb.NewsItem{
			Id: it.ID, SourceId: it.SourceID, OriginalId: it.OriginalID, Title: it.Title,
			Url: it.URL, MobileUrl: it.MobileURL, PubDateUnix: it.PubDate, ExtraJson: string(extra),
		})
	}
	return out
}

func (h *NewsHandler) UpsertItems(ctx context.Context, req *newspb.UpsertItemsRequest) (*newspb.UpsertItemsResponse, error) {
	ev := events.SourceFetchedEvent{SourceID: req.GetSourceId()}
	for _, it := range req.GetItems() {
		ev.Items = append(ev.Items, events.SourceNewsItem{
			ID: it.GetOriginalId(), Title: it.GetTitle(), URL: it.GetUrl(), MobileURL: it.GetMobileUrl(),
			PubDate: it.GetPubDateUnix(), ExtraJSON: it.GetExtraJson(),
		})
	}
	n, err := h.svc.UpsertFromEvent(ctx, ev)
	if err != nil {
		return nil, err
	}
	return &newspb.UpsertItemsResponse{Upserted: int32(n)}, nil
}

func (h *NewsHandler) ListBySource(ctx context.Context, req *newspb.ListBySourceRequest) (*newspb.ListBySourceResponse, error) {
	items, err := h.svc.ListBySource(ctx, req.GetSourceId(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &newspb.ListBySourceResponse{Items: toPB(items)}, nil
}

func (h *NewsHandler) SearchNews(ctx context.Context, req *newspb.SearchNewsRequest) (*newspb.SearchNewsResponse, error) {
	items, err := h.svc.Search(ctx, req.GetQuery(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &newspb.SearchNewsResponse{Items: toPB(items)}, nil
}
""",
    )
    write(
        "modules/news/interfaces/grpc/server.go",
        GRPC_SERVER.format(
            module="news", appsvc="newsapp", modpath="nfxnews/modules/news/application/news",
            pbalias="newspb", pbimport="nfxnews/protos/gen/news", Handler="News", Svc="NewsService",
        ),
    )
    write(
        "modules/news/server/wiring.go",
        WIRING.format(
            module="news", appsvc="newsapp", modpath="nfxnews/modules/news/application/news",
            extra_imports='"google.golang.org/grpc"',
            dial_and_svc="d.appSvc = newsapp.NewService(postgres.DB(), cacheConn)\n\t_ = provider",
        ),
    )
    write(
        "modules/news/interfaces/pipeline/server.go",
        """
package pipeline

import (
	"context"
	"time"

	"nfxnews/events"
	newsapp "nfxnews/modules/news/application/news"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/logx"

	"github.com/ThreeDotsLabs/watermill/message"
	wmMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Deps interface {
	KafkaConfig() *kafkax.Config
	BusPublisher() *eventbus.BusPublisher
	AppSvc() *newsapp.Service
}

type Router struct {
	*eventbus.EventRouter
	svc *newsapp.Service
}

func NewServer(d Deps) (*Router, error) {
	sub, err := kafkax.NewSubscriber(d.KafkaConfig())
	if err != nil {
		return nil, err
	}
	router, err := eventbus.NewEventRouter(sub, eventbus.EventRouterConfig{
		CloseTimeout: 10 * time.Second,
		Logger:       logx.NewZapWatermillLogger(logx.L()),
	})
	if err != nil {
		return nil, err
	}
	r := &Router{EventRouter: router, svc: d.AppSvc()}
	router.AddMiddleware(
		wmMiddleware.CorrelationID,
		wmMiddleware.Recoverer,
		wmMiddleware.Retry{MaxRetries: 3, InitialInterval: 200 * time.Millisecond, MaxInterval: 2 * time.Second, Multiplier: 2.0}.Middleware,
		wmMiddleware.Timeout(10*time.Second),
	)
	return r, nil
}

func (r *Router) RegisterRoutes() {
	eventbus.RegisterHandler(r.EventRouter, func(ctx context.Context, evt events.SourceFetchedEvent, msg *message.Message) error {
		_, err := r.svc.UpsertFromEvent(ctx, evt)
		return err
	})
}

func (r *Router) Run(ctx context.Context) error {
	logx.S().Info("Starting pipeline router...")
	return r.Router.Run(ctx)
}

func (r *Router) Close() error { return r.Router.Close() }
""",
    )

    # --- crawl ---
    write(
        "modules/crawl/interfaces/http/handler/crawl.go",
        """
package handler

import (
	"strconv"

	crawlapp "nfxnews/modules/crawl/application/crawl"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type CrawlHandler struct{ svc *crawlapp.Service }

func NewCrawlHandler(svc *crawlapp.Service) *CrawlHandler { return &CrawlHandler{svc: svc} }

type triggerBody struct {
	SourceID string `json:"source_id"`
}

func (h *CrawlHandler) Trigger(c fiber.Ctx) error {
	var body triggerBody
	_ = c.Bind().Body(&body)
	sess, err := h.svc.Trigger(c.Context(), body.SourceID)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "triggered", httpx.SuccessOptions{Data: sess})
}

func (h *CrawlHandler) List(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.svc.List(c.Context(), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *CrawlHandler) Get(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errx.ErrInvalidParams.WithMsg("id required")
	}
	row, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: row})
}
""",
    )
    write(
        "modules/crawl/interfaces/http/router.go",
        """
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
	g := r.app.Group("/crawl")
	g.Get("/i18n/errors/:lang", r.handlers.I18n.GetErrorTranslations)
	g.Get("/sessions", r.handlers.App.List)
	g.Get("/sessions/:id", r.handlers.App.Get)
	protected := g.Group("", middleware.TokenAuth(r.tokenVerifier))
	protected.Post("/sessions", r.handlers.App.Trigger)
}
""",
    )
    write(
        "modules/crawl/interfaces/http/registry.go",
        REGISTRY.format(appsvc="crawlapp", modpath="nfxnews/modules/crawl/application/crawl", module="crawl", Handler="Crawl"),
    )
    write(
        "modules/crawl/interfaces/http/server.go",
        HTTP_SERVER.format(appsvc="crawlapp", modpath="nfxnews/modules/crawl/application/crawl"),
    )
    write(
        "modules/crawl/interfaces/grpc/handler/crawl.go",
        """
package handler

import (
	"context"

	crawlapp "nfxnews/modules/crawl/application/crawl"
	crawlpb "nfxnews/protos/gen/crawl"
)

type CrawlHandler struct {
	crawlpb.UnimplementedCrawlServiceServer
	svc *crawlapp.Service
}

func NewCrawlHandler(svc *crawlapp.Service) *CrawlHandler { return &CrawlHandler{svc: svc} }

func (h *CrawlHandler) TriggerCrawl(ctx context.Context, req *crawlpb.TriggerCrawlRequest) (*crawlpb.TriggerCrawlResponse, error) {
	sess, err := h.svc.Trigger(ctx, req.GetSourceId())
	if err != nil {
		return nil, err
	}
	return &crawlpb.TriggerCrawlResponse{SessionId: sess.ID.String()}, nil
}

func (h *CrawlHandler) GetSession(ctx context.Context, req *crawlpb.GetSessionRequest) (*crawlpb.GetSessionResponse, error) {
	sess, err := h.svc.Get(ctx, req.GetSessionId())
	if err != nil {
		return nil, err
	}
	sid := ""
	if sess.SourceID != nil {
		sid = *sess.SourceID
	}
	errMsg := ""
	if sess.ErrorMessage != nil {
		errMsg = *sess.ErrorMessage
	}
	return &crawlpb.GetSessionResponse{Session: &crawlpb.CrawlSession{
		Id: sess.ID.String(), SourceId: sid, Status: sess.Status, ItemCount: int32(sess.ItemCount), ErrorMessage: errMsg,
	}}, nil
}
""",
    )
    write(
        "modules/crawl/interfaces/grpc/server.go",
        GRPC_SERVER.format(
            module="crawl", appsvc="crawlapp", modpath="nfxnews/modules/crawl/application/crawl",
            pbalias="crawlpb", pbimport="nfxnews/protos/gen/crawl", Handler="Crawl", Svc="CrawlService",
        ),
    )
    write(
        "modules/crawl/server/wiring.go",
        WIRING.format(
            module="crawl", appsvc="crawlapp", modpath="nfxnews/modules/crawl/application/crawl",
            extra_imports='sourceconn "nfxnews/connections/source"\n\tsourcepb "nfxnews/protos/gen/source"\n\t"nfxnews/pkgs/grpcx"\n\n\t"google.golang.org/grpc"',
            dial_and_svc="""var sourceClient sourcepb.SourceServiceClient
	if conn, err := grpcx.Dial(cfg.GRPCClient.SourceAddr, "crawl", cfg.Token); err == nil {
		d.conns = append(d.conns, conn)
		sourceClient = sourceconn.New(conn)
	}
	d.appSvc = crawlapp.NewService(postgres.DB(), sourceClient)
	_ = provider""",
        ),
    )
    write(
        "modules/crawl/interfaces/pipeline/server.go",
        """
package pipeline

import (
	"context"
	"os"
	"strconv"
	"time"

	crawlapp "nfxnews/modules/crawl/application/crawl"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/logx"

	wmMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Deps interface {
	KafkaConfig() *kafkax.Config
	BusPublisher() *eventbus.BusPublisher
	AppSvc() *crawlapp.Service
}

type Router struct {
	*eventbus.EventRouter
	svc *crawlapp.Service
}

func NewServer(d Deps) (*Router, error) {
	sub, err := kafkax.NewSubscriber(d.KafkaConfig())
	if err != nil {
		return nil, err
	}
	router, err := eventbus.NewEventRouter(sub, eventbus.EventRouterConfig{
		CloseTimeout: 10 * time.Second,
		Logger:       logx.NewZapWatermillLogger(logx.L()),
	})
	if err != nil {
		return nil, err
	}
	r := &Router{EventRouter: router, svc: d.AppSvc()}
	router.AddMiddleware(
		wmMiddleware.CorrelationID,
		wmMiddleware.Recoverer,
		wmMiddleware.Retry{MaxRetries: 3, InitialInterval: 200 * time.Millisecond, MaxInterval: 2 * time.Second, Multiplier: 2.0}.Middleware,
		wmMiddleware.Timeout(10*time.Second),
	)
	return r, nil
}

func (r *Router) RegisterRoutes() {}

func (r *Router) Run(ctx context.Context) error {
	go r.schedule(ctx)
	logx.S().Info("Starting pipeline router...")
	return r.Router.Run(ctx)
}

func (r *Router) schedule(ctx context.Context) {
	sec := 600
	if v := os.Getenv("CRAWL_SCHEDULE_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			sec = n
		}
	}
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	defer ticker.Stop()
	_, _ = r.svc.TriggerAll(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := r.svc.TriggerAll(ctx); err != nil {
				logx.S().Warnf("scheduled crawl failed: %v", err)
			}
		}
	}
}

func (r *Router) Close() error { return r.Router.Close() }
""",
    )

    # --- report ---
    write(
        "modules/report/interfaces/http/handler/report.go",
        """
package handler

import (
	"strconv"

	reportapp "nfxnews/modules/report/application/report"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type ReportHandler struct{ svc *reportapp.Service }

func NewReportHandler(svc *reportapp.Service) *ReportHandler { return &ReportHandler{svc: svc} }

type keywordBody struct {
	GroupName  string `json:"group_name"`
	Word       string `json:"word"`
	Kind       string `json:"kind"`
	CountLimit int    `json:"count_limit"`
}

type generateBody struct {
	Mode string `json:"mode"`
}

func (h *ReportHandler) ListKeywords(c fiber.Ctx) error {
	rows, err := h.svc.ListKeywords(c.Context())
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *ReportHandler) AddKeyword(c fiber.Ctx) error {
	var body keywordBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	row, err := h.svc.AddKeyword(c.Context(), body.GroupName, body.Word, body.Kind, body.CountLimit)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "created", httpx.SuccessOptions{Data: row})
}

func (h *ReportHandler) Generate(c fiber.Ctx) error {
	var body generateBody
	_ = c.Bind().Body(&body)
	snap, err := h.svc.Generate(c.Context(), body.Mode)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "generated", httpx.SuccessOptions{Data: snap})
}

func (h *ReportHandler) ListSnapshots(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.svc.ListSnapshots(c.Context(), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *ReportHandler) Get(c fiber.Ctx) error {
	row, err := h.svc.Get(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: row})
}
""",
    )
    write(
        "modules/report/interfaces/http/router.go",
        """
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
""",
    )
    write(
        "modules/report/interfaces/http/registry.go",
        REGISTRY.format(appsvc="reportapp", modpath="nfxnews/modules/report/application/report", module="report", Handler="Report"),
    )
    write(
        "modules/report/interfaces/http/server.go",
        HTTP_SERVER.format(appsvc="reportapp", modpath="nfxnews/modules/report/application/report"),
    )
    write(
        "modules/report/interfaces/grpc/handler/report.go",
        """
package handler

import (
	"context"

	reportapp "nfxnews/modules/report/application/report"
	reportpb "nfxnews/protos/gen/report"
)

type ReportHandler struct {
	reportpb.UnimplementedReportServiceServer
	svc *reportapp.Service
}

func NewReportHandler(svc *reportapp.Service) *ReportHandler { return &ReportHandler{svc: svc} }

func (h *ReportHandler) GenerateReport(ctx context.Context, req *reportpb.GenerateReportRequest) (*reportpb.GenerateReportResponse, error) {
	snap, err := h.svc.Generate(ctx, req.GetMode())
	if err != nil {
		return nil, err
	}
	return &reportpb.GenerateReportResponse{ReportId: snap.ID.String(), ItemCount: int32(snap.ItemCount)}, nil
}

func (h *ReportHandler) GetReport(ctx context.Context, req *reportpb.GetReportRequest) (*reportpb.GetReportResponse, error) {
	snap, err := h.svc.Get(ctx, req.GetReportId())
	if err != nil {
		return nil, err
	}
	return &reportpb.GetReportResponse{
		Id: snap.ID.String(), Mode: snap.Mode, Title: snap.Title, PayloadJson: string(snap.Payload), ItemCount: int32(snap.ItemCount),
	}, nil
}

func (h *ReportHandler) ListKeywords(ctx context.Context, req *reportpb.ListKeywordsRequest) (*reportpb.ListKeywordsResponse, error) {
	rows, err := h.svc.ListKeywords(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*reportpb.Keyword, 0, len(rows))
	for _, k := range rows {
		out = append(out, &reportpb.Keyword{Id: k.ID.String(), GroupName: k.GroupName, Word: k.Word, Kind: k.Kind, CountLimit: int32(k.CountLimit)})
	}
	return &reportpb.ListKeywordsResponse{Keywords: out}, nil
}
""",
    )
    write(
        "modules/report/interfaces/grpc/server.go",
        GRPC_SERVER.format(
            module="report", appsvc="reportapp", modpath="nfxnews/modules/report/application/report",
            pbalias="reportpb", pbimport="nfxnews/protos/gen/report", Handler="Report", Svc="ReportService",
        ),
    )
    write(
        "modules/report/server/wiring.go",
        WIRING.format(
            module="report", appsvc="reportapp", modpath="nfxnews/modules/report/application/report",
            extra_imports='newsconn "nfxnews/connections/news"\n\tnewspb "nfxnews/protos/gen/news"\n\t"nfxnews/pkgs/grpcx"\n\n\t"google.golang.org/grpc"',
            dial_and_svc="""var newsClient newspb.NewsServiceClient
	if conn, err := grpcx.Dial(cfg.GRPCClient.NewsAddr, "report", cfg.Token); err == nil {
		d.conns = append(d.conns, conn)
		newsClient = newsconn.New(conn)
	}
	d.appSvc = reportapp.NewService(postgres.DB(), newsClient, busPublisher)
	_ = provider""",
        ),
    )

    # --- notify ---
    write(
        "modules/notify/interfaces/http/handler/notify.go",
        """
package handler

import (
	"nfxnews/events"
	notifyapp "nfxnews/modules/notify/application/notify"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type NotifyHandler struct{ svc *notifyapp.Service }

func NewNotifyHandler(svc *notifyapp.Service) *NotifyHandler { return &NotifyHandler{svc: svc} }

type channelBody struct {
	Kind    string         `json:"kind"`
	Name    string         `json:"name"`
	Enabled bool           `json:"enabled"`
	Config  map[string]any `json:"config"`
}

type dispatchBody struct {
	ReportID  string `json:"report_id"`
	Mode      string `json:"mode"`
	Title     string `json:"title"`
	ItemCount int    `json:"item_count"`
}

func (h *NotifyHandler) ListChannels(c fiber.Ctx) error {
	rows, err := h.svc.ListChannels(c.Context())
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *NotifyHandler) UpsertChannel(c fiber.Ctx) error {
	var body channelBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	row, err := h.svc.UpsertChannel(c.Context(), body.Kind, body.Name, body.Enabled, body.Config)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "created", httpx.SuccessOptions{Data: row})
}

func (h *NotifyHandler) Dispatch(c fiber.Ctx) error {
	var body dispatchBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	n, err := h.svc.DispatchReport(c.Context(), events.ReportGeneratedEvent{
		ReportID: body.ReportID, Mode: body.Mode, Title: body.Title, ItemCount: body.ItemCount,
	})
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: map[string]any{"queued": n}})
}
""",
    )
    write(
        "modules/notify/interfaces/http/router.go",
        """
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
	g := r.app.Group("/notify")
	g.Get("/i18n/errors/:lang", r.handlers.I18n.GetErrorTranslations)
	protected := g.Group("", middleware.TokenAuth(r.tokenVerifier))
	protected.Get("/channels", r.handlers.App.ListChannels)
	protected.Post("/channels", r.handlers.App.UpsertChannel)
	protected.Post("/dispatch", r.handlers.App.Dispatch)
}
""",
    )
    write(
        "modules/notify/interfaces/http/registry.go",
        REGISTRY.format(appsvc="notifyapp", modpath="nfxnews/modules/notify/application/notify", module="notify", Handler="Notify"),
    )
    write(
        "modules/notify/interfaces/http/server.go",
        HTTP_SERVER.format(appsvc="notifyapp", modpath="nfxnews/modules/notify/application/notify"),
    )
    write(
        "modules/notify/interfaces/grpc/handler/notify.go",
        """
package handler

import (
	"context"

	"nfxnews/events"
	notifyapp "nfxnews/modules/notify/application/notify"
	notifypb "nfxnews/protos/gen/notify"
)

type NotifyHandler struct {
	notifypb.UnimplementedNotifyServiceServer
	svc *notifyapp.Service
}

func NewNotifyHandler(svc *notifyapp.Service) *NotifyHandler { return &NotifyHandler{svc: svc} }

func (h *NotifyHandler) DispatchReport(ctx context.Context, req *notifypb.DispatchReportRequest) (*notifypb.DispatchReportResponse, error) {
	n, err := h.svc.DispatchReport(ctx, events.ReportGeneratedEvent{ReportID: req.GetReportId(), Title: "report", Mode: "manual"})
	if err != nil {
		return nil, err
	}
	return &notifypb.DispatchReportResponse{Queued: int32(n)}, nil
}
""",
    )
    write(
        "modules/notify/interfaces/grpc/server.go",
        GRPC_SERVER.format(
            module="notify", appsvc="notifyapp", modpath="nfxnews/modules/notify/application/notify",
            pbalias="notifypb", pbimport="nfxnews/protos/gen/notify", Handler="Notify", Svc="NotifyService",
        ),
    )
    write(
        "modules/notify/server/wiring.go",
        WIRING.format(
            module="notify", appsvc="notifyapp", modpath="nfxnews/modules/notify/application/notify",
            extra_imports='"google.golang.org/grpc"',
            dial_and_svc="d.appSvc = notifyapp.NewService(postgres.DB())\n\t_ = provider",
        ),
    )
    write(
        "modules/notify/interfaces/pipeline/server.go",
        """
package pipeline

import (
	"context"
	"time"

	"nfxnews/events"
	notifyapp "nfxnews/modules/notify/application/notify"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/logx"

	"github.com/ThreeDotsLabs/watermill/message"
	wmMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Deps interface {
	KafkaConfig() *kafkax.Config
	BusPublisher() *eventbus.BusPublisher
	AppSvc() *notifyapp.Service
}

type Router struct {
	*eventbus.EventRouter
	svc *notifyapp.Service
}

func NewServer(d Deps) (*Router, error) {
	sub, err := kafkax.NewSubscriber(d.KafkaConfig())
	if err != nil {
		return nil, err
	}
	router, err := eventbus.NewEventRouter(sub, eventbus.EventRouterConfig{
		CloseTimeout: 10 * time.Second,
		Logger:       logx.NewZapWatermillLogger(logx.L()),
	})
	if err != nil {
		return nil, err
	}
	r := &Router{EventRouter: router, svc: d.AppSvc()}
	router.AddMiddleware(
		wmMiddleware.CorrelationID,
		wmMiddleware.Recoverer,
		wmMiddleware.Retry{MaxRetries: 3, InitialInterval: 200 * time.Millisecond, MaxInterval: 2 * time.Second, Multiplier: 2.0}.Middleware,
		wmMiddleware.Timeout(15*time.Second),
	)
	return r, nil
}

func (r *Router) RegisterRoutes() {
	eventbus.RegisterHandler(r.EventRouter, func(ctx context.Context, evt events.ReportGeneratedEvent, msg *message.Message) error {
		_, err := r.svc.DispatchReport(ctx, evt)
		return err
	})
}

func (r *Router) Run(ctx context.Context) error {
	logx.S().Info("Starting pipeline router...")
	return r.Router.Run(ctx)
}

func (r *Router) Close() error { return r.Router.Close() }
""",
    )

    # --- mcp ---
    write(
        "modules/mcp/interfaces/http/handler/mcp.go",
        """
package handler

import (
	"encoding/json"

	mcpapp "nfxnews/modules/mcp/application/mcp"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type MCPHandler struct{ svc *mcpapp.Service }

func NewMCPHandler(svc *mcpapp.Service) *MCPHandler { return &MCPHandler{svc: svc} }

type runBody struct {
	Arguments json.RawMessage `json:"arguments"`
}

func (h *MCPHandler) RunTool(c fiber.Ctx) error {
	name := c.Params("name")
	var body runBody
	_ = c.Bind().Body(&body)
	result, err := h.svc.RunTool(c.Context(), name, string(body.Arguments))
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: result})
}

func (h *MCPHandler) Tools(c fiber.Ctx) error {
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: []string{
		"get_latest_news", "search_news", "list_sources", "generate_report", "list_keywords",
	}})
}

func (h *MCPHandler) RunJSON(c fiber.Ctx) error {
	var body map[string]any
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	name, _ := body["tool_name"].(string)
	args, _ := json.Marshal(body["arguments"])
	result, err := h.svc.RunTool(c.Context(), name, string(args))
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: result})
}
""",
    )
    write(
        "modules/mcp/interfaces/http/router.go",
        """
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
""",
    )
    write(
        "modules/mcp/interfaces/http/registry.go",
        REGISTRY.format(appsvc="mcpapp", modpath="nfxnews/modules/mcp/application/mcp", module="mcp", Handler="MCP"),
    )
    write(
        "modules/mcp/interfaces/http/server.go",
        HTTP_SERVER.format(appsvc="mcpapp", modpath="nfxnews/modules/mcp/application/mcp"),
    )
    write(
        "modules/mcp/interfaces/grpc/handler/mcp.go",
        """
package handler

import (
	"context"
	"encoding/json"

	mcpapp "nfxnews/modules/mcp/application/mcp"
	mcppb "nfxnews/protos/gen/mcp"
)

type MCPHandler struct {
	mcppb.UnimplementedMCPServiceServer
	svc *mcpapp.Service
}

func NewMCPHandler(svc *mcpapp.Service) *MCPHandler { return &MCPHandler{svc: svc} }

func (h *MCPHandler) RunTool(ctx context.Context, req *mcppb.RunToolRequest) (*mcppb.RunToolResponse, error) {
	result, err := h.svc.RunTool(ctx, req.GetToolName(), req.GetArgumentsJson())
	if err != nil {
		return &mcppb.RunToolResponse{Ok: false, ErrorMessage: err.Error()}, nil
	}
	raw, _ := json.Marshal(result)
	return &mcppb.RunToolResponse{Ok: true, ResultJson: string(raw)}, nil
}
""",
    )
    write(
        "modules/mcp/interfaces/grpc/server.go",
        GRPC_SERVER.format(
            module="mcp", appsvc="mcpapp", modpath="nfxnews/modules/mcp/application/mcp",
            pbalias="mcppb", pbimport="nfxnews/protos/gen/mcp", Handler="MCP", Svc="MCPService",
        ),
    )
    write(
        "modules/mcp/server/wiring.go",
        WIRING.format(
            module="mcp", appsvc="mcpapp", modpath="nfxnews/modules/mcp/application/mcp",
            extra_imports='newsconn "nfxnews/connections/news"\n\treportconn "nfxnews/connections/report"\n\tsourceconn "nfxnews/connections/source"\n\tnewspb "nfxnews/protos/gen/news"\n\treportpb "nfxnews/protos/gen/report"\n\tsourcepb "nfxnews/protos/gen/source"\n\t"nfxnews/pkgs/grpcx"\n\n\t"google.golang.org/grpc"',
            dial_and_svc="""var newsClient newspb.NewsServiceClient
	var reportClient reportpb.ReportServiceClient
	var sourceClient sourcepb.SourceServiceClient
	if conn, err := grpcx.Dial(cfg.GRPCClient.NewsAddr, "mcp", cfg.Token); err == nil {
		d.conns = append(d.conns, conn)
		newsClient = newsconn.New(conn)
	}
	if conn, err := grpcx.Dial(cfg.GRPCClient.ReportAddr, "mcp", cfg.Token); err == nil {
		d.conns = append(d.conns, conn)
		reportClient = reportconn.New(conn)
	}
	if conn, err := grpcx.Dial(cfg.GRPCClient.SourceAddr, "mcp", cfg.Token); err == nil {
		d.conns = append(d.conns, conn)
		sourceClient = sourceconn.New(conn)
	}
	d.appSvc = mcpapp.NewService(postgres.DB(), newsClient, reportClient, sourceClient)
	_ = provider""",
        ),
    )

    # --- system ---
    write(
        "modules/system/interfaces/http/handler/system.go",
        """
package handler

import (
	systemapp "nfxnews/modules/system/application/system"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type SystemHandler struct{ svc *systemapp.Service }

func NewSystemHandler(svc *systemapp.Service) *SystemHandler { return &SystemHandler{svc: svc} }

type initBody struct {
	Version string `json:"version"`
}

func (h *SystemHandler) Latest(c fiber.Ctx) error {
	row, err := h.svc.Latest(c.Context())
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: row})
}

func (h *SystemHandler) Initialize(c fiber.Ctx) error {
	var body initBody
	_ = c.Bind().Body(&body)
	row, err := h.svc.Initialize(c.Context(), body.Version)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "initialized", httpx.SuccessOptions{Data: row})
}
""",
    )
    write(
        "modules/system/interfaces/http/router.go",
        """
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
	g := r.app.Group("/system")
	g.Get("/i18n/errors/:lang", r.handlers.I18n.GetErrorTranslations)
	g.Get("/system-state/latest", r.handlers.App.Latest)
	g.Post("/system-state/initialize", r.handlers.App.Initialize)
}
""",
    )
    write(
        "modules/system/interfaces/http/registry.go",
        REGISTRY.format(appsvc="systemapp", modpath="nfxnews/modules/system/application/system", module="system", Handler="System"),
    )
    write(
        "modules/system/interfaces/http/server.go",
        HTTP_SERVER.format(appsvc="systemapp", modpath="nfxnews/modules/system/application/system"),
    )
    write(
        "modules/system/interfaces/grpc/handler/system.go",
        """
package handler

import (
	"context"

	systemapp "nfxnews/modules/system/application/system"
	systemstatepb "nfxnews/protos/gen/system/system_state"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type SystemHandler struct {
	systemstatepb.UnimplementedSystemStateServiceServer
	svc *systemapp.Service
}

func NewSystemHandler(svc *systemapp.Service) *SystemHandler { return &SystemHandler{svc: svc} }

func toPB(row *systemapp.State) *systemstatepb.SystemState {
	if row == nil {
		return &systemstatepb.SystemState{}
	}
	st := &systemstatepb.SystemState{Id: row.ID.String(), Initialized: row.Initialized, ResetCount: int32(row.ResetCount)}
	if row.InitializedAt != nil {
		st.InitializedAt = timestamppb.New(*row.InitializedAt)
	}
	st.InitializationVersion = row.InitializationVersion
	return st
}

func (h *SystemHandler) GetLatestSystemState(ctx context.Context, req *systemstatepb.GetLatestSystemStateRequest) (*systemstatepb.GetLatestSystemStateResponse, error) {
	row, err := h.svc.Latest(ctx)
	if err != nil {
		return nil, err
	}
	return &systemstatepb.GetLatestSystemStateResponse{SystemState: toPB(row)}, nil
}

func (h *SystemHandler) InitializeSystem(ctx context.Context, req *systemstatepb.InitializeSystemRequest) (*systemstatepb.InitializeSystemResponse, error) {
	row, err := h.svc.Initialize(ctx, req.GetVersion())
	if err != nil {
		return nil, err
	}
	return &systemstatepb.InitializeSystemResponse{SystemState: toPB(row)}, nil
}
""",
    )
    write(
        "modules/system/interfaces/grpc/server.go",
        """
package grpc

import (
	"nfxnews/modules/system/application/resource"
	systemapp "nfxnews/modules/system/application/system"
	grpcHandler "nfxnews/modules/system/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	systemstatepb "nfxnews/protos/gen/system/system_state"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *systemapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	systemstatepb.RegisterSystemStateServiceServer(s, grpcHandler.NewSystemHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "system"))
	return s
}
""",
    )
    write(
        "modules/system/server/wiring.go",
        WIRING.format(
            module="system", appsvc="systemapp", modpath="nfxnews/modules/system/application/system",
            extra_imports='"google.golang.org/grpc"',
            dial_and_svc="d.appSvc = systemapp.NewService(postgres.DB())\n\t_ = provider",
        ),
    )

    print("generated remaining backend files")


if __name__ == "__main__":
    main()
