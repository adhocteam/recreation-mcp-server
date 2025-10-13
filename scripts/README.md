# Test Scripts

Quick test scripts for verifying MCP server functionality.

## Scripts

- **`test-server.sh`** - Main test suite that runs 4 tests: initialization, tool listing, search_parks, and get_weather
- **`test-search-parks.sh`** - Focused test for the search_parks tool
- **`test-search-detailed.sh`** - Detailed search_parks test with verbose output

## Usage

Make sure the server is built first:

```bash
go build -o mcp-server ./cmd/server
```

Then run any test script:

```bash
./scripts/test-server.sh
```

These scripts communicate with the MCP server using JSON-RPC over stdio, simulating how a client like Claude Desktop would interact with the server.
