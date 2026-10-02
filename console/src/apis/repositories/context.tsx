import { createContext, useContext } from "react";

import { newsRepositories } from "./newsRepositories";

export type NewsRepositories = typeof newsRepositories;

export const NewsRepositoriesContext = createContext<NewsRepositories>(newsRepositories);

export function useNewsRepositories() {
  return useContext(NewsRepositoriesContext);
}
