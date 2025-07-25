# Pure Library Example

This example demonstrates how to use MCP components as a **pure library** without any client-server infrastructure. This approach provides maximum performance and flexibility for embedding MCP functionality directly into applications.

## What is Pure Library Mode?

Pure library mode enables you to use MCP components (tools, memory, data processing) as **native Go library calls** without:
- ❌ JSON-RPC overhead
- ❌ Network communication  
- ❌ Transport layers
- ❌ Client-server setup

Instead, you get:
- ✅ **Direct function calls** - Zero serialization overhead
- ✅ **Maximum performance** - No network or protocol delays
- ✅ **Type safety** - Native Go interfaces and types
- ✅ **Easy integration** - Drop into any Go application
- ✅ **In-process memory** - No external dependencies

This makes it perfect for:
- **Embedded applications**
- **High-performance services**
- **Desktop applications**
- **CLI tools**
- **Microservices** requiring tool functionality

## Features

- ✅ **Tool Registry** - Register and call tools directly
- ✅ **Memory Backend** - In-memory key-value storage with stats
- ✅ **Data Processing** - Transform and validate data
- ✅ **Error Handling** - Native Go error patterns
- ✅ **Zero Dependencies** - No external services required
- ✅ **Type Safety** - Full Go type system support

## Quick Start

### 1. Build and Run

```bash
cd examples/pure_library
go build .
./pure_library
```

### 2. Run Integration Example

```bash
go run -tags=integration integration_example.go
```

### 3. Run Benchmarks

```bash
go test -bench=. -benchmem
```

Expected benchmark results:
```
BenchmarkToolCall-8           	 7516413	       156.9 ns/op	     128 B/op	       2 allocs/op
BenchmarkMemoryOperations/Set-8         	  402762	      2833 ns/op	     407 B/op	       3 allocs/op
BenchmarkMemoryOperations/Get-8         	42615894	        28.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkTextTransform-8                	 1959620	       543.8 ns/op	     192 B/op	       3 allocs/op
BenchmarkCalculator-8                   	 2408766	       534.6 ns/op	     144 B/op	       4 allocs/op
```

```
=== Tool Demonstration ===
Available tools: [calculator greet text_transform]

1. Text Transformation:
   Input: 'Hello, World!' -> Output: 'HELLO, WORLD!'

2. Calculator:
   15.5 * 2.0 = 31.00

3. Greeting:
   Good morning, Alice!

4. Error Handling:
   Expected error: division by zero

=== Memory Demonstration ===
Retrieved user:1: map[age:30 email:alice@example.com name:Alice]
All keys: [user:1 user:2 config:theme]
Memory stats: 3 active keys, backend: inmemory
Keys after deletion: [user:1 user:2]
Pure library example completed!
```

## Project Structure

```
pure_library/
├── main.go                    # Main example demonstrating library usage
├── integration_example.go     # Integration example for applications  
├── benchmark_test.go          # Performance benchmarks
├── test_library.sh           # Test script
└── README.md                 # This documentation
```

## Core Components

### 1. Component Registry

The `ComponentRegistry` is the entry point for pure library usage:

```go
// Create registry
registry := library.NewComponentRegistry()

// Access components
tools := registry.Tools()
memory := registry.Memory()
processor := registry.Processor()
```

### 2. Tool Registry

Register and call tools without any protocol overhead:

```go
// Register a tool
tools.Register("my_tool", func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
    // Tool implementation
    return &protocol.ToolResult{
        Content: []protocol.Content{{
            Type: "text",
            Text: "Tool result",
        }},
    }, nil
})

// Call the tool
result, err := tools.Call(ctx, "my_tool", map[string]interface{}{
    "param1": "value1",
})
```

### 3. Memory Backend

High-performance in-memory storage:

```go
// Store data
memory.Set("key", map[string]interface{}{
    "name": "Alice",
    "age":  30,
})

// Retrieve data
value, err := memory.Get("key")

// List all keys
keys, err := memory.List()

// Get statistics
stats, err := memory.Stats()
```

