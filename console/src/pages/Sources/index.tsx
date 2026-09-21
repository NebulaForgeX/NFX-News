import { GlobeIcon, RefreshIcon } from "nfx-ui/icons";
import { memo, useMemo, useState } from "react";
import { Badge, Box, Button, Flex, Link, Text, TextField } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { DataTable, PageHeader } from "@/components";
import { PageFrame } from "@/layouts";
import { sourceColorVar } from "@/enums/newsEnum";
import { useFetchSource, useSources } from "@/hooks/news";
import type { SourceMeta } from "@/types/domain";

const SourcesPage = memo(() => {
  const { t } = useTranslation("pages.Sources");
  const { data, isLoading } = useSources();
  const fetchSource = useFetchSource();
  const [q, setQ] = useState("");
  const [column, setColumn] = useState("");

  const columns = useMemo(() => {
    const set = new Set<string>();
    for (const source of data ?? []) {
      if (source.column) set.add(source.column);
    }
    return [...set];
  }, [data]);

  const rows = useMemo(() => {
    const needle = q.trim().toLowerCase();
    return (data ?? []).filter((source) => {
      if (column && source.column !== column) return false;
      if (!needle) return true;
      return [source.id, source.name, source.title, source.type].some((value) => value?.toLowerCase().includes(needle));
    });
  }, [data, q, column]);

  return (
    <PageFrame>
      <PageHeader
        icon={GlobeIcon}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2" wrap="wrap">
            <TextField.Root value={q} onChange={(e) => setQ(e.target.value)} placeholder={t("search")} />
          </Flex>
        }
      />
      <Box pb="3">
        <Flex gap="2" wrap="wrap">
          <Button size="1" variant={column === "" ? "solid" : "outline"} onClick={() => setColumn("")}>
            {t("allColumns")}
          </Button>
          {columns.map((item) => (
            <Button key={item} size="1" variant={column === item ? "solid" : "outline"} onClick={() => setColumn(item)}>
              {item}
            </Button>
          ))}
        </Flex>
      </Box>
      <DataTable
        loading={isLoading}
        empty={t("empty")}
        rows={rows}
        rowKey={(row) => row.id}
        columns={[
          {
            key: "name",
            header: t("name"),
            render: (row: SourceMeta) => (
              <Flex align="center" gap="2">
                <span style={{ width: 8, height: 8, borderRadius: 99, background: sourceColorVar(row.color), flexShrink: 0 }} />
                <Text size="2" weight="medium">
                  {row.name || row.id}
                </Text>
              </Flex>
            ),
          },
          { key: "id", header: t("id") },
          { key: "title", header: t("feedTitle") },
          { key: "column", header: t("column") },
          { key: "type", header: t("type") },
          {
            key: "intervalMs",
            header: t("interval"),
            render: (row) => (row.intervalMs ? `${Math.round(row.intervalMs / 1000)}s` : "—"),
          },
          {
            key: "redirect",
            header: t("redirect"),
            render: (row) =>
              row.redirect ? (
                <Badge color="gray" variant="outline">
                  {row.redirect}
                </Badge>
              ) : (
                "—"
              ),
          },
          {
            key: "home",
            header: t("home"),
            render: (row) =>
              row.home ? (
                <Link href={row.home} target="_blank" rel="noreferrer">
                  {t("open")}
                </Link>
              ) : (
                "—"
              ),
          },
          {
            key: "fetch",
            header: t("fetch"),
            render: (row) => (
              <Button
                size="1"
                variant="outline"
                disabled={Boolean(row.redirect) || fetchSource.isPending}
                onClick={() => void fetchSource.mutateAsync(row.redirect || row.id)}
              >
                <RefreshIcon size={12} />
                {t("fetch")}
              </Button>
            ),
          },
        ]}
      />
    </PageFrame>
  );
});

SourcesPage.displayName = "SourcesPage";
export default SourcesPage;
