package swarm

import (
	"context"
	"time"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/mcp"
)

// Message represents a conversation message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// Agent represents a swarm agent with instructions and functions
type Agent struct {
	Name         string          `json:"name"`
	Instructions string          `json:"instructions"`
	Functions    []AgentFunction `json:"functions"`
	Model        string          `json:"model"`

	// NEW: Optional per-agent LLM configuration (backward compatible)
	ModelFunc   mcp.ModelFunc        `json:"-"`                      // Individual LLM function
	ModelConfig *conduit.ModelConfig `json:"model_config,omitempty"` // Model-specific configuration
}

// AgentFunction represents a function that an agent can call
type AgentFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Function    func(args map[string]interface{}, contextVars map[string]interface{}) Result
}

// Result represents the result of a function call
type Result struct {
	Value           string                 `json:"value"`
	Agent           *Agent                 `json:"agent,omitempty"`
	ContextVars     map[string]interface{} `json:"context_variables,omitempty"`
	ToolCallID      string                 `json:"tool_call_id,omitempty"`
	Error           error                  `json:"error,omitempty"`
	Success         bool                   `json:"success"`
	ResponseMessage *Message               `json:"response_message,omitempty"`
}

// Response represents a swarm execution response
type Response struct {
	Messages       []Message              `json:"messages"`
	Agent          *Agent                 `json:"agent"`
	ContextVars    map[string]interface{} `json:"context_variables"`
	ExecutionTime  time.Duration          `json:"execution_time"`
	TotalTurns     int                    `json:"total_turns"`
	ToolCallsCount int                    `json:"tool_calls_count"`
	HandoffsCount  int                    `json:"handoffs_count"`
	Error          error                  `json:"error,omitempty"`
	Success        bool                   `json:"success"`
}

// SwarmConfig holds configuration for the swarm client
type SwarmConfig struct {
	MaxTurns      int           `json:"max_turns"`
	ExecuteTools  bool          `json:"execute_tools"`
	Debug         bool          `json:"debug"`
	Stream        bool          `json:"stream"`
	ModelOverride string        `json:"model_override,omitempty"`
	Timeout       time.Duration `json:"timeout"`
	EnableMemory  bool          `json:"enable_memory"`
	EnableLogging bool          `json:"enable_logging"`
}

// DefaultSwarmConfig returns a default swarm configuration
func DefaultSwarmConfig() *SwarmConfig {
	return &SwarmConfig{
		MaxTurns:      10,
		ExecuteTools:  true,
		Debug:         false,
		Stream:        false,
		Timeout:       30 * time.Second,
		EnableMemory:  true,
		EnableLogging: true,
	}
}

// SwarmClient interface defines the swarm client capabilities
type SwarmClient interface {
	CreateAgent(name, instructions string, tools []string) *Agent
	CreateAgentWithModel(name, instructions string, tools []string, modelConfig *conduit.ModelConfig) *Agent
	CreateAgentWithLLM(name, instructions string, tools []string, modelFunc mcp.ModelFunc, modelName string) *Agent
	RegisterFunction(agentName string, fn AgentFunction) error
	Run(agent *Agent, messages []Message, contextVars map[string]interface{}) *Response
	RunWithContext(ctx context.Context, agent *Agent, messages []Message, contextVars map[string]interface{}) *Response
	GetAvailableTools() []string
	GetMemory() *mcp.Memory
	SetModel(modelFunc mcp.ModelFunc, modelName string)
	HasLLM() bool
}

// HandoffFunction creates a function that transfers to another agent
type HandoffFunction func() *Agent

// ContextVariables represents shared state between agents
type ContextVariables map[string]interface{}

// ExecutionContext provides context for agent execution
type ExecutionContext struct {
	SessionID     string
	CurrentAgent  *Agent
	MessageCount  int
	ToolCallCount int
	HandoffCount  int
	StartTime     time.Time
	Memory        *mcp.Memory
	Debug         bool
}

// SwarmEvent represents events that occur during swarm execution
type SwarmEvent struct {
	Type      EventType              `json:"type"`
	AgentName string                 `json:"agent_name"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventType represents different types of swarm events
type EventType string

const (
	EventAgentStart    EventType = "agent_start"
	EventAgentComplete EventType = "agent_complete"
	EventHandoff       EventType = "handoff"
	EventFunctionCall  EventType = "function_call"
	EventError         EventType = "error"
)

// Logger interface for swarm logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}
