#!/usr/bin/env python3
"""Generate Identity-shaped Go module boilerplate for NFX-News."""
from pathlib import Path

ROOT = Path("/volume1/Projects/NebulaForgeX/NFX-News")
MODULES = [
    ("source", "Source", "GRPC_PORT_SOURCE"),
    ("news", "News", "GRPC_PORT_NEWS"),
    ("crawl", "Crawl", "GRPC_PORT_CRAWL"),
    ("report", "Report", "GRPC_PORT_REPORT"),
    ("notify", "Notify", "GRPC_PORT_NOTIFY"),
    ("mcp", "MCP", "GRPC_PORT_MCP"),
    ("system", "System", "GRPC_PORT_SYSTEM"),
]
MODES = ["api", "connection", "pipeline", "messaging", "base"]


def write(rel: str, content: str) -> None:
    path = ROOT / rel
    path.parent.mkdir(parents=True, exist_ok=True)
    if not content.endswith("\n"):
        content += "\n"
    path.write_text(content)


def toml_for(name: str, title: str, grpc_var: str, env_name: str) -> str:
    db = "${POSTGRESQL_NAME_DEV}" if env_name == "dev" else "${POSTGRESQL_NAME_PROD}"
    log_level = "debug" if env_name == "dev" else "info"
    log_fmt = "console" if env_name == "dev" else "json"
    return f'''# NFX-News {title} Service ({env_name})

[server]
    name = "{title} Service"
    host = "0.0.0.0"
    http_port = "${{HTTP_PORT}}"
    grpc_port = "${{{grpc_var}}}"

    [server.access_log]
    mode = "logger"

[postgresql]
    host = "${{POSTGRESQL_HOST}}"
    port = "${{POSTGRESQL_PORT}}"
    user = "${{POSTGRESQL_USER}}"
    password = "${{POSTGRESQL_PASSWORD}}"
    dbname = "{db}"
    sslmode = "disable"
    timezone = "UTC"
    logger_level = "warn"
    auto_migrate = false

    [postgresql.connection]
        timeout = "5s"
        max_retries = 5
        retry_interval = "2s"
        max_idle_connections = 10
        max_open_connections = 100
        conn_max_idle_time = "15m"
        conn_max_lifetime = "1h"

[cache]
    host = "${{REDIS_HOST}}"
    port = "${{REDIS_PORT}}"
    password = "${{REDIS_PASSWORD}}"

    [cache.connection]
        dial_timeout = "3s"
        write_timeout = "3s"
        read_timeout = "3s"
        max_retries = 5
        retry_interval = "2s"

    [cache.tls]
        enabled = false
        server_name = ""

[logger]
    level = "{log_level}"
    format = "{log_fmt}"
    output = "stdout"
    file_path = "logs/{name}.log"
    max_size_mb = 100
    max_backups = 10
    max_age_day = 30
    compress = false

[kafka]
    brokers = ["${{KAFKA_BROKERS}}"]
    client_id = "nfxnews-{name}-service"

    [kafka.network]
        max_open_requests = 1

    [kafka.producer]
        acks = "all"
        compression = "snappy"
        retries = 3
        batch_bytes = 1048576
        linger_ms = 5
        idempotent = true

    [kafka.consumer]
        group_id = "nfxnews-{name}"
        initial_offset = "latest"
        session_timeout_ms = 10000
        heartbeat_interval_ms = 3000
        fetch_min_bytes = 1
        fetch_max_bytes = 5242880
        return_errors = true

    [kafka.producer_topics]
        {name} = "nfxnews.{name}"
        {name}_poison = "nfxnews.{name}_poison"

    [kafka.consumer_topics]
        {name} = "nfxnews.{name}"
        {name}_poison = "nfxnews.{name}_poison"

    [kafka.security]
        enabled = false
        mechanism = "PLAIN"
        username = ""
        password = ""
        tls_insecure_skip_verify = false

[token]
    secret_key = "${{TOKEN_SECRET_KEY}}"
    issuer = "${{TOKEN_ISSUER}}"
    access_token_ttl = "${{TOKEN_ACCESS_TTL}}"
    refresh_token_ttl = "${{TOKEN_REFRESH_TTL}}"
    algorithm = "${{TOKEN_ALGORITHM}}"

[i18n]
    errors_langs_path = "errors/langs"

[grpc_client]
    auth_addr = "${{GRPC_HOST_AUTH}}:${{GRPC_PORT_AUTH}}"
    source_addr = "${{GRPC_HOST_SOURCE}}:${{GRPC_PORT_SOURCE}}"
    news_addr = "${{GRPC_HOST_NEWS}}:${{GRPC_PORT_NEWS}}"
    crawl_addr = "${{GRPC_HOST_CRAWL}}:${{GRPC_PORT_CRAWL}}"
    report_addr = "${{GRPC_HOST_REPORT}}:${{GRPC_PORT_REPORT}}"
    notify_addr = "${{GRPC_HOST_NOTIFY}}:${{GRPC_PORT_NOTIFY}}"
    mcp_addr = "${{GRPC_HOST_MCP}}:${{GRPC_PORT_MCP}}"
    system_addr = "${{GRPC_HOST_SYSTEM}}:${{GRPC_PORT_SYSTEM}}"

[otel]
    enabled = true
    endpoint = "${{OTEL_EXPORTER_OTLP_ENDPOINT}}"
    insecure = true
    sampler_arg = 1.0
    traces = true
    metrics = true
    logs = true
    export_timeout = "10s"
'''


