package swarm

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// executeSequential executes nodes one after another in order
func (we *WorkflowExecutor) executeSequential(ctx context.Context, workflow *Workflow) (*WorkflowResult, error) {
	result := &WorkflowResult{
		WorkflowID:   workflow.ID,
		NodeResults:  make(map[string]*NodeResult),
		FinalContext: make(map[string]interface{}),
	}

	// Sort nodes by dependencies (simple topological sort for sequential)
	orderedNodes := we.getExecutionOrder(workflow)

	for _, nodeID := range orderedNodes {
		node := workflow.Nodes[nodeID]

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		// Check conditions
		if !we.evaluateConditions(node, workflow.Context, result.NodeResults) {
			node.Status = NodeStatusSkipped
			continue
		}

		nodeResult := we.executeNode(ctx, node, workflow.Context)
		result.NodeResults[nodeID] = nodeResult

		// Update workflow context with node output
		for k, v := range nodeResult.Output {
			workflow.Context[k] = v
		}

		if nodeResult.Error != nil && nodeResult.Status == NodeStatusFailed {
			result.Status = WorkflowStatusFailed
			result.Error = nodeResult.Error
			return result, nodeResult.Error
		}
	}

	// Copy final context
	for k, v := range workflow.Context {
		result.FinalContext[k] = v
	}

	result.Status = WorkflowStatusCompleted
	return result, nil
}

// executeParallel executes all nodes concurrently
func (we *WorkflowExecutor) executeParallel(ctx context.Context, workflow *Workflow) (*WorkflowResult, error) {
	result := &WorkflowResult{
		WorkflowID:   workflow.ID,
		NodeResults:  make(map[string]*NodeResult),
		FinalContext: make(map[string]interface{}),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errorsChan := make(chan error, len(workflow.Nodes))

	// Semaphore to limit concurrency
	semaphore := make(chan struct{}, workflow.MaxConcurrency)

	for nodeID, node := range workflow.Nodes {
		wg.Add(1)
		go func(nodeID string, node *WorkflowNode) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Check conditions
			if !we.evaluateConditions(node, workflow.Context, result.NodeResults) {
				node.Status = NodeStatusSkipped
				return
			}

			nodeResult := we.executeNode(ctx, node, workflow.Context)

			mu.Lock()
			result.NodeResults[nodeID] = nodeResult
			// Update shared context with node output
			for k, v := range nodeResult.Output {
				workflow.Context[k] = v
			}
			mu.Unlock()

			if nodeResult.Error != nil {
				errorsChan <- nodeResult.Error
			}
		}(nodeID, node)
	}

	wg.Wait()
	close(errorsChan)

	// Check for errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		result.Status = WorkflowStatusFailed
		result.Error = fmt.Errorf("multiple node failures: %v", errors)
		return result, result.Error
	}

	// Copy final context
	for k, v := range workflow.Context {
		result.FinalContext[k] = v
	}

	result.Status = WorkflowStatusCompleted
	return result, nil
}

// executeDAG executes nodes based on dependency graph
func (we *WorkflowExecutor) executeDAG(ctx context.Context, workflow *Workflow) (*WorkflowResult, error) {
	result := &WorkflowResult{
		WorkflowID:   workflow.ID,
		NodeResults:  make(map[string]*NodeResult),
		FinalContext: make(map[string]interface{}),
	}

	// Build dependency graph
	dependencyGraph := make(map[string][]string)
	inDegree := make(map[string]int)

	for nodeID, node := range workflow.Nodes {
		inDegree[nodeID] = len(node.Dependencies)
		for _, dep := range node.Dependencies {
			dependencyGraph[dep] = append(dependencyGraph[dep], nodeID)
		}
	}

	// Find nodes with no dependencies (ready to execute)
	readyQueue := make([]string, 0)
	for nodeID, degree := range inDegree {
		if degree == 0 {
			readyQueue = append(readyQueue, nodeID)
			workflow.Nodes[nodeID].Status = NodeStatusReady
		}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, workflow.MaxConcurrency)
	hasError := false

	for len(readyQueue) > 0 && !hasError {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		// Process all ready nodes in parallel
		currentBatch := make([]string, len(readyQueue))
		copy(currentBatch, readyQueue)
		readyQueue = readyQueue[:0] // Clear ready queue

		for _, nodeID := range currentBatch {
			wg.Add(1)
			go func(nodeID string) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				node := workflow.Nodes[nodeID]

				// Check conditions
				if !we.evaluateConditions(node, workflow.Context, result.NodeResults) {
					node.Status = NodeStatusSkipped
					mu.Lock()
					we.updateDependents(nodeID, dependencyGraph, inDegree, &readyQueue, workflow)
					mu.Unlock()
					return
				}

				nodeResult := we.executeNode(ctx, node, workflow.Context)

				mu.Lock()
				result.NodeResults[nodeID] = nodeResult

				if nodeResult.Error != nil {
					hasError = true
					result.Status = WorkflowStatusFailed
					result.Error = nodeResult.Error
				} else {
					// Update context with node output
					for k, v := range nodeResult.Output {
						workflow.Context[k] = v
					}
					// Update dependents
					we.updateDependents(nodeID, dependencyGraph, inDegree, &readyQueue, workflow)
				}
				mu.Unlock()
			}(nodeID)
		}

		wg.Wait()
	}

	if !hasError {
		// Copy final context
		for k, v := range workflow.Context {
			result.FinalContext[k] = v
		}
		result.Status = WorkflowStatusCompleted
	}

	return result, result.Error
}

