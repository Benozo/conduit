package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/benozo/conduit/mcp"
)

// LLMAgent represents an agent that uses LLM for reasoning and decision-making
type LLMAgent struct {
	*Agent
	modelFunc mcp.ModelFunc
	modelName string
}

// LLMAgentManager extends MCPAgentManager with LLM capabilities
type LLMAgentManager struct {
	*MCPAgentManager
	modelFunc mcp.ModelFunc
	modelName string
}

// NewLLMAgentManager creates a new agent manager with LLM integration
func NewLLMAgentManager(mcpServer interface {
	GetToolRegistry() *mcp.ToolRegistry
	GetMemory() *mcp.Memory
}, modelFunc mcp.ModelFunc, modelName string) *LLMAgentManager {
	return &LLMAgentManager{
		MCPAgentManager: NewMCPAgentManager(mcpServer),
		modelFunc:       modelFunc,
		modelName:       modelName,
	}
}

// CreateLLMAgent creates an agent that uses LLM for reasoning
func (lam *LLMAgentManager) CreateLLMAgent(id, name, description, systemPrompt string, tools []string, config *AgentConfig) (*LLMAgent, error) {
	agent, err := lam.CreateAgent(id, name, description, systemPrompt, tools, config)
	if err != nil {
		return nil, err
	}

	return &LLMAgent{
		Agent:     agent,
		modelFunc: lam.modelFunc,
		modelName: lam.modelName,
	}, nil
}

// ExecuteTaskWithLLM executes a task using LLM reasoning
func (lam *LLMAgentManager) ExecuteTaskWithLLM(taskID string) error {
	task, exists := lam.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	agent, err := lam.GetAgent(task.AgentID)
	if err != nil {
		return err
	}

	// Update task status
	task.Status = TaskStatusRunning
	task.Progress = 0.0
	now := time.Now()
	task.StartedAt = &now
	agent.State = StateThinking

	// Create execution context
	execCtx := &ExecutionContext{
		TaskID:    task.ID,
		AgentID:   agent.ID,
		SessionID: fmt.Sprintf("session_%d", time.Now().UnixNano()),
		Context:   lam.ctx,
		Memory:    agent.Memory,
		Logger:    &defaultLogger{},
	}

	// Execute task with LLM reasoning
	err = lam.executeTaskWithLLMReasoning(execCtx, task, agent)
	if err != nil {
		task.Status = TaskStatusFailed
		task.Error = err.Error()
		agent.State = StateError
		return err
	}

	// Mark task as completed
	task.Status = TaskStatusCompleted
	task.Progress = 1.0
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	agent.State = StateCompleted

	return nil
}

// executeTaskWithLLMReasoning uses LLM to reason about and execute tasks
func (lam *LLMAgentManager) executeTaskWithLLMReasoning(execCtx *ExecutionContext, task *Task, agent *Agent) error {
	// Step 1: LLM Analysis and Planning
	reasoningStep := TaskStep{
		ID:          fmt.Sprintf("step_%d", time.Now().UnixNano()),
		Name:        "llm_reasoning",
		Description: "LLM analyzes task and creates execution plan",
		Input:       task.Input,
		Status:      TaskStatusRunning,
		StartedAt:   timePtr(time.Now()),
	}
	task.Steps = append(task.Steps, reasoningStep)

	execCtx.Logger.Info("LLM is analyzing the task", "task_id", task.ID, "agent_id", agent.ID)

	// Create reasoning prompt for the LLM
	reasoningPrompt := lam.createReasoningPrompt(task, agent)

	// Get LLM analysis
	analysis, err := lam.getLLMAnalysis(reasoningPrompt, execCtx)
	if err != nil {
		reasoningStep.Status = TaskStatusFailed
		reasoningStep.Error = err.Error()
		return fmt.Errorf("LLM reasoning failed: %w", err)
	}

	//	log.Printf("LLM Analysis: %s", analysis)

	// Parse LLM response to extract action plan
	actionPlan, err := lam.parseLLMActionPlan(analysis)
	if err != nil {
		reasoningStep.Status = TaskStatusFailed
		reasoningStep.Error = err.Error()
		return fmt.Errorf("failed to parse LLM action plan: %w", err)
	}

	reasoningStep.Output = map[string]interface{}{
		"llm_analysis": analysis,
		"action_plan":  actionPlan,
		"reasoning":    "LLM analyzed task and created execution plan",
	}
	reasoningStep.Status = TaskStatusCompleted
	reasoningStep.CompletedAt = timePtr(time.Now())

	// Step 2: Execute LLM-planned actions
	agent.State = StateActing
	for i, action := range actionPlan {
		// Resolve step dependencies before execution
		resolvedAction := lam.resolveStepDependencies(action, task.Steps)

		stepID := fmt.Sprintf("step_%d_%d", time.Now().UnixNano(), i)
		actionStep := TaskStep{
			ID:          stepID,
			Name:        resolvedAction.Name,
			Description: resolvedAction.Description,
			Input:       resolvedAction.Input,
			Status:      TaskStatusRunning,
			StartedAt:   timePtr(time.Now()),
		}

		task.Steps = append(task.Steps, actionStep)

		// Execute the resolved action using MCP tools
		result, err := lam.executeAction(execCtx, resolvedAction, agent)
		if err != nil {
			actionStep.Status = TaskStatusFailed
			actionStep.Error = err.Error()

			// Ask LLM for error recovery strategy
			if lam.shouldRetryWithLLM(err, resolvedAction) {
				recoveryAction, recoveryErr := lam.getLLMErrorRecovery(err, resolvedAction, execCtx)
				if recoveryErr == nil {
					// Resolve dependencies for recovery action too
					resolvedRecoveryAction := lam.resolveStepDependencies(recoveryAction, task.Steps)
					// Try recovery action
					result, err = lam.executeAction(execCtx, resolvedRecoveryAction, agent)
					if err == nil {
						actionStep.Status = TaskStatusCompleted
						actionStep.Output = result
						actionStep.CompletedAt = timePtr(time.Now())
						task.Progress = float64(i+1) / float64(len(actionPlan))
						continue
					}
				}
			}

			return err
		}

		actionStep.Output = result
		actionStep.Status = TaskStatusCompleted
		actionStep.CompletedAt = timePtr(time.Now())

		// Update task progress
		task.Progress = float64(i+1) / float64(len(actionPlan))
	}

	return nil
}

