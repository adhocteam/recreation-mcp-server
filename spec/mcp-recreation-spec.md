# Recreation Opportunities MCP Server - Software Specification

## 1. Project Overview

### 1.1 Purpose
Build an MCP (Model Context Protocol) server in Go that enables LLMs like Claude to discover and learn about recreation opportunities by integrating three public REST APIs:
- National Park Service API
- Recreation.gov API  
- OpenWeatherMap API

### 1.2 Goals
- Provide a unified interface for querying recreation data from multiple sources
- Enable natural language exploration of parks, campgrounds, and outdoor activities
- Include weather information to help users plan visits
- Run locally via Docker for easy demonstration and portability
- Integrate seamlessly with Claude Desktop

### 1.3 Reference Implementation
A similar TypeScript implementation exists at: https://github.com/Kyle-Ski/nps-explorer-mcp-server
Key differences for this project:
- Language: Go (not TypeScript)
- Deployment: Docker + docker-compose (not CloudFlare Workers)
- Authentication: None for initial version (not GitHub OAuth)

## 2. Technical Architecture

### 2.1 Technology Stack
- **Language**: Go 1.21+
- **Protocol**: MCP (Model Context Protocol)
- **Containerization**: Docker
- **Orchestration**: docker-compose
- **MCP SDK**: Use official Go MCP SDK if available, or implement MCP protocol directly

### 2.2 System Components
```
┌─────────────────┐
│  Claude Desktop │
│   (MCP Client)  │
└────────┬────────┘
         │ MCP Protocol (stdio)
         ▼
┌─────────────────────────────────┐
│   MCP Server (Go)               │
│   ┌─────────────────────────┐   │
│   │  MCP Protocol Handler   │   │
│   └──────────┬──────────────┘   │
│              │                   │
│   ┌──────────▼──────────────┐   │
│   │   Tool Implementations  │   │
│   │  - NPS API Client       │   │
│   │  - Recreation.gov Client│   │
│   │  - OpenWeather Client   │   │
│   └─────────────────────────┘   │
└─────────────────────────────────┘
         │ HTTPS
         ▼
┌─────────────────────────────────┐
│   External REST APIs            │
│   - api.nps.gov                 │
│   - ridb.recreation.gov         │
│   - api.openweathermap.org      │
└─────────────────────────────────┘
```

### 2.3 Deployment Architecture
- Single Docker container running the MCP server
- Configuration via environment variables
- Volume mounts for persistent configuration (optional)
- Communicates with Claude Desktop via stdio

## 3. API Integration Requirements

### 3.1 National Park Service API
**Base URL**: `https://developer.nps.gov/api/v1/`

**Authentication**: API key via query parameter `?api_key=YOUR_KEY`

**Key Endpoints to Implement**:
- `GET /parks` - Search and list parks
- `GET /parks/{parkCode}` - Get detailed park information
- `GET /activities` - List available activities
- `GET /activities/parks` - Find parks by activity
- `GET /alerts` - Get park alerts and closures
- `GET /campgrounds` - Search campgrounds in parks
- `GET /visitorcenters` - List visitor centers

**Required Parameters**:
- `stateCode` - Filter by state (e.g., CA, CO)
- `limit` - Results per page
- `start` - Pagination offset
- `q` - Search query

### 3.2 Recreation.gov API (RIDB)
**Base URL**: `https://ridb.recreation.gov/api/v1/`

**Authentication**: API key via header `apikey: YOUR_KEY`

**Key Endpoints to Implement**:
- `GET /recareas` - Search recreation areas
- `GET /recareas/{recAreaId}` - Get recreation area details
- `GET /facilities` - Search facilities (campgrounds, day use areas)
- `GET /facilities/{facilityId}` - Get facility details
- `GET /campsites` - Search individual campsites
- `GET /activities` - List activities

**Required Parameters**:
- `query` - Search text
- `state` - State abbreviation
- `activity` - Activity ID filter
- `limit` - Results per page
- `offset` - Pagination offset

