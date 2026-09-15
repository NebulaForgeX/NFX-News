import { useState } from "react";
import { Button, Card, Flex, Heading, Switch, Text, TextField } from "@radix-ui/themes";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { useChannels, useDeliveries } from "@/hooks/news";
import { useNewsRepositories } from "@/apis/repositories";

export default function SettingsPage() {
  const { t } = useTranslation("SettingsPage");
  const { data: channels } = useChannels();
  const { data: deliveries } = useDeliveries();
  const repos = useNewsRepositories();
  const qc = useQueryClient();
  const [kind, setKind] = useState("feishu");
  const [name, setName] = useState("");
  const [enabled, setEnabled] = useState(true);
  const save = useMutation({
    mutationFn: () => repos.notify.UpsertChannel({ kind, name, enabled, config: {} }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notify-channels"] }),
  });

  return (
    <Flex direction="column" gap="4" p="4">
      <Heading size="6">{t("title")}</Heading>
      <Text color="gray">{t("webhookHint")}</Text>
      <Card>
        <Heading size="4">{t("channels")}</Heading>
        <Flex gap="2" mt="3" wrap="wrap">
          <TextField.Root value={kind} onChange={(e) => setKind(e.target.value)} placeholder="feishu / dingtalk / telegram" />
          <TextField.Root value={name} onChange={(e) => setName(e.target.value)} placeholder={t("name")} />
          <Flex align="center" gap="2">
            <Switch checked={enabled} onCheckedChange={setEnabled} />
            <Text size="2">{t("enabled")}</Text>
          </Flex>
          <Button onClick={() => save.mutate()} disabled={!name}>
            {t("save")}
          </Button>
        </Flex>
        <Flex direction="column" gap="2" mt="3">
          {(channels ?? []).map((ch) => (
            <Text key={ch.id} size="2">
              {ch.kind} · {ch.name} · {ch.enabled ? "on" : "off"}
            </Text>
          ))}
        </Flex>
      </Card>
      <Card>
        <Heading size="4">{t("deliveries", { defaultValue: "Deliveries" })}</Heading>
        <Flex direction="column" gap="2" mt="3">
          {(deliveries ?? []).length === 0 ? (
            <Text size="2" color="gray">
              —
            </Text>
          ) : (
            (deliveries ?? []).map((row) => (
              <Text key={row.id} size="2">
                {row.status} · {row.channel_id} · {row.created_at}
              </Text>
            ))
          )}
        </Flex>
      </Card>
    </Flex>
  );
}
