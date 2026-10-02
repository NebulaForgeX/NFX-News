import { FileDescriptionIcon, MagnifierIcon, RefreshIcon } from "nfx-ui/icons";
import { memo, useEffect, useMemo, useState } from "react";
import { Container, Box, Section, Button, Flex, Text, TextField } from "@radix-ui/themes";
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
    <Flex ref={setNodeRef} direction="column" className={`${styles.columnSize} ${styles.columnEdge}`} style={{ transform: CSS.Transform.toString(transform), transition }} {...attributes}>
      <Flex direction="column" minHeight="0" className={styles.column}>
      <Flex direction="column" flexShrink="0" className={styles.columnHeader}>
        <Container className={styles.columnHeaderPx}>
          <Section className={styles.columnHeaderPy}>
            <Flex align="start" gap="2">
              <Section className={`${styles.railItem} ${styles.railOffset}`}>
                <Box className={`${styles.railSize} ${styles.railRadius}`} style={{ background: color }} />
              </Section>
              <Flex direction="column" gap="1" minWidth="0" flexGrow="1" className={styles.headText} {...listeners}>
                <span className={styles.headName}>{source.name || source.id}</span>
                <Section className={styles.headMetaSpace}>
                  <span className={styles.headMeta}>
                    {[source.title || source.type, source.column].filter(Boolean).join(" · ")}
                  </span>
                </Section>
              </Flex>
              <Flex gap="1" flexShrink="0">
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
              </Flex>
            </Flex>
          </Section>
        </Container>
      </Flex>
      <Flex direction="column" minHeight="0" className={styles.list}>
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
              <Container className={styles.itemPx}>
                <Section className={styles.itemPy}>
                  <Flex gap="2">
                    <Flex className={styles.rank}>
                      <span>{String(index + 1).padStart(2, "0")}</span>
                    </Flex>
                    <Flex direction="column" gap="1" minWidth="0" className={styles.body}>
                      <Flex align="center" gap="1">
                        {item.extra?.icon?.url ? (
                          <span className={styles.flagSize}>
                            <img className={styles.flagImage} src={item.extra.icon.url} alt="" />
                          </span>
                        ) : null}
                        <span className={styles.title}>{item.title}</span>
                      </Flex>
                      <Section className={styles.metaSpace}>
                        <span className={styles.meta}>{[item.extra?.info, itemTime(item)].filter(Boolean).join(" · ")}</span>
                      </Section>
                    </Flex>
                  </Flex>
                </Section>
              </Container>
            </a>
          ))
        )}
      </Flex>
      </Flex>
    </Flex>
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
      <Flex direction="column" gap="3" className={styles.shell}>
        <PageHeader
          icon={FileDescriptionIcon}
          title={t("title")}
          description={t("subtitle")}
          actions={
            <Flex gap="2" align="center">
              <TextField.Root className={styles.search} value={q} onChange={(e) => setQ(e.target.value)} placeholder={t("search")} />
              <Button
                variant="outline"
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
        <Flex wrap="wrap" align="center" gap="2">
          <Flex wrap="wrap" gap="2" flexGrow="1">
            <Button size="1" variant={columnFilter === "" ? "solid" : "outline"} onClick={() => onFilter("")}>
              {t("allColumns")}
            </Button>
            {columns.map((column) => (
              <Button key={column} size="1" variant={columnFilter === column ? "solid" : "outline"} onClick={() => onFilter(column)}>
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
          </Flex>
        </Flex>
        {q ? (
          <Box className={styles.searchPanelHairline}>
            <Section className={styles.searchPanelPy}>
              <Box className={styles.searchPanel}>
            <Section mb="2"><Text size="1" color="gray">
              {t("results")}
            </Text></Section>
            {searchRows.length === 0 ? (
              <EmptyState icon={MagnifierIcon} title={t("emptySearch")} />
            ) : (
              searchRows.map((item) => (
                <a key={item.id} className={styles.searchRow} href={item.mobileUrl || item.url} target="_blank" rel="noreferrer">
                  <Section className={styles.searchRowPy}>
                    <Flex gap="2" align="center">
                      <Text size="1" color="gray" style={{ minWidth: 88 }}>
                        {sourceById.get(item.sourceId)?.name || item.sourceId}
                      </Text>
                      <Text size="2">{item.title}</Text>
                      <Text size="1" color="gray">
                        {itemTime(item)}
                      </Text>
                    </Flex>
                  </Section>
                </a>
              ))
            )}
              </Box>
            </Section>
          </Box>
        ) : visible.length === 0 ? (
          <EmptyState icon={FileDescriptionIcon} title={t("emptySources")} description={t("emptySourcesHint")} />
        ) : (
          <DndContext collisionDetection={closestCenter} onDragEnd={onDragEnd}>
            <SortableContext items={ids} strategy={horizontalListSortingStrategy}>
              <Flex className={styles.board}>
                {visible.map((source) => (
                  <SortableColumn key={source.id} source={source} onFetch={(id) => void refresh.mutateAsync(id)} onHide={onHide} />
                ))}
              </Flex>
            </SortableContext>
          </DndContext>
        )}
      </Flex>
    </PageFrame>
  );
});

ReaderPage.displayName = "ReaderPage";
export default ReaderPage;
