import type { CreateI18nResourcesResult, NameSpacesMap, Resources } from "nfx-ui/languages";

import enHooks from "./en/hooks.json";
import enLanguage from "./en/language.json";
import enAuthShell from "./en/pages/Account/AuthShell.json";
import enLogin from "./en/pages/Account/Login.json";
import enSignup from "./en/pages/Account/Signup.json";
import enCrawl from "./en/pages/Crawl.json";
import enMCP from "./en/pages/MCP.json";
import enNotFound from "./en/pages/NotFound.json";
import enReader from "./en/pages/Reader.json";
import enReports from "./en/pages/Reports.json";
import enSources from "./en/pages/Sources.json";
import enSystem from "./en/pages/System.json";
import enSettings from "./en/pages/Notify.json";
import enUserProfileEdit from "./en/pages/User/Profile/Edit.json";
import enUserProfileIdentities from "./en/pages/User/Profile/Identities.json";
import enUserProfileOverview from "./en/pages/User/Profile/Overview.json";
import enUserSetting from "./en/pages/User/Setting.json";
import frHooks from "./fr/hooks.json";
import frLanguage from "./fr/language.json";
import frAuthShell from "./fr/pages/Account/AuthShell.json";
import frLogin from "./fr/pages/Account/Login.json";
import frSignup from "./fr/pages/Account/Signup.json";
import frCrawl from "./fr/pages/Crawl.json";
import frMCP from "./fr/pages/MCP.json";
import frNotFound from "./fr/pages/NotFound.json";
import frReader from "./fr/pages/Reader.json";
import frReports from "./fr/pages/Reports.json";
import frSources from "./fr/pages/Sources.json";
import frSystem from "./fr/pages/System.json";
import frSettings from "./fr/pages/Notify.json";
import frUserProfileEdit from "./fr/pages/User/Profile/Edit.json";
import frUserProfileIdentities from "./fr/pages/User/Profile/Identities.json";
import frUserProfileOverview from "./fr/pages/User/Profile/Overview.json";
import frUserSetting from "./fr/pages/User/Setting.json";
import zhHooks from "./zh/hooks.json";
import zhLanguage from "./zh/language.json";
import zhAuthShell from "./zh/pages/Account/AuthShell.json";
import zhLogin from "./zh/pages/Account/Login.json";
import zhSignup from "./zh/pages/Account/Signup.json";
import zhCrawl from "./zh/pages/Crawl.json";
import zhMCP from "./zh/pages/MCP.json";
import zhNotFound from "./zh/pages/NotFound.json";
import zhReader from "./zh/pages/Reader.json";
import zhReports from "./zh/pages/Reports.json";
import zhSources from "./zh/pages/Sources.json";
import zhSystem from "./zh/pages/System.json";
import zhSettings from "./zh/pages/Notify.json";
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
  Reader: "pages.Reader",
  Reports: "pages.Reports",
  Crawl: "pages.Crawl",
  MCP: "pages.MCP",
  Notify: "pages.Notify",
  Sources: "pages.Sources",
  System: "pages.System",
  NotFound: "pages.NotFound",
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
      [PAGE.Reader]: enReader,
      [PAGE.Reports]: enReports,
      [PAGE.Crawl]: enCrawl,
      [PAGE.MCP]: enMCP,
      [PAGE.Notify]: enSettings,
      [PAGE.Sources]: enSources,
      [PAGE.System]: enSystem,
      [PAGE.NotFound]: enNotFound,
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
      [PAGE.Reader]: zhReader,
      [PAGE.Reports]: zhReports,
      [PAGE.Crawl]: zhCrawl,
      [PAGE.MCP]: zhMCP,
      [PAGE.Notify]: zhSettings,
      [PAGE.Sources]: zhSources,
      [PAGE.System]: zhSystem,
      [PAGE.NotFound]: zhNotFound,
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
      [PAGE.Reader]: frReader,
      [PAGE.Reports]: frReports,
      [PAGE.Crawl]: frCrawl,
      [PAGE.MCP]: frMCP,
      [PAGE.Notify]: frSettings,
      [PAGE.Sources]: frSources,
      [PAGE.System]: frSystem,
      [PAGE.NotFound]: frNotFound,
    },
  };
  return {
    RESOURCES,
    NAME_SPACES_MAP: BUILTIN_I18N_NAMESPACES_MAP,
    NAME_SPACES: Object.values(BUILTIN_I18N_NAMESPACES_MAP),
  };
}
