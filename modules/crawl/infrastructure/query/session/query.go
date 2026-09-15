package session

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nfxnews/modules/crawl/infrastructure/rdb/views"
	sessionQuery "nfxnews/modules/crawl/query/session"
	"nfxnews/pkgs/errx"
)

type handler struct{ db *gorm.DB }

func NewQuery(db *gorm.DB) *sessionQuery.Query { return &sessionQuery.Query{List: &handler{db: db}} }
func toVO(r views.SessionsActiveView) sessionQuery.SessionVO {
	return sessionQuery.SessionVO{ID: r.ID, SourceID: r.SourceID, Status: r.Status, ItemCount: r.ItemCount, ErrorMessage: r.ErrorMessage, StartedAt: r.StartedAt, FinishedAt: r.FinishedAt}
}
func (h *handler) Recent(ctx context.Context, limit int) ([]sessionQuery.SessionVO, error) {
	var rows []views.SessionsActiveView
	if err := h.db.WithContext(ctx).Table(views.SessionsActiveView{}.TableName()).Order("started_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]sessionQuery.SessionVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, toVO(r))
	}
	return out, nil
}
func (h *handler) ByID(ctx context.Context, id uuid.UUID) (*sessionQuery.SessionVO, error) {
	var row views.SessionsActiveView
	if err := h.db.WithContext(ctx).Table(views.SessionsActiveView{}.TableName()).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.NotFound("CRAWL_NOT_FOUND", "crawl session not found")
		}
		return nil, err
	}
	vo := toVO(row)
	return &vo, nil
}
