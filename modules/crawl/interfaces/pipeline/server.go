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
	sec := 30
	if v := os.Getenv("CRAWL_SCHEDULE_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			sec = n
		}
	}
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	defer ticker.Stop()
	if err := r.svc.FetchDue(ctx); err != nil {
		logx.S().Warnf("scheduled crawl failed: %v", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.svc.FetchDue(ctx); err != nil {
				logx.S().Warnf("scheduled crawl failed: %v", err)
			}
		}
	}
}

func (r *Router) Close() error { return r.Router.Close() }
