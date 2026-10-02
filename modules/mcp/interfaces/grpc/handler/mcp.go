package handler

import (
	"context"
	"encoding/json"

	mcpapp "nfxnews/modules/mcp/application/mcp"
	mcppb "nfxnews/protos/gen/mcp"
)

type MCPHandler struct {
	mcppb.UnimplementedMCPServiceServer
	svc *mcpapp.Service
}

func NewMCPHandler(svc *mcpapp.Service) *MCPHandler { return &MCPHandler{svc: svc} }

func (h *MCPHandler) RunTool(ctx context.Context, req *mcppb.RunToolRequest) (*mcppb.RunToolResponse, error) {
	result, err := h.svc.RunTool(ctx, req.GetToolName(), req.GetArgumentsJson())
	if err != nil {
		return &mcppb.RunToolResponse{Ok: false, ErrorMessage: err.Error()}, nil
	}
	raw, _ := json.Marshal(result)
	return &mcppb.RunToolResponse{Ok: true, ResultJson: string(raw)}, nil
}
