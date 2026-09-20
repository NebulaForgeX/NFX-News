import type { ReactNode } from "react";

import { useCallback, useEffect } from "react";
import { BrowserRouter, useNavigate } from "react-router";
import { authEventEmitter, authEvents } from "nfx-ui/events";
import { clearAuth } from "nfx-ui/stores";
import { routerEventEmitter, routerEvents } from "@/events/router";
import { ROUTES } from "@/navigations";

export interface RouterProviderProps {
  children: ReactNode;
}

function RouterEventsHandler({ children }: { children: ReactNode }) {
  const navigate = useNavigate();

  const handleNavigate = useCallback(
    (payload: { to: string; replace?: boolean; state?: unknown }) => {
      navigate(payload.to, { replace: payload.replace, state: payload.state });
    },
    [navigate],
  );

  const handleNavigateBack = useCallback(() => {
    navigate(-1);
  }, [navigate]);

  const handleNavigateToLogin = useCallback(() => {
    navigate(ROUTES.LOGIN, { replace: true });
  }, [navigate]);

  const handleNavigateToDashboard = useCallback(() => {
    navigate(ROUTES.READER, { replace: true });
  }, [navigate]);

  const handleLoginSuccess = useCallback(() => {
    navigate(ROUTES.READER, { replace: true });
  }, [navigate]);

  const handleLogout = useCallback(() => {
    clearAuth();
    navigate(ROUTES.LOGIN, { replace: true });
  }, [navigate]);

  useEffect(() => {
    routerEventEmitter.on(routerEvents.NAVIGATE, handleNavigate);
    routerEventEmitter.on(routerEvents.NAVIGATE_BACK, handleNavigateBack);
    routerEventEmitter.on(routerEvents.NAVIGATE_TO_LOGIN, handleNavigateToLogin);
    routerEventEmitter.on(routerEvents.NAVIGATE_TO_DASHBOARD, handleNavigateToDashboard);
    authEventEmitter.on(authEvents.LOGIN_SUCCESS, handleLoginSuccess);
    authEventEmitter.on(authEvents.LOGOUT, handleLogout);
    return () => {
      routerEventEmitter.off(routerEvents.NAVIGATE, handleNavigate);
      routerEventEmitter.off(routerEvents.NAVIGATE_BACK, handleNavigateBack);
      routerEventEmitter.off(routerEvents.NAVIGATE_TO_LOGIN, handleNavigateToLogin);
      routerEventEmitter.off(routerEvents.NAVIGATE_TO_DASHBOARD, handleNavigateToDashboard);
      authEventEmitter.off(authEvents.LOGIN_SUCCESS, handleLoginSuccess);
      authEventEmitter.off(authEvents.LOGOUT, handleLogout);
    };
  }, [handleNavigate, handleNavigateBack, handleNavigateToDashboard, handleNavigateToLogin, handleLoginSuccess, handleLogout]);

  return <>{children}</>;
}

function RouterProvider({ children }: RouterProviderProps) {
  return (
    <BrowserRouter>
      <RouterEventsHandler>{children}</RouterEventsHandler>
    </BrowserRouter>
  );
}

export default RouterProvider;
