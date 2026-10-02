package session

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/crawl/domain/session"
	"nfxnews/modules/crawl/infrastructure/rdb/models"
)

type handler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *session.Repo {
	impl := &handler{db: db}
	return &session.Repo{Create: impl, Update: impl}
}
func toModel(s *session.Session) *models.Session {
	st := s.State()
	return &models.Session{ID: st.ID, AccountID: st.AccountID, ProfileID: st.ProfileID, SourceID: st.SourceID, Status: st.Status, ItemCount: st.ItemCount, ErrorMessage: st.ErrorMessage, StartedAt: st.StartedAt, FinishedAt: st.FinishedAt}
}
func (handler *handler) New(ctx context.Context, s *session.Session) error {
	return handler.db.WithContext(ctx).Create(toModel(s)).Error
}
func (handler *handler) Generic(ctx context.Context, s *session.Session) error {
	return handler.db.WithContext(ctx).Save(toModel(s)).Error
}
