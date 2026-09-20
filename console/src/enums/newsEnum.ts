export const NOTIFY_KINDS = ["feishu", "dingtalk", "wework", "telegram", "email", "ntfy", "bark", "slack"] as const;
export type NotifyKind = (typeof NOTIFY_KINDS)[number];

export const REPORT_MODES = ["daily", "current", "incremental"] as const;
export type ReportMode = (typeof REPORT_MODES)[number];

export const KEYWORD_KINDS = ["include", "exclude", "required"] as const;
export type KeywordKind = (typeof KEYWORD_KINDS)[number];

export const CRAWL_STATUSES = ["running", "ok", "failed"] as const;
export type CrawlStatus = (typeof CRAWL_STATUSES)[number];

export type NotifyConfigField = {
  key: string;
  kind: "string" | "number" | "boolean";
  secret?: boolean;
};

const sharedNotifyFields: NotifyConfigField[] = [
  { key: "batchSendInterval", kind: "number" },
  { key: "pushWindowEnabled", kind: "boolean" },
  { key: "pushWindowStart", kind: "string" },
  { key: "pushWindowEnd", kind: "string" },
  { key: "oncePerDay", kind: "boolean" },
];

export function notifyConfigFields(kind: NotifyKind): NotifyConfigField[] {
  switch (kind) {
    case "telegram":
      return [
        { key: "botToken", kind: "string", secret: true },
        { key: "chatId", kind: "string" },
        ...sharedNotifyFields,
      ];
    case "ntfy":
      return [
        { key: "topic", kind: "string" },
        { key: "serverUrl", kind: "string" },
        { key: "token", kind: "string", secret: true },
        ...sharedNotifyFields,
      ];
    case "bark":
      return [{ key: "url", kind: "string" }, ...sharedNotifyFields];
    case "email":
      return [
        { key: "to", kind: "string" },
        { key: "from", kind: "string" },
        { key: "smtpHost", kind: "string" },
        { key: "smtpPort", kind: "number" },
        { key: "smtpUser", kind: "string" },
        { key: "smtpPassword", kind: "string", secret: true },
        ...sharedNotifyFields,
      ];
    case "wework":
      return [
        { key: "webhookUrl", kind: "string" },
        { key: "msgType", kind: "string" },
        ...sharedNotifyFields,
      ];
    default:
      return [{ key: "webhookUrl", kind: "string" }, ...sharedNotifyFields];
  }
}

export type MCPArgKind = "string" | "number";

export type MCPToolArg = {
  key: string;
  kind: MCPArgKind;
};

export type MCPToolDef = {
  name: string;
  args: MCPToolArg[];
};

export const MCP_TOOL_DEFS: MCPToolDef[] = [
  { name: "get_latest_news", args: [{ key: "limit", kind: "number" }, { key: "platforms", kind: "string" }] },
  { name: "get_trending_topics", args: [{ key: "top_n", kind: "number" }, { key: "mode", kind: "string" }] },
  { name: "get_news_by_date", args: [{ key: "date_query", kind: "string" }, { key: "platforms", kind: "string" }, { key: "limit", kind: "number" }] },
  { name: "analyze_topic_trend", args: [{ key: "topic", kind: "string" }, { key: "analysis_type", kind: "string" }] },
  { name: "analyze_data_insights", args: [{ key: "insight_type", kind: "string" }, { key: "topic", kind: "string" }] },
  { name: "analyze_sentiment", args: [{ key: "topic", kind: "string" }, { key: "limit", kind: "number" }, { key: "platforms", kind: "string" }] },
  { name: "find_similar_news", args: [{ key: "reference_title", kind: "string" }, { key: "threshold", kind: "string" }, { key: "limit", kind: "number" }] },
  { name: "generate_summary_report", args: [{ key: "report_type", kind: "string" }] },
  { name: "search_news", args: [{ key: "query", kind: "string" }, { key: "search_mode", kind: "string" }, { key: "limit", kind: "number" }, { key: "platforms", kind: "string" }] },
  { name: "search_related_news_history", args: [{ key: "reference_text", kind: "string" }, { key: "time_preset", kind: "string" }] },
  { name: "get_current_config", args: [{ key: "section", kind: "string" }] },
  { name: "get_system_status", args: [] },
  { name: "trigger_crawl", args: [{ key: "platforms", kind: "string" }] },
  { name: "list_sources", args: [] },
  { name: "list_keywords", args: [] },
];

export const MCP_TOOL_MAP = Object.fromEntries(MCP_TOOL_DEFS.map((tool) => [tool.name, tool])) as Record<string, MCPToolDef>;

export const SOURCE_COLOR_VARS: Record<string, string> = {
  slate: "var(--gray-9)",
  gray: "var(--gray-9)",
  blue: "var(--blue-9)",
  red: "var(--red-9)",
  green: "var(--green-9)",
  indigo: "var(--indigo-9)",
  orange: "var(--orange-9)",
  amber: "var(--amber-9)",
  violet: "var(--violet-9)",
  cyan: "var(--cyan-9)",
};

export function sourceColorVar(color?: string): string {
  if (!color) return "var(--accent-9)";
  return SOURCE_COLOR_VARS[color] ?? "var(--accent-9)";
}
