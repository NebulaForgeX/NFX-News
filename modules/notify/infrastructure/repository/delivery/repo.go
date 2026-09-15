package delivery

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/notify/domain/delivery"
	"nfxnews/modules/notify/infrastructure/rdb/models"
)

type h struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *delivery.Repo {
	x := &h{db: db}
	return &delivery.Repo{Create: x, Update: x}
}
func toM(d *delivery.Delivery) *models.Delivery {
	st := d.State()
	return &models.Delivery{ID: st.ID, ChannelID: st.ChannelID, ReportID: st.ReportID, Status: st.Status, ErrorMessage: st.ErrorMessage, CreatedAt: st.CreatedAt, SentAt: st.SentAt}
}
func (x *h) New(ctx context.Context, d *delivery.Delivery) error {
	return x.db.WithContext(ctx).Create(toM(d)).Error
}
func (x *h) Generic(ctx context.Context, d *delivery.Delivery) error {
	return x.db.WithContext(ctx).Save(toM(d)).Error
}