def dockerfile(name: str, mode: str) -> str:
    binary = f"{name}-{mode}-service"
    expose_grpc = "EXPOSE 10170" if mode in ("connection", "base") else ""
    return f'''FROM golang:1.26-alpine AS base-builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download

FROM base-builder AS prod-build
COPY . ./
ARG APP_VERSION=dev
ENV CGO_ENABLED=0 GOOS=linux
RUN --mount=type=cache,target=/go/pkg/mod \\
  --mount=type=cache,target=/root/.cache/go-build \\
  go build -trimpath -ldflags="-s -w -X main.version=${{APP_VERSION}}" \\
  -o /out/{binary} ./inputs/{name}/{mode}

FROM alpine:latest AS prod
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
COPY --from=prod-build /out/{binary} /app/{binary}
COPY inputs/{name}/configuration/prod.toml /app/config.toml
COPY .env /app/.env
ENV APP_ENV=prod
ENV CONFIG_FILE=/app/config.toml
EXPOSE 8080
{expose_grpc}
ENTRYPOINT ["/app/{binary}"]

FROM base-builder AS dev
WORKDIR /app
ENV PATH="/go/bin:$PATH"
RUN go install github.com/air-verse/air@v1.61.7
COPY . ./
CMD ["air", "-c", "inputs/{name}/{mode}/.air.toml"]
'''


def air_toml(name: str, mode: str) -> str:
    binary = f"{name}-{mode}-service"
    return f'''root = "."
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "tmp/{binary}"
  cmd = "go build -trimpath -buildvcs=false -ldflags='-s -w' -o ./tmp/{binary} ./inputs/{name}/{mode}"
  delay = 500
  exclude_file = []
  exclude_unchanged = true
  follow_symlink = true
  include_dir = ["inputs/{name}/{mode}", "inputs/{name}/configuration", "modules/{name}", "enums", "errors/src", "events", "pkgs", "protos"]
  include_ext = ["go","mod","sum","toml"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = true
  poll_interval = 1000
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = true
  silent = false
  time = true

[misc]
  clean_on_exit = false

[proxy]
  app_port = 0
  enabled = false
  proxy_port = 0

[screen]
  clear_on_rebuild = false
  keep_scroll = true
'''


