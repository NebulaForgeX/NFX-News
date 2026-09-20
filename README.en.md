# NFX-News

[中文](README.md)

Go news microservices. **No local auth**; login is [NFX-Identity](https://github.com/NebulaForgeX/NFX-Identity).

Modules, routes, tables: [NFX-Documentation chapter 8](https://github.com/NebulaForgeX/NFX-Documentation/blob/main/books/en/chapter-08-nfx-news-deployment.md). Stack, Edge, and Identity first. Host gRPC **10204–10210**, console **10211**. `GRPC_PORT_AUTH=50071` is the Identity client.

```bash
cp .example.env .env
task proto:gen
task errors:gen-langs
task atlas:pipeline:run:sh
task console:i
sudo docker compose -f docker-compose.dev.yml up --build
```
