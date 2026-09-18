export const NOTIFY_KINDS = ["feishu", "dingtalk", "wework", "telegram", "email", "ntfy", "bark", "slack"] as const;
export type NotifyKind = (typeof NOTIFY_KINDS)[number];

export const REPORT_MODES = ["daily", "current", "incremental"] as const;
export type ReportMode = (typeof REPORT_MODES)[number];

export const KEYWORD_KINDS = ["include", "exclude", "required"] as const;
export type KeywordKind = (typeof KEYWORD_KINDS)[number];

export function notifyChannelConfig(kind: NotifyKind, value: string): Record<string, unknown> {
  const v = value.trim();
  if (!v) return {};
  switch (kind) {
    case "telegram":
      return { chatId: v };
    case "ntfy":
      return { topic: v };
    case "bark":
      return { url: v };
    case "email":
      return { to: v };
    default:
      return { webhookUrl: v };
  }
}

export function notifyConfigPlaceholder(kind: NotifyKind): string {
  switch (kind) {
    case "telegram":
      return "chat_id";
    case "ntfy":
      return "topic";
    case "bark":
      return "url";
    case "email":
      return "to";
    default:
      return "webhook_url";
  }
}
