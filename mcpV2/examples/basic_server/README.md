# Basic MCP Server - STDIO Example

This example demonstrates how to create a simple MCP server using the official Go SDK with STDIO transport. The server provides basic tools like text transformation, mathematical operations, and echo functionality.

## What is STDIO Transport?

STDIO (Standard Input/Output) transport is the most common way to integrate MCP servers with AI assistants and development tools. The server communicates via:
- **stdin** - Receives JSON-RPC requests
- **stdout** - Sends JSON-RPC responses  
- **stderr** - Logs and debug information

This makes it easy to integrate with tools like:
- **VS Code Copilot** 
- **Claude Desktop**
- **Cline**
- **Any MCP-compatible client**

## Features

- ✅ **Text Transform Tool** - uppercase, lowercase, reverse operations
- ✅ **Calculator Tool** - basic math operations (add, subtract, multiply, divide)
- ✅ **Echo Tool** - simple message echoing for testing
- ✅ **Graceful Shutdown** - proper signal handling
- ✅ **Comprehensive Logging** - debug information and error handling
- ✅ **Type Safety** - robust parameter validation

## Quick Start

### 1. Build the Server

```bash
cd examples/basic_server
go build .
```

### 2. Test Manually

You can test the server manually by sending JSON-RPC messages:

```bash
# Start the server
./basic_server

# In another terminal, send messages:
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | nc localhost 0
```

### 3. Use with MCP Clients

Add to your MCP client configuration:

#### Claude Desktop
```json
{
  "mcpServers": {
    "basic-server": {
      "command": "/path/to/basic_server",
      "args": []
    }
  }
}
```

#### VS Code Copilot
```json
{
  "mcp.servers": [
    {
      "name": "basic-server",
      "command": "/path/to/basic_server"
    }
  ]
}
```

## Available Tools

### 1. Text Transform (`text_transform`)
Transform text using various operations.

**Parameters:**
- `text` (string, required) - The text to transform
- `operation` (string, required) - Operation: "uppercase", "lowercase", "reverse"

**Example:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "text_transform",
    "arguments": {
      "text": "Hello World",
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

### 2. Calculator (`calculator`)
Perform basic mathematical operations.

**Parameters:**
- `operation` (string, required) - Operation: "add", "subtract", "multiply", "divide"
- `a` (number, required) - First number
- `b` (number, required) - Second number

**Example:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "calculator",
    "arguments": {
      "operation": "add",
      "a": 5,
      "b": 3
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
        "text": "8.00"
      }
    ]
  }
}
```

### 3. Echo (`echo`)
Echo back the provided message (useful for testing).

**Parameters:**
- `message` (string, required) - Message to echo back

**Example:**
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "message": "Hello MCP!"
    }
  }
}
```

## Testing

### Automated Testing

Use the provided test script:

```bash
./test_stdio.sh
```

### Manual Testing Sequence

1. **Initialize the connection:**
   ```bash
   echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | ./basic_server
   ```

2. **List available tools:**
   ```bash
   echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' | ./basic_server
   ```

3. **Call a tool:**
   ```bash
   echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "echo", "arguments": {"message": "Hello!"}}}' | ./basic_server
   ```

## Code Structure

```
basic_server/
├── main.go           # Main server implementation
├── test_stdio.sh     # Testing script
└── README.md         # This file
```

### Key Components

- **Server Setup** - Creates MCP server with STDIO transport
- **Tool Registration** - Registers tools with schema validation
- **Request Handlers** - Implements tool logic with error handling
- **Signal Handling** - Graceful shutdown on SIGINT/SIGTERM
- **Logging** - Comprehensive debug and error logging

## Error Handling

The server includes robust error handling:

- **Parameter Validation** - Type checking and required parameter validation
- **Division by Zero** - Proper error for calculator division by zero
- **Unknown Operations** - Clear error messages for invalid operations
- **Signal Handling** - Graceful shutdown on interruption

## Extending the Server

To add new tools:

1. **Define the tool schema:**
   ```go
   tool := &protocol.Tool{
       Name: "my_tool",
       Description: "My custom tool",
       InputSchema: protocol.JSONSchema{
           // Schema definition
       },
   }
   ```

2. **Implement the handler:**
   ```go
   func handleMyTool(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
       // Tool implementation
       return &protocol.CallToolResult{
           Content: []protocol.Content{{
               Type: "text",
               Text: "result",
           }},
       }, nil
   }
   ```

3. **Register the tool:**
   ```go
   srv.RegisterTool(tool, handleMyTool)
   ```

## Related Examples

- **[`advanced_server/`](../advanced_server/)** - Full-featured server with middleware, resources, and prompts
- **[`advanced_client/`](../advanced_client/)** - Client example showing how to connect to servers
- **[`pure_library/`](../pure_library/)** - Pure library usage without networking
- **[`http_server/`](../http_server/)** - HTTP transport example

## Production Considerations

For production use, consider:

1. **Logging** - Use structured logging (logrus, zap)
2. **Configuration** - External configuration files
3. **Monitoring** - Health checks and metrics
4. **Security** - Input validation and sanitization
5. **Performance** - Connection pooling and optimization

This basic server provides a solid foundation for building MCP-enabled applications with STDIO transport.
