package keyword

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/report/domain/keyword"
	"nfxnews/modules/report/infrastructure/rdb/models"
)

type h struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *keyword.Repo { return &keyword.Repo{Create: &h{db: db}} }
func (x *h) New(ctx context.Context, k *keyword.Keyword) error {
	st := k.State()
	return x.db.WithContext(ctx).Create(&models.Keyword{ID: st.ID, GroupName: st.GroupName, Word: st.Word, Kind: st.Kind, CountLimit: st.CountLimit, CreatedAt: st.CreatedAt}).Error
}
