package reportapp

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"nfxnews/events"
	keywordDomain "nfxnews/modules/report/domain/keyword"
	snapshotDomain "nfxnews/modules/report/domain/snapshot"
	repofactory "nfxnews/modules/report/infrastructure/repository/factory"
	reportQuery "nfxnews/modules/report/query/report"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/kafkax/eventbus"
	"nfxnews/pkgs/transaction"
	newspb "nfxnews/protos/gen/news"

	"github.com/google/uuid"
)

type Keyword = reportQuery.KeywordVO
type Snapshot = reportQuery.SnapshotVO

type Service struct {
	tx    transaction.TxManager
	repos *repofactory.TxRepoFactory
	query *reportQuery.Query
	news  newspb.NewsServiceClient
	pub   *eventbus.BusPublisher
}

func NewService(
	tx transaction.TxManager,
	repos *repofactory.TxRepoFactory,
	query *reportQuery.Query,
	news newspb.NewsServiceClient,
	pub *eventbus.BusPublisher,
) *Service {
	return &Service{tx: tx, repos: repos, query: query, news: news, pub: pub}
}

func (s *Service) ListKeywords(ctx context.Context) ([]Keyword, error) {
	rows, err := s.query.Keywords.All(ctx)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) AddKeyword(ctx context.Context, group, word, kind string, limit int) (*Keyword, error) {
	word = strings.TrimSpace(word)
	if word == "" {
		return nil, errx.ErrInvalidBody.WithMsg("word required")
	}
	if kind == "" {
		kind = "include"
	}
	if group == "" {
		group = "default"
	}
	now := time.Now()
	id := uuid.Must(uuid.NewV7())
	ent := keywordDomain.NewFromState(keywordDomain.State{
		ID: id, GroupName: group, Word: word, Kind: kind, CountLimit: limit, CreatedAt: now,
	})
	if err := s.repos.Keyword(transaction.UoW{}).Create.New(ctx, ent); err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return &Keyword{ID: id, GroupName: group, Word: word, Kind: kind, CountLimit: limit, CreatedAt: now}, nil
}

func (s *Service) Generate(ctx context.Context, mode string) (*Snapshot, error) {
	if mode == "" {
		mode = "daily"
	}
	keywords, err := s.ListKeywords(ctx)
	if err != nil {
		return nil, err
	}
	matched := []map[string]any{}
	if s.news != nil {
		resp, err := s.news.SearchNews(ctx, &newspb.SearchNewsRequest{Query: "", Limit: 200})
		if err == nil {
			for _, item := range resp.GetItems() {
				if matchKeywords(item.GetTitle(), keywords) {
					matched = append(matched, map[string]any{
						"id": item.GetId(), "title": item.GetTitle(), "url": item.GetUrl(), "source_id": item.GetSourceId(),
					})
				}
			}
		}
	}
	payload, _ := json.Marshal(map[string]any{"mode": mode, "items": matched})
	now := time.Now()
	id := uuid.Must(uuid.NewV7())
	title := strings.ToUpper(mode) + " report"
	ent := snapshotDomain.NewFromState(snapshotDomain.State{
		ID: id, Mode: mode, Title: title, Payload: payload, ItemCount: len(matched), CreatedAt: now,
	})
	if err := s.repos.Snapshot(transaction.UoW{}).Create.New(ctx, ent); err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	var payloadObj any
	_ = json.Unmarshal(payload, &payloadObj)
	snap := Snapshot{ID: id, Mode: mode, Title: title, Payload: payload, PayloadObj: payloadObj, ItemCount: len(matched), CreatedAt: now}
	if s.pub != nil {
		_ = eventbus.PublishEvent(ctx, s.pub, events.ReportGeneratedEvent{
			ReportID: id.String(), Mode: mode, Title: title, ItemCount: len(matched), Payload: string(payload),
		})
	}
	return &snap, nil
}

func (s *Service) ListSnapshots(ctx context.Context, limit int) ([]Snapshot, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.query.Snapshots.Recent(ctx, limit)
}

func (s *Service) Get(ctx context.Context, id string) (*Snapshot, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errx.ErrInvalidParams.WithCause(err)
	}
	return s.query.Snapshots.ByID(ctx, uid)
}

func matchKeywords(title string, keywords []Keyword) bool {
	if len(keywords) == 0 {
		return true
	}
	t := strings.ToLower(title)
	excluded := false
	included := false
	hasInclude := false
	for _, k := range keywords {
		w := strings.ToLower(strings.TrimPrefix(k.Word, "+"))
		w = strings.TrimPrefix(w, "!")
		if k.Kind == "exclude" || strings.HasPrefix(k.Word, "!") {
			if strings.Contains(t, w) {
				excluded = true
			}
			continue
		}
		hasInclude = true
		if strings.Contains(t, w) {
			included = true
		}
	}
	if excluded {
		return false
	}
	if !hasInclude {
		return true
	}
	return included
}
