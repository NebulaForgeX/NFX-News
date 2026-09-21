package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"nfxnews/modules/system/config"
	"nfxnews/modules/system/server"
	"nfxnews/pkgs/connections/otelx"
	"nfxnews/pkgs/env"
	"nfxnews/pkgs/logx"

	"go.uber.org/zap"
)

func main() {
	envStr := flag.String("env", "dev", "Environment (dev/secure)")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(ctx, env.Env(*envStr))
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	if err := logx.Init(cfg.Logger, "system-base-service", env.Env(*envStr)); err != nil {
		log.Fatalf("logger init failed: %v", err)
	}
	defer logx.Sync()

	otelShutdown, err := otelx.Init(ctx, cfg.OTEL, "system", env.Env(*envStr))
	if err != nil {
		log.Fatalf("otel init failed: %v", err)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

	if err := server.RunServer(ctx, cfg); err != nil && !errors.Is(err, context.Canceled) {
		logx.L().Fatal("base server stopped with error", zap.Error(err))
	}

	logx.L().Info("base server shutdown gracefully")
}