// executeSupervisor executes workflow with supervisor oversight
func (we *WorkflowExecutor) executeSupervisor(ctx context.Context, workflow *Workflow) (*WorkflowResult, error) {
	result := &WorkflowResult{
		WorkflowID:   workflow.ID,
		NodeResults:  make(map[string]*NodeResult),
		FinalContext: make(map[string]interface{}),
	}

	if workflow.SupervisorAgent == nil {
		return nil, fmt.Errorf("supervisor workflow requires a supervisor agent")
	}

	// Supervisor analyzes the workflow and creates execution plan
	supervisorMessages := []Message{
		{
			Role:    "user",
			Content: fmt.Sprintf("Analyze workflow '%s' with %d nodes and create execution plan", workflow.Name, len(workflow.Nodes)),
		},
	}

	supervisorContext := map[string]interface{}{
		"workflow_id":   workflow.ID,
		"workflow_type": "supervisor",
		"total_nodes":   len(workflow.Nodes),
		"node_names":    we.getNodeNames(workflow),
	}

	// Run supervisor planning
	supervisorResponse := we.client.RunWithContext(ctx, workflow.SupervisorAgent, supervisorMessages, supervisorContext)
	if supervisorResponse.Error != nil {
		return result, fmt.Errorf("supervisor planning failed: %w", supervisorResponse.Error)
	}

	// Execute nodes based on supervisor's plan (for now, use sequential)
	orderedNodes := we.getExecutionOrder(workflow)

	for _, nodeID := range orderedNodes {
		node := workflow.Nodes[nodeID]

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		// Supervisor monitors each node execution
		nodeResult := we.executeNode(ctx, node, workflow.Context)
		result.NodeResults[nodeID] = nodeResult

		// Supervisor reviews result
		reviewMessages := []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Review execution of node '%s': status=%s", node.Name, nodeResult.Status),
			},
		}

		supervisorReview := we.client.RunWithContext(ctx, workflow.SupervisorAgent, reviewMessages, supervisorContext)
		if supervisorReview.Error == nil {
			// Update context based on supervisor review
			for k, v := range supervisorReview.ContextVars {
				workflow.Context[k] = v
			}
		}

		// Update workflow context with node output
		for k, v := range nodeResult.Output {
			workflow.Context[k] = v
		}

		if nodeResult.Error != nil && nodeResult.Status == NodeStatusFailed {
			result.Status = WorkflowStatusFailed
			result.Error = nodeResult.Error
			return result, nodeResult.Error
		}
	}

	// Copy final context
	for k, v := range workflow.Context {
		result.FinalContext[k] = v
	}

	result.Status = WorkflowStatusCompleted
	return result, nil
}

