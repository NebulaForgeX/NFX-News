package notifyapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"nfxnews/events"
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
	repos  *repofactory.TxRepoFactory
	query  *notifyQuery.Query
	client *http.Client
}

func NewService(repos *repofactory.TxRepoFactory, query *notifyQuery.Query) *Service {
	return &Service{repos: repos, query: query, client: &http.Client{Timeout: 15 * time.Second}}
}

func (s *Service) ListChannels(ctx context.Context) ([]Channel, error) {
	rows, err := s.query.Channels.All(ctx)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) ListDeliveries(ctx context.Context, limit int) ([]Delivery, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.query.Deliveries.Recent(ctx, limit)
}

func (s *Service) UpsertChannel(ctx context.Context, kind, name string, enabled bool, cfg map[string]any) (*Channel, error) {
	raw, _ := json.Marshal(cfg)
	now := time.Now()
	id := uuid.Must(uuid.NewV7())
	if err := s.repos.Channel(transaction.UoW{}).Create.New(ctx, channelDomain.NewFromState(channelDomain.State{
		ID: id, Kind: kind, Name: name, Enabled: enabled, Config: raw, CreatedAt: now, UpdatedAt: now,
	})); err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return &Channel{ID: id, Kind: kind, Name: name, Enabled: enabled, Config: raw, ConfigObj: cfg, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *Service) DispatchReport(ctx context.Context, ev events.ReportGeneratedEvent) (int, error) {
	channels, err := s.ListChannels(ctx)
	if err != nil {
		return 0, err
	}
	queued := 0
	body := fmt.Sprintf("%s\n%s\nitems=%d", ev.Title, ev.Mode, ev.ItemCount)
	deliveries := s.repos.Delivery(transaction.UoW{})
	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}
		now := time.Now()
		d := deliveryDomain.NewFromState(deliveryDomain.State{
			ID: uuid.Must(uuid.NewV7()), ChannelID: ch.ID, Status: "pending", CreatedAt: now,
		})
		if ev.ReportID != "" {
			if id, err := uuid.Parse(ev.ReportID); err == nil {
				st := d.State()
				st.ReportID = &id
				d = deliveryDomain.NewFromState(st)
			}
		}
		_ = deliveries.Create.New(ctx, d)
		if err := s.send(ctx, ch, body); err != nil {
			msg := err.Error()
			d.Mark("failed", &msg, nil)
		} else {
			sent := time.Now()
			d.Mark("sent", nil, &sent)
			queued++
		}
		_ = deliveries.Update.Generic(ctx, d)
	}
	return queued, nil
}

func (s *Service) send(ctx context.Context, ch Channel, text string) error {
	cfg := ch.ConfigObj
	webhook := ""
	if cfg != nil {
		if v, ok := cfg["webhook_url"].(string); ok {
			webhook = v
		}
	}
	if webhook == "" {
		webhook = envWebhook(ch.Kind)
	}
	if webhook == "" {
		return fmt.Errorf("no webhook for %s", ch.Kind)
	}
	payload, _ := json.Marshal(map[string]any{"text": text, "msg_type": "text", "content": map[string]string{"text": text}})
	if ch.Kind == "telegram" {
		chat := os.Getenv("NOTIFY_TELEGRAM_CHAT_ID")
		if cfg != nil {
			if v, ok := cfg["chat_id"].(string); ok {
				chat = v
			}
		}
		token := os.Getenv("NOTIFY_TELEGRAM_BOT_TOKEN")
		if token != "" && chat != "" {
			webhook = "https://api.telegram.org/bot" + token + "/sendMessage"
			payload, _ = json.Marshal(map[string]any{"chat_id": chat, "text": text})
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook %s: %s", ch.Kind, resp.Status)
	}
	return nil
}

func envWebhook(kind string) string {
	switch kind {
	case "feishu":
		return os.Getenv("NOTIFY_FEISHU_WEBHOOK_URL")
	case "dingtalk":
		return os.Getenv("NOTIFY_DINGTALK_WEBHOOK_URL")
	case "wework":
		return os.Getenv("NOTIFY_WEWORK_WEBHOOK_URL")
	case "slack":
		return os.Getenv("NOTIFY_SLACK_WEBHOOK_URL")
	case "bark":
		return os.Getenv("NOTIFY_BARK_URL")
	case "ntfy":
		topic := os.Getenv("NOTIFY_NTFY_TOPIC")
		if topic == "" {
			return ""
		}
		base := os.Getenv("NOTIFY_NTFY_SERVER_URL")
		if base == "" {
			base = "https://ntfy.sh"
		}
		return strings.TrimRight(base, "/") + "/" + topic
	default:
		return ""
	}
}
