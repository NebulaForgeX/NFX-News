package handler

import (
	"strconv"

	reportapp "nfxnews/modules/report/application/report"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type ReportHandler struct{ svc *reportapp.Service }

func NewReportHandler(svc *reportapp.Service) *ReportHandler { return &ReportHandler{svc: svc} }

type keywordBody struct {
	GroupName  string `json:"group_name"`
	Word       string `json:"word"`
	Kind       string `json:"kind"`
	CountLimit int    `json:"count_limit"`
}

type generateBody struct {
	Mode string `json:"mode"`
}

func (h *ReportHandler) ListKeywords(c fiber.Ctx) error {
	rows, err := h.svc.ListKeywords(c.Context())
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *ReportHandler) AddKeyword(c fiber.Ctx) error {
	var body keywordBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	row, err := h.svc.AddKeyword(c.Context(), body.GroupName, body.Word, body.Kind, body.CountLimit)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "created", httpx.SuccessOptions{Data: row})
}

func (h *ReportHandler) Generate(c fiber.Ctx) error {
	var body generateBody
	_ = c.Bind().Body(&body)
	snap, err := h.svc.Generate(c.Context(), body.Mode)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "generated", httpx.SuccessOptions{Data: snap})
}

func (h *ReportHandler) ListSnapshots(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.svc.ListSnapshots(c.Context(), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *ReportHandler) Get(c fiber.Ctx) error {
	row, err := h.svc.Get(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: row})
}
