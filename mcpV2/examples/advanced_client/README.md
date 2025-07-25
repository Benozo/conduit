# Advanced MCP Client Example

This example demonstrates how to use the Go MCP client SDK to interact with MCP servers using all available features:

- **Tool Calling**: Execute tools with complex parameters and error handling
- **Resource Access**: Read resources with different content types
- **Prompt Usage**: Use dynamic prompts with arguments
- **Error Handling**: Proper error handling and timeout management
- **Connection Management**: Robust connection setup and teardown

## Features Demonstrated

### 1. Tool Operations

The client shows how to:
- List available tools and their schemas
- Call tools with complex parameters
- Handle tool results and errors
- Process different content types

```go
// List tools
tools, err := client.ListTools(ctx)

// Call tool with parameters
params := map[string]interface{}{
    "text": "Hello World",
    "operations": []string{"uppercase", "reverse"},
}
result, err := client.CallTool(ctx, "advanced_text_transform", params)
```

### 2. Resource Management

Demonstrates resource access patterns:
- List available resources
- Read resources with different MIME types
- Handle structured data (JSON) and text content
- Process resource metadata

```go
// List resources
resources, err := client.ListResources(ctx)

// Read specific resource
response, err := client.ReadResource(ctx, "file:///config/app.json")
```

### 3. Prompt Management

Shows prompt usage with dynamic arguments:
- List available prompts and their parameters
- Generate prompts with specific arguments
- Handle prompt responses with multiple messages
- Validate required arguments

```go
// List prompts
prompts, err := client.ListPrompts(ctx)

// Get prompt with arguments
args := map[string]interface{}{
    "language": "go",
    "complexity": "high",
}
response, err := client.GetPrompt(ctx, "code_review", args)
```

### 4. Error Handling

Comprehensive error handling examples:
- Connection errors
- Tool parameter validation errors
- Missing resource errors
- Prompt argument validation errors

## Running the Example

1. **Start the advanced server** (in another terminal):
   ```bash
   cd examples/advanced_server
   go build . && ./advanced_server
   ```

2. **Build and run the client:**
   ```bash
   cd examples/advanced_client
   go build .
   ./advanced_client
   ```

## Expected Output

When run against the advanced server, the client will output:

```
✓ Connected to MCP server
✓ Server: advanced-example-server 1.0.0

=== Tool Demonstrations ===
Available tools: 2
  - advanced_text_transform: Advanced text transformation with multiple operations
  - file_operations: Perform file system operations

1. Testing Advanced Text Transform:
Result:
  Character count: 17
Original: dlrow olleh

2. Testing File Operations:
Directory listing:
  file1.txt
  dir1/
  file2.log

3. Testing Error Handling:
  Expected error: text parameter must be a string

=== Resource Demonstrations ===
Available resources: 2
  - Application Configuration (application/json): Application configuration file
  - Documentation (text/markdown): Application documentation

1. Reading Configuration Resource:
  URI: file:///config/app.json
  Type: application/json
  Content:
  {
    "app_name": "Advanced MCP Server",
    "version": "1.0.0",
    "debug": true,
    "max_connections": 100,
    "features": [
      "middleware",
      "resources", 
      "prompts",
      "tools"
    ]
  }

2. Reading Documentation Resource:
  URI: file:///docs/README.md
  Type: text/markdown
  Content (first 200 chars):
  # Advanced MCP Server

This is an advanced example of an MCP server that demonstrates:

## Features

- **Middleware Integration**: Request/response logging and metrics
- **Resource Manag...

=== Prompt Demonstrations ===
Available prompts: 2
  - code_review: Generate a code review prompt with context
    Arguments:
      - language: Programming language (required)
      - complexity: Code complexity level
  - documentation: Generate documentation prompts
    Arguments:
      - type: Documentation type (api, guide, reference) (required)
      - audience: Target audience

1. Getting Code Review Prompt:
  Description: Code review prompt for go (complexity: high)
  Messages: 1
  Message 1 (user):
    You are a senior software engineer conducting a code review for go code.

Context:
- Programming Language: go
- Code Complexity: high
- Review Focus: Best pr...

2. Getting Documentation Prompt:
  Description: api documentation prompt for developers
  Messages: 1
  Message 1 (user):
    Create comprehensive API documentation for developers.

Include:
- Endpoint descriptions
- Parameter specifications
- Response formats
- Exam...

3. Testing Error Handling with Missing Arguments:
  Expected error: language argument is required

✓ All demonstrations completed
```

## Code Structure

The example is organized into demonstration functions:

- `demonstrateTools()` - Shows tool listing and calling
- `demonstrateResources()` - Shows resource access patterns
- `demonstratePrompts()` - Shows prompt usage and argument handling
- `main()` - Orchestrates the demonstrations with proper setup/teardown

## Configuration Options

The client demonstrates various configuration options:

```go
opts := &client.ClientOptions{
    Timeout:        30 * time.Second,  // Request timeout
    ConnectTimeout: 10 * time.Second,  // Connection timeout
    ClientInfo: protocol.Implementation{
        Name:    "advanced-example-client",
        Version: "1.0.0",
    },
    Logger: &simpleLogger{},           // Custom logger
}
```

## Integration with Other Examples

This client works with:
- [`advanced_server`](../advanced_server/): Full-featured server
- [`basic_server`](../basic_server/): Basic tool functionality
- [`http_server`](../http_server/): HTTP transport testing

## Production Considerations

This example shows patterns for production use:

1. **Timeouts**: Proper timeout configuration for reliability
2. **Error Handling**: Comprehensive error checking and recovery
3. **Logging**: Structured logging for debugging and monitoring
4. **Resource Management**: Proper connection lifecycle management
5. **Type Safety**: Strong typing for parameters and responses
