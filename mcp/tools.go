package mcp

import "fmt"

type ToolFunc func(params map[string]interface{}, memory *Memory) (interface{}, error)

type ToolRegistry struct {
	tools map[string]ToolFunc
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[string]ToolFunc)}
}

func (r *ToolRegistry) Register(name string, fn ToolFunc) {
	r.tools[name] = fn
}

func (r *ToolRegistry) Call(name string, params map[string]interface{}, memory *Memory) (interface{}, error) {
	if tool, ok := r.tools[name]; ok {
		return tool(params, memory)
	}
	return nil, ErrToolNotFound(name)
}

func (r *ToolRegistry) GetRegisteredTools() []string {
	var names []string
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

func ErrToolNotFound(name string) error {
	return fmt.Errorf("tool not found: %s", name)
}
