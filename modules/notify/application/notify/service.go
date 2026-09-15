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
	"nfxnews/pkgs/errx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Channel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Kind      string     `gorm:"type:varchar(32)" json:"kind"`
	Name      string     `gorm:"type:varchar(128)" json:"name"`
	Enabled   bool       `json:"enabled"`
	Config    []byte     `gorm:"type:jsonb" json:"-"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	ConfigObj map[string]any `gorm:"-" json:"config"`
}

func (Channel) TableName() string { return "notify.channels" }

type Delivery struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ChannelID    uuid.UUID  `json:"channel_id"`
	ReportID     *uuid.UUID `json:"report_id"`
	Status       string      `json:"status"`
	ErrorMessage *string     `json:"error_message"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	SentAt       *time.Time  `json:"sent_at"`
}

func (Delivery) TableName() string { return "notify.deliveries" }

type Service struct {
	db     *gorm.DB
	client *http.Client
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db, client: &http.Client{Timeout: 15 * time.Second}}
}

func (s *Service) ListChannels(ctx context.Context) ([]Channel, error) {
	var rows []Channel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	for i := range rows {
		_ = json.Unmarshal(rows[i].Config, &rows[i].ConfigObj)
	}
	return rows, nil
}

func (s *Service) ListDeliveries(ctx context.Context, limit int) ([]Delivery, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var rows []Delivery
	if err := s.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) UpsertChannel(ctx context.Context, kind, name string, enabled bool, cfg map[string]any) (*Channel, error) {
	raw, _ := json.Marshal(cfg)
	row := Channel{ID: uuid.Must(uuid.NewV7()), Kind: kind, Name: name, Enabled: enabled, Config: raw}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	row.ConfigObj = cfg
	return &row, nil
}

func (s *Service) DispatchReport(ctx context.Context, ev events.ReportGeneratedEvent) (int, error) {
	channels, err := s.ListChannels(ctx)
	if err != nil {
		return 0, err
	}
	queued := 0
	body := fmt.Sprintf("%s\n%s\nitems=%d", ev.Title, ev.Mode, ev.ItemCount)
	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}
		d := Delivery{ID: uuid.Must(uuid.NewV7()), ChannelID: ch.ID, Status: "pending"}
		if ev.ReportID != "" {
			if id, err := uuid.Parse(ev.ReportID); err == nil {
				d.ReportID = &id
			}
		}
		_ = s.db.WithContext(ctx).Create(&d).Error
		if err := s.send(ctx, ch, body); err != nil {
			msg := err.Error()
			d.Status = "failed"
			d.ErrorMessage = &msg
		} else {
			now := time.Now()
			d.Status = "sent"
			d.SentAt = &now
			queued++
		}
		_ = s.db.WithContext(ctx).Save(&d).Error
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
