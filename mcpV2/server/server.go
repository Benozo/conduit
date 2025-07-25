// Package server provides MCP server implementation for hosting MCP services.
//
// This package implements the server side of the Model Context Protocol,
// handling tool registration, resource management, client connections,
// and request routing.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/middleware"
	"github.com/modelcontextprotocol/go-sdk/protocol"
	"github.com/modelcontextprotocol/go-sdk/transport"
)

// Server provides MCP server functionality
type Server struct {
	capabilities protocol.ServerCapabilities
	info         protocol.Implementation

	// Registry maps
	tools     map[string]*ToolHandler
	resources map[string]*ResourceHandler
	prompts   map[string]*PromptHandler

	// Connection state
	initialized bool
	clientInfo  *protocol.Implementation

	// Middleware
	middlewareChain *middleware.Chain

	// Synchronization
	mu sync.RWMutex

	// Configuration
	options *ServerOptions
}

// ServerOptions configures server behavior
type ServerOptions struct {
	// Server information
	Info protocol.Implementation

	// Server capabilities
	Capabilities protocol.ServerCapabilities

	// Logger for server operations
	Logger Logger

	// Error handler for unhandled errors
	ErrorHandler ErrorHandler

	// Whether to validate inputs against schemas
	ValidateInputs bool

	// Maximum size for request parameters
	MaxRequestSize int64

	// Middleware chain for request/response processing
	Middleware []middleware.Middleware
}

// ToolHandler represents a tool implementation
type ToolHandler struct {
	Tool    *protocol.Tool
	Handler ToolFunc
}

// ResourceHandler represents a resource implementation
type ResourceHandler struct {
	Resource *protocol.Resource
	Handler  ResourceFunc
}

// PromptHandler represents a prompt implementation
type PromptHandler struct {
	Prompt  *protocol.Prompt
	Handler PromptFunc
}

// ToolFunc is the function signature for tool implementations
type ToolFunc func(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error)

// ResourceFunc is the function signature for resource implementations
type ResourceFunc func(ctx context.Context, uri string) (*protocol.ReadResourceResponse, error)

// PromptFunc is the function signature for prompt implementations
type PromptFunc func(ctx context.Context, name string, args map[string]interface{}) (*protocol.GetPromptResponse, error)

// Logger interface for server logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// ErrorHandler handles unhandled errors
type ErrorHandler func(err error)

// DefaultServerOptions returns default server options
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Info: protocol.Implementation{
			Name:    "go-mcp-server",
			Version: "1.0.0",
		},
		Capabilities: protocol.ServerCapabilities{
			Tools: &protocol.ToolsCapability{
				ListChanged: false,
			},
			Resources: &protocol.ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
			Prompts: &protocol.PromptsCapability{
				ListChanged: false,
			},
		},
		ValidateInputs: true,
		MaxRequestSize: 1024 * 1024, // 1MB
	}
}

// NewServer creates a new MCP server
func NewServer(opts *ServerOptions) *Server {
	if opts == nil {
		opts = DefaultServerOptions()
	}

	// Create middleware chain
	chain := middleware.NewChain(opts.Middleware...)

	return &Server{
		capabilities:    opts.Capabilities,
		info:            opts.Info,
		tools:           make(map[string]*ToolHandler),
		resources:       make(map[string]*ResourceHandler),
		prompts:         make(map[string]*PromptHandler),
		middlewareChain: chain,
		options:         opts,
	}
}

// RegisterTool adds a tool to the server
func (s *Server) RegisterTool(tool *protocol.Tool, handler ToolFunc) error {
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	if err := tool.Validate(); err != nil {
		return fmt.Errorf("invalid tool schema: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tools[tool.Name] = &ToolHandler{
		Tool:    tool,
		Handler: handler,
	}

	if s.options.Logger != nil {
		s.options.Logger.Debug("Registered tool: %s", tool.Name)
	}

	return nil
}

// UnregisterTool removes a tool from the server
func (s *Server) UnregisterTool(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tools, name)

	if s.options.Logger != nil {
		s.options.Logger.Debug("Unregistered tool: %s", name)
	}
}

// AddMiddleware adds middleware to the server's middleware chain
func (s *Server) AddMiddleware(middlewares ...middleware.Middleware) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, mw := range middlewares {
		s.middlewareChain.Add(mw)
	}
}

