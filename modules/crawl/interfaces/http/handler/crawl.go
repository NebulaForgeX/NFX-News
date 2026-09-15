package handler

import (
	"strconv"

	crawlapp "nfxnews/modules/crawl/application/crawl"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type CrawlHandler struct{ svc *crawlapp.Service }

func NewCrawlHandler(svc *crawlapp.Service) *CrawlHandler { return &CrawlHandler{svc: svc} }

type triggerBody struct {
	SourceID string `json:"source_id"`
}

func (h *CrawlHandler) Trigger(c fiber.Ctx) error {
	var body triggerBody
	_ = c.Bind().Body(&body)
	sess, err := h.svc.Trigger(c.Context(), body.SourceID)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "triggered", httpx.SuccessOptions{Data: sess})
}

func (h *CrawlHandler) List(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.svc.List(c.Context(), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *CrawlHandler) Get(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errx.ErrInvalidParams.WithMsg("id required")
	}
	row, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: row})
}
