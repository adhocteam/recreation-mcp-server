# AI Agent Guidelines for Recreation MCP Server

This document provides guidance for AI coding assistants working on this project. It captures the conventions, patterns, and practices established during development.

## Project Philosophy

**Clean, Simple, Fast** - The guiding principle for all code:
- **Clean**: Easy to read and understand
- **Simple**: Minimal complexity, straightforward solutions
- **Fast**: Performance matters, but clarity comes first

## Architecture Patterns

### Project Structure
```
recreation-mcp-server/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── api/            # External API clients (NPS, Recreation.gov, Weather)
│   ├── cache/          # Caching layer with TTL/LRU
│   ├── config/         # Configuration management
│   ├── mcp/            # MCP protocol implementation
│   └── models/         # Shared data types
├── pkg/util/           # Public utility packages
├── testdata/           # Test fixtures (JSON responses)
├── docs/               # GitHub Pages website
└── scripts/            # Helper scripts for testing
```

### Code Organization Principles

1. **Internal vs Public**: All application logic lives in `internal/`. Only truly reusable utilities go in `pkg/`.

2. **API Clients**: Each external API gets its own file in `internal/api/`:
   - `nps.go` - National Park Service
   - `recreation.go` - Recreation.gov
   - `weather.go` - OpenWeatherMap

3. **MCP Layer**: Separated into two files:
   - `server.go` - Server setup and tool registration
   - `handlers.go` - Tool implementations and schemas

4. **Test Files**: Co-located with source files using `_test.go` suffix

## Code Conventions

### Error Handling
- Always return errors, don't panic
- Wrap errors with context: `fmt.Errorf("failed to fetch parks: %w", err)`
- Log errors before returning them to users
- Use MCP-compliant error responses with error codes

### Naming Conventions
- **Files**: Snake case (e.g., `http_client.go`)
- **Packages**: Single word, lowercase (e.g., `cache`, `config`)
- **Structs**: PascalCase (e.g., `SearchParksInput`)
- **Functions**: camelCase for private, PascalCase for exported
- **Constants**: PascalCase (e.g., `DefaultCacheTTL`)

### Configuration Management
- Environment variables for secrets and runtime config
- Optional YAML for advanced configuration
- All variables documented in `.env.example`
- Use `internal/config` package for centralized access

### API Client Patterns
```go
// All API clients follow this pattern:
type NPSClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
    cache      *cache.Cache
    logger     *slog.Logger
}

func NewNPSClient(cfg *config.Config, cache *cache.Cache, logger *slog.Logger) *NPSClient
func (c *NPSClient) SearchParks(ctx context.Context, params SearchParams) ([]Park, error)
```

### HTTP Client Standards
- Always use context for cancellation
- 30-second timeout for NPS/Recreation.gov
- 15-second timeout for weather APIs
- Retry with exponential backoff on 429/5xx
- Cache GET requests by default

### Caching Strategy
- Cache key format: `{endpoint}:{params_hash}`
- Default TTL: 1 hour
- LRU eviction when size limit reached
- Bypass cache via configuration

## Testing Standards

### Unit Testing Philosophy
Tests must be **fast, simple, and reliable**:
- Entire test suite runs in < 10 seconds
- No network calls (use mocks/fixtures)
- No filesystem dependencies (except reading fixtures)
- No flaky tests - remove or fix them

### Test Organization
- Use table-driven tests for multiple scenarios
- Test fixtures in `testdata/` directory (JSON files)
- Mock external APIs using interfaces
- Integration tests tagged with `// +build integration`

### Test Coverage Targets
- Overall: 25%+ (achieved 25.3%)
- Critical paths: 70%+
- Don't test trivial getters/setters

### Example Test Pattern
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

## Docker & Security

### Container Security
- Multi-stage builds (golang:1.25.10-alpine3.22 -> alpine:3.22.2)
- Non-root user (UID 1000)
- Read-only filesystem with tmpfs for /tmp
- Security flags: `-ldflags='-w -s -extldflags "-static"'`
- Minimal final image (~30MB)

### docker-compose Security
```yaml
security_opt:
  - no-new-privileges:true
read_only: true
tmpfs:
  - /tmp
user: "1000:1000"
```

## Documentation Standards

### Code Documentation
- Package-level documentation for each package
- Exported functions must have doc comments
- Complex logic requires inline comments
- Use full sentences in doc comments

### README Structure
1. Project description with badges
2. Table of contents
3. Features
4. Quick start
5. Configuration (with examples)
6. Claude Desktop integration
7. Development guide
8. API documentation (or link to API.md)
9. Troubleshooting
10. License

### API Documentation
- Each tool documented in `API.md`
- Include input schema with examples
- Show sample output
- Document error scenarios
- Provide usage tips

