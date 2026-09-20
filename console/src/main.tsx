import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "@radix-ui/themes/styles.css";
import "nfx-ui/themes/fonts";
import "nfx-ui/themes/index.css";

import { LanguageEnum } from "nfx-ui/enums";
import { LanguageProvider, ThemeProvider } from "nfx-ui/providers";
import { ensureDeviceIdStorage } from "nfx-ui/stores";

import { newsRepositories } from "@/apis/repositories";
import { getBuiltinI18nBundles } from "@/assets/language";
import { syncDocumentLogo } from "@/constants";
import { DataProvider, ModalProvider, QueryProvider, RouterProvider } from "@/providers";

import App from "./App";

import "./index.css";

void ensureDeviceIdStorage();

async function onLoadExtraBundles(lng: LanguageEnum) {
  const lang = lng.toString();
  try {
    const [newsErrors, newsMessages] = await Promise.all([
      newsRepositories.system.getErrorTranslations(lang),
      newsRepositories.system.getMessageTranslations(lang),
    ]);
    const bundles = [];
    if (newsErrors) bundles.push({ namespace: "errors", bundle: newsErrors as Record<string, unknown> });
    if (newsMessages) bundles.push({ namespace: "messages", bundle: newsMessages as Record<string, unknown> });
    return bundles.length > 0 ? bundles : null;
  } catch (error) {
    console.error("Failed to load product translation bundles", error);
    return null;
  }
}

function bootstrap() {
  createRoot(document.getElementById("root")!).render(
    <StrictMode>
      <QueryProvider>
        <LanguageProvider fallbackLng={LanguageEnum.ZH} getBuiltinBundles={getBuiltinI18nBundles} onLoadExtraBundles={onLoadExtraBundles}>
          <ThemeProvider onAppearanceChange={syncDocumentLogo}>
            <DataProvider>
              <RouterProvider>
                <ModalProvider>
                  <App />
                </ModalProvider>
              </RouterProvider>
            </DataProvider>
          </ThemeProvider>
        </LanguageProvider>
      </QueryProvider>
    </StrictMode>,
  );
}

void bootstrap();
