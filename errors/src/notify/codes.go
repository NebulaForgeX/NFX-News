package notify

import "nfxnews/pkgs/errx"

var (
	ErrNotifyChannelNotFound = errx.NotFound("NOTIFY_CHANNEL_NOT_FOUND", "notify channel not found")
	ErrNotifyKindInvalid     = errx.InvalidArg("NOTIFY_KIND_INVALID", "unsupported notify channel kind")
	ErrNotifyWebhookMissing  = errx.InvalidArg("NOTIFY_WEBHOOK_MISSING", "webhook or credentials missing")
	ErrNotifySendFailed      = errx.Internal("NOTIFY_SEND_FAILED", "failed to send notification")
	ErrNotifyWindowClosed    = errx.Conflict("NOTIFY_WINDOW_CLOSED", "outside the configured push window")
)

/*
!NOTIFY_CHANNEL_NOT_FOUND
*en<notify channel not found>
*zh<通知渠道不存在>
*fr<canal de notification introuvable>

!NOTIFY_KIND_INVALID
*en<unsupported notify channel kind>
*zh<不支持的通知渠道类型>
*fr<type de canal de notification non pris en charge>

!NOTIFY_WEBHOOK_MISSING
*en<webhook or credentials missing>
*zh<缺少 webhook 或凭证>
*fr<webhook ou identifiants manquants>

!NOTIFY_SEND_FAILED
*en<failed to send notification>
*zh<发送通知失败>
*fr<échec de l'envoi de la notification>

!NOTIFY_WINDOW_CLOSED
*en<outside the configured push window>
*zh<不在推送时间窗内>
*fr<en dehors de la fenêtre de push configurée>
*/
