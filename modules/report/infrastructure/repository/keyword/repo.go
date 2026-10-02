package keyword

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/report/domain/keyword"
	"nfxnews/modules/report/infrastructure/rdb/models"
)

type handler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *keyword.Repo { return &keyword.Repo{Create: &handler{db: db}} }
func (handler *handler) New(ctx context.Context, k *keyword.Keyword) error {
	st := k.State()
	return handler.db.WithContext(ctx).Create(&models.Keyword{ID: st.ID, AccountID: st.AccountID, ProfileID: st.ProfileID, GroupName: st.GroupName, Word: st.Word, Kind: st.Kind, CountLimit: st.CountLimit, CreatedAt: st.CreatedAt}).Error
}