// executePipeline executes nodes in a data pipeline pattern
func (we *WorkflowExecutor) executePipeline(ctx context.Context, workflow *Workflow) (*WorkflowResult, error) {
	result := &WorkflowResult{
		WorkflowID:   workflow.ID,
		NodeResults:  make(map[string]*NodeResult),
		FinalContext: make(map[string]interface{}),
	}

	// Pipeline passes data from one stage to the next
	orderedNodes := we.getExecutionOrder(workflow)
	pipelineData := make(map[string]interface{})

	// Initialize pipeline with workflow context
	for k, v := range workflow.Context {
		pipelineData[k] = v
	}

	for _, nodeID := range orderedNodes {
		node := workflow.Nodes[nodeID]

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		// Check conditions
		if !we.evaluateConditions(node, pipelineData, result.NodeResults) {
			node.Status = NodeStatusSkipped
			continue
		}

		// Execute node with pipeline data
		nodeResult := we.executeNode(ctx, node, pipelineData)
		result.NodeResults[nodeID] = nodeResult

		if nodeResult.Error != nil && nodeResult.Status == NodeStatusFailed {
			result.Status = WorkflowStatusFailed
			result.Error = nodeResult.Error
			return result, nodeResult.Error
		}

		// Pipeline: output becomes input for next stage
		pipelineData = nodeResult.Output

		// Also update workflow context
		for k, v := range nodeResult.Output {
			workflow.Context[k] = v
		}
	}

	// Copy final context
	for k, v := range workflow.Context {
		result.FinalContext[k] = v
	}

	result.Status = WorkflowStatusCompleted
	return result, nil
}

// executeConditional executes nodes based on conditional logic
func (we *WorkflowExecutor) executeConditional(ctx context.Context, workflow *Workflow) (*WorkflowResult, error) {
	result := &WorkflowResult{
		WorkflowID:   workflow.ID,
		NodeResults:  make(map[string]*NodeResult),
		FinalContext: make(map[string]interface{}),
	}

	// Build conditional execution tree
	executedNodes := make(map[string]bool)

	// Start with nodes that have no dependencies
	candidateNodes := make([]string, 0)
	for nodeID, node := range workflow.Nodes {
		if len(node.Dependencies) == 0 {
			candidateNodes = append(candidateNodes, nodeID)
		}
	}

	for len(candidateNodes) > 0 {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		nextCandidates := make([]string, 0)

		for _, nodeID := range candidateNodes {
			node := workflow.Nodes[nodeID]

			// Check if all dependencies are satisfied
			allDepsSatisfied := true
			for _, dep := range node.Dependencies {
				if !executedNodes[dep] {
					allDepsSatisfied = false
					break
				}
			}

			if !allDepsSatisfied {
				nextCandidates = append(nextCandidates, nodeID)
				continue
			}

			// Evaluate conditions
			if !we.evaluateConditions(node, workflow.Context, result.NodeResults) {
				node.Status = NodeStatusSkipped
				executedNodes[nodeID] = true

				// Add dependent nodes to candidates
				for dependentID, dependentNode := range workflow.Nodes {
					for _, dep := range dependentNode.Dependencies {
						if dep == nodeID && !executedNodes[dependentID] {
							nextCandidates = append(nextCandidates, dependentID)
						}
					}
				}
				continue
			}

			// Execute node
			nodeResult := we.executeNode(ctx, node, workflow.Context)
			result.NodeResults[nodeID] = nodeResult
			executedNodes[nodeID] = true

			// Update context
			for k, v := range nodeResult.Output {
				workflow.Context[k] = v
			}

			if nodeResult.Error != nil && nodeResult.Status == NodeStatusFailed {
				result.Status = WorkflowStatusFailed
				result.Error = nodeResult.Error
				return result, nodeResult.Error
			}

			// Add dependent nodes to candidates
			for dependentID, dependentNode := range workflow.Nodes {
				for _, dep := range dependentNode.Dependencies {
					if dep == nodeID && !executedNodes[dependentID] {
						nextCandidates = append(nextCandidates, dependentID)
					}
				}
			}
		}

		candidateNodes = nextCandidates
	}

	// Copy final context
	for k, v := range workflow.Context {
		result.FinalContext[k] = v
	}

	result.Status = WorkflowStatusCompleted
	return result, nil
}

// Helper methods

