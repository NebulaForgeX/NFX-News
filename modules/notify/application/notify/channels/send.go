package channels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	notifyerr "nfxnews/errors/src/notify"
	"nfxnews/pkgs/email"
)

type Sender func(ctx context.Context, client *http.Client, cfg map[string]any, text string) error

func ForKind(kind string) (Sender, error) {
	switch kind {
	case "feishu":
		return sendFeishu, nil
	case "dingtalk":
		return sendDingTalk, nil
	case "wework":
		return sendWeWork, nil
	case "telegram":
		return sendTelegram, nil
	case "slack":
		return sendSlack, nil
	case "ntfy":
		return sendNtfy, nil
	case "bark":
		return sendBark, nil
	case "email":
		return sendEmail, nil
	default:
		return nil, notifyerr.ErrNotifyKindInvalid
	}
}

func postJSON(ctx context.Context, client *http.Client, endpoint string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("%s: %s %s", endpoint, resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func sendFeishu(ctx context.Context, client *http.Client, cfg map[string]any, text string) error {
	webhook := cfgString(cfg, "webhook_url", "NOTIFY_FEISHU_WEBHOOK_URL")
	if webhook == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	return postJSON(ctx, client, webhook, map[string]any{
		"msg_type": "text",
		"content":  map[string]string{"text": text},
	})
}

func sendDingTalk(ctx context.Context, client *http.Client, cfg map[string]any, text string) error {
	webhook := cfgString(cfg, "webhook_url", "NOTIFY_DINGTALK_WEBHOOK_URL")
	if webhook == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	return postJSON(ctx, client, webhook, map[string]any{
		"msgtype": "text",
		"text":    map[string]string{"content": text},
	})
}

func sendWeWork(ctx context.Context, client *http.Client, cfg map[string]any, text string) error {
	webhook := cfgString(cfg, "webhook_url", "NOTIFY_WEWORK_WEBHOOK_URL")
	if webhook == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	msgType := cfgString(cfg, "msg_type", "NOTIFY_WEWORK_MSG_TYPE")
	if msgType == "" {
		msgType = "markdown"
	}
	if msgType == "text" {
		return postJSON(ctx, client, webhook, map[string]any{
			"msgtype": "text",
			"text":    map[string]string{"content": text},
		})
	}
	return postJSON(ctx, client, webhook, map[string]any{
		"msgtype":  "markdown",
		"markdown": map[string]string{"content": text},
	})
}

func sendTelegram(ctx context.Context, client *http.Client, cfg map[string]any, text string) error {
	token := cfgString(cfg, "bot_token", "NOTIFY_TELEGRAM_BOT_TOKEN")
	chat := cfgString(cfg, "chat_id", "NOTIFY_TELEGRAM_CHAT_ID")
	if token == "" || chat == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	endpoint := "https://api.telegram.org/bot" + token + "/sendMessage"
	return postJSON(ctx, client, endpoint, map[string]any{"chat_id": chat, "text": text})
}

func sendSlack(ctx context.Context, client *http.Client, cfg map[string]any, text string) error {
	webhook := cfgString(cfg, "webhook_url", "NOTIFY_SLACK_WEBHOOK_URL")
	if webhook == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	return postJSON(ctx, client, webhook, map[string]any{"text": text})
}

func sendNtfy(ctx context.Context, client *http.Client, cfg map[string]any, text string) error {
	topic := cfgString(cfg, "topic", "NOTIFY_NTFY_TOPIC")
	if topic == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	base := cfgString(cfg, "server_url", "NOTIFY_NTFY_SERVER_URL")
	if base == "" {
		base = "https://ntfy.sh"
	}
	endpoint := strings.TrimRight(base, "/") + "/" + topic
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(text))
	if err != nil {
		return err
	}
	req.Header.Set("Title", "NFX-News")
	if token := cfgString(cfg, "token", "NOTIFY_NTFY_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: %s", resp.Status)
	}
	return nil
}

