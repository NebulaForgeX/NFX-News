package handler

import (
	"encoding/json"
	"strconv"

	authconn "nfxnews/connections/auth"
	newsmsg "nfxnews/messages/src/news"
	newsapp "nfxnews/modules/news/application/news"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type NewsHandler struct {
	svc      *newsapp.Service
	identity *authconn.Client
}

func NewNewsHandler(svc *newsapp.Service, identity *authconn.Client) *NewsHandler {
	return &NewsHandler{svc: svc, identity: identity}
}

func (h *NewsHandler) List(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, err := h.svc.ListBySource(c.Context(), c.Query("source_id"), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: items})
}

func (h *NewsHandler) Search(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, err := h.svc.Search(c.Context(), c.Query("q"), limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: items})
}

func (h *NewsHandler) accountProfile(c fiber.Ctx) (accountID, profileID string, ferr *errx.Error) {
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

func (h *NewsHandler) GetPreferences(c fiber.Ctx) error {
	aid, pid, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	pref, err := h.svc.GetPreferences(c.Context(), aid, pid)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: pref})
}

func (h *NewsHandler) SetPreferences(c fiber.Ctx) error {
	aid, pid, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	var req struct {
		ColumnOrder json.RawMessage `json:"column_order"`
		Payload     json.RawMessage `json:"payload"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	if len(req.ColumnOrder) == 0 {
		req.ColumnOrder = json.RawMessage("[]")
	}
	if len(req.Payload) == 0 {
		req.Payload = json.RawMessage("{}")
	}
	if err := h.svc.SetPreferences(c.Context(), aid, pid, req.ColumnOrder, req.Payload); err != nil {
		return err
	}
	return fiberx.OK(c, newsmsg.PREFERENCES_UPDATED, httpx.SuccessOptions{Data: nil})
}
