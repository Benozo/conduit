package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPTool represents an MCP tool definition
type MCPTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// MCPToolsListResult represents the result of tools/list
type MCPToolsListResult struct {
	Tools []MCPTool `json:"tools"`
}

// MCPToolCallParams represents parameters for tools/call
type MCPToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// MCPToolCallResult represents the result of tools/call
type MCPToolCallResult struct {
	Content []MCPContent `json:"content"`
}

// MCPContent represents content in MCP responses
type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MCPInitializeParams represents parameters for initialize
type MCPInitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]interface{} `json:"clientInfo"`
}

// MCPInitializeResult represents the result of initialize
type MCPInitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      map[string]interface{} `json:"serverInfo"`
}

// EnhancedSchemaProvider interface for servers that provide custom tool schemas
type EnhancedSchemaProvider interface {
	GetToolMetadata() map[string]interface{}           // Returns custom tool metadata
	GetToolSchema(toolName string) (interface{}, bool) // Returns schema for specific tool
}

// StdioServer handles MCP over stdio (for Copilot integration)
type StdioServer struct {
	tools          *ToolRegistry
	memory         *Memory
	input          io.Reader
	output         io.Writer
	logger         *log.Logger
	schemaProvider EnhancedSchemaProvider // Optional enhanced schema provider
}

// NewStdioServer creates a new stdio-based MCP server
func NewStdioServer(tools *ToolRegistry, memory *Memory) *StdioServer {
	return &StdioServer{
		tools:  tools,
		memory: memory,
		input:  os.Stdin,
		output: os.Stdout,
		logger: log.New(os.Stderr, "[MCP] ", log.LstdFlags),
	}
}

// NewStdioServerWithSchemaProvider creates a new stdio server with enhanced schema support
func NewStdioServerWithSchemaProvider(tools *ToolRegistry, memory *Memory, provider EnhancedSchemaProvider) *StdioServer {
	return &StdioServer{
		tools:          tools,
		memory:         memory,
		input:          os.Stdin,
		output:         os.Stdout,
		logger:         log.New(os.Stderr, "[MCP] ", log.LstdFlags),
		schemaProvider: provider,
	}
}

// SetIO allows customizing input/output streams (useful for testing)
func (s *StdioServer) SetIO(input io.Reader, output io.Writer) {
	s.input = input
	s.output = output
}

// Run starts the stdio server
func (s *StdioServer) Run() error {
	s.logger.Println("MCP Stdio Server starting...")

	scanner := bufio.NewScanner(s.input)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(req.ID, -32700, "Parse error")
			continue
		}

		s.handleRequest(req)
	}

	return scanner.Err()
}

// handleRequest processes individual JSON-RPC requests
func (s *StdioServer) handleRequest(req JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "notifications/initialized":
		// No response needed for notifications
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolCall(req)
	default:
		s.sendError(req.ID, -32601, "Method not found")
	}
}

// handleInitialize processes initialize requests
func (s *StdioServer) handleInitialize(req JSONRPCRequest) {
	var params MCPInitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}

	result := MCPInitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		ServerInfo: map[string]interface{}{
			"name":    "conduit-server",
			"version": "1.0.0",
		},
	}
	s.sendResult(req.ID, result)
}

// handleToolsList processes tools/list requests
func (s *StdioServer) handleToolsList(req JSONRPCRequest) {
	// Get tool schemas dynamically from the registry
	tools := s.getToolSchemas()

	result := MCPToolsListResult{Tools: tools}
	s.sendResult(req.ID, result)
}

// handleToolCall processes tools/call requests
func (s *StdioServer) handleToolCall(req JSONRPCRequest) {
	var params MCPToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}

	result, err := s.tools.Call(params.Name, params.Arguments, s.memory)
	if err != nil {
		s.sendError(req.ID, -32601, fmt.Sprintf("Tool error: %v", err))
		return
	}

	// Convert result to MCP format
	resultText := s.formatToolResult(result)

	mcpResult := MCPToolCallResult{
		Content: []MCPContent{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}
	s.sendResult(req.ID, mcpResult)
}

// getToolSchemas dynamically generates tool schemas from registered tools
func (s *StdioServer) getToolSchemas() []MCPTool {
	var mcpTools []MCPTool

	// Get all registered tool names
	toolNames := s.tools.GetRegisteredTools()

	for _, name := range toolNames {
		// Create a basic schema for each tool
		tool := MCPTool{
			Name:        name,
			Description: s.getToolDescription(name),
			InputSchema: s.getToolInputSchema(name),
		}
		mcpTools = append(mcpTools, tool)
	}

	return mcpTools
}

