package handler

import (
	"encoding/json"

	mcpapp "nfxnews/modules/mcp/application/mcp"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type MCPHandler struct{ svc *mcpapp.Service }

func NewMCPHandler(svc *mcpapp.Service) *MCPHandler { return &MCPHandler{svc: svc} }

type runBody struct {
	Arguments json.RawMessage `json:"arguments"`
}

func (h *MCPHandler) RunTool(c fiber.Ctx) error {
	name := c.Params("name")
	var body runBody
	_ = c.Bind().Body(&body)
	result, err := h.svc.RunTool(c.Context(), name, string(body.Arguments))
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: result})
}

func (h *MCPHandler) Tools(c fiber.Ctx) error {
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: h.svc.Tools()})
}

func (h *MCPHandler) RunJSON(c fiber.Ctx) error {
	var body map[string]any
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	name, _ := body["tool_name"].(string)
	args, _ := json.Marshal(body["arguments"])
	result, err := h.svc.RunTool(c.Context(), name, string(args))
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: result})
}
