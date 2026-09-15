package session

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/crawl/domain/session"
	"nfxnews/modules/crawl/infrastructure/rdb/models"
)

type repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *session.Repo {
	h := &repo{db: db}
	return &session.Repo{Create: h, Update: h}
}
func toModel(s *session.Session) *models.Session {
	st := s.State()
	return &models.Session{ID: st.ID, SourceID: st.SourceID, Status: st.Status, ItemCount: st.ItemCount, ErrorMessage: st.ErrorMessage, StartedAt: st.StartedAt, FinishedAt: st.FinishedAt}
}
func (h *repo) New(ctx context.Context, s *session.Session) error {
	return h.db.WithContext(ctx).Create(toModel(s)).Error
}
func (h *repo) Generic(ctx context.Context, s *session.Session) error {
	return h.db.WithContext(ctx).Save(toModel(s)).Error
}
