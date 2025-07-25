// Package protocol implements the core JSON-RPC 2.0 and MCP protocol types.
//
// This package provides the fundamental message types, error codes, and
// validation logic for the Model Context Protocol (MCP) implementation.
// It serves as the foundation for both client and server implementations.
package protocol

import (
	"fmt"
	"time"
)

// JSONRPCVersion is the supported JSON-RPC version
const JSONRPCVersion = "2.0"

// MCPProtocolVersion is the supported MCP protocol version
const MCPProtocolVersion = "2025-03-26"

// JSONRPCMessage represents a JSON-RPC 2.0 message
type JSONRPCMessage struct {
	Version string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"` // string, number, or null
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	Meta    *Meta       `json:"_meta,omitempty"` // MCP extension
}

// Meta provides context for progress tracking and tracing
type Meta struct {
	ProgressToken string                 `json:"progressToken,omitempty"`
	TraceID       string                 `json:"traceId,omitempty"`
	SpanID        string                 `json:"spanId,omitempty"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}

// RPCError represents a JSON-RPC error object
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error implements the error interface
func (e *RPCError) Error() string {
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// Standard JSON-RPC error codes
const (
	// JSON-RPC standard errors
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603

	// MCP-specific error codes
	NotInitialized       = -32001
	RequestFailed        = -32002
	InvalidTool          = -32003
	InvalidResource      = -32004
	MethodDisabled       = -32005
	InvalidPrompt        = -32006
	AuthenticationFailed = -32007
	PermissionDenied     = -32008
	RateLimitExceeded    = -32009
)

// Implementation represents client or server implementation information
type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ClientCapabilities defines what a client can do
type ClientCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Sampling     *SamplingCapability    `json:"sampling,omitempty"`
	Roots        *RootsCapability       `json:"roots,omitempty"`
}

// ServerCapabilities defines what a server supports
type ServerCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Logging      *LoggingCapability     `json:"logging,omitempty"`
	Prompts      *PromptsCapability     `json:"prompts,omitempty"`
	Resources    *ResourcesCapability   `json:"resources,omitempty"`
	Tools        *ToolsCapability       `json:"tools,omitempty"`
}

// SamplingCapability indicates sampling support
type SamplingCapability struct{}

// RootsCapability indicates roots support
type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// LoggingCapability indicates logging support
type LoggingCapability struct{}

// PromptsCapability indicates prompts support
type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCapability indicates resources support
type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// ToolsCapability indicates tools support
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// InitializeRequest represents the MCP initialization request
type InitializeRequest struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      Implementation     `json:"clientInfo"`
	Meta            *Meta              `json:"_meta,omitempty"`
}

// InitializeResult represents the MCP initialization response
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      Implementation     `json:"serverInfo"`
	Instructions    string             `json:"instructions,omitempty"`
}

// JSONSchema represents a JSON Schema definition
type JSONSchema struct {
	Type                 string                 `json:"type,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Properties           map[string]*JSONSchema `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	Items                *JSONSchema            `json:"items,omitempty"`
	AdditionalProperties interface{}            `json:"additionalProperties,omitempty"`
	Enum                 []interface{}          `json:"enum,omitempty"`
	Const                interface{}            `json:"const,omitempty"`
	Default              interface{}            `json:"default,omitempty"`
	Examples             []interface{}          `json:"examples,omitempty"`
	Format               string                 `json:"format,omitempty"`
	Pattern              string                 `json:"pattern,omitempty"`
	MinLength            *int                   `json:"minLength,omitempty"`
	MaxLength            *int                   `json:"maxLength,omitempty"`
	Minimum              *float64               `json:"minimum,omitempty"`
	Maximum              *float64               `json:"maximum,omitempty"`
	MultipleOf           *float64               `json:"multipleOf,omitempty"`
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	InputSchema JSONSchema `json:"inputSchema"`
}

// Validate performs basic validation on the tool
func (t *Tool) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("tool name is required")
	}
	return nil
}

// ToolCallRequest represents a tool call request
type ToolCallRequest struct {
	Name      string      `json:"name"`
	Arguments interface{} `json:"arguments,omitempty"`
	Meta      *Meta       `json:"_meta,omitempty"`
}

