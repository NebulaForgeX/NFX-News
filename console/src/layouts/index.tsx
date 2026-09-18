import type { SidebarMenuItem } from "nfx-ui/layouts";
import type { ReactNode } from "react";

import { memo, useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useLocation } from "react-router";
import { FileText, Home, Settings, Wand2 } from "@/assets/icons/lucide";

import { LayoutFrame } from "nfx-ui/layouts";
import { Logo, PreferencesPopover } from "nfx-ui/components";

import { authEventEmitter, authEvents } from "@/events/auth";
import { routerEventEmitter } from "@/events/router";
import { ROUTES } from "@/navigations";

export const ConsoleLayout = memo(({ children }: { children: ReactNode }) => {
  const { t } = useTranslation("components");
  const location = useLocation();
  const sidebarItems: SidebarMenuItem[] = useMemo(
    () => [
      { label: t("sidebar.reader"), path: ROUTES.READER, icon: <Home size={20} /> },
      { label: t("sidebar.reports"), path: ROUTES.REPORTS, icon: <FileText size={20} /> },
      { label: t("sidebar.crawl"), path: ROUTES.CRAWL, icon: <FileText size={20} /> },
      { label: t("sidebar.mcp"), path: ROUTES.MCP, icon: <Wand2 size={20} /> },
      { label: t("sidebar.settings"), path: ROUTES.SETTINGS, icon: <Settings size={20} /> },
    ],
    [t],
  );

  const onSidebarNavigate = useCallback((path: string) => {
    routerEventEmitter.navigate({ to: path });
  }, []);

  return (
    <LayoutFrame
      headerLeft={<Logo title="NFX" subtitle="News" alt="NFX" onClick={() => routerEventEmitter.navigateToDashboard()} />}
      headerRight={<PreferencesPopover />}
      sidebarItems={sidebarItems}
      sidebarCurrentPathname={location.pathname}
      onSidebarNavigate={onSidebarNavigate}
      sidebarLogoutLabel={t("header.logout")}
      onSidebarLogout={() => authEventEmitter.emit(authEvents.LOGOUT)}
    >
      {children}
    </LayoutFrame>
  );
});

ConsoleLayout.displayName = "ConsoleLayout";
export default ConsoleLayout;