def main_go(name: str, mode: str) -> str:
    run = {
        "api": "RunHTTP",
        "connection": "RunGRPC",
        "pipeline": "RunPipeline",
        "messaging": "RunMessaging",
        "base": "RunServer",
    }[mode]
    default_env = "dev" if mode == "base" else "prod"
    return f'''package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"nfxnews/modules/{name}/config"
	"nfxnews/modules/{name}/server"
	"nfxnews/pkgs/connections/otelx"
	"nfxnews/pkgs/env"
	"nfxnews/pkgs/logx"

	"go.uber.org/zap"
)

func main() {{
	envStr := flag.String("env", "{default_env}", "Environment (dev/prod)")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(ctx, env.Env(*envStr))
	if err != nil {{
		log.Fatalf("load config failed: %v", err)
	}}

	if err := logx.Init(cfg.Logger, "{name}-{mode}-service", env.Env(*envStr)); err != nil {{
		log.Fatalf("logger init failed: %v", err)
	}}
	defer logx.Sync()

	otelShutdown, err := otelx.Init(ctx, cfg.OTEL, "{name}", env.Env(*envStr))
	if err != nil {{
		log.Fatalf("otel init failed: %v", err)
	}}
	defer func() {{ _ = otelShutdown(context.Background()) }}()

	if err := server.{run}(ctx, cfg); err != nil && !errors.Is(err, context.Canceled) {{
		logx.L().Fatal("{mode} server stopped with error", zap.Error(err))
	}}

	logx.L().Info("{mode} server shutdown gracefully")
}}
'''


CONFIG_GO = '''package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"nfxnews/pkgs/configx"
	"nfxnews/pkgs/env"
)

const ServiceName = "{name}"

func Load(ctx context.Context, env env.Env) (*Config, error) {{
	wd, err := os.Getwd()
	if err != nil {{
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}}
	configPath := filepath.Join(wd, "inputs", ServiceName, "configuration", fmt.Sprintf("%s.toml", env))

	loader, err := configx.NewLoader[Config](ctx, configx.WithPath(configPath))
	if err != nil {{
		return nil, err
	}}
	cfg := loader.Config()
	cfg.Env = env

	if err := cfg.KafkaConfig.Validate(); err != nil {{
		return nil, fmt.Errorf("invalid kafka configuration: %w", err)
	}}

	return cfg, nil
}}
'''

TYPES_GO = '''package config

import (
	"nfxnews/pkgs/cachex"
	"nfxnews/pkgs/connections/otelx"
	"nfxnews/pkgs/env"
	"nfxnews/pkgs/httpx"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/logx"
	"nfxnews/pkgs/postgresqlx"
	"nfxnews/pkgs/tokenx"
)

type Config struct {{
	Env            env.Env
	Server         ServerConfig       `koanf:"server"`
	PostgreSQL     postgresqlx.Config `koanf:"postgresql"`
	Cache          cachex.ConnConfig  `koanf:"cache"`
	Logger         logx.LoggerConfig  `koanf:"logger"`
	KafkaConfig    kafkax.Config      `koanf:"kafka"`
	GRPCClient     GRPCClientConfig   `koanf:"grpc_client"`
	Token          tokenx.Config      `koanf:"token"`
	I18n           I18nConfig         `koanf:"i18n"`
	OTEL           otelx.Config       `koanf:"otel"`
}}

type I18nConfig struct {{
	ErrorsLangsPath string `koanf:"errors_langs_path"`
}}

type GRPCClientConfig struct {{
	AuthAddr   string `koanf:"auth_addr"`
	SourceAddr string `koanf:"source_addr"`
	NewsAddr   string `koanf:"news_addr"`
	CrawlAddr  string `koanf:"crawl_addr"`
	ReportAddr string `koanf:"report_addr"`
	NotifyAddr string `koanf:"notify_addr"`
	MCPAddr    string `koanf:"mcp_addr"`
	SystemAddr string `koanf:"system_addr"`
}}

type ServerConfig struct {{
	Name      string                `koanf:"name"`
	Host      string                `koanf:"host"`
	HTTPPort  int                   `koanf:"http_port"`
	GRPCPort  int                   `koanf:"grpc_port"`
	AccessLog httpx.AccessLogConfig `koanf:"access_log"`
}}
'''

