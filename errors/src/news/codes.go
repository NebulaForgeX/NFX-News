package news

import "nfxnews/pkgs/errx"

var (
	ErrNewsNotFound      = errx.NotFound("NEWS_NOT_FOUND", "news item not found")
	ErrNewsQueryRequired = errx.InvalidArg("NEWS_QUERY_REQUIRED", "search query is required")
	ErrNewsPersistFailed = errx.Internal("NEWS_PERSIST_FAILED", "failed to persist news items")
)

/*
!NEWS_NOT_FOUND
*en<news item not found>
*zh<新闻不存在>
*fr<actualité introuvable>

!NEWS_QUERY_REQUIRED
*en<search query is required>
*zh<搜索关键词必填>
*fr<la requête de recherche est requise>

!NEWS_PERSIST_FAILED
*en<failed to persist news items>
*zh<新闻写入失败>
*fr<échec de l'enregistrement des actualités>
*/
