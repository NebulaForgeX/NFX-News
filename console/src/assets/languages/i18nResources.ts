import en_components from "./en/components.json";
import en_LoginPage from "./en/LoginPage.json";
import en_ReaderPage from "./en/ReaderPage.json";
import en_ReportsPage from "./en/ReportsPage.json";
import en_CrawlPage from "./en/CrawlPage.json";
import en_MCPPage from "./en/MCPPage.json";
import en_SettingsPage from "./en/SettingsPage.json";
import zh_components from "./zh/components.json";
import zh_LoginPage from "./zh/LoginPage.json";
import zh_ReaderPage from "./zh/ReaderPage.json";
import zh_ReportsPage from "./zh/ReportsPage.json";
import zh_CrawlPage from "./zh/CrawlPage.json";
import zh_MCPPage from "./zh/MCPPage.json";
import zh_SettingsPage from "./zh/SettingsPage.json";
import fr_components from "./fr/components.json";
import fr_LoginPage from "./fr/LoginPage.json";
import fr_ReaderPage from "./fr/ReaderPage.json";
import fr_ReportsPage from "./fr/ReportsPage.json";
import fr_CrawlPage from "./fr/CrawlPage.json";
import fr_MCPPage from "./fr/MCPPage.json";
import fr_SettingsPage from "./fr/SettingsPage.json";
import en_NotFoundPage from "./en/NotFoundPage.json";
import en_SelectProfilePage from "./en/SelectProfilePage.json";
import zh_NotFoundPage from "./zh/NotFoundPage.json";
import zh_SelectProfilePage from "./zh/SelectProfilePage.json";
import fr_NotFoundPage from "./fr/NotFoundPage.json";
import fr_SelectProfilePage from "./fr/SelectProfilePage.json";

export const RESOURCES = {
  en: {
    components: en_components,
    LoginPage: en_LoginPage,
    ReaderPage: en_ReaderPage,
    ReportsPage: en_ReportsPage,
    CrawlPage: en_CrawlPage,
    MCPPage: en_MCPPage,
    SettingsPage: en_SettingsPage,
    NotFoundPage: en_NotFoundPage,
    SelectProfilePage: en_SelectProfilePage,
  },
  zh: {
    components: zh_components,
    LoginPage: zh_LoginPage,
    ReaderPage: zh_ReaderPage,
    ReportsPage: zh_ReportsPage,
    CrawlPage: zh_CrawlPage,
    MCPPage: zh_MCPPage,
    SettingsPage: zh_SettingsPage,
    NotFoundPage: zh_NotFoundPage,
    SelectProfilePage: zh_SelectProfilePage,
  },
  fr: {
    components: fr_components,
    LoginPage: fr_LoginPage,
    ReaderPage: fr_ReaderPage,
    ReportsPage: fr_ReportsPage,
    CrawlPage: fr_CrawlPage,
    MCPPage: fr_MCPPage,
    SettingsPage: fr_SettingsPage,
    NotFoundPage: fr_NotFoundPage,
    SelectProfilePage: fr_SelectProfilePage,
  },
};

export const NAME_SPACES_MAP = {
  components: "components",
  LoginPage: "LoginPage",
  ReaderPage: "ReaderPage",
  ReportsPage: "ReportsPage",
  CrawlPage: "CrawlPage",
  MCPPage: "MCPPage",
  SettingsPage: "SettingsPage",
  NotFoundPage: "NotFoundPage",
  SelectProfilePage: "SelectProfilePage",
};

export const NAME_SPACES = Object.values(NAME_SPACES_MAP);

export function getBuiltinI18nBundles() {
  return { RESOURCES, NAME_SPACES_MAP, NAME_SPACES };
}