### 3.3 OpenWeatherMap API
**Base URL**: `https://api.openweathermap.org/data/2.5/`

**Authentication**: API key via query parameter `?appid=YOUR_KEY`

**Key Endpoints to Implement**:
- `GET /weather` - Current weather by coordinates
- `GET /forecast` - 5-day forecast by coordinates

**Required Parameters**:
- `lat` - Latitude
- `lon` - Longitude
- `units` - metric or imperial

## 4. MCP Tools Specification

### 4.1 Tool: search_parks
**Description**: Search for national parks by name, state, or activity

**Input Schema**:
```json
{
  "query": "string (optional) - Search text",
  "state": "string (optional) - Two-letter state code",
  "activity": "string (optional) - Activity name",
  "limit": "integer (optional, default: 10) - Max results"
}
```

**Output**: List of parks with name, description, state, park code, and URL

**Implementation**: Queries NPS API `/parks` endpoint

### 4.2 Tool: get_park_details
**Description**: Get detailed information about a specific park

**Input Schema**:
```json
{
  "park_code": "string (required) - Four-letter park code (e.g., 'yose' for Yosemite)"
}
```

**Output**: Detailed park information including hours, fees, activities, contacts, and directions

**Implementation**: Queries NPS API `/parks/{parkCode}` endpoint

### 4.3 Tool: get_park_alerts
**Description**: Get current alerts and closures for a park

**Input Schema**:
```json
{
  "park_code": "string (required) - Four-letter park code"
}
```

**Output**: List of alerts with title, description, category, and dates

**Implementation**: Queries NPS API `/alerts` endpoint

### 4.4 Tool: search_campgrounds
**Description**: Search for campgrounds in national parks or recreation areas

**Input Schema**:
```json
{
  "park_code": "string (optional) - NPS park code",
  "state": "string (optional) - State code",
  "query": "string (optional) - Search text",
  "limit": "integer (optional, default: 10)"
}
```

**Output**: List of campgrounds with name, description, amenities, and reservation info

**Implementation**: Queries both NPS API `/campgrounds` and Recreation.gov API `/facilities`

### 4.5 Tool: search_recreation_areas
**Description**: Search recreation areas on Recreation.gov

**Input Schema**:
```json
{
  "query": "string (optional) - Search text",
  "state": "string (optional) - State code",
  "activity": "string (optional) - Activity name",
  "limit": "integer (optional, default: 10)"
}
```

**Output**: List of recreation areas with name, description, activities, and facilities

**Implementation**: Queries Recreation.gov API `/recareas` endpoint

### 4.6 Tool: get_facility_details
**Description**: Get detailed information about a specific facility

**Input Schema**:
```json
{
  "facility_id": "string (required) - Recreation.gov facility ID"
}
```

**Output**: Detailed facility information including amenities, contact, and reservation details

**Implementation**: Queries Recreation.gov API `/facilities/{facilityId}` endpoint

### 4.7 Tool: get_weather
**Description**: Get current weather conditions for a location

**Input Schema**:
```json
{
  "latitude": "number (required)",
  "longitude": "number (required)",
  "units": "string (optional, default: 'imperial') - 'metric' or 'imperial'"
}
```

**Output**: Current weather with temperature, conditions, humidity, wind

**Implementation**: Queries OpenWeatherMap API `/weather` endpoint

### 4.8 Tool: get_weather_forecast
**Description**: Get 5-day weather forecast for a location

**Input Schema**:
```json
{
  "latitude": "number (required)",
  "longitude": "number (required)",
  "units": "string (optional, default: 'imperial')"
}
```

**Output**: Weather forecast with daily high/low, conditions, precipitation

**Implementation**: Queries OpenWeatherMap API `/forecast` endpoint

### 4.9 Tool: list_activities
**Description**: List all available activities across NPS and Recreation.gov

**Input Schema**:
```json
{
  "source": "string (optional, default: 'all') - 'nps', 'recreation_gov', or 'all'",
  "limit": "integer (optional, default: 50) - Max results"
}
```

**Output**: Comprehensive list of activity names and IDs from specified source(s)

