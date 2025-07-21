package swarm

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkflowType defines different workflow execution patterns
type WorkflowType string

const (
	WorkflowTypeSequential  WorkflowType = "sequential"
	WorkflowTypeParallel    WorkflowType = "parallel"
	WorkflowTypeDAG         WorkflowType = "dag"
	WorkflowTypeSupervisor  WorkflowType = "supervisor"
	WorkflowTypePipeline    WorkflowType = "pipeline"
	WorkflowTypeConditional WorkflowType = "conditional"
)

// WorkflowNode represents a node in a workflow
type WorkflowNode struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Agent        *Agent                 `json:"agent"`
	Dependencies []string               `json:"dependencies"`
	Conditions   []WorkflowCondition    `json:"conditions,omitempty"`
	Input        map[string]interface{} `json:"input,omitempty"`
	Output       map[string]interface{} `json:"output,omitempty"`
	Status       WorkflowNodeStatus     `json:"status"`
	StartTime    *time.Time             `json:"start_time,omitempty"`
	EndTime      *time.Time             `json:"end_time,omitempty"`
	Error        error                  `json:"error,omitempty"`
	RetryCount   int                    `json:"retry_count"`
	MaxRetries   int                    `json:"max_retries"`
}

// WorkflowNodeStatus represents the status of a workflow node
type WorkflowNodeStatus string

const (
	NodeStatusPending   WorkflowNodeStatus = "pending"
	NodeStatusReady     WorkflowNodeStatus = "ready"
	NodeStatusRunning   WorkflowNodeStatus = "running"
	NodeStatusCompleted WorkflowNodeStatus = "completed"
	NodeStatusFailed    WorkflowNodeStatus = "failed"
	NodeStatusSkipped   WorkflowNodeStatus = "skipped"
)

// WorkflowCondition defines when a node should execute
type WorkflowCondition struct {
	Type      ConditionType          `json:"type"`
	Field     string                 `json:"field"`
	Operator  ConditionOperator      `json:"operator"`
	Value     interface{}            `json:"value"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// ConditionType defines types of conditions
type ConditionType string

const (
	ConditionTypeContextVar ConditionType = "context_var"
	ConditionTypeNodeOutput ConditionType = "node_output"
	ConditionTypeNodeStatus ConditionType = "node_status"
	ConditionTypeCustom     ConditionType = "custom"
)

// ConditionOperator defines condition operators
type ConditionOperator string

const (
	OperatorEquals      ConditionOperator = "eq"
	OperatorNotEquals   ConditionOperator = "ne"
	OperatorGreaterThan ConditionOperator = "gt"
	OperatorLessThan    ConditionOperator = "lt"
	OperatorContains    ConditionOperator = "contains"
	OperatorExists      ConditionOperator = "exists"
	OperatorNotExists   ConditionOperator = "not_exists"
)

// Workflow represents a complete workflow definition
type Workflow struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Type            WorkflowType             `json:"type"`
	Nodes           map[string]*WorkflowNode `json:"nodes"`
	SupervisorAgent *Agent                   `json:"supervisor_agent,omitempty"`
	Context         map[string]interface{}   `json:"context"`
	Status          WorkflowStatus           `json:"status"`
	StartTime       *time.Time               `json:"start_time,omitempty"`
	EndTime         *time.Time               `json:"end_time,omitempty"`
	Error           error                    `json:"error,omitempty"`
	MaxConcurrency  int                      `json:"max_concurrency"`
	Timeout         time.Duration            `json:"timeout"`
}

// WorkflowStatus represents the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

// WorkflowExecutor executes workflows
type WorkflowExecutor struct {
	client        SwarmClient
	workflows     map[string]*Workflow
	mutex         sync.RWMutex
	logger        Logger
	eventHandlers map[WorkflowEventType][]WorkflowEventHandler
}

// WorkflowEventType represents different workflow events
type WorkflowEventType string

const (
	EventWorkflowStart    WorkflowEventType = "workflow_start"
	EventWorkflowComplete WorkflowEventType = "workflow_complete"
	EventWorkflowFailed   WorkflowEventType = "workflow_failed"
	EventNodeStart        WorkflowEventType = "node_start"
	EventNodeComplete     WorkflowEventType = "node_complete"
	EventNodeFailed       WorkflowEventType = "node_failed"
	EventNodeRetry        WorkflowEventType = "node_retry"
)

// WorkflowEventHandler handles workflow events
type WorkflowEventHandler func(event WorkflowEvent)

// WorkflowEvent represents a workflow event
type WorkflowEvent struct {
	Type       WorkflowEventType      `json:"type"`
	WorkflowID string                 `json:"workflow_id"`
	NodeID     string                 `json:"node_id,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Error      error                  `json:"error,omitempty"`
}

// NewWorkflowExecutor creates a new workflow executor
func NewWorkflowExecutor(client SwarmClient, logger Logger) *WorkflowExecutor {
	if logger == nil {
		logger = &defaultLogger{}
	}

	return &WorkflowExecutor{
		client:        client,
		workflows:     make(map[string]*Workflow),
		logger:        logger,
		eventHandlers: make(map[WorkflowEventType][]WorkflowEventHandler),
	}
}

