import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { useNewsRepositories } from "@/apis/repositories";
import type { SourceMeta } from "@/apis/news.api";

export type { SourceMeta };

export function useSources() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["sources"], queryFn: repos.news.ListSources });
}

export function useSource(id: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: ["source", id],
    queryFn: () => repos.news.GetSource(id),
    enabled: id.trim().length > 0,
  });
}

export function useNotifyKinds() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["notify-kinds"], queryFn: repos.notify.ListNotifyKinds });
}

export function useNewsItems(sourceId?: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: ["news", sourceId],
    queryFn: () => repos.news.ListNews({ sourceId, limit: 40 }),
  });
}

export function useSearchNews(q: string) {
  const repos = useNewsRepositories();
  return useQuery({
    queryKey: ["news-search", q],
    queryFn: () => repos.news.SearchNews(q),
    enabled: q.trim().length > 0,
  });
}

export function useKeywords() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["keywords"], queryFn: repos.report.ListKeywords });
}

export function useSnapshots() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["snapshots"], queryFn: () => repos.report.ListSnapshots(20) });
}

export function useCrawlSessions() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["crawl-sessions"], queryFn: () => repos.crawl.ListCrawlSessions(30) });
}

export function useChannels() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["notify-channels"], queryFn: repos.notify.ListChannels });
}

export function useDeliveries() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["notify-deliveries"], queryFn: () => repos.notify.ListDeliveries(30) });
}

export function useColumnPreferences() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["news-preferences"], queryFn: repos.news.GetPreferences });
}

export function useFetchSource() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => repos.news.FetchSource(id),
    onSuccess: (items, id) => {
      qc.setQueryData(["news", id], items);
    },
  });
}

export function useSaveColumnOrder() {
  const repos = useNewsRepositories();
  return useMutation({
    mutationFn: (columnOrder: string[]) => repos.news.SetPreferences({ columnOrder }),
  });
}

export function useTriggerCrawl() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (sourceId?: string) => repos.crawl.TriggerCrawl(sourceId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["crawl-sessions"] }),
  });
}

export function useAddKeyword() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { word: string; kind: string; groupName?: string }) => repos.report.AddKeyword(body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["keywords"] }),
  });
}

export function useDispatchReport() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const snap = await repos.report.GetSnapshot(id);
      const payload = snap.payload ? JSON.stringify(snap.payload) : "";
      return repos.notify.DispatchReport({
        reportId: snap.id,
        mode: snap.mode,
        title: snap.title,
        itemCount: snap.itemCount,
        payloadJson: payload,
      });
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notify-deliveries"] }),
  });
}

export function useMCPTools() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["mcp-tools"], queryFn: repos.mcp.ListMCPTools });
}

export function useRunMCPTool() {
  const repos = useNewsRepositories();
  return useMutation({
    mutationFn: (body: { name: string; arguments?: Record<string, unknown> }) => repos.mcp.RunMCPTool(body.name, body.arguments),
  });
}

export function useSystemState() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["system-state"], queryFn: repos.system.getLatestSystemState });
}

export function useInitializeSystem() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => repos.system.initializeSystem(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["system-state"] }),
  });
}

export function useGetCrawlSession() {
  const repos = useNewsRepositories();
  return useMutation({
    mutationFn: (id: string) => repos.crawl.GetCrawlSession(id),
  });
}

export function useGenerateReport() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (mode: string) => repos.report.GenerateReport(mode),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["snapshots"] }),
  });
}

export function useUpsertChannel() {
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { kind: string; name: string; enabled: boolean; config: Record<string, unknown> }) =>
      repos.notify.UpsertChannel(body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notify-channels"] }),
  });
}
