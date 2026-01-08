# TrendRadar MCP Server

**Version**: 1.0.1  
**Language**: Python 3.8+  
**Framework**: FastMCP 2.0

## 📋 Overview

MCP Server provides AI-powered analysis capabilities for TrendRadar using the Model Context Protocol (MCP). It enables natural language interaction with news data for deep analysis and insights.

### Core Features

- ✅ **MCP Protocol Support**: Standard Model Context Protocol implementation
- ✅ **Multiple Analysis Tools**: 13+ AI analysis tools
- ✅ **Dual Transport Modes**: STDIO and HTTP support
- ✅ **Natural Language Queries**: Query news data using natural language
- ✅ **Trend Analysis**: Track topic popularity and trends over time
- ✅ **Data Insights**: Cross-platform data comparison and analysis

## 🏗️ Architecture

### Directory Structure

```
mcp_server/
├── server.py              # Main MCP server (FastMCP 2.0)
├── tools/                 # Analysis tools
│   ├── data_query.py     # Data query tools
│   ├── analytics.py      # Analytics tools
│   ├── search_tools.py   # Search tools
│   ├── config_mgmt.py    # Configuration management
│   └── system.py         # System management
├── services/              # Business logic services
├── utils/                 # Utility functions
└── routes/                # HTTP routes (if HTTP mode)
```

### Tool Categories

#### 1. Data Query Tools
- `get_latest_news`: Get latest crawled news
- `get_news_by_date`: Query news by date
- `get_trending_topics`: Get trending topics statistics

#### 2. Analytics Tools
- `analyze_topic_trend`: Analyze topic popularity trends
- `analyze_data_insights`: Cross-platform data insights
- `analyze_sentiment`: Sentiment analysis
- `generate_summary_report`: Generate summary reports

#### 3. Search Tools
- `search_news`: Search news by keywords
- `search_related_news_history`: Search related historical news
- `find_similar_news`: Find similar news articles

#### 4. Configuration Management
- `get_current_config`: Get current configuration
- `get_system_status`: Get system status

#### 5. System Management
- `trigger_crawl`: Trigger crawl task manually

## 🚀 Quick Start

### 1. Install Dependencies

```bash
cd mcp_server
pip install -r requirements.txt
```

### 2. Configure Environment

Create `.env` file or set environment variables:

```bash
# Project root directory (where config/ and output/ are located)
TRENDRADAR_PROJECT_ROOT=/path/to/TrendRadar

# HTTP mode configuration (optional)
MCP_HTTP_PORT=3333
MCP_HTTP_HOST=0.0.0.0
```

### 3. Start Server

#### STDIO Mode (Recommended)

```bash
# Direct execution
python -m mcp_server.server

# Or using uv (recommended)
uv run python -m mcp_server.server
```

#### HTTP Mode

```bash
# Start HTTP server
python -m mcp_server.server --transport http --port 3333

# Or using provided scripts
# Windows
start-http.bat

# Mac/Linux
./start-http.sh
```

## ⚙️ Configuration

### Project Structure Requirements

The MCP server expects the following directory structure:

```
TrendRadar/
├── config/
│   ├── config.yaml           # Main configuration
│   └── frequency_words.txt   # Keyword configuration
├── output/                    # Crawled news data
│   ├── 2025-01-15/
│   │   ├── zhihu.json
│   │   ├── weibo.json
│   │   └── ...
│   └── ...
└── mcp_server/               # MCP server code
```

### Configuration File

The server reads configuration from `config/config.yaml`:

```yaml
platforms:
  - id: "zhihu"
    name: "Zhihu"
  - id: "weibo"
    name: "Weibo"
  # ... more platforms

report:
  mode: "current"  # daily, current, incremental
```

## 📚 Tool Usage Examples

### Get Latest News

```python
# Tool: get_latest_news
# Parameters:
#   platforms: Optional[List[str]] - e.g., ["zhihu", "weibo"]
#   limit: int - Default 50, max 1000
#   include_url: bool - Default False

# Example query:
# "Get the latest 100 news articles from Zhihu and Weibo"
```

### Search News

