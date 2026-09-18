import { memo, useState } from "react";
import { Button, Card, Flex, Select, Switch, Text, TextField } from "@radix-ui/themes";
import { Bell, Inbox, Server } from "lucide-react";
import { useTranslation } from "react-i18next";
import { PageFrame } from "@/layouts";
import { CardHeader, EmptyState, PageHeader } from "@/components";

import { NOTIFY_KINDS, notifyChannelConfig, notifyConfigPlaceholder, type NotifyKind } from "@/enums/newsEnum";
import { useChannels, useDeliveries, useInitializeSystem, useNotifyKinds, useSystemState, useUpsertChannel } from "@/hooks/news";

const SettingsPage = memo(() => {
  const { t } = useTranslation("SettingsPage");
  const { data: channels } = useChannels();
  const { data: deliveries } = useDeliveries();
  const { data: kinds } = useNotifyKinds();
  const system = useSystemState();
  const init = useInitializeSystem();
  const [kind, setKind] = useState<NotifyKind>("feishu");
  const [name, setName] = useState("");
  const [webhook, setWebhook] = useState("");
  const [enabled, setEnabled] = useState(true);
  const save = useUpsertChannel();
  const kindOptions = (kinds && kinds.length > 0 ? kinds : [...NOTIFY_KINDS]) as NotifyKind[];

  return (
    <PageFrame>
      <PageHeader icon={Bell} title={t("title")} description={t("webhookHint")} />
      <Flex direction="column" gap="4">
        <Card>
          <CardHeader icon={<Server size={18} />} title={t("system")} />
          <Flex gap="2" align="center" wrap="wrap">
            <Text size="2">{system.data?.initialized ? t("initialized") : t("notInitialized")}</Text>
            <Button variant="soft" onClick={() => init.mutate()} disabled={init.isPending}>
              {t("initialize")}
            </Button>
          </Flex>
        </Card>
        <Card>
          <CardHeader icon={<Bell size={18} />} title={t("channels")} />
          <Flex gap="2" wrap="wrap">
            <Select.Root value={kind} onValueChange={(v) => setKind(v as NotifyKind)}>
              <Select.Trigger />
              <Select.Content>
                {kindOptions.map((k) => (
                  <Select.Item key={k} value={k}>
                    {k}
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Root>
            <TextField.Root value={name} onChange={(e) => setName(e.target.value)} placeholder={t("name")} />
            <TextField.Root value={webhook} onChange={(e) => setWebhook(e.target.value)} placeholder={notifyConfigPlaceholder(kind)} />
            <Flex align="center" gap="2">
              <Switch checked={enabled} onCheckedChange={setEnabled} />
              <Text size="2">{t("enabled")}</Text>
            </Flex>
            <Button
              onClick={() =>
                save.mutate({
                  kind,
                  name,
                  enabled,
                  config: notifyChannelConfig(kind, webhook),
                })
              }
              disabled={!name}
            >
              {t("save")}
            </Button>
          </Flex>
          {(channels ?? []).length === 0 ? (
            <EmptyState icon={Bell} title={t("emptyChannels")} />
          ) : (
            <Flex direction="column" gap="2" mt="3">
              {(channels ?? []).map((ch) => (
                <Text key={ch.id} size="2">
                  {ch.kind} · {ch.name} · {ch.enabled ? "on" : "off"}
                </Text>
              ))}
            </Flex>
          )}
        </Card>
        <Card>
          <CardHeader icon={<Inbox size={18} />} title={t("deliveries")} />
          {(deliveries ?? []).length === 0 ? (
            <EmptyState icon={Inbox} title={t("emptyDeliveries")} />
          ) : (
            <Flex direction="column" gap="2">
              {(deliveries ?? []).map((row) => (
                <Text key={row.id} size="2">
                  {row.status} · {row.channelId} · {row.createdAt}
                </Text>
              ))}
            </Flex>
          )}
        </Card>
      </Flex>
    </PageFrame>
  );
});

SettingsPage.displayName = "SettingsPage";
export default SettingsPage;
