import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { useNewsRepositories } from "@/apis/repositories";
import type { SourceMeta } from "@/apis/news.api";

export type { SourceMeta };

export function useSources() {
  const repos = useNewsRepositories();
  return useQuery({ queryKey: ["sources"], queryFn: repos.news.ListSources });
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
    onSuccess: () => qc.invalidateQueries({ queryKey: ["news"] }),
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
    mutationFn: (body: { word: string; kind: string }) => repos.report.AddKeyword(body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["keywords"] }),
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
