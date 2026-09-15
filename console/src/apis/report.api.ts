import type { DataResponse } from "nfx-ui/types";

import { protectedClient, publicClient } from "./clients";
import { URL_PATHS } from "./ip";

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

export const ListKeywords = async (): Promise<Keyword[]> => {
  const { data } = await publicClient.get<DataResponse<Keyword[]>>(URL_PATHS.REPORT.keywords);
  return data.data;
};

export const AddKeyword = async (params: { groupName?: string; word: string; kind?: string; countLimit?: number }): Promise<Keyword> => {
  const { data } = await protectedClient.post<DataResponse<Keyword>>(URL_PATHS.REPORT.keywords, params);
  return data.data;
};

export const GenerateReport = async (mode: string): Promise<Snapshot> => {
  const { data } = await protectedClient.post<DataResponse<Snapshot>>(URL_PATHS.REPORT.snapshots, { mode });
  return data.data;
};

export const ListSnapshots = async (limit = 20): Promise<Snapshot[]> => {
  const { data } = await publicClient.get<DataResponse<Snapshot[]>>(URL_PATHS.REPORT.snapshots, { params: { limit } });
  return data.data;
};