## Available Tools

### 1. Text Transform (`text_transform`)
Transform text using various operations.

**Parameters:**
- `text` (string, required) - The text to transform
- `operation` (string, required) - Operation: "uppercase", "lowercase", "reverse", "title"

**Example:**
```go
result, err := tools.Call(ctx, "text_transform", map[string]interface{}{
    "text":      "hello world",
    "operation": "uppercase",
})
// Result: "HELLO WORLD"
```

### 2. Calculator (`calculator`)
Perform mathematical operations.

**Parameters:**
- `operation` (string, required) - Operation: "add", "subtract", "multiply", "divide"
- `a` (number, required) - First number
- `b` (number, required) - Second number

**Example:**
```go
result, err := tools.Call(ctx, "calculator", map[string]interface{}{
    "operation": "multiply",
    "a":         15.5,
    "b":         2.0,
})
// Result: "31.00"
```

### 3. Greeting (`greet`)
Generate personalized greetings.

**Parameters:**
- `name` (string, required) - Person's name
- `time` (string, optional) - Time of day: "morning", "afternoon", "evening"

**Example:**
```go
result, err := tools.Call(ctx, "greet", map[string]interface{}{
    "name": "Alice",
    "time": "morning",
})
// Result: "Good morning, Alice!"
```

## Integration Patterns

### 1. Embedded Application

```go
package main

import (
    "context"
    "github.com/modelcontextprotocol/go-sdk/library"
)

type App struct {
    registry *library.ComponentRegistry
}

func NewApp() *App {
    return &App{
        registry: library.NewComponentRegistry(),
    }
}

func (a *App) ProcessText(text, operation string) (string, error) {
    result, err := a.registry.Tools().Call(context.Background(), "text_transform", map[string]interface{}{
        "text":      text,
        "operation": operation,
    })
    if err != nil {
        return "", err
    }
    return result.Content[0].Text, nil
}
```

### 2. Custom Tool Registration

```go
// Register your own tools
err := registry.Tools().Register("custom_tool", func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
    // Your tool logic here
    return &protocol.ToolResult{
        Content: []protocol.Content{{
            Type: "text",
            Text: "Custom result",
        }},
    }, nil
})
```

### 3. Memory with Custom Backend

```go
// Use custom memory backend (when available)
customMemory := NewCustomMemoryBackend()
registry.SetMemory(customMemory)
```

## Performance Characteristics

### Benchmarks
Pure library mode provides significant performance advantages:

- **Direct calls**: ~100ns per tool call
- **No serialization**: Zero JSON marshaling overhead
- **No network**: No TCP/HTTP stack involvement
- **Memory locality**: All data in process memory

### Memory Usage
- **Minimal overhead**: Only Go structs and interfaces
- **Configurable**: Memory backend can be tuned
- **Garbage collection**: Standard Go GC behavior

## API Reference

### ComponentRegistry

```go
type ComponentRegistry struct{}

// Create new registry
func NewComponentRegistry() *ComponentRegistry

// Access components
func (cr *ComponentRegistry) Tools() ToolRegistry
func (cr *ComponentRegistry) Memory() Memory
func (cr *ComponentRegistry) Processor() *Processor
func (cr *ComponentRegistry) SetMemory(memory Memory)
```

### ToolRegistry

```go
type ToolRegistry interface {
    Register(name string, handler ToolFunc) error
    RegisterWithSchema(name string, handler ToolFunc, schema *protocol.JSONSchema) error
    Call(ctx context.Context, name string, params map[string]interface{}) (*protocol.ToolResult, error)
    List() []string
    Remove(name string)
    Clear()
    Stats() *ToolStats
}
```

### Memory

```go
type Memory interface {
    Set(key string, value interface{}) error
    Get(key string) (interface{}, error)
    Delete(key string) error
    List() ([]string, error)
    Clear() error
    Stats() (*MemoryStats, error)
    Close() error
}
```

