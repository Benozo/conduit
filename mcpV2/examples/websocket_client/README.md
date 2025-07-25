# WebSocket Client Example

This example demonstrates how to create an MCP client using **WebSocket transport** with progress tracking capabilities. WebSocket transport provides real-time, full-duplex communication ideal for interactive applications and long-running operations.

## What is WebSocket Transport?

WebSocket transport for MCP provides:
- ✅ **Real-time Communication** - Full-duplex messaging
- ✅ **Progress Tracking** - Live updates for long-running operations  
- ✅ **Interactive Sessions** - Bidirectional client-server communication
- ✅ **Persistent Connection** - Maintains connection state
- ✅ **Event Notifications** - Server-initiated messages
- ✅ **Low Latency** - Minimal protocol overhead

This makes it perfect for:
- **Interactive Dashboards** - Real-time data updates
- **Desktop Applications** - Rich client interfaces
- **Development Tools** - Live debugging and monitoring
- **Gaming Applications** - Real-time multiplayer coordination
- **AI Assistants** - Streaming responses and progress updates

## Features

- ✅ **WebSocket Transport** - Full-duplex real-time communication
- ✅ **Progress Tracking** - Live progress updates for long operations
- ✅ **Tool Management** - List and call server tools
- ✅ **Error Handling** - Comprehensive error recovery
- ✅ **Connection Management** - Automatic reconnection support
- ✅ **Custom Headers** - Authentication and authorization support
- ✅ **Timeout Configuration** - Configurable connection and operation timeouts

## Quick Start

### 1. Start a WebSocket Server

First, you need a WebSocket-enabled MCP server. You can use any of these approaches:

#### Option A: Use the Included Mock Server (Recommended)
```bash
# Terminal 1: Start the mock WebSocket server
cd cmd && go run mock_server.go
# Server will be available at ws://localhost:8082/mcp

# Terminal 2: Run the WebSocket client (from main directory)
cd .. && go run main.go
```

#### Option B: Use Your Own WebSocket Server
```bash
# Start your own WebSocket MCP server on ws://localhost:8082/mcp
```

#### Option C: Use the Test Script (Automated)
```bash
# Automatically starts mock server and runs client
./test_websocket.sh
```

### 2. Build and Run the Client

```bash
cd examples/websocket_client
go build .
./websocket_client
```

### 3. Expected Output

```
Connecting to WebSocket MCP server...
Connected to MCP server successfully!

=== Tool Listing ===
Available tools (2):
  - text_transform: Transform text using various operations
  - echo: Echo back the input with a timestamp

=== Tool Calling ===
Text transform result: HELLO, WEBSOCKET!
Echo result: Echo at 2025-07-25T16:30:45+10:00: Testing WebSocket transport

=== Progress Tracking ===
Progress [demo_task_001]: 10.0% (1/10)
Progress [demo_task_001]: 20.0% (2/10)
Progress [demo_task_001]: 30.0% (3/10)
...
Progress [demo_task_001]: 100.0% (10/10)
Long task completed: Task completed with 10 steps

WebSocket client example completed!
```

## Core Components

### 1. WebSocket Transport

The WebSocket transport handles the connection and messaging:

```go
// Create WebSocket transport
wsTransport := transport.NewWebSocketTransport("ws://localhost:8082/mcp")

// Optional: Set custom headers for authentication
wsTransport.SetHeader("Authorization", "Bearer token")

// Connect to server
err := wsTransport.Connect(ctx)
```

### 2. Client Configuration

Configure the client with timeouts and progress handling:

```go
mcpClient := client.NewClient(wsTransport, &client.ClientOptions{
    Timeout:        30 * time.Second,    // Tool call timeout
    ConnectTimeout: 10 * time.Second,    // Connection timeout
    ClientInfo: protocol.Implementation{
        Name:    "websocket-mcp-client",
        Version: "1.0.0",
    },
    ProgressHandler: handleProgress,     // Progress callback
})
```

### 3. Progress Tracking

