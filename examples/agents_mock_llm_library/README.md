# AI Agents with Mock LLM - Library Mode Example

This example demonstrates how to use ConduitMCP's AI agents with LLM integration in **library mode** - without starting an HTTP server. This is perfect for embedding intelligent agents directly into your Go applications.

## What This Example Shows

- 🧠 **LLM-Powered Reasoning**: Agents use a mock LLM to analyze problems and create intelligent action plans
- 📚 **Pure Library Usage**: No HTTP server required - agents work as a pure Go library
- 🔧 **Tool Integration**: Agents can use all available MCP tools (text, memory, utility, custom tools)
- ⚡ **Direct Function Calls**: Maximum performance with direct library calls
- 🎯 **Easy Integration**: Shows how to embed agents in your own applications

## Key Features Demonstrated

### 1. Mathematical Problem Solving
The agent analyzes a word problem ("Calculate the total cost of 3 items at $15 each") and:
- Uses LLM reasoning to understand it's a multiplication problem
- Plans to use the `multiply` tool
- Executes the calculation
- Stores the result in memory

### 2. Text Analysis with Context
The agent analyzes text about AI and:
- Counts words to understand text length
- Stores important content for future reference
- Generates a unique session ID for tracking

### 3. Custom Task Execution
Shows how to create and execute custom tasks with the agent manager.

## How It Works

### 1. No Server Required
```go
// Create MCP core for library mode (no server)
config := conduit.DefaultConfig()
config.EnableLogging = false // Disable server logging for library mode

server := conduit.NewEnhancedServer(config)
// Note: No server.Start() call - pure library usage!
```

### 2. Mock LLM Integration
```go
// Create a mock LLM that demonstrates intelligent reasoning
mockLLM := createMockLLM()

// Create LLM-powered agent manager (library mode)
llmAgentManager := agents.NewLLMAgentManager(server, mockLLM, "mock-llm-v1")
```

### 3. Direct Task Execution
```go
// Create task
task, _ := llmAgentManager.CreateTask(
    "intelligent_agent",
    "Calculate Total Cost",
    "Calculate total cost for multiple items",
    map[string]interface{}{
        "problem": "Calculate the total cost of 3 items at $15 each",
        "items":   3.0,
        "price":   15.0,
    },
)

// Execute with LLM reasoning
if err := llmAgentManager.ExecuteTaskWithLLM(task.ID); err != nil {
    fmt.Printf("❌ Failed: %v\\n", err)
} else {
    fmt.Println("✅ Problem solved successfully!")
}
```

## Running the Example

```bash
cd /path/to/ConduitMCP
go run examples/agents_mock_llm_library/main.go
```

## Expected Output

```
🧠 AI Agents with Mock LLM - Library Mode Demo
===============================================
This demo shows LLM-powered agents in library mode
(No HTTP server - pure library usage with mock LLM)

✓ Registered tool 'add': Add two numbers
✓ Registered tool 'multiply': Multiply two numbers
✅ Created: Intelligent Problem Solver (Library Mode)

🧪 LLM-Powered Agent Library Demonstrations
==========================================

🧮 Demo 1: Mathematical Problem Solving
Problem: Calculate the total cost of 3 items at $15 each
🧠 LLM analyzing the problem...
   🧮 Tool executed: 3.0 × 15.0 = 45.0
✅ Problem solved successfully!
  📋 Task: Calculate Total Cost
  📊 Status: completed
    Step 1: llm_reasoning
    Step 2: calculate_total
    Step 3: store_result

📝 Demo 2: Text Analysis with Context
Text: 'Artificial Intelligence is transforming technology'
🧠 LLM analyzing the text...
✅ Text analyzed successfully!
  📋 Task: Analyze Important Text
  📊 Status: completed
    Step 1: llm_reasoning
    Step 2: count_words
    Step 3: store_content
    Step 4: generate_id

⚡ Demo 3: Simple Task Execution
Creating and executing a custom calculation task
🧠 LLM planning the calculation...
   🧮 Tool executed: 7.0 × 8.0 = 56.0
✅ Calculation completed successfully!

📊 Library Mode Status
=====================
🔧 Available Tools: [add multiply word_count char_count uppercase lowercase title_case trim remember recall forget list_memories clear_memory base64_encode base64_decode hash_md5 hash_sha256 uuid timestamp random_number random_string]
🤖 Active Agents: 1
📋 Total Tasks: 3

🎓 Library Mode Benefits Demonstrated:
====================================
✅ No HTTP server required - pure Go library
✅ Direct function calls for maximum performance
✅ LLM-powered intelligent reasoning and planning
✅ Agents can be embedded in any Go application
✅ Memory and tool state managed internally
✅ Easy integration with existing Go codebases
```

## Library Mode Benefits

### 🚀 Performance
- **Direct function calls** instead of HTTP requests
- **Lower memory footprint** - no web server overhead
- **Faster execution** - no network serialization/deserialization

### 🔧 Integration
- **Easy to embed** in existing Go applications
- **No external dependencies** for server infrastructure
- **Full access** to all MCP tools and agent capabilities

### 📦 Deployment
- **Single binary** deployment
- **No port configuration** or networking setup
- **Container-friendly** - no exposed ports needed

## How to Integrate in Your App

### Basic Integration
```go
package main

import (
    "github.com/benozo/conduit/agents"
    conduit "github.com/benozo/conduit/lib"
    "github.com/benozo/conduit/lib/tools"
)

func main() {
    // 1. Create server for library mode
    config := conduit.DefaultConfig()
    config.EnableLogging = false
    server := conduit.NewEnhancedServer(config)

    // 2. Register tools
    tools.RegisterTextTools(server)
    tools.RegisterMemoryTools(server)
    tools.RegisterUtilityTools(server)

    // 3. Create LLM (replace with real LLM)
    llm := createYourLLM()

    // 4. Create agent manager
    manager := agents.NewLLMAgentManager(server, llm, "your-model")

    // 5. Create and use agents
    agent, _ := manager.CreateLLMAgent(
        "my_agent",
        "My Agent",
        "Description",
        "System prompt",
        []string{"word_count", "remember"},
        &agents.AgentConfig{
            MaxTokens:     1000,
            Temperature:   0.3,
            EnableMemory:  true,
            EnableLogging: true,
        },
    )

    // 6. Execute tasks
    task, _ := manager.CreateTask(
        "my_agent",
        "Task Title",
        "Task Description",
        map[string]interface{}{
            "input": "your input data",
        },
    )

    err := manager.ExecuteTaskWithLLM(task.ID)
    // Handle result...
}
```

### With Real Ollama LLM
```go
// Replace mock LLM with real Ollama
ollamaURL := "http://localhost:11434"
modelName := "llama3.2"
ollamaLLM := conduit.CreateOllamaModel(ollamaURL)

manager := agents.NewLLMAgentManager(server, ollamaLLM, modelName)
```

## Related Examples

- [`examples/agents_mock_llm/`](../agents_mock_llm/) - Same functionality but with HTTP server
- [`examples/agents_ollama/`](../agents_ollama/) - Real Ollama LLM integration with server
- [`examples/agents_library_mode/`](../agents_library_mode/) - Basic library mode without LLM
- [`examples/ai_agents/`](../ai_agents/) - Traditional rule-based agents

## Next Steps

1. **Replace Mock LLM**: Integrate with real LLM (Ollama, OpenAI, etc.)
2. **Add Custom Tools**: Register your own domain-specific tools
3. **Extend Agents**: Create specialized agents for your use case
4. **Scale Up**: Use multiple agents working together

This example provides a solid foundation for building intelligent, LLM-powered applications with ConduitMCP in library mode!
