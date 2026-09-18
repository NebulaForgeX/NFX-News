package systemstate

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/system/domain/systemstate"
	"nfxnews/modules/system/infrastructure/rdb/models"
)

type handler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *systemstate.Repo { return &systemstate.Repo{Create: &handler{db: db}} }
func (handler *handler) New(ctx context.Context, s *systemstate.State) error {
	st := s.Inner()
	return handler.db.WithContext(ctx).Create(&models.State{ID: st.ID, Initialized: st.Initialized, InitializedAt: st.InitializedAt, InitializationVersion: st.InitializationVersion, ResetCount: st.ResetCount, CreatedAt: st.CreatedAt, UpdatedAt: st.UpdatedAt}).Error
}