**Implementation**: Queries both NPS API `/activities` and Recreation.gov API `/activities`, or just one based on source parameter

## 5. Configuration Management

### 5.1 Environment Variables
```bash
# API Keys (required)
NPS_API_KEY=your_nps_api_key_here
RECREATION_GOV_API_KEY=your_recreation_gov_api_key_here
OPENWEATHER_API_KEY=your_openweather_api_key_here

# Server Configuration
LOG_LEVEL=info  # debug, info, warn, error
CACHE_ENABLED=true
CACHE_TTL_SECONDS=3600

# API Rate Limiting (optional)
MAX_REQUESTS_PER_MINUTE=60
```

### 5.2 Configuration File
Optional `config.yaml` support for advanced configuration:
```yaml
apis:
  nps:
    base_url: "https://developer.nps.gov/api/v1"
    timeout: 30s
    retry_attempts: 3
  recreation_gov:
    base_url: "https://ridb.recreation.gov/api/v1"
    timeout: 30s
    retry_attempts: 3
  openweather:
    base_url: "https://api.openweathermap.org/data/2.5"
    timeout: 15s
    retry_attempts: 2

cache:
  enabled: true
  ttl: 1h
  max_size: 100MB

logging:
  level: info
  format: json
```

## 6. Docker Implementation

### 6.1 Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o mcp-server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 mcp && \
    adduser -D -u 1000 -G mcp mcp

# Set up working directory
WORKDIR /home/mcp

# Copy binary from builder with correct ownership
COPY --from=builder --chown=mcp:mcp /app/mcp-server .

# Switch to non-root user
USER mcp

# Health check (optional, can be removed if not needed)
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD ["./mcp-server", "--health"] || exit 1

# Run the MCP server
CMD ["./mcp-server"]
```

### 6.2 docker-compose.yml
```yaml
version: '3.8'

services:
  mcp-recreation-server:
    build: .
    container_name: recreation-mcp-server
    environment:
      - NPS_API_KEY=${NPS_API_KEY}
      - RECREATION_GOV_API_KEY=${RECREATION_GOV_API_KEY}
      - OPENWEATHER_API_KEY=${OPENWEATHER_API_KEY}
      - LOG_LEVEL=info
      - CACHE_ENABLED=true
      - CACHE_TTL_SECONDS=3600
    volumes:
      - ./config.yaml:/home/mcp/config.yaml:ro
      - cache-data:/home/mcp/.cache
    stdin_open: true
    tty: true
    # Security options
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
    user: "1000:1000"

volumes:
  cache-data:
    driver: local
```

### 6.3 .env.example
```bash
# Copy this file to .env and fill in your API keys

# National Park Service API Key
# Get yours at: https://www.nps.gov/subjects/developer/get-started.htm
NPS_API_KEY=

# Recreation.gov API Key  
# Get yours at: https://ridb.recreation.gov/
RECREATION_GOV_API_KEY=

# OpenWeatherMap API Key
# Get yours at: https://openweathermap.org/api
OPENWEATHER_API_KEY=
```

## 7. Project Structure

```
recreation-mcp-server/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── mcp/
│   │   ├── server.go         # MCP server implementation
│   │   ├── tools.go          # Tool definitions
│   │   └── handlers.go       # Tool handlers
│   ├── api/
│   │   ├── nps.go           # NPS API client
│   │   ├── recreation.go    # Recreation.gov client
│   │   └── weather.go       # OpenWeather client
│   ├── models/
│   │   └── types.go         # Shared data structures
│   ├── cache/
│   │   └── cache.go         # Response caching
│   └── config/
│       └── config.go        # Configuration management
├── pkg/
│   └── util/
│       └── helpers.go       # Utility functions
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── .dockerignore
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

## 8. Error Handling

### 8.1 API Error Responses
- Handle rate limiting (429 responses) with exponential backoff
- Handle API key validation errors clearly
- Handle network timeouts gracefully
- Log all API errors with context

### 8.2 MCP Error Responses
Follow MCP protocol error format:
```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable error message"
  }
}
```

