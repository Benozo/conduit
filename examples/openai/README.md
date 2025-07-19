# OpenAI MCP Server Example

This example demonstrates how to create an OpenAI-powered MCP (Model Context Protocol) server using the Conduit library.

## Features

- **OpenAI Integration**: Compatible with OpenAI API and OpenAI-compatible services
- **Full MCP Tool Support**: Includes text processing, memory management, and utility tools
- **HTTP API**: RESTful endpoints for easy integration
- **Diagnostic Tools**: Built-in tools for testing and monitoring
- **Production Ready**: Proper error handling, logging, and configuration

## Quick Start

### Prerequisites

1. Go 1.19 or later
2. OpenAI API key (or compatible service)

### Setup

1. Set your OpenAI API key:
```bash
export OPENAI_API_KEY="your-openai-api-key-here"
```

2. Optional: Configure custom settings:
```bash
export OPENAI_API_URL="https://api.openai.com"  # Default
export OPENAI_MODEL="gpt-4o-mini"               # Default
```

3. Build and run:
```bash
cd /path/to/gomcp
go build -o bin/openai examples/openai/main.go
./bin/openai
```

The server will start on port 9090 by default.

## Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/health` | Health check |
| GET    | `/schema` | Tool schema |
| POST   | `/tool`   | Execute tool |
| POST   | `/chat`   | Chat with AI |
| POST   | `/mcp`    | MCP protocol |
| POST   | `/react`  | ReAct reasoning |

## Available Tools

### Standard MCP Tools

**Text Processing:**
- `uppercase` - Convert text to uppercase
- `lowercase` - Convert text to lowercase  
- `trim` - Remove leading/trailing whitespace
- `reverse` - Reverse text

**Memory Management:**
- `remember` - Store key-value pairs
- `recall` - Retrieve stored values
- `clear_memory` - Clear all memory
- `list_memories` - List all stored keys

**Utility Tools:**
- `timestamp` - Get current timestamp
- `uuid` - Generate UUID
- `hash_md5` - Generate MD5 hash
- `hash_sha256` - Generate SHA256 hash

### Custom OpenAI Tools

- `model_info` - Get model and connection information
- `chat_history` - Manage chat history
- `openai_test` - Diagnostic tool for testing

## Usage Examples

### Test Server Health
```bash
curl http://localhost:9090/health
```

### Execute a Tool
```bash
# Convert text to uppercase
curl -X POST http://localhost:9090/tool \
  -H 'Content-Type: application/json' \
  -d '{"name":"uppercase","params":{"text":"hello world"}}'

# Get model information
curl -X POST http://localhost:9090/tool \
  -H 'Content-Type: application/json' \
  -d '{"name":"model_info","params":{}}'

# Store and retrieve data
curl -X POST http://localhost:9090/tool \
  -H 'Content-Type: application/json' \
  -d '{"name":"remember","params":{"key":"user_name","value":"Alice"}}'

curl -X POST http://localhost:9090/tool \
  -H 'Content-Type: application/json' \
  -d '{"name":"recall","params":{"key":"user_name"}}'
```

### Get Tool Schema
```bash
curl http://localhost:9090/schema
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENAI_API_KEY` | OpenAI API key (required) | - |
| `OPENAI_API_URL` | OpenAI API base URL | `https://api.openai.com` |
| `OPENAI_MODEL` | Model to use | `gpt-4o-mini` |

### OpenAI-Compatible Services

This example works with any OpenAI-compatible API. Just set the `OPENAI_API_URL`:

```bash
# For local models via text-generation-webui
export OPENAI_API_URL="http://localhost:5000"

# For other OpenAI-compatible services
export OPENAI_API_URL="https://your-service.com"
```

## Testing

Run the comprehensive test suite:

```bash
./test_openai_example.sh
```

This will test:
- Server startup and health
- All available tools
- Custom OpenAI-specific tools
- Memory persistence
- Schema endpoint

## Architecture

```
OpenAI Example
├── Server Configuration
├── Standard MCP Tools
│   ├── Text Processing
│   ├── Memory Management
│   └── Utility Functions
├── Custom OpenAI Tools
│   ├── Model Information
│   ├── Chat History
│   └── Diagnostics
└── HTTP API Endpoints
    ├── Tool Execution
    ├── Chat Interface
    └── MCP Protocol
```

## Error Handling

The server includes comprehensive error handling:

- **Invalid API Key**: Clear error messages with validation
- **Network Issues**: Proper timeout and retry handling
- **Tool Errors**: Detailed error responses
- **Invalid Requests**: JSON validation and helpful error messages

## Logging

The server provides detailed logging:

- **Startup**: Configuration and tool registration
- **Requests**: HTTP request details
- **Tool Execution**: Tool calls and results
- **Errors**: Detailed error information

## Production Deployment

For production use:

1. Use a real OpenAI API key
2. Configure appropriate environment variables
3. Set up proper logging and monitoring
4. Consider rate limiting and authentication
5. Use HTTPS in production

## Integration

This example can be integrated with:

- **Web Applications**: Use the HTTP API
- **CLI Tools**: Direct tool execution
- **Other Go Applications**: Import the conduit library
- **MCP Clients**: Full MCP protocol support

## Troubleshooting

### Common Issues

1. **"OPENAI_API_KEY is required"**
   - Set the `OPENAI_API_KEY` environment variable

2. **"Tool not found"**
   - Check tool name spelling
   - Use `curl http://localhost:9090/schema` to see available tools

3. **"Invalid request"**
   - Ensure JSON format: `{"name":"tool_name","params":{}}`
   - Check Content-Type header: `application/json`

4. **Connection refused**
   - Verify server is running on port 9090
   - Check for port conflicts

### Debug Mode

Enable verbose logging by checking the server output for detailed request/response information.

## Contributing

To add new tools or features:

1. Register tools using `server.RegisterTool()`
2. Follow the `mcp.ToolFunc` signature
3. Add proper error handling
4. Update documentation and tests

## License

This example is part of the Conduit MCP library project.
