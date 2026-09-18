import type { DataResponse } from "nfx-ui/types";

import type { Keyword, Snapshot } from "@/types/domain";

import { protectedClient } from "./clients";
import { URL_PATHS } from "./ip";

export type { Keyword, Snapshot };

export const ListKeywords = async (): Promise<Keyword[]> => {
  const { data } = await protectedClient.get<DataResponse<Keyword[]>>(URL_PATHS.REPORT.keywords);
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
  const { data } = await protectedClient.get<DataResponse<Snapshot[]>>(URL_PATHS.REPORT.snapshots, { params: { limit } });
  return data.data;
};

export const GetSnapshot = async (id: string): Promise<Snapshot> => {
  const { data } = await protectedClient.get<DataResponse<Snapshot>>(URL_PATHS.REPORT.byId(id));
  return data.data;
};

export const OpenSnapshotHTML = async (id: string): Promise<void> => {
  const { data } = await protectedClient.get<string>(URL_PATHS.REPORT.html(id), { responseType: "text" });
  const blob = new Blob([data], { type: "text/html;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  window.open(url, "_blank", "noopener,noreferrer");
};
