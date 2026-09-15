package crawlapp

import (
	"context"
	"time"

	"nfxnews/pkgs/errx"
	sourcepb "nfxnews/protos/gen/source"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	SourceID     *string     `gorm:"type:varchar(64)" json:"source_id"`
	Status       string       `gorm:"type:varchar(32)" json:"status"`
	ItemCount    int          `json:"item_count"`
	ErrorMessage *string      `gorm:"type:text" json:"error_message"`
	StartedAt    time.Time    `gorm:"autoCreateTime" json:"started_at"`
	FinishedAt   *time.Time   `json:"finished_at"`
}

func (Session) TableName() string { return "crawl.sessions" }

type Service struct {
	db     *gorm.DB
	source sourcepb.SourceServiceClient
}

func NewService(db *gorm.DB, source sourcepb.SourceServiceClient) *Service {
	return &Service{db: db, source: source}
}

func (s *Service) Trigger(ctx context.Context, sourceID string) (*Session, error) {
	sess := Session{ID: uuid.Must(uuid.NewV7()), Status: "running"}
	if sourceID != "" {
		sess.SourceID = &sourceID
	}
	if err := s.db.WithContext(ctx).Create(&sess).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	if s.source != nil && sourceID != "" {
		resp, err := s.source.FetchSource(ctx, &sourcepb.FetchSourceRequest{SourceId: sourceID})
		now := time.Now()
		sess.FinishedAt = &now
		if err != nil {
			msg := err.Error()
			sess.Status = "failed"
			sess.ErrorMessage = &msg
		} else {
			sess.Status = "ok"
			sess.ItemCount = len(resp.GetItems())
		}
		_ = s.db.WithContext(ctx).Save(&sess).Error
	}
	return &sess, nil
}

func (s *Service) TriggerAll(ctx context.Context) (*Session, error) {
	if s.source == nil {
		return s.Trigger(ctx, "")
	}
	list, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
	if err != nil {
		return s.Trigger(ctx, "")
	}
	var last *Session
	for _, src := range list.GetSources() {
		last, err = s.Trigger(ctx, src.GetId())
		if err != nil {
			return last, err
		}
	}
	if last == nil {
		return s.Trigger(ctx, "")
	}
	return last, nil
}

func (s *Service) List(ctx context.Context, limit int) ([]Session, error) {
	if limit <= 0 {
		limit = 20
	}
	var rows []Session
	if err := s.db.WithContext(ctx).Order("started_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Session, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errx.ErrInvalidParams.WithCause(err)
	}
	var row Session
	if err := s.db.WithContext(ctx).First(&row, "id = ?", uid).Error; err != nil {
		return nil, errx.NotFound("CRAWL_NOT_FOUND", "crawl session not found")
	}
	return &row, nil
}
