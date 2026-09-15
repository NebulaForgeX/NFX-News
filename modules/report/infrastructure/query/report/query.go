package report

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nfxnews/modules/report/infrastructure/rdb/views"
	reportQuery "nfxnews/modules/report/query/report"
	"nfxnews/pkgs/errx"
)

type kw struct{ db *gorm.DB }
type snap struct{ db *gorm.DB }

func NewQuery(db *gorm.DB) *reportQuery.Query {
	return &reportQuery.Query{Keywords: &kw{db: db}, Snapshots: &snap{db: db}}
}
func (h *kw) All(ctx context.Context) ([]reportQuery.KeywordVO, error) {
	var rows []views.KeywordsActiveView
	if err := h.db.WithContext(ctx).Table(views.KeywordsActiveView{}.TableName()).Order("created_at").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]reportQuery.KeywordVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, reportQuery.KeywordVO{ID: r.ID, GroupName: r.GroupName, Word: r.Word, Kind: r.Kind, CountLimit: r.CountLimit, CreatedAt: r.CreatedAt})
	}
	return out, nil
}
func decodeSnap(r views.SnapshotsActiveView) reportQuery.SnapshotVO {
	vo := reportQuery.SnapshotVO{ID: r.ID, Mode: r.Mode, Title: r.Title, Payload: r.Payload, ItemCount: r.ItemCount, CreatedAt: r.CreatedAt}
	_ = json.Unmarshal(r.Payload, &vo.PayloadObj)
	return vo
}
func (h *snap) Recent(ctx context.Context, limit int) ([]reportQuery.SnapshotVO, error) {
	var rows []views.SnapshotsActiveView
	if err := h.db.WithContext(ctx).Table(views.SnapshotsActiveView{}.TableName()).Order("created_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]reportQuery.SnapshotVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, decodeSnap(r))
	}
	return out, nil
}
func (h *snap) ByID(ctx context.Context, id uuid.UUID) (*reportQuery.SnapshotVO, error) {
	var row views.SnapshotsActiveView
	if err := h.db.WithContext(ctx).Table(views.SnapshotsActiveView{}.TableName()).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.NotFound("REPORT_NOT_FOUND", "report not found")
		}
		return nil, err
	}
	vo := decodeSnap(row)
	return &vo, nil
}
