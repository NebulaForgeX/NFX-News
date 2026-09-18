package crawlapp

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	crawlerr "nfxnews/errors/src/crawl"
	sessionDomain "nfxnews/modules/crawl/domain/session"
	repofactory "nfxnews/modules/crawl/infrastructure/repository/factory"
	sessionQuery "nfxnews/modules/crawl/query/session"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/logx"
	"nfxnews/pkgs/transaction"
	reportpb "nfxnews/protos/gen/report"
	sourcepb "nfxnews/protos/gen/source"

	"github.com/google/uuid"
)

type Session = sessionQuery.SessionVO

type Service struct {
	tx          transaction.TxManager
	repoFactory *repofactory.TxRepoFactory
	query       *sessionQuery.Query
	source      sourcepb.SourceServiceClient
	report      reportpb.ReportServiceClient
	mu          sync.Mutex
	lastFetch   map[string]time.Time
	lastTitles  map[string]map[string]struct{}
}

func NewService(
	tx transaction.TxManager,
	repoFactory *repofactory.TxRepoFactory,
	query *sessionQuery.Query,
	source sourcepb.SourceServiceClient,
	report reportpb.ReportServiceClient,
) *Service {
	return &Service{
		tx: tx, repoFactory: repoFactory, query: query, source: source, report: report,
		lastFetch:  map[string]time.Time{},
		lastTitles: map[string]map[string]struct{}{},
	}
}

func ownerPtrs(accountID, profileID string) (aid, pid *string) {
	if accountID != "" {
		aid = &accountID
	}
	if profileID != "" {
		pid = &profileID
	}
	return aid, pid
}

func (s *Service) Trigger(ctx context.Context, accountID, profileID, sourceID string) (*Session, error) {
	if sourceID == "" {
		return s.TriggerAll(ctx, accountID, profileID)
	}
	aid, pid := ownerPtrs(accountID, profileID)
	sid := sourceID
	sess := sessionDomain.New(aid, pid, &sid)
	if err := s.persistNew(ctx, sess); err != nil {
		return nil, err
	}
	count, newTitles, fetchErr := s.fetchOne(ctx, sourceID)
	finished := time.Now().UTC()
	if fetchErr != nil {
		msg := fetchErr.Error()
		sess.Finish("failed", count, &msg, finished)
	} else {
		sess.Finish("ok", count, nil, finished)
	}
	_ = s.persistUpdate(ctx, sess)
	if len(newTitles) > 0 {
		s.emitIncremental(ctx)
	}
	return sessionVO(sess), nil
}

func (s *Service) TriggerAll(ctx context.Context, accountID, profileID string) (*Session, error) {
	aid, pid := ownerPtrs(accountID, profileID)
	sess := sessionDomain.New(aid, pid, nil)
	if err := s.persistNew(ctx, sess); err != nil {
		return nil, err
	}
	if s.source == nil {
		msg := "source client unavailable"
		sess.Finish("failed", 0, &msg, time.Now().UTC())
		_ = s.persistUpdate(ctx, sess)
		return sessionVO(sess), crawlerr.ErrCrawlTriggerFailed
	}
	list, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
	if err != nil {
		msg := err.Error()
		sess.Finish("failed", 0, &msg, time.Now().UTC())
		_ = s.persistUpdate(ctx, sess)
		return sessionVO(sess), crawlerr.ErrCrawlTriggerFailed.WithCause(err)
	}
	total := 0
	anyNew := false
	var firstErr error
	for _, src := range list.GetSources() {
		if src.GetRedirect() != "" {
			continue
		}
		n, news, ferr := s.fetchOne(ctx, src.GetId())
		total += n
		if len(news) > 0 {
			anyNew = true
		}
		if ferr != nil && firstErr == nil {
			firstErr = ferr
		}
	}
	finished := time.Now().UTC()
	if firstErr != nil && total == 0 {
		msg := firstErr.Error()
		sess.Finish("failed", total, &msg, finished)
	} else {
		sess.Finish("ok", total, nil, finished)
	}
	_ = s.persistUpdate(ctx, sess)
	if anyNew {
		s.emitIncremental(ctx)
	}
	return sessionVO(sess), nil
}

