#!/usr/bin/env python3
from pathlib import Path
import shutil

ROOT = Path("/volume1/Projects/NebulaForgeX/NFX-News/console")
IDENT = Path("/volume1/Projects/NebulaForgeX/NFX-Identity/console")


def w(rel: str, content: str) -> None:
    p = ROOT / rel
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(content if content.endswith("\n") else content + "\n")


def main() -> None:
    shutil.copy(IDENT / "package.json", ROOT / "package.json")
    pkg = (ROOT / "package.json").read_text().replace('"name": "nfx-console"', '"name": "nfx-news-console"')
    pkg = pkg.replace('"prepare": "husky install"', '"prepare": "echo skip-husky"')
    (ROOT / "package.json").write_text(pkg)

    shutil.copy(IDENT / "tsconfig.json", ROOT / "tsconfig.json")
    shutil.copy(IDENT / "tsconfig.app.json", ROOT / "tsconfig.app.json")
    shutil.copy(IDENT / "tsconfig.node.json", ROOT / "tsconfig.node.json")
    shutil.copy(IDENT / "nginx.conf", ROOT / "nginx.conf")
    shutil.copy(IDENT / "src/index.css", ROOT / "src/index.css")
    shutil.copy(IDENT / "src/App.module.css", ROOT / "src/App.module.css")
    shutil.copy(IDENT / "src/vite-env.d.ts", ROOT / "src/vite-env.d.ts")
    shutil.copytree(IDENT / "public", ROOT / "public", dirs_exist_ok=True)
    news = Path("/volume1/Projects/NebulaForgeX/NFX-News")
    shutil.copy(news / "logo_g.ico", ROOT / "public/logo_g.ico")
    shutil.copy(news / "image.png", ROOT / "public/image.png")

    w(
        "Dockerfile",
        """# syntax=docker/dockerfile:1
FROM node:22-alpine AS builder
WORKDIR /app
ARG VITE_API_URL=http://127.0.0.1:10188
ENV VITE_API_URL=${VITE_API_URL}
COPY package.json package-lock.json* ./
COPY --from=nfx-ui . /NFX-UI
RUN npm ci || npm install
COPY . .
RUN npm run build

FROM nginx:alpine AS prod
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
""",
    )
    w(
        ".env.dev",
        """VITE_PORT=5174
VITE_API_URL=http://127.0.0.1:10188
""",
    )
    w(
        ".env.prod",
        """VITE_API_URL=http://127.0.0.1:10188
""",
    )
    w(
        "vite.config.ts",
        """import path from "path";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const port = Number(env.VITE_PORT) || 5174;
  return {
    plugins: [react()],
    base: "/",
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    css: {
      modules: {
        localsConvention: "camelCase",
        generateScopedName: "[name]__[local]___[hash:base64:5]",
      },
    },
    server: {
      port,
      host: "0.0.0.0",
      open: true,
    },
    preview: { port, host: "0.0.0.0" },
  };
});
""",
    )
    w(
        "index.html",
        """<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>NFX News</title>
    <meta name="description" content="NFX News reader, reports, crawl and notifications." />
    <link rel="icon" type="image/x-icon" href="/logo_g.ico" />
    <link rel="apple-touch-icon" href="/image.png" />
    <meta name="theme-color" content="#ef4444" />
    <meta name="mobile-web-app-capable" content="yes" />
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
""",
    )
    print("console scaffold copied")


if __name__ == "__main__":
    main()