## Error Handling

Pure library mode uses standard Go error patterns:

```go
// Tool call errors
result, err := tools.Call(ctx, "tool_name", params)
if err != nil {
    // Handle tool execution error
    log.Printf("Tool error: %v", err)
    return
}

// Memory errors
value, err := memory.Get("key")
if err == library.ErrKeyNotFound {
    // Handle missing key
} else if err != nil {
    // Handle other errors
}
```

## Testing

### Unit Testing

```go
func TestToolIntegration(t *testing.T) {
    registry := library.NewComponentRegistry()
    
    // Register test tool
    registry.Tools().Register("test_tool", func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
        return &protocol.ToolResult{
            Content: []protocol.Content{{Type: "text", Text: "test"}},
        }, nil
    })
    
    // Test tool call
    result, err := registry.Tools().Call(context.Background(), "test_tool", nil)
    assert.NoError(t, err)
    assert.Equal(t, "test", result.Content[0].Text)
}
```

### Benchmark Testing

```go
func BenchmarkToolCall(b *testing.B) {
    registry := library.NewComponentRegistry()
    registry.Tools().Register("bench_tool", simpleTool)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        registry.Tools().Call(context.Background(), "bench_tool", nil)
    }
}
```

## Advanced Usage

### 1. Custom Tool Validation

```go
// Register tool with schema validation
schema := &protocol.JSONSchema{
    Type: "object",
    Properties: map[string]*protocol.JSONSchema{
        "input": {Type: "string", Description: "Input text"},
    },
    Required: []string{"input"},
}

err := tools.RegisterWithSchema("validated_tool", handler, schema)
```

### 2. Memory Statistics Monitoring

```go
// Monitor memory usage
go func() {
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        stats, err := memory.Stats()
        if err == nil {
            log.Printf("Memory: %d keys, %d bytes", stats.ActiveKeys, stats.MemoryUsage)
        }
    }
}()
```

### 3. Tool Call Metrics

```go
// Wrap tool calls with metrics
originalCall := tools.Call
tools.Call = func(ctx context.Context, name string, params map[string]interface{}) (*protocol.ToolResult, error) {
    start := time.Now()
    result, err := originalCall(ctx, name, params)
    duration := time.Since(start)
    
    // Record metrics
    metrics.RecordToolCall(name, duration, err == nil)
    
    return result, err
}
```

## Production Considerations

### 1. Memory Management

```go
// Implement memory cleanup
func (app *App) Cleanup() {
    app.registry.Memory().Clear()
    app.registry.Tools().Clear()
}
```

### 2. Error Recovery

```go
// Wrap tool calls for error recovery
func SafeToolCall(tools ToolRegistry, name string, params map[string]interface{}) (result *protocol.ToolResult, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("tool panic: %v", r)
        }
    }()
    
    return tools.Call(context.Background(), name, params)
}
```

### 3. Resource Limits

```go
// Add timeouts and limits
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := tools.Call(ctx, "tool_name", params)
```

## Migration from Client-Server

Converting from client-server to pure library mode:

### Before (Client-Server)
```go
// Create client
client := mcp.NewClient(transport)

// Call tool via JSON-RPC
response, err := client.CallTool(ctx, "tool_name", params)
```

### After (Pure Library)
```go
// Create registry
registry := library.NewComponentRegistry()

// Call tool directly
result, err := registry.Tools().Call(ctx, "tool_name", params)
```

## Related Examples

- **[`basic_server/`](../basic_server/)** - STDIO server example
- **[`advanced_server/`](../advanced_server/)** - Full-featured server
- **[`http_server/`](../http_server/)** - HTTP/SSE transport
- **[`websocket_client/`](../websocket_client/)** - WebSocket client

This pure library example provides the **highest performance** approach to using MCP functionality, perfect for applications that need tool and memory capabilities without transport overhead.
