# 🧠 AI Agents with LLM Integration

## Overview

This document explains how AI agents integrate with Large Language Models (LLMs) like Ollama to provide intelligent reasoning, planning, and decision-making capabilities.

## How Agents Use LLMs

### 1. **Traditional Agents vs LLM-Powered Agents**

#### Traditional Agents (Rule-Based)
```go
// Simple rule-based planning
if task.Input["operation"] == "add" {
    return []Action{{Tool: "add", Input: task.Input}}
}
```

#### LLM-Powered Agents (Intelligent Reasoning)
```go
// LLM analyzes the task and creates intelligent plans
prompt := "Analyze this task and create an execution plan..."
llmResponse := llm.Analyze(prompt)
actionPlan := parseLLMResponse(llmResponse)
```

### 2. **LLM Integration Architecture**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Task Input    │───▶│  LLM Reasoning  │───▶│ Action Planning │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Tool Execution  │◀───│ Context Memory  │───▶│  Error Recovery │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 3. **LLM Reasoning Process**

#### Step 1: Task Analysis
The LLM receives a structured prompt containing:
- Task description and requirements
- Available tools and their capabilities
- Context from previous actions
- Agent's system prompt and persona

#### Step 2: Intelligent Planning
The LLM creates a JSON response with:
```json
{
  "analysis": "Problem breakdown and reasoning",
  "steps": [
    {
      "name": "step_name",
      "description": "what this step accomplishes",
      "tool": "tool_to_use",
      "input": {"param": "value"}
    }
  ],
  "reasoning": "why this approach was chosen"
}
```

#### Step 3: Adaptive Execution
- Execute planned actions using MCP tools
- Monitor results and adapt if needed
- Use LLM for error recovery when actions fail

## Implementation Examples

### Basic LLM Agent Creation

```go
// Create Ollama model function
ollamaModel := conduit.CreateOllamaModel("http://localhost:11434")

// Create LLM-powered agent manager
llmAgentManager := agents.NewLLMAgentManager(server, ollamaModel, "llama3.2")

// Create intelligent agent
agent, err := llmAgentManager.CreateLLMAgent(
    "smart_agent",
    "Intelligent Assistant",
    "An agent that uses LLM reasoning",
    `You are an intelligent assistant that can analyze problems and use tools effectively.`,
    []string{"add", "multiply", "word_count", "remember"},
    &agents.AgentConfig{
        Temperature: 0.3,
        MaxTokens:   1000,
    },
)
```

### Task Execution with LLM

```go
// Create a complex task
task, _ := llmAgentManager.CreateTask(
    "smart_agent",
    "Calculate Project Cost",
    "Calculate total cost for a project with multiple components",
    map[string]interface{}{
        "problem": "Calculate cost: 5 developers × $100/hour × 40 hours",
        "developers": 5.0,
        "hourly_rate": 100.0,
        "hours": 40.0,
    },
)

// Execute with LLM reasoning
err := llmAgentManager.ExecuteTaskWithLLM(task.ID)
```

### LLM Reasoning Example

**Input:** "Calculate the area of a rectangle with length 15 and width 8"

**LLM Analysis:**
```json
{
  "analysis": "This is a geometric calculation. Area of rectangle = length × width. I need to multiply 15 by 8.",
  "steps": [
    {
      "name": "calculate_area",
      "description": "Multiply length by width to get area",
      "tool": "multiply",
      "input": {"a": 15.0, "b": 8.0}
    },
    {
      "name": "store_result",
      "description": "Remember this calculation for future reference",
      "tool": "remember",
      "input": {"key": "rectangle_area", "value": "15 × 8 = 120 square units"}
    }
  ],
  "reasoning": "I identified this as an area calculation and chose to use multiplication followed by storing the result."
}
```

## Key Benefits of LLM Integration

### 🧠 **Intelligent Reasoning**
- **Context Understanding**: LLMs understand the meaning and context of tasks
- **Problem Decomposition**: Break complex problems into logical steps
- **Tool Selection**: Choose the most appropriate tools for each situation

### 🎯 **Adaptive Planning**
- **Dynamic Strategy**: Plans adapt based on task complexity and requirements
- **Multi-step Coordination**: Coordinate multiple tools in logical sequences
- **Goal-oriented Execution**: Focus on achieving the desired outcome

### 🔄 **Error Recovery**
- **Intelligent Retry**: Analyze failures and suggest alternative approaches
- **Context-aware Solutions**: Use previous context to inform recovery strategies
- **Learning from Mistakes**: Adapt plans based on execution results

