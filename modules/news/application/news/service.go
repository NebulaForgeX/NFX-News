package newsapp

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"nfxnews/events"
	"nfxnews/modules/news/infrastructure/rdb/models"
	"nfxnews/pkgs/cachex"
	"nfxnews/pkgs/errx"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	db    *gorm.DB
	cache *cachex.Connection
}

func NewService(db *gorm.DB, cache *cachex.Connection) *Service {
	return &Service{db: db, cache: cache}
}

type ItemView struct {
	ID         string         `json:"id"`
	SourceID   string         `json:"source_id"`
	OriginalID string         `json:"original_id"`
	Title      string         `json:"title"`
	URL        string         `json:"url"`
	MobileURL  string         `json:"mobile_url,omitempty"`
	PubDate    int64          `json:"pub_date,omitempty"`
	Extra      map[string]any `json:"extra,omitempty"`
}

func (s *Service) UpsertFromEvent(ctx context.Context, ev events.SourceFetchedEvent) (int, error) {
	count := 0
	for _, it := range ev.Items {
		original := it.ID
		id := ev.SourceID + ":" + original
		row := models.Item{
			ID: id, SourceID: ev.SourceID, OriginalID: original, Title: it.Title, URL: it.URL, Extra: []byte(it.ExtraJSON),
		}
		if it.MobileURL != "" {
			row.MobileURL = &it.MobileURL
		}
		if it.PubDate > 0 {
			t := time.Unix(it.PubDate, 0)
			row.PubDate = &t
		}
		if len(row.Extra) == 0 {
			row.Extra = []byte("{}")
		}
		if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&row).Error; err != nil {
			return count, errx.ErrInternal.WithCause(err)
		}
		count++
	}
	if s.cache != nil && s.cache.Client() != nil && ev.SourceID != "" {
		_ = s.cache.Client().Del(ctx, "news:source:"+ev.SourceID).Err()
		items, _ := s.ListBySource(ctx, ev.SourceID, 50)
		raw, _ := json.Marshal(items)
		_ = s.cache.Client().Set(ctx, "news:source:"+ev.SourceID, raw, 10*time.Minute).Err()
	}
	return count, nil
}

func (s *Service) ListBySource(ctx context.Context, sourceID string, limit int) ([]ItemView, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if sourceID != "" && s.cache != nil && s.cache.Client() != nil {
		raw, err := s.cache.Client().Get(ctx, "news:source:"+sourceID).Bytes()
		if err == nil && len(raw) > 0 {
			var cached []ItemView
			if json.Unmarshal(raw, &cached) == nil {
				return cached, nil
			}
		}
	}
	var rows []models.Item
	q := s.db.WithContext(ctx).Order("updated_at desc").Limit(limit)
	if sourceID != "" {
		q = q.Where("source_id = ?", sourceID)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return toViews(rows), nil
}

func (s *Service) Search(ctx context.Context, query string, limit int) ([]ItemView, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query = strings.TrimSpace(query)
	var rows []models.Item
	tx := s.db.WithContext(ctx).Order("updated_at desc").Limit(limit)
	if query != "" {
		tx = tx.Where("title ILIKE ?", "%"+query+"%")
	}
	if err := tx.Find(&rows).Error; err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return toViews(rows), nil
}

func toViews(rows []models.Item) []ItemView {
	out := make([]ItemView, 0, len(rows))
	for _, r := range rows {
		v := ItemView{ID: r.ID, SourceID: r.SourceID, OriginalID: r.OriginalID, Title: r.Title, URL: r.URL}
		if r.MobileURL != nil {
			v.MobileURL = *r.MobileURL
		}
		if r.PubDate != nil {
			v.PubDate = r.PubDate.Unix()
		}
		if len(r.Extra) > 0 {
			_ = json.Unmarshal(r.Extra, &v.Extra)
		}
		out = append(out, v)
	}
	return out
}

type Preference struct {
	AccountID   string          `gorm:"column:account_id;type:uuid;primaryKey"`
	ProfileID   string          `gorm:"column:profile_id;type:uuid;primaryKey"`
	ColumnOrder json.RawMessage `gorm:"column:column_order;type:jsonb"`
	Payload     json.RawMessage `gorm:"column:payload;type:jsonb"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (Preference) TableName() string { return "news.profile_preferences" }

type PreferenceView struct {
	ColumnOrder json.RawMessage `json:"column_order"`
	Payload     json.RawMessage `json:"payload"`
}

func (s *Service) GetPreferences(ctx context.Context, accountID, profileID string) (PreferenceView, error) {
	var pref Preference
	err := s.db.WithContext(ctx).First(&pref, "account_id = ? AND profile_id = ?", accountID, profileID).Error
	if err == gorm.ErrRecordNotFound {
		return PreferenceView{ColumnOrder: json.RawMessage("[]"), Payload: json.RawMessage("{}")}, nil
	}
	if err != nil {
		return PreferenceView{}, err
	}
	return PreferenceView{ColumnOrder: pref.ColumnOrder, Payload: pref.Payload}, nil
}

func (s *Service) SetPreferences(ctx context.Context, accountID, profileID string, columnOrder, payload json.RawMessage) error {
	pref := Preference{
		AccountID: accountID, ProfileID: profileID,
		ColumnOrder: columnOrder, Payload: payload, UpdatedAt: time.Now(),
	}
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "profile_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"column_order", "payload", "updated_at"}),
	}).Create(&pref).Error
}
