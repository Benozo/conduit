# Advanced MCP Server Example

This example demonstrates the advanced features of the Go MCP SDK, including:

- **Middleware Integration**: Request/response logging and metrics collection
- **Resource Management**: File system and configuration access
- **Prompt Management**: Dynamic prompt generation with parameters
- **Advanced Tools**: Complex text processing and file operations
- **Error Handling**: Comprehensive error handling and validation

## Features Demonstrated

### 1. Middleware System

The server uses two middleware components:

- **Logging Middleware**: Logs all incoming requests and outgoing responses
- **Metrics Middleware**: Collects performance metrics for requests

```go
// Create logging middleware
loggingMiddleware := middleware.LoggingMiddleware(&simpleLogger{})

// Create metrics middleware  
metricsMiddleware := middleware.MetricsMiddleware(&simpleMetrics{})

// Add to server options
opts := &server.ServerOptions{
    Middleware: []middleware.Middleware{
        loggingMiddleware,
        metricsMiddleware,
    },
}
```

### 2. Advanced Tools

#### Text Transform Tool
Performs multiple text operations in sequence:

```bash
# Example call
echo '{"method": "tools/call", "params": {"name": "advanced_text_transform", "arguments": {"text": "hello world", "operations": ["uppercase", "reverse"]}}}' | ./advanced_server
```

Operations available:
- `uppercase` - Convert to uppercase
- `lowercase` - Convert to lowercase  
- `reverse` - Reverse the string
- `sort_words` - Sort words alphabetically
- `count_chars` - Count characters and display

#### File Operations Tool
Performs file system operations:

```bash
# List directory contents
echo '{"method": "tools/call", "params": {"name": "file_operations", "arguments": {"operation": "list", "path": "/tmp"}}}' | ./advanced_server

# Get file information
echo '{"method": "tools/call", "params": {"name": "file_operations", "arguments": {"operation": "stat", "path": "/etc/passwd"}}}' | ./advanced_server

# Create directory
echo '{"method": "tools/call", "params": {"name": "file_operations", "arguments": {"operation": "mkdir", "path": "/tmp/testdir"}}}' | ./advanced_server
```

### 3. Resource Management

#### Configuration Resource
Access application configuration:

```bash
# Read configuration
echo '{"method": "resources/read", "params": {"uri": "file:///config/app.json"}}' | ./advanced_server
```

#### Documentation Resource
Access documentation:

```bash
# Read documentation
echo '{"method": "resources/read", "params": {"uri": "file:///docs/README.md"}}' | ./advanced_server
```

### 4. Dynamic Prompts

#### Code Review Prompt
Generate code review prompts with language-specific context:

```bash
# Generate code review prompt for Python
echo '{"method": "prompts/get", "params": {"name": "code_review", "arguments": {"language": "python", "complexity": "high"}}}' | ./advanced_server
```

#### Documentation Prompt
Generate documentation prompts for different types:

```bash
# Generate API documentation prompt
echo '{"method": "prompts/get", "params": {"name": "documentation", "arguments": {"type": "api", "audience": "developers"}}}' | ./advanced_server

# Generate user guide prompt
echo '{"method": "prompts/get", "params": {"name": "documentation", "arguments": {"type": "guide", "audience": "end-users"}}}' | ./advanced_server
```

## Running the Example

1. **Build the server:**
   ```bash
   cd examples/advanced_server
   go build .
   ```

2. **Run the server:**
   ```bash
   ./advanced_server
   ```

3. **Connect using any MCP client** or test with manual JSON-RPC calls

## Architecture Overview

```
┌─────────────────┐
│   MCP Client    │
└─────────┬───────┘
          │ JSON-RPC
┌─────────▼───────┐
│   Transport     │ (STDIO/HTTP/WebSocket)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Middleware    │ (Logging, Metrics, Auth)
│     Chain       │
└─────────┬───────┘
          │
┌─────────▼───────┐
│  Request Router │
└─────┬───┬───┬───┘
      │   │   │
   ┌──▼┐ ┌▼─┐ ┌▼────┐
   │Tool│ │Resource│ │Prompt│
   │Mgr │ │ Mgr    │ │ Mgr  │
   └───┘ └────┘ └─────┘
```

## Middleware Development

To create custom middleware:

```go
func CustomMiddleware() middleware.Middleware {
    return func(next middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
            // Pre-processing
            start := time.Now()
            
            // Call next middleware/handler
            resp, err := next(ctx, req)
            
            // Post-processing
            duration := time.Since(start)
            log.Printf("Request %s took %v", req.Method, duration)
            
            return resp, err
        }
    }
}
```

## Error Handling

The example demonstrates comprehensive error handling:

- **Tool Errors**: Graceful handling of tool execution failures
- **Resource Errors**: Proper error responses for missing resources
- **Validation Errors**: Parameter validation with clear error messages
- **System Errors**: File system and other system-level error handling

## Production Considerations

This example shows patterns that can be extended for production use:

1. **Security**: Add authentication middleware
2. **Monitoring**: Extend metrics collection
3. **Scaling**: Add connection pooling and load balancing
4. **Configuration**: Use external configuration files
5. **Logging**: Integrate with structured logging systems

## Related Examples

- [`basic_server`](../basic_server/): Simple server implementation
- [`http_server`](../http_server/): HTTP transport with middleware
- [`websocket_client`](../websocket_client/): WebSocket client with progress tracking
