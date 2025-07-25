#!/bin/bash
echo "Starting basic_server for STDIO test..."

# Start server and send messages
(
    echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}'
    sleep 0.1
    echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'
    sleep 0.1
    echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "echo", "arguments": {"message": "Hello STDIO MCP!"}}}'
    sleep 0.1
) | timeout 5s ./basic_server
