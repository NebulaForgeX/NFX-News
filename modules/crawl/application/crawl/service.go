package crawlapp

import (
	"context"
	"time"

	sessionDomain "nfxnews/modules/crawl/domain/session"
	repofactory "nfxnews/modules/crawl/infrastructure/repository/factory"
	sessionQuery "nfxnews/modules/crawl/query/session"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/transaction"
	sourcepb "nfxnews/protos/gen/source"

	"github.com/google/uuid"
)

type Session = sessionQuery.SessionVO

type Service struct {
	tx     transaction.TxManager
	repos  *repofactory.TxRepoFactory
	query  *sessionQuery.Query
	source sourcepb.SourceServiceClient
}

func NewService(
	tx transaction.TxManager,
	repos *repofactory.TxRepoFactory,
	query *sessionQuery.Query,
	source sourcepb.SourceServiceClient,
) *Service {
	return &Service{tx: tx, repos: repos, query: query, source: source}
}

func (s *Service) Trigger(ctx context.Context, sourceID string) (*Session, error) {
	now := time.Now()
	sess := sessionDomain.NewFromState(sessionDomain.State{
		ID: uuid.Must(uuid.NewV7()), Status: "running", StartedAt: now,
	})
	if sourceID != "" {
		id := sourceID
		st := sess.State()
		st.SourceID = &id
		sess = sessionDomain.NewFromState(st)
	}
	if err := s.tx.WithUoW(ctx, func(ctx context.Context, uow transaction.UoW) error {
		return s.repos.Session(uow).Create.New(ctx, sess)
	}); err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	if s.source != nil && sourceID != "" {
		resp, err := s.source.FetchSource(ctx, &sourcepb.FetchSourceRequest{SourceId: sourceID})
		finished := time.Now()
		if err != nil {
			msg := err.Error()
			sess.Finish("failed", 0, &msg, finished)
		} else {
			sess.Finish("ok", len(resp.GetItems()), nil, finished)
		}
		_ = s.tx.WithUoW(ctx, func(ctx context.Context, uow transaction.UoW) error {
			return s.repos.Session(uow).Update.Generic(ctx, sess)
		})
	}
	st := sess.State()
	vo := Session{ID: st.ID, SourceID: st.SourceID, Status: st.Status, ItemCount: st.ItemCount, ErrorMessage: st.ErrorMessage, StartedAt: st.StartedAt, FinishedAt: st.FinishedAt}
	return &vo, nil
}

func (s *Service) TriggerAll(ctx context.Context) (*Session, error) {
	if s.source == nil {
		return s.Trigger(ctx, "")
	}
	list, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
	if err != nil {
		return s.Trigger(ctx, "")
	}
	var last *Session
	for _, src := range list.GetSources() {
		last, err = s.Trigger(ctx, src.GetId())
		if err != nil {
			return last, err
		}
	}
	if last == nil {
		return s.Trigger(ctx, "")
	}
	return last, nil
}

func (s *Service) List(ctx context.Context, limit int) ([]Session, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.query.List.Recent(ctx, limit)
	if err != nil {
		return nil, errx.ErrInternal.WithCause(err)
	}
	return rows, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Session, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errx.ErrInvalidParams.WithCause(err)
	}
	row, err := s.query.List.ByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return row, nil
}
