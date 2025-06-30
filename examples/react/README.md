# ReAct Agent Example (WIP)

This example demonstrates the ReAct (Reasoning + Acting) pattern using both conduit's high-level API and the low-level MCP package directly.

## What is ReAct?

ReAct combines **Reasoning** and **Acting** in language models:
- **Reasoning**: The model thinks step-by-step about the problem
- **Acting**: The model uses tools to gather information or perform actions
- **Observation**: The model processes the results and continues reasoning

## Implementation Patterns

This example shows two approaches:

### 1. High-Level conduit API
Uses the `lib/conduit` package for easy server setup with built-in ReAct tools.

### 2. Direct MCP Package Usage
Uses the `mcp` package directly with `mcp.ReActAgent()` for custom ReAct workflows.

## Pattern Flow

```
Thought: I need to analyze this text for sentiment
Action: analyze_sentiment(text="I love this product")
Observation: Sentiment is positive with 70% confidence
Thought: Now I should remember this result
Action: remember(key="last_sentiment", value="positive")
Observation: Successfully stored the sentiment result
Final Answer: The text has a positive sentiment
```

## Running the Example

### HTTP Mode (Default)
```bash
cd examples/react
go run main.go
```
Server will start on http://localhost:8085

### Stdio Mode (VS Code Copilot)
```bash
cd examples/react
go run main.go --stdio
```

### Direct MCP Package Usage
To see how to use the MCP package directly without a server:
```bash
cd examples/react/direct_mcp
go run main.go
```
This demonstrates raw MCP ReAct agent functionality.

## Available Tools for ReAct Agent

### Analysis Tools
- `analyze_sentiment` - Analyze text sentiment
- `decision_maker` - Make decisions based on criteria

### Calculation Tools  
- `math_calculate` - Perform mathematical operations

### Search Tools
- `web_search_mock` - Mock web search (replace with real API)

### Standard Tools
- Text processing (uppercase, lowercase, reverse, etc.)
- Memory management (remember, recall, forget, etc.)
- Utility functions (timestamp, uuid, hashing, etc.)

## API Usage Examples

### 1. Basic ReAct Request
```bash
curl -X POST http://localhost:8085/react \
  -H "Content-Type: application/json" \
  -d '{
    "thoughts": "I need to analyze some text and remember the results"
  }'
```

### 2. Complex ReAct Workflow
```bash
curl -X POST http://localhost:8085/react \
  -H "Content-Type: application/json" \
  -d '{
    "thoughts": "Help me decide between options based on sentiment analysis",
    "context": {
      "options": ["I love this product", "This product is terrible"],
      "task": "pick the positive option"
    }
  }'
```

### 3. Tool Usage via Schema
```bash
# Get available tools
curl http://localhost:8085/schema

# Use specific tool
curl -X POST http://localhost:8085/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "tool_choice": {
      "name": "analyze_sentiment",
      "params": {"text": "I love conduit!"}
    }
  }'
```

### 4. Mathematical Reasoning
```bash
curl -X POST http://localhost:8085/react \
  -H "Content-Type: application/json" \
  -d '{
    "thoughts": "Calculate the area of a rectangle and remember it",
    "context": {
      "length": 10,
      "width": 5,
      "task": "find area"
    }
  }'
```

## Tool Examples

### Sentiment Analysis
```bash
curl -X POST http://localhost:8085/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "tool_choice": {
      "name": "analyze_sentiment", 
      "params": {"text": "This is amazing!"}
    }
  }'
```

### Math Calculation
```bash
curl -X POST http://localhost:8085/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "tool_choice": {
      "name": "math_calculate",
      "params": {
        "operation": "multiply",
        "a": 10,
        "b": 5
      }
    }
  }'
```

### Decision Making
```bash
curl -X POST http://localhost:8085/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "tool_choice": {
      "name": "decision_maker",
      "params": {
        "options": ["Option A", "Option B", "Option C"],
        "criteria": "choose the shortest one"
      }
    }
  }'
```

## Extending the Agent

### Adding Custom Tools
```go
server.RegisterTool("my_tool", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
    // Your tool implementation
    return result, nil
})
```

### Integrating Real APIs
Replace mock tools with real implementations:
- `web_search_mock` → Real web search API
- `analyze_sentiment` → ML sentiment analysis service
- Add database tools, file system tools, etc.

### Custom Models
Replace the simple model with:
- Ollama integration: `conduit.CreateOllamaModel("http://localhost:11434")`
- OpenAI API
- Other LLM providers

## ReAct Pattern Benefits

1. **Transparency**: See the reasoning process
2. **Tool Integration**: Seamlessly use external tools
3. **Error Recovery**: Can retry or adjust approach
4. **Composability**: Chain multiple tools together
5. **Memory**: Remember context across interactions

## VS Code Copilot Integration

Add to your VS Code settings for ReAct agent:
```json
{
  "mcp.mcpServers": {
    "conduit-react": {
      "command": "/path/to/examples/react/main",
      "args": ["--stdio"]
    }
  }
}
```

This gives you a ReAct-powered assistant in VS Code that can reason through problems and use tools to help with your development tasks.