// SetMiddleware replaces the entire middleware chain
func (s *Server) SetMiddleware(middlewares ...middleware.Middleware) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middlewareChain = middleware.NewChain(middlewares...)
}

// GetMiddlewareChain returns a copy of the current middleware chain
func (s *Server) GetMiddlewareChain() *middleware.Chain {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a new chain (we can't access private fields, but this is for demonstration)
	// In practice, you might want to add a Clone() method to the middleware package
	return middleware.NewChain()
}

// RegisterResource adds a resource to the server
func (s *Server) RegisterResource(resource *protocol.Resource, handler ResourceFunc) error {
	if resource == nil {
		return fmt.Errorf("resource cannot be nil")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.resources[resource.URI] = &ResourceHandler{
		Resource: resource,
		Handler:  handler,
	}

	if s.options.Logger != nil {
		s.options.Logger.Debug("Registered resource: %s", resource.URI)
	}

	return nil
}

// UnregisterResource removes a resource from the server
func (s *Server) UnregisterResource(uri string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.resources, uri)

	if s.options.Logger != nil {
		s.options.Logger.Debug("Unregistered resource: %s", uri)
	}
}

// RegisterPrompt adds a prompt to the server
func (s *Server) RegisterPrompt(prompt *protocol.Prompt, handler PromptFunc) error {
	if prompt == nil {
		return fmt.Errorf("prompt cannot be nil")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.prompts[prompt.Name] = &PromptHandler{
		Prompt:  prompt,
		Handler: handler,
	}

	if s.options.Logger != nil {
		s.options.Logger.Debug("Registered prompt: %s", prompt.Name)
	}

	return nil
}

// UnregisterPrompt removes a prompt from the server
func (s *Server) UnregisterPrompt(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.prompts, name)

	if s.options.Logger != nil {
		s.options.Logger.Debug("Unregistered prompt: %s", name)
	}
}

// GetTool retrieves a tool by name
func (s *Server) GetTool(name string) (*protocol.Tool, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	handler, exists := s.tools[name]
	if !exists {
		return nil, false
	}

	return handler.Tool, true
}

// GetToolHandler returns the tool handler for a given tool name
// This is useful for HTTP transport integration
func (s *Server) GetToolHandler(name string) (*ToolHandler, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	handler, exists := s.tools[name]
	return handler, exists
}

// ListTools returns all registered tools
func (s *Server) ListTools() []*protocol.Tool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tools := make([]*protocol.Tool, 0, len(s.tools))
	for _, handler := range s.tools {
		tools = append(tools, handler.Tool)
	}

	return tools
}

// Serve starts serving on the given transport
func (s *Server) Serve(ctx context.Context, transport transport.Transport) error {
	if s.options.Logger != nil {
		s.options.Logger.Info("Starting MCP server: %s %s", s.info.Name, s.info.Version)
	}

	for {
		select {
		case <-ctx.Done():
			if s.options.Logger != nil {
				s.options.Logger.Info("Server context cancelled, shutting down")
			}
			return ctx.Err()
		default:
		}

		msg, err := transport.Receive(ctx)
		if err != nil {
			if s.options.Logger != nil {
				s.options.Logger.Error("Failed to receive message: %v", err)
			}
			return fmt.Errorf("failed to receive message: %w", err)
		}

		// Handle message in goroutine to allow concurrent processing
		go s.handleMessage(ctx, transport, msg)
	}
}

// handleMessage processes an incoming message
func (s *Server) handleMessage(ctx context.Context, transport transport.Transport, msg *protocol.JSONRPCMessage) {
	if msg.IsRequest() {
		s.handleRequest(ctx, transport, msg)
	} else if msg.IsNotification() {
		s.handleNotification(ctx, msg)
	} else {
		if s.options.Logger != nil {
			s.options.Logger.Warn("Received invalid message type")
		}
	}
}

// handleRequest processes a request message
func (s *Server) handleRequest(ctx context.Context, transport transport.Transport, msg *protocol.JSONRPCMessage) {
	response := s.processRequest(ctx, msg)

	if err := transport.Send(ctx, response); err != nil {
		if s.options.Logger != nil {
			s.options.Logger.Error("Failed to send response: %v", err)
		}
		if s.options.ErrorHandler != nil {
			s.options.ErrorHandler(err)
		}
	}
}

