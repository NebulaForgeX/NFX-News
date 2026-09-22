# syntax=docker/dockerfile:1
# Dev (`task run`): docker-compose.dev.yml bind-mounts host ./console (node_modules included).
# This file is the prod image only (ENV=secure / docker-compose.yml).
FROM node:22-alpine AS deps
WORKDIR /app
ARG VITE_API_URL=/nfx-news
ARG VITE_IDENTITY_API_URL=/nfx-identity
ARG VITE_BASE=/console/nfx-news/
ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_IDENTITY_API_URL=${VITE_IDENTITY_API_URL}
ENV VITE_BASE=${VITE_BASE}
COPY package.json package-lock.json ./
RUN npm install
COPY . .

FROM deps AS builder
ENV DOCKER_BUILD=1
RUN npm run build

FROM nginx:alpine AS secure
COPY --from=builder /app/dist /usr/share/nginx/html/console/nfx-news
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
