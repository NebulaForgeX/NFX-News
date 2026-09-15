import type { ReactNode } from "react";

import { BrowserRouter } from "react-router";

import { useRouterEvents } from "./hooks/useRouterEvents";

function RouterEventsHandler({ children }: { children: ReactNode }) {
  useRouterEvents();
  return <>{children}</>;
}

export function RouterProvider({ children }: { children: ReactNode }) {
  return (
    <BrowserRouter>
      <RouterEventsHandler>{children}</RouterEventsHandler>
    </BrowserRouter>
  );
}

export default RouterProvider;
