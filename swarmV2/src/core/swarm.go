package core

import (
	"errors"
	"sync"
)

// Swarm manages a collection of agents.
type Swarm struct {
	agents map[string]Agent
	mu     sync.RWMutex
}

// NewSwarm creates a new instance of Swarm.
func NewSwarm() *Swarm {
	return &Swarm{
		agents: make(map[string]Agent),
	}
}

// AddAgent adds a new agent to the swarm.
func (s *Swarm) AddAgent(agent Agent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agentName := agent.GetName()
	if _, exists := s.agents[agentName]; exists {
		return errors.New("agent already exists")
	}
	s.agents[agentName] = agent
	return nil
}

// RemoveAgent removes an agent from the swarm by name.
func (s *Swarm) RemoveAgent(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.agents[name]; !exists {
		return errors.New("agent not found")
	}
	delete(s.agents, name)
	return nil
}

// Execute runs a specified action on all agents in the swarm.
func (s *Swarm) Execute(action func(agent Agent)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, agent := range s.agents {
		action(agent)
	}
}

// GetAgent retrieves an agent by name.
func (s *Swarm) GetAgent(name string) (Agent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agent, exists := s.agents[name]
	if !exists {
		return nil, errors.New("agent not found")
	}
	return agent, nil
}

// ListAgents returns a list of all agent names in the swarm.
func (s *Swarm) ListAgents() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.agents))
	for name := range s.agents {
		names = append(names, name)
	}
	return names
}

// GetAgentCount returns the number of agents in the swarm.
func (s *Swarm) GetAgentCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.agents)
}
