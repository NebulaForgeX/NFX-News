import type { DataResponse } from "nfx-ui/types";

import { protectedClient, publicClient } from "./clients";
import { URL_PATHS } from "./ip";

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

export const ListSources = async (): Promise<SourceMeta[]> => {
  const { data } = await publicClient.get<DataResponse<SourceMeta[]>>(URL_PATHS.SOURCE.sources);
  return data.data;
};

export const GetSource = async (id: string): Promise<SourceMeta> => {
  const { data } = await publicClient.get<DataResponse<SourceMeta>>(URL_PATHS.SOURCE.byId(id));
  return data.data;
};

export const FetchSource = async (id: string): Promise<NewsItem[]> => {
  const { data } = await publicClient.post<DataResponse<NewsItem[]>>(URL_PATHS.SOURCE.fetch(id));
  return data.data;
};

export const ListNews = async (params?: { sourceId?: string; limit?: number }): Promise<NewsItem[]> => {
  const { data } = await publicClient.get<DataResponse<NewsItem[]>>(URL_PATHS.NEWS.items, { params });
  return data.data;
};

export const SearchNews = async (q: string, limit = 50): Promise<NewsItem[]> => {
  const { data } = await publicClient.get<DataResponse<NewsItem[]>>(URL_PATHS.NEWS.search, { params: { q, limit } });
  return data.data;
};

export const GetPreferences = async () => {
  const { data } = await protectedClient.get<DataResponse<{ columnOrder: string[]; payload: Record<string, unknown> }>>(
    URL_PATHS.NEWS.preferences,
  );
  return data.data;
};

export const SetPreferences = async (params: { columnOrder?: string[]; payload?: Record<string, unknown> }) => {
  const { data } = await protectedClient.put<DataResponse<unknown>>(URL_PATHS.NEWS.preferences, params);
  return data.data;
};
