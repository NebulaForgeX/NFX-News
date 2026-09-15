import { memo, useState } from "react";
import { Button, Card, Flex, Switch, Text, TextField } from "@radix-ui/themes";
import { Bell, Inbox } from "lucide-react";
import { useTranslation } from "react-i18next";
import { PageFrame } from "nfx-ui/layouts";
import { CardHeader, EmptyState, PageHeader, ThemeSettings } from "nfx-ui/components";

import { useChannels, useDeliveries, useUpsertChannel } from "@/hooks/news";

const SettingsPage = memo(() => {
  const { t } = useTranslation("SettingsPage");
  const { data: channels } = useChannels();
  const { data: deliveries } = useDeliveries();
  const [kind, setKind] = useState("feishu");
  const [name, setName] = useState("");
  const [enabled, setEnabled] = useState(true);
  const save = useUpsertChannel();

  return (
    <PageFrame>
      <PageHeader icon={Bell} title={t("title")} description={t("webhookHint")} />
      <Flex direction="column" gap="4">
        <Card>
          <CardHeader icon={<Bell size={18} />} title={t("channels")} />
          <Flex gap="2" wrap="wrap">
            <TextField.Root value={kind} onChange={(e) => setKind(e.target.value)} placeholder="feishu / dingtalk / telegram" />
            <TextField.Root value={name} onChange={(e) => setName(e.target.value)} placeholder={t("name")} />
            <Flex align="center" gap="2">
              <Switch checked={enabled} onCheckedChange={setEnabled} />
              <Text size="2">{t("enabled")}</Text>
            </Flex>
            <Button onClick={() => save.mutate({ kind, name, enabled, config: {} })} disabled={!name}>
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
                  {row.status} · {row.channel_id} · {row.created_at}
                </Text>
              ))}
            </Flex>
          )}
        </Card>
        <ThemeSettings />
      </Flex>
    </PageFrame>
  );
});

SettingsPage.displayName = "SettingsPage";
export default SettingsPage;