func sendBark(ctx context.Context, client *http.Client, cfg map[string]any, text string) error {
	base := cfgString(cfg, "url", "NOTIFY_BARK_URL")
	if base == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	endpoint := strings.TrimRight(base, "/")
	if !strings.Contains(endpoint, "/push") {
		endpoint = endpoint + "/push"
	}
	return postJSON(ctx, client, endpoint, map[string]any{
		"title": "NFX-News",
		"body":  text,
		"url":   "",
	})
}

func sendEmail(_ context.Context, _ *http.Client, cfg map[string]any, text string) error {
	host := cfgString(cfg, "smtp_host", "EMAIL_SMTP_HOST")
	user := cfgString(cfg, "smtp_user", "EMAIL_SMTP_USER")
	pass := cfgString(cfg, "smtp_password", "EMAIL_SMTP_PASSWORD")
	from := cfgString(cfg, "from", "EMAIL_SMTP_FROM")
	to := cfgString(cfg, "to", "NOTIFY_EMAIL_TO")
	if to == "" {
		to = cfgString(cfg, "to", "EMAIL_SMTP_USER")
	}
	portStr := cfgString(cfg, "smtp_port", "EMAIL_SMTP_PORT")
	port := 587
	if n, err := strconv.Atoi(portStr); err == nil && n > 0 {
		port = n
	}
	if host == "" || from == "" || to == "" {
		return notifyerr.ErrNotifyWebhookMissing
	}
	svc := email.NewEmailService(email.SMTPConfig{
		Host: host, Port: port, Username: user, Password: pass, From: from,
	})
	subject := "NFX-News report"
	html := "<pre>" + htmlEscape(text) + "</pre>"
	return svc.Send(email.EmailMessage{
		To: strings.Split(to, ","), Subject: subject, Body: html, IsHTML: true,
	})
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(s)
}

func EnvKinds() []string {
	var out []string
	if os.Getenv("NOTIFY_FEISHU_WEBHOOK_URL") != "" {
		out = append(out, "feishu")
	}
	if os.Getenv("NOTIFY_DINGTALK_WEBHOOK_URL") != "" {
		out = append(out, "dingtalk")
	}
	if os.Getenv("NOTIFY_WEWORK_WEBHOOK_URL") != "" {
		out = append(out, "wework")
	}
	if os.Getenv("NOTIFY_TELEGRAM_BOT_TOKEN") != "" && os.Getenv("NOTIFY_TELEGRAM_CHAT_ID") != "" {
		out = append(out, "telegram")
	}
	if os.Getenv("NOTIFY_SLACK_WEBHOOK_URL") != "" {
		out = append(out, "slack")
	}
	if os.Getenv("NOTIFY_NTFY_TOPIC") != "" {
		out = append(out, "ntfy")
	}
	if os.Getenv("NOTIFY_BARK_URL") != "" {
		out = append(out, "bark")
	}
	if os.Getenv("NOTIFY_EMAIL_TO") != "" || os.Getenv("EMAIL_SMTP_FROM") != "" {
		if os.Getenv("EMAIL_SMTP_HOST") != "" {
			out = append(out, "email")
		}
	}
	return out
}

func BatchInterval(cfg map[string]any) time.Duration {
	if n := cfgInt(cfg, "batch_send_interval"); n > 0 {
		return time.Duration(n) * time.Second
	}
	if v := os.Getenv("NOTIFY_BATCH_SEND_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return 3 * time.Second
}

func SendAll(ctx context.Context, client *http.Client, kind string, cfg map[string]any, text string) error {
	sender, err := ForKind(kind)
	if err != nil {
		return err
	}
	batches := SplitBatches(text, MaxBytes(kind, cfg))
	wait := BatchInterval(cfg)
	for i, batch := range batches {
		body := batch
		if len(batches) > 1 {
			body = fmt.Sprintf("[batch %d/%d]\n%s", i+1, len(batches), batch)
		}
		if err := sender(ctx, client, cfg, body); err != nil {
			return notifyerr.ErrNotifySendFailed.WithCause(err)
		}
		if i < len(batches)-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		}
	}
	return nil
}
