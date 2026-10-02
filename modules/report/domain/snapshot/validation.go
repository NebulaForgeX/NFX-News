package snapshot

import (
	"strings"

	reporterr "nfxnews/errors/src/report"
)

func NormalizeMode(mode string) (string, error) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		mode = "daily"
	}
	switch mode {
	case "daily", "current", "incremental":
		return mode, nil
	default:
		return "", reporterr.ErrReportModeInvalid
	}
}
