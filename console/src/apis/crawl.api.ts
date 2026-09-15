import type { DataResponse } from "nfx-ui/types";

import { protectedClient, publicClient } from "./clients";
import { URL_PATHS } from "./ip";

export type CrawlSession = {
  id: string;
  sourceId?: string;
  status: string;
  itemCount: number;
  errorMessage?: string;
  startedAt: string;
  finishedAt?: string;
};

export const TriggerCrawl = async (sourceId?: string): Promise<CrawlSession> => {
  const { data } = await protectedClient.post<DataResponse<CrawlSession>>(URL_PATHS.CRAWL.sessions, { sourceId });
  return data.data;
};

export const ListCrawlSessions = async (limit = 20): Promise<CrawlSession[]> => {
  const { data } = await publicClient.get<DataResponse<CrawlSession[]>>(URL_PATHS.CRAWL.sessions, { params: { limit } });
  return data.data;
};
