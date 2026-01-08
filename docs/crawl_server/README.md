# TrendRadar Crawl Server

**Version**: 3.4.0  
**Language**: Python 3.8+  
**Architecture**: Microservices architecture with scheduled tasks and event-driven support

## 📋 Overview

Crawl Server is the core crawler service of the TrendRadar project, responsible for fetching trending content from multiple news platforms and storing it in the database. Built with microservices architecture, it supports:

- ✅ Scheduled automatic crawling
- ✅ Kafka event-driven crawling
- ✅ Connection pool management for multiple data sources (PostgreSQL, Redis, Kafka)
- ✅ Graceful shutdown mechanism
- ✅ Comprehensive error handling and logging

## 🏗️ Architecture

### Core Components

```
crawl_server/
├── main.py                  # Main entry point (server mode)
├── crawl_task.py            # Crawl task execution logic
├── configs/                 # Configuration management
│   ├── config.py           # Config loader
│   ├── crawl_config.py     # Crawler configuration
│   └── database_config.py  # Database configuration
├── connections.py           # Connection pool initialization
├── controllers/             # Controller layer
│   ├── crawl_controller.py # Crawler controller
│   ├── data_controller.py  # Data controller
│   └── frequency_controller.py # Keyword controller
├── services/                # Business logic layer
├── repositories/            # Data access layer
├── models/                  # Data models
└── routers/                 # API routes (if needed)
```

### Running Modes

#### 1. Server Mode (Recommended)

Continuously running with scheduled tasks and event listening:

```bash
python -m crawl_server.main
```

**Features**:
- Scheduled crawl tasks (configurable interval)
- Kafka event listening (`operation.crawl`, `operation.clear`)
- Graceful shutdown (responds to SIGTERM/SIGINT)

#### 2. One-time Execution Mode

Execute a single crawl task and exit:

```bash
python -m crawl_server.crawl_task --count 1
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SCHEDULE_MINUTES` | Scheduled task interval (minutes) | `30` | No |
| `KAFKA_BOOTSTRAP_SERVERS` | Kafka server address | - | No |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` | Yes |
| `POSTGRES_PORT` | PostgreSQL port | `5432` | Yes |
| `POSTGRES_DB` | Database name | `trendradar` | Yes |
| `POSTGRES_USER` | Database user | `postgres` | Yes |
| `POSTGRES_PASSWORD` | Database password | - | Yes |
| `REDIS_HOST` | Redis host | `localhost` | No |
| `REDIS_PORT` | Redis port | `6379` | No |

### Configuration File

Configuration file location: `config/config.yaml`

Main configuration items:

```yaml
crawler:
  enable_crawler: true
  schedule_minutes: 30
  
platforms:
  - id: "zhihu"
    name: "Zhihu"
  - id: "weibo"
    name: "Weibo"
  # ... more platforms

database:
  host: localhost
  port: 5432
  dbname: trendradar
  user: postgres
  password: ${POSTGRES_PASSWORD}
```

## 🚀 Quick Start

### 1. Install Dependencies

```bash
cd crawl_server
pip install -r requirements.txt
```

### 2. Configure Environment

Create `.env` file or set environment variables:

```bash
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_DB=trendradar
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=your_password
export SCHEDULE_MINUTES=30
```

### 3. Start Service

```bash
# Server mode (recommended)
python -m crawl_server.main

# Or use uvicorn (if exposing HTTP API)
uvicorn crawl_server.main:app --host 0.0.0.0 --port 8001
```

## 📊 Features

### 1. Multi-platform Crawling

Supports fetching trending content from the following platforms:

- Zhihu (知乎)
- Weibo (微博)
- Douyin (抖音)
- Baidu (百度)
- Bilibili (哔哩哔哩)
- Toutiao (今日头条)
- Tieba (贴吧)
- The Paper (澎湃新闻)
- Yicai (财联社)
- Ifeng (凤凰网)
- Wallstreetcn (华尔街见闻)

### 2. Event-driven Crawling

Supports triggering crawls via Kafka events:

**Event Types**:
- `operation.crawl`: Trigger crawl task
- `operation.clear`: Clear old data

**Event Format**:
```json
{
  "type": "operation.crawl",
  "payload": {
    "count": 1,
    "platforms": ["zhihu", "weibo"]
  }
}
```

### 3. Connection Pool Management

Automatically manages the following connection pools:

- **PostgreSQL**: Database connection pool (using SQLAlchemy)
- **Redis**: Cache connection pool
- **Kafka**: Producer/consumer connections

### 4. Graceful Shutdown

Supports signal handling to ensure data integrity:

```bash
# Send stop signal
kill -SIGTERM <pid>

# Or use Ctrl+C
```

The service will:
1. Stop accepting new tasks
2. Complete currently executing tasks
3. Close all connections
4. Exit

## 🔧 Development Guide

### Adding a New Platform

1. Create a new crawler class in `services/crawlers/`
2. Add platform configuration in `config/config.yaml`
3. Register the new crawler in `controllers/crawl_controller.py`

### Customizing Crawl Logic

Modify the `crawl()` method in `controllers/crawl_controller.py`.

### Adding a New Controller

Create a new controller file in the `controllers/` directory and register it in `connections.py`.

## 📝 API Documentation (if HTTP mode is enabled)

If HTTP API is enabled (using FastAPI), you can access:

- Swagger UI: `http://localhost:8001/docs`
- ReDoc: `http://localhost:8001/redoc`

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check if PostgreSQL is running
   - Verify connection parameters
   - Check firewall settings

2. **Kafka Connection Failed**
   - Check if Kafka service is started
   - Verify `KAFKA_BOOTSTRAP_SERVERS` configuration

3. **Crawling Failed**
   - Check log files
   - Check network connection
   - Verify if target websites are accessible

### Viewing Logs

Logs are output to stdout, it's recommended to use log collection tools (e.g., journald, logrotate):

```bash
# View real-time logs
python -m crawl_server.main 2>&1 | tee crawl.log

# Using systemd (production environment)
journalctl -u trendradar-crawl -f
```

## 📦 Docker Deployment

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY crawl_server/ ./crawl_server/
COPY config/ ./config/
COPY requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt

CMD ["python", "-m", "crawl_server.main"]
```

```bash
docker build -t trendradar-crawl:latest .
docker run -d \
  --name trendradar-crawl \
  -e POSTGRES_HOST=postgres \
  -e POSTGRES_PASSWORD=password \
  -v $(pwd)/config:/app/config:ro \
  trendradar-crawl:latest
```

## 🔗 Related Services

- **news_server**: News aggregation API service (TypeScript/Fastify)
- **web_server**: Web report service (Python/FastAPI)
- **mcp_server**: AI analysis service (Python/FastMCP)

## 📄 License

GPL-3.0 License