### 📚 **Natural Language Processing**
- **Human-like Understanding**: Process natural language task descriptions
- **Conversational Interface**: Interact with tasks in natural language
- **Explanation Generation**: Provide clear explanations of reasoning and actions

## Comparison: Traditional vs LLM Agents

| Aspect | Traditional Agents | LLM-Powered Agents |
|--------|-------------------|-------------------|
| **Planning** | Rule-based, static | Intelligent, adaptive |
| **Understanding** | Keyword matching | Natural language comprehension |
| **Error Handling** | Predefined rules | Contextual analysis and recovery |
| **Flexibility** | Limited to programmed scenarios | Handles novel situations |
| **Explanation** | Basic status updates | Detailed reasoning and justification |
| **Learning** | No adaptation | Learns from context and feedback |

## Supported LLM Integrations

### 🦙 **Ollama Integration**
```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama3.2

# Set environment variables
export OLLAMA_URL=http://localhost:11434
export OLLAMA_MODEL=llama3.2
```

### ⚙️ **Configuration Options**
```go
config := &agents.AgentConfig{
    MaxTokens:   1000,        // Maximum tokens for LLM responses
    Temperature: 0.3,         // Creativity level (0.0-1.0)
    TopK:       40,           // Token selection diversity
    EnableMemory: true,       // Enable context memory
}
```

## Real-World Use Cases

### 1. **Data Analysis Pipeline**
```go
task := map[string]interface{}{
    "description": "Analyze sales data and generate insights",
    "data_file": "sales_q3.csv",
    "required_insights": ["trends", "anomalies", "forecasts"],
}
```
**LLM Plans:** Load data → Calculate statistics → Identify patterns → Store insights

### 2. **Content Processing**
```go
task := map[string]interface{}{
    "content": "Long article about AI developments...",
    "requirements": "Extract key points and create summary",
}
```
**LLM Plans:** Count words → Extract keywords → Generate summary → Store results

### 3. **System Automation**
```go
task := map[string]interface{}{
    "operation": "Setup user account with preferences",
    "user_data": {...},
    "preferences": {...},
}
```
**LLM Plans:** Validate data → Generate IDs → Store preferences → Create session

## Error Handling and Recovery

### Intelligent Error Recovery
When a tool execution fails, the LLM can:

1. **Analyze the Error**
   ```
   Error: Invalid parameter 'x' for tool 'add'
   LLM: "The tool expects 'a' and 'b' parameters, not 'x'"
   ```

2. **Suggest Corrections**
   ```json
   {
     "corrected_action": {
       "tool": "add",
       "input": {"a": 5.0, "b": 3.0}
     },
     "reasoning": "Corrected parameter name from 'x' to 'a'"
   }
   ```

3. **Alternative Approaches**
   ```
   If primary tool fails, LLM suggests alternative tools or approaches
   ```

## Best Practices

### 1. **Prompt Design**
- **Clear Context**: Provide clear task descriptions and requirements
- **Tool Documentation**: Include comprehensive tool descriptions
- **Examples**: Show examples of successful task completions

### 2. **Temperature Settings**
- **Math/Logic Tasks**: Low temperature (0.1-0.3) for precision
- **Creative Tasks**: Higher temperature (0.5-0.8) for creativity
- **General Tasks**: Medium temperature (0.3-0.5) for balance

### 3. **Memory Management**
- **Context Preservation**: Use memory tools to maintain context
- **Session Tracking**: Store important intermediate results
- **Error Learning**: Remember successful patterns for reuse

### 4. **Error Handling**
- **Graceful Degradation**: Fallback to simpler approaches when needed
- **Retry Logic**: Implement intelligent retry with LLM guidance
- **User Feedback**: Provide clear explanations when tasks fail

## Running the Examples

### Mock LLM Demo (No external dependencies)
```bash
cd examples/agents_mock_llm
go run main.go
```

### Real Ollama Integration
```bash
# Start Ollama
ollama serve

# Pull model
ollama pull llama3.2

# Run example
cd examples/agents_ollama
export OLLAMA_URL=http://localhost:11434
export OLLAMA_MODEL=llama3.2
go run main.go
```

## Future Enhancements

- **Multi-Model Support**: Integration with multiple LLM providers
- **Learning Memory**: Persistent learning from successful patterns
- **Collaborative Agents**: Multiple agents working together with LLM coordination
- **Streaming Responses**: Real-time LLM reasoning with progress updates
- **Custom Prompt Templates**: Specialized prompts for different domains

The LLM integration transforms AI agents from simple rule-based systems into intelligent, adaptive assistants capable of human-like reasoning and problem-solving! 🚀
