# HTTP Server Example

This example demonstrates how to create an MCP server using HTTP/SSE transport with middleware integration. Unlike STDIO transport, HTTP transport allows for web-based integration and provides better scalability for multiple clients.

## What is HTTP/SSE Transport?

HTTP/SSE (Server-Sent Events) transport provides:
- **HTTP POST** endpoints for JSON-RPC requests
- **Server-Sent Events** for real-time communication
- **Web integration** capabilities
- **Multiple client support** 
- **Load balancing** compatibility

This makes it suitable for:
- **Web applications**
- **Microservices architectures**
- **Multi-client scenarios**
- **Cloud deployments**

## Features

- ✅ **HTTP/SSE Transport** - Web-based MCP communication
- ✅ **Middleware Integration** - Logging, metrics, rate limiting, validation
- ✅ **Text Transform Tool** - Text manipulation operations
- ✅ **Echo Tool** - Testing tool with timestamps
- ✅ **Rate Limiting** - 10 requests per minute per client
- ✅ **Graceful Shutdown** - Proper HTTP server lifecycle management
- ✅ **Error Handling** - Comprehensive error responses

## Quick Start

### 1. Build and Run the Server

```bash
cd examples/http_server
go build .
./http_server
```

The server will start on `http://localhost:8081` with these endpoints:
- `POST /mcp` - JSON-RPC requests
- `GET /mcp/events` - Server-Sent Events

### 2. Test with curl

#### Initialize Connection
```bash
curl -X POST http://localhost:8081/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2025-03-26",
      "capabilities": {},
      "clientInfo": {"name": "test-client", "version": "1.0.0"}
    }
  }'
```

#### List Tools
```bash
curl -X POST http://localhost:8081/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }'
```

#### Call Text Transform Tool
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "text_transform",
      "arguments": {
        "text": "Hello World",
        "operation": "uppercase"
      }
    }
  }'
```

#### Call Echo Tool
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {
        "message": "Hello HTTP MCP!"
      }
    }
  }'
```

### 3. Test Server-Sent Events
```bash
curl -N http://localhost:8081/mcp/events
```

## Available Tools

### 1. Text Transform (`text_transform`)
Transform text using various operations.

**Parameters:**
- `text` (string, required) - The text to transform
- `operation` (string, required) - Operation: "uppercase", "lowercase", "reverse"

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "text_transform",
    "arguments": {
      "text": "hello world",
      "operation": "uppercase"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "HELLO WORLD"
      }
    ]
  }
}
```

### 2. Echo (`echo`)
Echo back the provided message with a timestamp.

**Parameters:**
- `message` (string, required) - Message to echo back

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "message": "Hello MCP!"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Echo at 2025-07-25T15:30:45Z: Hello MCP!"
      }
    ]
  }
}
```

## Middleware Features

The HTTP server demonstrates several middleware components:

### 1. Logging Middleware
Logs all incoming requests and outgoing responses:
```
[INFO] Request received: tools/call
[DEBUG] Processing request with ID: 3
[INFO] Response sent: success
```

### 2. Metrics Middleware
Collects performance metrics:
```
[METRICS] Counter request_count: {method: "tools/call"}
[METRICS] Duration request_duration: 15ms ({method: "tools/call"})
[METRICS] Gauge active_connections: 2.0
```

### 3. Rate Limiting Middleware
Limits requests to 10 per minute per client:
- Uses client IP for identification
- Returns HTTP 429 when limit exceeded
- Sliding window implementation

### 4. Validation Middleware
Validates JSON-RPC format and parameters:
- Schema validation for tool parameters
- Required field checking
- Type validation

### 5. Error Handling Middleware
Provides consistent error responses:
- Structured error messages
- HTTP status code mapping
- Error logging and tracking

## Architecture

```
┌─────────────────┐
│   HTTP Client   │
│ (Browser/curl)  │
└─────────┬───────┘
          │ HTTP/JSON-RPC
┌─────────▼───────┐
│  HTTP Transport │ (:8080)
│  POST /mcp      │
│  GET /mcp/events│
└─────────┬───────┘
          │
┌─────────▼───────┐
│  Middleware     │ (Logging, Metrics, 
│     Chain       │  Rate Limit, etc.)
└─────────┬───────┘
          │
┌─────────▼───────┐
│  MCP Server     │ (Tool Registration
│    Core         │  & Request Routing)
└─────────────────┘
```

## Configuration

