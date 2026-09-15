package item

import (
	"context"
	"encoding/json"
	"errors"

	"nfxnews/modules/news/infrastructure/rdb/views"
	itemQuery "nfxnews/modules/news/query/item"

	"gorm.io/gorm"
)

type listHandler struct{ db *gorm.DB }
type prefHandler struct{ db *gorm.DB }

func NewQuery(db *gorm.DB) *itemQuery.Query {
	return &itemQuery.Query{List: &listHandler{db: db}, Pref: &prefHandler{db: db}}
}

func toVO(r views.ItemsActiveView) itemQuery.ItemVO {
	v := itemQuery.ItemVO{ID: r.ID, SourceID: r.SourceID, OriginalID: r.OriginalID, Title: r.Title, URL: r.URL, UpdatedAt: r.UpdatedAt}
	if r.MobileURL != nil {
		v.MobileURL = *r.MobileURL
	}
	if r.PubDate != nil {
		v.PubDate = r.PubDate.Unix()
	}
	if len(r.Extra) > 0 {
		_ = json.Unmarshal(r.Extra, &v.Extra)
	}
	return v
}

func (h *listHandler) BySource(ctx context.Context, sourceID string, limit int) ([]itemQuery.ItemVO, error) {
	q := h.db.WithContext(ctx).Table(views.ItemsActiveView{}.TableName()).Order("updated_at desc").Limit(limit)
	if sourceID != "" {
		q = q.Where("source_id = ?", sourceID)
	}
	var rows []views.ItemsActiveView
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]itemQuery.ItemVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, toVO(r))
	}
	return out, nil
}

func (h *listHandler) Search(ctx context.Context, qstr string, limit int) ([]itemQuery.ItemVO, error) {
	q := h.db.WithContext(ctx).Table(views.ItemsActiveView{}.TableName()).Order("updated_at desc").Limit(limit)
	if qstr != "" {
		q = q.Where("title ILIKE ?", "%"+qstr+"%")
	}
	var rows []views.ItemsActiveView
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]itemQuery.ItemVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, toVO(r))
	}
	return out, nil
}

func (h *prefHandler) ByAccountProfile(ctx context.Context, accountID, profileID string) (*itemQuery.PreferenceVO, error) {
	var row views.ProfilePreferencesActiveView
	err := h.db.WithContext(ctx).Table(views.ProfilePreferencesActiveView{}.TableName()).
		Where("account_id = ? AND profile_id = ?", accountID, profileID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &itemQuery.PreferenceVO{ColumnOrder: row.ColumnOrder, Payload: row.Payload}, nil
}
