# Stdio MCP Example

This example demonstrates how to create a Conduit MCP server that runs in stdio mode for integration with MCP clients.

## What is Stdio Mode?

Stdio mode uses standard input/output for communication following the MCP (Model Context Protocol) specification. This is the standard way to integrate with MCP clients like:

- VS Code Copilot
- Cline
- Claude Desktop
- Continue.dev
- Cursor IDE
- Any MCP-compatible client

## Running the Example

```bash
# Build the example
go build -o stdio-mcp-server .

# Test manually (optional)
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./stdio-mcp-server

# Use with MCP clients (see configuration below)
```

## MCP Client Configuration

### VS Code Copilot

Add to your settings:

```json
{
  "mcp.mcpServers": {
    "conduit-stdio": {
      "command": "/path/to/examples/stdio_example/stdio-mcp-server",
      "args": []
    }
  }
}
```

### Cline

Add to your MCP settings:

```json
{
  "mcpServers": {
    "conduit-stdio": {
      "type": "stdio",
      "command": "/path/to/examples/stdio_example/stdio-mcp-server",
      "args": []
    }
  }
}
```

### Claude Desktop

Add to your configuration:

```json
{
  "mcpServers": {
    "conduit-stdio": {
      "command": "/path/to/examples/stdio_example/stdio-mcp-server",
      "args": []
    }
  }
}
```

## Available Tools

This example provides all 31 standard Conduit tools plus 2 demo tools:

### Standard Tools (31)
- Text manipulation tools (uppercase, lowercase, trim, etc.)
- Memory management tools (remember, recall, forget, etc.)
- Utility tools (base64, JSON, hashing, etc.)

### Demo Tools (2)
- `stdio_demo` - Demonstrates stdio MCP integration
- `client_info` - Shows supported MCP clients and protocol info

## Testing

```bash
# Test tool listing
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./stdio-mcp-server

# Test a specific tool
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "stdio_demo", "arguments": {}}}' | ./stdio-mcp-server

# Test text transformation
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "uppercase", "arguments": {"text": "hello world"}}}' | ./stdio-mcp-server
```

## Key Features

- **Standard MCP Protocol**: Full JSON-RPC 2.0 over stdio
- **Tool Discovery**: Dynamic tool registration and discovery
- **Universal Compatibility**: Works with any MCP client
- **Rich Tool Set**: 33 total tools available
- **Error Handling**: Proper MCP error responses
- **Logging**: Optional logging for debugging

This example shows the simplest way to create an MCP server that any MCP client can use.
