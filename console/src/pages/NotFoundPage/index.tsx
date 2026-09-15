import { Button, Flex } from "@radix-ui/themes";
import { Info } from "@/assets/icons/lucide";
import { PageFrame } from "nfx-ui/layouts";
import { EmptyState } from "nfx-ui/components";

import { routerEventEmitter } from "@/events/router";

export default function NotFoundPage() {
  return (
    <PageFrame>
      <EmptyState
        icon={Info}
        title="Page Not Found"
        description="The page might have been moved or deleted."
        action={
          <Flex gap="3">
            <Button onClick={() => routerEventEmitter.navigateToDashboard()}>Go to Reader</Button>
          </Flex>
        }
      />
    </PageFrame>
  );
}
