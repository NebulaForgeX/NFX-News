package handler

import (
	"context"

	reportapp "nfxnews/modules/report/application/report"
	reportpb "nfxnews/protos/gen/report"
)

type ReportHandler struct {
	reportpb.UnimplementedReportServiceServer
	svc *reportapp.Service
}

func NewReportHandler(svc *reportapp.Service) *ReportHandler { return &ReportHandler{svc: svc} }

func (h *ReportHandler) GenerateReport(ctx context.Context, req *reportpb.GenerateReportRequest) (*reportpb.GenerateReportResponse, error) {
	snap, err := h.svc.Generate(ctx, "", "", req.GetMode())
	if err != nil {
		return nil, err
	}
	return &reportpb.GenerateReportResponse{ReportId: snap.ID.String(), ItemCount: int32(snap.ItemCount)}, nil
}

func (h *ReportHandler) GetReport(ctx context.Context, req *reportpb.GetReportRequest) (*reportpb.GetReportResponse, error) {
	snap, err := h.svc.Get(ctx, "", req.GetReportId())
	if err != nil {
		return nil, err
	}
	return &reportpb.GetReportResponse{
		Id: snap.ID.String(), Mode: snap.Mode, Title: snap.Title, PayloadJson: string(snap.Payload), ItemCount: int32(snap.ItemCount),
	}, nil
}

func (h *ReportHandler) ListKeywords(ctx context.Context, req *reportpb.ListKeywordsRequest) (*reportpb.ListKeywordsResponse, error) {
	rows, err := h.svc.ListKeywords(ctx, "")
	if err != nil {
		return nil, err
	}
	out := make([]*reportpb.Keyword, 0, len(rows))
	for _, k := range rows {
		out = append(out, &reportpb.Keyword{Id: k.ID.String(), GroupName: k.GroupName, Word: k.Word, Kind: k.Kind, CountLimit: int32(k.CountLimit)})
	}
	return &reportpb.ListKeywordsResponse{Keywords: out}, nil
}
