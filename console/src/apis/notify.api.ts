import type { DataResponse } from "nfx-ui/types";

import { protectedClient } from "./clients";
import { URL_PATHS } from "./ip";

export type Channel = {
  id: string;
  kind: string;
  name: string;
  enabled: boolean;
  config?: Record<string, unknown>;
};

export const ListChannels = async (): Promise<Channel[]> => {
  const { data } = await protectedClient.get<DataResponse<Channel[]>>(URL_PATHS.NOTIFY.channels);
  return data.data;
};

export const UpsertChannel = async (params: { kind: string; name: string; enabled: boolean; config?: Record<string, unknown> }): Promise<Channel> => {
  const { data } = await protectedClient.post<DataResponse<Channel>>(URL_PATHS.NOTIFY.channels, params);
  return data.data;
};

export const DispatchReport = async (params: { reportId: string; mode?: string; title?: string }): Promise<{ queued: number }> => {
  const { data } = await protectedClient.post<DataResponse<{ queued: number }>>(URL_PATHS.NOTIFY.dispatch, params);
  return data.data;
};

export type Delivery = {
  id: string;
  channel_id: string;
  report_id?: string | null;
  status: string;
  error_message?: string | null;
  created_at: string;
  sent_at?: string | null;
};

export const ListDeliveries = async (limit = 50): Promise<Delivery[]> => {
  const { data } = await protectedClient.get<DataResponse<Delivery[]>>(URL_PATHS.NOTIFY.deliveries, { params: { limit } });
  return data.data;
};
