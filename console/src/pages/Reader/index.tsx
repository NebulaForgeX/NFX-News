import { FileDescriptionIcon, MagnifierIcon, RefreshIcon } from "nfx-ui/icons";
import { memo, useEffect, useMemo, useState } from "react";
import { Button, Flex, Text, TextField } from "@radix-ui/themes";
import { DndContext, closestCenter, type DragEndEvent } from "@dnd-kit/core";
import { SortableContext, arrayMove, horizontalListSortingStrategy, useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useTranslation } from "react-i18next";
import { PageFrame } from "@/layouts";
import { EmptyState, PageHeader } from "@/components";
import { sourceColorVar } from "@/enums/newsEnum";
import {
  useColumnPreferences,
  useNewsItems,
  useSearchNews,
  useSources,
  useFetchSource,
  useSavePreferences,
  type SourceMeta,
} from "@/hooks/news";
import type { NewsItem, ReaderPrefsPayload } from "@/types/domain";
import { formatRelativeTime } from "@/utils";

import styles from "./s.module.css";

function itemTime(item: NewsItem): string {
  if (item.pubDate && item.pubDate > 0) return formatRelativeTime(new Date(item.pubDate).toISOString());
  if (item.updatedAt) return formatRelativeTime(item.updatedAt);
  return "";
}

function SortableColumn({
  source,
  onFetch,
  onHide,
}: {
  source: SourceMeta;
  onFetch: (id: string) => void;
  onHide: (id: string) => void;
}) {
  const { t } = useTranslation("pages.Reader");
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: source.id });
  const { data: items } = useNewsItems(source.id);
  const rows = items ?? [];
  const color = sourceColorVar(source.color);
  return (
    <div ref={setNodeRef} className={styles.column} style={{ transform: CSS.Transform.toString(transform), transition }} {...attributes}>
      <div className={styles.columnHeader}>
        <span className={styles.rail} style={{ background: color }} />
        <div className={styles.headText} {...listeners}>
          <span className={styles.headName}>{source.name || source.id}</span>
          <span className={styles.headMeta}>
            {[source.title || source.type, source.column].filter(Boolean).join(" · ")}
          </span>
        </div>
        <div className={styles.headActions}>
          {source.home ? (
            <Button size="1" variant="ghost" asChild>
              <a href={source.home} target="_blank" rel="noreferrer">
                {t("home")}
              </a>
            </Button>
          ) : null}
          <Button size="1" variant="ghost" onClick={() => onFetch(source.id)}>
            {t("refresh")}
          </Button>
          <Button size="1" variant="ghost" onClick={() => onHide(source.id)}>
            {t("hide")}
          </Button>
        </div>
      </div>
      <div className={styles.list}>
        {rows.length === 0 ? (
          <EmptyState icon={FileDescriptionIcon} title={t("emptyColumn")} />
        ) : (
          rows.map((item, index) => (
            <a
              key={item.id}
              className={styles.item}
              href={item.mobileUrl || item.url}
              target="_blank"
              rel="noreferrer"
              title={item.extra?.hover || item.title}
            >
              <span className={styles.rank}>{String(index + 1).padStart(2, "0")}</span>
              <span className={styles.body}>
                <span className={styles.title}>
                  {item.extra?.icon?.url ? <img className={styles.flag} src={item.extra.icon.url} alt="" /> : null}
                  {item.title}
                </span>
                <span className={styles.meta}>{[item.extra?.info, itemTime(item)].filter(Boolean).join(" · ")}</span>
              </span>
            </a>
          ))
        )}
      </div>
    </div>
  );
}

