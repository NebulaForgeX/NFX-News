import { useEffect, useMemo, useState } from "react";
import { Badge, Box, Button, Card, Flex, Heading, Text, TextField } from "@radix-ui/themes";
import { DndContext, closestCenter, type DragEndEvent } from "@dnd-kit/core";
import { SortableContext, arrayMove, horizontalListSortingStrategy, useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useTranslation } from "react-i18next";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import { useColumnPreferences, useNewsItems, useSearchNews, useSources } from "@/hooks/news";
import { useNewsRepositories } from "@/apis/repositories";
import type { SourceMeta } from "@/apis/news.api";

function SortableColumn({ source }: { source: SourceMeta }) {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: source.id });
  const { data: items } = useNewsItems(source.id);
  const style = { transform: CSS.Transform.toString(transform), transition, minWidth: 280, maxWidth: 320 };
  return (
    <Box ref={setNodeRef} style={style} {...attributes}>
      <Card>
        <Flex direction="column" gap="3">
          <Flex justify="between" align="center" {...listeners} style={{ cursor: "grab" }}>
            <Heading size="3">{source.name || source.id}</Heading>
            <Badge style={{ background: source.color || undefined }}>{source.column || source.type}</Badge>
          </Flex>
          {(items ?? []).slice(0, 12).map((item) => (
            <a key={item.id} href={item.url} target="_blank" rel="noreferrer">
              <Text size="2">{item.title}</Text>
            </a>
          ))}
        </Flex>
      </Card>
    </Box>
  );
}

export default function ReaderPage() {
  const { t } = useTranslation("ReaderPage");
  const { data: sources } = useSources();
  const prefs = useColumnPreferences();
  const [order, setOrder] = useState<string[]>([]);
  const [q, setQ] = useState("");
  const search = useSearchNews(q);
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  const refresh = useMutation({
    mutationFn: async (id: string) => repos.news.FetchSource(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["news"] }),
  });
  const saveOrder = useMutation({
    mutationFn: (columnOrder: string[]) => repos.news.SetPreferences({ columnOrder }),
  });

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
    <Flex direction="column" gap="4" p="4">
      <Flex justify="between" align="center" wrap="wrap" gap="3">
        <Heading size="6">{t("title")}</Heading>
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
      </Flex>
      {q && (
        <Card>
          <Heading size="3">{t("results")}</Heading>
          <Flex direction="column" gap="2" mt="2">
            {(search.data ?? []).map((item) => (
              <a key={item.id} href={item.url} target="_blank" rel="noreferrer">
                {item.title}
              </a>
            ))}
          </Flex>
        </Card>
      )}
      <DndContext collisionDetection={closestCenter} onDragEnd={onDragEnd}>
        <SortableContext items={ids} strategy={horizontalListSortingStrategy}>
          <Flex gap="3" overflowX="auto" pb="4">
            {visible.map((source) => (
              <SortableColumn key={source.id} source={source} />
            ))}
          </Flex>
        </SortableContext>
      </DndContext>
    </Flex>
  );
}
