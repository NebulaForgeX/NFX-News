import { RouterIcon } from "nfx-ui/icons";
import { memo, useState } from "react";
import { Badge, Box, Button, Flex, Text, TextField } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { DataTable, PageHeader, SectionBlock } from "@/components";
import { PageFrame } from "@/layouts";
import { useInitializeSystem, useSystemState } from "@/hooks/news";
import { formatDateTime } from "@/utils";

const SystemPage = memo(() => {
  const { t } = useTranslation("pages.System");
  const { data, isLoading } = useSystemState();
  const init = useInitializeSystem();
  const [version, setVersion] = useState("1.0.0");

  const rows = data
    ? [
        { key: "initialized", label: t("initialized"), value: data.initialized ? t("yes") : t("no") },
        { key: "version", label: t("version"), value: data.initializationVersion || "—" },
        { key: "initializedAt", label: t("initializedAt"), value: data.initializedAt ? formatDateTime(data.initializedAt) : "—" },
        { key: "resetCount", label: t("resetCount"), value: String(data.resetCount ?? 0) },
        { key: "updatedAt", label: t("updatedAt"), value: data.updatedAt ? formatDateTime(data.updatedAt) : "—" },
      ]
    : [];

  return (
    <PageFrame>
      <PageHeader
        icon={RouterIcon}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <Flex gap="2" align="center">
            <TextField.Root value={version} onChange={(e) => setVersion(e.target.value)} placeholder={t("version")} />
            <Button onClick={() => init.mutate(version)} disabled={init.isPending}>
              {t("initialize")}
            </Button>
          </Flex>
        }
      />
      <SectionBlock title={t("state")} description={t("stateHint")}>
        <Box pb="3">
          {data?.initialized ? (
            <Badge color="green" variant="outline">
              {t("ready")}
            </Badge>
          ) : (
            <Badge color="orange" variant="outline">
              {t("pending")}
            </Badge>
          )}
        </Box>
        <DataTable
          loading={isLoading}
          empty={t("empty")}
          rows={rows}
          rowKey={(row) => row.key}
          columns={[
            { key: "label", header: t("field"), render: (row) => <Text size="2">{row.label}</Text> },
            { key: "value", header: t("value"), render: (row) => <Text size="2">{row.value}</Text> },
          ]}
        />
      </SectionBlock>
    </PageFrame>
  );
});

SystemPage.displayName = "SystemPage";
export default SystemPage;
