import { FileDescriptionIcon, HashtagIcon } from "nfx-ui/icons";
import { memo, useState } from "react";
import { Button, Card, Flex, Select, Text, TextField } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { PageFrame } from "@/layouts";
import { CardHeader, EmptyState, PageHeader } from "@/components";

import { OpenSnapshotHTML } from "@/apis/report.api";
import { KEYWORD_KINDS, REPORT_MODES } from "@/enums/newsEnum";
import { useKeywords, useSnapshots, useAddKeyword, useGenerateReport, useDispatchReport } from "@/hooks/news";

const ReportsPage = memo(() => {
  const { t } = useTranslation("pages.Reports");
  const { data: keywords } = useKeywords();
  const { data: snapshots } = useSnapshots();
  const [word, setWord] = useState("");
  const [groupName, setGroupName] = useState("default");
  const [kind, setKind] = useState<(typeof KEYWORD_KINDS)[number]>("include");
  const add = useAddKeyword();
  const generate = useGenerateReport();
  const dispatch = useDispatchReport();

  return (
    <PageFrame>
      <PageHeader
        icon={FileDescriptionIcon}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2">
            {REPORT_MODES.map((mode) => (
              <Button key={mode} variant="soft" onClick={() => generate.mutate(mode)}>
                {t(mode)}
              </Button>
            ))}
          </Flex>
        }
      />
      <Flex direction="column" gap="4">
        <Card>
          <CardHeader icon={<HashtagIcon size={18} />} title={t("keywords")} />
          <Flex gap="2" wrap="wrap">
            <TextField.Root value={groupName} onChange={(e) => setGroupName(e.target.value)} placeholder={t("group")} />
            <TextField.Root value={word} onChange={(e) => setWord(e.target.value)} placeholder="+required / !exclude" />
            <Select.Root value={kind} onValueChange={(v) => setKind(v as (typeof KEYWORD_KINDS)[number])}>
              <Select.Trigger />
              <Select.Content>
                {KEYWORD_KINDS.map((k) => (
                  <Select.Item key={k} value={k}>
                    {k}
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Root>
            <Button
              onClick={() =>
                add.mutate(
                  { word, kind, groupName },
                  {
                    onSuccess: () => setWord(""),
                  },
                )
              }
              disabled={!word}
            >
              {t("add")}
            </Button>
          </Flex>
          {(keywords ?? []).length === 0 ? (
            <EmptyState icon={HashtagIcon} title={t("emptyKeywords")} />
          ) : (
            <Flex direction="column" gap="2" mt="3">
              {(keywords ?? []).map((k) => (
                <Text key={k.id} size="2">
                  {k.groupName} · {k.kind} {k.word}
                </Text>
              ))}
            </Flex>
          )}
        </Card>
        <Card>
          <CardHeader icon={<FileDescriptionIcon size={18} />} title={t("snapshots")} />
          {(snapshots ?? []).length === 0 ? (
            <EmptyState icon={FileDescriptionIcon} title={t("emptySnapshots")} action={<Button variant="soft" onClick={() => generate.mutate("current")}>{t("current")}</Button>} />
          ) : (
            <Flex direction="column" gap="2">
              {(snapshots ?? []).map((s) => (
                <Flex key={s.id} align="center" justify="between" gap="2" wrap="wrap">
                  <Text size="2">
                    {s.title} · {s.itemCount} · {s.createdAt}
                  </Text>
                  <Flex gap="2">
                    <Button size="1" variant="soft" onClick={() => void OpenSnapshotHTML(s.id)}>
                      {t("openHtml")}
                    </Button>
                    <Button size="1" variant="soft" onClick={() => dispatch.mutate(s.id)} disabled={dispatch.isPending}>
                      {t("dispatch")}
                    </Button>
                  </Flex>
                </Flex>
              ))}
            </Flex>
          )}
        </Card>
      </Flex>
    </PageFrame>
  );
});

ReportsPage.displayName = "ReportsPage";
export default ReportsPage;