SERVER_GO = '''package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"

	"nfxnews/modules/{name}/config"
	grpcInterfaces "nfxnews/modules/{name}/interfaces/grpc"
	httpInterfaces "nfxnews/modules/{name}/interfaces/http"
	messagingInterfaces "nfxnews/modules/{name}/interfaces/messaging"
	eventbusInterfaces "nfxnews/modules/{name}/interfaces/pipeline"
	"nfxnews/pkgs/logx"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func RunHTTP(ctx context.Context, cfg *config.Config) error {{
	deps, err := NewDeps(ctx, cfg)
	if err != nil {{
		return err
	}}
	defer deps.Cleanup()

	httpSrv := httpInterfaces.NewHTTPServer(deps, cfg.Server.AccessLog)
	httpAddr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.HTTPPort))
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {{
		logx.S().Infof("HTTP server listening on %s", httpAddr)
		if err := httpSrv.Listen(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {{
			return err
		}}
		return nil
	}})

	g.Go(func() error {{
		<-gctx.Done()
		_ = httpSrv.Shutdown()
		return gctx.Err()
	}})

	return g.Wait()
}}

func RunGRPC(ctx context.Context, cfg *config.Config) error {{
	deps, err := NewDeps(ctx, cfg)
	if err != nil {{
		return err
	}}
	defer deps.Cleanup()

	grpcSrv := grpcInterfaces.NewServer(deps)
	grpcAddr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.GRPCPort))

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {{
		return err
	}}
	defer lis.Close()

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {{
		logx.S().Infof("gRPC server listening on %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {{
			return err
		}}
		return nil
	}})

	g.Go(func() error {{
		<-gctx.Done()
		grpcSrv.GracefulStop()
		return gctx.Err()
	}})

	return g.Wait()
}}

func RunPipeline(ctx context.Context, cfg *config.Config) error {{
	deps, err := NewDeps(ctx, cfg)
	if err != nil {{
		return err
	}}
	defer deps.Cleanup()

	eventbusSrv, err := eventbusInterfaces.NewServer(deps)
	if err != nil {{
		return err
	}}
	eventbusSrv.RegisterRoutes()

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {{
		return eventbusSrv.Run(ctx)
	}})
	g.Go(func() error {{
		<-gctx.Done()
		_ = eventbusSrv.Close()
		return gctx.Err()
	}})
	return g.Wait()
}}

func RunMessaging(ctx context.Context, cfg *config.Config) error {{
	deps, err := NewDeps(ctx, cfg)
	if err != nil {{
		return err
	}}
	defer deps.Cleanup()

	messagingSrv, err := messagingInterfaces.NewServer(deps)
	if err != nil {{
		return err
	}}
	messagingSrv.RegisterRoutes()

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {{
		return messagingSrv.Run(ctx)
	}})
	g.Go(func() error {{
		<-gctx.Done()
		_ = messagingSrv.Close()
		return gctx.Err()
	}})
	return g.Wait()
}}

func RunServer(ctx context.Context, cfg *config.Config) error {{
	deps, err := NewDeps(ctx, cfg)
	if err != nil {{
		return err
	}}
	defer deps.Cleanup()

	httpSrv := httpInterfaces.NewHTTPServer(deps, cfg.Server.AccessLog)
	grpcSrv := grpcInterfaces.NewServer(deps)
	eventbusSrv, err := eventbusInterfaces.NewServer(deps)
	if err != nil {{
		return err
	}}
	messagingSrv, err := messagingInterfaces.NewServer(deps)
	if err != nil {{
		return err
	}}

	httpAddr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.HTTPPort))
	grpcAddr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.GRPCPort))
	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {{
		return err
	}}
	defer grpcLis.Close()

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {{
		logx.S().Infof("HTTP server listening on %s", httpAddr)
		if err := httpSrv.Listen(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {{
			return err
		}}
		return nil
	}})
	g.Go(func() error {{
		logx.S().Infof("gRPC server listening on %s", grpcAddr)
		if err := grpcSrv.Serve(grpcLis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {{
			return err
		}}
		return nil
	}})
	g.Go(func() error {{
		return eventbusSrv.Run(ctx)
	}})
	g.Go(func() error {{
		return messagingSrv.Run(ctx)
	}})
	g.Go(func() error {{
		<-gctx.Done()
		_ = httpSrv.Shutdown()
		grpcSrv.GracefulStop()
		_ = eventbusSrv.Close()
		_ = messagingSrv.Close()
		return gctx.Err()
	}})
	return g.Wait()
}}
'''