const ReaderPage = memo(() => {
  const { t } = useTranslation("pages.Reader");
  const { data: sources } = useSources();
  const prefs = useColumnPreferences();
  const [order, setOrder] = useState<string[]>([]);
  const [hidden, setHidden] = useState<string[]>([]);
  const [columnFilter, setColumnFilter] = useState("");
  const [q, setQ] = useState("");
  const search = useSearchNews(q);
  const refresh = useFetchSource();
  const savePrefs = useSavePreferences();

  useEffect(() => {
    const saved = prefs.data;
    if (!saved) return;
    if (saved.columnOrder.length > 0) setOrder(saved.columnOrder);
    setHidden(saved.payload.hiddenSourceIds ?? []);
    setColumnFilter(saved.payload.columnFilter ?? "");
  }, [prefs.data]);

  const persist = (nextOrder: string[], nextHidden: string[], nextFilter: string) => {
    const payload: ReaderPrefsPayload = { hiddenSourceIds: nextHidden, columnFilter: nextFilter };
    void savePrefs.mutateAsync({ columnOrder: nextOrder, payload });
  };

  const catalog = useMemo(() => (sources ?? []).filter((source) => !source.redirect), [sources]);
  const columns = useMemo(() => {
    const set = new Set<string>();
    for (const source of catalog) {
      if (source.column) set.add(source.column);
    }
    return [...set];
  }, [catalog]);

  const sourceById = useMemo(() => new Map(catalog.map((source) => [source.id, source])), [catalog]);

  const visible = useMemo(() => {
    const hiddenSet = new Set(hidden);
    const list = catalog.filter((source) => !hiddenSet.has(source.id));
    const filtered = columnFilter ? list.filter((source) => source.column === columnFilter) : list;
    if (order.length === 0) return filtered;
    const map = new Map(filtered.map((source) => [source.id, source]));
    const ordered = order.map((id) => map.get(id)).filter((source): source is SourceMeta => Boolean(source));
    const rest = filtered.filter((source) => !order.includes(source.id));
    return [...ordered, ...rest];
  }, [catalog, order, hidden, columnFilter]);

  const ids = visible.map((source) => source.id);
  const searchRows = search.data ?? [];

  const onDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const oldIndex = ids.indexOf(String(active.id));
    const newIndex = ids.indexOf(String(over.id));
    const next = arrayMove(ids, oldIndex, newIndex);
    setOrder(next);
    persist(next, hidden, columnFilter);
  };

  const onHide = (id: string) => {
    const nextHidden = hidden.includes(id) ? hidden : [...hidden, id];
    setHidden(nextHidden);
    persist(order, nextHidden, columnFilter);
  };

  const onFilter = (value: string) => {
    setColumnFilter(value);
    persist(order, hidden, value);
  };

  return (
    <PageFrame maxWidth="100%" fullHeight>
      <div className={styles.shell}>
        <PageHeader
          icon={FileDescriptionIcon}
          title={t("title")}
          description={t("subtitle")}
          actions={
            <Flex gap="2" align="center">
              <TextField.Root className={styles.search} value={q} onChange={(e) => setQ(e.target.value)} placeholder={t("search")} />
              <Button
                variant="soft"
                onClick={() => {
                  for (const source of visible) void refresh.mutateAsync(source.id);
                }}
              >
                <RefreshIcon size={14} />
                {t("refreshAll")}
              </Button>
            </Flex>
          }
        />
        <div className={styles.toolbar}>
          <div className={styles.chips}>
            <Button size="1" variant={columnFilter === "" ? "solid" : "soft"} onClick={() => onFilter("")}>
              {t("allColumns")}
            </Button>
            {columns.map((column) => (
              <Button key={column} size="1" variant={columnFilter === column ? "solid" : "soft"} onClick={() => onFilter(column)}>
                {column}
              </Button>
            ))}
            {hidden.length > 0 ? (
              <Button
                size="1"
                variant="ghost"
                onClick={() => {
                  setHidden([]);
                  persist(order, [], columnFilter);
                }}
              >
                {t("showHidden", { count: hidden.length })}
              </Button>
            ) : null}
          </div>
        </div>
        {q ? (
          <div className={styles.searchPanel}>
            <Text size="1" color="gray" mb="2">
              {t("results")}
            </Text>
            {searchRows.length === 0 ? (
              <EmptyState icon={MagnifierIcon} title={t("emptySearch")} />
            ) : (
              searchRows.map((item) => (
                <a key={item.id} className={styles.searchRow} href={item.mobileUrl || item.url} target="_blank" rel="noreferrer">
                  <Text size="1" color="gray" style={{ minWidth: 88 }}>
                    {sourceById.get(item.sourceId)?.name || item.sourceId}
                  </Text>
                  <Text size="2">{item.title}</Text>
                  <Text size="1" color="gray">
                    {itemTime(item)}
                  </Text>
                </a>
              ))
            )}
          </div>
        ) : visible.length === 0 ? (
          <EmptyState icon={FileDescriptionIcon} title={t("emptySources")} description={t("emptySourcesHint")} />
        ) : (
          <DndContext collisionDetection={closestCenter} onDragEnd={onDragEnd}>
            <SortableContext items={ids} strategy={horizontalListSortingStrategy}>
              <div className={styles.board}>
                {visible.map((source) => (
                  <SortableColumn key={source.id} source={source} onFetch={(id) => void refresh.mutateAsync(id)} onHide={onHide} />
                ))}
              </div>
            </SortableContext>
          </DndContext>
        )}
      </div>
    </PageFrame>
  );
});

ReaderPage.displayName = "ReaderPage";
export default ReaderPage;
