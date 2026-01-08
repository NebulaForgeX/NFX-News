# TrendRadar News Server

**Version**: 1.0.0  
**Language**: TypeScript (Node.js 18+)  
**Framework**: Fastify 5.x  
**Architecture**: DDD (Domain-Driven Design)

## 📋 Overview

News Server is a news aggregation service built with Fastify and TypeScript, following DDD (Domain-Driven Design) architecture principles. It provides high-performance RESTful APIs and event-driven data processing capabilities.

### Core Features

- ✅ **Microservices Architecture**: Supports multiple independent API and Pipeline services
- ✅ **DDD Design**: Clear domain models and business logic separation
- ✅ **High Performance**: Built on Fastify for excellent performance
- ✅ **Type Safety**: Full TypeScript support
- ✅ **Event-Driven**: Supports Kafka message queue
- ✅ **Database**: PostgreSQL + Drizzle ORM
- ✅ **Caching**: Redis support

## 🏗️ Architecture

### Directory Structure

```
news_server/
├── inputs/                  # Entry points (similar to Go's cmd/)
│   ├── auth/
│   │   ├── api/            # Authentication API service
│   │   └── pipeline/       # Authentication Pipeline service
│   ├── news/
│   │   ├── api/            # News API service
│   │   └── pipeline/       # News Pipeline service
│   └── source/
│       ├── api/            # Data source API service
│       └── pipeline/       # Data source Pipeline service
├── modules/                 # Domain modules (DDD)
│   ├── auth/               # Authentication domain
│   ├── news/               # News domain
│   └── source/             # Data source domain
├── packages/                # Shared packages
├── config/                  # Configuration files
├── events/                  # Event definitions
├── models/                  # Data models (Drizzle)
└── enums/                   # Enum definitions
```

### Service Types

#### 1. API Services

Provide HTTP RESTful APIs:

- `auth/api`: Authentication-related APIs
- `news/api`: News query and management APIs
- `source/api`: Data source management APIs

#### 2. Pipeline Services

Event-driven data processing services:

- `auth/pipeline`: Authentication event processing
- `news/pipeline`: News data processing
- `source/pipeline`: Data source synchronization

## 🚀 Quick Start

### 1. Requirements

- Node.js >= 18.0.0
- npm >= 9.0.0 or pnpm >= 8.0.0
- PostgreSQL >= 14.0
- Redis >= 6.0 (optional)

### 2. Install Dependencies

```bash
cd news_server
npm install
# or
pnpm install
```

### 3. Configure Environment Variables

Create `.env` file:

```bash
# Server configuration
NODE_ENV=development
PORT=3000
HOST=0.0.0.0

# Database configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=trendradar
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password

# Redis configuration (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Kafka configuration (optional)
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_CLIENT_ID=news-server

# JWT configuration (required for auth service)
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=7d
```

### 4. Database Migration

```bash
# Generate migration files
npm run db:generate

# Run migrations
npm run db:migrate

# Or use Drizzle Kit
npx drizzle-kit push
```

### 5. Start Services

#### Development Mode

```bash
# Start News API service
npm run dev:news:api

# Start News Pipeline service
npm run dev:news:pipeline

# Start Auth API service
npm run dev:auth:api

# Start Source API service
npm run dev:source:api
```

#### Production Mode

```bash
# Build project
npm run build

# Start services
npm run start:news:api
npm run start:news:pipeline
npm run start:auth:api
npm run start:source:api
```

## 📚 API Documentation

### News API Service

**Port**: Default 3000

#### Get Latest News

```http
GET /api/news/latest?limit=50&platform=zhihu
```

**Response Example**:
```json
{
  "success": true,
  "data": [
    {
      "id": "123",
      "title": "News Title",
      "platform": "zhihu",
      "url": "https://...",
      "rank": 1,
      "timestamp": "2025-01-15T10:00:00Z"
    }
  ],
  "total": 50
}
```

#### Search News

```http
POST /api/news/search
Content-Type: application/json

{
  "keywords": ["AI", "ChatGPT"],
  "platforms": ["zhihu", "weibo"],
  "dateRange": {
    "start": "2025-01-01",
    "end": "2025-01-15"
  }
}
```

### Auth API Service

**Port**: Default 3001

#### User Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}
```

#### Refresh Token

```http
POST /api/auth/refresh
Authorization: Bearer <refresh_token>
```

## 🔧 Development Guide

### Adding a New Domain Module

1. Create a new module directory in `modules/`
2. Create domain models and business logic
3. Create corresponding API or Pipeline entry point in `inputs/`

### Module Structure Example

```
modules/news/
├── domain/              # Domain layer
│   ├── entities/        # Entities
│   ├── value-objects/   # Value objects
│   └── services/        # Domain services
├── application/         # Application layer
│   ├── use-cases/       # Use cases
│   └── dto/             # Data transfer objects
├── infrastructure/      # Infrastructure layer
│   ├── repositories/    # Repository implementations
│   └── external/        # External services
└── presentation/        # Presentation layer
    ├── routes/          # Routes
    └── controllers/     # Controllers
```

### Adding New Routes

Create route files in the corresponding module's `presentation/routes/` directory:

```typescript
import { FastifyInstance } from 'fastify'

export async function newsRoutes(app: FastifyInstance) {
  app.get('/news', async (request, reply) => {
    // Handle logic
  })
}
```

### Using Drizzle ORM

```typescript
import { db } from '@infrastructure/database'
import { news } from '@models/news'

// Query
const allNews = await db.select().from(news)

// Insert
await db.insert(news).values({
  title: 'Title',
  platform: 'zhihu',
  url: 'https://...'
})
```

## 🔌 Kafka Event Processing

### Pipeline Service Event Listening

Pipeline services automatically listen to Kafka events and process them:

```typescript
// modules/news/infrastructure/events/handlers.ts
export async function handleNewsCrawled(event: NewsCrawledEvent) {
  // Process crawled news data
  await newsService.processCrawledNews(event.data)
}
```

### Event Types

- `news.crawled`: News crawling completed
- `news.processed`: News processing completed
- `source.updated`: Data source updated

## 🧪 Testing

```bash
# Run tests
npm test

# Run tests (watch mode)
npm run test:watch

# Type checking
npm run typecheck
```

## 📦 Docker Deployment

### Dockerfile Example

```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
COPY tsconfig.json ./
COPY news_server/ ./news_server/

RUN npm ci
RUN npm run build

FROM node:18-alpine

WORKDIR /app

COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/news_server/package.json ./

ENV NODE_ENV=production

CMD ["node", "dist/inputs/news/api/main.js"]
```

### docker-compose Example

```yaml
version: '3.8'

services:
  news-api:
    build: ./news_server
    ports:
      - "3000:3000"
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PASSWORD=password
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  news-pipeline:
    build: ./news_server
    command: npm run start:news:pipeline
    environment:
      - KAFKA_BOOTSTRAP_SERVERS=kafka:9092
    depends_on:
      - kafka
```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check if PostgreSQL is running
   - Verify connection parameters
   - Check firewall settings

2. **Kafka Connection Failed**
   - Check Kafka service status
   - Verify `KAFKA_BOOTSTRAP_SERVERS` configuration

3. **Type Errors**
   - Run `npm run typecheck` to check types
   - Ensure all dependencies are installed

### Viewing Logs

Use Consola for logging:

```typescript
import { consola } from 'consola'

consola.info('Info')
consola.error('Error')
consola.warn('Warning')
```

## 📄 License

MIT License
