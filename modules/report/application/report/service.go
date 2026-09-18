package reportapp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	reporterr "nfxnews/errors/src/report"
	"nfxnews/events"
	"nfxnews/modules/report/application/frequency"
	"nfxnews/modules/report/application/htmlreport"
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
	tx          transaction.TxManager
	repoFactory *repofactory.TxRepoFactory
	query       *reportQuery.Query
	news        newspb.NewsServiceClient
	pub         *eventbus.BusPublisher
}

func NewService(
	tx transaction.TxManager,
	repoFactory *repofactory.TxRepoFactory,
	query *reportQuery.Query,
	news newspb.NewsServiceClient,
	pub *eventbus.BusPublisher,
) *Service {
	return &Service{tx: tx, repoFactory: repoFactory, query: query, news: news, pub: pub}
}

func (s *Service) ListKeywords(ctx context.Context, accountID string) ([]Keyword, error) {
	rows, err := s.query.Keywords.All(ctx, accountID)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) AddKeyword(ctx context.Context, accountID, profileID, group, word, kind string, limit int) (*Keyword, error) {
	if accountID == "" {
		return nil, errx.Unauthorized("INVALID_TOKEN", "missing account")
	}
	var aid, pid *string
	if accountID != "" {
		aid = &accountID
	}
	if profileID != "" {
		pid = &profileID
	}
	ent, err := keywordDomain.New(aid, pid, group, word, kind, limit)
	if err != nil {
		return nil, err
	}
	keywordRepo := s.repoFactory.Keyword(transaction.UoW{})
	if err := keywordRepo.Create.New(ctx, ent); err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	st := ent.State()
	return &Keyword{ID: st.ID, GroupName: st.GroupName, Word: st.Word, Kind: st.Kind, CountLimit: st.CountLimit, CreatedAt: st.CreatedAt}, nil
}

func (s *Service) Generate(ctx context.Context, accountID, profileID, mode string) (*Snapshot, error) {
	mode, err := snapshotDomain.NormalizeMode(mode)
	if err != nil {
		return nil, err
	}
	dict, err := s.loadDictionary(ctx, accountID)
	if err != nil {
		return nil, err
	}
	limit := int32(500)
	if mode == "current" {
		limit = 80
	}
	var items []*newspb.NewsItem
	if s.news != nil {
		resp, err := s.news.SearchNews(ctx, &newspb.SearchNewsRequest{Query: "", Limit: limit})
		if err != nil {
			return nil, reporterr.ErrReportGenerateFailed.WithCause(err)
		}
		items = resp.GetItems()
	}
	if mode == "daily" {
		start := time.Now().UTC().Truncate(24 * time.Hour).Unix()
		filtered := items[:0]
		for _, it := range items {
			ts := pubUnix(it)
			if ts == 0 || ts >= start {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}
	prevTitles := map[string]struct{}{}
	if snaps, err := s.query.Snapshots.Recent(ctx, accountID, 20); err == nil {
		for _, snap := range snaps {
			for _, t := range titlesFromPayload(snap.PayloadObj) {
				prevTitles[t] = struct{}{}
			}
			break
		}
	}
	type row struct {
		ID, Title, URL, SourceID, Group string
		IsNew                           bool
	}
	matched := []row{}
	groupCount := map[string]int{}
	for _, item := range items {
		g, ok := frequency.MatchTitle(item.GetTitle(), dict)
		if !ok {
			continue
		}
		if g.MaxCount > 0 && groupCount[g.GroupKey] >= g.MaxCount {
			continue
		}
		_, seen := prevTitles[item.GetTitle()]
		if mode == "incremental" && seen {
			continue
		}
		groupCount[g.GroupKey]++
		matched = append(matched, row{
			ID: item.GetId(), Title: item.GetTitle(), URL: item.GetUrl(), SourceID: item.GetSourceId(),
			Group: g.GroupKey, IsNew: !seen,
		})
	}
	htmlItems := make([]htmlreport.Item, 0, len(matched))
	payloadItems := make([]map[string]any, 0, len(matched))
	for _, m := range matched {
		htmlItems = append(htmlItems, htmlreport.Item{Title: m.Title, URL: m.URL, SourceID: m.SourceID, Group: m.Group, IsNew: m.IsNew})
		payloadItems = append(payloadItems, map[string]any{
			"id": m.ID, "title": m.Title, "url": m.URL, "source_id": m.SourceID, "group": m.Group, "is_new": m.IsNew,
		})
	}
	title := strings.ToUpper(mode) + " report"
	now := time.Now().UTC()
	html := htmlreport.Render(title, mode, now, htmlItems)
	payload, _ := json.Marshal(map[string]any{"mode": mode, "items": payloadItems, "html": html})
	var said, spid *string
	if accountID != "" {
		said = &accountID
	}
	if profileID != "" {
		spid = &profileID
	}
	ent, err := snapshotDomain.New(said, spid, mode, title, payload, len(matched))
	if err != nil {
		return nil, err
	}
	snapshotRepo := s.repoFactory.Snapshot(transaction.UoW{})
	if err := snapshotRepo.Create.New(ctx, ent); err != nil {
		return nil, reporterr.ErrReportGenerateFailed.WithCause(err)
	}
	var payloadObj any
	_ = json.Unmarshal(payload, &payloadObj)
	st := ent.State()
	snap := Snapshot{ID: st.ID, Mode: st.Mode, Title: st.Title, Payload: payload, PayloadObj: payloadObj, ItemCount: st.ItemCount, CreatedAt: st.CreatedAt}
	if s.pub != nil {
		_ = eventbus.PublishEvent(ctx, s.pub, events.ReportGeneratedEvent{
			ReportID: st.ID.String(), AccountID: accountID, ProfileID: profileID, Mode: mode, Title: title, ItemCount: len(matched), Payload: string(payload),
		})
	}
	return &snap, nil
}

func (s *Service) ListSnapshots(ctx context.Context, accountID string, limit int) ([]Snapshot, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.query.Snapshots.Recent(ctx, accountID, limit)
}

func (s *Service) Get(ctx context.Context, accountID, id string) (*Snapshot, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errx.ErrInvalidParams.WithCause(err)
	}
	row, err := s.query.Snapshots.ByID(ctx, accountID, uid)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, reporterr.ErrReportNotFound
	}
	return row, nil
}

