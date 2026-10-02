package preference

import (
	"context"

	"nfxnews/modules/news/domain/preference"
	"nfxnews/modules/news/infrastructure/rdb/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type createHandler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *preference.Repo { return &preference.Repo{Create: &createHandler{db: db}} }

func (h *createHandler) Upsert(ctx context.Context, p *preference.Preference) error {
	st := p.State()
	row := models.Preference{AccountID: st.AccountID, ProfileID: st.ProfileID, ColumnOrder: st.ColumnOrder, Payload: st.Payload, UpdatedAt: st.UpdatedAt}
	return h.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "profile_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"column_order", "payload", "updated_at"}),
	}).Create(&row).Error
}
