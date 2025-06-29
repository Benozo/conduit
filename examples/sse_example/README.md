# HTTP/SSE Example

This example demonstrates how to create a Conduit server that runs in HTTP mode with Server-Sent Events (SSE) support for web applications and real-time integrations.

## What is HTTP/SSE Mode?

HTTP/SSE mode provides:
- **HTTP API**: RESTful endpoints for tool execution
- **Server-Sent Events (SSE)**: Real-time streaming capabilities
- **Web Integration**: Perfect for web applications and dashboards
- **CORS Support**: Cross-origin requests enabled
- **Real-time Communication**: Live updates and streaming responses

## Running the Example

```bash
# Build the example
go build -o sse-server .

# Run the server
./sse-server

# Open your browser
open http://localhost:8090/demo
```

## Available Endpoints

### Core MCP Endpoints

- **`GET /schema`** - List all available tools and their schemas
- **`POST /mcp`** - MCP protocol endpoint with SSE support
- **`GET /health`** - Health check endpoint

### Demo Endpoints

- **`GET /demo`** - Interactive demo page with JavaScript examples
- **`GET /sse-test`** - SSE streaming test endpoint

## API Usage Examples

### 1. Get Available Tools

```bash
curl http://localhost:8090/schema
```

### 2. Call a Tool via HTTP

```bash
curl -X POST http://localhost:8090/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "sse_demo",
      "arguments": {}
    }
  }'
```

### 3. Test Server-Sent Events

```bash
curl -N http://localhost:8090/sse-test
```

### 4. Health Check

```bash
curl http://localhost:8090/health
```

## Web Application Integration

### JavaScript Example

```javascript
// Tool execution via fetch
async function callTool(name, args) {
  const response = await fetch('/mcp', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: 1,
      method: 'tools/call',
      params: {name: name, arguments: args}
    })
  });
  return response.json();
}

// SSE streaming
const eventSource = new EventSource('/sse-test');
eventSource.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};
```

### React Integration Example

```jsx
import { useState, useEffect } from 'react';

function ConduitTools() {
  const [tools, setTools] = useState([]);
  const [result, setResult] = useState(null);

  useEffect(() => {
    // Load available tools
    fetch('/schema')
      .then(r => r.json())
      .then(data => setTools(data.tools));
  }, []);

  const callTool = async (toolName, args) => {
    const response = await fetch('/mcp', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        jsonrpc: '2.0',
        id: 1,
        method: 'tools/call',
        params: {name: toolName, arguments: args}
      })
    });
    const result = await response.json();
    setResult(result);
  };

  return (
    <div>
      <h2>Available Tools</h2>
      {tools.map(tool => (
        <button key={tool.name} onClick={() => callTool(tool.name, {})}>
          {tool.name}
        </button>
      ))}
      {result && <pre>{JSON.stringify(result, null, 2)}</pre>}
    </div>
  );
}
```

## Available Tools

This example provides all 33 tools:

### Standard Tools (31)
- Text manipulation tools (uppercase, lowercase, trim, etc.)
- Memory management tools (remember, recall, forget, etc.)
- Utility tools (base64, JSON, hashing, etc.)

### Demo Tools (2)
- `sse_demo` - Demonstrates HTTP/SSE integration
- `streaming_demo` - Shows data preparation for streaming

## Configuration Options

```go
config := conduit.DefaultConfig()
config.Mode = mcp.ModeHTTP
config.Port = 8090
config.EnableCORS = true      // Allow cross-origin requests
config.EnableHTTPS = false    // Use HTTP (set to true for HTTPS)
config.EnableLogging = true   // Enable request logging
```

## Use Cases

### Web Applications
- Interactive dashboards
- Real-time data processing
- Tool integration in web UIs
- Live updates and notifications

### API Integrations
- Webhook endpoints
- External service integration
- Batch processing APIs
- Monitoring and alerting

### Development & Testing
- Tool testing and validation
- Interactive debugging
- Performance monitoring
- Integration testing

## Advanced Features

### Custom Endpoints
You can add custom HTTP endpoints alongside the MCP endpoints:

```go
http.HandleFunc("/custom", func(w http.ResponseWriter, r *http.Request) {
    // Custom endpoint logic
})
```

### Streaming Responses
The SSE support allows for real-time streaming of tool results and updates.

### CORS Configuration
Cross-origin requests are supported for web application integration.

This example demonstrates the power of Conduit's HTTP/SSE mode for building web-integrated AI tool servers.
