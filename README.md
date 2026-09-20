# NFX-News

> 部署、网络、配置与安全的详细说明见 [NFX-Documentation](https://github.com/NebulaForgeX/NFX-Documentation)。
> Deploy, network, config, and security: [NFX-Documentation](https://github.com/NebulaForgeX/NFX-Documentation).

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

HTTP 走 **NFX-Edge**（先启动 Edge）。本仓不跑 Traefik。Avoid Stack MinIO Console `10188`.

容器内 `GRPC_PORT_*` 为 50071+（对齐 PulsoLink-API）；主机只映射 `GRPC_EXT_*`（10200 起）。

| Use | Variable | Port |
|-----|----------|------|
| Identity auth gRPC (container / via nfx-edge) | `GRPC_PORT_AUTH` | 50071 (`GRPC_HOST_AUTH=NFX-Identity-Auth-Base-Dev`) |
| gRPC source…system (container) | `GRPC_PORT_*` | 50072–50078 |
| gRPC source…system (host) | `GRPC_EXT_PORT_*` | 10204–10210 |
| Console (optional host map) | `CONSOLE_EXTERNAL_PORT` | 10211 |

HTTP prefixes via Edge Host `TRAEFIK_API_HOST`: `/source` `/news` `/crawl` `/report` `/notify` `/mcp` `/system`.

Console: `VITE_API_URL=https://news-api.example.com`, `VITE_IDENTITY_API_URL=https://api.example.com`.

Token secret/issuer **must match NFX-Identity** (`TOKEN_ISSUER=nfxidentity`) so product APIs can verify user JWTs locally.

Stack: Postgres `192.168.1.64:10104`, Redis `10106`, Kafka `10108`, OTEL OTLP gRPC `10116`. Databases: `nfxnews_dev` / `nfxnews` / `nfxnews_diff`.

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

MCP Streamable HTTP: Edge Host `TRAEFIK_API_HOST` + PathPrefix `/mcp`（`POST/GET /mcp/protocol`）。REST tools: `/mcp/tools`, `/mcp/run`.

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
