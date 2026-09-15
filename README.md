# NFX-News

中文 / English

Identity-shaped Go microservices for news reading, crawling, keyword reports, MCP, and multi-channel notify. Browser talks **HTTP REST**; services talk **native gRPC**; Kafka sits on **NFX-Stack**. Login is **NFX-Identity** via `nfx-ui` AuthRepository + product `connections/auth` gRPC.

以 [NFX-Identity](../NFX-Identity) 为目录/协议模板：`console/` 对浏览器走 Fiber HTTP；模块间走原生 gRPC；异步用 Kafka `nfxnews.*`。登录走 Identity（不再有本地 auth 模块）。数据面全部接已运行的 NFX-Stack。

## Layout / 目录

```
NFX-News/
├── console/                 # React + Vite SPA (login, select-profile, reader, reports, crawl, settings)
├── modules/{source,news,crawl,report,notify,mcp,system}/
├── inputs/{module}/{api,connection,pipeline,messaging,base}/
├── connections/auth         # Identity gRPC client (copied gen; do not hand-edit)
├── protos/                  # buf src → gen (module nfxnews)
├── databases/               # Atlas SQL (`news.profile_preferences` keyed by account_id+profile_id)
├── events/ pkgs/ errors/
├── docker-compose.yml
├── docker-compose.dev.yml
└── Taskfile.yml
```

## Ports / 端口

Avoid Stack MinIO Console `10188`. News Traefik HTTP is `10178`.

| Use | Variable | Port |
|-----|----------|------|
| Identity auth gRPC | `GRPC_PORT_AUTH` | 10156 (Identity) |
| gRPC source…system | `GRPC_PORT_*` | 10171–10177 |
| Traefik HTTP/HTTPS | `TRAEFIK_HTTP/HTTPS_PORT` | 10178 / 10179 |
| Console | `CONSOLE_EXTERNAL_PORT` | 10190 |
| Traefik dashboard | `TRAEFIK_DASHBOARD_PORT` | 10191 |

HTTP prefixes: `/source` `/news` `/crawl` `/report` `/notify` `/mcp` `/system`.

Console: `VITE_API_URL=http://127.0.0.1:10178`, `VITE_IDENTITY_API_URL=http://127.0.0.1:10166`.

Token secret/issuer **must match NFX-Identity** (`TOKEN_ISSUER=nfxidentity`) so product APIs can verify user JWTs locally.

Stack: Postgres `192.168.1.64:10105`, Redis `10181`, Kafka `10183`, OTEL `10192`. Databases: `nfxnews_dev` / `nfxnews` / `nfxnews_diff`.

## Run / 运行

```bash
cp .example.env .env   # fill Stack passwords + the same TOKEN_SECRET_KEY as Identity
# Create nfxnews_dev / nfxnews / nfxnews_diff on Stack Postgres, then:
task proto:gen
task errors:gen-langs
task atlas:pipeline:run:sh
task console:i
task console                  # Vite on VITE_PORT (5174)
sudo docker compose -f docker-compose.dev.yml up --build
```

MCP Streamable HTTP: `POST/GET http://127.0.0.1:10178/mcp/protocol`. REST tools: `/mcp/tools`, `/mcp/run`.

Console Docker uses `additional_contexts.nfx-ui: ../NFX-UI` and `"nfx-ui": "file:../../NFX-UI"`. Build NFX-UI (`npm run build` in that repo) first.

## Modules / 模块

1. **source** — catalog + Go getters (RSS / RSSHub / major sites).
2. **news** — persist items, search, Redis source snapshot TTL, column preferences; consumes `nfxnews.source`.
3. **crawl** — sessions + scheduled FetchSource.
4. **report** — keyword DSL in Postgres (`+required` / `!exclude`), daily/current/incremental snapshots.
5. **notify** — feishu/dingtalk/wework/telegram/ntfy/bark/slack; webhooks only in `.env`.
6. **mcp** — MCP protocol + gRPC to news/report/source.
7. **system** — health / bootstrap state.

Auth is not a News module. Consoles log in through Identity; product APIs verify the user JWT and call Identity `EnsureOwnedProfile` / `HasForgerRole` over gRPC.

## NFX-UI

Until publish: `"nfx-ui": "file:../../NFX-UI"` only. Do not add a registry version alongside `file:`.
