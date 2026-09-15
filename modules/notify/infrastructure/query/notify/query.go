package notify

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"nfxnews/modules/notify/infrastructure/rdb/views"
	notifyQuery "nfxnews/modules/notify/query/notify"
)

type ch struct{ db *gorm.DB }
type del struct{ db *gorm.DB }

func NewQuery(db *gorm.DB) *notifyQuery.Query {
	return &notifyQuery.Query{Channels: &ch{db: db}, Deliveries: &del{db: db}}
}
func (h *ch) All(ctx context.Context) ([]notifyQuery.ChannelVO, error) {
	var rows []views.ChannelsActiveView
	if err := h.db.WithContext(ctx).Table(views.ChannelsActiveView{}.TableName()).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]notifyQuery.ChannelVO, 0, len(rows))
	for _, r := range rows {
		vo := notifyQuery.ChannelVO{ID: r.ID, Kind: r.Kind, Name: r.Name, Enabled: r.Enabled, Config: r.Config, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
		_ = json.Unmarshal(r.Config, &vo.ConfigObj)
		out = append(out, vo)
	}
	return out, nil
}
func (h *del) Recent(ctx context.Context, limit int) ([]notifyQuery.DeliveryVO, error) {
	var rows []views.DeliveriesActiveView
	if err := h.db.WithContext(ctx).Table(views.DeliveriesActiveView{}.TableName()).Order("created_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]notifyQuery.DeliveryVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, notifyQuery.DeliveryVO{ID: r.ID, ChannelID: r.ChannelID, ReportID: r.ReportID, Status: r.Status, ErrorMessage: r.ErrorMessage, CreatedAt: r.CreatedAt, SentAt: r.SentAt})
	}
	return out, nil
}
