// Package client provides MCP client implementation for connecting to MCP servers.
//
// This package implements the client side of the Model Context Protocol,
// handling connection management, tool calling, resource access, and all
// client-server communication.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/protocol"
	"github.com/modelcontextprotocol/go-sdk/transport"
)

// Client provides MCP client functionality
type Client struct {
	transport    transport.Transport
	capabilities protocol.ClientCapabilities
	serverInfo   *protocol.Implementation

	// Cached server state
	tools     map[string]*protocol.Tool
	resources map[string]*protocol.Resource
	prompts   map[string]*protocol.Prompt

	// Connection state
	initialized bool
	mu          sync.RWMutex

	// Configuration
	options *ClientOptions
}

// ClientOptions configures client behavior
type ClientOptions struct {
	// Timeout for individual requests
	Timeout time.Duration

	// Timeout for the initial connection
	ConnectTimeout time.Duration

	// Client information
	ClientInfo protocol.Implementation

	// Progress handler for tracking long operations
	ProgressHandler ProgressHandler

	// Logger for client operations
	Logger Logger

	// Whether to automatically cache tools/resources/prompts on connect
	AutoCache bool
}

// ProgressHandler receives progress updates
type ProgressHandler func(token string, progress float64, total int64)

// Logger interface for client logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// DefaultClientOptions returns default client options
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Timeout:        30 * time.Second,
		ConnectTimeout: 10 * time.Second,
		ClientInfo: protocol.Implementation{
			Name:    "go-mcp-client",
			Version: "1.0.0",
		},
		AutoCache: true,
	}
}

// NewClient creates a new MCP client
func NewClient(transport transport.Transport, opts *ClientOptions) *Client {
	if opts == nil {
		opts = DefaultClientOptions()
	}

	return &Client{
		transport: transport,
		tools:     make(map[string]*protocol.Tool),
		resources: make(map[string]*protocol.Resource),
		prompts:   make(map[string]*protocol.Prompt),
		options:   opts,
	}
}

// Connect establishes connection and performs MCP handshake
func (c *Client) Connect(ctx context.Context, capabilities protocol.ClientCapabilities) error {
	if c.options.ConnectTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.ConnectTimeout)
		defer cancel()
	}

	c.mu.Lock()
	c.capabilities = capabilities
	c.mu.Unlock()

	// Send initialize request
	req := &protocol.InitializeRequest{
		ProtocolVersion: protocol.MCPProtocolVersion,
		Capabilities:    capabilities,
		ClientInfo:      c.options.ClientInfo,
	}

	if c.options.Logger != nil {
		c.options.Logger.Info("Connecting to MCP server")
	}

	result, err := c.sendRequest(ctx, "initialize", req)
	if err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}

	// Parse initialize result
	var initResult protocol.InitializeResult
	if err := json.Unmarshal(result, &initResult); err != nil {
		return fmt.Errorf("invalid initialize response: %w", err)
	}

	c.mu.Lock()
	c.serverInfo = &initResult.ServerInfo
	c.initialized = true
	c.mu.Unlock()

	if c.options.Logger != nil {
		c.options.Logger.Info("Connected to server: %s %s", initResult.ServerInfo.Name, initResult.ServerInfo.Version)
	}

	// Send initialized notification
	if err := c.sendNotification(ctx, "notifications/initialized", struct{}{}); err != nil {
		return fmt.Errorf("failed to send initialized notification: %w", err)
	}

	// Auto-cache server capabilities if enabled
	if c.options.AutoCache {
		if err := c.cacheServerCapabilities(ctx, initResult.Capabilities); err != nil {
			if c.options.Logger != nil {
				c.options.Logger.Warn("Failed to cache server capabilities: %v", err)
			}
		}
	}

	return nil
}

// IsConnected returns true if the client is connected and initialized
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.initialized && c.transport.IsConnected()
}

