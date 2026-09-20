import { PlayerIcon } from "nfx-ui/icons";
import { memo, useMemo, useState } from "react";
import { Badge, Button, Flex, Select, Text } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { DataTable, PageHeader, SectionBlock } from "@/components";
import { PageFrame } from "@/layouts";
import { useCrawlSession, useCrawlSessions, useSources, useTriggerCrawl } from "@/hooks/news";
import type { CrawlSession } from "@/types/domain";
import { formatDateTime } from "@/utils";

function statusColor(status: string): "green" | "red" | "orange" | "gray" {
  if (status === "ok") return "green";
  if (status === "failed") return "red";
  if (status === "running") return "orange";
  return "gray";
}

const CrawlPage = memo(() => {
  const { t } = useTranslation("pages.Crawl");
  const { data: sessions, isLoading } = useCrawlSessions();
  const { data: sources } = useSources();
  const [sourceId, setSourceId] = useState("all");
  const [selectedId, setSelectedId] = useState("");
  const trigger = useTriggerCrawl();
  const detail = useCrawlSession(selectedId);
  const crawlable = useMemo(() => (sources ?? []).filter((source) => !source.redirect), [sources]);
  const sourceName = (id?: string) => crawlable.find((source) => source.id === id)?.name || id || t("allSources");

  const onTrigger = () => {
    void trigger.mutateAsync(sourceId === "all" ? undefined : sourceId).then((session) => {
      if (session?.id) setSelectedId(session.id);
    });
  };

  const selected = detail.data ?? (sessions ?? []).find((row) => row.id === selectedId);

  return (
    <PageFrame>
      <PageHeader
        icon={PlayerIcon}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2" wrap="wrap">
            <Select.Root value={sourceId} onValueChange={setSourceId}>
              <Select.Trigger placeholder={t("sourceId")} />
              <Select.Content>
                <Select.Item value="all">{t("allSources")}</Select.Item>
                {crawlable.map((source) => (
                  <Select.Item key={source.id} value={source.id}>
                    {source.name || source.id}
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Root>
            <Button onClick={onTrigger} disabled={trigger.isPending}>
              {sourceId === "all" ? t("triggerAll") : t("trigger")}
            </Button>
          </Flex>
        }
      />
      <Flex direction="column" gap="6">
        <SectionBlock title={t("sessions")}>
          <DataTable
            loading={isLoading}
            empty={t("emptySessions")}
            rows={sessions ?? []}
            rowKey={(row: CrawlSession) => row.id}
            onRowClick={(row) => setSelectedId(row.id)}
            columns={[
              {
                key: "status",
                header: t("status"),
                render: (row) => (
                  <Badge color={statusColor(row.status)} variant="soft">
                    {t(`statusValue.${row.status}`, { defaultValue: row.status })}
                  </Badge>
                ),
              },
              {
                key: "sourceId",
                header: t("source"),
                render: (row) => sourceName(row.sourceId),
              },
              { key: "itemCount", header: t("itemCount") },
              {
                key: "startedAt",
                header: t("startedAt"),
                render: (row) => formatDateTime(row.startedAt),
              },
              {
                key: "finishedAt",
                header: t("finishedAt"),
                render: (row) => (row.finishedAt ? formatDateTime(row.finishedAt) : "—"),
              },
              {
                key: "errorMessage",
                header: t("error"),
                render: (row) => (
                  <Text size="1" color={row.errorMessage ? "red" : "gray"}>
                    {row.errorMessage || "—"}
                  </Text>
                ),
              },
            ]}
          />
        </SectionBlock>
        {selected ? (
          <SectionBlock title={t("detail")}>
            <Flex direction="column" gap="2">
              <Text size="2">
                {t("status")}: {selected.status}
              </Text>
              <Text size="2">
                {t("source")}: {sourceName(selected.sourceId)}
              </Text>
              <Text size="2">
                {t("itemCount")}: {selected.itemCount}
              </Text>
              <Text size="2">
                {t("startedAt")}: {formatDateTime(selected.startedAt)}
              </Text>
              <Text size="2">
                {t("finishedAt")}: {selected.finishedAt ? formatDateTime(selected.finishedAt) : "—"}
              </Text>
              <Text size="2" color={selected.errorMessage ? "red" : "gray"}>
                {t("error")}: {selected.errorMessage || "—"}
              </Text>
            </Flex>
          </SectionBlock>
        ) : null}
      </Flex>
    </PageFrame>
  );
});

CrawlPage.displayName = "CrawlPage";
export default CrawlPage;
