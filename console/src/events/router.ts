import type { EventNamesOf } from "nfx-ui/events";
import { defineEvents, EventEmitter } from "nfx-ui/events";
import { singleton } from "nfx-ui/utils";

export const routerEvents = defineEvents({
  NAVIGATE: "ROUTER:NAVIGATE",
  NAVIGATE_TO_LOGIN: "ROUTER:NAVIGATE_TO_LOGIN",
  NAVIGATE_TO_DASHBOARD: "ROUTER:NAVIGATE_TO_DASHBOARD",
  NAVIGATE_TO_HOME: "ROUTER:NAVIGATE_TO_HOME",
});

type RouterEvent = EventNamesOf<typeof routerEvents>;

interface NavigatePayload {
  to: string;
  replace?: boolean;
  state?: unknown;
}

class RouterEventEmitter extends EventEmitter<RouterEvent> {
  constructor() {
    super(routerEvents);
  }

  navigate(payload: NavigatePayload) {
    this.emit(routerEvents.NAVIGATE, payload);
  }

  navigateToLogin() {
    this.emit(routerEvents.NAVIGATE_TO_LOGIN);
  }

  navigateToDashboard() {
    this.emit(routerEvents.NAVIGATE_TO_DASHBOARD);
  }
}

export const routerEventEmitter = new (singleton(RouterEventEmitter))();
