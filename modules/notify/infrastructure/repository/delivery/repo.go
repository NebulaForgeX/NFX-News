package delivery

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/notify/domain/delivery"
	"nfxnews/modules/notify/infrastructure/rdb/models"
)

type handler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *delivery.Repo {
	impl := &handler{db: db}
	return &delivery.Repo{Create: impl, Update: impl}
}
func toM(d *delivery.Delivery) *models.Delivery {
	st := d.State()
	return &models.Delivery{ID: st.ID, AccountID: st.AccountID, ProfileID: st.ProfileID, ChannelID: st.ChannelID, ReportID: st.ReportID, Status: st.Status, ErrorMessage: st.ErrorMessage, CreatedAt: st.CreatedAt, SentAt: st.SentAt}
}
func (handler *handler) New(ctx context.Context, d *delivery.Delivery) error {
	return handler.db.WithContext(ctx).Create(toM(d)).Error
}
func (handler *handler) Generic(ctx context.Context, d *delivery.Delivery) error {
	return handler.db.WithContext(ctx).Save(toM(d)).Error
}
