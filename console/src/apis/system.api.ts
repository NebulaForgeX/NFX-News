import type { DataResponse } from "nfx-ui/types";

import type { SystemState } from "@/types/domain";

import { protectedClient, publicClient } from "./clients";
import { URL_PATHS } from "./ip";

export type NewsModule = "source" | "news" | "crawl" | "report" | "notify" | "mcp" | "system";

function moduleLocales(module: NewsModule) {
  return {
    source: URL_PATHS.SOURCE.locales,
    news: URL_PATHS.NEWS.locales,
    crawl: URL_PATHS.CRAWL.locales,
    report: URL_PATHS.REPORT.locales,
    notify: URL_PATHS.NOTIFY.locales,
    mcp: URL_PATHS.MCP.locales,
    system: URL_PATHS.SYSTEM.locales,
  }[module];
}

function moduleMessages(module: NewsModule) {
  return {
    source: URL_PATHS.SOURCE.messages,
    news: URL_PATHS.NEWS.messages,
    crawl: URL_PATHS.CRAWL.messages,
    report: URL_PATHS.REPORT.messages,
    notify: URL_PATHS.NOTIFY.messages,
    mcp: URL_PATHS.MCP.messages,
    system: URL_PATHS.SYSTEM.messages,
  }[module];
}

export const getErrorTranslations = async (lang: string, module: NewsModule = "system"): Promise<Record<string, unknown>> => {
  const { data } = await publicClient.get<Record<string, unknown>>(moduleLocales(module)(lang));
  return data;
};

export const getMessageTranslations = async (lang: string, module: NewsModule = "system"): Promise<Record<string, unknown>> => {
  const { data } = await publicClient.get<Record<string, unknown>>(moduleMessages(module)(lang));
  return data;
};

export const getLatestSystemState = async (): Promise<SystemState> => {
  const { data } = await protectedClient.get<DataResponse<SystemState>>(URL_PATHS.SYSTEM.latest);
  return data.data;
};

export const initializeSystem = async (version = "1.0.0"): Promise<SystemState> => {
  const { data } = await protectedClient.post<DataResponse<SystemState>>(URL_PATHS.SYSTEM.initialize, { version });
  return data.data;
};
