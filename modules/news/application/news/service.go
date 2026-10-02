package newsapp

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"nfxnews/events"
	itemDomain "nfxnews/modules/news/domain/item"
	prefDomain "nfxnews/modules/news/domain/preference"
	repofactory "nfxnews/modules/news/infrastructure/repository/factory"
	itemQuery "nfxnews/modules/news/query/item"
	"nfxnews/pkgs/cachex"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/transaction"
)

type Service struct {
	tx          transaction.TxManager
	repoFactory *repofactory.TxRepoFactory
	query       *itemQuery.Query
	cache       *cachex.Connection
}

func NewService(
	tx transaction.TxManager,
	repoFactory *repofactory.TxRepoFactory,
	query *itemQuery.Query,
	cache *cachex.Connection,
) *Service {
	return &Service{tx: tx, repoFactory: repoFactory, query: query, cache: cache}
}

type ItemView = itemQuery.ItemVO
type PreferenceView = itemQuery.PreferenceVO

func (s *Service) UpsertFromEvent(ctx context.Context, ev events.SourceFetchedEvent) (int, error) {
	count := 0
	err := s.tx.WithUoW(ctx, func(ctx context.Context, uow transaction.UoW) error {
		itemRepo := s.repoFactory.Item(uow)
		for _, it := range ev.Items {
			original := it.ID
			id := ev.SourceID + ":" + original
			extra := []byte(it.ExtraJSON)
			if len(extra) == 0 {
				extra = []byte("{}")
			}
			st := itemDomain.State{ID: id, SourceID: ev.SourceID, OriginalID: original, Title: it.Title, URL: it.URL, Extra: extra}
			if it.MobileURL != "" {
				st.MobileURL = &it.MobileURL
			}
			if it.PubDate > 0 {
				t := time.UnixMilli(it.PubDate)
				st.PubDate = &t
			}
			if err := itemRepo.Create.Upsert(ctx, itemDomain.NewFromState(st)); err != nil {
				return err
			}
			count++
		}
		return nil
	})
	if err != nil {
		return count, errx.ErrInternal.WithCause(err)
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
	rows, err := s.query.List.BySource(ctx, sourceID, limit)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) Search(ctx context.Context, query string, limit int) ([]ItemView, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	rows, err := s.query.List.Search(ctx, strings.TrimSpace(query), limit)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) GetPreferences(ctx context.Context, accountID, profileID string) (PreferenceView, error) {
	pref, err := s.query.Pref.ByAccountProfile(ctx, accountID, profileID)
	if err != nil {
		return PreferenceView{}, err
	}
	if pref == nil {
		return PreferenceView{ColumnOrder: json.RawMessage("[]"), Payload: json.RawMessage("{}")}, nil
	}
	return *pref, nil
}

func (s *Service) SetPreferences(ctx context.Context, accountID, profileID string, columnOrder, payload json.RawMessage) error {
	preferenceRepo := s.repoFactory.Preference(transaction.UoW{})
	return preferenceRepo.Create.Upsert(ctx, prefDomain.NewFromState(prefDomain.State{
		AccountID: accountID, ProfileID: profileID, ColumnOrder: columnOrder, Payload: payload, UpdatedAt: time.Now(),
	}))
}