func (s *Service) FetchDue(ctx context.Context) error {
	if s.source == nil {
		return nil
	}
	list, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
	if err != nil {
		return err
	}
	now := time.Now()
	jitterPct := jitterPercent()
	anyNew := false
	for _, src := range list.GetSources() {
		if src.GetRedirect() != "" {
			continue
		}
		interval := time.Duration(src.GetIntervalMs()) * time.Millisecond
		if interval <= 0 {
			interval = 10 * time.Minute
		}
		s.mu.Lock()
		last := s.lastFetch[src.GetId()]
		s.mu.Unlock()
		wait := withJitter(interval, jitterPct)
		if !last.IsZero() && now.Before(last.Add(wait)) {
			continue
		}
		_, news, ferr := s.fetchOne(ctx, src.GetId())
		if ferr != nil {
			logx.S().Warnf("scheduled crawl %s: %v", src.GetId(), ferr)
			continue
		}
		if len(news) > 0 {
			anyNew = true
		}
	}
	if anyNew {
		s.emitIncremental(ctx)
	}
	return nil
}

func (s *Service) FetchCatalog(ctx context.Context) error {
	return s.FetchDue(ctx)
}

func (s *Service) List(ctx context.Context, accountID string, limit int) ([]Session, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.query.List.Recent(ctx, accountID, limit)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) Get(ctx context.Context, accountID, id string) (*Session, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errx.ErrInvalidParams.WithCause(err)
	}
	row, err := s.query.List.ByID(ctx, accountID, uid)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, crawlerr.ErrCrawlSessionNotFound
	}
	return row, nil
}

func (s *Service) fetchOne(ctx context.Context, sourceID string) (int, []string, error) {
	if s.source == nil {
		return 0, nil, crawlerr.ErrCrawlTriggerFailed
	}
	resp, err := s.source.FetchSource(ctx, &sourcepb.FetchSourceRequest{SourceId: sourceID})
	s.mu.Lock()
	s.lastFetch[sourceID] = time.Now()
	s.mu.Unlock()
	if err != nil {
		return 0, nil, err
	}
	titles := make([]string, 0, len(resp.GetItems()))
	for _, it := range resp.GetItems() {
		if it.GetTitle() != "" {
			titles = append(titles, it.GetTitle())
		}
	}
	return len(resp.GetItems()), s.detectNewTitles(sourceID, titles), nil
}

func (s *Service) detectNewTitles(sourceID string, titles []string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	prev := s.lastTitles[sourceID]
	next := make(map[string]struct{}, len(titles))
	var news []string
	for _, t := range titles {
		next[t] = struct{}{}
		if prev != nil {
			if _, seen := prev[t]; seen {
				continue
			}
		}
		news = append(news, t)
	}
	if prev == nil {
		news = nil
	}
	s.lastTitles[sourceID] = next
	return news
}

func (s *Service) emitIncremental(ctx context.Context) {
	if s.report == nil {
		return
	}
	if _, err := s.report.GenerateReport(ctx, &reportpb.GenerateReportRequest{Mode: "incremental"}); err != nil {
		logx.S().Warnf("incremental report after new titles: %v", err)
	}
}

func (s *Service) persistNew(ctx context.Context, sess *sessionDomain.Session) error {
	if err := s.tx.WithUoW(ctx, func(ctx context.Context, uow transaction.UoW) error {
		return s.repoFactory.Session(uow).Create.New(ctx, sess)
	}); err != nil {
		return errx.ErrInternal.WithCause(err)
	}
	return nil
}

func (s *Service) persistUpdate(ctx context.Context, sess *sessionDomain.Session) error {
	return s.tx.WithUoW(ctx, func(ctx context.Context, uow transaction.UoW) error {
		return s.repoFactory.Session(uow).Update.Generic(ctx, sess)
	})
}

func sessionVO(sess *sessionDomain.Session) *Session {
	st := sess.State()
	vo := Session{
		ID: st.ID, AccountID: st.AccountID, ProfileID: st.ProfileID, SourceID: st.SourceID,
		Status: st.Status, ItemCount: st.ItemCount, ErrorMessage: st.ErrorMessage, StartedAt: st.StartedAt, FinishedAt: st.FinishedAt,
	}
	return &vo
}

func jitterPercent() int {
	if v := os.Getenv("CRAWL_INTERVAL_JITTER_PERCENT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n <= 90 {
			return n
		}
	}
	return 15
}

func withJitter(base time.Duration, pct int) time.Duration {
	if pct <= 0 {
		return base
	}
	span := int64(base) * int64(pct) / 100
	if span <= 0 {
		return base
	}
	delta := rand.Int63n(span*2+1) - span
	out := base + time.Duration(delta)
	if out < time.Second {
		return time.Second
	}
	return out
}
