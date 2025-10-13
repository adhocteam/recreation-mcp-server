#!/bin/bash

# Test search_parks with proper JSON-RPC protocol
echo "Testing search_parks with CA state code..."

# Send requests to server
RESPONSE=$(cat <<'EOF' | timeout 5 ./mcp-server 2>&1 | grep -A 50 '"id": 2'
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "0.1.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "search_parks", "arguments": {"state_code": "CA", "limit": 3}}, "id": 2}
EOF
)

echo "$RESPONSE" | head -20
