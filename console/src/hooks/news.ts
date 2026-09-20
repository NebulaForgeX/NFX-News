import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { useNewsRepositories } from "@/apis/repositories";
import { NEWS_QUERY_KEYS } from "@/constants";
import type { ReaderPreferences, ReaderPrefsPayload, SourceMeta } from "@/types/domain";

export type { SourceMeta };

function asStringArray(value: unknown): string[] {
  if (typeof value === "string") {
    try {
      const parsed: unknown = JSON.parse(value);
      if (Array.isArray(parsed)) return parsed.filter((item): item is string => typeof item === "string");
    } catch {
      return [];
    }
  }
  if (!Array.isArray(value)) return [];
  return value.filter((item): item is string => typeof item === "string");
}

function normalizePreferences(raw: ReaderPreferences | undefined): ReaderPreferences {
  const payload = raw?.payload ?? {};
  return {
    columnOrder: asStringArray(raw?.columnOrder),
    payload: {
      hiddenSourceIds: asStringArray(payload.hiddenSourceIds),
      columnFilter: typeof payload.columnFilter === "string" ? payload.columnFilter : "",
    },
  };
}

export function useSources() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.sources, queryFn: repos.news.ListSources });
}

export function useSource(id: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: NEWS_QUERY_KEYS.source(id),
    queryFn: () => repos.news.GetSource(id),
    enabled: id.trim().length > 0,
  });
}

export function useNotifyKinds() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.notifyKinds, queryFn: repos.notify.ListNotifyKinds });
}

export function useNewsItems(sourceId?: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: NEWS_QUERY_KEYS.items(sourceId),
    queryFn: () => repos.news.ListNews({ sourceId, limit: 40 }),
    enabled: Boolean(sourceId),
  });
}

export function useSearchNews(q: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: NEWS_QUERY_KEYS.search(q),
    queryFn: () => repos.news.SearchNews(q, 80),
    enabled: q.trim().length > 0,
  });
}

export function useKeywords() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.keywords, queryFn: repos.report.ListKeywords });
}

export function useSnapshots() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.snapshots, queryFn: () => repos.report.ListSnapshots(40) });
}

export function useSnapshot(id: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: NEWS_QUERY_KEYS.snapshot(id),
    queryFn: () => repos.report.GetSnapshot(id),
    enabled: id.trim().length > 0,
  });
}

export function useCrawlSessions() {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: NEWS_QUERY_KEYS.crawlSessions,
    queryFn: () => repos.crawl.ListCrawlSessions(40),
    refetchInterval: (query) => ((query.state.data ?? []).some((row) => row.status === "running") ? 2000 : false),
  });
}

export function useCrawlSession(id: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: NEWS_QUERY_KEYS.crawlSession(id),
    queryFn: () => repos.crawl.GetCrawlSession(id),
    enabled: id.trim().length > 0,
    refetchInterval: (query) => (query.state.data?.status === "running" ? 2000 : false),
  });
}

export function useChannels() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.notifyChannels, queryFn: repos.notify.ListChannels });
}

export function useDeliveries() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.notifyDeliveries, queryFn: () => repos.notify.ListDeliveries(50) });
}

export function useColumnPreferences() {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: NEWS_QUERY_KEYS.preferences,
    queryFn: async () => normalizePreferences(await repos.news.GetPreferences()),
  });
}

export function useFetchSource() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => repos.news.FetchSource(id),
    onSuccess: (items, id) => {
      qc.setQueryData(NEWS_QUERY_KEYS.items(id), items);
    },
  });
}

export function useSavePreferences() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (params: { columnOrder?: string[]; payload?: ReaderPrefsPayload }) => repos.news.SetPreferences(params),
    onSuccess: () => qc.invalidateQueries({ queryKey: NEWS_QUERY_KEYS.preferences }),
  });
}

export function useTriggerCrawl() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (sourceId?: string) => repos.crawl.TriggerCrawl(sourceId),
    onSuccess: () => qc.invalidateQueries({ queryKey: NEWS_QUERY_KEYS.crawlSessions }),
  });
}

export function useAddKeyword() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { word: string; kind: string; groupName?: string; countLimit?: number }) =>
      repos.report.AddKeyword(body),
    onSuccess: () => qc.invalidateQueries({ queryKey: NEWS_QUERY_KEYS.keywords }),
  });
}

function snapshotPayloadJson(payload: unknown): string {
  if (!payload || typeof payload !== "object") return "";
  const raw = payload as {
    mode?: string;
    items?: Array<{ id?: string; title?: string; url?: string; sourceId?: string; group?: string; isNew?: boolean }>;
  };
  return JSON.stringify({
    mode: raw.mode,
    items: (raw.items ?? []).map((item) => ({
      id: item.id,
      title: item.title,
      url: item.url,
      source_id: item.sourceId,
      group: item.group,
      is_new: item.isNew,
    })),
  });
}

export function useOpenSnapshotHTML() {
  const repos = useNewsRepositories();
  return useMutation({
    mutationFn: (id: string) => repos.report.OpenSnapshotHTML(id),
  });
}

export function useDispatchReport() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const snap = await repos.report.GetSnapshot(id);
      return repos.notify.DispatchReport({
        reportId: snap.id,
        mode: snap.mode,
        title: snap.title,
        itemCount: snap.itemCount,
        payloadJson: snapshotPayloadJson(snap.payload),
      });
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: NEWS_QUERY_KEYS.notifyDeliveries }),
  });
}

export function useMCPTools() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.mcpTools, queryFn: repos.mcp.ListMCPTools });
}

export function useRunMCPTool() {
  const repos = useNewsRepositories();
  return useMutation({
    mutationFn: (body: { name: string; arguments?: Record<string, unknown> }) =>
      repos.mcp.RunMCPTool(body.name, body.arguments),
  });
}

export function useSystemState() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: NEWS_QUERY_KEYS.systemState, queryFn: repos.system.getLatestSystemState });
}

export function useInitializeSystem() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (version?: string) => repos.system.initializeSystem(version),
    onSuccess: () => qc.invalidateQueries({ queryKey: NEWS_QUERY_KEYS.systemState }),
  });
}

export function useGenerateReport() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (mode: string) => repos.report.GenerateReport(mode),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: NEWS_QUERY_KEYS.snapshots });
    },
  });
}

export function useUpsertChannel() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { kind: string; name: string; enabled: boolean; config: Record<string, unknown> }) =>
      repos.notify.UpsertChannel(body),
    onSuccess: () => qc.invalidateQueries({ queryKey: NEWS_QUERY_KEYS.notifyChannels }),
  });
}
