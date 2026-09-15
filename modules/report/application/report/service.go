package reportapp

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"nfxnews/events"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/kafkax/eventbus"
	newspb "nfxnews/protos/gen/news"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Keyword struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	GroupName  string     `gorm:"type:varchar(128)" json:"group_name"`
	Word       string     `gorm:"type:varchar(255)" json:"word"`
	Kind       string     `gorm:"type:varchar(16)" json:"kind"`
	CountLimit int        `json:"count_limit"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (Keyword) TableName() string { return "report.keywords" }

type Snapshot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Mode      string     `gorm:"type:varchar(32)" json:"mode"`
	Title     string     `gorm:"type:varchar(255)" json:"title"`
	Payload   []byte    `gorm:"type:jsonb" json:"-"`
	ItemCount int        `json:"item_count"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	PayloadObj any      `gorm:"-" json:"payload"`
}

func (Snapshot) TableName() string { return "report.snapshots" }

type Service struct {
	db   *gorm.DB
	news newspb.NewsServiceClient
	pub  *eventbus.BusPublisher
}

func NewService(db *gorm.DB, news newspb.NewsServiceClient, pub *eventbus.BusPublisher) *Service {
	return &Service{db: db, news: news, pub: pub}
}

func (s *Service) ListKeywords(ctx context.Context) ([]Keyword, error) {
	var rows []Keyword
	if err := s.db.WithContext(ctx).Order("created_at").Find(&rows).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) AddKeyword(ctx context.Context, group, word, kind string, limit int) (*Keyword, error) {
	word = strings.TrimSpace(word)
	if word == "" {
		return nil, errx.ErrInvalidBody.WithMsg("word required")
	}
	if kind == "" {
		kind = "include"
	}
	row := Keyword{ID: uuid.Must(uuid.NewV7()), GroupName: group, Word: word, Kind: kind, CountLimit: limit}
	if row.GroupName == "" {
		row.GroupName = "default"
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return &row, nil
}

func (s *Service) Generate(ctx context.Context, mode string) (*Snapshot, error) {
	if mode == "" {
		mode = "daily"
	}
	keywords, err := s.ListKeywords(ctx)
	if err != nil {
		return nil, err
	}
	matched := []map[string]any{}
	if s.news != nil {
		resp, err := s.news.SearchNews(ctx, &newspb.SearchNewsRequest{Query: "", Limit: 200})
		if err == nil {
			for _, item := range resp.GetItems() {
				if matchKeywords(item.GetTitle(), keywords) {
					matched = append(matched, map[string]any{
						"id": item.GetId(), "title": item.GetTitle(), "url": item.GetUrl(), "source_id": item.GetSourceId(),
					})
				}
			}
		}
	}
	payload, _ := json.Marshal(map[string]any{"mode": mode, "items": matched})
	snap := Snapshot{
		ID: uuid.Must(uuid.NewV7()), Mode: mode, Title: strings.ToUpper(mode) + " report",
		Payload: payload, ItemCount: len(matched),
	}
	if err := s.db.WithContext(ctx).Create(&snap).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	_ = json.Unmarshal(snap.Payload, &snap.PayloadObj)
	if s.pub != nil {
		_ = eventbus.PublishEvent(ctx, s.pub, events.ReportGeneratedEvent{
			ReportID: snap.ID.String(), Mode: mode, Title: snap.Title, ItemCount: snap.ItemCount, Payload: string(payload),
		})
	}
	return &snap, nil
}

func (s *Service) ListSnapshots(ctx context.Context, limit int) ([]Snapshot, error) {
	if limit <= 0 {
		limit = 20
	}
	var rows []Snapshot
	if err := s.db.WithContext(ctx).Order("created_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	for i := range rows {
		_ = json.Unmarshal(rows[i].Payload, &rows[i].PayloadObj)
	}
	return rows, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Snapshot, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errx.ErrInvalidParams.WithCause(err)
	}
	var row Snapshot
	if err := s.db.WithContext(ctx).First(&row, "id = ?", uid).Error; err != nil {
		return nil, errx.NotFound("REPORT_NOT_FOUND", "report not found")
	}
	_ = json.Unmarshal(row.Payload, &row.PayloadObj)
	return &row, nil
}

func matchKeywords(title string, keywords []Keyword) bool {
	if len(keywords) == 0 {
		return true
	}
	t := strings.ToLower(title)
	excluded := false
	included := false
	hasInclude := false
	for _, k := range keywords {
		w := strings.ToLower(strings.TrimPrefix(k.Word, "+"))
		w = strings.TrimPrefix(w, "!")
		if k.Kind == "exclude" || strings.HasPrefix(k.Word, "!") {
			if strings.Contains(t, w) {
				excluded = true
			}
			continue
		}
		hasInclude = true
		if strings.Contains(t, w) {
			included = true
		}
	}
	if excluded {
		return false
	}
	if !hasInclude {
		return true
	}
	return included
}
