<div align="center" id="nfx-news">

<a href="https://github.com/NebulaForgeX/NFX-News" title="NFX-News">
  <img src="logo_g.png" alt="NFX-News Logo" width="120" height="120">
</a>

<h1>NFX-News</h1>

🚀 **Your Smart Trending News Aggregation Platform** — Deploy in minutes, track what matters

[![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg?style=flat-square)](LICENSE)
[![Python](https://img.shields.io/badge/Python-3.8+-blue.svg?style=flat-square&logo=python)](https://www.python.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-blue.svg?style=flat-square&logo=typescript)](https://www.typescriptlang.org/)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.115+-green.svg?style=flat-square&logo=fastapi)](https://fastapi.tiangolo.com/)
[![Fastify](https://img.shields.io/badge/Fastify-5.6+-green.svg?style=flat-square&logo=fastify)](https://www.fastify.io/)

</div>

<div align="center">

**[English](README.md)** | **中文**

</div>

---

## 📋 Overview

**NFX-News** is a comprehensive trending news aggregation platform built with microservices architecture. It automatically crawls trending content from multiple news platforms, processes and analyzes the data, and provides various interfaces for viewing and interacting with trending news.

### ✨ Key Features

- 🔍 **Multi-Platform Aggregation**: Automatically crawl trending content from 11+ major platforms (Zhihu, Weibo, Douyin, Bilibili, Baidu, etc.)
- 🤖 **AI-Powered Analysis**: Natural language query and analysis using MCP (Model Context Protocol)
- 📊 **Smart Filtering**: Keyword-based filtering with advanced syntax (required words, filter words, count limits)
- 🌐 **Multiple Interfaces**: RESTful API, HTML reports, and AI conversational interface
- ⚡ **High Performance**: Built with FastAPI, Fastify, and optimized data processing
- 🏗️ **Microservices Architecture**: Independent, scalable services for different functions
- 🔔 **Event-Driven**: Kafka-based event system for real-time data processing

## 🏗️ Architecture

NFX-News consists of four main services:

```
NFX-News/
├── crawl_server/      # Python - News crawler service
├── news_server/       # TypeScript - News aggregation API service  
├── web_server/        # Python - HTML report web service
├── mcp_server/        # Python - AI analysis service (MCP)
├── config/            # Configuration files
├── output/            # Crawled news data storage
└── docs/              # Documentation
    ├── crawl_server/  # Detailed docs for crawl server
    ├── news_server/   # Detailed docs for news server
    ├── web_server/    # Detailed docs for web server
    └── mcp_server/    # Detailed docs for MCP server
```

### Service Overview

| Service | Language | Framework | Purpose | Port |
|---------|----------|-----------|---------|------|
| **Crawl Server** | Python | Custom | Crawl trending news from multiple platforms | - |
| **News Server** | TypeScript | Fastify | RESTful API for news data access | 3000+ |
| **Web Server** | Python | FastAPI | HTML report viewing interface | 10199 |
| **MCP Server** | Python | FastMCP | AI analysis and query interface | 3333 (HTTP) |

## 🚀 Quick Start

### Prerequisites

- **Python 3.8+** for Python services
- **Node.js 18+** for TypeScript services
- **PostgreSQL 14+** for data storage
- **Redis 6+** (optional) for caching
- **Kafka** (optional) for event streaming

### 1. Clone Repository

```bash
git clone https://github.com/NebulaForgeX/NFX-News.git
cd NFX-News
```

### 2. Configure Environment

Create `.env` files for each service:

**Crawl Server** (`crawl_server/.env`):
```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=nfx_news
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
SCHEDULE_MINUTES=30
```

**News Server** (`news_server/.env`):
```bash
NODE_ENV=development
PORT=3000
HOST=0.0.0.0
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=nfx_news
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
```

**Web Server** (`web_server/.env`):
```bash
PROJECT_ROOT=/path/to/NFX-News
HOST=0.0.0.0
PORT=10199
DEBUG=false
```

**MCP Server** (`mcp_server/.env`):
```bash
NFX_NEWS_PROJECT_ROOT=/path/to/NFX-News
```

### 3. Install Dependencies

```bash
# Install Python dependencies
cd crawl_server && pip install -r requirements.txt && cd ..
cd web_server && pip install -r requirements.txt && cd ..
cd mcp_server && pip install -r requirements.txt && cd ..

# Install Node.js dependencies
cd news_server && npm install && cd ..
```

### 4. Initialize Database

```bash
# Run database migrations (if applicable)
cd news_server
npm run db:migrate
cd ..
```

### 5. Start Services

```bash
# Terminal 1: Start Crawl Server
cd crawl_server
python -m crawl_server.main

# Terminal 2: Start News Server
cd news_server
npm run dev:news:api

# Terminal 3: Start Web Server
cd web_server
python -m web_server.main

# Terminal 4: Start MCP Server (optional, for AI features)
cd mcp_server
python -m mcp_server.server --transport http --port 3333
```

### 6. Access Services

- **Web Interface**: http://localhost:10199/report
- **News API**: http://localhost:3000/api/news
- **MCP Server**: http://localhost:3333/mcp (if HTTP mode enabled)
- **API Documentation**:
  - Web Server: http://localhost:10199/docs (if debug enabled)
  - News Server: Check service-specific documentation

## 📚 Documentation

Detailed documentation for each service:

### Core Services

- 📖 **[Crawl Server Documentation](docs/crawl_server/README.md)** - News crawler service
  - Architecture and configuration
  - Multi-platform crawling
  - Event-driven crawling
  - Docker deployment

- 📖 **[News Server Documentation](docs/news_server/README.md)** - News aggregation API
  - DDD architecture
  - RESTful API endpoints
  - Kafka event processing
  - Database schema

- 📖 **[Web Server Documentation](docs/web_server/README.md)** - HTML report service
  - MVC architecture
  - Report generation
  - API endpoints

- 📖 **[MCP Server Documentation](docs/mcp_server/README.md)** - AI analysis service
  - MCP protocol implementation
  - Natural language queries
  - Client integration (Claude Desktop, Cursor, etc.)
  - Tool usage examples

## 🎯 Features

### 1. Multi-Platform Crawling

Automatically crawl trending content from:
- Zhihu (知乎)
- Weibo (微博)
- Douyin (抖音)
- Bilibili (哔哩哔哩)
- Baidu (百度)
- Toutiao (今日头条)
- Tieba (贴吧)
- The Paper (澎湃新闻)
- Yicai (财联社)
- Ifeng (凤凰网)
- Wallstreetcn (华尔街见闻)

### 2. Smart Filtering

Advanced keyword filtering syntax:
- **Normal keywords**: Basic matching
- **Required words** (`+keyword`): Narrow scope
- **Filter words** (`!keyword`): Exclude noise
- **Count limits** (`@number`): Control display count

### 3. AI-Powered Analysis

Query and analyze news data using natural language:
- "Get the latest news about AI from Zhihu"
- "Analyze the popularity trend of 'Bitcoin' over the past 30 days"
- "Search for news related to 'Tesla' and 'Musk'"

### 4. Multiple Interfaces

- **RESTful API**: Programmatic access to news data
- **HTML Reports**: Beautiful, mobile-responsive web interface
- **AI Interface**: Natural language query interface via MCP

## 🔧 Configuration

### Main Configuration

Edit `config/config.yaml`:

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

report:
  mode: "current"  # daily, current, incremental

database:
  host: localhost
  port: 5432
  dbname: nfx_news
  user: postgres
  password: ${POSTGRES_PASSWORD}
```

### Keyword Configuration

Edit `config/frequency_words.txt`:

```txt
AI
ChatGPT
OpenAI
+technology

Tesla
Musk
@10

Bitcoin
Cryptocurrency
!advertisement
```

## 📦 Docker Deployment

### Using Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/NebulaForgeX/NFX-News.git
cd NFX-News

# Copy and configure .env files
cp docker/.env.example docker/.env
# Edit docker/.env with your configuration

# Start all services
cd docker
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Individual Service Deployment

See individual service documentation for Docker deployment instructions:
- [Crawl Server Docker](docs/crawl_server/README.md#docker-deployment)
- [News Server Docker](docs/news_server/README.md#docker-deployment)
- [Web Server Docker](docs/web_server/README.md#docker-deployment)
- [MCP Server Docker](docs/mcp_server/README.md#docker-deployment)

## 🤖 AI Integration

### Supported Clients

- **Claude Desktop**: Native MCP support
- **Cursor**: HTTP or STDIO mode
- **VSCode (Cline/Continue)**: Full MCP integration
- **Any MCP-compatible client**: Standard MCP protocol

### Quick Setup (Claude Desktop)

1. Edit Claude Desktop config:
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
   - **Mac**: `~/Library/Application Support/Claude/claude_desktop_config.json`

2. Add configuration:
```json
{
  "mcpServers": {
    "nfx-news": {
      "command": "uv",
      "args": [
        "--directory",
        "/path/to/NFX-News",
        "run",
        "python",
        "-m",
        "mcp_server.server"
      ]
    }
  }
}
```

3. Restart Claude Desktop and start querying!

See [MCP Server Documentation](docs/mcp_server/README.md) for detailed setup instructions.

## 🔄 Data Flow

```
┌─────────────────┐
│  News Platforms │ (Zhihu, Weibo, etc.)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Crawl Server   │ ──► Crawls trending content
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │ ──► Stores news data
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌─────────┐ ┌──────────┐
│  News   │ │   Web    │
│ Server  │ │  Server  │
│  (API)  │ │ (Reports)│
└─────────┘ └──────────┘
    │
    ▼
┌──────────┐
│   MCP    │ ──► AI Analysis
│  Server  │
└──────────┘
```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify PostgreSQL is running
   - Check connection credentials
   - Ensure database exists

2. **Service Won't Start**
   - Check port availability
   - Verify dependencies are installed
   - Check logs for detailed errors

3. **No Data in Reports**
   - Ensure crawl server has run
   - Check `output/` directory exists
   - Verify crawler configuration

### Getting Help

- Check service-specific documentation in `docs/`
- Review logs: Each service outputs logs to stdout
- Check configuration files: Ensure all required settings are present

## 📊 Project Structure

```
NFX-News/
├── config/                 # Configuration files
│   ├── config.yaml        # Main configuration
│   └── frequency_words.txt # Keyword configuration
├── crawl_server/           # Crawler service (Python)
├── news_server/            # News API service (TypeScript)
├── web_server/             # Web report service (Python)
├── mcp_server/             # AI analysis service (Python)
├── output/                 # Crawled news data (generated)
├── docs/                   # Documentation
│   ├── crawl_server/
│   ├── news_server/
│   ├── web_server/
│   └── mcp_server/
├── docker/                 # Docker configuration
├── data/                   # Additional data files
└── README.md              # This file
```

## 🔗 Related Projects

- **NewsNow**: News aggregation service that powers platform data
- **FastMCP**: MCP protocol implementation framework

## 📄 License

GPL-3.0 License

See [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [FastAPI](https://fastapi.tiangolo.com/) and [Fastify](https://www.fastify.io/)
- Powered by [PostgreSQL](https://www.postgresql.org/) and [Redis](https://redis.io/)
- AI capabilities via [Model Context Protocol](https://modelcontextprotocol.io/)

---

<div align="center">

**[🔝 Back to Top](#nfx-news)**

Made with ❤️ for the open-source community

</div>