### Server Options
```go
srv := server.NewServer(&server.ServerOptions{
    Info: protocol.Implementation{
        Name:    "http-mcp-server",
        Version: "1.0.0",
    },
    Capabilities: protocol.ServerCapabilities{
        Tools: &protocol.ToolsCapability{
            ListChanged: true,
        },
    },
})
```

### Middleware Chain
```go
middlewareChain := middleware.NewChain(
    middleware.LoggingMiddleware(logger),
    middleware.MetricsMiddleware(metrics),
    middleware.RateLimitMiddleware(rateLimiter),
    middleware.ValidationMiddleware(),
    middleware.ErrorHandlingMiddleware(logger),
)
```

### Rate Limiting
```go
rateLimiter := NewSimpleRateLimiter(10, time.Minute) // 10 requests per minute
```

## HTTP Endpoints

### POST /mcp
Main JSON-RPC endpoint for tool calls and operations.

**Headers:**
- `Content-Type: application/json`
- `Accept: application/json`

**Body:** JSON-RPC 2.0 message

### GET /mcp/events
Server-Sent Events endpoint for real-time updates.

**Headers:**
- `Accept: text/event-stream`
- `Cache-Control: no-cache`

**Response:** SSE stream

## Error Handling

The server provides comprehensive error handling:

### HTTP Status Codes
- `200` - Success
- `400` - Bad Request (invalid JSON-RPC)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

### JSON-RPC Error Codes
- `-32700` - Parse error
- `-32600` - Invalid request
- `-32601` - Method not found
- `-32602` - Invalid params
- `-32603` - Internal error

### Example Error Response
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "details": "text parameter must be a string"
    }
  }
}
```

## Testing

### Manual Testing Script
Create a test script to verify functionality:

```bash
#!/bin/bash
echo "Testing HTTP MCP Server"

# Start server in background
./http_server &
SERVER_PID=$!

sleep 2

# Test tools/list
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}'

# Test echo tool
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "echo", "arguments": {"message": "test"}}}'

# Cleanup
kill $SERVER_PID
```

### Load Testing
Use tools like `ab` or `wrk` for load testing:

```bash
# Apache Bench
ab -n 1000 -c 10 -T application/json -p request.json http://localhost:8081/mcp

# wrk
wrk -t10 -c100 -d30s --script=test.lua http://localhost:8081/mcp
```

## Production Considerations

### Security
1. **HTTPS/TLS** - Use TLS in production
2. **Authentication** - Add auth middleware
3. **CORS** - Configure CORS for web clients
4. **Input Validation** - Sanitize all inputs

### Performance
1. **Connection Pooling** - Reuse connections
2. **Caching** - Cache responses where appropriate
3. **Compression** - Enable gzip compression
4. **Load Balancing** - Use reverse proxy

### Monitoring
1. **Health Checks** - Add `/health` endpoint
2. **Metrics Export** - Prometheus integration
3. **Structured Logging** - JSON logs for parsing
4. **Distributed Tracing** - OpenTelemetry integration

### Deployment
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o http_server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/http_server .
EXPOSE 8080
CMD ["./http_server"]
```

### Environment Variables
```bash
export MCP_PORT=8080
export MCP_LOG_LEVEL=info
export MCP_RATE_LIMIT=100
export MCP_ENABLE_METRICS=true
```

## Integration Examples

### JavaScript/Browser
```javascript
async function callMCPTool(tool, args) {
  const response = await fetch('/mcp', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: Date.now(),
      method: 'tools/call',
      params: { name: tool, arguments: args }
    })
  });
  return response.json();
}

// Usage
const result = await callMCPTool('echo', { message: 'Hello!' });
```

### Python
```python
import requests

def call_mcp_tool(tool, args):
    response = requests.post('http://localhost:8080/mcp', json={
        'jsonrpc': '2.0',
        'id': 1,
        'method': 'tools/call',
        'params': {'name': tool, 'arguments': args}
    })
    return response.json()

# Usage
result = call_mcp_tool('echo', {'message': 'Hello!'})
```

## Related Examples

- **[`basic_server/`](../basic_server/)** - STDIO server for simple integration
- **[`advanced_server/`](../advanced_server/)** - Full-featured server with all capabilities
- **[`websocket_client/`](../websocket_client/)** - WebSocket client example
- **[`pure_library/`](../pure_library/)** - Pure library usage

This HTTP server example provides a solid foundation for building web-integrated MCP services with comprehensive middleware support and production-ready features.