// Content represents content in various formats
type Content struct {
	Type        string                 `json:"type"`
	Text        string                 `json:"text,omitempty"`
	Data        string                 `json:"data,omitempty"` // base64 encoded
	MimeType    string                 `json:"mimeType,omitempty"`
	Annotations map[string]interface{} `json:"annotations,omitempty"`
}

// ToolResult represents the result of a tool call
type ToolResult struct {
	Content []Content `json:"content,omitempty"`
	IsError bool      `json:"isError,omitempty"`
	Meta    *Meta     `json:"_meta,omitempty"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string                 `json:"uri"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ResourceTemplate represents a resource template
type ResourceTemplate struct {
	URITemplate string                 `json:"uriTemplate"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ResourceContent represents the content of a resource
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"` // base64 encoded
}

// Prompt represents an MCP prompt
type Prompt struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Arguments   []PromptArgument       `json:"arguments,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PromptArgument represents a prompt argument
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptMessage represents a message in a prompt
type PromptMessage struct {
	Role    MessageRole `json:"role"`
	Content Content     `json:"content"`
}

// MessageRole represents the role of a message
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// === MCP Request/Response Types ===

// ListResourcesRequest represents a resources/list request
type ListResourcesRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListResourcesResponse represents a resources/list response
type ListResourcesResponse struct {
	Resources  []Resource `json:"resources"`
	NextCursor string     `json:"nextCursor,omitempty"`
}

// ReadResourceRequest represents a resources/read request
type ReadResourceRequest struct {
	URI string `json:"uri"`
}

// ReadResourceResponse represents a resources/read response
type ReadResourceResponse struct {
	Contents []ResourceContent `json:"contents"`
}

// SubscribeRequest represents a resources/subscribe request
type SubscribeRequest struct {
	URI string `json:"uri"`
}

// UnsubscribeRequest represents a resources/unsubscribe request
type UnsubscribeRequest struct {
	URI string `json:"uri"`
}

// ResourceUpdatedNotification represents a notifications/resources/updated notification
type ResourceUpdatedNotification struct {
	URI string `json:"uri"`
}

// ListPromptsRequest represents a prompts/list request
type ListPromptsRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListPromptsResponse represents a prompts/list response
type ListPromptsResponse struct {
	Prompts    []Prompt `json:"prompts"`
	NextCursor string   `json:"nextCursor,omitempty"`
}

// GetPromptRequest represents a prompts/get request
type GetPromptRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// GetPromptResponse represents a prompts/get response
type GetPromptResponse struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}

// === Tool Operations ===

// ListToolsRequest represents a tools/list request
type ListToolsRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListToolsResult represents a tools/list response
type ListToolsResult struct {
	Tools      []Tool `json:"tools"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// CallToolRequest represents a tools/call request
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult represents a tools/call response
type CallToolResult struct {
	Content  []Content              `json:"content"`
	IsError  bool                   `json:"isError,omitempty"`
	Metadata map[string]interface{} `json:"_meta,omitempty"`
}

// CreateMessageOptions represents options for creating messages
func NewJSONRPCRequest(method string, params interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		Version: JSONRPCVersion,
		ID:      generateRequestID(),
		Method:  method,
		Params:  params,
	}
}

// NewJSONRPCNotification creates a new JSON-RPC notification
func NewJSONRPCNotification(method string, params interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		Version: JSONRPCVersion,
		Method:  method,
		Params:  params,
	}
}

// NewJSONRPCResponse creates a new JSON-RPC response
func NewJSONRPCResponse(id interface{}, result interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		Version: JSONRPCVersion,
		ID:      id,
		Result:  result,
	}
}

// NewJSONRPCError creates a new JSON-RPC error response
func NewJSONRPCError(id interface{}, code int, message string, data interface{}) *JSONRPCMessage {
	return &JSONRPCMessage{
		Version: JSONRPCVersion,
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// IsRequest checks if the message is a request
func (m *JSONRPCMessage) IsRequest() bool {
	return m.Method != "" && m.ID != nil
}

// IsNotification checks if the message is a notification
func (m *JSONRPCMessage) IsNotification() bool {
	return m.Method != "" && m.ID == nil
}

// IsResponse checks if the message is a response
func (m *JSONRPCMessage) IsResponse() bool {
	return m.Method == "" && (m.Result != nil || m.Error != nil)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