// GetServerInfo returns information about the connected server
func (c *Client) GetServerInfo() *protocol.Implementation {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.serverInfo == nil {
		return nil
	}

	// Return a copy to prevent modification
	info := *c.serverInfo
	return &info
}

// ListTools retrieves available tools from server
func (c *Client) ListTools(ctx context.Context) ([]*protocol.Tool, error) {
	if !c.IsConnected() {
		return nil, protocol.ErrNotInitialized
	}

	result, err := c.sendRequest(ctx, "tools/list", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	var response protocol.ListToolsResult
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("invalid tools list response: %w", err)
	}

	// Cache tools for quick access
	c.mu.Lock()
	for i := range response.Tools {
		c.tools[response.Tools[i].Name] = &response.Tools[i]
	}
	c.mu.Unlock()

	if c.options.Logger != nil {
		c.options.Logger.Debug("Listed %d tools", len(response.Tools))
	}

	// Convert to slice of pointers
	tools := make([]*protocol.Tool, len(response.Tools))
	for i := range response.Tools {
		tools[i] = &response.Tools[i]
	}

	return tools, nil
}

// CallTool executes a tool with given parameters
func (c *Client) CallTool(ctx context.Context, name string, params interface{}) (*protocol.ToolResult, error) {
	if !c.IsConnected() {
		return nil, protocol.ErrNotInitialized
	}

	req := &protocol.ToolCallRequest{
		Name:      name,
		Arguments: params,
	}

	// Add progress tracking if supported
	if c.capabilities.Experimental != nil {
		// Progress tracking is experimental for now
		if c.options.ProgressHandler != nil {
			meta := &protocol.Meta{
				ProgressToken: generateProgressToken(),
			}
			req.Meta = meta
		}
	}

	if c.options.Logger != nil {
		c.options.Logger.Debug("Calling tool: %s", name)
	}

	result, err := c.sendRequest(ctx, "tools/call", req)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	var toolResult protocol.ToolResult
	if err := json.Unmarshal(result, &toolResult); err != nil {
		return nil, fmt.Errorf("invalid tool result: %w", err)
	}

	if c.options.Logger != nil {
		c.options.Logger.Debug("Tool call completed: %s", name)
	}

	return &toolResult, nil
}

// ListResources retrieves all available resources from the server
func (c *Client) ListResources(ctx context.Context) ([]*protocol.Resource, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("client not connected")
	}

	// Apply timeout if configured
	if c.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}

	req := &protocol.ListResourcesRequest{}

	result, err := c.sendRequest(ctx, "resources/list", req)
	if err != nil {
		return nil, fmt.Errorf("list resources failed: %w", err)
	}

	var response protocol.ListResourcesResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("invalid resources list response: %w", err)
	}

	// Update cached resources
	c.mu.Lock()
	c.resources = make(map[string]*protocol.Resource)
	resources := make([]*protocol.Resource, len(response.Resources))
	for i, resource := range response.Resources {
		resourceCopy := resource
		c.resources[resource.URI] = &resourceCopy
		resources[i] = &resourceCopy
	}
	c.mu.Unlock()

	if c.options.Logger != nil {
		c.options.Logger.Debug("Listed %d resources", len(resources))
	}

	return resources, nil
}

// ReadResource reads the content of a specific resource
func (c *Client) ReadResource(ctx context.Context, uri string) (*protocol.ReadResourceResponse, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("client not connected")
	}

	// Apply timeout if configured
	if c.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}

	req := &protocol.ReadResourceRequest{
		URI: uri,
	}

	result, err := c.sendRequest(ctx, "resources/read", req)
	if err != nil {
		return nil, fmt.Errorf("read resource failed: %w", err)
	}

	var response protocol.ReadResourceResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("invalid resource read response: %w", err)
	}

	if c.options.Logger != nil {
		c.options.Logger.Debug("Read resource: %s", uri)
	}

	return &response, nil
}

