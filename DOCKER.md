# Docker Setup Guide

This guide explains how to build and run the Recreation MCP Server using Docker.

## Prerequisites

- Docker installed ([Install Docker](https://docs.docker.com/get-docker/))
- Docker Compose installed (included with Docker Desktop)
- API keys for:
  - [National Park Service](https://www.nps.gov/subjects/developer/get-started.htm)
  - [Recreation.gov](https://ridb.recreation.gov/)
  - [OpenWeatherMap](https://openweathermap.org/api)

## Quick Start

### 1. Set up environment variables

Copy the example file and add your API keys:

```bash
cp .env.example .env
# Edit .env and add your API keys
```

Your `.env` file should look like:

```bash
NPS_API_KEY=your_nps_api_key_here
RECREATION_GOV_API_KEY=your_recreation_gov_api_key_here
OPENWEATHER_API_KEY=your_openweather_api_key_here
```

### 2. Build the Docker image

```bash
docker-compose build
```

Or manually:

```bash
docker build -t recreation-mcp-server .
```

### 3. Run the server

```bash
docker-compose up
```

The server will start and wait for MCP protocol commands on stdin/stdout.

## Configuration

### Environment Variables

The following environment variables can be set in your `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `NPS_API_KEY` | *(required)* | National Park Service API key |
| `RECREATION_GOV_API_KEY` | *(required)* | Recreation.gov API key |
| `OPENWEATHER_API_KEY` | *(required)* | OpenWeatherMap API key |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `CACHE_ENABLED` | `true` | Enable response caching |
| `CACHE_TTL_SECONDS` | `3600` | Cache time-to-live in seconds (1 hour) |

### Optional Configuration File

You can use a `config.yaml` file for advanced configuration. Uncomment the volume mount in `docker-compose.yml`:

```yaml
volumes:
  - ./config.yaml:/home/mcp/config.yaml:ro  # Uncomment this line
  - cache-data:/home/mcp/.cache
```

See the main README for configuration file format.

## Docker Image Details

### Multi-Stage Build

The Dockerfile uses a two-stage build process:

1. **Builder stage**: Compiles the Go binary
  - Based on `golang:1.25.9-alpine3.22`
   - Downloads dependencies
   - Builds static binary with security flags
   
2. **Runtime stage**: Creates minimal runtime image
  - Based on `alpine:3.22.2`
   - Only includes necessary runtime dependencies
   - Creates non-root user (`mcp:mcp` with UID/GID 1000)
   - Binary size: ~20MB (optimized with `-ldflags='-w -s'`)

### Security Features

The container runs with several security hardening options:

- **Non-root user**: Runs as user `mcp` (UID 1000)
- **Read-only filesystem**: Container filesystem is read-only
- **No new privileges**: Prevents privilege escalation
- **Minimal dependencies**: Only includes `ca-certificates` and `tzdata`
- **Static binary**: No dynamic linking reduces attack surface

### Image Size

- Builder image: ~500MB (discarded after build)
- Final runtime image: ~30MB
- Binary size: ~20MB

## Usage with Claude Desktop

Add to your Claude Desktop configuration file:

**Location:**
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

**Configuration:**

```json
{
  "mcpServers": {
    "recreation": {
      "command": "docker",
      "args": [
        "compose",
        "-f",
        "/absolute/path/to/recreation-mcp-server/docker-compose.yml",
        "run",
        "--rm",
        "mcp-recreation-server"
      ]
    }
  }
}
```

**Important:** Replace `/absolute/path/to/recreation-mcp-server/` with the actual absolute path to your project directory.

## Common Commands

### Build the image

```bash
docker-compose build
```

### Run interactively

```bash
docker-compose up
```

### Run in detached mode

```bash
docker-compose up -d
```

### View logs

```bash
docker-compose logs -f
```

### Stop the server

```bash
docker-compose down
```

### Remove volumes (cache data)

```bash
docker-compose down -v
```

### Rebuild from scratch

```bash
docker-compose build --no-cache
```

### Run without docker-compose

```bash
docker run --rm -i \
  --env-file .env \
  -v $(pwd)/.cache:/home/mcp/.cache \
  recreation-mcp-server
```

## Troubleshooting

### Build fails with "permission denied"

Ensure your user has permission to run Docker commands:

```bash
# macOS/Linux
sudo usermod -aG docker $USER
# Log out and back in for changes to take effect
```

### Container exits immediately

Check that your API keys are set correctly in `.env`:

```bash
docker-compose config | grep API_KEY
```

### Cannot connect from Claude Desktop

1. Verify the container is running: `docker-compose ps`
2. Check logs for errors: `docker-compose logs`
3. Ensure the path in Claude config is absolute
4. Restart Claude Desktop after config changes

### "config.yaml not found" error

The `config.yaml` volume is optional. Comment it out in `docker-compose.yml` if you're not using it:

```yaml
volumes:
  # - ./config.yaml:/home/mcp/config.yaml:ro  # Commented out
  - cache-data:/home/mcp/.cache
```

### Container runs but no response

The MCP server communicates via stdin/stdout. If running interactively, it should respond to MCP protocol messages. For normal use, let Claude Desktop manage the container lifecycle.

## Development

### Testing the Docker build

```bash
# Build
docker build -t recreation-mcp-server:test .

# Run with test environment
docker run --rm -it \
  -e NPS_API_KEY=test \
  -e RECREATION_GOV_API_KEY=test \
  -e OPENWEATHER_API_KEY=test \
  -e LOG_LEVEL=debug \
  recreation-mcp-server:test
```

### Inspecting the image

```bash
# Check image size
docker images recreation-mcp-server

# Inspect image layers
docker history recreation-mcp-server:test

# Run a shell in the container
docker run --rm -it --entrypoint /bin/sh recreation-mcp-server:test
```

### Build for multiple architectures

```bash
# Enable BuildKit
export DOCKER_BUILDKIT=1

# Build for multiple platforms
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t recreation-mcp-server:latest .
```

## Best Practices

1. **Always use `.env` file**: Don't hardcode API keys in docker-compose.yml
2. **Pin image versions**: Update base image versions periodically
3. **Review security**: Run `docker scan recreation-mcp-server` to check for vulnerabilities
4. **Clean up**: Remove unused images with `docker image prune`
5. **Monitor resources**: Check resource usage with `docker stats`

## Production Considerations

For production deployments, consider:

- Using a container orchestration platform (Kubernetes, ECS, etc.)
- Implementing health checks
- Setting up log aggregation
- Configuring resource limits (CPU, memory)
- Using secrets management (not `.env` files)
- Setting up monitoring and alerting
- Implementing automated backups of cache data

---

**Note:** This is a development/demo setup. For production use, additional security hardening and operational practices should be implemented.
