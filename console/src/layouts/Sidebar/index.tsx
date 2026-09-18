import type { ReactNode } from "react";

import { useEffect, useRef, useState } from "react";
import { Avatar, Box, Button, Flex, IconButton, Text } from "@radix-ui/themes";
import {
  Bell,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  Contact,
  FileText,
  IdCard,
  LayoutGrid,
  LogOut,
  LucideIcon as LucideIconType,
  Menu as MenuIcon,
  Newspaper,
  Pencil,
  Settings as SettingsIcon,
  UserRound,
  Wand2,
} from "lucide-react";
import { authEventEmitter, authEvents } from "nfx-ui/events";
import { useCurrentProfile } from "nfx-ui/hooks";
import { AuthStore, clearAuth } from "nfx-ui/stores";
import { useTranslation } from "react-i18next";
import { Menu, Sidebar as ProSidebar } from "react-pro-sidebar";
import { Link, Outlet, useLocation } from "react-router";

import { LucideIcon } from "@/components";
import UserTopBar from "@/layouts/UserTopBar";
import { ROUTES } from "@/navigations";
import { buildImageUrl, resolveAccountDisplayName, safeNullable } from "@/utils";

import { MenuItem, SidebarMenuState, SubMenu } from "./menu";
import styles from "./s.module.css";

const SIDEBAR_WIDTH = "234px";
const SIDEBAR_COLLAPSED_WIDTH = "84px";

function MenuLabel({ children, active = false }: { children: ReactNode; active?: boolean }) {
  return (
    <Text as="span" size="2" weight={active ? "bold" : "medium"}>
      {children}
    </Text>
  );
}

function SectionTitle({ label, icon }: { label: string; icon: LucideIconType }) {
  return (
    <Flex align="center" justify="between" gap="2" className={styles.sectionTitle}>
      <Text as="span" size="2" weight="bold">
        {label}
      </Text>
      <LucideIcon icon={icon} size={16} className={styles.sectionTitleIcon} />
    </Flex>
  );
}

function createMenuItemStyles(collapsed: boolean) {
  return {
    button: ({ active, level = 0 }: { active: boolean; level?: number }) => ({
      height: level > 0 ? "34px" : "40px",
      margin: level > 0 ? (collapsed ? "3px 7px" : "3px 0 3px 32px") : "6px 0",
      borderRadius: "var(--radius-chip)",
      paddingLeft: level > 0 ? "11px" : "10px",
      paddingRight: "10px",
      fontSize: level > 0 ? "14px" : "15px",
      fontWeight: 400,
      color: active ? "var(--accent-11)" : "var(--gray-11)",
      backgroundColor: active ? "var(--accent-a3)" : "transparent",
      transition: "background-color 150ms ease, color 150ms ease",
      "&:hover": { backgroundColor: active ? "var(--accent-a3)" : "var(--gray-a3)", color: active ? "var(--accent-11)" : "var(--gray-12)" },
      "&:focus-visible": {
        outline: "2px solid var(--accent-8)",
        outlineOffset: "2px",
      },
    }),
    icon: ({ level = 0 }: { level?: number }) => ({
      width: "20px",
      minWidth: "20px",
      height: "20px",
      marginRight: "10px",
      color: "inherit",
      ...(level > 0 ? { display: "none" } : {}),
    }),
    subMenuContent: {
      backgroundColor: collapsed ? "var(--color-panel-solid)" : "transparent",
      padding: collapsed ? "5px 0" : "0",
      ...(collapsed
        ? {
            zIndex: 1000,
            minWidth: "156px",
            width: "max-content",
            maxWidth: "260px",
            maxHeight: "calc(100dvh - 32px)",
            overflowY: "auto" as const,
            borderRadius: "var(--radius-4)",
            border: "1px solid var(--gray-a5)",
            boxShadow: "var(--shadow-5)",
          }
        : {}),
    },
  };
}

interface SectionProps {
  collapsed: boolean;
  broken: boolean;
  onMobileClose: () => void;
}

function OverviewSection({ collapsed, broken, onMobileClose }: SectionProps) {
  const { t } = useTranslation("language");
  const location = useLocation();
  const active = location.pathname === ROUTES.READER || location.pathname.startsWith(`${ROUTES.READER}/`);

  return (
    <Menu renderExpandIcon={({ open }) => <LucideIcon icon={open ? ChevronUp : ChevronDown} size={14} />} menuItemStyles={createMenuItemStyles(collapsed)} closeOnClick>
      <MenuItem component={<Link to={ROUTES.READER} />} icon={<LucideIcon icon={Newspaper} size={18} />} active={active} onClick={() => broken && onMobileClose()}>
        <MenuLabel active={active}>{t("sidebar.reader")}</MenuLabel>
      </MenuItem>
    </Menu>
  );
}

