package keyword

import (
	"strings"

	reporterr "nfxnews/errors/src/report"
)

func validateWord(word, kind string) error {
	if strings.TrimSpace(word) == "" {
		return reporterr.ErrReportKeywordRequired
	}
	switch kind {
	case "", "include", "exclude", "required":
		return nil
	default:
		return reporterr.ErrReportKeywordRequired
	}
}
