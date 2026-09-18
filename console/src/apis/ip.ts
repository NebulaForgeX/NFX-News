const HTTP_BASE_URL = import.meta.env.VITE_API_URL || "http://127.0.0.1:10178";

export const URL_PATHS = {
  SOURCE: {
    sources: "/source/sources",
    byId: (id: string) => `/source/sources/${id}`,
    fetch: (id: string) => `/source/sources/${id}/fetch`,
  },
  NEWS: {
    items: "/news/items",
    search: "/news/search",
    preferences: "/news/preferences",
  },
  CRAWL: {
    sessions: "/crawl/sessions",
    byId: (id: string) => `/crawl/sessions/${id}`,
  },
  REPORT: {
    keywords: "/report/keywords",
    snapshots: "/report/snapshots",
    byId: (id: string) => `/report/snapshots/${id}`,
    html: (id: string) => `/report/snapshots/${id}/html`,
  },
  NOTIFY: {
    kinds: "/notify/kinds",
    channels: "/notify/channels",
    deliveries: "/notify/deliveries",
    dispatch: "/notify/dispatch",
  },
  MCP: {
    tools: "/mcp/tools",
    run: "/mcp/run",
    byName: (name: string) => `/mcp/tools/${name}`,
  },
  SYSTEM: {
    latest: "/system/system-state/latest",
    initialize: "/system/system-state/initialize",
    i18nErrors: (lang: string) => `/system/i18n/errors/${lang}`,
  },
} as const;

export const API_ENDPOINTS = {
  PURE: HTTP_BASE_URL,
} as const;
