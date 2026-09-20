import { SparklesIcon } from "nfx-ui/icons";
import { memo, useState } from "react";
import { Button, Card, Flex, Text, TextArea, TextField } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { PageFrame } from "@/layouts";
import { CardHeader, EmptyState, PageHeader } from "@/components";

import { useMCPTools, useRunMCPTool } from "@/hooks/news";

const MCPPage = memo(() => {
  const { t } = useTranslation("pages.MCP");
  const { data: tools } = useMCPTools();
  const run = useRunMCPTool();
  const [name, setName] = useState("");
  const [args, setArgs] = useState("{}");
  const [result, setResult] = useState("");

  return (
    <PageFrame>
      <PageHeader icon={SparklesIcon} title={t("title")} description={t("subtitle")} />
      <Flex direction="column" gap="4">
        <Card>
          <CardHeader icon={<SparklesIcon size={18} />} title={t("tools")} />
          {(tools ?? []).length === 0 ? (
            <EmptyState icon={SparklesIcon} title={t("emptyTools")} />
          ) : (
            <Flex gap="2" wrap="wrap" mt="3">
              {(tools ?? []).map((tool) => (
                <Button
                  key={tool}
                  size="1"
                  variant={name === tool ? "solid" : "soft"}
                  onClick={() => setName(tool)}
                >
                  {tool}
                </Button>
              ))}
            </Flex>
          )}
        </Card>
        <Card>
          <CardHeader icon={<SparklesIcon size={18} />} title={t("run")} />
          <Flex direction="column" gap="2">
            <TextField.Root value={name} onChange={(e) => setName(e.target.value)} placeholder="get_latest_news" />
            <TextArea value={args} onChange={(e) => setArgs(e.target.value)} rows={6} />
            <Button
              disabled={!name || run.isPending}
              onClick={() => {
                let parsed: Record<string, unknown> = {};
                try {
                  parsed = args.trim() ? (JSON.parse(args) as Record<string, unknown>) : {};
                } catch {
                  setResult(t("invalidJson"));
                  return;
                }
                void run.mutateAsync({ name, arguments: parsed }).then((data) => {
                  setResult(JSON.stringify(data, null, 2));
                }).catch((err: unknown) => {
                  setResult(String(err));
                });
              }}
            >
              {t("execute")}
            </Button>
            {result ? (
              <Text size="1" style={{ whiteSpace: "pre-wrap", fontFamily: "ui-monospace, monospace" }}>
                {result}
              </Text>
            ) : null}
          </Flex>
        </Card>
      </Flex>
    </PageFrame>
  );
});

MCPPage.displayName = "MCPPage";
export default MCPPage;
