import { memo, useState } from "react";
import { Button, Card, Flex, Text, TextField } from "@radix-ui/themes";
import { Play, List } from "lucide-react";
import { useTranslation } from "react-i18next";
import { PageFrame } from "nfx-ui/layouts";
import { CardHeader, EmptyState, PageHeader } from "nfx-ui/components";

import { useCrawlSessions, useGetCrawlSession, useSource, useSources, useTriggerCrawl } from "@/hooks/news";

const CrawlPage = memo(() => {
  const { t } = useTranslation("CrawlPage");
  const { data: sessions } = useCrawlSessions();
  const { data: sources } = useSources();
  const [sourceId, setSourceId] = useState("");
  const trigger = useTriggerCrawl();
  const detail = useGetCrawlSession();
  const meta = useSource(sourceId);

  return (
    <PageFrame>
      <PageHeader
        icon={Play}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2">
            <TextField.Root value={sourceId} onChange={(e) => setSourceId(e.target.value)} placeholder={t("sourceId")} list="source-ids" />
            <datalist id="source-ids">
              {(sources ?? []).map((s) => (
                <option key={s.id} value={s.id} />
              ))}
            </datalist>
            <Button onClick={() => trigger.mutate(sourceId || undefined)}>{t("trigger")}</Button>
            <Button variant="soft" onClick={() => trigger.mutate(undefined)}>
              {t("triggerAll")}
            </Button>
          </Flex>
        }
      />
      {meta.data ? (
        <Text size="2">
          {meta.data.name} · {meta.data.home} · {meta.data.intervalMs}ms
        </Text>
      ) : null}
      <Card>
        <CardHeader icon={<List size={18} />} title={t("sessions")} />
        {(sessions ?? []).length === 0 ? (
          <EmptyState icon={Play} title={t("emptySessions")} description={t("emptySessionsHint")} />
        ) : (
          <Flex direction="column" gap="2">
            {(sessions ?? []).map((s) => (
              <Flex key={s.id} align="center" justify="between" gap="2">
                <Text size="2">
                  {s.status} · {s.sourceId || "*"} · {s.itemCount} · {s.startedAt}
                </Text>
                <Button size="1" variant="ghost" onClick={() => detail.mutate(s.id)}>
                  {t("detail")}
                </Button>
              </Flex>
            ))}
            {detail.data ? (
              <Text size="2">
                {detail.data.id} · {detail.data.status} · {detail.data.errorMessage || "-"}
              </Text>
            ) : null}
          </Flex>
        )}
      </Card>
    </PageFrame>
  );
});

CrawlPage.displayName = "CrawlPage";
export default CrawlPage;
