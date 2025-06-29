# Pure MCP Library Usage

This example demonstrates how to use conduit as a pure library without any built-in server. Users import the MCP package and use its components directly in their own applications and servers.

## Usage Pattern

```go
import "github.com/benozo/conduit/mcp"

// Create components
memory := mcp.NewMemory()
tools := mcp.NewToolRegistry()
processor := mcp.NewProcessor(modelFunc, tools)

// Use directly
result, err := tools.Call("tool_name", params, memory)
```

## Key Components Available

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