Handle real-time progress updates for long-running operations:

```go
func handleProgress(token string, progress float64, total int64) {
    percentage := progress * 100
    log.Printf("Progress [%s]: %.1f%% (%d/%d)", 
        token, percentage, 
        int64(progress*float64(total)), total)
}
```

## Client Operations

### 1. Connection Management

```go
// Connect to server
err := wsTransport.Connect(ctx)
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}

// Perform MCP handshake
capabilities := protocol.ClientCapabilities{
    Experimental: map[string]interface{}{
        "progressNotifications": true,
    },
}

err = mcpClient.Connect(ctx, capabilities)
if err != nil {
    log.Fatalf("Failed to connect to MCP server: %v", err)
}
```

### 2. Tool Discovery

```go
// List available tools
tools, err := client.ListTools(ctx)
if err != nil {
    log.Printf("Failed to list tools: %v", err)
    return
}

log.Printf("Available tools (%d):", len(tools))
for _, tool := range tools {
    log.Printf("  - %s: %s", tool.Name, tool.Description)
}
```

### 3. Tool Execution

```go
// Call a tool
result, err := client.CallTool(ctx, "text_transform", map[string]interface{}{
    "text":      "Hello, WebSocket!",
    "operation": "uppercase",
})

if err != nil {
    log.Printf("Tool call failed: %v", err)
    return
}

log.Printf("Result: %s", result.Content[0].Text)
```

### 4. Progress-Enabled Operations

```go
// Create context with progress tracking
progressCtx := context.WithValue(ctx, "progressToken", "task_001")

// Call tool with progress tracking
result, err := client.CallTool(progressCtx, "long_task", map[string]interface{}{
    "duration": 5,  // 5 seconds
    "steps":    10, // 10 progress updates
})
```

## Available Operations

The WebSocket client supports all standard MCP operations:

### Tool Management
- **List Tools** - `ListTools(ctx)` - Get available tools
- **Call Tool** - `CallTool(ctx, name, params)` - Execute tools
- **Get Tool Schema** - Access tool input/output schemas

### Resource Management (when supported by server)
- **List Resources** - `ListResources(ctx)` - Get available resources
- **Read Resource** - `ReadResource(ctx, uri)` - Read resource content

### Prompt Management (when supported by server)
- **List Prompts** - `ListPrompts(ctx)` - Get available prompts
- **Get Prompt** - `GetPrompt(ctx, name, args)` - Get prompt content

## WebSocket Configuration

### Connection Options

```go
// Basic connection
transport := transport.NewWebSocketTransport("ws://localhost:8082/mcp")

// With custom headers
transport.SetHeader("Authorization", "Bearer your-token")
transport.SetHeader("X-Client-Version", "1.0.0")

// With TLS (wss://)
transport := transport.NewWebSocketTransport("wss://secure.example.com/mcp")
```

### Client Timeouts

```go
clientOptions := &client.ClientOptions{
    Timeout:        30 * time.Second,    // Individual operation timeout
    ConnectTimeout: 10 * time.Second,    // Connection establishment timeout
    RetryAttempts:  3,                   // Connection retry attempts
    RetryDelay:     time.Second,         // Delay between retries
}
```

## Progress Tracking Deep Dive

### Progress Handler

The progress handler receives updates for long-running operations:

```go
func handleProgress(token string, progress float64, total int64) {
    // token: unique identifier for the operation
    // progress: completion ratio (0.0 to 1.0)
    // total: total number of items/steps
    
    percentage := progress * 100
    completed := int64(progress * float64(total))
    
    log.Printf("Operation %s: %.1f%% complete (%d/%d)", 
        token, percentage, completed, total)
}
```

### Progress-Enabled Context

```go
// Create context with progress token
progressCtx := context.WithValue(ctx, "progressToken", "unique_task_id")

// The server can use this token to send progress updates
result, err := client.CallTool(progressCtx, "long_running_tool", params)
```

### Progress Visualization

