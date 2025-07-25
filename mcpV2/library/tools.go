// Package library provides pure library components for MCP without requiring
// any transport or client-server infrastructure.
//
// This package enables embedding MCP functionality directly into applications
// as library calls, providing maximum performance and flexibility without
// the overhead of JSON-RPC communication.
package library

import (
	"context"
	"fmt"
	"sync"

	"github.com/benozo/neuron-mcp/protocol"
)

// ComponentRegistry provides pure library access to MCP components
type ComponentRegistry struct {
	tools     ToolRegistry
	memory    Memory
	processor *Processor
}

// NewComponentRegistry creates library-only MCP components
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		tools:     NewToolRegistry(),
		memory:    NewInMemoryBackend(nil),
		processor: NewProcessor(),
	}
}

// Tools returns the tool registry
func (cr *ComponentRegistry) Tools() ToolRegistry {
	return cr.tools
}

// Memory returns the memory interface
func (cr *ComponentRegistry) Memory() Memory {
	return cr.memory
}

// Processor returns the data processor
func (cr *ComponentRegistry) Processor() *Processor {
	return cr.processor
}

// SetMemory sets a custom memory backend
func (cr *ComponentRegistry) SetMemory(memory Memory) {
	cr.memory = memory
}

// ToolRegistry manages tools without server dependency
type ToolRegistry interface {
	// Register adds a tool to the registry
	Register(name string, handler ToolFunc) error

	// RegisterWithSchema adds a tool with explicit schema
	RegisterWithSchema(name string, handler ToolFunc, schema *protocol.JSONSchema) error

	// Unregister removes a tool from the registry
	Unregister(name string) error

	// List returns all registered tool names
	List() []string

	// Get retrieves a tool by name
	Get(name string) (*protocol.Tool, error)

	// Call executes a tool directly
	Call(ctx context.Context, name string, params map[string]interface{}) (*protocol.ToolResult, error)

	// GetSchema returns the schema for a tool
	GetSchema(name string) (*protocol.JSONSchema, error)

	// Validate validates parameters against a tool's schema
	Validate(name string, params map[string]interface{}) error
}

// ToolFunc is the function signature for tools in library mode
type ToolFunc func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error)

// toolRegistry implements ToolRegistry
type toolRegistry struct {
	tools map[string]*toolEntry
	mu    sync.RWMutex
}

// toolEntry represents a registered tool
type toolEntry struct {
	tool    *protocol.Tool
	handler ToolFunc
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() ToolRegistry {
	return &toolRegistry{
		tools: make(map[string]*toolEntry),
	}
}

// Register implements ToolRegistry.Register
func (tr *toolRegistry) Register(name string, handler ToolFunc) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	// Create a basic tool with minimal schema
	tool := &protocol.Tool{
		Name:        name,
		Description: fmt.Sprintf("Tool: %s", name),
		InputSchema: protocol.JSONSchema{
			Type: "object",
		},
	}

	tr.tools[name] = &toolEntry{
		tool:    tool,
		handler: handler,
	}

	return nil
}

// RegisterWithSchema implements ToolRegistry.RegisterWithSchema
func (tr *toolRegistry) RegisterWithSchema(name string, handler ToolFunc, schema *protocol.JSONSchema) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	tool := &protocol.Tool{
		Name:        name,
		Description: fmt.Sprintf("Tool: %s", name),
		InputSchema: *schema,
	}

	tr.tools[name] = &toolEntry{
		tool:    tool,
		handler: handler,
	}

	return nil
}

// Unregister implements ToolRegistry.Unregister
func (tr *toolRegistry) Unregister(name string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if _, exists := tr.tools[name]; !exists {
		return fmt.Errorf("tool not found: %s", name)
	}

	delete(tr.tools, name)
	return nil
}

// List implements ToolRegistry.List
func (tr *toolRegistry) List() []string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	names := make([]string, 0, len(tr.tools))
	for name := range tr.tools {
		names = append(names, name)
	}

	return names
}

// Get implements ToolRegistry.Get
func (tr *toolRegistry) Get(name string) (*protocol.Tool, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	entry, exists := tr.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	// Return a copy to prevent modification
	tool := *entry.tool
	return &tool, nil
}

// Call implements ToolRegistry.Call
func (tr *toolRegistry) Call(ctx context.Context, name string, params map[string]interface{}) (*protocol.ToolResult, error) {
	tr.mu.RLock()
	entry, exists := tr.tools[name]
	tr.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	// Validate parameters if schema is available
	if err := tr.validateParams(entry.tool.InputSchema, params); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}

	// Call the handler
	return entry.handler(ctx, params)
}

// GetSchema implements ToolRegistry.GetSchema
func (tr *toolRegistry) GetSchema(name string) (*protocol.JSONSchema, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	entry, exists := tr.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	// Return a copy to prevent modification
	schema := entry.tool.InputSchema
	return &schema, nil
}

// Validate implements ToolRegistry.Validate
func (tr *toolRegistry) Validate(name string, params map[string]interface{}) error {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	entry, exists := tr.tools[name]
	if !exists {
		return fmt.Errorf("tool not found: %s", name)
	}

	return tr.validateParams(entry.tool.InputSchema, params)
}

// validateParams performs basic schema validation
func (tr *toolRegistry) validateParams(schema protocol.JSONSchema, params map[string]interface{}) error {
	// Basic validation - check required fields
	for _, required := range schema.Required {
		if _, exists := params[required]; !exists {
			return fmt.Errorf("required parameter missing: %s", required)
		}
	}

	// Additional validation can be added here
	return nil
}

// Processor provides data processing utilities
type Processor struct {
	// Future: data transformation, filtering, etc.
}

// NewProcessor creates a new data processor
func NewProcessor() *Processor {
	return &Processor{}
}

// Transform applies a transformation to data (placeholder)
func (p *Processor) Transform(data interface{}, operation string) (interface{}, error) {
	// Placeholder for data transformation logic
	return data, nil
}
