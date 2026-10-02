package handler

import (
	"context"

	"nfxnews/events"
	notifyapp "nfxnews/modules/notify/application/notify"
	notifypb "nfxnews/protos/gen/notify"
)

type NotifyHandler struct {
	notifypb.UnimplementedNotifyServiceServer
	svc *notifyapp.Service
}

func NewNotifyHandler(svc *notifyapp.Service) *NotifyHandler { return &NotifyHandler{svc: svc} }

func (h *NotifyHandler) DispatchReport(ctx context.Context, req *notifypb.DispatchReportRequest) (*notifypb.DispatchReportResponse, error) {
	n, err := h.svc.DispatchReport(ctx, events.ReportGeneratedEvent{ReportID: req.GetReportId(), Title: "report", Mode: "manual"})
	if err != nil {
		return nil, err
	}
	return &notifypb.DispatchReportResponse{Queued: int32(n)}, nil
}
