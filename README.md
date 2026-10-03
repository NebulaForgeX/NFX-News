# NFX-News

[English](README.en.md)

<div align="center">
  <img src="image.png" alt="NFX-News" width="200">
</div>

Go 资讯微服务，目录对齐 Identity。浏览器走 Fiber HTTP，模块之间走 gRPC，Kafka 前缀 `nfxnews.*`。**没有** auth 模块，也 **没有**本地账号表。Console 用 nfx-ui 调 Identity；本仓 API 校验同一套 JWT（`TOKEN_ISSUER=nfxidentity`）。旧的 Fastify / TrendRadar / Node `news_server` 文档已经作废，以 `modules/` 为准。

先有 Stack（Postgres **10004**、Redis **10006**、Kafka `NAS_IP:10008`、OTLP `NAS_IP:10016`）、Edge、Identity。本仓 `.env` **没有** `GRPC_PORT_AUTH`，也 **没有** `API_GATEWAY_PREFIX`。拨 Identity 是 `GRPC_HOST_AUTH:GRPC_EXT_PORT_AUTH`（dev **10031**，secure **10036**）。`50071` 只是 Identity 容器内的 listen。

## 模块

| 模块 | 做什么 | Fiber 前缀 |
|------|--------|------------|
| source | 源目录和拉取 | `/source` |
| news | 条目、搜索、栏目偏好 | `/news` |
| crawl | 抓取会话 | `/crawl` |
| report | 关键词 DSL（`+` 必须、`!` 排除）和快照 | `/report` |
| notify | 飞书 / 钉钉 / 企微 / Telegram webhook，密钥只放 `.env` | `/notify` |
| mcp | 工具列表和调用 | `/mcp` |

公开：`GET /source/sources`、`GET /source/sources/:id`、`POST /source/sources/:id/fetch`；`GET /news/items`、`GET /news/search`；各模块的 `GET /locales/:lang`、`GET /messages/:lang`。要 Token：`GET/PUT /news/preferences`；crawl 的 `/sessions`；report 的 `/keywords`、`/snapshots`；notify 的 `/kinds`、`/channels`、`/deliveries`、`POST /dispatch`；mcp 的 `GET /tools`、`POST /tools/:name`、`POST /run`。

Edge PathPrefix：dev `/dev/nfx-news`，secure `/nfx-news`，去掉前缀后 Fiber 看到的是上面的 `/source` 等。浏览器：dev `VITE_API_URL=/dev/nfx-news`、`VITE_BASE=/dev/console/nfx-news/`、`VITE_IDENTITY_API_URL=/dev/nfx-identity`；secure 去掉 `/dev`。

## 端口

容器 HTTP `8080`。容器 gRPC 从 source 的 `50072` 排到 mcp 的 `50077`。Vite `5174`。console 是 `CONSOLE_EXTERNAL_PORT`，不是 gRPC。

| 模块 | dev HTTP / gRPC | secure HTTP / gRPC |
|------|-----------------|--------------------|
| source | **10050 / 10051** | **10063 / 10064** |
| news | **10052 / 10053** | **10065 / 10066** |
| crawl | **10054 / 10055** | **10067 / 10068** |
| report | **10056 / 10057** | **10069 / 10070** |
| notify | **10058 / 10059** | **10071 / 10072** |
| mcp | **10060 / 10061** | **10073 / 10074** |
| console | **10062** | **10075** |

局域网打开 `/console/nfx-news` 会 302 到 secure console **10075**。域名仍走 HTTPS。

## Console 和库

访客：`/`、`/auth/login`、`/auth/signup`。登录后资料页仍是 `/user/profile/overview|edit|identities` 和 `/user/settings`，还没改成 Identity 的 `/forger/*`。`App.tsx` 把旧的 `/user`、`/user/overview` 重定向到 `/reader`。业务页：`/reader`、`/sources`、`/reports`、`/reports/:id`、`/crawl`、`/mcp`、`/notify`。`/crawl` 选中单个源时显示该源的 `home` 和 `intervalMs`，选「全部来源」时不显示。

库：`nfxnews_dev` / `nfxnews` / `nfxnews_diff`。表：`source.sources`、`news.items`、`news.profile_preferences`（主键 `account_id, profile_id`）、`crawl.sessions`、`report.keywords`、`report.snapshots`、`notify.channels`、`notify.deliveries`、`mcp.tool_calls`。Kafka 示例：`nfxnews.news` / `nfxnews.news_poison`，并消费 `nfxnews.source`。`nfx-ui` 钉 **0.36.0**。

```bash
cp .example.env .env
task proto:gen
task errors:gen-langs
task atlas:pipeline:run
task console:i
sudo docker compose -f docker-compose.dev.yml up --build
```

没有 `task atlas:pipeline:run:sh`。改表先改 `databases/src`，再跑 Atlas。

详细信息见 [NFX-Documentation 第八章](https://github.com/NebulaForgeX/NFX-Documentation/blob/main/books/zh/chapter-08-nfx-news-deployment.md)。
