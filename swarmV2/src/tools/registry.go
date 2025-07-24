package tools

import (
	"errors"
	"sync"
)

// ToolRegistry manages the registration and retrieval of tools used by agents.
type ToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]interface{}
}

// NewToolRegistry creates a new instance of ToolRegistry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]interface{}),
	}
}

// Register adds a new tool to the registry.
func (tr *ToolRegistry) Register(name string, tool interface{}) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if _, exists := tr.tools[name]; exists {
		return errors.New("tool already registered")
	}

	tr.tools[name] = tool
	return nil
}

// Get retrieves a tool by its name.
func (tr *ToolRegistry) Get(name string) (interface{}, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tool, exists := tr.tools[name]
	if !exists {
		return nil, errors.New("tool not found")
	}

	return tool, nil
}

// List returns all registered tools.
func (tr *ToolRegistry) List() []string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	names := make([]string, 0, len(tr.tools))
	for name := range tr.tools {
		names = append(names, name)
	}
	return names
}