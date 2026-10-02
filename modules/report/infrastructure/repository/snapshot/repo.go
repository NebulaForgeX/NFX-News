package snapshot

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/report/domain/snapshot"
	"nfxnews/modules/report/infrastructure/rdb/models"
)

type handler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *snapshot.Repo { return &snapshot.Repo{Create: &handler{db: db}} }
func (handler *handler) New(ctx context.Context, s *snapshot.Snapshot) error {
	st := s.State()
	return handler.db.WithContext(ctx).Create(&models.Snapshot{ID: st.ID, AccountID: st.AccountID, ProfileID: st.ProfileID, Mode: st.Mode, Title: st.Title, Payload: st.Payload, ItemCount: st.ItemCount, CreatedAt: st.CreatedAt}).Error
}