func (s *Service) HTML(ctx context.Context, accountID, id string) (string, error) {
	row, err := s.Get(ctx, accountID, id)
	if err != nil {
		return "", err
	}
	if m, ok := row.PayloadObj.(map[string]any); ok {
		if h, ok := m["html"].(string); ok && h != "" {
			return h, nil
		}
	}
	return htmlreport.Render(row.Title, row.Mode, row.CreatedAt, nil), nil
}

func (s *Service) loadDictionary(ctx context.Context, accountID string) (frequency.Dictionary, error) {
	path := os.Getenv("FREQUENCY_WORDS_PATH")
	if path == "" {
		path = filepath.Join("config", "frequency_words.txt")
	}
	dict, err := frequency.LoadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return frequency.Dictionary{}, reporterr.ErrFrequencyFileMissing.WithCause(err)
	}
	keywords, err := s.ListKeywords(ctx, accountID)
	if err != nil {
		return dict, err
	}
	if len(keywords) == 0 {
		return dict, nil
	}
	byGroup := map[string]*frequency.Group{}
	for _, k := range keywords {
		g := byGroup[k.GroupName]
		if g == nil {
			g = &frequency.Group{Name: k.GroupName, MaxCount: k.CountLimit}
			byGroup[k.GroupName] = g
		}
		switch k.Kind {
		case "exclude":
			dict.Exclude = append(dict.Exclude, k.Word)
		case "required":
			g.Required = append(g.Required, k.Word)
		default:
			g.Normal = append(g.Normal, k.Word)
		}
	}
	for _, g := range byGroup {
		if len(g.Required) == 0 && len(g.Normal) == 0 {
			continue
		}
		g.GroupKey = strings.Join(g.Normal, " ")
		if g.GroupKey == "" {
			g.GroupKey = strings.Join(g.Required, " ")
		}
		dict.Groups = append(dict.Groups, *g)
	}
	return dict, nil
}

func pubUnix(it *newspb.NewsItem) int64 {
	v := it.GetPubDateUnix()
	if v > 1_000_000_000_000 {
		return v / 1000
	}
	return v
}

func titlesFromPayload(obj any) []string {
	m, ok := obj.(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := m["items"].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, it := range raw {
		row, ok := it.(map[string]any)
		if !ok {
			continue
		}
		if t, ok := row["title"].(string); ok {
			out = append(out, t)
		}
	}
	return out
}
