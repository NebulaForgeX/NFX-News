package handler

import (
	"context"
	"encoding/json"

	sourceapp "nfxnews/modules/source/application/source"
	sourcepb "nfxnews/protos/gen/source"
)

type SourceHandler struct {
	sourcepb.UnimplementedSourceServiceServer
	svc *sourceapp.Service
}

func NewSourceHandler(svc *sourceapp.Service) *SourceHandler { return &SourceHandler{svc: svc} }

func (h *SourceHandler) FetchSource(ctx context.Context, req *sourcepb.FetchSourceRequest) (*sourcepb.FetchSourceResponse, error) {
	items, err := h.svc.Fetch(ctx, req.GetSourceId())
	if err != nil {
		return nil, err
	}
	out := make([]*sourcepb.NewsItem, 0, len(items))
	for _, it := range items {
		extra, _ := json.Marshal(it.Extra)
		out = append(out, &sourcepb.NewsItem{Id: it.ID, Title: it.Title, Url: it.URL, MobileUrl: it.MobileURL, PubDateUnix: it.PubDate, ExtraJson: string(extra)})
	}
	return &sourcepb.FetchSourceResponse{SourceId: req.GetSourceId(), Items: out}, nil
}

func (h *SourceHandler) ListSources(ctx context.Context, req *sourcepb.ListSourcesRequest) (*sourcepb.ListSourcesResponse, error) {
	list := h.svc.List()
	out := make([]*sourcepb.SourceMeta, 0, len(list))
	for _, m := range list {
		out = append(out, &sourcepb.SourceMeta{Id: m.ID, Name: m.Name, Title: m.Title, Column: m.Column, Home: m.Home, Color: m.Color, IntervalMs: int32(m.IntervalMS), Type: m.Type, Redirect: m.Redirect})
	}
	return &sourcepb.ListSourcesResponse{Sources: out}, nil
}

func (h *SourceHandler) GetSourceMeta(ctx context.Context, req *sourcepb.GetSourceMetaRequest) (*sourcepb.GetSourceMetaResponse, error) {
	m, err := h.svc.Meta(req.GetSourceId())
	if err != nil {
		return nil, err
	}
	return &sourcepb.GetSourceMetaResponse{Source: &sourcepb.SourceMeta{Id: m.ID, Name: m.Name, Title: m.Title, Column: m.Column, Home: m.Home, Color: m.Color, IntervalMs: int32(m.IntervalMS), Type: m.Type, Redirect: m.Redirect}}, nil
}
