# Conduit - Universal MCP Server

![Conduit Banner](ConduitBanner.png)

Conduit is a versatile, embeddable MCP (Model Context Protocol) server implementation in Go that supports both standalone and library usage. It provides dual protocol support (stdio for MCP clients and HTTP/SSE for web applications) and comes with a comprehensive set of tools for text manipulation, memory management, and utility functions.

## Universal MCP Compatibility

Conduit implements the standard MCP (Model Context Protocol) specification and works with any MCP-compatible client.

### ✅ **Tested and Verified Clients**

- **VS Code Copilot** - Full integration with all 31 tools
- **Cline** - Complete tool discovery and functionality  
- **Claude Desktop** - Standard MCP stdio support

### 🔄 **Compatible with Any MCP Client**

Since Conduit follows the MCP specification, it should work with:
- **Anthropic Claude Desktop**
- **Continue.dev**
- **Cursor IDE**
- **Any custom MCP client implementation**

Tested and confirmed with
 - **Vs Code Co-Pilot**
 - **Vs Code Cline**

All clients will have access to Conduit's complete toolkit of 31 tools for enhanced AI assistance.

## Getting Started

### 1. Install Conduit

```bash
go get github.com/benozo/conduit
```

### 2. Create Your MCP Server

```go
// main.go
package main

import (
    "log"
    
    conduit "github.com/benozo/conduit/lib"
    "github.com/benozo/conduit/lib/tools"
    "github.com/benozo/conduit/mcp"
)

func main() {
    config := conduit.DefaultConfig()
    config.Mode = mcp.ModeStdio
    
    server := conduit.NewServer(config)
    tools.RegisterTextTools(server)
    tools.RegisterMemoryTools(server)
    tools.RegisterUtilityTools(server)
    
    log.Fatal(server.Start())
}
```

### 3. Build and Configure

```bash
go build -o my-mcp-server .
```

Then add to your MCP client configuration (VS Code, Cline, etc.):
```json
{
  "command": "/path/to/my-mcp-server",
  "args": ["--stdio"]
}
```

## Features

- **Universal MCP Compatibility**: Works with any MCP client (VS Code Copilot, Cline, Claude Desktop, and more)
- **Dual Protocol Support**: stdio (for MCP clients) and HTTP/SSE (for web applications)
- **Embeddable Design**: Use as a standalone server or embed in your Go applications
- **Enhanced Tool Registration**: Rich schema support with type validation and detailed documentation
- **Modular Tool System**: Comprehensive text, memory, and utility tools (31+ tools) - fully tested and verified
- **LLM Integration**: Built-in Ollama support with streaming
- **Memory Management**: Persistent memory system for tool context
- **ReAct Agent**: Built-in reasoning and action capabilities
- **Configurable**: Flexible configuration options for different use cases

## Quick Start

### As a Library (Recommended)

Install Conduit in your Go project:

```bash
go get github.com/benozo/conduit
```

Then create your own MCP server:

```go
package main

import (
    "log"
    
    conduit "github.com/benozo/conduit/lib"
    "github.com/benozo/conduit/lib/tools"
    "github.com/benozo/conduit/mcp"
)

func main() {
    // Create configuration
    config := conduit.DefaultConfig()
    config.Port = 8080
    config.Mode = mcp.ModeStdio // For MCP clients
    
    // Create server
    server := conduit.NewServer(config)
    
    // Register tool packages
    tools.RegisterTextTools(server)
    tools.RegisterMemoryTools(server)
    tools.RegisterUtilityTools(server)
    
    // Register custom tools
    server.RegisterTool("my_tool", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
        return map[string]string{"result": "Hello from my tool!"}, nil
    })
    
    // Start server
    log.Fatal(server.Start())
}
```

#### With Enhanced Tool Registration

For tools that need rich parameter validation and documentation:

```go
package main

import (
    "log"
    conduit "github.com/benozo/conduit/lib"
    "github.com/benozo/conduit/lib/tools"
    "github.com/benozo/conduit/mcp"
)

func main() {
    config := conduit.DefaultConfig()
    config.Mode = mcp.ModeStdio
    
    // Create enhanced server
    server := conduit.NewEnhancedServer(config)
    
    // Register standard tools
    tools.RegisterTextTools(server.Server)
    tools.RegisterMemoryTools(server.Server)
    
    // Register custom tool with rich schema
    server.RegisterToolWithSchema("weather",
        func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
            city := params["city"].(string)
            return map[string]interface{}{
                "result": fmt.Sprintf("Weather in %s: Sunny, 72°F", city),
                "city": city, "temperature": "72°F", "condition": "Sunny",
            }, nil
        },
        conduit.CreateToolMetadata("weather", "Get weather for a city", map[string]interface{}{
            "city": conduit.StringParam("City name to get weather for"),
        }, []string{"city"}))
    
    log.Fatal(server.Start())
}
```

