export type SourceMeta = {
  id: string;
  name: string;
  title: string;
  column: string;
  home: string;
  color: string;
  intervalMs: number;
  type: string;
  redirect: string;
};

export type NewsItem = {
  id: string;
  sourceId: string;
  originalId: string;
  title: string;
  url: string;
  mobileUrl?: string;
  pubDate?: number;
  extra?: Record<string, unknown>;
};

export type Keyword = {
  id: string;
  groupName: string;
  word: string;
  kind: string;
  countLimit: number;
};

export type Snapshot = {
  id: string;
  mode: string;
  title: string;
  itemCount: number;
  createdAt: string;
  payload?: unknown;
};

export type Channel = {
  id: string;
  kind: string;
  name: string;
  enabled: boolean;
  config?: Record<string, unknown>;
};

export type Delivery = {
  id: string;
  channelId: string;
  reportId?: string | null;
  status: string;
  errorMessage?: string | null;
  createdAt: string;
  sentAt?: string | null;
};

export type CrawlSession = {
  id: string;
  sourceId?: string;
  status: string;
  itemCount: number;
  errorMessage?: string;
  startedAt: string;
  finishedAt?: string;
};
