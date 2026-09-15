import { useState } from "react";
import { Button, Card, Flex, Heading, Select, Text, TextField } from "@radix-ui/themes";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { useKeywords, useSnapshots } from "@/hooks/news";
import { useNewsRepositories } from "@/apis/repositories";

export default function ReportsPage() {
  const { t } = useTranslation("ReportsPage");
  const { data: keywords } = useKeywords();
  const { data: snapshots } = useSnapshots();
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  const [word, setWord] = useState("");
  const [kind, setKind] = useState("include");
  const add = useMutation({
    mutationFn: () => repos.report.AddKeyword({ word, kind }),
    onSuccess: () => {
      setWord("");
      void qc.invalidateQueries({ queryKey: ["keywords"] });
    },
  });
  const generate = useMutation({
    mutationFn: (mode: string) => repos.report.GenerateReport(mode),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["snapshots"] }),
  });

  return (
    <Flex direction="column" gap="4" p="4">
      <Heading size="6">{t("title")}</Heading>
      <Card>
        <Heading size="4">{t("keywords")}</Heading>
        <Flex gap="2" mt="3" wrap="wrap">
          <TextField.Root value={word} onChange={(e) => setWord(e.target.value)} placeholder="+required / !exclude" />
          <Select.Root value={kind} onValueChange={setKind}>
            <Select.Trigger />
            <Select.Content>
              <Select.Item value="include">include</Select.Item>
              <Select.Item value="exclude">exclude</Select.Item>
            </Select.Content>
          </Select.Root>
          <Button onClick={() => add.mutate()} disabled={!word}>
            {t("add")}
          </Button>
        </Flex>
        <Flex direction="column" gap="2" mt="3">
          {(keywords ?? []).map((k) => (
            <Text key={k.id} size="2">
              {k.kind} {k.word}
            </Text>
          ))}
        </Flex>
      </Card>
      <Flex gap="2">
        {["daily", "current", "incremental"].map((mode) => (
          <Button key={mode} variant="soft" onClick={() => generate.mutate(mode)}>
            {t(mode)}
          </Button>
        ))}
      </Flex>
      <Card>
        <Heading size="4">{t("snapshots")}</Heading>
        <Flex direction="column" gap="2" mt="3">
          {(snapshots ?? []).map((s) => (
            <Text key={s.id} size="2">
              {s.title} · {s.itemCount} · {s.createdAt}
            </Text>
          ))}
        </Flex>
      </Card>
    </Flex>
  );
}