function MainMenuSection({ collapsed, broken, onMobileClose }: SectionProps) {
  const { t } = useTranslation("language");
  const location = useLocation();
  const [profileOpen, setProfileOpen] = useState(() => location.pathname.startsWith(ROUTES.PROFILE));

  const isActive = (to: string) => location.pathname === to || location.pathname.startsWith(`${to}/`);

  const profileSubItems = [
    {
      key: "profileOverview",
      to: ROUTES.USER_PROFILE_OVERVIEW,
      icon: <LucideIcon icon={Contact} size={16} />,
      label: t("sidebar.profileOverview"),
    },
    {
      key: "profileEdit",
      to: ROUTES.USER_PROFILE_EDIT,
      icon: <LucideIcon icon={Pencil} size={16} />,
      label: t("sidebar.profileEdit"),
    },
    {
      key: "profileIdentities",
      to: ROUTES.USER_PROFILE_IDENTITIES,
      icon: <LucideIcon icon={IdCard} size={16} />,
      label: t("sidebar.profileIdentities"),
    },
  ];

  const isProfileChildActive = profileSubItems.some((item) => isActive(item.to));

  return (
    <Menu renderExpandIcon={({ open }) => <LucideIcon icon={open ? ChevronUp : ChevronDown} size={14} />} menuItemStyles={createMenuItemStyles(collapsed)} closeOnClick>
      <SectionTitle label={t("sidebar.mainMenu")} icon={LayoutGrid} />
      <SubMenu label={t("sidebar.profile")} icon={<LucideIcon icon={UserRound} size={18} />} open={profileOpen} onOpenChange={setProfileOpen} active={isProfileChildActive}>
        {profileSubItems.map((item) => (
          <MenuItem key={item.key} component={<Link to={item.to} />} icon={item.icon} active={isActive(item.to)} onClick={() => broken && onMobileClose()}>
            <MenuLabel active={isActive(item.to)}>{item.label}</MenuLabel>
          </MenuItem>
        ))}
      </SubMenu>
      <MenuItem component={<Link to={ROUTES.REPORTS} />} icon={<LucideIcon icon={FileText} size={18} />} active={isActive(ROUTES.REPORTS)} onClick={() => broken && onMobileClose()}>
        <MenuLabel active={isActive(ROUTES.REPORTS)}>{t("sidebar.reports")}</MenuLabel>
      </MenuItem>
      <MenuItem component={<Link to={ROUTES.CRAWL} />} icon={<LucideIcon icon={FileText} size={18} />} active={isActive(ROUTES.CRAWL)} onClick={() => broken && onMobileClose()}>
        <MenuLabel active={isActive(ROUTES.CRAWL)}>{t("sidebar.crawl")}</MenuLabel>
      </MenuItem>
      <MenuItem component={<Link to={ROUTES.MCP} />} icon={<LucideIcon icon={Wand2} size={18} />} active={isActive(ROUTES.MCP)} onClick={() => broken && onMobileClose()}>
        <MenuLabel active={isActive(ROUTES.MCP)}>{t("sidebar.mcp")}</MenuLabel>
      </MenuItem>
      <MenuItem component={<Link to={ROUTES.NOTIFY} />} icon={<LucideIcon icon={Bell} size={18} />} active={isActive(ROUTES.NOTIFY)} onClick={() => broken && onMobileClose()}>
        <MenuLabel active={isActive(ROUTES.NOTIFY)}>{t("sidebar.notify")}</MenuLabel>
      </MenuItem>
    </Menu>
  );
}

function SettingsSection({ collapsed, broken, onMobileClose }: SectionProps) {
  const { t } = useTranslation("language");
  const location = useLocation();
  const isActive = (to: string) => location.pathname === to || location.pathname.startsWith(`${to}/`);

  return (
    <Menu renderExpandIcon={({ open }) => <LucideIcon icon={open ? ChevronUp : ChevronDown} size={14} />} menuItemStyles={createMenuItemStyles(collapsed)} closeOnClick>
      <SectionTitle label={t("sidebar.settings")} icon={SettingsIcon} />
      <MenuItem
        component={<Link to={ROUTES.USER_SETTINGS} />}
        icon={<LucideIcon icon={SettingsIcon} size={18} />}
        active={isActive(ROUTES.USER_SETTINGS)}
        onClick={() => broken && onMobileClose()}
      >
        <MenuLabel active={isActive(ROUTES.USER_SETTINGS)}>{t("sidebar.settingsItem")}</MenuLabel>
      </MenuItem>
    </Menu>
  );
}

