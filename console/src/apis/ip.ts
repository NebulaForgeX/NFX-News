const HTTP_BASE_URL = import.meta.env.VITE_API_URL || "http://127.0.0.1:10166";

export const URL_PATHS = {
  SOURCE: {
    sources: "/source/sources",
    byId: (id: string) => `/source/sources/${id}`,
    fetch: (id: string) => `/source/sources/${id}/fetch`,
    i18nErrors: (lang: string) => `/source/i18n/errors/${lang}`,
  },
  NEWS: {
    items: "/news/items",
    search: "/news/search",
    preferences: "/news/preferences",
    i18nErrors: (lang: string) => `/news/i18n/errors/${lang}`,
  },
  CRAWL: {
    sessions: "/crawl/sessions",
    byId: (id: string) => `/crawl/sessions/${id}`,
    i18nErrors: (lang: string) => `/crawl/i18n/errors/${lang}`,
  },
  REPORT: {
    keywords: "/report/keywords",
    snapshots: "/report/snapshots",
    byId: (id: string) => `/report/snapshots/${id}`,
    html: (id: string) => `/report/snapshots/${id}/html`,
    i18nErrors: (lang: string) => `/report/i18n/errors/${lang}`,
  },
  NOTIFY: {
    kinds: "/notify/kinds",
    channels: "/notify/channels",
    deliveries: "/notify/deliveries",
    dispatch: "/notify/dispatch",
    i18nErrors: (lang: string) => `/notify/i18n/errors/${lang}`,
  },
  MCP: {
    tools: "/mcp/tools",
    run: "/mcp/run",
    byName: (name: string) => `/mcp/tools/${name}`,
    i18nErrors: (lang: string) => `/mcp/i18n/errors/${lang}`,
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
