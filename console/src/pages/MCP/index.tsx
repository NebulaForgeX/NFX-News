import { SparklesIcon } from "nfx-ui/icons";
import { memo, useMemo, useState } from "react";
import { Button, Flex, Link, Text, TextArea, TextField } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { DataTable, PageHeader, SectionBlock } from "@/components";
import { PageFrame } from "@/layouts";
import { MCP_TOOL_MAP } from "@/enums/newsEnum";
import { useMCPTools, useRunMCPTool } from "@/hooks/news";

function asRecord(value: unknown): Record<string, unknown> | null {
  if (value && typeof value === "object" && !Array.isArray(value)) return value as Record<string, unknown>;
  return null;
}

function asList(value: unknown): Record<string, unknown>[] {
  if (!Array.isArray(value)) return [];
  return value.map((item) => (typeof item === "object" && item ? (item as Record<string, unknown>) : { value: item }));
}

function titleOf(row: Record<string, unknown>): string {
  return String(row.title ?? row.topic ?? row.name ?? row.tool ?? row.id ?? "");
}

function ResultView({ result }: { result: unknown }) {
  const { t } = useTranslation("pages.MCP");
  const record = asRecord(result);
  if (!record) {
    return (
      <Text size="1" style={{ whiteSpace: "pre-wrap", fontFamily: "ui-monospace, monospace" }}>
        {String(result)}
      </Text>
    );
  }
  const items = asList(record.items ?? record.topics ?? record.sessions);
  if (items.length > 0) {
    const hasUrl = items.some((row) => typeof row.url === "string");
    const hasCount = items.some((row) => row.count != null);
    return (
      <DataTable
        empty={t("emptyResult")}
        rows={items}
        rowKey={(row) => String(row.id ?? row.topic ?? row.url ?? titleOf(row))}
        columns={[
          {
            key: "title",
            header: t("resultTitle"),
            render: (row) =>
              hasUrl && typeof row.url === "string" ? (
                <Link href={row.url} target="_blank" rel="noreferrer">
                  {titleOf(row) || row.url}
                </Link>
              ) : (
                <Text size="2">{titleOf(row) || "—"}</Text>
              ),
          },
          ...(hasCount
            ? [{ key: "count", header: t("count"), render: (row: Record<string, unknown>) => String(row.count ?? "") }]
            : [{ key: "sourceId", header: t("source"), render: (row: Record<string, unknown>) => String(row.sourceId ?? row.source_id ?? "—") }]),
        ]}
      />
    );
  }
  return (
    <Text size="1" style={{ whiteSpace: "pre-wrap", fontFamily: "ui-monospace, monospace" }}>
      {JSON.stringify(result, null, 2)}
    </Text>
  );
}

const MCPPage = memo(() => {
  const { t } = useTranslation("pages.MCP");
  const { data: tools } = useMCPTools();
  const run = useRunMCPTool();
  const [name, setName] = useState("get_latest_news");
  const [values, setValues] = useState<Record<string, string>>({});
  const [rawJson, setRawJson] = useState("");
  const [result, setResult] = useState<unknown>(null);
  const listed = tools && tools.length > 0 ? tools : Object.keys(MCP_TOOL_MAP);
  const def = MCP_TOOL_MAP[name];

  const builtArgs = useMemo(() => {
    const args: Record<string, unknown> = {};
    for (const field of def?.args ?? []) {
      const raw = values[field.key]?.trim();
      if (!raw) continue;
      args[field.key] = field.kind === "number" ? Number(raw) : raw;
    }
    return args;
  }, [def, values]);

  const execute = () => {
    let parsed: Record<string, unknown> = builtArgs;
    if (rawJson.trim()) {
      try {
        parsed = JSON.parse(rawJson) as Record<string, unknown>;
      } catch {
        setResult(t("invalidJson"));
        return;
      }
    }
    void run
      .mutateAsync({ name, arguments: parsed })
      .then((data) => setResult(data))
      .catch((err: unknown) => setResult(String(err)));
  };

  return (
    <PageFrame>
      <PageHeader icon={SparklesIcon} title={t("title")} description={t("subtitle")} />
      <Flex direction="column" gap="6">
        <SectionBlock title={t("tools")}>
          <Flex gap="2" wrap="wrap">
            {listed.map((tool) => (
              <Button
                key={tool}
                size="1"
                variant={name === tool ? "solid" : "soft"}
                onClick={() => {
                  setName(tool);
                  setValues({});
                  setRawJson("");
                  setResult(null);
                }}
              >
                {tool}
              </Button>
            ))}
          </Flex>
        </SectionBlock>
        <SectionBlock title={t("run")} description={t(`desc.${name}`, { defaultValue: name })}>
          <Flex direction="column" gap="2">
            {(def?.args ?? []).map((field) => (
              <TextField.Root
                key={field.key}
                type={field.kind === "number" ? "number" : "text"}
                value={values[field.key] ?? ""}
                onChange={(e) => setValues((prev) => ({ ...prev, [field.key]: e.target.value }))}
                placeholder={t(`args.${field.key}`, { defaultValue: field.key })}
              />
            ))}
            <TextArea value={rawJson} onChange={(e) => setRawJson(e.target.value)} rows={4} placeholder={t("rawJson")} />
            <Button disabled={!name || run.isPending} onClick={execute}>
              {t("execute")}
            </Button>
            {result != null ? <ResultView result={result} /> : null}
          </Flex>
        </SectionBlock>
      </Flex>
    </PageFrame>
  );
});

MCPPage.displayName = "MCPPage";
export default MCPPage;
