package agents

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/benozo/conduit/mcp"
)

// NewAgentManager creates a new agent manager
func NewAgentManager(mcpServer interface{}) *AgentManager {
	return &AgentManager{
		agents:    make(map[string]*Agent),
		tasks:     make(map[string]*Task),
		mcpServer: mcpServer,
		ctx:       context.Background(),
	}
}

// CreateAgent creates a new agent with the given configuration
func (am *AgentManager) CreateAgent(id, name, description, systemPrompt string, tools []string, config *AgentConfig) (*Agent, error) {
	if config == nil {
		config = DefaultAgentConfig()
	}

	agent := &Agent{
		ID:           id,
		Name:         name,
		Description:  description,
		SystemPrompt: systemPrompt,
		Tools:        tools,
		Config:       config,
		State:        StateIdle,
		Memory:       mcp.NewMemory(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	am.agents[id] = agent
	return agent, nil
}

// GetAgent retrieves an agent by ID
func (am *AgentManager) GetAgent(id string) (*Agent, error) {
	agent, exists := am.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}
	return agent, nil
}

// ListAgents returns all agents
func (am *AgentManager) ListAgents() []*Agent {
	agents := make([]*Agent, 0, len(am.agents))
	for _, agent := range am.agents {
		agents = append(agents, agent)
	}
	return agents
}

// DeleteAgent removes an agent
func (am *AgentManager) DeleteAgent(id string) error {
	delete(am.agents, id)
	return nil
}

// CreateTask creates a new task for an agent
func (am *AgentManager) CreateTask(agentID, title, description string, input map[string]interface{}) (*Task, error) {
	agent, err := am.GetAgent(agentID)
	if err != nil {
		return nil, err
	}

	taskID := fmt.Sprintf("task_%d", time.Now().UnixNano())
	task := &Task{
		ID:          taskID,
		AgentID:     agent.ID,
		Title:       title,
		Description: description,
		Input:       input,
		Status:      TaskStatusPending,
		Progress:    0.0,
		Steps:       []TaskStep{},
		CreatedAt:   time.Now(),
	}

	am.tasks[taskID] = task
	return task, nil
}

// ExecuteTask executes a task using the assigned agent
func (am *AgentManager) ExecuteTask(taskID string) error {
	task, exists := am.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	agent, err := am.GetAgent(task.AgentID)
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
		Context:   am.ctx,
		Memory:    agent.Memory,
		Logger:    &defaultLogger{},
	}

	// Execute the task
	err = am.executeTaskSteps(execCtx, task, agent)
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

// executeTaskSteps executes the individual steps of a task
func (am *AgentManager) executeTaskSteps(execCtx *ExecutionContext, task *Task, agent *Agent) error {
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
	actionPlan, err := am.createActionPlan(execCtx, task, agent)
	if err != nil {
		return err
	}

	reasoningStep.Output = map[string]interface{}{
		"action_plan": actionPlan,
		"reasoning":   "Analyzed task requirements and created execution plan",
	}
	reasoningStep.Status = TaskStatusCompleted
	reasoningStep.CompletedAt = timePtr(time.Now())

	// Execute action steps
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

		// Execute the action
		result, err := am.executeAction(execCtx, action, agent)
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

// Action represents an action that an agent can take
type Action struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Tool        string                 `json:"tool"`
	Input       map[string]interface{} `json:"input"`
}

// createActionPlan creates a plan of actions based on the task and available tools
func (am *AgentManager) createActionPlan(execCtx *ExecutionContext, task *Task, agent *Agent) ([]Action, error) {
	// Simple action planning based on task input
	var actions []Action

	// Check if task requires text processing
	if query, ok := task.Input["query"].(string); ok {
		// Add text analysis action
		actions = append(actions, Action{
			Name:        "analyze_text",
			Description: "Analyze the input query",
			Tool:        "word_count",
			Input: map[string]interface{}{
				"text": query,
			},
		})

		// Add memory storage action
		actions = append(actions, Action{
			Name:        "store_context",
			Description: "Store the query context in memory",
			Tool:        "remember",
			Input: map[string]interface{}{
				"key":   fmt.Sprintf("task_%s_query", task.ID),
				"value": query,
			},
		})
	}

	// Check if task requires calculation
	if a, aOk := task.Input["a"].(float64); aOk {
		if b, bOk := task.Input["b"].(float64); bOk {
			operation := "add"
			if op, opOk := task.Input["operation"].(string); opOk {
				operation = op
			}

			actions = append(actions, Action{
				Name:        "perform_calculation",
				Description: fmt.Sprintf("Perform %s operation", operation),
				Tool:        operation,
				Input: map[string]interface{}{
					"a": a,
					"b": b,
				},
			})
		}
	}

	// If no specific actions found, add a default analysis action
	if len(actions) == 0 {
		actions = append(actions, Action{
			Name:        "analyze_task",
			Description: "Analyze the task input",
			Tool:        "timestamp",
			Input:       map[string]interface{}{},
		})
	}

	return actions, nil
}

// executeAction executes a single action using the MCP server
func (am *AgentManager) executeAction(execCtx *ExecutionContext, action Action, agent *Agent) (map[string]interface{}, error) {
	execCtx.Logger.Info("Executing action", "action", action.Name, "tool", action.Tool)

	// TODO: Integrate with actual MCP server tool execution
	// For now, simulate tool execution
	result := map[string]interface{}{
		"action":      action.Name,
		"tool":        action.Tool,
		"input":       action.Input,
		"executed_at": time.Now(),
		"success":     true,
	}

	// Simulate different tool responses
	switch action.Tool {
	case "word_count":
		if text, ok := action.Input["text"].(string); ok {
			result["word_count"] = len([]rune(text))
		}
	case "add":
		if a, aOk := action.Input["a"].(float64); aOk {
			if b, bOk := action.Input["b"].(float64); bOk {
				result["result"] = a + b
			}
		}
	case "multiply":
		if a, aOk := action.Input["a"].(float64); aOk {
			if b, bOk := action.Input["b"].(float64); bOk {
				result["result"] = a * b
			}
		}
	case "timestamp":
		result["timestamp"] = time.Now().Unix()
	}

	return result, nil
}

// GetTask retrieves a task by ID
func (am *AgentManager) GetTask(taskID string) (*Task, error) {
	task, exists := am.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}
	return task, nil
}

