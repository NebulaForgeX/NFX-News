package handler

import (
	"strconv"

	authconn "nfxnews/connections/auth"
	reportapp "nfxnews/modules/report/application/report"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type ReportHandler struct {
	svc      *reportapp.Service
	identity *authconn.Client
}

func NewReportHandler(svc *reportapp.Service, identity *authconn.Client) *ReportHandler {
	return &ReportHandler{svc: svc, identity: identity}
}

type keywordBody struct {
	GroupName  string `json:"group_name"`
	Word       string `json:"word"`
	Kind       string `json:"kind"`
	CountLimit int    `json:"count_limit"`
}

type generateBody struct {
	Mode string `json:"mode"`
}

func (h *ReportHandler) accountProfile(c fiber.Ctx) (accountID, profileID string, ferr *errx.Error) {
	aid, ok := fiberx.AccountIDFromContext(c.Context())
	if !ok {
		return "", "", errx.Unauthorized("INVALID_TOKEN", "missing account")
	}
	pid, ok := fiberx.ProfileIDFromContext(c.Context())
	if !ok {
		return "", "", errx.Unauthorized("INVALID_TOKEN", "missing profile")
	}
	scope, _ := fiberx.ProfileScopeFromContext(c.Context())
	if h.identity != nil {
		allowed, err := h.identity.Account.EnsureOwnedProfile(c.Context(), aid, pid, scope)
		if err != nil {
			return "", "", errx.Unauthorized("IDENTITY_UNAVAILABLE", "identity lookup failed").WithCause(err)
		}
		if !allowed {
			return "", "", errx.Unauthorized("PROFILE_NOT_OWNED", "profile does not belong to account")
		}
	}
	return aid.String(), pid.String(), nil
}

func (h *ReportHandler) ListKeywords(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	rows, err := h.svc.ListKeywords(c.Context(), aid)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *ReportHandler) AddKeyword(c fiber.Ctx) error {
	aid, pid, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	var body keywordBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	row, err := h.svc.AddKeyword(c.Context(), aid, pid, body.GroupName, body.Word, body.Kind, body.CountLimit)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "created", httpx.SuccessOptions{Data: row})
}

func (h *ReportHandler) Generate(c fiber.Ctx) error {
	aid, pid, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	var body generateBody
	_ = c.Bind().Body(&body)
	snap, err := h.svc.Generate(c.Context(), aid, pid, body.Mode)
	if err != nil {
		return err
	}
	return fiberx.Created(c, "generated", httpx.SuccessOptions{Data: snap})
}

func (h *ReportHandler) ListSnapshots(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.svc.ListSnapshots(c.Context(), aid, limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *ReportHandler) Get(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	row, err := h.svc.Get(c.Context(), aid, c.Params("id"))
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: row})
}

func (h *ReportHandler) HTML(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	html, err := h.svc.HTML(c.Context(), aid, c.Params("id"))
	if err != nil {
		return err
	}
	c.Type("html", "utf-8")
	return c.SendString(html)
}
