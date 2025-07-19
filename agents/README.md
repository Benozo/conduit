# AI Agents Package

The AI Agents package provides a high-level abstraction layer for creating and managing autonomous AI agents that use your MCP (Model Context Protocol) server capabilities.

## Features

- **Agent Management**: Create, configure, and manage multiple AI agents
- **Task Execution**: Assign tasks to agents and track their execution
- **MCP Integration**: Direct integration with your MCP server tools
- **Specialized Agents**: Pre-built agent templates for common use cases
- **Memory Management**: Per-agent memory and context management
- **Event System**: Track agent activities and task progress

## Quick Start

```go
package main

import (
    "github.com/benozo/conduit/agents"
    conduit "github.com/benozo/conduit/lib"
    "github.com/benozo/conduit/lib/tools"
    "github.com/benozo/conduit/mcp"
)

func main() {
    // Create MCP server
    config := conduit.DefaultConfig()
    server := conduit.NewEnhancedServer(config)
    
    // Register tools
    tools.RegisterTextTools(server)
    tools.RegisterMemoryTools(server)
    tools.RegisterUtilityTools(server)
    
    // Create agent manager
    agentManager := agents.NewMCPAgentManager(server)
    
    // Create specialized agents
    agentManager.CreateSpecializedAgents()
    
    // Create and execute a task
    task, _ := agentManager.CreateTaskForAgent("math_agent", agents.TaskTypeMath, map[string]interface{}{
        "a": 10.0,
        "b": 5.0,
        "operation": "add",
    })
    
    // Start MCP server
    go server.Start()
    
    // Execute task
    agentManager.ExecuteTask(task.ID)
}
```

## Agent Types

### Pre-built Specialized Agents

1. **Math Agent** (`math_agent`)
   - Tools: `add`, `multiply`
   - Specialized in mathematical calculations
   - Low temperature for precise results

2. **Text Processing Agent** (`text_agent`)
   - Tools: `word_count`, `char_count`, `uppercase`, `lowercase`, `title_case`, `trim`
   - Specialized in text analysis and transformation

3. **Memory Management Agent** (`memory_agent`)
   - Tools: `remember`, `recall`, `forget`, `list_memories`, `clear_memory`
   - Specialized in data storage and retrieval

4. **Utility Agent** (`utility_agent`)
   - Tools: `base64_encode`, `base64_decode`, `hash_md5`, `hash_sha256`, `uuid`, `timestamp`
   - Specialized in utility functions and transformations

5. **General Purpose Agent** (`general_agent`)
   - Mixed tools from all categories
   - Versatile agent for various tasks

## Creating Custom Agents

### From Scratch

```go
agent, err := agentManager.CreateAgent(
    "custom_agent",
    "Custom Assistant",
    "A custom agent for specific tasks",
    "You are a specialized assistant...",
    []string{"tool1", "tool2", "tool3"},
    &agents.AgentConfig{
        MaxTokens:     1000,
        Temperature:   0.5,
        EnableMemory:  true,
        EnableLogging: true,
    },
)
```

### From Templates

```go
templates := agents.GetAgentTemplates()
agent, err := agentManager.CreateAgentFromTemplate(templates[0], "my_custom_agent")
```

## Task Management

### Creating Tasks

```go
// Create a specific task type
task, err := agentManager.CreateTaskForAgent("agent_id", agents.TaskTypeMath, map[string]interface{}{
    "a": 25.0,
    "b": 15.0,
    "operation": "multiply",
})

// Create a custom task
task, err := agentManager.CreateTask("agent_id", "Custom Task", "Description", map[string]interface{}{
    "custom_param": "value",
})
```

### Executing Tasks

```go
// Synchronous execution
err := agentManager.ExecuteTask(task.ID)

// Asynchronous execution
err := agentManager.ExecuteTaskAsync(task.ID)

// Wait for completion with timeout
err := agentManager.WaitForTask(task.ID, 30*time.Second)
```

### Monitoring Tasks

```go
// Get task status
task, err := agentManager.GetTask(task.ID)
fmt.Printf("Status: %s, Progress: %.1f%%\n", task.Status, task.Progress*100)

// List all tasks
tasks := agentManager.ListTasks()

// List tasks for specific agent
agentTasks := agentManager.ListTasksForAgent("agent_id")
```

## Agent Configuration

```go
config := &agents.AgentConfig{
    MaxTokens:     1000,        // Maximum tokens for responses
    Temperature:   0.7,         // Creativity level (0.0-1.0)
    TopK:          40,          // Top-K sampling
    AutoRetry:     true,        // Retry failed operations
    MaxRetries:    3,           // Maximum retry attempts
    Timeout:       30*time.Second, // Operation timeout
    EnableMemory:  true,        // Enable memory storage
    EnableLogging: true,        // Enable detailed logging
}
```

## Task Types

- `TaskTypeMath`: Mathematical calculations
- `TaskTypeTextProcessing`: Text analysis and transformation
- `TaskTypeMemoryManagement`: Data storage and retrieval
- `TaskTypeUtility`: Encoding, hashing, generation tasks
- `TaskTypeGeneral`: Multi-purpose tasks

## Agent States

- `StateIdle`: Agent is ready for tasks
- `StateThinking`: Agent is analyzing the task
- `StateActing`: Agent is executing actions
- `StateCompleted`: Agent finished the task
- `StateError`: Agent encountered an error

## Task Status

- `TaskStatusPending`: Task is waiting to be executed
- `TaskStatusRunning`: Task is currently being executed
- `TaskStatusCompleted`: Task completed successfully
- `TaskStatusFailed`: Task failed with an error
- `TaskStatusCancelled`: Task was cancelled

## Integration with MCP Tools

The agents automatically use your registered MCP tools:

```go
// Your MCP tools are available to agents
server.RegisterToolWithSchema("custom_tool", toolFunc, metadata)

// Agents can use them in their action plans
agent.Tools = []string{"custom_tool", "add", "multiply"}
```

## Advanced Usage

### Custom Action Planning

Override the action planning logic by extending the agent manager:

```go
type CustomAgentManager struct {
    *agents.MCPAgentManager
}

func (cam *CustomAgentManager) createActionPlan(execCtx *agents.ExecutionContext, task *agents.Task, agent *agents.Agent) ([]agents.Action, error) {
    // Custom planning logic
    return customActions, nil
}
```

### Event Handling

```go
// Monitor agent events (future feature)
agentManager.OnEvent(func(event agents.AgentEvent) {
    log.Printf("Event: %s for agent %s", event.Type, event.AgentID)
})
```

## Examples

See the `examples/ai_agents/` directory for complete working examples:

- Basic agent creation and task execution
- Specialized agent usage
- Custom agent development
- Task monitoring and management

## Running the Example

```bash
cd examples/ai_agents
go run main.go

# For interactive mode (keeps server running)
go run main.go --interactive
```

This will create multiple agents, execute various tasks, and demonstrate the full capabilities of the AI Agents package.
