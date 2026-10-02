package server

import (
	"context"
	"fmt"
	"time"

	authconn "nfxnews/connections/auth"
	resourceApp "nfxnews/modules/source/application/resource"
	sourceapp "nfxnews/modules/source/application/source"
	"nfxnews/modules/source/config"
	"nfxnews/pkgs/cachex"
	"nfxnews/pkgs/health"
	"nfxnews/pkgs/kafkax"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/postgresqlx"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	"nfxnews/pkgs/tokenx"
)

type Dependencies struct {
	healthMgr           *health.Manager
	cache               *cachex.Connection
	postgres            *postgresqlx.Connection
	kafkaConfig         *kafkax.Config
	busPublisher        *eventbus.BusPublisher
	appSvc              *sourceapp.Service
	resourceSvc         *resourceApp.Service
	userTokenVerifier   token.Verifier
	serverTokenVerifier token.Verifier
	errorsLangsPath     string
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
	reg, err := sourceapp.NewRegistryFromWD(".")
	if err != nil {
		return nil, fmt.Errorf("load source catalog: %w", err)
	}
	identityClient, err := authconn.Dial(authconn.GRPCConfig{
		Addr:           cfg.GRPCClient.AuthAddr,
		TokenSecretKey: cfg.Token.SecretKey,
		TokenIssuer:    cfg.Token.Issuer,
		CallerService:  "source",
	})
	if err != nil {
		return nil, fmt.Errorf("dial identity auth: %w", err)
	}
	errorsLangsPath := cfg.I18n.ErrorsLangsPath
	if errorsLangsPath == "" {
		errorsLangsPath = "./errors/langs"
	}
	return &Dependencies{
		healthMgr: healthMgr, postgres: postgres, cache: cacheConn, kafkaConfig: &kafkaConfig,
		busPublisher:      busPublisher,
		appSvc:            sourceapp.NewService(reg, busPublisher),
		resourceSvc:       resourceApp.NewService(postgres, cacheConn, &kafkaConfig),
		userTokenVerifier: userTokenVerifier, serverTokenVerifier: serverTokenVerifier, errorsLangsPath: errorsLangsPath,
		identityAuth: identityClient,
	}, nil
}

func (d *Dependencies) Cleanup() {
	d.healthMgr.Stop()
	d.postgres.Close()
	d.cache.Close()
	if d.identityAuth != nil {
		_ = d.identityAuth.Close()
	}
}

func (d *Dependencies) AppSvc() *sourceapp.Service           { return d.appSvc }
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
