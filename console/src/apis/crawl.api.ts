import type { DataResponse } from "nfx-ui/types";

import type { CrawlSession } from "@/types/domain";

import { protectedClient } from "./clients";
import { URL_PATHS } from "./ip";

export type { CrawlSession };

export const TriggerCrawl = async (sourceId?: string): Promise<CrawlSession> => {
  const { data } = await protectedClient.post<DataResponse<CrawlSession>>(URL_PATHS.CRAWL.sessions, { sourceId: sourceId ?? "" });
  return data.data;
};

export const ListCrawlSessions = async (limit = 20): Promise<CrawlSession[]> => {
  const { data } = await protectedClient.get<DataResponse<CrawlSession[]>>(URL_PATHS.CRAWL.sessions, { params: { limit } });
  return data.data;
};

export const GetCrawlSession = async (id: string): Promise<CrawlSession> => {
  const { data } = await protectedClient.get<DataResponse<CrawlSession>>(URL_PATHS.CRAWL.byId(id));
  return data.data;
};
