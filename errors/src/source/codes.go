package source

import "nfxnews/pkgs/errx"

var (
	ErrSourceNotFound     = errx.NotFound("SOURCE_NOT_FOUND", "source not found")
	ErrSourceFetchFailed  = errx.Internal("SOURCE_FETCH_FAILED", "source fetch failed")
	ErrSourceIDRequired   = errx.InvalidArg("SOURCE_ID_REQUIRED", "source id is required")
	ErrSourceCatalogEmpty = errx.Internal("SOURCE_CATALOG_EMPTY", "source catalog is empty")
)

/*
!SOURCE_NOT_FOUND
*en<source not found>
*zh<来源不存在>
*fr<source introuvable>

!SOURCE_FETCH_FAILED
*en<source fetch failed>
*zh<来源抓取失败>
*fr<échec de récupération de la source>

!SOURCE_ID_REQUIRED
*en<source id is required>
*zh<来源 ID 必填>
*fr<l'identifiant de source est requis>

!SOURCE_CATALOG_EMPTY
*en<source catalog is empty>
*zh<来源目录为空>
*fr<le catalogue de sources est vide>
*/
