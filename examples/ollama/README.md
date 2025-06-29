# Ollama Integration Example

This example demonstrates how to integrate conduit with Ollama for local LLM support.

## Prerequisites

1. **Ollama installed and running**:
   ```bash
   # Install Ollama (if not already installed)
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Start Ollama service
   ollama serve
   
   # Pull a model (in another terminal)
   ollama pull llama2
   # or
   ollama pull mistral
   # or any other model you prefer
   ```

2. **Verify Ollama is working**:
   ```bash
   curl http://localhost:11434/api/tags
   ```

## Running the Example

### Default Configuration
```bash
cd examples/ollama
go run main.go
```

### Custom Configuration
```bash
# Set custom Ollama URL
export OLLAMA_URL=http://192.168.10.10:11434

# Set custom model
export OLLAMA_MODEL=llama3.2

# Run with custom settings
go run main.go
```

### Direct Ollama Usage (No Server)
To see how to use Ollama models directly without a server:
```bash
cd examples/ollama/direct_ollama
go run main.go
```
This demonstrates direct Ollama model integration in your applications.

## Features Demonstrated

### 1. Ollama Model Integration
- Connects to local Ollama instance
- Configurable model selection
- Streaming response support
- Error handling and fallbacks

### 2. Enhanced Tools
- `model_info` - Get current Ollama configuration
- `chat_history` - Store and retrieve conversation history
- All standard text, memory, and utility tools

### 3. HTTP Endpoints
- `GET /health` - Server health check
- `GET /schema` - Available tools and schemas
- `POST /mcp` - MCP protocol with Ollama backend
- `POST /react` - ReAct agent with Ollama reasoning

## API Usage Examples

### 1. Check Server Health
```bash
curl http://localhost:8084/health
```

### 2. Get Available Tools
```bash
curl http://localhost:8084/schema
```

### 3. Chat with Ollama via MCP
```bash
curl -X POST http://localhost:8084/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama2",
    "contexts": [{
      "context_id": "chat",
      "inputs": {"query": "Hello, how are you?"}
    }]
  }'
```

### 4. Use ReAct Agent
```bash
curl -X POST http://localhost:8084/react \
  -H "Content-Type: application/json" \
  -d '{
    "thoughts": "I need to process some text and remember it",
    "model": "llama2"
  }'
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OLLAMA_URL` | `http://localhost:11434` | Ollama server URL |
| `OLLAMA_MODEL` | `llama2` | Default model to use |

## Troubleshooting

### Ollama Not Running
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# If not, start it
ollama serve
```

### Model Not Found
```bash
# List available models
ollama list

# Pull a new model
ollama pull llama2
```

### Connection Issues
```bash
# Check Ollama logs
ollama logs

# Test with different URL
export OLLAMA_URL=http://127.0.0.1:11434
```

## Integration Patterns

This example shows how to:
1. **Configure Ollama Backend**: Set custom URLs and models
2. **Enhance Context**: Add metadata to requests
3. **Memory Integration**: Store conversation history
4. **Error Handling**: Graceful fallbacks for connection issues
5. **Environment Config**: Runtime configuration via environment variables

The pattern can be extended for other LLM providers by replacing the Ollama model function with your preferred backend.
