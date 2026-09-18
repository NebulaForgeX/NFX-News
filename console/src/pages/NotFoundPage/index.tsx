import { Button, Flex } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { Info } from "@/assets/icons/lucide";
import { PageFrame } from "nfx-ui/layouts";
import { EmptyState } from "nfx-ui/components";

import { routerEventEmitter } from "@/events/router";

export default function NotFoundPage() {
  const { t } = useTranslation("NotFoundPage");
  return (
    <PageFrame>
      <EmptyState
        icon={Info}
        title={t("title")}
        description={t("description")}
        action={
          <Flex gap="3">
            <Button onClick={() => routerEventEmitter.navigateToDashboard()}>{t("goReader")}</Button>
          </Flex>
        }
      />
    </PageFrame>
  );
}
