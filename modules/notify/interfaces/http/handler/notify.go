package handler

import (
	"strconv"

	"nfxnews/events"
	notifyapp "nfxnews/modules/notify/application/notify"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type NotifyHandler struct{ svc *notifyapp.Service }

func NewNotifyHandler(svc *notifyapp.Service) *NotifyHandler { return &NotifyHandler{svc: svc} }

type channelBody struct {
	Kind    string         `json:"kind"`
	Name    string         `json:"name"`
	Enabled bool           `json:"enabled"`
	Config  map[string]any `json:"config"`
}

type dispatchBody struct {
	ReportID  string `json:"report_id"`
	Mode      string `json:"mode"`
	Title     string `json:"title"`
	ItemCount int    `json:"item_count"`
}

func (h *NotifyHandler) ListChannels(c fiber.Ctx) error {
	rows, err := h.svc.ListChannels(c.Context())
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *NotifyHandler) ListDeliveries(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	rows, err := h.svc.ListDeliveries(c.Context(), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *NotifyHandler) UpsertChannel(c fiber.Ctx) error {
	var body channelBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	row, err := h.svc.UpsertChannel(c.Context(), body.Kind, body.Name, body.Enabled, body.Config)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "created", httpx.SuccessOptions{Data: row})
}

func (h *NotifyHandler) Dispatch(c fiber.Ctx) error {
	var body dispatchBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	n, err := h.svc.DispatchReport(c.Context(), events.ReportGeneratedEvent{
		ReportID: body.ReportID, Mode: body.Mode, Title: body.Title, ItemCount: body.ItemCount,
	})
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: map[string]any{"queued": n}})
}
