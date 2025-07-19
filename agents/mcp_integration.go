package agents

import (
	"fmt"
	"time"

	"github.com/benozo/conduit/mcp"
)

// MCPAgentManager extends AgentManager with direct MCP server integration
type MCPAgentManager struct {
	*AgentManager
	mcpServer interface {
		GetToolRegistry() *mcp.ToolRegistry
		GetMemory() *mcp.Memory
	}
}

// NewMCPAgentManager creates a new agent manager that integrates with MCP server
func NewMCPAgentManager(mcpServer interface {
	GetToolRegistry() *mcp.ToolRegistry
	GetMemory() *mcp.Memory
}) *MCPAgentManager {
	return &MCPAgentManager{
		AgentManager: NewAgentManager(mcpServer),
		mcpServer:    mcpServer,
	}
}

// ExecuteToolWithMCP executes a tool using the actual MCP server
func (mam *MCPAgentManager) ExecuteToolWithMCP(toolName string, params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	// Use the MCP server's tool registry to execute the tool
	return mam.mcpServer.GetToolRegistry().Call(toolName, params, memory)
}

// executeAction overrides the base implementation to use actual MCP tools
func (mam *MCPAgentManager) executeAction(execCtx *ExecutionContext, action Action, agent *Agent) (map[string]interface{}, error) {
	execCtx.Logger.Info("Executing MCP action", "action", action.Name, "tool", action.Tool)

	// Execute the tool using the MCP server
	result, err := mam.ExecuteToolWithMCP(action.Tool, action.Input, agent.Memory)
	if err != nil {
		execCtx.Logger.Error("Tool execution failed", "tool", action.Tool, "error", err)
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	// Wrap the result in a standard format
	wrappedResult := map[string]interface{}{
		"action":      action.Name,
		"tool":        action.Tool,
		"input":       action.Input,
		"output":      result,
		"executed_at": time.Now(),
		"success":     true,
	}

	execCtx.Logger.Info("Action executed successfully", "action", action.Name, "tool", action.Tool)
	return wrappedResult, nil
}

// ExecuteTask overrides the base implementation to use MCP-specific execution
func (mam *MCPAgentManager) ExecuteTask(taskID string) error {
	task, exists := mam.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	agent, err := mam.GetAgent(task.AgentID)
	if err != nil {
		return err
	}

	// Update task status
	task.Status = TaskStatusRunning
	task.Progress = 0.0
	now := time.Now()
	task.StartedAt = &now

	// Update agent state
	agent.State = StateThinking

	// Create execution context
	execCtx := &ExecutionContext{
		TaskID:    task.ID,
		AgentID:   agent.ID,
		SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
		Context:   mam.ctx,
		Memory:    agent.Memory,
		Logger:    &defaultLogger{},
	}

	// Execute the task using MCP-enabled execution
	err = mam.executeTaskStepsWithMCP(execCtx, task, agent)
	if err != nil {
		task.Status = TaskStatusFailed
		task.Error = err.Error()
		agent.State = StateError
		return err
	}

	// Mark task as completed
	task.Status = TaskStatusCompleted
	task.Progress = 1.0
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	agent.State = StateCompleted

	return nil
}

// executeTaskStepsWithMCP executes task steps using actual MCP tools
func (mam *MCPAgentManager) executeTaskStepsWithMCP(execCtx *ExecutionContext, task *Task, agent *Agent) error {
	// Create reasoning step
	reasoningStep := TaskStep{
		ID:          fmt.Sprintf("step_%d", time.Now().UnixNano()),
		Name:        "reasoning",
		Description: "Analyze the task and plan execution",
		Input:       task.Input,
		Status:      TaskStatusRunning,
		StartedAt:   timePtr(time.Now()),
	}

	task.Steps = append(task.Steps, reasoningStep)

	// Simulate reasoning process
	agent.State = StateThinking
	execCtx.Logger.Info("Agent is analyzing the task", "task_id", task.ID, "agent_id", agent.ID)

	// Create action plan based on available tools
	actionPlan, err := mam.createActionPlan(execCtx, task, agent)
	if err != nil {
		return err
	}

	reasoningStep.Output = map[string]interface{}{
		"action_plan": actionPlan,
		"reasoning":   "Analyzed task requirements and created execution plan",
	}
	reasoningStep.Status = TaskStatusCompleted
	reasoningStep.CompletedAt = timePtr(time.Now())

	// Execute action steps using actual MCP tools
	agent.State = StateActing
	for i, action := range actionPlan {
		stepID := fmt.Sprintf("step_%d_%d", time.Now().UnixNano(), i)
		actionStep := TaskStep{
			ID:          stepID,
			Name:        action.Name,
			Description: action.Description,
			Input:       action.Input,
			Status:      TaskStatusRunning,
			StartedAt:   timePtr(time.Now()),
		}

		task.Steps = append(task.Steps, actionStep)

		// Execute the action using MCP tools
		result, err := mam.executeAction(execCtx, action, agent)
		if err != nil {
			actionStep.Status = TaskStatusFailed
			actionStep.Error = err.Error()
			return err
		}

		actionStep.Output = result
		actionStep.Status = TaskStatusCompleted
		actionStep.CompletedAt = timePtr(time.Now())

		// Update task progress
		task.Progress = float64(i+1) / float64(len(actionPlan))
	}

	return nil
}

// CreateSpecializedAgents creates common types of agents
func (mam *MCPAgentManager) CreateSpecializedAgents() error {
	// Math Agent
	_, err := mam.CreateAgent(
		"math_agent",
		"Mathematical Calculator",
		"An agent specialized in mathematical calculations and operations",
		"You are a mathematical assistant. You can perform calculations using the available tools like add, multiply, and other mathematical operations. Always explain your calculations step by step.",
		[]string{"add", "multiply"},
		&AgentConfig{
			MaxTokens:     500,
			Temperature:   0.1, // Lower temperature for more precise calculations
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create math agent: %w", err)
	}

	// Text Processing Agent
	_, err = mam.CreateAgent(
		"text_agent",
		"Text Processor",
		"An agent specialized in text processing and analysis",
		"You are a text processing assistant. You can analyze text, count words and characters, transform case, and perform various text operations. Always provide detailed analysis.",
		[]string{"word_count", "char_count", "uppercase", "lowercase", "title_case", "trim"},
		&AgentConfig{
			MaxTokens:     1000,
			Temperature:   0.3,
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create text agent: %w", err)
	}

	// Memory Management Agent
	_, err = mam.CreateAgent(
		"memory_agent",
		"Memory Manager",
		"An agent specialized in memory operations and data storage",
		"You are a memory management assistant. You can store, retrieve, and manage information using memory tools. You help organize and recall important data.",
		[]string{"remember", "recall", "forget", "list_memories", "clear_memory"},
		&AgentConfig{
			MaxTokens:     800,
			Temperature:   0.2,
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create memory agent: %w", err)
	}

	// Utility Agent
	_, err = mam.CreateAgent(
		"utility_agent",
		"Utility Assistant",
		"A general-purpose agent for various utility tasks",
		"You are a utility assistant capable of performing various tasks like encoding/decoding, hashing, generating UUIDs, timestamps, and other utility functions.",
		[]string{"base64_encode", "base64_decode", "hash_md5", "hash_sha256", "uuid", "timestamp", "random_number", "random_string"},
		&AgentConfig{
			MaxTokens:     600,
			Temperature:   0.5,
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create utility agent: %w", err)
	}

	// General Purpose Agent
	_, err = mam.CreateAgent(
		"general_agent",
		"General Assistant",
		"A versatile agent capable of handling various types of tasks",
		"You are a general-purpose assistant with access to a wide range of tools. You can perform calculations, process text, manage memory, and handle utility tasks. Always choose the most appropriate tools for each task.",
		[]string{"add", "multiply", "word_count", "char_count", "remember", "recall", "uuid", "timestamp"},
		&AgentConfig{
			MaxTokens:     1200,
			Temperature:   0.7,
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create general agent: %w", err)
	}

	return nil
}

// CreateTaskForAgent creates a task optimized for a specific agent type
func (mam *MCPAgentManager) CreateTaskForAgent(agentID string, taskType TaskType, input map[string]interface{}) (*Task, error) {
	var title, description string

	switch taskType {
	case TaskTypeMath:
		title = "Mathematical Calculation"
		description = "Perform mathematical operations using available tools"
	case TaskTypeTextProcessing:
		title = "Text Processing"
		description = "Process and analyze text using text manipulation tools"
	case TaskTypeMemoryManagement:
		title = "Memory Operation"
		description = "Perform memory storage, retrieval, or management operations"
	case TaskTypeUtility:
		title = "Utility Task"
		description = "Perform utility operations like encoding, hashing, or generation"
	case TaskTypeGeneral:
		title = "General Task"
		description = "General purpose task that may require multiple types of tools"
	default:
		title = "Custom Task"
		description = "Custom task with user-defined requirements"
	}

	return mam.CreateTask(agentID, title, description, input)
}

// TaskType represents different categories of tasks
type TaskType string

const (
	TaskTypeMath             TaskType = "math"
	TaskTypeTextProcessing   TaskType = "text_processing"
	TaskTypeMemoryManagement TaskType = "memory_management"
	TaskTypeUtility          TaskType = "utility"
	TaskTypeGeneral          TaskType = "general"
)

// GetAvailableTools returns the list of tools available on the MCP server
func (mam *MCPAgentManager) GetAvailableTools() []string {
	// This would ideally query the MCP server for available tools
	// For now, return the tools we know are registered
	return []string{
		"add", "multiply",
		"word_count", "char_count", "uppercase", "lowercase", "title_case", "trim",
		"remember", "recall", "forget", "list_memories", "clear_memory",
		"base64_encode", "base64_decode", "hash_md5", "hash_sha256",
		"uuid", "timestamp", "random_number", "random_string",
	}
}

// CreateAgentFromTemplate creates an agent from a predefined template
func (mam *MCPAgentManager) CreateAgentFromTemplate(template AgentTemplate, customID string) (*Agent, error) {
	config := DefaultAgentConfig()

	// Apply template-specific configurations
	switch template.Type {
	case "math":
		config.Temperature = 0.1
		config.MaxTokens = 500
	case "text":
		config.Temperature = 0.3
		config.MaxTokens = 1000
	case "memory":
		config.Temperature = 0.2
		config.MaxTokens = 800
	case "utility":
		config.Temperature = 0.5
		config.MaxTokens = 600
	}

	return mam.CreateAgent(
		customID,
		template.Name,
		template.Description,
		template.SystemPrompt,
		template.Tools,
		config,
	)
}

// AgentTemplate represents a template for creating agents
type AgentTemplate struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	SystemPrompt string   `json:"system_prompt"`
	Tools        []string `json:"tools"`
}

// GetAgentTemplates returns predefined agent templates
func GetAgentTemplates() []AgentTemplate {
	return []AgentTemplate{
		{
			Type:         "math",
			Name:         "Math Specialist",
			Description:  "Specialized in mathematical calculations",
			SystemPrompt: "You are a mathematical specialist. Perform accurate calculations and explain your work.",
			Tools:        []string{"add", "multiply"},
		},
		{
			Type:         "text",
			Name:         "Text Processor",
			Description:  "Specialized in text processing and analysis",
			SystemPrompt: "You are a text processing expert. Analyze and transform text efficiently.",
			Tools:        []string{"word_count", "char_count", "uppercase", "lowercase", "title_case", "trim"},
		},
		{
			Type:         "memory",
			Name:         "Memory Manager",
			Description:  "Specialized in data storage and retrieval",
			SystemPrompt: "You are a memory management expert. Store and organize information effectively.",
			Tools:        []string{"remember", "recall", "forget", "list_memories", "clear_memory"},
		},
		{
			Type:         "utility",
			Name:         "Utility Helper",
			Description:  "Specialized in utility functions and transformations",
			SystemPrompt: "You are a utility expert. Perform encoding, hashing, and generation tasks.",
			Tools:        []string{"base64_encode", "base64_decode", "hash_md5", "hash_sha256", "uuid", "timestamp"},
		},
	}
}