Build and use with any MCP client:

```bash
go build -o my-mcp-server .
./my-mcp-server --stdio    # For MCP clients
./my-mcp-server --http     # For HTTP API
./my-mcp-server --both     # For both protocols
```

### Standalone Usage (Development/Testing)

For development or testing, you can also clone and run directly:

```bash
git clone https://github.com/benozo/conduit
cd conduit
go run main.go --stdio    # For MCP clients (VS Code Copilot, Cline, etc.)
go run main.go --http     # For HTTP API and web applications
go run main.go --both     # For both protocols simultaneously
```

### MCP Client Configuration

**VS Code Copilot:**
```json
{
  "mcp.mcpServers": {
    "my-mcp-server": {
      "command": "/path/to/my-mcp-server",
      "args": ["--stdio"]
    }
  }
}
```

**Cline:**
```json
{
  "mcpServers": {
    "my-mcp-server": {
      "type": "stdio",
      "command": "/path/to/my-mcp-server",
      "args": ["--stdio"]
    }
  }
}
```

**Claude Desktop:**
```json
{
  "mcpServers": {
    "my-mcp-server": {
      "command": "/path/to/my-mcp-server",
      "args": ["--stdio"]
    }
  }
}
```

### Embedded Usage

Embed Conduit directly in your existing Go application:

```go
package main

import (
    "log"
    
    conduit "github.com/benozo/conduit/lib"
    "github.com/benozo/conduit/lib/tools"
    "github.com/benozo/conduit/mcp"
)

func main() {
    // Create configuration
    config := conduit.DefaultConfig()
    config.Port = 8081
    config.Mode = mcp.ModeHTTP
    
    // Create server
    server := conduit.NewServer(config)
    
    // Register tool packages
    tools.RegisterTextTools(server)
    tools.RegisterMemoryTools(server)
    tools.RegisterUtilityTools(server)
    
    // Register custom tools
    server.RegisterTool("my_tool", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
        return map[string]string{"result": "Hello from my tool!"}, nil
    })
    
    // Start server
    log.Fatal(server.Start())
}
```

### Pure Library Usage

Use MCP components directly without any server (you implement your own):

```go
package main

import "github.com/benozo/conduit/mcp"

func main() {
    // Create components
    memory := mcp.NewMemory()
    tools := mcp.NewToolRegistry()
    
    // Register tools
    tools.Register("my_tool", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
        return map[string]string{"result": "Hello!"}, nil
    })
    
    // Use directly
    result, err := tools.Call("my_tool", map[string]interface{}{}, memory)
    // Integrate into your own web server, CLI, gRPC service, etc.
}
```

## Library API

### Server Configuration

```go
config := &conduit.Config{
    Port:          8080,              // HTTP server port
    OllamaURL:     "http://localhost:11434", // Ollama API URL
    Mode:          mcp.ModeBoth,      // Server mode (Stdio/HTTP/Both)
    EnableCORS:    true,              // Enable CORS for HTTP mode
    EnableHTTPS:   false,             // Enable HTTPS
    EnableLogging: true,              // Enable logging
}
```

### Creating a Server

```go
// Standard server with default config
server := conduit.NewServer(nil)

// Standard server with custom config
server := conduit.NewServer(config)

// Enhanced server with rich schema support
server := conduit.NewEnhancedServer(config)

// Standard server with custom model
server := conduit.NewServerWithModel(config, myModelFunc)
```

### Registering Tools

```go
// Register tool packages
tools.RegisterTextTools(server)    // Text manipulation tools
tools.RegisterMemoryTools(server)  // Memory management tools
tools.RegisterUtilityTools(server) // Utility tools

// Register individual tools
server.RegisterTool("my_tool", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
    // Tool implementation
    return result, nil
})
```

### Enhanced Tool Registration (with Rich Schemas)

For tools that need rich parameter validation and documentation, use the enhanced registration system:

```go
// Create enhanced server
server := conduit.NewEnhancedServer(config)

// Register standard tools (optional)
tools.RegisterTextTools(server.Server)
tools.RegisterMemoryTools(server.Server)
tools.RegisterUtilityTools(server.Server)

// Register custom tools with full schema metadata
server.RegisterToolWithSchema("calculate",
    func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
        operation := params["operation"].(string)
        a := params["a"].(float64)
        b := params["b"].(float64)
        
        var result float64
        switch operation {
        case "add":
            result = a + b
        case "multiply":
            result = a * b
        default:
            return nil, fmt.Errorf("unknown operation: %s", operation)
        }
        
        return map[string]interface{}{"result": result}, nil
    },
    conduit.CreateToolMetadata("calculate", "Perform mathematical operations", map[string]interface{}{
        "operation": conduit.EnumParam("Mathematical operation", []string{"add", "multiply"}),
        "a":         conduit.NumberParam("First number"),
        "b":         conduit.NumberParam("Second number"),
    }, []string{"operation", "a", "b"}))

// Start with enhanced schema support
server.Start()
```

#### Schema Helper Functions

```go
// Parameter type helpers
conduit.NumberParam("Description")                           // Numbers
conduit.StringParam("Description")                           // Strings  
conduit.BoolParam("Description")                            // Booleans
conduit.ArrayParam("Description", "itemType")               // Arrays
conduit.EnumParam("Description", []string{"opt1", "opt2"})  // Enums

// Complete metadata builder
conduit.CreateToolMetadata(name, description, properties, required)
```

#### Benefits of Enhanced Registration

- **🔍 Rich Schemas**: Full JSON Schema validation with parameter types and descriptions
- **📖 Better Documentation**: MCP clients show detailed parameter information
- **✅ Type Safety**: Automatic parameter validation and error handling
- **🎯 IDE Support**: Better autocomplete and hints in MCP clients
- **🔧 Professional**: Production-ready tool definitions

### Model Integration

```go
// Use default Ollama model
server := conduit.NewServer(config) // Uses default Ollama

// Set custom model
server.SetModel(func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
    // Custom model implementation
    return response, nil
})

// Use built-in model helpers
ollamaModel := conduit.CreateOllamaModel("http://localhost:11434")
server.SetModel(ollamaModel)
```

## Available Tools

### Text Tools

- `uppercase` - Convert text to uppercase
- `lowercase` - Convert text to lowercase
- `reverse` - Reverse text
- `word_count` - Count words in text
- `trim` - Trim whitespace
- `title_case` - Convert to title case
- `snake_case` - Convert to snake_case
- `camel_case` - Convert to camelCase
- `replace` - Replace text patterns
- `extract_words` - Extract words from text
- `sort_words` - Sort words alphabetically
- `char_count` - Count characters
- `remove_whitespace` - Remove whitespace

### Memory Tools

- `remember` - Store information in memory
- `recall` - Retrieve stored information
- `forget` - Remove information from memory
- `list_memories` - List all stored memories
- `clear_memory` - Clear all memories
- `memory_stats` - Get memory statistics

### Utility Tools

- `timestamp` - Generate timestamps
- `uuid` - Generate UUIDs
- `base64_encode` - Base64 encoding
- `base64_decode` - Base64 decoding
- `url_encode` - URL encoding
- `url_decode` - URL decoding
- `hash_md5` - MD5 hashing
- `hash_sha256` - SHA256 hashing
- `json_format` - Format JSON
- `json_minify` - Minify JSON
- `random_number` - Generate random numbers
- `random_string` - Generate random strings

## HTTP API

When running in HTTP mode, the server exposes these endpoints:

- `GET /schema` - List available tools and their schemas
- `POST /mcp` - MCP protocol endpoint with Server-Sent Events (SSE)
- `POST /react` - ReAct agent endpoint for reasoning and action
- `GET /health` - Health check endpoint

Example usage:
```bash
# Get available tools
curl http://localhost:8080/schema

# Health check
curl http://localhost:8080/health
```

## Server Modes

- **ModeStdio**: Runs stdio MCP server for universal MCP client integration
- **ModeHTTP**: Runs HTTP/SSE server for web applications and custom integrations
- **ModeBoth**: Runs both protocols simultaneously for maximum compatibility

## Configuration Options

```go
type Config struct {
    Port         int                // HTTP server port (default: 8080)
    OllamaURL    string            // Ollama API URL
    Mode         mcp.ServerMode    // Server mode
    Environment  map[string]string // Environment variables
    EnableCORS   bool              // Enable CORS
    EnableHTTPS  bool              // Enable HTTPS
    CertFile     string            // HTTPS certificate file
    KeyFile      string            // HTTPS key file
    EnableLogging bool             // Enable logging
}
```

## Examples

Check the `examples/` directory for more usage examples:

