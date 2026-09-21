package handler

import (
	"strconv"

	authconn "nfxnews/connections/auth"
	crawlmsg "nfxnews/messages/src/crawl"
	crawlapp "nfxnews/modules/crawl/application/crawl"
	"nfxnews/pkgs/errx"
	"nfxnews/pkgs/fiberx"
	"nfxnews/pkgs/httpx"

	"github.com/gofiber/fiber/v3"
)

type CrawlHandler struct {
	svc      *crawlapp.Service
	identity *authconn.Client
}

func NewCrawlHandler(svc *crawlapp.Service, identity *authconn.Client) *CrawlHandler {
	return &CrawlHandler{svc: svc, identity: identity}
}

type triggerBody struct {
	SourceID string `json:"source_id"`
}

func (h *CrawlHandler) accountProfile(c fiber.Ctx) (accountID, profileID string, ferr *errx.Error) {
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

func (h *CrawlHandler) Trigger(c fiber.Ctx) error {
	aid, pid, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	var body triggerBody
	_ = c.Bind().Body(&body)
	sess, err := h.svc.Trigger(c.Context(), aid, pid, body.SourceID)
	if err != nil {
		return err
	}
	return fiberx.Created(c, crawlmsg.CRAWL_TRIGGERED, httpx.SuccessOptions{Data: sess})
}

func (h *CrawlHandler) List(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.svc.List(c.Context(), aid, limit)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: rows})
}

func (h *CrawlHandler) Get(c fiber.Ctx) error {
	aid, _, ferr := h.accountProfile(c)
	if ferr != nil {
		return fiberx.ErrorFromErrx(c, ferr)
	}
	id := c.Params("id")
	if id == "" {
		return errx.ErrInvalidParams.WithMsg("id required")
	}
	row, err := h.svc.Get(c.Context(), aid, id)
	if err != nil {
		return err
	}
	return fiberx.OK(c, "ok", httpx.SuccessOptions{Data: row})
}