// CreateWorkflow creates a new workflow
func (we *WorkflowExecutor) CreateWorkflow(name string, workflowType WorkflowType) *Workflow {
	workflow := &Workflow{
		ID:             fmt.Sprintf("workflow_%d", time.Now().UnixNano()),
		Name:           name,
		Type:           workflowType,
		Nodes:          make(map[string]*WorkflowNode),
		Context:        make(map[string]interface{}),
		Status:         WorkflowStatusPending,
		MaxConcurrency: 5,
		Timeout:        30 * time.Minute,
	}

	we.mutex.Lock()
	we.workflows[workflow.ID] = workflow
	we.mutex.Unlock()

	return workflow
}

// AddNode adds a node to the workflow
func (wf *Workflow) AddNode(id, name string, agent *Agent, dependencies []string) *WorkflowNode {
	node := &WorkflowNode{
		ID:           id,
		Name:         name,
		Agent:        agent,
		Dependencies: dependencies,
		Status:       NodeStatusPending,
		MaxRetries:   3,
		Input:        make(map[string]interface{}),
		Output:       make(map[string]interface{}),
	}

	wf.Nodes[id] = node
	return node
}

// AddCondition adds a condition to a workflow node
func (node *WorkflowNode) AddCondition(condType ConditionType, field string, operator ConditionOperator, value interface{}) {
	condition := WorkflowCondition{
		Type:     condType,
		Field:    field,
		Operator: operator,
		Value:    value,
	}
	node.Conditions = append(node.Conditions, condition)
}

// SetSupervisor sets a supervisor agent for the workflow
func (wf *Workflow) SetSupervisor(agent *Agent) {
	wf.SupervisorAgent = agent
}

// ExecuteWorkflow executes a workflow
func (we *WorkflowExecutor) ExecuteWorkflow(ctx context.Context, workflowID string, initialContext map[string]interface{}) (*WorkflowResult, error) {
	we.mutex.RLock()
	workflow, exists := we.workflows[workflowID]
	we.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("workflow %s not found", workflowID)
	}

	// Initialize workflow context
	if initialContext != nil {
		for k, v := range initialContext {
			workflow.Context[k] = v
		}
	}

	// Set workflow timeout
	if workflow.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, workflow.Timeout)
		defer cancel()
	}

	// Execute based on workflow type
	var result *WorkflowResult
	var err error

	workflow.Status = WorkflowStatusRunning
	startTime := time.Now()
	workflow.StartTime = &startTime

	we.emitEvent(EventWorkflowStart, workflowID, "", nil, nil)

	switch workflow.Type {
	case WorkflowTypeSequential:
		result, err = we.executeSequential(ctx, workflow)
	case WorkflowTypeParallel:
		result, err = we.executeParallel(ctx, workflow)
	case WorkflowTypeDAG:
		result, err = we.executeDAG(ctx, workflow)
	case WorkflowTypeSupervisor:
		result, err = we.executeSupervisor(ctx, workflow)
	case WorkflowTypePipeline:
		result, err = we.executePipeline(ctx, workflow)
	case WorkflowTypeConditional:
		result, err = we.executeConditional(ctx, workflow)
	default:
		err = fmt.Errorf("unsupported workflow type: %s", workflow.Type)
	}

	endTime := time.Now()
	workflow.EndTime = &endTime

	if err != nil {
		workflow.Status = WorkflowStatusFailed
		workflow.Error = err
		we.emitEvent(EventWorkflowFailed, workflowID, "", nil, err)
	} else {
		workflow.Status = WorkflowStatusCompleted
		we.emitEvent(EventWorkflowComplete, workflowID, "", nil, nil)
	}

	return result, err
}

// WorkflowResult represents the result of workflow execution
type WorkflowResult struct {
	WorkflowID    string                 `json:"workflow_id"`
	Status        WorkflowStatus         `json:"status"`
	ExecutionTime time.Duration          `json:"execution_time"`
	NodeResults   map[string]*NodeResult `json:"node_results"`
	FinalContext  map[string]interface{} `json:"final_context"`
	Error         error                  `json:"error,omitempty"`
}

// NodeResult represents the result of a node execution
type NodeResult struct {
	NodeID        string                 `json:"node_id"`
	Status        WorkflowNodeStatus     `json:"status"`
	Output        map[string]interface{} `json:"output"`
	ExecutionTime time.Duration          `json:"execution_time"`
	RetryCount    int                    `json:"retry_count"`
	Error         error                  `json:"error,omitempty"`
}

// OnEvent registers an event handler
func (we *WorkflowExecutor) OnEvent(eventType WorkflowEventType, handler WorkflowEventHandler) {
	we.eventHandlers[eventType] = append(we.eventHandlers[eventType], handler)
}

// emitEvent emits a workflow event
func (we *WorkflowExecutor) emitEvent(eventType WorkflowEventType, workflowID, nodeID string, data map[string]interface{}, err error) {
	event := WorkflowEvent{
		Type:       eventType,
		WorkflowID: workflowID,
		NodeID:     nodeID,
		Data:       data,
		Timestamp:  time.Now(),
		Error:      err,
	}

	if handlers, exists := we.eventHandlers[eventType]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}
