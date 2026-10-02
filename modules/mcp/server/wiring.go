package server

import (
	"context"
	"fmt"
	"time"

	authconn "nfxnews/connections/auth"
	crawlconn "nfxnews/connections/crawl"
	newsconn "nfxnews/connections/news"
	reportconn "nfxnews/connections/report"
	sourceconn "nfxnews/connections/source"
	mcpapp "nfxnews/modules/mcp/application/mcp"
	resourceApp "nfxnews/modules/mcp/application/resource"
	"nfxnews/modules/mcp/config"
	toolcallRepo "nfxnews/modules/mcp/infrastructure/repository/toolcall"
	"nfxnews/pkgs/cachex"
	"nfxnews/pkgs/grpcx"
	"nfxnews/pkgs/health"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/postgresqlx"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	"nfxnews/pkgs/tokenx"
	crawlpb "nfxnews/protos/gen/crawl"
	newspb "nfxnews/protos/gen/news"
	reportpb "nfxnews/protos/gen/report"
	sourcepb "nfxnews/protos/gen/source"

	"google.golang.org/grpc"
)

type Dependencies struct {
	healthMgr           *health.Manager
	cache               *cachex.Connection
	postgres            *postgresqlx.Connection
	kafkaConfig         *kafkax.Config
	busPublisher        *eventbus.BusPublisher
	appSvc              *mcpapp.Service
	resourceSvc         *resourceApp.Service
	userTokenVerifier   token.Verifier
	serverTokenVerifier token.Verifier
	errorsLangsPath     string
	conns               []*grpc.ClientConn
	identityAuth        *authconn.Client
}

func NewDeps(ctx context.Context, cfg *config.Config) (*Dependencies, error) {
	postgres, err := postgresqlx.Init(ctx, cfg.PostgreSQL)
	if err != nil {
		return nil, fmt.Errorf("init PostgreSQL: %w", err)
	}
	cacheConn, err := cachex.InitConn(ctx, cfg.Cache)
	if err != nil {
		return nil, fmt.Errorf("init Redis: %w", err)
	}
	healthMgr := health.NewManager(ctx, 30*time.Second)
	healthMgr.Register(postgres)
	healthMgr.Register(cacheConn)
	kafkaConfig := cfg.KafkaConfig
	busPublisher, err := kafkax.NewPublisher(&kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("kafka publisher: %w", err)
	}
	tokenxInstance := tokenx.New(cfg.Token)
	userTokenVerifier := &tokenxVerifierAdapter{tokenx: tokenxInstance}
	serverTokenVerifier := servertoken.NewVerifier(
		&servertoken.HMACSigner{Key: []byte(cfg.Token.SecretKey)},
		cfg.Token.Issuer,
		servertoken.WithAllowedSkew(5*time.Second),
	)
	provider := servertoken.NewProvider(
		&servertoken.HMACSigner{Key: []byte(cfg.Token.SecretKey)},
		cfg.Token.Issuer,
		"mcp",
	)
	identityClient, err := authconn.Dial(authconn.GRPCConfig{
		Addr:           cfg.GRPCClient.AuthAddr,
		TokenSecretKey: cfg.Token.SecretKey,
		TokenIssuer:    cfg.Token.Issuer,
		CallerService:  "mcp",
	})
	if err != nil {
		return nil, fmt.Errorf("dial identity auth: %w", err)
	}
	errorsLangsPath := cfg.I18n.ErrorsLangsPath
	if errorsLangsPath == "" {
		errorsLangsPath = "./errors/langs"
	}
	d := &Dependencies{
		healthMgr: healthMgr, postgres: postgres, cache: cacheConn, kafkaConfig: &kafkaConfig,
		busPublisher:      busPublisher,
		resourceSvc:       resourceApp.NewService(postgres, cacheConn, &kafkaConfig),
		userTokenVerifier: userTokenVerifier, serverTokenVerifier: serverTokenVerifier, errorsLangsPath: errorsLangsPath,
		identityAuth: identityClient,
	}
	var newsClient newspb.NewsServiceClient
	var reportClient reportpb.ReportServiceClient
	var sourceClient sourcepb.SourceServiceClient
	var crawlClient crawlpb.CrawlServiceClient
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
	if conn, err := grpcx.Dial(cfg.GRPCClient.CrawlAddr, "mcp", cfg.Token); err == nil {
		d.conns = append(d.conns, conn)
		crawlClient = crawlconn.New(conn)
	}
	d.appSvc = mcpapp.NewService(toolcallRepo.NewRepo(postgres.DB()), newsClient, reportClient, sourceClient, crawlClient)
	_ = provider
	return d, nil
}

func (d *Dependencies) Cleanup() {
	d.healthMgr.Stop()
	d.postgres.Close()
	d.cache.Close()
	if d.identityAuth != nil {
		_ = d.identityAuth.Close()
	}
	for _, c := range d.conns {
		_ = c.Close()
	}
}

func (d *Dependencies) AppSvc() *mcpapp.Service              { return d.appSvc }
func (d *Dependencies) ResourceSvc() *resourceApp.Service    { return d.resourceSvc }
func (d *Dependencies) UserTokenVerifier() token.Verifier    { return d.userTokenVerifier }
func (d *Dependencies) ServerTokenVerifier() token.Verifier  { return d.serverTokenVerifier }
func (d *Dependencies) KafkaConfig() *kafkax.Config          { return d.kafkaConfig }
func (d *Dependencies) BusPublisher() *eventbus.BusPublisher { return d.busPublisher }
func (d *Dependencies) ErrorsLangsPath() string              { return d.errorsLangsPath }
func (d *Dependencies) AuthClient() *authconn.Client         { return d.identityAuth }

type tokenxVerifierAdapter struct{ tokenx *tokenx.Tokenx }

func (a *tokenxVerifierAdapter) Verify(ctx context.Context, tokenStr string) (*token.Claims, error) {
	claims, err := a.tokenx.VerifyAccessToken(tokenStr)
	if err != nil {
		return nil, err
	}
	return &token.Claims{Registered: claims.RegisteredClaims, Raw: map[string]any{
		"account_id":    claims.AccountID,
		"profile_id":    claims.ProfileID,
		"profile_scope": claims.ProfileScope,
	}}, nil
}
