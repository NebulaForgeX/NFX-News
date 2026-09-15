import { memo, useState } from "react";
import { Button, Card, Flex, Select, Text, TextField } from "@radix-ui/themes";
import { FileText, Tag } from "lucide-react";
import { useTranslation } from "react-i18next";
import { PageFrame } from "nfx-ui/layouts";
import { CardHeader, EmptyState, PageHeader } from "nfx-ui/components";

import { useKeywords, useSnapshots, useAddKeyword, useGenerateReport } from "@/hooks/news";

const ReportsPage = memo(() => {
  const { t } = useTranslation("ReportsPage");
  const { data: keywords } = useKeywords();
  const { data: snapshots } = useSnapshots();
  const [word, setWord] = useState("");
  const [kind, setKind] = useState("include");
  const add = useAddKeyword();
  const generate = useGenerateReport();

  return (
    <PageFrame>
      <PageHeader
        icon={FileText}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2">
            {["daily", "current", "incremental"].map((mode) => (
              <Button key={mode} variant="soft" onClick={() => generate.mutate(mode)}>
                {t(mode)}
              </Button>
            ))}
          </Flex>
        }
      />
      <Flex direction="column" gap="4">
        <Card>
          <CardHeader icon={<Tag size={18} />} title={t("keywords")} />
          <Flex gap="2" wrap="wrap">
            <TextField.Root value={word} onChange={(e) => setWord(e.target.value)} placeholder="+required / !exclude" />
            <Select.Root value={kind} onValueChange={setKind}>
              <Select.Trigger />
              <Select.Content>
                <Select.Item value="include">include</Select.Item>
                <Select.Item value="exclude">exclude</Select.Item>
              </Select.Content>
            </Select.Root>
            <Button
              onClick={() =>
                add.mutate(
                  { word, kind },
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
            <EmptyState icon={Tag} title={t("emptyKeywords")} />
          ) : (
            <Flex direction="column" gap="2" mt="3">
              {(keywords ?? []).map((k) => (
                <Text key={k.id} size="2">
                  {k.kind} {k.word}
                </Text>
              ))}
            </Flex>
          )}
        </Card>
        <Card>
          <CardHeader icon={<FileText size={18} />} title={t("snapshots")} />
          {(snapshots ?? []).length === 0 ? (
            <EmptyState icon={FileText} title={t("emptySnapshots")} action={<Button variant="soft" onClick={() => generate.mutate("current")}>{t("current")}</Button>} />
          ) : (
            <Flex direction="column" gap="2">
              {(snapshots ?? []).map((s) => (
                <Text key={s.id} size="2">
                  {s.title} · {s.itemCount} · {s.createdAt}
                </Text>
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
