import { memo, useState } from "react";
import { Button, Card, Flex, Text, TextField } from "@radix-ui/themes";
import { Play, List } from "lucide-react";
import { useTranslation } from "react-i18next";
import { PageFrame } from "nfx-ui/layouts";
import { CardHeader, EmptyState, PageHeader } from "nfx-ui/components";

import { useCrawlSessions, useSources, useTriggerCrawl } from "@/hooks/news";

const CrawlPage = memo(() => {
  const { t } = useTranslation("CrawlPage");
  const { data: sessions } = useCrawlSessions();
  const { data: sources } = useSources();
  const [sourceId, setSourceId] = useState("");
  const trigger = useTriggerCrawl();

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
          </Flex>
        }
      />
      <Card>
        <CardHeader icon={<List size={18} />} title={t("sessions")} />
        {(sessions ?? []).length === 0 ? (
          <EmptyState icon={Play} title={t("emptySessions")} description={t("emptySessionsHint")} />
        ) : (
          <Flex direction="column" gap="2">
            {(sessions ?? []).map((s) => (
              <Text key={s.id} size="2">
                {s.status} · {s.sourceId || "*"} · {s.itemCount} · {s.startedAt}
              </Text>
            ))}
          </Flex>
        )}
      </Card>
    </PageFrame>
  );
});

CrawlPage.displayName = "CrawlPage";
export default CrawlPage;