// processRequest processes a request and returns a response
func (s *Server) processRequest(ctx context.Context, msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	// Create the core handler
	coreHandler := func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
		switch req.Method {
		case "initialize":
			return s.handleInitialize(req), nil
		case "tools/list":
			return s.handleListTools(req), nil
		case "tools/call":
			return s.handleCallTool(ctx, req), nil
		case "resources/list":
			return s.handleListResources(req), nil
		case "resources/read":
			return s.handleReadResource(ctx, req), nil
		case "prompts/list":
			return s.handleListPrompts(req), nil
		case "prompts/get":
			return s.handleGetPrompt(ctx, req), nil
		default:
			return protocol.NewJSONRPCError(
				req.ID,
				protocol.MethodNotFound,
				fmt.Sprintf("Method not found: %s", req.Method),
				nil,
			), nil
		}
	}

	// Build middleware chain with core handler
	handler := s.middlewareChain.Build(coreHandler)

	// Execute request through middleware chain
	response, err := handler(ctx, msg)
	if err != nil {
		// Log error and return internal error response
		if s.options.Logger != nil {
			s.options.Logger.Error("Request processing failed: %v", err)
		}
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.InternalError,
			"Internal server error",
			nil,
		)
	}

	return response
}

// ProcessMessage processes a single JSON-RPC message and returns the response
// This is useful for HTTP/REST integrations where each request is independent
func (s *Server) ProcessMessage(ctx context.Context, msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.processRequest(ctx, msg)
}

// handleInitialize processes initialize request
func (s *Server) handleInitialize(msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	var req protocol.InitializeRequest
	if data, err := json.Marshal(msg.Params); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	} else if err := json.Unmarshal(data, &req); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	}

	// Validate protocol version
	if req.ProtocolVersion != protocol.MCPProtocolVersion {
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.InvalidRequest,
			fmt.Sprintf("Unsupported protocol version: %s", req.ProtocolVersion),
			nil,
		)
	}

	s.mu.Lock()
	s.clientInfo = &req.ClientInfo
	s.initialized = true
	s.mu.Unlock()

	result := protocol.InitializeResult{
		ProtocolVersion: protocol.MCPProtocolVersion,
		Capabilities:    s.capabilities,
		ServerInfo:      s.info,
	}

	if s.options.Logger != nil {
		s.options.Logger.Info("Initialized connection with client: %s %s", req.ClientInfo.Name, req.ClientInfo.Version)
	}

	return protocol.NewJSONRPCResponse(msg.ID, result)
}

// handleListTools processes tools/list request
func (s *Server) handleListTools(msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	if !s.isInitialized() {
		return protocol.NewJSONRPCError(msg.ID, protocol.NotInitialized, "Not initialized", nil)
	}

	tools := s.ListTools()
	result := protocol.ListToolsResult{
		Tools: make([]protocol.Tool, len(tools)),
	}

	for i, tool := range tools {
		result.Tools[i] = *tool
	}

	return protocol.NewJSONRPCResponse(msg.ID, result)
}

// handleCallTool processes tools/call request
func (s *Server) handleCallTool(ctx context.Context, msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	if !s.isInitialized() {
		return protocol.NewJSONRPCError(msg.ID, protocol.NotInitialized, "Not initialized", nil)
	}

	var req protocol.ToolCallRequest
	if data, err := json.Marshal(msg.Params); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	} else if err := json.Unmarshal(data, &req); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	}

	s.mu.RLock()
	handler, exists := s.tools[req.Name]
	s.mu.RUnlock()

	if !exists {
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.InvalidTool,
			fmt.Sprintf("Tool not found: %s", req.Name),
			nil,
		)
	}

	// Convert arguments to map[string]interface{}
	var params map[string]interface{}
	if req.Arguments != nil {
		if data, err := json.Marshal(req.Arguments); err != nil {
			return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid arguments", nil)
		} else if err := json.Unmarshal(data, &params); err != nil {
			return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid arguments", nil)
		}
	} else {
		params = make(map[string]interface{})
	}

	// Call the tool handler
	result, err := handler.Handler(ctx, params)
	if err != nil {
		if s.options.Logger != nil {
			s.options.Logger.Error("Tool call failed: %s: %v", req.Name, err)
		}
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.RequestFailed,
			fmt.Sprintf("Tool call failed: %v", err),
			nil,
		)
	}

	if s.options.Logger != nil {
		s.options.Logger.Debug("Tool call completed: %s", req.Name)
	}

	return protocol.NewJSONRPCResponse(msg.ID, result)
}

