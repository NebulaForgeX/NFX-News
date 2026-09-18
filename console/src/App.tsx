import { GuestRoute, ProtectedRoute } from "nfx-ui/navigations";
import { Navigate, Route, Routes } from "react-router";

import { Main, Sidebar } from "@/layouts";
import { ROUTES } from "@/navigations";
import {
  CrawlPage,
  GitHubCallbackPage,
  LoginPage,
  MCPPage,
  NotFoundPage,
  NotifyPage,
  ProfileEditPage,
  ProfileIdentitiesPage,
  ProfileOverviewPage,
  ReaderPage,
  ReportsPage,
  SettingsPage,
  SignupPage,
} from "@/pages";

export default function App() {
  return (
    <Routes>
      <Route element={<Main />}>
        <Route path={ROUTES.HOME} element={<Navigate to={ROUTES.LOGIN} replace />} />
      </Route>

      <Route element={<ProtectedRoute redirectTo={ROUTES.LOGIN} />}>
        <Route element={<Sidebar />}>
          <Route path={ROUTES.USER} element={<Navigate to={ROUTES.READER} replace />} />
          <Route path={ROUTES.USER_OVERVIEW} element={<Navigate to={ROUTES.READER} replace />} />
          <Route path={ROUTES.READER} element={<ReaderPage />} />
          <Route path={ROUTES.REPORTS} element={<ReportsPage />} />
          <Route path={ROUTES.CRAWL} element={<CrawlPage />} />
          <Route path={ROUTES.MCP} element={<MCPPage />} />
          <Route path={ROUTES.NOTIFY} element={<NotifyPage />} />
          <Route path={ROUTES.PROFILE} element={<Navigate to={ROUTES.USER_PROFILE_OVERVIEW} replace />} />
          <Route path={ROUTES.USER_PROFILE_OVERVIEW} element={<ProfileOverviewPage />} />
          <Route path={ROUTES.USER_PROFILE_EDIT} element={<ProfileEditPage />} />
          <Route path={ROUTES.USER_PROFILE_IDENTITIES} element={<ProfileIdentitiesPage />} />
          <Route path={ROUTES.USER_SETTINGS} element={<SettingsPage />} />
        </Route>
      </Route>

      <Route element={<GuestRoute redirectTo={ROUTES.READER} />}>
        <Route path={ROUTES.LOGIN} element={<LoginPage />} />
        <Route path={ROUTES.SIGNUP} element={<SignupPage />} />
      </Route>

      <Route path={ROUTES.LOGIN_GITHUB_CALLBACK} element={<GitHubCallbackPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