// ListTasks returns all tasks
func (am *AgentManager) ListTasks() []*Task {
	tasks := make([]*Task, 0, len(am.tasks))
	for _, task := range am.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// ListTasksForAgent returns all tasks for a specific agent
func (am *AgentManager) ListTasksForAgent(agentID string) []*Task {
	var tasks []*Task
	for _, task := range am.tasks {
		if task.AgentID == agentID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// DefaultAgentConfig returns default configuration for agents
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		MaxTokens:     1000,
		Temperature:   0.7,
		TopK:          40,
		AutoRetry:     true,
		MaxRetries:    3,
		Timeout:       30 * time.Second,
		EnableMemory:  true,
		EnableLogging: true,
	}
}

// defaultLogger provides a simple logger implementation
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *defaultLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

// timePtr returns a pointer to a time value
func timePtr(t time.Time) *time.Time {
	return &t
}

// ExecuteTaskAsync executes a task asynchronously
func (am *AgentManager) ExecuteTaskAsync(taskID string) error {
	go func() {
		if err := am.ExecuteTask(taskID); err != nil {
			log.Printf("Task execution failed: %v", err)
		}
	}()
	return nil
}

// CancelTask cancels a running task
func (am *AgentManager) CancelTask(taskID string) error {
	task, exists := am.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	if task.Status == TaskStatusRunning {
		task.Status = TaskStatusCancelled
		return nil
	}

	return fmt.Errorf("task %s is not running", taskID)
}

// WaitForTask waits for a task to complete
func (am *AgentManager) WaitForTask(taskID string, timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		task, err := am.GetTask(taskID)
		if err != nil {
			return err
		}

		if task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed || task.Status == TaskStatusCancelled {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for task %s", taskID)
}