RESOURCE_SERVICE = '''package resource

import (
	"context"

	"nfxnews/pkgs/cachex"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/postgresqlx"
)

type Service struct {{
	postgres *postgresqlx.Connection
	cache    *cachex.Connection
	kafkaCfg *kafkax.Config
}}

func NewService(
	postgres *postgresqlx.Connection,
	cache *cachex.Connection,
	kafkaCfg *kafkax.Config,
) *Service {{
	return &Service{{
		postgres: postgres,
		cache:    cache,
		kafkaCfg: kafkaCfg,
	}}
}}

func (s *Service) CheckPostgres(ctx context.Context) error {{
	if s.postgres == nil {{
		return nil
	}}
	return s.postgres.Check(ctx)
}}

func (s *Service) CheckRedis(ctx context.Context) error {{
	if s.cache == nil {{
		return nil
	}}
	return s.cache.Check(ctx)
}}

func (s *Service) CheckKafka(ctx context.Context) error {{
	return nil
}}
'''

PIPELINE = '''package pipeline

import (
	"context"

	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/logx"
	"time"

	wmMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Deps interface {{
	KafkaConfig() *kafkax.Config
	BusPublisher() *eventbus.BusPublisher
}}

type Router struct {{
	*eventbus.EventRouter
}}

func NewServer(d Deps) (*Router, error) {{
	sub, err := kafkax.NewSubscriber(d.KafkaConfig())
	if err != nil {{
		return nil, err
	}}
	router, err := eventbus.NewEventRouter(sub, eventbus.EventRouterConfig{{
		CloseTimeout: 10 * time.Second,
		Logger:       logx.NewZapWatermillLogger(logx.L()),
	}})
	if err != nil {{
		return nil, err
	}}
	r := &Router{{EventRouter: router}}
	router.AddMiddleware(
		wmMiddleware.CorrelationID,
		wmMiddleware.Recoverer,
		wmMiddleware.Retry{{MaxRetries: 3, InitialInterval: 200 * time.Millisecond, MaxInterval: 2 * time.Second, Multiplier: 2.0}}.Middleware,
		wmMiddleware.Timeout(10*time.Second),
	)
	return r, nil
}}

func (r *Router) RegisterRoutes() {{}}

func (r *Router) Run(ctx context.Context) error {{
	logx.S().Info("Starting pipeline router...")
	return r.Router.Run(ctx)
}}

func (r *Router) Close() error {{
	return r.Router.Close()
}}
'''

MESSAGING = '''package messaging

import (
	"context"

	"nfxnews/pkgs/logx"
)

type Deps interface{{}}

type Router struct{{}}

func NewServer(_ Deps) (*Router, error) {{
	return &Router{{}}, nil
}}

func (r *Router) RegisterRoutes() {{}}

func (r *Router) Run(ctx context.Context) error {{
	logx.S().Info("Kafka messaging router idle (no RabbitMQ); pipeline owns event consumption")
	<-ctx.Done()
	return ctx.Err()
}}

func (r *Router) Close() error {{ return nil }}
'''

