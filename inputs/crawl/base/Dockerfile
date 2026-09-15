FROM golang:1.26-alpine AS base-builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download

FROM base-builder AS prod-build
COPY . ./
ARG APP_VERSION=dev
ENV CGO_ENABLED=0 GOOS=linux
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -trimpath -ldflags="-s -w -X main.version=${APP_VERSION}" \
  -o /out/crawl-base-service ./inputs/crawl/base

FROM alpine:latest AS prod
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
COPY --from=prod-build /out/crawl-base-service /app/crawl-base-service
COPY inputs/crawl/configuration/prod.toml /app/config.toml
COPY .env /app/.env
ENV APP_ENV=prod
ENV CONFIG_FILE=/app/config.toml
EXPOSE 8080
EXPOSE 10170
ENTRYPOINT ["/app/crawl-base-service"]

FROM base-builder AS dev
WORKDIR /app
ENV PATH="/go/bin:$PATH"
RUN go install github.com/air-verse/air@v1.61.7
COPY . ./
CMD ["air", "-c", "inputs/crawl/base/.air.toml"]