## MCP Tool Design

### Tool Descriptions
- Clear, concise descriptions
- Include usage guidance and examples
- Suggest parameters (e.g., "Use 3-5 limit for detailed planning")
- Help LLMs choose the right tool

### Schema Documentation
```go
// Good: Includes examples and guidance
"query": {
    Type:        "string",
    Description: "Search text for parks (e.g., 'yose' for Yosemite, 'grand canyon')",
}

// Bad: Too vague
"query": {
    Type:        "string",
    Description: "Search query",
}
```

### Tool Composition
- Design tools to work together
- Encourage multi-tool queries
- Return manageable result sets (default limit: 10)
- Include coordinates for weather lookups

## Git Workflow

### Commit Messages
- Use conventional format: `Phase X: Brief description`
- Examples:
  - `Phase 1: Initialize Go module and project structure`
  - `Phase 6: Add Docker containerization with security hardening`
  - `Phase 9 complete: Polish and validation`

### Version Tags
- Follow semantic versioning (MAJOR.MINOR.PATCH)
- Annotated tags: `git tag -a v1.0.0 -m "Initial release"`
- Tag format: `vX.Y.Z`

### What to Ignore
Always in `.gitignore`:
- Compiled binaries (`/mcp-server`, `/recreation-mcp-server`, `/test-server`)
- Coverage files (`coverage.out`, `*.html`)
- Environment files (`.env`)
- IDE files (`.vscode/`, `.idea/`)
- OS files (`.DS_Store`)

## Quality Checks

### Pre-Release Checklist
Run these commands before any release:
```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint code
golangci-lint run ./...

# Run tests
go test -v -race ./...

# Check coverage
go test -coverprofile=coverage.out ./...
```

All must pass with zero issues.

### Static Analysis
- Use golangci-lint with version pinning (v2.5.0)
- Configuration in `.golangci.yml`
- No warnings allowed in main branch

## Performance Considerations

### Optimization Strategy
1. **Correctness first**: Make it work
2. **Clarity second**: Make it readable
3. **Performance third**: Make it fast (if needed)

### Known Performance Metrics
- Binary size: ~11MB (static build)
- Container size: ~30MB (alpine-based)
- Test suite: < 1 second (0.847s achieved)
- API response time: < 5 seconds (with caching)

## Common Pitfalls to Avoid

1. **Don't use root in containers** - Always run as non-root user
2. **Don't skip context** - All HTTP requests need context.Context
3. **Don't ignore errors** - Check every error return value
4. **Don't hardcode timeouts** - Use configuration
5. **Don't test implementation details** - Test behavior, not internals
6. **Don't commit secrets** - Use environment variables
7. **Don't make breaking changes** - Follow semantic versioning

## Adding New Features

### Adding a New API Endpoint
1. Add method to appropriate client in `internal/api/`
2. Add test with fixture in `testdata/`
3. Update `internal/models/types.go` if needed
4. Document in `API.md`

### Adding a New MCP Tool
1. Register tool in `internal/mcp/server.go`
2. Implement handler in `internal/mcp/handlers.go`
3. Add input/output types to `internal/mcp/handlers.go`
4. Add tests with mocked API responses
5. Document in `API.md` with examples
6. Update README tool count

### Adding Dependencies
```bash
go get <package>
go mod tidy
```
- Prefer standard library when possible
- Vet dependencies for security/maintenance
- Pin versions for reproducible builds

## Debugging Tips

### Testing Individual Tools
Use test scripts in `scripts/` directory:
```bash
./scripts/test-search-parks.sh
./scripts/test-server.sh
```

### Viewing Logs
```bash
# Set log level
export LOG_LEVEL=debug

# Run with verbose output
go run cmd/server/main.go
```

### Common Issues
- **"API key not found"**: Check `.env` file
- **"Connection timeout"**: Check network/firewall
- **"Cache miss"**: Expected on first request
- **"Tool not found"**: Check MCP tool registration

## Resources

- **MCP Protocol**: https://modelcontextprotocol.io/
- **Go MCP SDK**: https://github.com/modelcontextprotocol/go-sdk
- **NPS API Docs**: https://www.nps.gov/subjects/developer/api-documentation.htm
- **Recreation.gov API**: https://ridb.recreation.gov/docs
- **OpenWeather API**: https://openweathermap.org/api

## Questions?

When in doubt:
1. Check existing code for patterns
2. Favor simplicity over cleverness
3. Write tests that demonstrate intent
4. Document non-obvious decisions
5. Ask for clarification before making architectural changes

---

**Last Updated**: October 14, 2025 (v1.0.0)
**Maintained By**: Recreation MCP Server Contributors