I18N_HANDLER = '''package handler

import (
	"os"
	"path/filepath"

	"nfxnews/pkgs/errx"

	"github.com/gofiber/fiber/v3"
)

var supportedLangs = map[string]bool{{"en": true, "zh": true, "fr": true}}

type I18nHandler struct {{
	errorsLangsPath string
}}

func NewI18nHandler(errorsLangsPath string) *I18nHandler {{
	return &I18nHandler{{errorsLangsPath: errorsLangsPath}}
}}

func (h *I18nHandler) GetErrorTranslations(c fiber.Ctx) error {{
	lang := c.Params("lang")
	if lang == "" || !supportedLangs[lang] {{
		return errx.ErrInvalidParams.WithMsg("lang must be one of: en, zh, fr")
	}}
	name := lang + ".json"
	fpath := filepath.Join(h.errorsLangsPath, name)
	data, err := os.ReadFile(fpath)
	if err != nil {{
		if os.IsNotExist(err) {{
			return c.Status(200).JSON(map[string]any{{}})
		}}
		return errx.ErrInternal.WithCause(err)
	}}
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Send(data)
}}
'''

HEALTH_GRPC = '''package handler

import (
	"context"
	"time"

	"nfxnews/modules/{name}/application/resource"
	healthpb "nfxnews/protos/gen/common/health"
)

type HealthHandler struct {{
	healthpb.UnimplementedHealthServiceServer
	resourceSvc *resource.Service
	serviceName string
}}

func NewHealthHandler(resourceSvc *resource.Service, serviceName string) *HealthHandler {{
	return &HealthHandler{{resourceSvc: resourceSvc, serviceName: serviceName}}
}}

func (h *HealthHandler) GetHealth(ctx context.Context, req *healthpb.GetHealthRequest) (*healthpb.GetHealthResponse, error) {{
	infra := &healthpb.InfrastructureHealth{{Others: map[string]*healthpb.ResourceHealth{{}}}}
	allHealthy := true
	now := time.Now().Unix()

	postgresErr := h.resourceSvc.CheckPostgres(ctx)
	dbHealth := &healthpb.ResourceHealth{{Healthy: postgresErr == nil, CheckedAt: now}}
	if postgresErr != nil {{
		errMsg := postgresErr.Error()
		dbHealth.ErrorMessage = &errMsg
		allHealthy = false
	}}
	infra.Database = dbHealth

	redisErr := h.resourceSvc.CheckRedis(ctx)
	redisHealth := &healthpb.ResourceHealth{{Healthy: redisErr == nil, CheckedAt: now}}
	if redisErr != nil {{
		errMsg := redisErr.Error()
		redisHealth.ErrorMessage = &errMsg
		allHealthy = false
	}}
	infra.Redis = redisHealth

	return &healthpb.GetHealthResponse{{
		Healthy: allHealthy, Infrastructure: infra, ServiceName: h.serviceName, CheckedAt: now,
	}}, nil
}}
'''


