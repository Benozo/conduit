package agents

import (
	"context"
	"time"

	"github.com/benozo/conduit/mcp"
)

// Agent represents an autonomous AI agent that can perform tasks
type Agent struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	SystemPrompt string       `json:"system_prompt"`
	Tools        []string     `json:"tools"`
	Model        string       `json:"model"`
	Memory       *mcp.Memory  `json:"-"`
	Config       *AgentConfig `json:"config"`
	State        AgentState   `json:"state"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// AgentConfig holds configuration for an agent
type AgentConfig struct {
	MaxTokens     int           `json:"max_tokens"`
	Temperature   float64       `json:"temperature"`
	TopK          int           `json:"top_k"`
	AutoRetry     bool          `json:"auto_retry"`
	MaxRetries    int           `json:"max_retries"`
	Timeout       time.Duration `json:"timeout"`
	EnableMemory  bool          `json:"enable_memory"`
	EnableLogging bool          `json:"enable_logging"`
}

// AgentState represents the current state of an agent
type AgentState string

const (
	StateIdle      AgentState = "idle"
	StateThinking  AgentState = "thinking"
	StateActing    AgentState = "acting"
	StateCompleted AgentState = "completed"
	StateError     AgentState = "error"
)

// Task represents a task that can be assigned to an agent
type Task struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output"`
	Status      TaskStatus             `json:"status"`
	Progress    float64                `json:"progress"`
	Steps       []TaskStep             `json:"steps"`
	Error       string                 `json:"error,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// TaskStep represents a single step in task execution
type TaskStep struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output"`
	ToolCalls   []mcp.ToolCall         `json:"tool_calls"`
	Status      TaskStatus             `json:"status"`
	Error       string                 `json:"error,omitempty"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// AgentManager manages multiple agents
type AgentManager struct {
	agents    map[string]*Agent
	tasks     map[string]*Task
	mcpServer interface{} // MCP server interface
	ctx       context.Context
}

// ExecutionContext provides context for agent execution
type ExecutionContext struct {
	TaskID    string
	AgentID   string
	SessionID string
	Context   context.Context
	Memory    *mcp.Memory
	Logger    Logger
}

// Logger interface for agent logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// AgentEvent represents events that occur during agent execution
type AgentEvent struct {
	Type      EventType              `json:"type"`
	AgentID   string                 `json:"agent_id"`
	TaskID    string                 `json:"task_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventType represents different types of agent events
type EventType string

const (
	EventAgentCreated  EventType = "agent_created"
	EventAgentUpdated  EventType = "agent_updated"
	EventAgentDeleted  EventType = "agent_deleted"
	EventTaskCreated   EventType = "task_created"
	EventTaskStarted   EventType = "task_started"
	EventTaskCompleted EventType = "task_completed"
	EventTaskFailed    EventType = "task_failed"
	EventToolCalled    EventType = "tool_called"
	EventStateChanged  EventType = "state_changed"
)

// EventCallback is called when agent events occur
type EventCallback func(event AgentEvent)