// ListPrompts retrieves all available prompts from the server
func (c *Client) ListPrompts(ctx context.Context) ([]*protocol.Prompt, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("client not connected")
	}

	// Apply timeout if configured
	if c.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}

	req := &protocol.ListPromptsRequest{}

	result, err := c.sendRequest(ctx, "prompts/list", req)
	if err != nil {
		return nil, fmt.Errorf("list prompts failed: %w", err)
	}

	var response protocol.ListPromptsResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("invalid prompts list response: %w", err)
	}

	// Update cached prompts
	c.mu.Lock()
	c.prompts = make(map[string]*protocol.Prompt)
	prompts := make([]*protocol.Prompt, len(response.Prompts))
	for i, prompt := range response.Prompts {
		promptCopy := prompt
		c.prompts[prompt.Name] = &promptCopy
		prompts[i] = &promptCopy
	}
	c.mu.Unlock()

	if c.options.Logger != nil {
		c.options.Logger.Debug("Listed %d prompts", len(prompts))
	}

	return prompts, nil
}

// GetPrompt retrieves a specific prompt with optional arguments
func (c *Client) GetPrompt(ctx context.Context, name string, args map[string]interface{}) (*protocol.GetPromptResponse, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("client not connected")
	}

	// Apply timeout if configured
	if c.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}

	req := &protocol.GetPromptRequest{
		Name:      name,
		Arguments: args,
	}

	result, err := c.sendRequest(ctx, "prompts/get", req)
	if err != nil {
		return nil, fmt.Errorf("get prompt failed: %w", err)
	}

	var response protocol.GetPromptResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, fmt.Errorf("invalid prompt get response: %w", err)
	}

	if c.options.Logger != nil {
		c.options.Logger.Debug("Got prompt: %s", name)
	}

	return &response, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initialized = false

	if c.options.Logger != nil {
		c.options.Logger.Info("Closing client connection")
	}

	return c.transport.Close()
}

// sendRequest sends a request and waits for a response
func (c *Client) sendRequest(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	if c.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}

	msg := protocol.NewJSONRPCRequest(method, params)

	if err := c.transport.Send(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for response
	for {
		response, err := c.transport.Receive(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to receive response: %w", err)
		}

		// Check if this is the response to our request
		if response.IsResponse() && response.ID == msg.ID {
			if response.Error != nil {
				return nil, response.Error
			}
			result, err := json.Marshal(response.Result)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal result: %w", err)
			}
			return result, nil
		}

		// Handle notifications or other messages
		if response.IsNotification() {
			c.handleNotification(response)
		}
	}
}

// sendNotification sends a notification (no response expected)
func (c *Client) sendNotification(ctx context.Context, method string, params interface{}) error {
	if c.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.Timeout)
		defer cancel()
	}

	msg := protocol.NewJSONRPCNotification(method, params)
	return c.transport.Send(ctx, msg)
}

// handleNotification handles incoming notifications
func (c *Client) handleNotification(msg *protocol.JSONRPCMessage) {
	switch msg.Method {
	case "notifications/progress":
		if c.options.ProgressHandler != nil {
			var progress protocol.ProgressNotification
			if data, err := json.Marshal(msg.Params); err == nil {
				if err := json.Unmarshal(data, &progress); err == nil {
					c.options.ProgressHandler(progress.ProgressToken, progress.Progress, progress.Total)
				}
			}
		}
	default:
		if c.options.Logger != nil {
			c.options.Logger.Debug("Received unhandled notification: %s", msg.Method)
		}
	}
}

// cacheServerCapabilities caches server capabilities on connection
func (c *Client) cacheServerCapabilities(ctx context.Context, capabilities protocol.ServerCapabilities) error {
	// Cache tools if supported
	if capabilities.Tools != nil {
		if _, err := c.ListTools(ctx); err != nil {
			return fmt.Errorf("failed to cache tools: %w", err)
		}
	}

	// Additional caching for resources and prompts can be added here
	return nil
}

// generateProgressToken generates a unique progress token
func generateProgressToken() string {
	return fmt.Sprintf("progress_%d", time.Now().UnixNano())
}