// createReasoningPrompt creates a prompt for LLM reasoning
func (lam *LLMAgentManager) createReasoningPrompt(task *Task, agent *Agent) string {
	availableTools := strings.Join(agent.Tools, ", ")

	// Check if this is a modification task
	var promptTemplate string
	if input, ok := task.Input["action"]; ok && input == "modify_existing" {
		promptTemplate = `%s

MODIFICATION TASK: %s
DESCRIPTION: %s
INPUT: %s

This is a MODIFICATION task. You must:
1. Use read_existing_html to get the current HTML content
2. Analyze the requested changes
3. Create updated HTML with create_html_page

AVAILABLE TOOLS: %s

IMPORTANT: Respond ONLY with valid JSON. No explanations, no thinking tags, no markdown.

Required JSON format:
{
  "analysis": "your reasoning about the modifications needed",
  "steps": [
    {
      "name": "read_current",
      "description": "Read existing HTML file",
      "tool": "read_existing_html",
      "input": {"filename": "FILENAME"}
    },
    {
      "name": "create_updated",
      "description": "Create updated HTML with modifications",
      "tool": "create_html_page",
      "input": {"filename": "FILENAME", "content": "COMPLETE_UPDATED_HTML"}
    }
  ],
  "reasoning": "why these modifications improve the page"
}`
	} else {
		promptTemplate = `%s

TASK: %s
DESCRIPTION: %s
INPUT: %s

AVAILABLE TOOLS: %s

IMPORTANT: Respond ONLY with valid JSON. No explanations, no thinking tags, no markdown.

Required JSON format:
{
  "analysis": "your reasoning about the task",
  "steps": [
    {
      "name": "step_name",
      "description": "what this step does",
      "tool": "tool_to_use",
      "input": {"param1": "value1", "param2": "value2"}
    }
  ],
  "reasoning": "why you chose this approach"
}`
	}

	prompt := fmt.Sprintf(promptTemplate,
		agent.SystemPrompt,
		task.Title, task.Description, formatInput(task.Input),
		availableTools)

	return prompt
}

// getLLMAnalysis gets analysis from the LLM
func (lam *LLMAgentManager) getLLMAnalysis(prompt string, execCtx *ExecutionContext) (string, error) {
	// Create context input for LLM
	ctx := mcp.ContextInput{
		ContextID: execCtx.SessionID,
		Inputs: map[string]interface{}{
			"query": prompt,
		},
	}

	// Create MCP request
	req := mcp.MCPRequest{
		SessionID:   execCtx.SessionID,
		Contexts:    []mcp.ContextInput{ctx},
		Model:       lam.modelName,
		Temperature: 0.3, // Lower temperature for more focused reasoning
		Stream:      false,
	}

	// Get LLM response
	response, err := lam.modelFunc(ctx, req, execCtx.Memory, nil)
	if err != nil {
		return "", fmt.Errorf("LLM request failed: %w", err)
	}

	return response, nil
}

