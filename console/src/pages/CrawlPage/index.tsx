import { Button, Card, Flex, Heading, Text, TextField } from "@radix-ui/themes";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { useCrawlSessions, useSources } from "@/hooks/news";
import { useNewsRepositories } from "@/apis/repositories";

export default function CrawlPage() {
  const { t } = useTranslation("CrawlPage");
  const { data: sessions } = useCrawlSessions();
  const { data: sources } = useSources();
  const [sourceId, setSourceId] = useState("");
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  const trigger = useMutation({
    mutationFn: () => repos.crawl.TriggerCrawl(sourceId || undefined),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["crawl-sessions"] }),
  });

  return (
    <Flex direction="column" gap="4" p="4">
      <Heading size="6">{t("title")}</Heading>
      <Flex gap="2">
        <TextField.Root value={sourceId} onChange={(e) => setSourceId(e.target.value)} placeholder={t("sourceId")} list="source-ids" />
        <datalist id="source-ids">
          {(sources ?? []).map((s) => (
            <option key={s.id} value={s.id} />
          ))}
        </datalist>
        <Button onClick={() => trigger.mutate()}>{t("trigger")}</Button>
      </Flex>
      <Card>
        <Flex direction="column" gap="2">
          {(sessions ?? []).map((s) => (
            <Text key={s.id} size="2">
              {s.status} · {s.sourceId || "*"} · {s.itemCount} · {s.startedAt}
            </Text>
          ))}
        </Flex>
      </Card>
    </Flex>
  );
}
