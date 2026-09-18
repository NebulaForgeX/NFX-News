package handler

import (
	sourceapp "nfxnews/modules/source/application/source"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type SourceHandler struct{ svc *sourceapp.Service }

func NewSourceHandler(svc *sourceapp.Service) *SourceHandler { return &SourceHandler{svc: svc} }

func (h *SourceHandler) List(c fiber.Ctx) error {
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: h.svc.List()})
}

func (h *SourceHandler) Get(c fiber.Ctx) error {
	id := c.Params("id")
	if err := sourceapp.RequireID(id); err != nil {
		return err
	}
	meta, err := h.svc.Meta(id)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: meta})
}

func (h *SourceHandler) Fetch(c fiber.Ctx) error {
	id := c.Params("id")
	if err := sourceapp.RequireID(id); err != nil {
		return err
	}
	items, err := h.svc.Fetch(c.Context(), id)
	if err != nil {
		return err
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		out = append(out, map[string]any{
			"id": it.ID, "source_id": id, "original_id": it.ID, "title": it.Title,
			"url": it.URL, "mobile_url": it.MobileURL, "pub_date": it.PubDate, "extra": it.Extra,
		})
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: out})
}
