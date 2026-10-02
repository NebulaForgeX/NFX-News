package notifyapp

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"nfxnews/events"
	"nfxnews/modules/notify/application/notify/channels"
	channelDomain "nfxnews/modules/notify/domain/channel"
	deliveryDomain "nfxnews/modules/notify/domain/delivery"
	repofactory "nfxnews/modules/notify/infrastructure/repository/factory"
	notifyQuery "nfxnews/modules/notify/query/notify"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/transaction"

	"github.com/google/uuid"
)

type Channel = notifyQuery.ChannelVO
type Delivery = notifyQuery.DeliveryVO

type Service struct {
	repoFactory *repofactory.TxRepoFactory
	query       *notifyQuery.Query
	client      *http.Client
	mu          sync.Mutex
	pushedDay   string
}

func NewService(repoFactory *repofactory.TxRepoFactory, query *notifyQuery.Query) *Service {
	return &Service{repoFactory: repoFactory, query: query, client: &http.Client{Timeout: 15 * time.Second}}
}

func (s *Service) ListKinds() []string {
	return channelDomain.Kinds()
}

func (s *Service) ListChannels(ctx context.Context, accountID string) ([]Channel, error) {
	rows, err := s.query.Channels.All(ctx, accountID)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) ListDeliveries(ctx context.Context, accountID string, limit int) ([]Delivery, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.query.Deliveries.Recent(ctx, accountID, limit)
}

func (s *Service) UpsertChannel(ctx context.Context, accountID, profileID, kind, name string, enabled bool, cfg map[string]any) (*Channel, error) {
	if accountID == "" {
		return nil, errx.Unauthorized("INVALID_TOKEN", "missing account")
	}
	normalized, err := channelDomain.NormalizeKind(kind)
	if err != nil {
		return nil, err
	}
	kind = normalized
	name = strings.TrimSpace(name)
	if name == "" {
		name = kind
	}
	if cfg == nil {
		cfg = map[string]any{}
	}
	raw, _ := json.Marshal(cfg)
	var aid, pid *string
	if accountID != "" {
		aid = &accountID
	}
	if profileID != "" {
		pid = &profileID
	}
	existing, err := s.ListChannels(ctx, accountID)
	if err != nil {
		return nil, err
	}
	channelRepo := s.repoFactory.Channel(transaction.UoW{})
	now := time.Now().UTC()
	for _, row := range existing {
		if row.Kind != kind || row.Name != name {
			continue
		}
		ent := channelDomain.NewFromState(channelDomain.State{
			ID: row.ID, AccountID: aid, ProfileID: pid, Kind: kind, Name: name, Enabled: enabled, Config: raw, CreatedAt: row.CreatedAt, UpdatedAt: now,
		})
		if err := channelRepo.Update.Generic(ctx, ent); err != nil {
			return nil, errx.ErrInternal.WithCause(err)
		}
		st := ent.State()
		return &Channel{ID: st.ID, Kind: st.Kind, Name: st.Name, Enabled: st.Enabled, Config: st.Config, ConfigObj: cfg, CreatedAt: st.CreatedAt, UpdatedAt: st.UpdatedAt}, nil
	}
	ent, err := channelDomain.New(aid, pid, kind, name, enabled, raw)
	if err != nil {
		return nil, err
	}
	if err := channelRepo.Create.New(ctx, ent); err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	st := ent.State()
	return &Channel{ID: st.ID, Kind: st.Kind, Name: st.Name, Enabled: st.Enabled, Config: st.Config, ConfigObj: cfg, CreatedAt: st.CreatedAt, UpdatedAt: st.UpdatedAt}, nil
}

func (s *Service) DispatchReport(ctx context.Context, ev events.ReportGeneratedEvent) (int, error) {
	channelsList, err := s.ListChannels(ctx, ev.AccountID)
	if err != nil {
		return 0, err
	}
	enabled := make([]Channel, 0, len(channelsList))
	for _, ch := range channelsList {
		if ch.Enabled {
			enabled = append(enabled, ch)
		}
	}
	if len(enabled) == 0 {
		for _, kind := range channels.EnvKinds() {
			enabled = append(enabled, Channel{Kind: kind, Name: kind, Enabled: true, ConfigObj: map[string]any{}})
		}
	}
	if len(enabled) == 0 {
		return 0, nil
	}
	windowCfg := map[string]any{}
	if len(enabled) > 0 && enabled[0].ConfigObj != nil {
		windowCfg = enabled[0].ConfigObj
	}
	if !channels.InPushWindow(time.Now(), windowCfg) {
		return 0, nil
	}
	if channels.OncePerDay(windowCfg) && s.alreadyPushedToday() {
		return 0, nil
	}
	body := channels.ReportBody(ev)
	deliveryRepo := s.repoFactory.Delivery(transaction.UoW{})
	queued := 0
	for _, ch := range enabled {
		now := time.Now().UTC()
		d := deliveryDomain.NewFromState(deliveryDomain.State{
			ID: uuid.Must(uuid.NewV7()), ChannelID: ch.ID, Status: "pending", CreatedAt: now,
		})
		if ev.AccountID != "" {
			st := d.State()
			aid := ev.AccountID
			st.AccountID = &aid
			if ev.ProfileID != "" {
				pid := ev.ProfileID
				st.ProfileID = &pid
			}
			d = deliveryDomain.NewFromState(st)
		}
		if ev.ReportID != "" {
			if id, err := uuid.Parse(ev.ReportID); err == nil {
				st := d.State()
				st.ReportID = &id
				d = deliveryDomain.NewFromState(st)
			}
		}
		if ch.ID != uuid.Nil {
			_ = deliveryRepo.Create.New(ctx, d)
		}
		cfg := ch.ConfigObj
		if cfg == nil {
			cfg = map[string]any{}
		}
		if err := channels.SendAll(ctx, s.client, ch.Kind, cfg, body); err != nil {
			msg := err.Error()
			d.Mark("failed", &msg, nil)
		} else {
			sent := time.Now().UTC()
			d.Mark("sent", nil, &sent)
			queued++
		}
		if ch.ID != uuid.Nil {
			_ = deliveryRepo.Update.Generic(ctx, d)
		}
	}
	if queued > 0 {
		s.markPushedToday()
	}
	return queued, nil
}

func (s *Service) alreadyPushedToday() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return s.pushedDay == time.Now().In(loc).Format("20060102")
}

func (s *Service) markPushedToday() {
	s.mu.Lock()
	defer s.mu.Unlock()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	s.pushedDay = time.Now().In(loc).Format("20060102")
}
