package messaging

import (
	"context"

	"nfxnews/pkgs/logx"
)

type Deps interface{}

type Router struct{}

func NewServer(_ Deps) (*Router, error) {
	return &Router{}, nil
}

func (r *Router) RegisterRoutes() {}

func (r *Router) Run(ctx context.Context) error {
	logx.S().Info("Kafka messaging router idle (no RabbitMQ); pipeline owns event consumption")
	<-ctx.Done()
	return ctx.Err()
}

func (r *Router) Close() error { return nil }
