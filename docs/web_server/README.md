# TrendRadar Web Server

**Version**: 1.0.0  
**Language**: Python 3.8+  
**Framework**: FastAPI  
**Architecture**: MVC (Model-View-Controller)

## 📋 Overview

Web Server provides HTML report viewing interface for TrendRadar, built with FastAPI and following MVC (Model-View-Controller) architecture principles. It serves as the web interface for browsing and viewing trending news reports.

### Core Features

- ✅ **MVC Architecture**: Clear separation of concerns
- ✅ **HTML Reports**: Beautiful, mobile-responsive HTML reports
- ✅ **RESTful API**: Standard RESTful API endpoints
- ✅ **Health Checks**: Health check endpoints for monitoring
- ✅ **CORS Support**: Cross-origin resource sharing enabled
- ✅ **Auto-reload**: Support for development mode with auto-reload

## 🏗️ Architecture

### Directory Structure

```
web_server/
├── main.py                  # Main entry point (FastAPI app)
├── config.py               # Configuration management
├── server.py               # Server initialization
├── controllers/            # Controller layer (MVC)
│   ├── root_controller.py  # Root routes
│   ├── health_controller.py # Health check routes
│   └── report_controller.py # Report routes
├── models/                 # Model layer (MVC)
├── views/                  # View layer (MVC) - Templates
├── utils/                  # Utility functions
└── version                 # Version file
```

### MVC Architecture

- **Model**: Data models and business logic (`models/`)
- **View**: HTML templates and rendering (`views/`)
- **Controller**: Request handling and routing (`controllers/`)

## 🚀 Quick Start

### 1. Install Dependencies

```bash
cd web_server
pip install -r requirements.txt
```

### 2. Configure Environment

Set environment variables or use defaults:

```bash
# Project root directory (where output/ is located)
export PROJECT_ROOT=/path/to/TrendRadar

# Server configuration
export HOST=0.0.0.0
export PORT=10199
export DEBUG=false
```

### 3. Start Server

#### Development Mode

```bash
# Run directly
python -m web_server.main

# Or using uvicorn
uvicorn web_server.main:app --host 0.0.0.0 --port 10199 --reload
```

#### Production Mode

```bash
# Using uvicorn with workers
uvicorn web_server.main:app --host 0.0.0.0 --port 10199 --workers 4
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PROJECT_ROOT` | Project root directory | `/app` | Yes |
| `HOST` | Server host | `0.0.0.0` | No |
| `PORT` | Server port | `10199` | No |
| `DEBUG` | Debug mode | `false` | No |

### Configuration File

Configuration is managed through environment variables and `config.py`:

```python
# config.py
@dataclass
class AppConfig:
    project_root: str = "/app"
    host: str = "0.0.0.0"
    port: int = 10199
    output_dir: str = "output"
    debug: bool = False
```

## 📚 API Documentation

### Root Endpoints

#### Health Check

```http
GET /health
```

**Response**:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "timestamp": "2025-01-15T10:00:00Z"
}
```

#### API Documentation

If debug mode is enabled, access:
- Swagger UI: `http://localhost:10199/docs`
- ReDoc: `http://localhost:10199/redoc`

### Report Endpoints

#### Get Latest Report

```http
GET /report
```

Returns the latest HTML report for browsing.

**Response**: HTML page with trending news report

#### Get Report by Date

```http
GET /report/{date}
```

**Parameters**:
- `date`: Date in YYYY-MM-DD format

**Example**:
```http
GET /report/2025-01-15
```

**Response**: HTML page with trending news report for the specified date

#### Get Report Data (JSON)

```http
GET /api/report
```

Returns JSON data of the latest report.

**Response**:
```json
{
  "success": true,
  "data": {
    "date": "2025-01-15",
    "total_news": 150,
    "platforms": [...],
    "trending_topics": [...]
  }
}
```

## 🔧 Development Guide

### Adding New Routes

1. Create a new controller file in `controllers/`
2. Define routes using FastAPI router
3. Register router in `main.py`

Example:

```python
# controllers/my_controller.py
from fastapi import APIRouter

router = APIRouter()

@router.get("/my-endpoint")
async def my_endpoint():
    return {"message": "Hello World"}
```

```python
# main.py
from .controllers import my_controller

app.include_router(my_controller.router)
```

### Adding New Views

1. Create HTML template in `views/` directory
2. Use Jinja2 for templating (if integrated)
3. Return rendered template from controller

### Using Models

Create data models in `models/` directory:

```python
# models/report.py
from pydantic import BaseModel
from datetime import datetime

class ReportModel(BaseModel):
    date: datetime
    total_news: int
    platforms: list
```

## 📊 Features

### 1. HTML Report Generation

Generates beautiful, mobile-responsive HTML reports from crawled news data.

### 2. Date-based Reports

Access historical reports by date for trend analysis.

### 3. JSON API

Provides JSON API for programmatic access to report data.

### 4. Health Monitoring

Health check endpoints for service monitoring and load balancing.

## 🐛 Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using the port
   lsof -i :10199  # Mac/Linux
   netstat -ano | findstr :10199  # Windows
   
   # Use a different port
   export PORT=10200
   ```

2. **Output Directory Not Found**
   - Ensure `output/` directory exists in project root
   - Check `PROJECT_ROOT` environment variable is set correctly

3. **Report Not Found**
   - Verify crawler has run and generated data in `output/` directory
   - Check date format (YYYY-MM-DD)
   - Ensure data files exist for requested date

### Debug Mode

Enable debug mode for detailed error messages:

```bash
export DEBUG=true
python -m web_server.main
```

Debug mode also enables:
- Swagger UI at `/docs`
- ReDoc at `/redoc`
- Detailed error tracebacks

## 📦 Docker Deployment

### Dockerfile Example

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY web_server/ ./web_server/
COPY config/ ./config/
COPY output/ ./output/
COPY requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt

ENV PROJECT_ROOT=/app
ENV PORT=10199

EXPOSE 10199

CMD ["uvicorn", "web_server.main:app", "--host", "0.0.0.0", "--port", "10199"]
```

### docker-compose Example

```yaml
version: '3.8'

services:
  web-server:
    build: ./web_server
    ports:
      - "10199:10199"
    environment:
      - PROJECT_ROOT=/app
      - PORT=10199
      - DEBUG=false
    volumes:
      - ./output:/app/output:ro
      - ./config:/app/config:ro
    depends_on:
      - postgres
```

## 🔗 Related Services

- **crawl_server**: Crawler service that generates data for reports
- **news_server**: News aggregation API service
- **mcp_server**: AI analysis service

## 📄 License

GPL-3.0 License
