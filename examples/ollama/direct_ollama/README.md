# Direct Ollama Usage

This example demonstrates how to use Ollama models directly through the conduit library without setting up a server.

## Key Components

### 1. Model Creation
```go
ollamaModel := conduit.CreateOllamaModel("http://localhost:11434")
```

### 2. Direct Model Calls
```go
ctx := mcp.ContextInput{
    ContextID: "test",
    Inputs: map[string]interface{}{
        "query": "Your question here",
    },
}

req := mcp.MCPRequest{
    Model: "llama3.2",
}

response, err := ollamaModel(ctx, req, memory, tokenCallback)
```

### 3. Streaming Support
```go
response, err := ollamaModel(ctx, req, memory, func(contextID string, token string) {
    fmt.Print(token) // Stream tokens as they arrive
})
```

## Prerequisites

1. Ollama running: `ollama serve`
2. Model available: `ollama pull llama3.2`

## Running the Example

```bash
# With default settings
go run main.go

# With custom configuration
OLLAMA_URL=http://192.168.1.100:11434 OLLAMA_MODEL=mistral go run main.go
```

## Use Cases

This pattern is ideal for:
- Embedding Ollama into existing applications
- Building custom AI workflows
- Testing model responses programmatically
- Creating specialized AI tools

No server overhead - just direct model integration.