```go
func visualProgress(token string, progress float64, total int64) {
    percentage := progress * 100
    barWidth := 50
    filled := int(float64(barWidth) * progress)
    
    bar := strings.Repeat("=", filled) + strings.Repeat("-", barWidth-filled)
    fmt.Printf("\r[%s] %.1f%% %s", token, percentage, bar)
    
    if progress >= 1.0 {
        fmt.Println(" Complete!")
    }
}
```

## Error Handling

### Connection Errors

```go
err := wsTransport.Connect(ctx)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "connection refused"):
        log.Fatal("Server not running or unreachable")
    case strings.Contains(err.Error(), "timeout"):
        log.Fatal("Connection timeout - server may be overloaded")
    default:
        log.Fatalf("Connection failed: %v", err)
    }
}
```

### Tool Call Errors

```go
result, err := client.CallTool(ctx, "tool_name", params)
if err != nil {
    var rpcErr *protocol.RPCError
    if errors.As(err, &rpcErr) {
        switch rpcErr.Code {
        case protocol.MethodNotFound:
            log.Printf("Tool not found: %s", rpcErr.Message)
        case protocol.InvalidParams:
            log.Printf("Invalid parameters: %s", rpcErr.Message)
        default:
            log.Printf("Tool error: %s", rpcErr.Message)
        }
    } else {
        log.Printf("Network error: %v", err)
    }
}
```

### Reconnection Strategy

```go
func connectWithRetry(transport *transport.WebSocketTransport, ctx context.Context, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := transport.Connect(ctx)
        if err == nil {
            return nil
        }
        
        log.Printf("Connection attempt %d failed: %v", i+1, err)
        
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * time.Second)
        }
    }
    
    return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}
```

## Testing and Development

### Local Testing

```bash
# Terminal 1: Start a test server (if available)
cd ../http_server
./http_server

# Terminal 2: Run the WebSocket client
cd ../websocket_client
./websocket_client
```

### Mock WebSocket Server

For testing, you can create a simple WebSocket server:

```go
func createTestServer() {
    http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
        upgrader := websocket.Upgrader{}
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            return
        }
        defer conn.Close()
        
        // Handle MCP protocol messages
        for {
            var msg protocol.JSONRPCMessage
            err := conn.ReadJSON(&msg)
            if err != nil {
                break
            }
            
            // Echo back or handle appropriately
            conn.WriteJSON(protocol.JSONRPCMessage{
                Version: "2.0",
                ID:      msg.ID,
                Result:  "test response",
            })
        }
    })
    
    log.Fatal(http.ListenAndServe(":8082", nil))
}
```

## Integration Examples

### Desktop Application

```go
type DesktopApp struct {
    client *client.Client
    ui     *ui.Window
}

func (app *DesktopApp) connectToMCP() {
    transport := transport.NewWebSocketTransport("ws://localhost:8082/mcp")
    
    app.client = client.NewClient(transport, &client.ClientOptions{
        ProgressHandler: func(token string, progress float64, total int64) {
            // Update UI progress bar
            app.ui.SetProgress(token, progress)
        },
    })
    
    // Connect and update UI
    err := app.client.Connect(context.Background(), protocol.ClientCapabilities{})
    if err != nil {
        app.ui.ShowError("Failed to connect to MCP server")
        return
    }
    
    app.ui.ShowStatus("Connected to MCP server")
}
```

### Web Dashboard Backend

```go
type Dashboard struct {
    client     *client.Client
    wsClients  map[string]*websocket.Conn
}

func (d *Dashboard) handleProgress(token string, progress float64, total int64) {
    // Broadcast progress to all connected web clients
    update := map[string]interface{}{
        "type":     "progress",
        "token":    token,
        "progress": progress,
        "total":    total,
    }
    
    for _, conn := range d.wsClients {
        conn.WriteJSON(update)
    }
}
```

## Performance Considerations

### Connection Pooling

