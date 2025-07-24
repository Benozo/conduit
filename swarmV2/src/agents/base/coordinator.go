package base

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CoordinatorStatus represents the current state of the coordinator
type CoordinatorStatus string

const (
	StatusIdle    CoordinatorStatus = "idle"
	StatusActive  CoordinatorStatus = "active"
	StatusError   CoordinatorStatus = "error"
	StatusStopped CoordinatorStatus = "stopped"
)

// Coordinator is an agent responsible for managing task delegation and workflow execution among other agents.
type Coordinator struct {
	agents      map[string]Agent
	workflowMap map[string]*Workflow
	status      CoordinatorStatus
	startTime   time.Time
	metrics     *CoordinatorMetrics
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// CoordinatorMetrics tracks coordinator performance
type CoordinatorMetrics struct {
	WorkflowsExecuted int           `json:"workflows_executed"`
	WorkflowsFailed   int           `json:"workflows_failed"`
	AgentsRegistered  int           `json:"agents_registered"`
	AverageExecTime   time.Duration `json:"average_execution_time"`
	Uptime            time.Duration `json:"uptime"`
	LastError         *time.Time    `json:"last_error,omitempty"`
}

// NewCoordinator creates a new Coordinator instance.
func NewCoordinator() *Coordinator {
	ctx, cancel := context.WithCancel(context.Background())
	return &Coordinator{
		agents:      make(map[string]Agent),
		workflowMap: make(map[string]*Workflow),
		status:      StatusIdle,
		startTime:   time.Now(),
		metrics:     &CoordinatorMetrics{},
		ctx:         ctx,
		cancel:      cancel,
	}
}

// RegisterAgent registers a new agent with the coordinator.
func (c *Coordinator) RegisterAgent(agent Agent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.agents[agent.GetName()] = agent
	c.metrics.AgentsRegistered++

	fmt.Printf("Agent '%s' registered successfully\n", agent.GetName())
}

// UnregisterAgent removes an agent from the coordinator.
func (c *Coordinator) UnregisterAgent(agentName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.agents[agentName]; exists {
		delete(c.agents, agentName)
		c.metrics.AgentsRegistered--
		fmt.Printf("Agent '%s' unregistered successfully\n", agentName)
	}
}

// ExecuteWorkflow executes a specified workflow using registered agents.
func (c *Coordinator) ExecuteWorkflow(workflowName string, context map[string]interface{}) error {
	startTime := time.Now()

	c.mu.Lock()
	c.status = StatusActive
	workflow, exists := c.workflowMap[workflowName]
	c.mu.Unlock()

	if !exists {
		c.setStatus(StatusError)
		return fmt.Errorf("workflow %s not found", workflowName)
	}

	fmt.Printf("Executing workflow: %s\n", workflowName)

	err := workflow.Execute()

	// Update metrics
	c.mu.Lock()
	if err != nil {
		c.metrics.WorkflowsFailed++
		c.status = StatusError
		now := time.Now()
		c.metrics.LastError = &now
	} else {
		c.metrics.WorkflowsExecuted++
		c.status = StatusIdle
	}

	// Update average execution time
	execTime := time.Since(startTime)
	totalWorkflows := c.metrics.WorkflowsExecuted + c.metrics.WorkflowsFailed
	if totalWorkflows > 0 {
		c.metrics.AverageExecTime = (c.metrics.AverageExecTime*time.Duration(totalWorkflows-1) + execTime) / time.Duration(totalWorkflows)
	}
	c.mu.Unlock()

	if err != nil {
		fmt.Printf("Workflow '%s' failed: %v\n", workflowName, err)
		return err
	}

	fmt.Printf("Workflow '%s' completed successfully in %v\n", workflowName, execTime)
	return nil
}

// ExecuteWorkflowAsync executes a workflow asynchronously and returns a channel for the result
func (c *Coordinator) ExecuteWorkflowAsync(workflowName string, context map[string]interface{}) <-chan error {
	resultChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		err := c.ExecuteWorkflow(workflowName, context)
		resultChan <- err
	}()

	return resultChan
}

// RegisterWorkflow registers a new workflow with the coordinator.
func (c *Coordinator) RegisterWorkflow(workflow *Workflow) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.workflowMap[workflow.Name] = workflow
	fmt.Printf("Workflow '%s' registered successfully\n", workflow.Name)
}

// UnregisterWorkflow removes a workflow from the coordinator
func (c *Coordinator) UnregisterWorkflow(workflowName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.workflowMap[workflowName]; exists {
		delete(c.workflowMap, workflowName)
		fmt.Printf("Workflow '%s' unregistered successfully\n", workflowName)
	}
}

// GetAgents returns a list of registered agents.
func (c *Coordinator) GetAgents() map[string]Agent {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy to prevent external modification
	agentsCopy := make(map[string]Agent)
	for name, agent := range c.agents {
		agentsCopy[name] = agent
	}
	return agentsCopy
}

// GetWorkflows returns a list of registered workflows
func (c *Coordinator) GetWorkflows() map[string]*Workflow {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy to prevent external modification
	workflowsCopy := make(map[string]*Workflow)
	for name, workflow := range c.workflowMap {
		workflowsCopy[name] = workflow
	}
	return workflowsCopy
}

// GetStatus returns the current coordinator status
func (c *Coordinator) GetStatus() CoordinatorStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

// GetMetrics returns the current coordinator metrics
func (c *Coordinator) GetMetrics() *CoordinatorMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Update uptime
	metrics := *c.metrics
	metrics.Uptime = time.Since(c.startTime)
	return &metrics
}

// HasAgent checks if an agent is registered
func (c *Coordinator) HasAgent(agentName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.agents[agentName]
	return exists
}

// HasWorkflow checks if a workflow is registered
func (c *Coordinator) HasWorkflow(workflowName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.workflowMap[workflowName]
	return exists
}

// Stop gracefully shuts down the coordinator
func (c *Coordinator) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.status = StatusStopped
	c.cancel()
	fmt.Println("Coordinator stopped")
}

// setStatus sets the coordinator status (internal method)
func (c *Coordinator) setStatus(status CoordinatorStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = status
}

// String returns a string representation of the coordinator
func (c *Coordinator) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return fmt.Sprintf("Coordinator{Status: %s, Agents: %d, Workflows: %d, Uptime: %v}",
		c.status, len(c.agents), len(c.workflowMap), time.Since(c.startTime))
}
