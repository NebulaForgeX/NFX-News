import { FilledBellIcon } from "nfx-ui/icons";
import { memo, useMemo, useState } from "react";
import { Badge, Button, Flex, Select, Switch, Text, TextField } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { DataTable, PageHeader, SectionBlock } from "@/components";
import { PageFrame } from "@/layouts";
import { NOTIFY_KINDS, notifyConfigFields, type NotifyKind } from "@/enums/newsEnum";
import { useChannels, useDeliveries, useNotifyKinds, useUpsertChannel } from "@/hooks/news";
import type { Channel, Delivery } from "@/types/domain";
import { formatDateTime } from "@/utils";

function emptyConfig(kind: NotifyKind): Record<string, unknown> {
  const next: Record<string, unknown> = {};
  for (const field of notifyConfigFields(kind)) {
    next[field.key] = field.kind === "boolean" ? false : field.kind === "number" ? "" : "";
  }
  return next;
}

function asRecord(value: unknown): Record<string, unknown> {
  if (value && typeof value === "object" && !Array.isArray(value)) return value as Record<string, unknown>;
  return {};
}

const NotifyPage = memo(() => {
  const { t } = useTranslation("pages.Notify");
  const { data: channels, isLoading: channelsLoading } = useChannels();
  const { data: deliveries, isLoading: deliveriesLoading } = useDeliveries();
  const { data: kinds } = useNotifyKinds();
  const kindOptions = (kinds && kinds.length > 0 ? kinds : [...NOTIFY_KINDS]) as NotifyKind[];
  const [kind, setKind] = useState<NotifyKind>("feishu");
  const [name, setName] = useState("");
  const [enabled, setEnabled] = useState(true);
  const [config, setConfig] = useState<Record<string, unknown>>(() => emptyConfig("feishu"));
  const save = useUpsertChannel();
  const fields = useMemo(() => notifyConfigFields(kind), [kind]);
  const channelName = (id: string) => (channels ?? []).find((row) => row.id === id)?.name || id;

  const onKindChange = (next: NotifyKind) => {
    setKind(next);
    setConfig(emptyConfig(next));
  };

  const loadChannel = (row: Channel) => {
    const nextKind = (NOTIFY_KINDS as readonly string[]).includes(row.kind) ? (row.kind as NotifyKind) : "feishu";
    setKind(nextKind);
    setName(row.name);
    setEnabled(row.enabled);
    setConfig({ ...emptyConfig(nextKind), ...asRecord(row.config) });
  };

  const setField = (key: string, value: unknown) => {
    setConfig((prev) => ({ ...prev, [key]: value }));
  };

  const builtConfig = (): Record<string, unknown> => {
    const out: Record<string, unknown> = {};
    for (const field of fields) {
      const value = config[field.key];
      if (field.kind === "boolean") {
        out[field.key] = Boolean(value);
        continue;
      }
      if (field.kind === "number") {
        const n = typeof value === "number" ? value : Number(value);
        if (Number.isFinite(n) && String(value).trim() !== "") out[field.key] = n;
        continue;
      }
      if (typeof value === "string" && value.trim()) out[field.key] = value.trim();
    }
    return out;
  };

  return (
    <PageFrame>
      <PageHeader icon={FilledBellIcon} title={t("title")} description={t("webhookHint")} />
      <Flex direction="column" gap="6">
        <SectionBlock title={t("channels")}>
          <Flex direction="column" gap="3" mb="4">
            <Flex gap="2" wrap="wrap" align="center">
              <Select.Root value={kind} onValueChange={(v) => onKindChange(v as NotifyKind)}>
                <Select.Trigger />
                <Select.Content>
                  {kindOptions.map((item) => (
                    <Select.Item key={item} value={item}>
                      {item}
                    </Select.Item>
                  ))}
                </Select.Content>
              </Select.Root>
              <TextField.Root value={name} onChange={(e) => setName(e.target.value)} placeholder={t("name")} />
              <Flex align="center" gap="2">
                <Switch checked={enabled} onCheckedChange={setEnabled} />
                <Text size="2">{t("enabled")}</Text>
              </Flex>
              <Button
                onClick={() => save.mutate({ kind, name, enabled, config: builtConfig() })}
                disabled={!name || save.isPending}
              >
                {t("save")}
              </Button>
            </Flex>
            <Flex gap="2" wrap="wrap">
              {fields.map((field) => {
                const value = config[field.key];
                if (field.kind === "boolean") {
                  return (
                    <Flex key={field.key} align="center" gap="2">
                      <Switch checked={Boolean(value)} onCheckedChange={(checked) => setField(field.key, checked)} />
                      <Text size="2">{t(`fields.${field.key}`)}</Text>
                    </Flex>
                  );
                }
                return (
                  <TextField.Root
                    key={field.key}
                    type={field.secret ? "password" : field.kind === "number" ? "number" : "text"}
                    value={value == null ? "" : String(value)}
                    onChange={(e) => setField(field.key, e.target.value)}
                    placeholder={t(`fields.${field.key}`)}
                  />
                );
              })}
            </Flex>
          </Flex>
          <DataTable
            loading={channelsLoading}
            empty={t("emptyChannels")}
            rows={channels ?? []}
            rowKey={(row: Channel) => row.id}
            onRowClick={loadChannel}
            columns={[
              { key: "kind", header: t("kind") },
              { key: "name", header: t("name") },
              {
                key: "enabled",
                header: t("enabled"),
                render: (row) => (
                  <Badge color={row.enabled ? "green" : "gray"} variant="soft">
                    {row.enabled ? t("on") : t("off")}
                  </Badge>
                ),
              },
              {
                key: "updatedAt",
                header: t("updatedAt"),
                render: (row) => (row.updatedAt ? formatDateTime(row.updatedAt) : "—"),
              },
            ]}
          />
        </SectionBlock>
        <SectionBlock title={t("deliveries")}>
          <DataTable
            loading={deliveriesLoading}
            empty={t("emptyDeliveries")}
            rows={deliveries ?? []}
            rowKey={(row: Delivery) => row.id}
            columns={[
              {
                key: "status",
                header: t("status"),
                render: (row) => (
                  <Badge color={row.status === "sent" ? "green" : row.status === "failed" ? "red" : "orange"} variant="soft">
                    {row.status}
                  </Badge>
                ),
              },
              {
                key: "channelId",
                header: t("channel"),
                render: (row) => channelName(row.channelId),
              },
              { key: "reportId", header: t("report"), render: (row) => row.reportId || "—" },
              {
                key: "errorMessage",
                header: t("error"),
                render: (row) => (
                  <Text size="1" color={row.errorMessage ? "red" : "gray"}>
                    {row.errorMessage || "—"}
                  </Text>
                ),
              },
              { key: "createdAt", header: t("createdAt"), render: (row) => formatDateTime(row.createdAt) },
              { key: "sentAt", header: t("sentAt"), render: (row) => (row.sentAt ? formatDateTime(row.sentAt) : "—") },
            ]}
          />
        </SectionBlock>
      </Flex>
    </PageFrame>
  );
});

NotifyPage.displayName = "NotifyPage";
export default NotifyPage;