def compose_dev() -> str:
    services = [
        '''  reverse-proxy:
    image: traefik:latest
    container_name: NFX-News-Reverse-Proxy-Dev
    command:
      - --api.insecure=true
      - --providers.docker
      - --providers.docker.exposedByDefault=false
      - --providers.docker.constraints=Label(`traefik.project`, `nfx-news`)
      - --entrypoints.web.address=:${TRAEFIK_HTTP_INTERNAL_PORT}
      - --entrypoints.traefik.address=:${TRAEFIK_DASHBOARD_INTERNAL_PORT}
    ports:
      - "${TRAEFIK_HTTP_PORT}:${TRAEFIK_HTTP_INTERNAL_PORT}"
      - "${TRAEFIK_HTTPS_PORT}:${TRAEFIK_HTTPS_INTERNAL_PORT}"
      - "${TRAEFIK_DASHBOARD_PORT}:${TRAEFIK_DASHBOARD_INTERNAL_PORT}"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    extra_hosts:
      - "host.docker.internal:host-gateway"
    networks:
      - nfx-news
    restart: unless-stopped
'''
    ]
    for name, _, grpc_var in MODULES:
        prefix = f"API_PREFIX_PATH_{name.upper()}"
        services.append(f'''  {name}-base:
    image: nfx-news-{name}-base:dev
    container_name: NFX-News-{name.capitalize()}-Base-Dev
    build:
      context: .
      dockerfile: inputs/{name}/base/Dockerfile
      target: dev
    environment:
      - ENV=dev
    command: [ "air", "-c", "inputs/{name}/base/.air.toml" ]
    expose:
      - "${{HTTP_PORT}}"
      - "${{{grpc_var}}}"
    ports:
      - "${{{grpc_var}}}:${{{grpc_var}}}"
    labels:
      traefik.enable: "true"
      traefik.project: "nfx-news"
      traefik.http.routers.{name}-base.entrypoints: web
      traefik.http.routers.{name}-base.rule: PathPrefix(`${{{prefix}}}`)
      traefik.http.services.{name}-base.loadbalancer.server.port: "${{HTTP_PORT}}"
    volumes:
      - ./static:/app/static
      - ./data:/app/data
      - ./assets:/app/assets
      - ./errors/langs:/app/errors/langs
    extra_hosts:
      - "host.docker.internal:host-gateway"
    networks:
      - nfx-news
    depends_on:
      reverse-proxy:
        condition: service_started
''')
    services.append('''  console:
    image: nfx-news-console:dev
    container_name: NFX-News-Console-Dev
    build:
      context: ./console
      dockerfile: Dockerfile
      additional_contexts:
        nfx-ui: ../NFX-UI
      target: prod
    ports:
      - "${CONSOLE_EXTERNAL_PORT}:${CONSOLE_PORT}"
    labels:
      traefik.enable: "true"
      traefik.project: "nfx-news"
      traefik.http.routers.console.entrypoints: web
      traefik.http.routers.console.rule: PathPrefix(`/`)
      traefik.http.routers.console.priority: "1"
      traefik.http.services.console.loadbalancer.server.port: "${CONSOLE_PORT}"
    networks:
      - nfx-news
''')
    return "name: nfx-news-dev\\n\\nservices:\\n" + "\\n".join(services) + '''
networks:
  nfx-news:
    name: nfx-news
    driver: bridge
'''


