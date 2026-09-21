import { FileDescriptionIcon, HashtagIcon } from "nfx-ui/icons";
import { memo, useMemo, useState } from "react";
import { Badge, Box, Button, Flex, Select, Text, TextField } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import { DataTable, EmptyState, PageHeader, SectionBlock } from "@/components";
import { PageFrame } from "@/layouts";
import { ROUTES } from "@/navigations";
import { KEYWORD_KINDS, REPORT_MODES } from "@/enums/newsEnum";
import { useAddKeyword, useDispatchReport, useGenerateReport, useKeywords, useOpenSnapshotHTML, useSnapshots } from "@/hooks/news";
import type { Keyword, Snapshot } from "@/types/domain";
import { formatDateTime } from "@/utils";

const ReportsPage = memo(() => {
  const { t } = useTranslation("pages.Reports");
  const navigate = useNavigate();
  const { data: keywords, isLoading: keywordsLoading } = useKeywords();
  const { data: snapshots, isLoading: snapshotsLoading } = useSnapshots();
  const [word, setWord] = useState("");
  const [groupName, setGroupName] = useState("default");
  const [kind, setKind] = useState<(typeof KEYWORD_KINDS)[number]>("include");
  const [countLimit, setCountLimit] = useState("10");
  const add = useAddKeyword();
  const generate = useGenerateReport();
  const dispatch = useDispatchReport();
  const openHtml = useOpenSnapshotHTML();

  const grouped = useMemo(() => {
    const map = new Map<string, Keyword[]>();
    for (const keyword of keywords ?? []) {
      const key = keyword.groupName || "default";
      const list = map.get(key) ?? [];
      list.push(keyword);
      map.set(key, list);
    }
    return [...map.entries()];
  }, [keywords]);

  return (
    <PageFrame>
      <PageHeader
        icon={FileDescriptionIcon}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2">
            {REPORT_MODES.map((mode) => (
              <Button key={mode} variant="outline" onClick={() => generate.mutate(mode)} disabled={generate.isPending}>
                {t(mode)}
              </Button>
            ))}
          </Flex>
        }
      />
      <Flex direction="column" gap="6">
        <SectionBlock title={t("keywords")}>
          <Box pb="3">
          <Flex gap="2" wrap="wrap">
            <TextField.Root value={groupName} onChange={(e) => setGroupName(e.target.value)} placeholder={t("group")} />
            <TextField.Root value={word} onChange={(e) => setWord(e.target.value)} placeholder={t("wordHint")} />
            <Select.Root value={kind} onValueChange={(v) => setKind(v as (typeof KEYWORD_KINDS)[number])}>
              <Select.Trigger />
              <Select.Content>
                {KEYWORD_KINDS.map((item) => (
                  <Select.Item key={item} value={item}>
                    {t(`kind.${item}`)}
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Root>
            <TextField.Root
              type="number"
              value={countLimit}
              onChange={(e) => setCountLimit(e.target.value)}
              placeholder={t("countLimit")}
              style={{ width: 96 }}
            />
            <Button
              onClick={() =>
                add.mutate(
                  { word, kind, groupName, countLimit: Number(countLimit) || 0 },
                  { onSuccess: () => setWord("") },
                )
              }
              disabled={!word || add.isPending}
            >
              {t("add")}
            </Button>
          </Flex>
          </Box>
          {keywordsLoading || grouped.length === 0 ? (
            <EmptyState icon={HashtagIcon} title={t("emptyKeywords")} />
          ) : (
            grouped.map(([group, rows]) => (
              <Box key={group} pb="4">
              <Flex direction="column" gap="2">
                <Text size="2" weight="medium">
                  {group}
                </Text>
                <DataTable
                  empty={t("emptyKeywords")}
                  rows={rows}
                  rowKey={(row) => row.id}
                  columns={[
                    { key: "word", header: t("word") },
                    {
                      key: "kind",
                      header: t("kindLabel"),
                      render: (row) => (
                        <Badge variant="outline" color={row.kind === "exclude" ? "red" : row.kind === "required" ? "amber" : "green"}>
                          {t(`kind.${row.kind}`, { defaultValue: row.kind })}
                        </Badge>
                      ),
                    },
                    { key: "countLimit", header: t("countLimit") },
                    {
                      key: "createdAt",
                      header: t("createdAt"),
                      render: (row) => (row.createdAt ? formatDateTime(row.createdAt) : "—"),
                    },
                    ]}
                  />
                </Flex>
                </Box>
              ))
            )}
        </SectionBlock>
        <SectionBlock title={t("snapshots")}>
          <DataTable
            loading={snapshotsLoading}
            empty={t("emptySnapshots")}
            rows={snapshots ?? []}
            rowKey={(row: Snapshot) => row.id}
            columns={[
              {
                key: "title",
                header: t("titleCol"),
                render: (row) => (
                  <Button variant="ghost" onClick={() => navigate(`${ROUTES.REPORTS}/${row.id}`)}>
                    {row.title}
                  </Button>
                ),
              },
              {
                key: "mode",
                header: t("mode"),
                render: (row) => <Badge variant="outline">{t(row.mode, { defaultValue: row.mode })}</Badge>,
              },
              { key: "itemCount", header: t("itemCount") },
              {
                key: "createdAt",
                header: t("createdAt"),
                render: (row) => formatDateTime(row.createdAt),
              },
              {
                key: "actions",
                header: t("actions"),
                render: (row) => (
                  <Flex gap="2">
                    <Button size="1" variant="outline" onClick={() => void openHtml.mutateAsync(row.id)}>
                      {t("openHtml")}
                    </Button>
                    <Button size="1" variant="outline" onClick={() => dispatch.mutate(row.id)} disabled={dispatch.isPending}>
                      {t("dispatch")}
                    </Button>
                  </Flex>
                ),
              },
            ]}
          />
        </SectionBlock>
      </Flex>
    </PageFrame>
  );
});

ReportsPage.displayName = "ReportsPage";
export default ReportsPage;
