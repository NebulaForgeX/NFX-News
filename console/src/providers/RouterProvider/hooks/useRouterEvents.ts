import { useCallback, useEffect } from "react";
import { useNavigate } from "react-router";
import { clearAuth } from "nfx-ui/stores";

import { authEventEmitter, authEvents } from "@/events/auth";
import { routerEventEmitter, routerEvents } from "@/events/router";
import { ROUTES } from "@/navigations";

export function useRouterEvents() {
  const navigate = useNavigate();

  const handleNavigate = useCallback((payload: { to: string; replace?: boolean; state?: unknown }) => {
    navigate(payload.to, { replace: payload.replace, state: payload.state });
  }, [navigate]);

  const handleLoginSuccess = useCallback(() => {
    navigate(ROUTES.SELECT_PROFILE, { replace: true });
  }, [navigate]);

  const handleLogout = useCallback(() => {
    clearAuth();
    navigate(ROUTES.LOGIN, { replace: true });
  }, [navigate]);

  useEffect(() => {
    routerEventEmitter.on(routerEvents.NAVIGATE, handleNavigate as (...args: unknown[]) => void);
    routerEventEmitter.on(routerEvents.NAVIGATE_TO_LOGIN, () => navigate(ROUTES.LOGIN, { replace: true }));
    routerEventEmitter.on(routerEvents.NAVIGATE_TO_DASHBOARD, () => navigate(ROUTES.READER, { replace: true }));
    routerEventEmitter.on(routerEvents.NAVIGATE_TO_HOME, () => navigate(ROUTES.HOME, { replace: true }));
    authEventEmitter.on(authEvents.LOGIN_SUCCESS, handleLoginSuccess);
    authEventEmitter.on(authEvents.LOGOUT, handleLogout);
    return () => {
      routerEventEmitter.off(routerEvents.NAVIGATE, handleNavigate as (...args: unknown[]) => void);
      authEventEmitter.off(authEvents.LOGIN_SUCCESS, handleLoginSuccess);
      authEventEmitter.off(authEvents.LOGOUT, handleLogout);
    };
  }, [handleNavigate, handleLoginSuccess, handleLogout, navigate]);
}
