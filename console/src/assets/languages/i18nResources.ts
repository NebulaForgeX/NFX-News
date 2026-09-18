import type { CreateI18nResourcesResult, NameSpacesMap, Resources } from "nfx-ui/languages";

import enHooks from "./en/hooks.json";
import enLanguage from "./en/language.json";
import enAuthShell from "./en/pages/Account/AuthShell.json";
import enLogin from "./en/pages/Account/Login.json";
import enSignup from "./en/pages/Account/Signup.json";
import enCrawl from "./en/CrawlPage.json";
import enMCP from "./en/MCPPage.json";
import enNotFound from "./en/NotFoundPage.json";
import enReader from "./en/ReaderPage.json";
import enReports from "./en/ReportsPage.json";
import enSettings from "./en/SettingsPage.json";
import enUserProfileEdit from "./en/pages/User/Profile/Edit.json";
import enUserProfileIdentities from "./en/pages/User/Profile/Identities.json";
import enUserProfileOverview from "./en/pages/User/Profile/Overview.json";
import enUserSetting from "./en/pages/User/Setting.json";
import frHooks from "./fr/hooks.json";
import frLanguage from "./fr/language.json";
import frAuthShell from "./fr/pages/Account/AuthShell.json";
import frLogin from "./fr/pages/Account/Login.json";
import frSignup from "./fr/pages/Account/Signup.json";
import frCrawl from "./fr/CrawlPage.json";
import frMCP from "./fr/MCPPage.json";
import frNotFound from "./fr/NotFoundPage.json";
import frReader from "./fr/ReaderPage.json";
import frReports from "./fr/ReportsPage.json";
import frSettings from "./fr/SettingsPage.json";
import frUserProfileEdit from "./fr/pages/User/Profile/Edit.json";
import frUserProfileIdentities from "./fr/pages/User/Profile/Identities.json";
import frUserProfileOverview from "./fr/pages/User/Profile/Overview.json";
import frUserSetting from "./fr/pages/User/Setting.json";
import zhHooks from "./zh/hooks.json";
import zhLanguage from "./zh/language.json";
import zhAuthShell from "./zh/pages/Account/AuthShell.json";
import zhLogin from "./zh/pages/Account/Login.json";
import zhSignup from "./zh/pages/Account/Signup.json";
import zhCrawl from "./zh/CrawlPage.json";
import zhMCP from "./zh/MCPPage.json";
import zhNotFound from "./zh/NotFoundPage.json";
import zhReader from "./zh/ReaderPage.json";
import zhReports from "./zh/ReportsPage.json";
import zhSettings from "./zh/SettingsPage.json";
import zhUserProfileEdit from "./zh/pages/User/Profile/Edit.json";
import zhUserProfileIdentities from "./zh/pages/User/Profile/Identities.json";
import zhUserProfileOverview from "./zh/pages/User/Profile/Overview.json";
import zhUserSetting from "./zh/pages/User/Setting.json";

const PAGE = {
  AuthShell: "pages.Account.AuthShell",
  Login: "pages.Account.Login",
  Signup: "pages.Account.Signup",
  UserSetting: "pages.User.Setting",
  UserProfileOverview: "pages.User.Profile.Overview",
  UserProfileEdit: "pages.User.Profile.Edit",
  UserProfileIdentities: "pages.User.Profile.Identities",
  ReaderPage: "ReaderPage",
  ReportsPage: "ReportsPage",
  CrawlPage: "CrawlPage",
  MCPPage: "MCPPage",
  SettingsPage: "SettingsPage",
  NotFoundPage: "NotFoundPage",
} as const;

const BUILTIN_I18N_NAMESPACES_MAP: NameSpacesMap = {
  language: "language",
  hooks: "hooks",
  ...PAGE,
};

export function getBuiltinI18nBundles(): CreateI18nResourcesResult {
  const RESOURCES: Resources = {
    en: {
      language: enLanguage,
      hooks: enHooks,
      [PAGE.AuthShell]: enAuthShell,
      [PAGE.Login]: enLogin,
      [PAGE.Signup]: enSignup,
      [PAGE.UserSetting]: enUserSetting,
      [PAGE.UserProfileOverview]: enUserProfileOverview,
      [PAGE.UserProfileEdit]: enUserProfileEdit,
      [PAGE.UserProfileIdentities]: enUserProfileIdentities,
      [PAGE.ReaderPage]: enReader,
      [PAGE.ReportsPage]: enReports,
      [PAGE.CrawlPage]: enCrawl,
      [PAGE.MCPPage]: enMCP,
      [PAGE.SettingsPage]: enSettings,
      [PAGE.NotFoundPage]: enNotFound,
    },
    zh: {
      language: zhLanguage,
      hooks: zhHooks,
      [PAGE.AuthShell]: zhAuthShell,
      [PAGE.Login]: zhLogin,
      [PAGE.Signup]: zhSignup,
      [PAGE.UserSetting]: zhUserSetting,
      [PAGE.UserProfileOverview]: zhUserProfileOverview,
      [PAGE.UserProfileEdit]: zhUserProfileEdit,
      [PAGE.UserProfileIdentities]: zhUserProfileIdentities,
      [PAGE.ReaderPage]: zhReader,
      [PAGE.ReportsPage]: zhReports,
      [PAGE.CrawlPage]: zhCrawl,
      [PAGE.MCPPage]: zhMCP,
      [PAGE.SettingsPage]: zhSettings,
      [PAGE.NotFoundPage]: zhNotFound,
    },
    fr: {
      language: frLanguage,
      hooks: frHooks,
      [PAGE.AuthShell]: frAuthShell,
      [PAGE.Login]: frLogin,
      [PAGE.Signup]: frSignup,
      [PAGE.UserSetting]: frUserSetting,
      [PAGE.UserProfileOverview]: frUserProfileOverview,
      [PAGE.UserProfileEdit]: frUserProfileEdit,
      [PAGE.UserProfileIdentities]: frUserProfileIdentities,
      [PAGE.ReaderPage]: frReader,
      [PAGE.ReportsPage]: frReports,
      [PAGE.CrawlPage]: frCrawl,
      [PAGE.MCPPage]: frMCP,
      [PAGE.SettingsPage]: frSettings,
      [PAGE.NotFoundPage]: frNotFound,
    },
  };
  return {
    RESOURCES,
    NAME_SPACES_MAP: BUILTIN_I18N_NAMESPACES_MAP,
    NAME_SPACES: Object.values(BUILTIN_I18N_NAMESPACES_MAP),
  };
}
