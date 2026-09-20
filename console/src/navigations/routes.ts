import { createRouter, defineRouter } from "@/utils";

const routeMap = defineRouter({
  HOME: "/",
  LOGIN: "/auth/login",
  SIGNUP: "/auth/signup",

  USER: "/user",
  USER_OVERVIEW: "/user/overview",
  PROFILE: "/user/profile",
  USER_PROFILE_OVERVIEW: "/user/profile/overview",
  USER_PROFILE_EDIT: "/user/profile/edit",
  USER_PROFILE_IDENTITIES: "/user/profile/identities",
  USER_SETTINGS: "/user/settings",

  READER: "/reader",
  REPORTS: "/reports",
  REPORT_DETAIL: "/reports/:id",
  CRAWL: "/crawl",
  MCP: "/mcp",
  NOTIFY: "/notify",
  SOURCES: "/sources",
  SYSTEM: "/system",
});

const { ROUTES, matchRoute, isActiveRoute, buildPath } = createRouter(routeMap);

export type RouteKey = keyof typeof ROUTES;
export { ROUTES, matchRoute, isActiveRoute, buildPath };
