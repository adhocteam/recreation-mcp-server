#!/bin/bash

# Test search_parks with proper state code
echo "Testing search_parks..."

# Start the server in the background
./mcp-server &
PID=$!

# Give it a moment to start
sleep 1

# Call the tool and capture response
RESPONSE=$(cat <<EOF | timeout 5 ./mcp-server 2>/dev/null
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "0.1.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "search_parks", "arguments": {"state_code": "CA", "limit": 5}}, "id": 2}
EOF
)

echo "Response:"
echo "$RESPONSE" | tail -1 | jq '.'

# Kill background server if still running
kill $PID 2>/dev/null
