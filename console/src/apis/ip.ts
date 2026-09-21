const HTTP_BASE_URL = import.meta.env.VITE_API_URL || "";

export const URL_PATHS = {
  SOURCE: {
    sources: "/source/sources",
    byId: (id: string) => `/source/sources/${id}`,
    fetch: (id: string) => `/source/sources/${id}/fetch`,
    locales: (lang: string) => `/source/locales/${lang}`,
    messages: (lang: string) => `/source/messages/${lang}`,
  },
  NEWS: {
    items: "/news/items",
    search: "/news/search",
    preferences: "/news/preferences",
    locales: (lang: string) => `/news/locales/${lang}`,
    messages: (lang: string) => `/news/messages/${lang}`,
  },
  CRAWL: {
    sessions: "/crawl/sessions",
    byId: (id: string) => `/crawl/sessions/${id}`,
    locales: (lang: string) => `/crawl/locales/${lang}`,
    messages: (lang: string) => `/crawl/messages/${lang}`,
  },
  REPORT: {
    keywords: "/report/keywords",
    snapshots: "/report/snapshots",
    byId: (id: string) => `/report/snapshots/${id}`,
    html: (id: string) => `/report/snapshots/${id}/html`,
    locales: (lang: string) => `/report/locales/${lang}`,
    messages: (lang: string) => `/report/messages/${lang}`,
  },
  NOTIFY: {
    kinds: "/notify/kinds",
    channels: "/notify/channels",
    deliveries: "/notify/deliveries",
    dispatch: "/notify/dispatch",
    locales: (lang: string) => `/notify/locales/${lang}`,
    messages: (lang: string) => `/notify/messages/${lang}`,
  },
  MCP: {
    tools: "/mcp/tools",
    run: "/mcp/run",
    byName: (name: string) => `/mcp/tools/${name}`,
    locales: (lang: string) => `/mcp/locales/${lang}`,
    messages: (lang: string) => `/mcp/messages/${lang}`,
  },
} as const;

export const API_ENDPOINTS = {
  PURE: HTTP_BASE_URL,
} as const;