// getToolDescription returns a description for a tool
func (s *StdioServer) getToolDescription(name string) string {
	// Check if enhanced schema provider has custom description
	if s.schemaProvider != nil {
		if schema, exists := s.schemaProvider.GetToolSchema(name); exists {
			if schemaMap, ok := schema.(map[string]interface{}); ok {
				if desc, ok := schemaMap["description"].(string); ok {
					return desc
				}
			}
		}
	}

	// Fallback to built-in descriptions
	descriptions := map[string]string{
		"uppercase":         "Convert text to uppercase",
		"lowercase":         "Convert text to lowercase",
		"reverse":           "Reverse a string",
		"word_count":        "Count words in text",
		"trim":              "Remove leading and trailing whitespace",
		"title_case":        "Convert text to title case",
		"snake_case":        "Convert text to snake_case",
		"camel_case":        "Convert text to camelCase",
		"replace":           "Replace text patterns",
		"extract_words":     "Extract words from text",
		"sort_words":        "Sort words alphabetically",
		"char_count":        "Count characters in text",
		"remove_whitespace": "Remove all whitespace from text",
		"remember":          "Store a value in memory",
		"recall":            "Retrieve a value from memory",
		"forget":            "Remove a value from memory",
		"list_memory":       "List all stored memory keys",
		"clear_memory":      "Clear all memory",
		"timestamp":         "Get current timestamp",
		"uuid":              "Generate a UUID",
		"hash":              "Generate hash of input",
		"base64_encode":     "Encode text to base64",
		"base64_decode":     "Decode base64 to text",
		"url_encode":        "URL encode text",
		"url_decode":        "URL decode text",
	}

	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return fmt.Sprintf("Tool: %s", name)
}

// getToolInputSchema returns input schema for a tool
func (s *StdioServer) getToolInputSchema(name string) interface{} {
	// Check if enhanced schema provider has custom schema
	if s.schemaProvider != nil {
		if schema, exists := s.schemaProvider.GetToolSchema(name); exists {
			if schemaMap, ok := schema.(map[string]interface{}); ok {
				if inputSchema, ok := schemaMap["inputSchema"]; ok {
					return inputSchema
				}
			}
		}
	}

	// Fallback to built-in schemas
	// Most tools use text parameter
	textSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"text": map[string]interface{}{
				"type":        "string",
				"description": "Input text",
			},
		},
		"required": []string{"text"},
	}

	// Special schemas for specific tools
	switch name {
	case "remember":
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"key": map[string]interface{}{
					"type":        "string",
					"description": "Memory key",
				},
				"value": map[string]interface{}{
					"type":        "string",
					"description": "Value to store",
				},
			},
			"required": []string{"key", "value"},
		}
	case "recall", "forget":
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"key": map[string]interface{}{
					"type":        "string",
					"description": "Memory key",
				},
			},
			"required": []string{"key"},
		}
	case "replace":
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"text": map[string]interface{}{
					"type":        "string",
					"description": "Input text",
				},
				"old": map[string]interface{}{
					"type":        "string",
					"description": "Text to replace",
				},
				"new": map[string]interface{}{
					"type":        "string",
					"description": "Replacement text",
				},
			},
			"required": []string{"text", "old", "new"},
		}
	case "timestamp":
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"format": map[string]interface{}{
					"type":        "string",
					"description": "Timestamp format (iso, unix, readable)",
				},
			},
		}
	case "hash":
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"text": map[string]interface{}{
					"type":        "string",
					"description": "Text to hash",
				},
				"algorithm": map[string]interface{}{
					"type":        "string",
					"description": "Hash algorithm (md5, sha256, sha1)",
				},
			},
			"required": []string{"text"},
		}
	default:
		return textSchema
	}
}

// formatToolResult formats tool results for MCP responses
func (s *StdioServer) formatToolResult(result interface{}) string {
	if resultMap, ok := result.(map[string]interface{}); ok {
		if res, exists := resultMap["result"]; exists {
			return fmt.Sprintf("%v", res)
		}
	}
	return fmt.Sprintf("%v", result)
}

// sendResult sends a successful JSON-RPC response
func (s *StdioServer) sendResult(id interface{}, result interface{}) {
	response := JSONRPCResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
	}
	data, _ := json.Marshal(response)
	fmt.Fprintln(s.output, string(data))
}

// sendError sends an error JSON-RPC response
func (s *StdioServer) sendError(id interface{}, code int, message string) {
	response := JSONRPCResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
	data, _ := json.Marshal(response)
	fmt.Fprintln(s.output, string(data))
}
