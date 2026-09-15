package sourceapp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"nfxnews/errors/src/common"
	"nfxnews/events"
	"nfxnews/modules/source/infrastructure/getters"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/kafkax/eventbus"
)

type Service struct {
	reg *getters.Registry
	pub *eventbus.BusPublisher
}

func NewService(reg *getters.Registry, pub *eventbus.BusPublisher) *Service {
	return &Service{reg: reg, pub: pub}
}

func NewRegistryFromWD(wd string) (*getters.Registry, error) {
	reg := getters.NewRegistry()
	path := filepath.Join(wd, "modules/source/infrastructure/getters/catalog.json")
	if err := reg.LoadCatalog(path); err != nil {
		return reg, err
	}
	return reg, nil
}

func (s *Service) List() []getters.Meta { return s.reg.List() }

func (s *Service) Fetch(ctx context.Context, sourceID string) ([]getters.Item, error) {
	items, err := s.reg.Fetch(ctx, sourceID)
	if err != nil {
		return nil, common.ErrSourceNotFound.WithCause(err).WithMsg(err.Error())
	}
	if s.pub != nil {
		evItems := make([]events.SourceNewsItem, 0, len(items))
		for _, it := range items {
			extra, _ := json.Marshal(it.Extra)
			evItems = append(evItems, events.SourceNewsItem{
				ID: it.ID, Title: it.Title, URL: it.URL, MobileURL: it.MobileURL, PubDate: it.PubDate, ExtraJSON: string(extra),
			})
		}
		_ = eventbus.PublishEvent(ctx, s.pub, events.SourceFetchedEvent{SourceID: sourceID, Items: evItems})
	}
	return items, nil
}

func (s *Service) Meta(id string) (getters.Meta, error) {
	m, ok := s.reg.Meta(id)
	if !ok {
		return getters.Meta{}, common.ErrSourceNotFound
	}
	return m, nil
}

func RequireID(id string) error {
	if id == "" {
		return errx.ErrInvalidParams.WithMsg("source_id is required")
	}
	return nil
}
