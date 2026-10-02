export type Tokens = {
  accessToken: string;
  refreshToken?: string;
};

export type {
  SourceMeta,
  NewsItem,
  Keyword,
  Snapshot,
  Channel,
  Delivery,
  CrawlSession,
} from "./domain";
