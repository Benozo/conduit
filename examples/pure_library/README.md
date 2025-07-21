# pure_library

## 🧠 What It Does

This example demonstrates how to use ConduitMCP as a pure Go library without any built-in server. It shows how to import MCP components and use them directly in your own applications, web servers, CLI tools, or gRPC services.

## ⚙️ Requirements

- **Go 1.21+** - For building and running
- **No external services** - Runs completely locally

## 🚀 How to Run

```bash
# Install dependencies (if needed)
go mod tidy

# Run the pure library demo
go run main.go
```

## 🔍 Components Used

- **Memory** — Store and retrieve key-value data
- **Tool Registry** — Register and execute custom tools
- **Processor** — Process requests with tool calling (optional)

## 💡 Sample Output

```bash
🔧 Pure Library Usage Demo
========================

✅ Memory Operations:
- Set 'user_name' = 'Alice'  
- Get 'user_name' = 'Alice'
- Memory stats: 1 items stored

✅ Tool Registry:
- Registered tool: greeting
- Tool result: {"message": "Hello Alice!", "timestamp": "2025-07-22T10:30:00Z"}

✅ Custom Tool Execution:
- Uppercase tool: "HELLO WORLD"
- Math tool: 25.0 + 15.0 = 40.0

🎯 Integration ready! Use these components in your app.
```

## 🧪 Test It

The demo shows these usage patterns:

### 1. Memory Management
```go
memory := mcp.NewMemory()
memory.Set("key", "value")
value := memory.Get("key")
```

### 2. Tool Registry
```go
tools := mcp.NewToolRegistry()
tools.Register("my_tool", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
    // Tool implementation
    return result, nil
})

// Use tool
result, err := tools.Call("my_tool", params, memory)
```

### 3. Model Processing
```go
model := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
    // Model implementation
    return response, nil
}

processor := mcp.NewProcessor(model, tools)
response, err := processor.Run(request)
```

## Integration Examples

### Web Server Integration
```go
func handler(w http.ResponseWriter, r *http.Request) {
    result, err := tools.Call("my_tool", params, memory)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(result)
}
```

### CLI Tool Integration
```go
func main() {
    result, err := tools.Call(os.Args[1], parseParams(os.Args[2:]), memory)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

### gRPC Service Integration
```go
func (s *server) ProcessTool(ctx context.Context, req *pb.ToolRequest) (*pb.ToolResponse, error) {
    result, err := s.tools.Call(req.ToolName, req.Params, s.memory)
    if err != nil {
        return nil, err
    }
    return &pb.ToolResponse{Result: result}, nil
}
```

## Benefits

- **No Server Overhead**: Use only the components you need
- **Maximum Flexibility**: Integrate into any architecture
- **Pure Library**: No forced server patterns or protocols
- **Lightweight**: Import only the MCP package
- **Custom Integration**: Build your own server/API layer

## Running the Example

```bash
cd examples/pure_library
go run main.go
```

This approach gives you complete control over how MCP components are exposed and used in your application.
