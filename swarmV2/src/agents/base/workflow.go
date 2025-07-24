package base

import (
	"fmt"
	"sync"
)

// Workflow represents a structured workflow that orchestrates the execution of various agents.
type Workflow struct {
	Name     string
	Agents   []Agent
	Mutex    sync.Mutex
}

// NewWorkflow creates a new Workflow instance with the specified name and agents.
func NewWorkflow(name string, agents []Agent) *Workflow {
	return &Workflow{
		Name:   name,
		Agents: agents,
	}
}

// Execute runs the workflow by executing each agent in sequence.
func (w *Workflow) Execute() error {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	fmt.Printf("Executing workflow: %s\n", w.Name)
	for _, agent := range w.Agents {
		if err := agent.Execute(); err != nil {
			return fmt.Errorf("failed to execute agent %s: %w", agent.GetName(), err)
		}
	}
	return nil
}

// AddAgent adds a new agent to the workflow.
func (w *Workflow) AddAgent(agent Agent) {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	w.Agents = append(w.Agents, agent)
}

// Agent interface defines the methods that an agent must implement.
type Agent interface {
	Execute() error
	GetName() string
}