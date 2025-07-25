#!/bin/bash

echo "=== MCP STDIO Server Demo ==="
echo "Building and testing the basic MCP server with STDIO transport..."
echo ""

# Build the server
cd /home/engineone/Projects/AI/ConduitMCP/mcpV2/examples/basic_server
go build -o basic_server .

echo "✅ Server built successfully!"
echo ""
echo "🔄 Testing STDIO communication..."
echo ""

# Test sequence function
run_test() {
    local test_name="$1"
    local json_request="$2"
    
    echo "📤 Sending: $test_name"
    echo "   Request: $json_request"
    echo "📥 Response:"
    
    # Send request and capture response
    echo "$json_request" | timeout 3s ./basic_server 2>/dev/null | head -1
    echo ""
}

# Test 1: Initialize
run_test "Initialize Connection" \
  '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "demo-client", "version": "1.0.0"}}}'

# Test 2: List Tools
run_test "List Available Tools" \
  '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'

# Test 3: Echo Tool
run_test "Echo Tool Test" \
  '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "echo", "arguments": {"message": "Hello STDIO MCP!"}}}'

# Test 4: Text Transform
run_test "Text Transform (Uppercase)" \
  '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "text_transform", "arguments": {"text": "hello world", "operation": "uppercase"}}}'

# Test 5: Calculator
run_test "Calculator (Addition)" \
  '{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "calculator", "arguments": {"operation": "add", "a": 10, "b": 5}}}'

echo "✅ STDIO Demo Complete!"
echo ""
echo "🎯 Key Points:"
echo "   • Server communicates via stdin/stdout using JSON-RPC"
echo "   • Each request gets a corresponding response"
echo "   • Tools are called by name with typed parameters"
echo "   • Error handling provides clear feedback"
echo ""
echo "🔗 Integration:"
echo "   • Use with Claude Desktop, VS Code Copilot, or any MCP client"
echo "   • Add the server path to your MCP client configuration"
echo "   • The client will automatically handle JSON-RPC communication"