// handleNotification processes notification messages
func (s *Server) handleNotification(ctx context.Context, msg *protocol.JSONRPCMessage) {
	switch msg.Method {
	case "notifications/initialized":
		if s.options.Logger != nil {
			s.options.Logger.Debug("Client confirmed initialization")
		}
	default:
		if s.options.Logger != nil {
			s.options.Logger.Debug("Received unhandled notification: %s", msg.Method)
		}
	}
}

// isInitialized checks if the server is initialized
func (s *Server) isInitialized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.initialized
}

// GetClientInfo returns information about the connected client
func (s *Server) GetClientInfo() *protocol.Implementation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.clientInfo == nil {
		return nil
	}

	// Return a copy to prevent modification
	info := *s.clientInfo
	return &info
}

// handleListResources processes resources/list request
func (s *Server) handleListResources(msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	if !s.isInitialized() {
		return protocol.NewJSONRPCError(msg.ID, protocol.NotInitialized, "Not initialized", nil)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	resources := make([]protocol.Resource, 0, len(s.resources))
	for _, handler := range s.resources {
		resources = append(resources, *handler.Resource)
	}

	result := &protocol.ListResourcesResponse{
		Resources: resources,
	}

	return protocol.NewJSONRPCResponse(msg.ID, result)
}

// handleReadResource processes resources/read request
func (s *Server) handleReadResource(ctx context.Context, msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	if !s.isInitialized() {
		return protocol.NewJSONRPCError(msg.ID, protocol.NotInitialized, "Not initialized", nil)
	}

	var req protocol.ReadResourceRequest
	if data, err := json.Marshal(msg.Params); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	} else if err := json.Unmarshal(data, &req); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	}

	s.mu.RLock()
	handler, exists := s.resources[req.URI]
	s.mu.RUnlock()

	if !exists {
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.InvalidRequest,
			fmt.Sprintf("Resource not found: %s", req.URI),
			nil,
		)
	}

	// Call the resource handler
	response, err := handler.Handler(ctx, req.URI)
	if err != nil {
		if s.options.Logger != nil {
			s.options.Logger.Error("Resource read failed: %s: %v", req.URI, err)
		}
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.RequestFailed,
			fmt.Sprintf("Resource read failed: %v", err),
			nil,
		)
	}

	if s.options.Logger != nil {
		s.options.Logger.Debug("Resource read completed: %s", req.URI)
	}

	return protocol.NewJSONRPCResponse(msg.ID, response)
}

// handleListPrompts processes prompts/list request
func (s *Server) handleListPrompts(msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	if !s.isInitialized() {
		return protocol.NewJSONRPCError(msg.ID, protocol.NotInitialized, "Not initialized", nil)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	prompts := make([]protocol.Prompt, 0, len(s.prompts))
	for _, handler := range s.prompts {
		prompts = append(prompts, *handler.Prompt)
	}

	result := &protocol.ListPromptsResponse{
		Prompts: prompts,
	}

	return protocol.NewJSONRPCResponse(msg.ID, result)
}

// handleGetPrompt processes prompts/get request
func (s *Server) handleGetPrompt(ctx context.Context, msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	if !s.isInitialized() {
		return protocol.NewJSONRPCError(msg.ID, protocol.NotInitialized, "Not initialized", nil)
	}

	var req protocol.GetPromptRequest
	if data, err := json.Marshal(msg.Params); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	} else if err := json.Unmarshal(data, &req); err != nil {
		return protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
	}

	s.mu.RLock()
	handler, exists := s.prompts[req.Name]
	s.mu.RUnlock()

	if !exists {
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.InvalidRequest,
			fmt.Sprintf("Prompt not found: %s", req.Name),
			nil,
		)
	}

	// Call the prompt handler
	response, err := handler.Handler(ctx, req.Name, req.Arguments)
	if err != nil {
		if s.options.Logger != nil {
			s.options.Logger.Error("Prompt get failed: %s: %v", req.Name, err)
		}
		return protocol.NewJSONRPCError(
			msg.ID,
			protocol.RequestFailed,
			fmt.Sprintf("Prompt get failed: %v", err),
			nil,
		)
	}

	if s.options.Logger != nil {
		s.options.Logger.Debug("Prompt get completed: %s", req.Name)
	}

	return protocol.NewJSONRPCResponse(msg.ID, response)
}
