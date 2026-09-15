package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"nfxnews/modules/notify/config"
	"nfxnews/modules/notify/server"
	"nfxnews/pkgs/connections/otelx"
	"nfxnews/pkgs/env"
	"nfxnews/pkgs/logx"

	"go.uber.org/zap"
)

func main() {
	envStr := flag.String("env", "prod", "Environment (dev/prod)")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(ctx, env.Env(*envStr))
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	if err := logx.Init(cfg.Logger, "notify-messaging-service", env.Env(*envStr)); err != nil {
		log.Fatalf("logger init failed: %v", err)
	}
	defer logx.Sync()

	otelShutdown, err := otelx.Init(ctx, cfg.OTEL, "notify", env.Env(*envStr))
	if err != nil {
		log.Fatalf("otel init failed: %v", err)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

	if err := server.RunMessaging(ctx, cfg); err != nil && !errors.Is(err, context.Canceled) {
		logx.L().Fatal("messaging server stopped with error", zap.Error(err))
	}

	logx.L().Info("messaging server shutdown gracefully")
}
