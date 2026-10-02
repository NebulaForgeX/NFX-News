package item

import (
	"context"

	"nfxnews/modules/news/domain/item"
	"nfxnews/modules/news/infrastructure/rdb/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type createHandler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *item.Repo { return &item.Repo{Create: &createHandler{db: db}} }

func (h *createHandler) Upsert(ctx context.Context, i *item.Item) error {
	st := i.State()
	row := models.Item{ID: st.ID, SourceID: st.SourceID, OriginalID: st.OriginalID, Title: st.Title, URL: st.URL, MobileURL: st.MobileURL, PubDate: st.PubDate, Extra: st.Extra}
	return h.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&row).Error
}
