package handler

import (
	"strconv"

	authconn "nfxnews/connections/auth"
	"nfxnews/events"
	notifymsg "nfxnews/messages/src/notify"
	notifyapp "nfxnews/modules/notify/application/notify"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type NotifyHandler struct {
	svc      *notifyapp.Service
	identity *authconn.Client
}

func NewNotifyHandler(svc *notifyapp.Service, identity *authconn.Client) *NotifyHandler {
	return &NotifyHandler{svc: svc, identity: identity}
}

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
	Payload   string `json:"payload_json"`
}

func (h *NotifyHandler) accountProfile(c fiber.Ctx) (accountID, profileID string, ferr *errx.Error) {
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

func (h *NotifyHandler) ListKinds(c fiber.Ctx) error {
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: h.svc.ListKinds()})
}

func (h *NotifyHandler) ListChannels(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	rows, err := h.svc.ListChannels(c.Context(), aid)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *NotifyHandler) ListDeliveries(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	rows, err := h.svc.ListDeliveries(c.Context(), aid, limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *NotifyHandler) UpsertChannel(c fiber.Ctx) error {
	aid, pid, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	var body channelBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	row, err := h.svc.UpsertChannel(c.Context(), aid, pid, body.Kind, body.Name, body.Enabled, body.Config)
	if err != nil {
		return err
	}
	return fiberx.Created(c, notifymsg.CHANNEL_CREATED, httpx.SuccessOptions{Data: row})
}

func (h *NotifyHandler) Dispatch(c fiber.Ctx) error {
	aid, pid, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	var body dispatchBody
	if err := c.Bind().Body(&body); err != nil {
		return errx.ErrInvalidBody.WithCause(err)
	}
	n, err := h.svc.DispatchReport(c.Context(), events.ReportGeneratedEvent{
		ReportID: body.ReportID, AccountID: aid, ProfileID: pid, Mode: body.Mode, Title: body.Title, ItemCount: body.ItemCount, Payload: body.Payload,
	})
	if err != nil {
		return err
	}
	return fiberx.OK(c, notifymsg.NOTIFY_DISPATCHED, httpx.SuccessOptions{Data: map[string]any{"queued": n}})
}
