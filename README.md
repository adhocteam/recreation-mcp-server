[![CI](https://github.com/adhocteam/recreation-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/adhocteam/recreation-mcp-server/actions/workflows/ci.yml)

# Recreation Opportunities MCP Server

An MCP (Model Context Protocol) server written in Go that enables LLMs like Claude to discover and learn about recreation opportunities by integrating three public REST APIs:
- National Park Service API - Information about national parks, activities, alerts, and campgrounds
- Recreation.gov API - Recreation areas, facilities, and campsites across federal lands
- OpenWeatherMap API - Current weather and forecasts for outdoor locations

## Features

- Search and explore national parks across the United States
- Find campgrounds and recreation facilities
- Get weather information for planning outdoor activities
- Query parks by state, activity, or keyword
- Check park alerts and closures
- Discover recreation areas and detailed facility information
- Easy deployment via Docker and docker-compose
- Seamless integration with Claude Desktop

## MCP Tools

This server provides 9 tools for exploring recreation opportunities:

1. **search_parks** - Search for national parks by name, state, or activity
2. **get_park_details** - Get detailed information about a specific park
3. **get_park_alerts** - Get current alerts and closures for a park
4. **search_campgrounds** - Search for campgrounds in national parks or recreation areas
5. **search_recreation_areas** - Search recreation areas on Recreation.gov
6. **get_facility_details** - Get detailed information about a specific facility
7. **get_weather** - Get current weather conditions for a location
8. **get_weather_forecast** - Get 5-day weather forecast for a location
9. **list_activities** - List all available activities across NPS and Recreation.gov

## Prerequisites

- Go 1.21+ for development and building
- Docker for containerized deployment
- Docker Compose for orchestration
- API Keys - Free API keys from:
  - [National Park Service](https://www.nps.gov/subjects/developer/get-started.htm)
  - [Recreation.gov (RIDB)](https://ridb.recreation.gov/)
  - [OpenWeatherMap](https://openweathermap.org/api)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/markheadd/recreation-mcp-server.git
cd recreation-mcp-server
```

### 2. Set Up API Keys

Copy the example environment file and add your API keys:

```bash
cp .env.example .env
# Edit .env and add your API keys
```

### 3. Build and Run with Docker

```bash
docker-compose build
docker-compose up
```

### 4. Configure Claude Desktop

Add the following configuration to your Claude Desktop config file (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "recreation": {
      "command": "docker",
      "args": [
        "compose",
        "-f",
        "/path/to/recreation-mcp-server/docker-compose.yml",
        "run",
        "--rm",
        "mcp-recreation-server"
      ]
    }
  }
}
```

Replace `/path/to/recreation-mcp-server` with the actual path to this repository.

## Development

### Build Locally

```bash
go build -o mcp-server ./cmd/server
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Locally (without Docker)

```bash
# Make sure .env file exists with your API keys
go run ./cmd/server
```

## Example Queries

Once integrated with Claude Desktop, you can ask questions like:

- "What national parks are in Colorado?"
- "Tell me about Yosemite National Park"
- "Are there any alerts for Grand Canyon?"
- "Find campgrounds near Yellowstone"
- "What's the weather like at Rocky Mountain National Park?"
- "Show me recreation areas in Utah with hiking"

## Project Structure

```
recreation-mcp-server/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── mcp/            # MCP server implementation
│   ├── api/            # API clients (NPS, Recreation.gov, Weather)
│   ├── models/         # Data models
│   ├── cache/          # Response caching
│   └── config/         # Configuration management
├── pkg/
│   └── util/           # Utility functions
├── test/
│   ├── fixtures/       # Test data
│   └── integration/    # Integration tests
├── spec/               # Project specification
├── Dockerfile          # Container definition
├── docker-compose.yml  # Container orchestration
└── README.md          # This file
```

## Configuration

The server can be configured through environment variables or an optional config.yaml file. See .env.example for available options.

### Environment Variables

- NPS_API_KEY - National Park Service API key (required)
- RECREATION_GOV_API_KEY - Recreation.gov API key (required)
- OPENWEATHER_API_KEY - OpenWeatherMap API key (required)
- LOG_LEVEL - Logging level: debug, info, warn, error (default: info)
- CACHE_ENABLED - Enable response caching (default: true)
- CACHE_TTL_SECONDS - Cache time-to-live in seconds (default: 3600)
- MAX_REQUESTS_PER_MINUTE - Rate limiting (default: 60)

## Troubleshooting

### Docker Container Won't Start

- Verify API keys are set in `.env` file
- Check Docker logs: `docker-compose logs`
- Ensure no other services are using the same ports

### Claude Desktop Can't Connect

- Verify the path in `claude_desktop_config.json` is correct
- Restart Claude Desktop after configuration changes
- Check that docker-compose runs successfully from the command line

### API Errors

- Verify your API keys are valid and active
- Check API rate limits haven't been exceeded
- Ensure network connectivity to external APIs

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Inspired by [nps-explorer-mcp-server](https://github.com/Kyle-Ski/nps-explorer-mcp-server)
- Built with the Model Context Protocol (MCP)
- Powered by Go and Docker