```python
# Tool: search_news
# Parameters:
#   keywords: List[str] - e.g., ["AI", "ChatGPT"]
#   platforms: Optional[List[str]]
#   date_range: Optional[Dict[str, str]]
#   limit: int - Default 50

# Example query:
# "Search for news about AI and ChatGPT in the last 7 days"
```

### Analyze Topic Trend

```python
# Tool: analyze_topic_trend
# Parameters:
#   topic: str - Topic keyword
#   days: int - Number of days to analyze
#   platforms: Optional[List[str]]

# Example query:
# "Analyze the popularity trend of 'Bitcoin' over the past 30 days"
```

### Generate Summary Report

```python
# Tool: generate_summary_report
# Parameters:
#   date: str - Date in YYYY-MM-DD format
#   platforms: Optional[List[str]]

# Example query:
# "Generate a summary report for 2025-01-15"
```

## 🔌 Client Integration

### Claude Desktop

Edit Claude Desktop's MCP config file:

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`  
**Mac**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "trendradar": {
      "command": "uv",
      "args": [
        "--directory",
        "/path/to/TrendRadar",
        "run",
        "python",
        "-m",
        "mcp_server.server"
      ],
      "env": {},
      "disabled": false,
      "alwaysAllow": []
    }
  }
}
```

### Cursor

Create `.cursor/mcp.json` in project root:

```json
{
  "mcpServers": {
    "trendradar": {
      "command": "uv",
      "args": [
        "--directory",
        "/path/to/TrendRadar",
        "run",
        "python",
        "-m",
        "mcp_server.server"
      ]
    }
  }
}
```

### HTTP Mode (Universal)

Any MCP-compatible client can connect via HTTP:

```json
{
  "name": "trendradar",
  "url": "http://localhost:3333/mcp",
  "type": "http",
  "description": "News Trending Aggregation Analysis"
}
```

Start HTTP server first:

```bash
python -m mcp_server.server --transport http --port 3333
```

## 🔧 Development Guide

### Adding a New Tool

1. Create tool function in appropriate tool file (e.g., `tools/data_query.py`)
2. Register with `@mcp.tool` decorator
3. Add tool to `_get_tools()` initialization

Example:

```python
from fastmcp import FastMCP

mcp = FastMCP('trendradar-news')

@mcp.tool
async def my_new_tool(param1: str, param2: int = 10) -> str:
    """
    Tool description for AI to understand usage.
    
    Args:
        param1: Description of param1
        param2: Description of param2 (default: 10)
    
    Returns:
        JSON string with results
    """
    # Tool logic here
    result = {"data": "example"}
    return json.dumps(result, ensure_ascii=False)
```

### Testing Tools

Use MCP Inspector for testing:

```bash
# Start HTTP server
python -m mcp_server.server --transport http --port 3333

# Start MCP Inspector
npx @modelcontextprotocol/inspector

# Connect to http://localhost:3333/mcp
```

## 🐛 Troubleshooting

### Common Issues

1. **Data Not Found**
   - Ensure crawler has run and generated data in `output/` directory
   - Check date range for available data
   - Verify `config/config.yaml` exists

2. **Import Errors**
   - Ensure all dependencies are installed: `pip install -r requirements.txt`
   - Check Python version (requires 3.8+)

3. **HTTP Server Won't Start**
   - Check if port 3333 is already in use
   - Verify firewall settings
   - Try a different port: `--port 33333`

4. **Client Connection Failed**
   - For STDIO mode: Verify command path and project root path
   - For HTTP mode: Ensure server is running and accessible
   - Check client logs for detailed error messages

### Debug Mode

Enable verbose logging:

```bash
# Set environment variable
export MCP_DEBUG=true

# Run server
python -m mcp_server.server
```

## 📦 Docker Deployment

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY mcp_server/ ./mcp_server/
COPY config/ ./config/
COPY requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt

# STDIO mode (default)
CMD ["python", "-m", "mcp_server.server"]

# HTTP mode
# CMD ["python", "-m", "mcp_server.server", "--transport", "http", "--port", "3333"]
```

## 🔗 Related Services

- **crawl_server**: Crawler service that generates data for analysis
- **news_server**: News aggregation API service
- **web_server**: Web report service

## 📄 License

GPL-3.0 License
