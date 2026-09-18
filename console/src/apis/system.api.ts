import type { DataResponse } from "nfx-ui/types";

import { protectedClient, publicClient } from "./clients";
import { URL_PATHS } from "./ip";

export type NewsModule = "source" | "news" | "crawl" | "report" | "notify" | "mcp" | "system";

export const getErrorTranslations = async (lang: string, module: NewsModule = "system"): Promise<Record<string, unknown>> => {
  const path = {
    source: URL_PATHS.SOURCE.i18nErrors,
    news: URL_PATHS.NEWS.i18nErrors,
    crawl: URL_PATHS.CRAWL.i18nErrors,
    report: URL_PATHS.REPORT.i18nErrors,
    notify: URL_PATHS.NOTIFY.i18nErrors,
    mcp: URL_PATHS.MCP.i18nErrors,
    system: URL_PATHS.SYSTEM.i18nErrors,
  }[module](lang);
  const { data } = await publicClient.get<Record<string, unknown>>(path);
  return data;
};

export const getLatestSystemState = async () => {
  const { data } = await protectedClient.get<DataResponse<{ initialized: boolean }>>(URL_PATHS.SYSTEM.latest);
  return data.data;
};

export const initializeSystem = async (version = "1.0.0") => {
  const { data } = await protectedClient.post<DataResponse<{ initialized: boolean }>>(URL_PATHS.SYSTEM.initialize, { version });
  return data.data;
};
