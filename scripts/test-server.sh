#!/bin/bash
# Test script to verify MCP server functionality

echo "Testing Recreation MCP Server"
echo "=============================="
echo ""

# Test 1: Initialize request
echo "Test 1: Testing MCP initialization..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' | ./mcp-server 2>/dev/null | jq -r 'select(.result != null) | "✓ Server initialized: \(.result.serverInfo.name) v\(.result.serverInfo.version)"'

echo ""
echo "Test 2: Listing available tools..."
# We need to initialize first, then send notifications/initialized, then list tools
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
  sleep 0.1
  echo '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}'
  sleep 0.1
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
  sleep 0.5
) | ./mcp-server 2>/dev/null | jq -r 'select(.result.tools != null) | "✓ Found \(.result.tools | length) tools:", (.result.tools[].name | "  - \(.)")'

echo ""
echo "Test 3: Testing a tool call (search_parks for California)..."
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
  sleep 0.1
  echo '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}'
  sleep 0.1
  echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"search_parks","arguments":{"state":"CA","limit":3}}}'
  sleep 2
) | ./mcp-server 2>/dev/null | jq -r 'select(.id == 3) | if .result then "✓ search_parks returned \(.result.content[0].text | fromjson | .count) parks" else "✗ Error: \(.error.message)" end'

echo ""
echo "Test 4: Testing weather tool (Yosemite coordinates)..."
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
  sleep 0.1
  echo '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}'
  sleep 0.1
  echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_weather","arguments":{"latitude":37.8651,"longitude":-119.5383,"units":"imperial"}}}'
  sleep 2
) | ./mcp-server 2>/dev/null | jq -r 'select(.id == 4) | if .result then "✓ get_weather returned temperature: \(.result.content[0].text | fromjson | .temp)°F" else "✗ Error: \(.error.message)" end'

echo ""
echo "=============================="
echo "All tests completed!"
