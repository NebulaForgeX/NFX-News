package report

import "nfxnews/pkgs/errx"

var (
	ErrReportNotFound       = errx.NotFound("REPORT_NOT_FOUND", "report not found")
	ErrReportKeywordRequired = errx.InvalidArg("REPORT_KEYWORD_REQUIRED", "keyword word is required")
	ErrReportModeInvalid    = errx.InvalidArg("REPORT_MODE_INVALID", "report mode must be daily, current, or incremental")
	ErrReportGenerateFailed = errx.Internal("REPORT_GENERATE_FAILED", "failed to generate report")
	ErrFrequencyFileMissing = errx.NotFound("FREQUENCY_FILE_MISSING", "frequency words file not found")
)

/*
!REPORT_NOT_FOUND
*en<report not found>
*zh<报告不存在>
*fr<rapport introuvable>

!REPORT_KEYWORD_REQUIRED
*en<keyword word is required>
*zh<关键词必填>
*fr<le mot-clé est requis>

!REPORT_MODE_INVALID
*en<report mode must be daily, current, or incremental>
*zh<报告模式必须是 daily、current 或 incremental>
*fr<le mode de rapport doit être daily, current ou incremental>

!REPORT_GENERATE_FAILED
*en<failed to generate report>
*zh<生成报告失败>
*fr<échec de la génération du rapport>

!FREQUENCY_FILE_MISSING
*en<frequency words file not found>
*zh<关注词文件不存在>
*fr<fichier de mots de fréquence introuvable>
*/