```go
type ClientPool struct {
    clients chan *client.Client
    factory func() (*client.Client, error)
}

func NewClientPool(size int, factory func() (*client.Client, error)) *ClientPool {
    pool := &ClientPool{
        clients: make(chan *client.Client, size),
        factory: factory,
    }
    
    // Pre-populate pool
    for i := 0; i < size; i++ {
        if client, err := factory(); err == nil {
            pool.clients <- client
        }
    }
    
    return pool
}

func (p *ClientPool) Get() (*client.Client, error) {
    select {
    case client := <-p.clients:
        return client, nil
    default:
        return p.factory()
    }
}

func (p *ClientPool) Put(client *client.Client) {
    select {
    case p.clients <- client:
    default:
        client.Close()
    }
}
```

### Message Compression

```go
// Enable compression for large messages
transport := transport.NewWebSocketTransport("ws://localhost:8082/mcp")
transport.EnableCompression(true)
```

## Security Considerations

### Authentication

```go
// Bearer token authentication
transport.SetHeader("Authorization", "Bearer " + token)

// API key authentication  
transport.SetHeader("X-API-Key", apiKey)

// Custom authentication
transport.SetHeader("X-Auth-Token", customToken)
```

### TLS Configuration

```go
// Use secure WebSocket (WSS)
transport := transport.NewWebSocketTransport("wss://secure.example.com/mcp")

// Custom TLS config (if needed)
transport.SetTLSConfig(&tls.Config{
    InsecureSkipVerify: false,
    ServerName:         "secure.example.com",
})
```

## Production Deployment

### Health Checks

```go
func (app *App) healthCheck() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Simple ping to verify connection
    _, err := app.client.ListTools(ctx)
    return err
}
```

### Monitoring

```go
func (app *App) monitorConnection() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if err := app.healthCheck(); err != nil {
            log.Printf("Connection health check failed: %v", err)
            // Trigger reconnection
            app.reconnect()
        }
    }
}
```

### Graceful Shutdown

```go
func (app *App) shutdown() {
    // Stop accepting new requests
    app.stopping = true
    
    // Wait for ongoing operations
    app.wg.Wait()
    
    // Close client connection
    app.client.Close()
    
    // Close transport
    app.transport.Close()
}
```

## JSON Payload Examples

If you want to test the WebSocket server directly with raw JSON payloads, here are examples for each tool:

### 1. Echo Tool
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "message": "Hello from JSON payload!"
    }
  }
}
```

### 2. Text Transform Tool
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "text_transform",
    "arguments": {
      "text": "Hello WebSocket World",
      "operation": "uppercase"
    }
  }
}
```

**Available operations for text_transform:**
- `"uppercase"` - converts to uppercase
- `"lowercase"` - converts to lowercase  
- `"reverse"` - reverses the string

### 3. Long Task Tool
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "long_task",
    "arguments": {
      "duration": 5,
      "steps": 10
    }
  }
}
```

### Example Response Format
The server will respond with:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Echo at 2025-07-25T17:01:55+10:00: Hello from JSON payload!"
      }
    ]
  }
}
```

### Tool Listing Payload
To list available tools:
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/list",
  "params": {}
}
```

### Initialize Connection Payload
To initialize the MCP connection:
```json
{
  "jsonrpc": "2.0",
  "id": 0,
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-03-26",
    "capabilities": {
      "experimental": {
        "progressNotifications": true
      }
    },
    "clientInfo": {
      "name": "websocket-test-client",
      "version": "1.0.0"
    }
  }
}
```

You can send these payloads directly through a WebSocket client (like `wscat`) or use them as reference for building your own tool calls programmatically.

## Related Examples

- **[`basic_server/`](../basic_server/)** - STDIO server for simple integration
- **[`http_server/`](../http_server/)** - HTTP/SSE transport server
- **[`advanced_client/`](../advanced_client/)** - Full-featured client example
- **[`pure_library/`](../pure_library/)** - Pure library usage

This WebSocket client example provides a **comprehensive foundation** for building real-time, interactive MCP client applications with full progress tracking and robust error handling capabilities.
