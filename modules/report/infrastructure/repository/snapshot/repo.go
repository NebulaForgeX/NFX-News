package snapshot

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/report/domain/snapshot"
	"nfxnews/modules/report/infrastructure/rdb/models"
)

type h struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *snapshot.Repo { return &snapshot.Repo{Create: &h{db: db}} }
func (x *h) New(ctx context.Context, s *snapshot.Snapshot) error {
	st := s.State()
	return x.db.WithContext(ctx).Create(&models.Snapshot{ID: st.ID, Mode: st.Mode, Title: st.Title, Payload: st.Payload, ItemCount: st.ItemCount, CreatedAt: st.CreatedAt}).Error
}
