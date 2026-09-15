package handler

import (
	"context"
	"encoding/json"

	"nfxnews/events"
	newsapp "nfxnews/modules/news/application/news"
	newspb "nfxnews/protos/gen/news"
)

type NewsHandler struct {
	newspb.UnimplementedNewsServiceServer
	svc *newsapp.Service
}

func NewNewsHandler(svc *newsapp.Service) *NewsHandler { return &NewsHandler{svc: svc} }

func toPB(items []newsapp.ItemView) []*newspb.NewsItem {
	out := make([]*newspb.NewsItem, 0, len(items))
	for _, it := range items {
		extra, _ := json.Marshal(it.Extra)
		out = append(out, &newspb.NewsItem{
			Id: it.ID, SourceId: it.SourceID, OriginalId: it.OriginalID, Title: it.Title,
			Url: it.URL, MobileUrl: it.MobileURL, PubDateUnix: it.PubDate, ExtraJson: string(extra),
		})
	}
	return out
}

func (h *NewsHandler) UpsertItems(ctx context.Context, req *newspb.UpsertItemsRequest) (*newspb.UpsertItemsResponse, error) {
	ev := events.SourceFetchedEvent{SourceID: req.GetSourceId()}
	for _, it := range req.GetItems() {
		ev.Items = append(ev.Items, events.SourceNewsItem{
			ID: it.GetOriginalId(), Title: it.GetTitle(), URL: it.GetUrl(), MobileURL: it.GetMobileUrl(),
			PubDate: it.GetPubDateUnix(), ExtraJSON: it.GetExtraJson(),
		})
	}
	n, err := h.svc.UpsertFromEvent(ctx, ev)
	if err != nil {
		return nil, err
	}
	return &newspb.UpsertItemsResponse{Upserted: int32(n)}, nil
}

func (h *NewsHandler) ListBySource(ctx context.Context, req *newspb.ListBySourceRequest) (*newspb.ListBySourceResponse, error) {
	items, err := h.svc.ListBySource(ctx, req.GetSourceId(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &newspb.ListBySourceResponse{Items: toPB(items)}, nil
}

func (h *NewsHandler) SearchNews(ctx context.Context, req *newspb.SearchNewsRequest) (*newspb.SearchNewsResponse, error) {
	items, err := h.svc.Search(ctx, req.GetQuery(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &newspb.SearchNewsResponse{Items: toPB(items)}, nil
}
