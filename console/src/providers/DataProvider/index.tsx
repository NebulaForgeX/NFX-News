import type { ReactNode } from "react";

import { NewsRepositoriesContext, newsRepositories } from "@/apis/repositories";

export function NewsDataProvider({ children }: { children: ReactNode }) {
  return <NewsRepositoriesContext.Provider value={newsRepositories}>{children}</NewsRepositoriesContext.Provider>;
}
