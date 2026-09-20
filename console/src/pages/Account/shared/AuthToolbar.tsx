import { Flex } from "@radix-ui/themes";
import { APP_NAME } from "nfx-ui/config";

import { Logo, PreferencesPopover } from "@/components";
import { ROUTES } from "@/navigations";

import styles from "./s.module.css";

export default function AuthToolbar() {
  return (
    <Flex align="center" gap="2" className={styles.wrap}>
      <Logo variant="plain" size="small" to={ROUTES.HOME} alt={`${APP_NAME} logo`} />
      <PreferencesPopover />
    </Flex>
  );
}
