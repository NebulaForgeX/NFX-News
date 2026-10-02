import type { ReactNode } from "react";

import { DataProvider as NfxDataProvider } from "nfx-ui/providers";

import { NewsRepositoriesContext, newsRepositories } from "@/apis/repositories";

export function DataProvider({ children }: { children: ReactNode }) {
  return (
    <NfxDataProvider>
      <NewsRepositoriesContext.Provider value={newsRepositories}>{children}</NewsRepositoriesContext.Provider>
    </NfxDataProvider>
  );
}
