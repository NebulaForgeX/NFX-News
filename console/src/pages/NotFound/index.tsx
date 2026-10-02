import { InfoCircleIcon } from "nfx-ui/icons";
import { Button, Flex } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { PageFrame } from "@/layouts";
import { EmptyState } from "@/components";

import { routerEventEmitter } from "@/events/router";

export default function NotFoundPage() {
  const { t } = useTranslation("pages.NotFound");
  return (
    <PageFrame>
      <EmptyState
        icon={InfoCircleIcon}
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
