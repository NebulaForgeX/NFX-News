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

export type NewsExtraIcon = {
  url?: string;
  scale?: number;
};

export type NewsExtra = {
  info?: string;
  hover?: string;
  date?: string;
  icon?: NewsExtraIcon;
};

export type NewsItem = {
  id: string;
  sourceId: string;
  originalId: string;
  title: string;
  url: string;
  mobileUrl?: string;
  pubDate?: number;
  extra?: NewsExtra;
  updatedAt?: string;
};

export type ReaderPrefsPayload = {
  hiddenSourceIds?: string[];
  columnFilter?: string;
};

export type ReaderPreferences = {
  columnOrder: string[];
  payload: ReaderPrefsPayload;
};

export type Keyword = {
  id: string;
  groupName: string;
  word: string;
  kind: string;
  countLimit: number;
  createdAt?: string;
};

export type SnapshotItem = {
  id: string;
  title: string;
  url: string;
  sourceId: string;
  group: string;
  isNew: boolean;
};

export type SnapshotPayload = {
  mode?: string;
  items?: SnapshotItem[];
  html?: string;
};

export type Snapshot = {
  id: string;
  mode: string;
  title: string;
  itemCount: number;
  createdAt: string;
  payload?: SnapshotPayload;
};

export type Channel = {
  id: string;
  kind: string;
  name: string;
  enabled: boolean;
  config?: Record<string, unknown>;
  createdAt?: string;
  updatedAt?: string;
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
