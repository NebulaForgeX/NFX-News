import { memo, useEffect, useMemo, useState } from "react";
import { Box, Button, Card, Flex, Text, TextField } from "@radix-ui/themes";
import { DndContext, closestCenter, type DragEndEvent } from "@dnd-kit/core";
import { SortableContext, arrayMove, horizontalListSortingStrategy, useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Newspaper, Search } from "lucide-react";
import { useTranslation } from "react-i18next";
import { PageFrame } from "nfx-ui/layouts";
import { CardHeader, EmptyState, PageHeader } from "nfx-ui/components";

import { useColumnPreferences, useNewsItems, useSearchNews, useSources, useFetchSource, useSaveColumnOrder, type SourceMeta } from "@/hooks/news";

function SortableColumn({ source }: { source: SourceMeta }) {
  const { t } = useTranslation("ReaderPage");
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: source.id });
  const { data: items } = useNewsItems(source.id);
  const style = { transform: CSS.Transform.toString(transform), transition, minWidth: 280, maxWidth: 320 };
  const rows = (items ?? []).slice(0, 12);
  return (
    <Box ref={setNodeRef} style={style} {...attributes}>
      <Card>
        <Flex direction="column" gap="3">
          <Box {...listeners} style={{ cursor: "grab" }}>
            <CardHeader icon={<Newspaper size={18} />} title={source.name || source.id} description={source.column || source.type} />
          </Box>
          {rows.length === 0 ? (
            <EmptyState icon={Newspaper} title={t("emptyColumn")} />
          ) : (
            rows.map((item) => (
              <a key={item.id} href={item.url} target="_blank" rel="noreferrer">
                <Text size="2">{item.title}</Text>
              </a>
            ))
          )}
        </Flex>
      </Card>
    </Box>
  );
}

const ReaderPage = memo(() => {
  const { t } = useTranslation("ReaderPage");
  const { data: sources } = useSources();
  const prefs = useColumnPreferences();
  const [order, setOrder] = useState<string[]>([]);
  const [q, setQ] = useState("");
  const search = useSearchNews(q);
  const refresh = useFetchSource();
  const saveOrder = useSaveColumnOrder();

  useEffect(() => {
    const saved = prefs.data?.columnOrder;
    if (saved && saved.length > 0) setOrder(saved);
  }, [prefs.data]);

  const visible = useMemo(() => {
    const list = sources ?? [];
    if (order.length === 0) return list.slice(0, 8);
    const map = new Map(list.map((s) => [s.id, s]));
    return order.map((id) => map.get(id)).filter((s): s is SourceMeta => Boolean(s));
  }, [sources, order]);

  const ids = visible.map((s) => s.id);
  const searchRows = search.data ?? [];

  const onDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const oldIndex = ids.indexOf(String(active.id));
    const newIndex = ids.indexOf(String(over.id));
    const next = arrayMove(ids, oldIndex, newIndex);
    setOrder(next);
    void saveOrder.mutateAsync(next);
  };

  return (
    <PageFrame>
      <PageHeader
        icon={Newspaper}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2">
            <TextField.Root value={q} onChange={(e) => setQ(e.target.value)} placeholder={t("search")} />
            <Button
              variant="soft"
              onClick={() => {
                if (visible[0]) void refresh.mutateAsync(visible[0].id);
              }}
            >
              {t("refresh")}
            </Button>
          </Flex>
        }
      />
      {q ? (
        <Card mb="4">
          <CardHeader icon={<Search size={18} />} title={t("results")} />
          {searchRows.length === 0 ? (
            <EmptyState icon={Search} title={t("emptySearch")} />
          ) : (
            <Flex direction="column" gap="2">
              {searchRows.map((item) => (
                <a key={item.id} href={item.url} target="_blank" rel="noreferrer">
                  {item.title}
                </a>
              ))}
            </Flex>
          )}
        </Card>
      ) : null}
      {visible.length === 0 ? (
        <EmptyState icon={Newspaper} title={t("emptySources")} description={t("emptySourcesHint")} />
      ) : (
        <DndContext collisionDetection={closestCenter} onDragEnd={onDragEnd}>
          <SortableContext items={ids} strategy={horizontalListSortingStrategy}>
            <Flex gap="3" overflowX="auto" pb="4">
              {visible.map((source) => (
                <SortableColumn key={source.id} source={source} />
              ))}
            </Flex>
          </SortableContext>
        </DndContext>
      )}
    </PageFrame>
  );
});

ReaderPage.displayName = "ReaderPage";
export default ReaderPage;