Common error codes:
- `invalid_params` - Missing or invalid input parameters
- `api_error` - External API returned error
- `rate_limit` - Rate limit exceeded
- `not_found` - Resource not found
- `internal_error` - Server internal error

## 9. Performance Requirements

### 9.1 Response Times
- MCP tool calls should respond within 5 seconds under normal conditions
- Implement request timeouts (30s for NPS/Recreation.gov, 15s for weather)
- Cache frequently accessed data to reduce API calls

### 9.2 Caching Strategy
- Cache GET requests for 1 hour by default
- Use composite cache keys: `{endpoint}:{params_hash}`
- Implement LRU eviction when cache reaches size limit
- Allow cache bypass via configuration

## 10. Testing Requirements

### 10.1 Testing Philosophy
**Clean, Simple, Fast** - Tests must be:
- **Clean**: Easy to read and understand
- **Simple**: Minimal setup/teardown complexity
- **Fast**: Complete test suite runs in under 30 seconds

Tests that take too long or require elaborate setup will not be run frequently enough to be valuable. The testing strategy prioritizes speed and simplicity over exhaustive coverage.

### 10.2 Unit Tests (Primary Focus)
**Target**: 70%+ code coverage, all tests run in < 10 seconds

**Approach**:
- Mock all external API calls using interfaces
- No network calls in unit tests
- No database or file system dependencies
- Use table-driven tests for multiple scenarios
- Test files should be co-located with source files

**What to test**:
- API client request construction (without actual HTTP calls)
- Response parsing with fixture data
- Tool input validation
- Error handling logic
- Cache key generation and lookup logic
- Configuration parsing

**Example test structure**:
```go
func TestSearchParks(t *testing.T) {
    tests := []struct {
        name     string
        input    SearchParksInput
        mockResp string
        want     []Park
        wantErr  bool
    }{
        // test cases here
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Simple, fast test with mocked client
        })
    }
}
```

### 10.3 Integration Tests (Optional, Manual)
**Target**: Verify live API connectivity when needed

**Approach**:
- Separate package or build tag (e.g., `// +build integration`)
- Run manually with `go test -tags=integration`
- Not part of standard CI/CD pipeline
- Only run when debugging API issues

