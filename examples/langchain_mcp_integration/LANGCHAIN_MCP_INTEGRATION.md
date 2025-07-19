# LangChain Go + MCP Integration

## Overview

This integration combines the power of LangChain Go's agent framework with our Model Context Protocol (MCP) tools, creating a sophisticated system where LLMs can reason about tasks and execute them using MCP tools.

## Architecture

```
LangChain Go Agent → MCP Tool Wrapper → MCP Server → Tool Execution
```

### Components

1. **LangChain Go Agent**: Handles LLM reasoning and task planning
2. **MCP Tool Wrappers**: Bridge between LangChain tools interface and MCP tools
3. **MCP Server**: Manages tool registry and execution
4. **MCP Tools**: Actual tool implementations (text processing, memory, HTML creation, etc.)

## Key Features

✅ **LLM-Driven Decision Making**: LangChain agents decide which tools to use  
✅ **Rich Tool Ecosystem**: Access to all MCP tools through LangChain  
✅ **Memory Management**: Persistent memory across tool executions  
✅ **HTML Generation**: Create landing pages through agent reasoning  
✅ **Extensible**: Easy to add new MCP tools to the agent framework  

## Implementation Details

### MCP Tool Wrapper

The `MCPTool` struct implements LangChain's `tools.Tool` interface:

```go
type MCPTool struct {
    name        string
    description string
    mcpServer   *conduit.EnhancedServer
    toolName    string
}

func (t MCPTool) Call(ctx context.Context, input string) (string, error) {
    // Convert input to MCP parameters
    params := make(map[string]interface{})
    params["text"] = input
    
    // Execute via MCP tool registry
    toolRegistry := t.mcpServer.GetToolRegistry()
    result, err := toolRegistry.Call(t.toolName, params, t.mcpServer.GetMemory())
    
    return formatResult(result), err
}
```

### Available Tools in Agent

The integration provides these MCP tools to LangChain agents:

- **Text Processing**: `uppercase`, `lowercase`, `trim`, `reverse`
- **Memory Management**: `remember`, `recall`
- **HTML Creation**: `create_html_page`
- **Calculations**: LangChain's built-in `Calculator`

### HTML Tool Integration

Special handling for complex HTML creation:

```go
type MCPHTMLTool struct {
    mcpServer *conduit.EnhancedServer
    outputDir string
}

func (t MCPHTMLTool) Call(ctx context.Context, input string) (string, error) {
    // Parse "filename|content" format
    parts := strings.SplitN(input, "|", 2)
    
    params := map[string]interface{}{
        "filename": parts[0],
        "content":  parts[1],
    }
    
    // Execute HTML creation tool
    toolRegistry := t.mcpServer.GetToolRegistry()
    _, err := toolRegistry.Call("create_html_page", params, t.mcpServer.GetMemory())
    
    return "HTML page created successfully", err
}
```

## Usage Examples

### Basic Text Processing

```go
question := "Convert 'hello world' to uppercase and remember it as greeting"
answer, err := chains.Run(context.Background(), executor, question)
```

**Agent Reasoning:**
1. LLM identifies need to use `uppercase` tool
2. Converts text to "HELLO WORLD"
3. Uses `remember` tool to store as "greeting"

### Memory + Calculation

```go
question := "Remember that the answer is 42, then calculate 42 * 2"
answer, err := chains.Run(context.Background(), executor, question)
```

**Agent Reasoning:**
1. Uses `remember` tool to store "answer=42"
2. Uses `Calculator` tool to compute 42 * 2 = 84
3. Returns combined result

### HTML Page Creation

```go
question := `Create an HTML page called "demo" with Tailwind CSS content`
answer, err := chains.Run(context.Background(), executor, question)
```

**Agent Reasoning:**
1. LLM understands need to create HTML
2. Uses `create_html_page` tool with filename and content
3. Saves complete HTML file with Tailwind CSS

## Benefits of Integration

### 1. **Intelligent Tool Selection**
- LangChain agents automatically choose the right tools
- Multi-step reasoning with tool chaining
- Context-aware decision making

### 2. **Seamless Tool Execution**
- MCP tools work transparently with LangChain
- Shared memory across tool calls
- Error handling and recovery

### 3. **Extensibility**
- Easy to add new MCP tools to agent toolbox
- Custom tool wrappers for complex operations
- Flexible input/output handling

### 4. **Production Ready**
- Built on proven LangChain Go framework
- Robust MCP tool infrastructure
- Clean separation of concerns

## Running the Integration

### Prerequisites

```bash
export OPENAI_API_KEY=your_openai_api_key
```

### Execute

```bash
cd /home/engineone/Downloads/gomcp
go run examples/langchain_mcp_integration/main.go
```

### Expected Output

```
🤖 LangChain Go + MCP Integration Demo
=====================================

📝 Test 1: Text Processing
Q: Convert 'hello world' to uppercase and remember it as greeting
A: I'll convert the text to uppercase and store it for you...

🧮 Test 2: Memory + Calculation  
Q: Remember that the answer is 42, then calculate 42 * 2
A: I've stored the answer as 42 and calculated 42 * 2 = 84

🎨 Test 3: HTML Creation
Q: Create HTML demo page
A: Created HTML file: demo.html

✅ Integration demo completed!
```

## Comparison with Direct MCP Usage

| Feature | Direct MCP | LangChain + MCP |
|---------|------------|-----------------|
| Tool Execution | Manual tool calling | LLM-driven selection |
| Multi-step Tasks | Manual orchestration | Automatic chaining |
| Error Handling | Manual | Built-in recovery |
| Reasoning | Programmatic | Natural language |
| Extensibility | Add tools to registry | Wrap tools for agent |

## Future Enhancements

- **Custom Agent Types**: Specialized agents for different domains
- **Tool Composition**: Complex tools built from MCP primitives  
- **Streaming Responses**: Real-time tool execution feedback
- **Multi-Agent Systems**: Collaborative agent networks using MCP tools
- **Visual Tool Flow**: Diagram agent reasoning and tool usage

This integration represents a powerful combination of LLM reasoning capabilities with practical tool execution, enabling sophisticated automated workflows through natural language interfaces.
