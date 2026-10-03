# NFX-News

[中文](README.md)

<div align="center">
  <img src="image.png" alt="NFX-News" width="200">
</div>

Go news microservices, laid out like Identity. Browsers use Fiber HTTP, modules use gRPC, and Kafka topics use the `nfxnews.*` prefix. There is **no** auth module and **no** local account table. The console calls Identity through nfx-ui. This API verifies the same JWT (`TOKEN_ISSUER=nfxidentity`). The old Fastify / TrendRadar / Node `news_server` docs are void. Trust `modules/`.

Stack (Postgres **10004**, Redis **10006**, Kafka `NAS_IP:10008`, OTLP `NAS_IP:10016`), Edge, and Identity come first. This `.env` has **no** `GRPC_PORT_AUTH` and **no** `API_GATEWAY_PREFIX`. Dial Identity at `GRPC_HOST_AUTH:GRPC_EXT_PORT_AUTH` (dev **10031**, secure **10036**). `50071` is only Identity's in-container listen.

## Modules

| Module | Role | Fiber prefix |
|--------|------|----------------|
| source | Source catalog and fetch | `/source` |
| news | Items, search, column preferences | `/news` |
| crawl | Crawl sessions | `/crawl` |
| report | Keyword DSL (`+` required, `!` excluded) and snapshots | `/report` |
| notify | Feishu / DingTalk / WeCom / Telegram webhooks. Secrets stay in `.env` | `/notify` |
| mcp | Tool list and calls | `/mcp` |

Public: `GET /source/sources`, `GET /source/sources/:id`, `POST /source/sources/:id/fetch`; `GET /news/items`, `GET /news/search`; and each module's `GET /locales/:lang`, `GET /messages/:lang`. Token required: `GET/PUT /news/preferences`; crawl `/sessions`; report `/keywords` and `/snapshots`; notify `/kinds`, `/channels`, `/deliveries`, `POST /dispatch`; mcp `GET /tools`, `POST /tools/:name`, `POST /run`.

Edge PathPrefix: dev `/dev/nfx-news`, secure `/nfx-news`. After the prefix is stripped, Fiber sees `/source` and the rest above. Browser: dev `VITE_API_URL=/dev/nfx-news`, `VITE_BASE=/dev/console/nfx-news/`, `VITE_IDENTITY_API_URL=/dev/nfx-identity`; secure drops `/dev`.

## Ports

Container HTTP is `8080`. Container gRPC runs from source `50072` through mcp `50077`. Vite is `5174`. The console port is `CONSOLE_EXTERNAL_PORT`, not gRPC.

| Module | dev HTTP / gRPC | secure HTTP / gRPC |
|--------|-----------------|--------------------|
| source | **10050 / 10051** | **10063 / 10064** |
| news | **10052 / 10053** | **10065 / 10066** |
| crawl | **10054 / 10055** | **10067 / 10068** |
| report | **10056 / 10057** | **10069 / 10070** |
| notify | **10058 / 10059** | **10071 / 10072** |
| mcp | **10060 / 10061** | **10073 / 10074** |
| console | **10062** | **10075** |

Opening `/console/nfx-news` on the LAN returns 302 to the secure console **10075**. The domain stays on HTTPS.

## Console and database

Guest: `/`, `/auth/login`, `/auth/signup`. After login the profile pages are still `/user/profile/overview|edit|identities` and `/user/settings`, not Identity's `/forger/*` yet. `App.tsx` redirects the old `/user` and `/user/overview` to `/reader`. Product pages: `/reader`, `/sources`, `/reports`, `/reports/:id`, `/crawl`, `/mcp`, `/notify`. On `/crawl`, one selected source shows that source's `home` and `intervalMs`. "All sources" shows neither.

Databases: `nfxnews_dev` / `nfxnews` / `nfxnews_diff`. Tables: `source.sources`, `news.items`, `news.profile_preferences` (primary key `account_id, profile_id`), `crawl.sessions`, `report.keywords`, `report.snapshots`, `notify.channels`, `notify.deliveries`, `mcp.tool_calls`. Kafka example: `nfxnews.news` / `nfxnews.news_poison`, and it consumes `nfxnews.source`. Pin **nfx-ui 0.36.0**.

```bash
cp .example.env .env
task proto:gen
task errors:gen-langs
task atlas:pipeline:run
task console:i
sudo docker compose -f docker-compose.dev.yml up --build
```

There is no `task atlas:pipeline:run:sh`. Change `databases/src` first, then run Atlas.

Full detail: [NFX-Documentation chapter 8](https://github.com/NebulaForgeX/NFX-Documentation/blob/main/books/en/chapter-08-nfx-news-deployment.md).