// parseLLMActionPlan parses the LLM response to extract action plan
func (lam *LLMAgentManager) parseLLMActionPlan(analysis string) ([]Action, error) {
	log.Printf("🔍 Parsing LLM response (first 500 chars): %s", truncateString(analysis, 500))

	// Clean up the response - remove thinking tags and other artifacts
	cleanedResponse := analysis

	// Remove <think>...</think> blocks
	if thinkStart := strings.Index(cleanedResponse, "<think>"); thinkStart != -1 {
		if thinkEnd := strings.Index(cleanedResponse, "</think>"); thinkEnd != -1 {
			cleanedResponse = cleanedResponse[:thinkStart] + cleanedResponse[thinkEnd+8:]
		}
	}

	// Remove markdown code blocks
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "```json", "")
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "```", "")

	// Trim whitespace
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	log.Printf("🔍 Cleaned response: %s", truncateString(cleanedResponse, 300))

	// Find the start of JSON
	jsonStart := strings.Index(cleanedResponse, "{")
	if jsonStart == -1 {
		log.Printf("❌ No JSON start found in response")
		return lam.createFallbackActionPlan(analysis)
	}

	// Extract JSON using brace counting for proper nesting
	jsonStr := lam.extractCompleteJSON(cleanedResponse[jsonStart:])
	if jsonStr == "" {
		log.Printf("❌ Failed to extract complete JSON")
		return lam.createFallbackActionPlan(analysis)
	}

	log.Printf("🔍 Extracted JSON length: %d chars", len(jsonStr))

	var llmResponse struct {
		Analysis string `json:"analysis"`
		Steps    []struct {
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			Tool        string                 `json:"tool"`
			Input       map[string]interface{} `json:"input"`
		} `json:"steps"`
		Reasoning string `json:"reasoning"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &llmResponse); err != nil {
		log.Printf("❌ JSON parsing failed: %v", err)
		// Fallback if JSON parsing fails
		return lam.createFallbackActionPlan(analysis)
	}

	log.Printf("✅ Successfully parsed %d steps from LLM response", len(llmResponse.Steps))

	// Convert to Action structs
	var actions []Action
	for _, step := range llmResponse.Steps {
		actions = append(actions, Action{
			Name:        step.Name,
			Description: step.Description,
			Tool:        step.Tool,
			Input:       step.Input,
		})
	}

	return actions, nil
}

// extractCompleteJSON extracts a complete JSON object using brace counting
func (lam *LLMAgentManager) extractCompleteJSON(text string) string {
	if !strings.HasPrefix(text, "{") {
		return ""
	}

	// Limit maximum JSON size to prevent memory issues
	maxJSONSize := 50000 // 50KB limit
	if len(text) > maxJSONSize {
		log.Printf("⚠️ Response too large (%d chars), truncating for JSON extraction", len(text))
		text = text[:maxJSONSize]
	}

	braceCount := 0
	inString := false
	escape := false

	for i, char := range text {
		if escape {
			escape = false
			continue
		}

		if char == '\\' && inString {
			escape = true
			continue
		}

		if char == '"' {
			inString = !inString
			continue
		}

		if !inString {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 {
					result := text[:i+1]
					log.Printf("🔍 Successfully extracted JSON: %d chars", len(result))
					return result
				}
			}
		}
	}

	log.Printf("❌ Could not find complete JSON in text")
	return ""
}

// createFallbackActionPlan creates a simple action plan when LLM parsing fails
func (lam *LLMAgentManager) createFallbackActionPlan(analysis string) ([]Action, error) {
	log.Printf("🔄 Creating intelligent fallback action plan")

	// Try to detect if this looks like an HTML creation task
	if strings.Contains(analysis, "html") || strings.Contains(analysis, "HTML") ||
		strings.Contains(analysis, "landing") || strings.Contains(analysis, "page") {

		// Look for what might be HTML content in the response
		if htmlStart := strings.Index(analysis, "<!DOCTYPE"); htmlStart != -1 {
			// Try to extract HTML content
			htmlContent := analysis[htmlStart:]

			// Find a reasonable end point (look for </html> or take a reasonable chunk)
			if htmlEnd := strings.Index(htmlContent, "</html>"); htmlEnd != -1 {
				htmlContent = htmlContent[:htmlEnd+7]

				log.Printf("🔍 Found HTML content in response, creating fallback HTML action")
				actions := []Action{
					{
						Name:        "create_html_fallback",
						Description: "Create HTML page from extracted content",
						Tool:        "create_html_page",
						Input: map[string]interface{}{
							"filename": "landing.html",
							"content":  htmlContent,
						},
					},
				}
				return actions, nil
			}
		}

		// If no HTML found, create a simple HTML template
		log.Printf("🔍 Creating simple HTML template as fallback")
		simpleHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Landing Page</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50">
    <div class="container mx-auto px-4 py-16 text-center">
        <h1 class="text-4xl font-bold mb-4">Welcome</h1>
        <p class="text-lg text-gray-600 mb-8">This page was created as a fallback when JSON parsing failed.</p>
        <a href="#" class="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition">Get Started</a>
    </div>
</body>
</html>`

		actions := []Action{
			{
				Name:        "create_simple_html",
				Description: "Create simple HTML template as fallback",
				Tool:        "create_html_page",
				Input: map[string]interface{}{
					"filename": "landing.html",
					"content":  simpleHTML,
				},
			},
		}
		return actions, nil
	}

	// Generic fallback for non-HTML tasks
	actions := []Action{
		{
			Name:        "fallback_analysis",
			Description: "Fallback analysis when LLM response parsing fails",
			Tool:        "timestamp",
			Input:       map[string]interface{}{},
		},
	}

	log.Printf("Using fallback action plan due to LLM parsing failure")
	return actions, nil
}

