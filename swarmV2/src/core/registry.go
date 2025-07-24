package core

import (
	"errors"
	"sync"
)

// Registry manages the registration of agents and tools within the swarm framework.
type Registry struct {
	agents map[string]Agent
	tools  map[string]Tool
	mu     sync.RWMutex
}

// NewRegistry creates a new instance of Registry.
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]Agent),
		tools:  make(map[string]Tool),
	}
}

// RegisterAgent registers a new agent in the registry.
func (r *Registry) RegisterAgent(name string, agent Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[name]; exists {
		return errors.New("agent already registered")
	}
	r.agents[name] = agent
	return nil
}

// GetAgent retrieves an agent by name from the registry.
func (r *Registry) GetAgent(name string) (Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[name]
	if !exists {
		return nil, errors.New("agent not found")
	}
	return agent, nil
}

// RegisterTool registers a new tool in the registry.
func (r *Registry) RegisterTool(name string, tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[name]; exists {
		return errors.New("tool already registered")
	}
	r.tools[name] = tool
	return nil
}

// GetTool retrieves a tool by name from the registry.
func (r *Registry) GetTool(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, errors.New("tool not found")
	}
	return tool, nil
}

// UnregisterAgent removes an agent from the registry.
func (r *Registry) UnregisterAgent(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[name]; !exists {
		return errors.New("agent not found")
	}
	delete(r.agents, name)
	return nil
}

// UnregisterTool removes a tool from the registry.
func (r *Registry) UnregisterTool(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[name]; !exists {
		return errors.New("tool not found")
	}
	delete(r.tools, name)
	return nil
}

// ListAgents returns a list of all registered agent names.
func (r *Registry) ListAgents() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.agents))
	for name := range r.agents {
		names = append(names, name)
	}
	return names
}

// ListTools returns a list of all registered tool names.
func (r *Registry) ListTools() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetAgentCount returns the number of registered agents.
func (r *Registry) GetAgentCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.agents)
}

// GetToolCount returns the number of registered tools.
func (r *Registry) GetToolCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tools)
}

// Clear removes all agents and tools from the registry.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents = make(map[string]Agent)
	r.tools = make(map[string]Tool)
}
