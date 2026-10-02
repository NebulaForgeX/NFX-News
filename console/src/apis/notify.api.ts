import type { DataResponse } from "nfx-ui/types";

import type { Channel, Delivery } from "@/types/domain";

import { protectedClient } from "./clients";
import { URL_PATHS } from "./ip";

export type { Channel, Delivery };

export const ListNotifyKinds = async (): Promise<string[]> => {
  const { data } = await protectedClient.get<DataResponse<string[]>>(URL_PATHS.NOTIFY.kinds);
  return data.data;
};

export const ListChannels = async (): Promise<Channel[]> => {
  const { data } = await protectedClient.get<DataResponse<Channel[]>>(URL_PATHS.NOTIFY.channels);
  return data.data;
};

export const UpsertChannel = async (params: { kind: string; name: string; enabled: boolean; config?: Record<string, unknown> }): Promise<Channel> => {
  const { data } = await protectedClient.post<DataResponse<Channel>>(URL_PATHS.NOTIFY.channels, params);
  return data.data;
};

export const DispatchReport = async (params: {
  reportId: string;
  mode?: string;
  title?: string;
  itemCount?: number;
  payloadJson?: string;
}): Promise<{ queued: number }> => {
  const { data } = await protectedClient.post<DataResponse<{ queued: number }>>(URL_PATHS.NOTIFY.dispatch, {
    report_id: params.reportId,
    mode: params.mode,
    title: params.title,
    item_count: params.itemCount,
    payload_json: params.payloadJson,
  });
  return data.data;
};

export const ListDeliveries = async (limit = 50): Promise<Delivery[]> => {
  const { data } = await protectedClient.get<DataResponse<Delivery[]>>(URL_PATHS.NOTIFY.deliveries, { params: { limit } });
  return data.data;
};
