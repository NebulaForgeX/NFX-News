package handler

import (
	"context"

	crawlapp "nfxnews/modules/crawl/application/crawl"
	crawlpb "nfxnews/protos/gen/crawl"
)

type CrawlHandler struct {
	crawlpb.UnimplementedCrawlServiceServer
	svc *crawlapp.Service
}

func NewCrawlHandler(svc *crawlapp.Service) *CrawlHandler { return &CrawlHandler{svc: svc} }

func (h *CrawlHandler) TriggerCrawl(ctx context.Context, req *crawlpb.TriggerCrawlRequest) (*crawlpb.TriggerCrawlResponse, error) {
	sess, err := h.svc.Trigger(ctx, "", "", req.GetSourceId())
	if err != nil {
		return nil, err
	}
	return &crawlpb.TriggerCrawlResponse{SessionId: sess.ID.String()}, nil
}

func (h *CrawlHandler) GetSession(ctx context.Context, req *crawlpb.GetSessionRequest) (*crawlpb.GetSessionResponse, error) {
	sess, err := h.svc.Get(ctx, "", req.GetSessionId())
	if err != nil {
		return nil, err
	}
	sid := ""
	if sess.SourceID != nil {
		sid = *sess.SourceID
	}
	errMsg := ""
	if sess.ErrorMessage != nil {
		errMsg = *sess.ErrorMessage
	}
	return &crawlpb.GetSessionResponse{Session: &crawlpb.CrawlSession{
		Id: sess.ID.String(), SourceId: sid, Status: sess.Status, ItemCount: int32(sess.ItemCount), ErrorMessage: errMsg,
	}}, nil
}