function Sidebar() {
  const { t } = useTranslation("language");
  const [desktopCollapsed, setCollapsed] = useState(false);
  const [toggled, setToggled] = useState(false);
  const [broken, setBroken] = useState(false);
  const collapsed = !broken && desktopCollapsed;
  const drawerRef = useRef<HTMLHtmlElement>(null);
  useEffect(() => {
    if (!broken || !toggled) return;
    const previousFocus = document.activeElement as HTMLElement | null;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    const focusable = () =>
      Array.from(drawerRef.current?.querySelectorAll<HTMLElement>("button:not([disabled]), a[href], [tabindex='0']") ?? []).filter(
        (node) => node.getClientRects().length && getComputedStyle(node).visibility !== "hidden",
      );
    focusable()[0]?.focus();
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.defaultPrevented) return;
      if (event.key === "Escape") {
        event.preventDefault();
        setToggled(false);
      }
      if (event.key !== "Tab") return;
      const nodes = focusable();
      const first = nodes[0];
      const last = nodes[nodes.length - 1];
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last?.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first?.focus();
      }
    };
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.body.style.overflow = previousOverflow;
      document.removeEventListener("keydown", onKeyDown);
      previousFocus?.focus();
    };
  }, [broken, toggled]);

  const { kind, data, profile } = useCurrentProfile();

  const accountId = safeNullable(data?.account.id);
  const displayName = resolveAccountDisplayName(profile?.displayName, accountId);
  const avatarImageId = safeNullable(profile?.avatars?.[0]?.imageId);

  const closeMobile = () => broken && setToggled(false);

  const handleLogout = () => {
    const aID = AuthStore.getState().currentAccountId;
    if (aID) authEventEmitter.emit(authEvents.LOGOUT, aID);
    clearAuth();
  };

  return (
    <Flex minHeight="100dvh" width="100%" className={styles.shell}>
      <SidebarMenuState collapsed={collapsed}>
        <ProSidebar
          ref={drawerRef}
          inert={broken && !toggled ? true : undefined}
          aria-hidden={broken && !toggled ? true : undefined}
          collapsed={collapsed}
          toggled={toggled}
          onBackdropClick={() => setToggled(false)}
          onBreakPoint={setBroken}
          breakPoint="md"
          width={SIDEBAR_WIDTH}
          collapsedWidth={SIDEBAR_COLLAPSED_WIDTH}
          rootStyles={{
            border: 0,
            height: "100dvh",
            flexShrink: 0,
            position: "sticky",
            top: 0,
            margin: 0,
            zIndex: 200,
            overflow: "visible",
            "&.ps-broken.ps-toggled": { left: 0 },
            "&.ps-broken": { height: "100dvh", position: "fixed", margin: 0, top: 0, bottom: 0 },
            "& .ps-sidebar-container": {
              background: "transparent",
              height: "100%",
              overflow: "visible",
            },
          }}
        >
          <Flex direction="column" height="100%" minHeight="0" className={styles.sidebar}>
            <div className={styles.header}>
              <button type="button" className={styles.accountCard} aria-label={displayName}>
                <Avatar
                  size="3"
                  className={styles.avatar}
                  src={avatarImageId ? buildImageUrl(avatarImageId) : undefined}
                  fallback={<UserRound size={20} />}
                  alt=""
                  aria-hidden="true"
                />
                {!collapsed && (
                  <span className={styles.accountInfo}>
                    <span className={styles.accountRole}>{t(kind === "authority" ? "sidebar.profileAuthority" : "sidebar.profileForger")}</span>
                    <span className={styles.accountName}>{displayName}</span>
                  </span>
                )}
              </button>
              <IconButton
                variant="soft"
                size="1"
                className={styles.toggle}
                aria-label={collapsed ? t("sidebar.expand") : t("sidebar.collapse")}
                aria-expanded={!collapsed}
                onClick={() => (broken ? setToggled(false) : setCollapsed((value) => !value))}
              >
                <LucideIcon icon={collapsed ? ChevronRight : ChevronLeft} size={14} />
              </IconButton>
            </div>

            <Box flexGrow="1" minHeight="0" py="2" className={`${styles.menuArea} ${collapsed ? styles.menuAreaCollapsed : ""}`}>
              <OverviewSection collapsed={collapsed} broken={broken} onMobileClose={closeMobile} />
              <MainMenuSection collapsed={collapsed} broken={broken} onMobileClose={closeMobile} />
              <SettingsSection collapsed={collapsed} broken={broken} onMobileClose={closeMobile} />
            </Box>

            <div className={`${styles.footer} ${collapsed ? styles.footerCollapsed : ""}`}>
              <Button
                variant="ghost"
                className={`${styles.logout} ${collapsed ? styles.logoutCollapsed : ""}`}
                onClick={handleLogout}
                aria-label={t("sidebar.logout")}
                title={collapsed ? t("sidebar.logout") : undefined}
              >
                <LucideIcon icon={LogOut} size={18} />
                {!collapsed && t("sidebar.logout")}
              </Button>
            </div>
          </Flex>
        </ProSidebar>
      </SidebarMenuState>

      <Flex direction="column" flexGrow="1" minWidth="0" inert={broken && toggled ? true : undefined} className={styles.content}>
        {broken ? (
          <IconButton variant="surface" color="gray" size="3" className={styles.mobileToggle} onClick={() => setToggled(true)} aria-label={t("sidebar.openMenu")}>
            <LucideIcon icon={MenuIcon} size={18} />
          </IconButton>
        ) : null}
        <UserTopBar />
        <Box flexGrow="1" minWidth="0" minHeight="0" className={styles.contentInner}>
          <Outlet />
        </Box>
      </Flex>
    </Flex>
  );
}

export default Sidebar;
