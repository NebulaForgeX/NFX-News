import { Navigate, Route, Routes } from "react-router";
import { useAuthStore, hasSelectedProfile } from "nfx-ui/stores";

import { ConsoleLayout } from "@/layouts";
import { CrawlPage, LoginPage, NotFoundPage, ReaderPage, ReportsPage, SelectProfilePage, SettingsPage } from "@/pages";
import { ROUTES } from "@/navigations";

import "./App.module.css";

function App() {
  const accessToken = useAuthStore((state) => state.accessToken);
  const profileId = useAuthStore((state) => state.currentProfileId);
  const isAuthValid = useAuthStore((state) => state.isAuthValid);

  if (!accessToken) {
    return (
      <Routes>
        <Route path={ROUTES.LOGIN} element={<LoginPage />} />
        <Route path="*" element={<Navigate to={ROUTES.LOGIN} replace />} />
      </Routes>
    );
  }

  if (!isAuthValid || !hasSelectedProfile(profileId)) {
    return (
      <Routes>
        <Route path={ROUTES.SELECT_PROFILE} element={<SelectProfilePage />} />
        <Route path="*" element={<Navigate to={ROUTES.SELECT_PROFILE} replace />} />
      </Routes>
    );
  }

  return (
    <ConsoleLayout>
      <Routes>
        <Route path={ROUTES.HOME} element={<ReaderPage />} />
        <Route path={ROUTES.READER} element={<ReaderPage />} />
        <Route path={ROUTES.REPORTS} element={<ReportsPage />} />
        <Route path={ROUTES.CRAWL} element={<CrawlPage />} />
        <Route path={ROUTES.SETTINGS} element={<SettingsPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </ConsoleLayout>
  );
}

export default App;
