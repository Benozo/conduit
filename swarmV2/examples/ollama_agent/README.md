# Ollama Agent Example

This example demonstrates how to create an agent that integrates with Ollama to use real LLM capabilities.

## What it demonstrates

- **Ollama Integration**: Direct connection to Ollama server
- **Real LLM Calls**: Actual AI model inference using llama3.2
- **Agent-LLM Hybrid**: Combining agent framework with language model capabilities
- **Error Handling**: Robust connection and response handling
- **Multi-prompt Testing**: Testing various types of prompts

## Prerequisites

1. **Ollama Server**: You need Ollama running at `192.168.10.10:11434`
2. **llama3.2 Model**: The model should be pulled and available

### Setting up Ollama

```bash
# Install Ollama (if not already installed)
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama server
ollama serve

# Pull the llama3.2 model
ollama pull llama3.2

# Verify the model is available
ollama list
```

## How to run

```bash
cd /home/engineone/Projects/AI/ConduitMCP/swarmV2/examples/ollama_agent
go run main.go
```

## Configuration

The example is configured to connect to:
- **Server**: `http://192.168.10.10:11434`
- **Model**: `llama3.2`

You can modify these in the `main()` function:

```go
ollamaURL := "http://192.168.10.10:11434"
model := "llama3.2"
```

## Expected output

The example will:
1. Create an Ollama-powered agent
2. Test connection to the Ollama server
3. Send multiple test prompts to the AI model
4. Display responses from llama3.2
5. Show timing and connection information

## Features demonstrated

- **Connection Testing**: Pings Ollama server before use
- **Model Information**: Displays model and provider details
- **Prompt Processing**: Sends various types of prompts
- **Response Handling**: Processes and displays AI responses
- **Error Management**: Handles connection and processing errors

## Code structure

- **OllamaProvider**: HTTP client for Ollama API (`src/llm/providers/ollama.go`)
- **OllamaAgent**: Agent wrapper with LLM capabilities
- **Interface Compliance**: Implements standard agent interfaces
- **Real AI Integration**: Actual communication with language model

## Troubleshooting

If you get connection errors:
1. Verify Ollama is running: `ps aux | grep ollama`
2. Check the server is accessible: `curl http://192.168.10.10:11434`
3. Verify the model is available: `ollama list`
4. Check firewall settings on the Ollama server

This example shows how to bridge the gap between the agent framework and real AI capabilities!