- `stdio_example/` - **MCP stdio server** for client integration (VS Code, Cline, etc.)
- `sse_example/` - **HTTP/SSE server** for web applications and real-time integration
- `embedded/` - Basic embedded usage with server wrapper
- `custom_tools/` - **Enhanced tool registration** with rich schemas and validation  
- `model_integration/` - Custom model integration patterns
- `pure_library/` - Pure library usage without any server
- `pure_library_cli/` - CLI tool using MCP components
- `pure_library_web/` - Custom web server using MCP components
- `pure_mcp/` - Direct MCP usage example
- `ollama/` - Ollama integration with local LLM support
  - `direct_ollama/` - Direct Ollama model usage without server
- `react/` - ReAct agent (Reasoning + Acting) pattern
  - `direct_mcp/` - Raw MCP package ReAct usage

### Quick Start Examples

**1. MCP Stdio Server (for VS Code, Cline, etc.):**
```bash
cd examples/stdio_example && go run main.go
# Stdio MCP server for client integration
```

**2. HTTP/SSE Server (for web applications):**
```bash
cd examples/sse_example && go run main.go
# HTTP server at http://localhost:8090 with SSE support
# Visit http://localhost:8090/demo for interactive demo
```

**3. Ollama Integration:**
```bash
cd examples/ollama && go run main.go
# Server at http://localhost:8084 with Ollama backend
```

**4. ReAct Agent:**
```bash
cd examples/react && go run main.go  
# ReAct pattern server at http://localhost:8085
```

**5. Direct MCP Usage:**
```bash
cd examples/react/direct_mcp && go run main.go
# Pure MCP package demo (no server)
```

## Building

```bash
# Install as library (recommended)
go get github.com/benozo/conduit

# Build your own MCP server
go build -o my-mcp-server .

# Run tests (if developing Conduit itself)
go test ./...
```

## MCP Client Integration

After building your own MCP server with Conduit (via `go get github.com/benozo/conduit`), configure any MCP-compatible client:

### VS Code Copilot

```json
{
  "mcp.mcpServers": {
    "my-conduit-server": {
        "command": "/path/to/my-mcp-server",
        "args": ["--stdio"],
        "env": {}
    }
  }
}
```

### Cline

```json
{
  "mcpServers": {
    "my-conduit-server": {
      "autoApprove": [],
      "disabled": false,
      "timeout": 60,
      "type": "stdio",
      "command": "/path/to/my-mcp-server",
      "args": ["--stdio"]
    }
  }
}
```

### Claude Desktop

```json
{
  "mcpServers": {
    "my-conduit-server": {
      "command": "/path/to/my-mcp-server",
      "args": ["--stdio"]
    }
  }
}
```

### Other MCP Clients

For any other MCP client, use the standard MCP stdio configuration:
- **Command**: `/path/to/my-mcp-server`
- **Args**: `["--stdio"]`
- **Protocol**: stdio

**Note:** Replace `/path/to/my-mcp-server` with the actual path to your built binary.

### Verify Integration

To test that all tools are available to any MCP client:

```bash
# Test your built MCP server
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | \
  ./my-mcp-server --stdio | jq '.result.tools | length'
```

Should show 31 tools available.

#### Test Specific Tools

To verify the encoding/formatting tools work correctly:

```bash
# Test base64_decode
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "base64_decode", "arguments": {"text": "SGVsbG8gV29ybGQ="}}}' | ./my-mcp-server --stdio

# Test url_decode  
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "url_decode", "arguments": {"text": "Hello%20World%21"}}}' | ./my-mcp-server --stdio

# Test json_format
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "json_format", "arguments": {"text": "{\"name\":\"test\",\"value\":123}"}}}' | ./my-mcp-server --stdio

# Test json_minify
echo '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "json_minify", "arguments": {"text": "{\n  \"name\": \"test\",\n  \"value\": 123\n}"}}}' | ./my-mcp-server --stdio
```

All tools should return proper JSON-RPC responses with results.

### Troubleshooting

**Error: "no required module provides package"**
- Make sure you've run `go get github.com/benozo/conduit`
- Ensure your `go.mod` file includes the Conduit dependency
- Run `go mod tidy` to resolve dependencies

**Error: "Connection closed"**
- Verify your binary builds correctly: `go build -o my-mcp-server .`
- Test the binary manually: `echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./my-mcp-server --stdio`

**MCP Client Not Detecting Tools**
- Verify the client supports MCP stdio protocol
- Check client configuration points to the correct binary path
- Test the stdio interface manually (see verification steps above)
- Ensure proper timeout settings (some clients may need 30-60 seconds)

## License

MIT License - see LICENSE file for details.