def compose_prod() -> str:
    chunks = [
        '''  reverse-proxy:
    image: traefik:latest
    container_name: NFX-News-Reverse-Proxy
    command:
      - --api.insecure=true
      - --providers.docker
      - --providers.docker.exposedByDefault=false
      - --providers.docker.constraints=Label(`traefik.project`, `nfx-news`)
      - --entrypoints.web.address=:${TRAEFIK_HTTP_INTERNAL_PORT}
      - --entrypoints.websecure.address=:${TRAEFIK_HTTPS_INTERNAL_PORT}
      - --entrypoints.traefik.address=:${TRAEFIK_DASHBOARD_INTERNAL_PORT}
    ports:
      - "${TRAEFIK_HTTP_PORT}:${TRAEFIK_HTTP_INTERNAL_PORT}"
      - "${TRAEFIK_HTTPS_PORT}:${TRAEFIK_HTTPS_INTERNAL_PORT}"
      - "${TRAEFIK_DASHBOARD_PORT}:${TRAEFIK_DASHBOARD_INTERNAL_PORT}"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ${TRAEFIK_DYNAMIC_DIR}:/etc/traefik/dynamic:ro
      - ${CERTS_DIR}:/certs:ro
    extra_hosts:
      - "host.docker.internal:host-gateway"
    networks:
      - nfx-news
    restart: unless-stopped
'''
    ]
    for name, _, grpc_var in MODULES:
        prefix = f"API_PREFIX_PATH_{name.upper()}"
        for mode in ("api", "connection", "pipeline", "messaging"):
            ports = ""
            expose = '      - "${HTTP_PORT}"'
            labels = ""
            if mode == "api":
                labels = f'''    labels:
      traefik.enable: "true"
      traefik.project: "nfx-news"
      traefik.http.routers.{name}-api.entrypoints: web
      traefik.http.routers.{name}-api.rule: PathPrefix(`${{{prefix}}}`)
      traefik.http.services.{name}-api.loadbalancer.server.port: "${{HTTP_PORT}}"
'''
            if mode == "connection":
                expose += f'\\n      - "${{{grpc_var}}}"'
                ports = f'''    ports:
      - "${{{grpc_var}}}:${{{grpc_var}}}"
'''
            chunks.append(f'''  {name}-{mode}:
    image: nfx-news-{name}-{mode}:prod
    container_name: NFX-News-{name.capitalize()}-{mode.capitalize()}
    build:
      context: .
      dockerfile: inputs/{name}/{mode}/Dockerfile
      target: prod
    environment:
      - ENV=prod
    expose:
{expose}
{ports}{labels}    extra_hosts:
      - "host.docker.internal:host-gateway"
    networks:
      - nfx-news
    restart: unless-stopped
    depends_on:
      reverse-proxy:
        condition: service_started
''')
    chunks.append('''  console:
    image: nfx-news-console:prod
    container_name: NFX-News-Console
    build:
      context: ./console
      dockerfile: Dockerfile
      additional_contexts:
        nfx-ui: ../NFX-UI
    ports:
      - "${CONSOLE_EXTERNAL_PORT}:${CONSOLE_PORT}"
    labels:
      traefik.enable: "true"
      traefik.project: "nfx-news"
      traefik.http.routers.console.entrypoints: web
      traefik.http.routers.console.rule: PathPrefix(`/`)
      traefik.http.routers.console.priority: "1"
      traefik.http.services.console.loadbalancer.server.port: "${CONSOLE_PORT}"
    networks:
      - nfx-news
    restart: unless-stopped
''')
    return "name: nfx-news\\n\\nservices:\\n" + "\\n".join(chunks) + '''
networks:
  nfx-news:
    name: nfx-news
    driver: bridge
'''


def main() -> None:
    write("docker-compose.dev.yml", compose_dev().replace("\\\\n", "\\n").replace("\\n", "\n"))
    write("docker-compose.yml", compose_prod().replace("\\\\n", "\\n").replace("\\n", "\n"))
    for name, title, grpc_var in MODULES:
        write(f"modules/{name}/config/config.go", CONFIG_GO.format(name=name))
        write(f"modules/{name}/config/types.go", TYPES_GO.format())
        write(f"modules/{name}/server/server.go", SERVER_GO.format(name=name))
        write(f"modules/{name}/application/resource/service.go", RESOURCE_SERVICE.format())
        write(f"modules/{name}/interfaces/pipeline/server.go", PIPELINE)
        write(f"modules/{name}/interfaces/messaging/server.go", MESSAGING)
        write(f"modules/{name}/interfaces/http/handler/i18n.go", I18N_HANDLER.format())
        write(f"modules/{name}/interfaces/grpc/handler/health.go", HEALTH_GRPC.format(name=name))
        write(f"inputs/{name}/configuration/dev.toml", toml_for(name, title, grpc_var, "dev"))
        write(f"inputs/{name}/configuration/prod.toml", toml_for(name, title, grpc_var, "prod"))
        for mode in MODES:
            write(f"inputs/{name}/{mode}/main.go", main_go(name, mode))
            write(f"inputs/{name}/{mode}/Dockerfile", dockerfile(name, mode))
            write(f"inputs/{name}/{mode}/.air.toml", air_toml(name, mode))
    print("boilerplate generated")


if __name__ == "__main__":
    main()
