package crawl

import "nfxnews/pkgs/errx"

var (
	ErrCrawlSessionNotFound = errx.NotFound("CRAWL_SESSION_NOT_FOUND", "crawl session not found")
	ErrCrawlTriggerFailed   = errx.Internal("CRAWL_TRIGGER_FAILED", "failed to trigger crawl")
	ErrCrawlAlreadyRunning  = errx.Conflict("CRAWL_ALREADY_RUNNING", "a crawl is already running")
)

/*
!CRAWL_SESSION_NOT_FOUND
*en<crawl session not found>
*zh<抓取会话不存在>
*fr<session de crawl introuvable>

!CRAWL_TRIGGER_FAILED
*en<failed to trigger crawl>
*zh<触发抓取失败>
*fr<échec du déclenchement du crawl>

!CRAWL_ALREADY_RUNNING
*en<a crawl is already running>
*zh<已有抓取正在进行>
*fr<un crawl est déjà en cours>
*/
