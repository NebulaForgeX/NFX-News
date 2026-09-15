import { useQuery } from "@tanstack/react-query";

import { useNewsRepositories } from "@/apis/repositories";

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
