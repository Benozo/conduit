package mcp

type ContextInput struct {
	ContextID string                 `json:"context_id"`
	Inputs    map[string]interface{} `json:"inputs"`
}

type ToolCall struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

type MCPRequest struct {
	SessionID   string         `json:"session_id"`
	Contexts    []ContextInput `json:"contexts"`
	Model       string         `json:"model"`
	ToolChoice  *ToolCall      `json:"tool_choice,omitempty"`
	Temperature float64        `json:"temperature,omitempty"`
	TopK        int            `json:"top_k,omitempty"`
	Stream      bool           `json:"stream,omitempty"`
}