package mcpapp

import (
	"context"
	"encoding/json"

	"nfxnews/pkgs/errx"
	newspb "nfxnews/protos/gen/news"
	reportpb "nfxnews/protos/gen/report"
	sourcepb "nfxnews/protos/gen/source"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ToolCall struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ToolName     string
	Arguments    []byte `gorm:"type:jsonb"`
	OK           bool
	ErrorMessage *string
}

func (ToolCall) TableName() string { return "mcp.tool_calls" }

type Service struct {
	db     *gorm.DB
	news   newspb.NewsServiceClient
	report reportpb.ReportServiceClient
	source sourcepb.SourceServiceClient
}

func NewService(db *gorm.DB, news newspb.NewsServiceClient, report reportpb.ReportServiceClient, source sourcepb.SourceServiceClient) *Service {
	return &Service{db: db, news: news, report: report, source: source}
}

func (s *Service) RunTool(ctx context.Context, name, argsJSON string) (map[string]any, error) {
	args := map[string]any{}
	if argsJSON != "" {
		_ = json.Unmarshal([]byte(argsJSON), &args)
	}
	raw, _ := json.Marshal(args)
	call := ToolCall{ID: uuid.Must(uuid.NewV7()), ToolName: name, Arguments: raw, OK: true}
	var result map[string]any
	var runErr error
	switch name {
	case "get_latest_news":
		if s.news == nil {
			runErr = errx.ErrInternal.WithMsg("news client unavailable")
			break
		}
		limit := int32(20)
		if v, ok := args["limit"].(float64); ok {
			limit = int32(v)
		}
		resp, err := s.news.SearchNews(ctx, &newspb.SearchNewsRequest{Limit: limit})
		if err != nil {
			runErr = err
		} else {
			result = map[string]any{"count": len(resp.GetItems()), "items": resp.GetItems()}
		}
	case "search_news":
		if s.news == nil {
			runErr = errx.ErrInternal.WithMsg("news client unavailable")
			break
		}
		q, _ := args["query"].(string)
		resp, err := s.news.SearchNews(ctx, &newspb.SearchNewsRequest{Query: q, Limit: 50})
		if err != nil {
			runErr = err
		} else {
			result = map[string]any{"count": len(resp.GetItems()), "items": resp.GetItems()}
		}
	case "list_sources":
		if s.source == nil {
			runErr = errx.ErrInternal.WithMsg("source client unavailable")
			break
		}
		resp, err := s.source.ListSources(ctx, &sourcepb.ListSourcesRequest{})
		if err != nil {
			runErr = err
		} else {
			result = map[string]any{"sources": resp.GetSources()}
		}
	case "generate_report":
		if s.report == nil {
			runErr = errx.ErrInternal.WithMsg("report client unavailable")
			break
		}
		mode, _ := args["mode"].(string)
		resp, err := s.report.GenerateReport(ctx, &reportpb.GenerateReportRequest{Mode: mode})
		if err != nil {
			runErr = err
		} else {
			result = map[string]any{"report_id": resp.GetReportId(), "item_count": resp.GetItemCount()}
		}
	case "list_keywords":
		if s.report == nil {
			runErr = errx.ErrInternal.WithMsg("report client unavailable")
			break
		}
		resp, err := s.report.ListKeywords(ctx, &reportpb.ListKeywordsRequest{})
		if err != nil {
			runErr = err
		} else {
			result = map[string]any{"keywords": resp.GetKeywords()}
		}
	default:
		runErr = errx.ErrInvalidParams.WithMsg("unknown tool: " + name)
	}
	if runErr != nil {
		msg := runErr.Error()
		call.OK = false
		call.ErrorMessage = &msg
		_ = s.db.WithContext(ctx).Create(&call).Error
		return nil, runErr
	}
	_ = s.db.WithContext(ctx).Create(&call).Error
	if result == nil {
		result = map[string]any{}
	}
	return result, nil
}