**What to test**:
- Live API connectivity with real credentials
- API contract compliance (response structure hasn't changed)
- Rate limiting behavior

**Keep it minimal**: Only 3-5 critical integration tests

### 10.4 End-to-End Testing (Manual)
**Target**: Validate MCP server works with Claude Desktop

**Approach**:
- Manual testing only
- Document test scenarios in README
- Create example conversations for validation

**Test scenarios**:
1. Search for parks in a specific state
2. Get detailed park information
3. Check weather for park coordinates
4. Search campgrounds with filters
5. Verify error messages are helpful

### 10.5 Test Organization
```
recreation-mcp-server/
├── internal/
│   ├── api/
│   │   ├── nps.go
│   │   ├── nps_test.go          # Fast unit tests with mocks
│   │   ├── recreation.go
│   │   └── recreation_test.go
│   └── mcp/
│       ├── handlers.go
│       └── handlers_test.go
├── test/
│   ├── fixtures/                # JSON fixture files for mocking
│   │   ├── nps_parks.json
│   │   ├── recreation_areas.json
│   │   └── weather.json
│   └── integration/             # Optional integration tests
│       └── api_integration_test.go  # +build integration
└── Makefile
```

### 10.6 Test Commands
```makefile
# Makefile targets for easy test execution

.PHONY: test
test:
	@echo "Running fast unit tests..."
	go test -v -race -timeout=30s ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration
test-integration:
	@echo "Running integration tests (requires API keys)..."
	go test -v -tags=integration -timeout=2m ./test/integration/...

.PHONY: test-quick
test-quick:
	@echo "Running quick smoke test..."
	go test -short ./...
```

### 10.7 Mocking Strategy
Use Go interfaces to enable simple mocking without frameworks:

```go
// API client interface for easy mocking
type NPSClient interface {
    SearchParks(ctx context.Context, params SearchParams) ([]Park, error)
    GetParkDetails(ctx context.Context, parkCode string) (*Park, error)
}

// Mock for testing
type MockNPSClient struct {
    SearchParksFunc func(ctx context.Context, params SearchParams) ([]Park, error)
}

func (m *MockNPSClient) SearchParks(ctx context.Context, params SearchParams) ([]Park, error) {
    if m.SearchParksFunc != nil {
        return m.SearchParksFunc(ctx, params)
    }
    return nil, nil
}
```

### 10.8 CI/CD Consideration
- Only unit tests run in CI/CD pipeline
- Fast feedback loop (< 30 seconds for full test suite)
- No external dependencies or API keys required
- Integration tests are developer-only, run manually

### 10.9 Success Criteria
- [ ] Unit test suite completes in under 10 seconds
- [ ] No test requires network access
- [ ] No test requires file system setup beyond reading fixtures
- [ ] 70%+ code coverage from unit tests alone
- [ ] All tests pass consistently without flakiness
- [ ] Tests are easy to understand and modify

## 11. Documentation Requirements

### 11.1 README.md
Must include:
- Project description and features
- Prerequisites (Go, Docker, API keys)
- Quick start guide
- API key setup instructions
- Claude Desktop configuration
- Example queries
- Troubleshooting section

### 11.2 Code Documentation
- Package-level documentation for each package
- Function documentation for exported functions
- Inline comments for complex logic

### 11.3 API Documentation
- Document each tool with examples
- Include sample inputs and outputs
- Document error scenarios

## 12. Claude Desktop Integration

### 12.1 Configuration
Add to Claude Desktop config (`claude_desktop_config.json`):
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

### 12.2 Usage Examples
Example prompts to test:
- "What national parks are in Colorado?"
- "Tell me about Yosemite National Park"
- "Are there any alerts for Grand Canyon?"
- "Find campgrounds near Yellowstone"
- "What's the weather like at Rocky Mountain National Park?"
- "Show me recreation areas in Utah with hiking"

## 13. Future Enhancements (Out of Scope for V1)

- Authentication (OAuth, API keys)
- User preferences and favorites
- Availability checking for campsite reservations
- Distance calculations between locations
- Trail information integration
- Photo gallery integration
- Multi-language support
- Advanced filtering (price, amenities, accessibility)
- Webhook notifications for alert updates
- Structure API responses in formats that are optimized for AI tool consumption, including summaries, metadata, and contextual information

## 14. Success Criteria

The project is successful when:
1. All 9 MCP tools are implemented and functional
2. Server runs in Docker container orchestrated by docker-compose
3. Successfully integrates with Claude Desktop
4. Can query all three external APIs
5. Returns accurate, well-formatted responses
6. Handles errors gracefully with helpful messages
7. Documentation is complete and examples work
8. Code passes basic testing requirements

## 15. Getting Started Checklist

For the AI coding agent, implement in this order:

### Phase 1: Project Setup
- [ ] Initialize Go module with `go mod init`
- [ ] Set up project directory structure (cmd/, internal/, pkg/, test/)
- [ ] Create `.gitignore` for Go projects
- [ ] Create `.dockerignore` for Docker builds
- [ ] Set up `.env.example` with all required API keys
- [ ] Create initial `README.md` with project description

### Phase 2: Core Infrastructure
- [ ] Implement configuration management (environment variables + optional YAML)
- [ ] Set up logging framework with configurable log levels
- [ ] Create shared data models/types in `internal/models/`
- [ ] Implement HTTP client utilities with timeout and retry logic
- [ ] Create cache layer with TTL and LRU eviction

### Phase 3: API Client Implementation
- [ ] Implement National Park Service API client
  - [ ] Search parks endpoint
  - [ ] Get park details endpoint
  - [ ] Get alerts endpoint
  - [ ] Get campgrounds endpoint
  - [ ] Get activities endpoint
  - [ ] Get visitor centers endpoint
- [ ] Implement Recreation.gov API client
  - [ ] Search recreation areas endpoint
  - [ ] Get recreation area details endpoint
  - [ ] Search facilities endpoint
  - [ ] Get facility details endpoint
  - [ ] Get campsites endpoint
  - [ ] Get activities endpoint
- [ ] Implement OpenWeatherMap API client
  - [ ] Current weather endpoint
  - [ ] 5-day forecast endpoint

### Phase 4: MCP Server Implementation
- [ ] Implement MCP protocol handler (stdio communication)
- [ ] Define all 9 MCP tools with proper schemas
- [ ] Implement tool handlers:
  - [ ] search_parks
  - [ ] get_park_details
  - [ ] get_park_alerts
  - [ ] search_campgrounds
  - [ ] search_recreation_areas
  - [ ] get_facility_details
  - [ ] get_weather
  - [ ] get_weather_forecast
  - [ ] list_activities
- [ ] Implement error handling with MCP-compliant error responses
- [ ] Add request validation for all tool inputs

### Phase 5: Testing
- [ ] Create test fixtures directory with sample API responses
- [ ] Write unit tests for all API clients (with mocks)
- [ ] Write unit tests for all tool handlers
- [ ] Write unit tests for cache functionality
- [ ] Write unit tests for configuration parsing
- [ ] Verify test suite runs in under 10 seconds
- [ ] Achieve 70%+ code coverage
- [ ] Create Makefile with test targets
- [ ] (Optional) Create integration tests with `+build integration` tag

### Phase 6: Containerization
- [ ] Create multi-stage Dockerfile with non-root user
- [ ] Create docker-compose.yml with security options
- [ ] Test Docker build process
- [ ] Test running container with docker-compose
- [ ] Verify environment variables are properly passed
- [ ] Test volume mounts for config and cache

### Phase 7: Documentation
- [ ] Write comprehensive README.md:
  - [ ] Project description and features
  - [ ] Prerequisites (Go version, Docker, API keys)
  - [ ] API key setup instructions with links
  - [ ] Quick start guide
  - [ ] Claude Desktop configuration instructions
  - [ ] Example queries and use cases
  - [ ] Troubleshooting section
  - [ ] Development guide
  - [ ] Testing instructions
- [ ] Document all exported functions and packages
- [ ] Add inline comments for complex logic
- [ ] Create API documentation for each tool with examples
- [ ] Document error codes and messages
- [x] Add LICENSE file

### Phase 8: Integration and Testing
- [ ] Configure Claude Desktop with the MCP server
- [ ] Test all 9 tools end-to-end with Claude Desktop
- [ ] Test error scenarios (invalid API keys, network failures)
- [ ] Test with various query patterns
- [ ] Verify error messages are helpful and actionable
- [ ] Test caching behavior
- [ ] Validate response formats are correct

### Phase 9: Polish and Validation
- [ ] Run `go fmt` on all code
- [ ] Run `go vet` to catch common issues
- [ ] Run `golint` or `staticcheck` for style issues
- [ ] Review all TODO comments and address or document
- [ ] Verify all environment variables are documented
- [ ] Test with all three API keys individually (ensure graceful degradation)
- [ ] Review security settings in Dockerfile and docker-compose
- [ ] Final review of README accuracy
- [ ] Tag initial release (v1.0.0)

### Verification Checklist
Before considering the project complete:
- [ ] All 9 MCP tools return valid responses
- [ ] Server successfully runs in Docker container
- [ ] docker-compose starts server without errors
- [ ] Claude Desktop successfully connects to the server
- [ ] Can query national parks by state
- [ ] Can get detailed park information
- [ ] Can retrieve park alerts
- [ ] Can search campgrounds
- [ ] Can search recreation areas
- [ ] Can get facility details
- [ ] Can get current weather
- [ ] Can get weather forecast
- [ ] Can list activities
- [ ] All API errors are handled gracefully
- [ ] Cache reduces redundant API calls
- [ ] Documentation examples work as written
- [ ] Test suite passes completely
- [ ] No root user in container
- [ ] Container runs with read-only filesystem