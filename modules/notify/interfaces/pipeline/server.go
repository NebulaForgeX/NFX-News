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
