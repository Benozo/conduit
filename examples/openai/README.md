# openai

## 🧠 What It Does

This example creates a production-ready MCP server that integrates with OpenAI's API (or OpenAI-compatible services). It demonstrates how to build a robust AI-powered tool system with comprehensive error handling, logging, and monitoring capabilities.

## ⚙️ Requirements

- **OpenAI API Key** - Get one from [OpenAI Platform](https://platform.openai.com/api-keys)
- **Go 1.21+** - For running the server
- **Internet connection** - For OpenAI API calls

## 🚀 How to Run

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="sk-your-actual-api-key-here"

# Optional: Configure custom settings
export OPENAI_API_URL="https://api.openai.com"  # Default
export OPENAI_MODEL="gpt-4o-mini"               # Default

# Run the server
go run main.go
```

## 🔍 Tools Used

**Standard MCP Tools:**
- `uppercase`, `lowercase`, `trim`, `reverse` — Text manipulation
- `remember`, `recall`, `clear_memory`, `list_memories` — Memory management  
- `timestamp`, `uuid`, `hash_md5`, `hash_sha256` — Utility functions

**Custom OpenAI Tools:**
- `model_info` — Get current model and connection status
- `chat_history` — Manage conversation history with timestamps
- `openai_test` — Test OpenAI integration with diagnostic info

## 💡 Sample Output

```bash
🧠 Using OpenAI at: https://api.openai.com
📦 Using model: gpt-4o-mini
🔧 Registering MCP tools...
✅ Registered standard MCP tools: text, memory, and utility tools
🔧 Registering custom tools...
🚀 Starting OpenAI-powered MCP server on port 9090...

📡 Available endpoints:
  GET  http://localhost:9090/health        - Health check
  GET  http://localhost:9090/schema        - Tool schema
  POST http://localhost:9090/tool          - Execute tool
  POST http://localhost:9090/chat          - Chat with AI
  POST http://localhost:9090/mcp           - MCP protocol
  POST http://localhost:9090/react         - ReAct reasoning

🔧 Environment configuration:
  OPENAI_API_KEY: sk-proj-...abc123
  OPENAI_API_URL: https://api.openai.com
  OPENAI_MODEL:   gpt-4o-mini
```

## 🧪 Test It
### 1. Health Check
```bash
curl http://localhost:9090/health
```
```json
{"status": "healthy", "timestamp": "2025-07-22T10:30:00Z"}
```

### 2. Model Information  
```bash
curl -X POST http://localhost:9090/tool \
  -H 'Content-Type: application/json' \
  -d '{"name":"model_info","params":{}}'
```
```json
{
  "openai_url": "https://api.openai.com",
  "model": "gpt-4o-mini", 
  "status": "connected",
  "provider": "OpenAI"
}
```

### 3. AI Chat with Tool Calling
```bash
curl -X POST http://localhost:9090/chat \
  -H 'Content-Type: application/json' \
  -d '{"message": "convert HELLO WORLD to lowercase and remember it"}'
```
```json
{
  "response": "I've converted 'HELLO WORLD' to lowercase as 'hello world' and stored it in memory for you.",
  "tools_used": ["lowercase", "remember"]
}
```

### 4. Chat History Management
```bash
curl -X POST http://localhost:9090/tool \
  -H 'Content-Type: application/json' \
  -d '{"name":"chat_history","params":{"message":"Testing the system"}}'
```

## 📝 Sample Prompts

Try these natural language requests:

- `"Generate a UUID and remember it as my session ID"`
- `"Convert 'Production Ready' to snake_case"`  
- `"What's the MD5 hash of 'OpenAI Integration'?"`
- `"Show me my chat history"`
- `"Clear all my memories and start fresh"`

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
