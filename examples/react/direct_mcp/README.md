# Direct MCP Package Usage

This example demonstrates how to use the `mcp` package directly for ReAct functionality without setting up a server.

## Key Components

### 1. Tool Registry
```go
toolRegistry := mcp.NewToolRegistry()
toolRegistry.Register("toolname", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
    // Tool implementation
    return result, nil
})
```

### 2. Memory Management
```go
memory := mcp.NewMemory()
memory.Set("key", value)
value := memory.Get("key")
```

### 3. ReAct Agent
```go
thoughts := []string{"action to take", "another action"}
steps, err := mcp.ReActAgent(thoughts, toolRegistry, memory)
```

## Running the Example

```bash
go run main.go
```

## Output

The example will show:
1. Tool registration and execution
2. ReAct agent step-by-step processing
3. Memory state management
4. Raw MCP package capabilities

This pattern is useful when you want to embed ReAct functionality directly into your application without the overhead of a server.
