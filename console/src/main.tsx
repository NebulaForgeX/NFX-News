import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "@radix-ui/themes/styles.css";
import "nfx-ui/themes/fonts";
import "nfx-ui/themes/styles.css";

import { LanguageEnum } from "nfx-ui/enums";
import { LanguageProvider, ThemeProvider, ModalProvider, DataProvider } from "nfx-ui/providers";
import { LayoutProvider } from "nfx-ui/layouts";

import "./index.css";

import { getBuiltinI18nBundles } from "@/assets/languages/i18nResources";
import { newsRepositories } from "@/apis/repositories";
import { NewsDataProvider, QueryProvider, RouterProvider } from "@/providers";

import App from "./App.tsx";

async function onLoadExtraBundles(lng: LanguageEnum) {
  try {
    const bundle = await newsRepositories.system.getErrorTranslations(lng);
    return { namespace: "errors", bundle: bundle as Record<string, unknown> };
  } catch {
    return null;
  }
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryProvider>
      <LanguageProvider getBuiltinBundles={getBuiltinI18nBundles} fallbackLng={LanguageEnum.ZH} onLoadExtraBundles={onLoadExtraBundles}>
        <ThemeProvider>
          <LayoutProvider>
            <DataProvider>
              <NewsDataProvider>
                <RouterProvider>
                  <ModalProvider>
                    <App />
                  </ModalProvider>
                </RouterProvider>
              </NewsDataProvider>
            </DataProvider>
          </LayoutProvider>
        </ThemeProvider>
      </LanguageProvider>
    </QueryProvider>
  </StrictMode>,
);