// shouldRetryWithLLM determines if we should ask LLM for error recovery
func (lam *LLMAgentManager) shouldRetryWithLLM(err error, action Action) bool {
	// Retry for tool-related errors but not for system errors
	return strings.Contains(err.Error(), "tool") || strings.Contains(err.Error(), "parameter")
}

// getLLMErrorRecovery asks LLM to suggest error recovery
func (lam *LLMAgentManager) getLLMErrorRecovery(err error, failedAction Action, execCtx *ExecutionContext) (Action, error) {
	prompt := fmt.Sprintf(`The following action failed:
Action: %s
Tool: %s
Input: %s
Error: %s

Please suggest a corrected action. Respond with JSON:
{
  "name": "corrected_action_name",
  "description": "what the corrected action does",
  "tool": "tool_to_use",
  "input": {"param1": "value1"}
}`,
		failedAction.Name, failedAction.Tool, formatInput(failedAction.Input), err.Error())

	response, err := lam.getLLMAnalysis(prompt, execCtx)
	if err != nil {
		return Action{}, err
	}

	// Parse recovery action
	actions, err := lam.parseLLMActionPlan(response)
	if err != nil || len(actions) == 0 {
		return Action{}, fmt.Errorf("failed to parse recovery action")
	}

	return actions[0], nil
}

// resolveStepDependencies resolves references to previous step results in action inputs
func (lam *LLMAgentManager) resolveStepDependencies(action Action, completedSteps []TaskStep) Action {
	// Create a copy of the action to modify
	resolvedAction := Action{
		Name:        action.Name,
		Description: action.Description,
		Tool:        action.Tool,
		Input:       make(map[string]interface{}),
	}

	// Process each input parameter
	for key, value := range action.Input {
		resolvedValue := lam.resolveParameterValue(value, completedSteps)
		resolvedAction.Input[key] = resolvedValue
	}

	return resolvedAction
}

// resolveParameterValue resolves a parameter value, substituting step results if needed
func (lam *LLMAgentManager) resolveParameterValue(value interface{}, completedSteps []TaskStep) interface{} {
	switch v := value.(type) {
	case string:
		// Check if this is a reference to a previous step result
		if resolved := lam.resolveStepReference(v, completedSteps); resolved != nil {
			return resolved
		}
		return v
	default:
		return v
	}
}

// resolveStepReference attempts to resolve step references like "result of step 1"
func (lam *LLMAgentManager) resolveStepReference(reference string, completedSteps []TaskStep) interface{} {
	// Skip llm_reasoning step when counting steps for dependencies
	var actionSteps []TaskStep
	for _, step := range completedSteps {
		if step.Name != "llm_reasoning" && step.Status == TaskStatusCompleted {
			actionSteps = append(actionSteps, step)
		}
	}

	// Common patterns for step references
	patterns := map[string]int{
		"result of step 1": 0,
		"result of step 2": 1,
		"result of step 3": 2,
		"step 1 result":    0,
		"step 2 result":    1,
		"step 3 result":    2,
		"previous result":  len(actionSteps) - 1,
		"last result":      len(actionSteps) - 1,
	}

	if stepIndex, exists := patterns[reference]; exists {
		// Make sure the step index is valid
		if stepIndex >= 0 && stepIndex < len(actionSteps) {
			step := actionSteps[stepIndex]

			// Extract numeric result if available
			if step.Output != nil {
				if output, ok := step.Output["output"]; ok {
					if outputMap, ok := output.(map[string]interface{}); ok {
						if result, ok := outputMap["result"]; ok {
							return result
						}
					}
				}

				// Also try direct result in output
				if result, ok := step.Output["result"]; ok {
					return result
				}
			}
		}
	}

	return nil
}

// formatInput formats input parameters for display
func formatInput(input map[string]interface{}) string {
	if len(input) == 0 {
		return "{}"
	}

	parts := []string{}
	for k, v := range input {
		parts = append(parts, fmt.Sprintf("%s: %v", k, v))
	}
	return strings.Join(parts, ", ")
}

// truncateString truncates a string to maxLength with ellipsis
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
