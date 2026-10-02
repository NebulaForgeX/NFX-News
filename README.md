# NFX-News

[English](README.en.md)

Go 资讯微服务。**没有本地 auth**；登录走 [NFX-Identity](https://github.com/NebulaForgeX/NFX-Identity)。

模块、路由、库表：[NFX-Documentation 第八章](https://github.com/NebulaForgeX/NFX-Documentation/blob/main/books/zh/chapter-08-nfx-news-deployment.md)。先 Stack、Edge、Identity。dev gRPC **10051–10061**（间隔为 HTTP），console **10062**。Identity 客户端拨 `GRPC_EXT_PORT_AUTH=10031`。

```bash
cp .example.env .env
task proto:gen
task errors:gen-langs
task atlas:pipeline:run:sh
task console:i
sudo docker compose -f docker-compose.dev.yml up --build
```