func (we *WorkflowExecutor) executeNode(ctx context.Context, node *WorkflowNode, workflowContext map[string]interface{}) *NodeResult {
	nodeResult := &NodeResult{
		NodeID: node.ID,
		Output: make(map[string]interface{}),
	}

	startTime := time.Now()
	node.StartTime = &startTime
	node.Status = NodeStatusRunning

	we.emitEvent(EventNodeStart, "", node.ID, nil, nil)

	// Retry logic
	for attempt := 0; attempt <= node.MaxRetries; attempt++ {
		if attempt > 0 {
			node.RetryCount++
			we.emitEvent(EventNodeRetry, "", node.ID, map[string]interface{}{
				"attempt":     attempt,
				"max_retries": node.MaxRetries,
			}, nil)
		}

		// Prepare messages for agent execution
		messages := []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Execute task: %s", node.Name),
			},
		}

		// Merge node input with workflow context
		contextVars := make(map[string]interface{})
		for k, v := range workflowContext {
			contextVars[k] = v
		}
		for k, v := range node.Input {
			contextVars[k] = v
		}

		// Execute node with agent
		response := we.client.RunWithContext(ctx, node.Agent, messages, contextVars)

		if response.Error == nil {
			// Success
			node.Status = NodeStatusCompleted
			nodeResult.Status = NodeStatusCompleted
			nodeResult.Output = response.ContextVars
			nodeResult.RetryCount = node.RetryCount

			endTime := time.Now()
			node.EndTime = &endTime
			nodeResult.ExecutionTime = endTime.Sub(startTime)

			we.emitEvent(EventNodeComplete, "", node.ID, nil, nil)
			break
		} else {
			// Failure
			if attempt == node.MaxRetries {
				node.Status = NodeStatusFailed
				node.Error = response.Error
				nodeResult.Status = NodeStatusFailed
				nodeResult.Error = response.Error
				nodeResult.RetryCount = node.RetryCount

				endTime := time.Now()
				node.EndTime = &endTime
				nodeResult.ExecutionTime = endTime.Sub(startTime)

				we.emitEvent(EventNodeFailed, "", node.ID, nil, response.Error)
			}
		}
	}

	return nodeResult
}

func (we *WorkflowExecutor) evaluateConditions(node *WorkflowNode, context map[string]interface{}, nodeResults map[string]*NodeResult) bool {
	if len(node.Conditions) == 0 {
		return true
	}

	for _, condition := range node.Conditions {
		if !we.evaluateCondition(condition, context, nodeResults) {
			return false
		}
	}

	return true
}

func (we *WorkflowExecutor) evaluateCondition(condition WorkflowCondition, context map[string]interface{}, nodeResults map[string]*NodeResult) bool {
	var actualValue interface{}

	switch condition.Type {
	case ConditionTypeContextVar:
		actualValue = context[condition.Field]
	case ConditionTypeNodeOutput:
		// Field format: "node_id.output_key"
		// Simple implementation for now
		actualValue = context[condition.Field]
	case ConditionTypeNodeStatus:
		if result, exists := nodeResults[condition.Field]; exists {
			actualValue = string(result.Status)
		}
	default:
		return false
	}

	return we.compareValues(actualValue, condition.Operator, condition.Value)
}

func (we *WorkflowExecutor) compareValues(actual interface{}, operator ConditionOperator, expected interface{}) bool {
	switch operator {
	case OperatorEquals:
		return actual == expected
	case OperatorNotEquals:
		return actual != expected
	case OperatorExists:
		return actual != nil
	case OperatorNotExists:
		return actual == nil
	case OperatorContains:
		if actualStr, ok := actual.(string); ok {
			if expectedStr, ok := expected.(string); ok {
				return fmt.Sprintf("%v", actualStr) == expectedStr
			}
		}
		return false
	default:
		return false
	}
}

func (we *WorkflowExecutor) getExecutionOrder(workflow *Workflow) []string {
	// Simple topological sort
	visited := make(map[string]bool)
	stack := make([]string, 0)

	var visit func(string)
	visit = func(nodeID string) {
		if visited[nodeID] {
			return
		}
		visited[nodeID] = true

		node := workflow.Nodes[nodeID]
		for _, dep := range node.Dependencies {
			visit(dep)
		}
		stack = append(stack, nodeID)
	}

	for nodeID := range workflow.Nodes {
		visit(nodeID)
	}

	return stack
}

func (we *WorkflowExecutor) updateDependents(completedNodeID string, dependencyGraph map[string][]string, inDegree map[string]int, readyQueue *[]string, workflow *Workflow) {
	for _, dependent := range dependencyGraph[completedNodeID] {
		inDegree[dependent]--
		if inDegree[dependent] == 0 {
			*readyQueue = append(*readyQueue, dependent)
			workflow.Nodes[dependent].Status = NodeStatusReady
		}
	}
}

func (we *WorkflowExecutor) getNodeNames(workflow *Workflow) []string {
	names := make([]string, 0, len(workflow.Nodes))
	for _, node := range workflow.Nodes {
		names = append(names, node.Name)
	}
	sort.Strings(names)
	return names
}
