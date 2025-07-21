package swarm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/mcp"
)

// swarmClient implements the SwarmClient interface
type swarmClient struct {
	mcpServer    interface{} // MCP server interface
	toolRegistry *mcp.ToolRegistry
	memory       *mcp.Memory
	config       *SwarmConfig
	functions    map[string]AgentFunction
	agents       map[string]*Agent
	logger       Logger
	modelFunc    mcp.ModelFunc // LLM model function for agent reasoning
	modelName    string        // Model name for LLM requests
}

// defaultLogger provides basic logging functionality
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}
func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}
func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}
func (l *defaultLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

// NewSwarmClient creates a new swarm client with MCP integration
func NewSwarmClient(mcpServer interface {
	GetToolRegistry() *mcp.ToolRegistry
	GetMemory() *mcp.Memory
}, config *SwarmConfig) SwarmClient {
	if config == nil {
		config = DefaultSwarmConfig()
	}

	return &swarmClient{
		mcpServer:    mcpServer,
		toolRegistry: mcpServer.GetToolRegistry(),
		memory:       mcpServer.GetMemory(),
		config:       config,
		functions:    make(map[string]AgentFunction),
		agents:       make(map[string]*Agent),
		logger:       &defaultLogger{},
		modelFunc:    nil, // Will be set via SetModel
		modelName:    "llama3.2",
	}
}

// NewSwarmClientWithLLM creates a new swarm client with LLM integration
func NewSwarmClientWithLLM(mcpServer interface {
	GetToolRegistry() *mcp.ToolRegistry
	GetMemory() *mcp.Memory
}, config *SwarmConfig, modelFunc mcp.ModelFunc, modelName string) SwarmClient {
	if config == nil {
		config = DefaultSwarmConfig()
	}

	return &swarmClient{
		mcpServer:    mcpServer,
		toolRegistry: mcpServer.GetToolRegistry(),
		memory:       mcpServer.GetMemory(),
		config:       config,
		functions:    make(map[string]AgentFunction),
		agents:       make(map[string]*Agent),
		logger:       &defaultLogger{},
		modelFunc:    modelFunc,
		modelName:    modelName,
	}
}

// CreateAgent creates a new agent with specified tools
func (sc *swarmClient) CreateAgent(name, instructions string, tools []string) *Agent {
	agent := &Agent{
		Name:         name,
		Instructions: instructions,
		Functions:    []AgentFunction{},
		Model:        "gpt-4o",
	}

	// Add MCP tools as agent functions
	for _, toolName := range tools {
		if strings.HasPrefix(toolName, "transfer_to_") {
			// Handle transfer functions specially
			continue
		}

		// Create agent function from MCP tool
		agentFunc := AgentFunction{
			Name:        toolName,
			Description: fmt.Sprintf("MCP tool: %s", toolName),
			Parameters:  map[string]interface{}{},
			Function: func(args map[string]interface{}, contextVars map[string]interface{}) Result {
				result, err := sc.toolRegistry.Call(toolName, args, sc.memory)
				if err != nil {
					return Result{
						Value:   fmt.Sprintf("Error calling tool %s: %v", toolName, err),
						Success: false,
						Error:   err,
					}
				}

				return Result{
					Value:   fmt.Sprintf("%v", result),
					Success: true,
				}
			},
		}

		agent.Functions = append(agent.Functions, agentFunc)
		sc.functions[toolName] = agentFunc
	}

	sc.agents[name] = agent
	return agent
}

// CreateAgentWithModel creates a new agent with specified tools and model configuration
func (sc *swarmClient) CreateAgentWithModel(name, instructions string, tools []string, modelConfig *conduit.ModelConfig) *Agent {
	if modelConfig == nil {
		modelConfig = &conduit.ModelConfig{
			Provider:    "ollama",
			Model:       "llama3.2",
			URL:         "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   1000,
			TopK:        40,
		}
	}

	// Create model function from configuration
	modelFunc, err := conduit.CreateModelFunction(modelConfig)
	if err != nil {
		sc.logger.Error("Failed to create model function", "error", err, "config", modelConfig)
		// Fallback to swarm default model
		return sc.CreateAgent(name, instructions, tools)
	}

	agent := &Agent{
		Name:         name,
		Instructions: instructions,
		Functions:    []AgentFunction{},
		Model:        modelConfig.Model,
		ModelFunc:    modelFunc,
		ModelConfig:  modelConfig,
	}

	// Add MCP tools as agent functions
	for _, toolName := range tools {
		if strings.HasPrefix(toolName, "transfer_to_") {
			// Handle transfer functions specially
			continue
		}

		// Create MCP tool function
		mcpToolFunc := AgentFunction{
			Name:        toolName,
			Description: fmt.Sprintf("MCP tool: %s", toolName),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"params": map[string]interface{}{
						"type": "object",
					},
				},
			},
			Function: func(args map[string]interface{}, contextVars map[string]interface{}) Result {
				params, ok := args["params"].(map[string]interface{})
				if !ok {
					params = args
				}

				result, err := sc.toolRegistry.Call(toolName, params, sc.memory)
				if err != nil {
					return Result{
						Error:   err,
						Success: false,
					}
				}

				return Result{
					Value:   fmt.Sprintf("Tool %s executed successfully: %v", toolName, result),
					Success: true,
				}
			},
		}

		agent.Functions = append(agent.Functions, mcpToolFunc)
	}

	sc.agents[name] = agent
	return agent
}

// CreateAgentWithLLM creates a new agent with specified tools and direct LLM function
func (sc *swarmClient) CreateAgentWithLLM(name, instructions string, tools []string, modelFunc mcp.ModelFunc, modelName string) *Agent {
	agent := &Agent{
		Name:         name,
		Instructions: instructions,
		Functions:    []AgentFunction{},
		Model:        modelName,
		ModelFunc:    modelFunc,
		ModelConfig:  nil, // No config when using direct model function
	}

	// Add MCP tools as agent functions
	for _, toolName := range tools {
		if strings.HasPrefix(toolName, "transfer_to_") {
			// Handle transfer functions specially
			continue
		}

		// Create MCP tool function
		mcpToolFunc := AgentFunction{
			Name:        toolName,
			Description: fmt.Sprintf("MCP tool: %s", toolName),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"params": map[string]interface{}{
						"type": "object",
					},
				},
			},
			Function: func(args map[string]interface{}, contextVars map[string]interface{}) Result {
				params, ok := args["params"].(map[string]interface{})
				if !ok {
					params = args
				}

				result, err := sc.toolRegistry.Call(toolName, params, sc.memory)
				if err != nil {
					return Result{
						Error:   err,
						Success: false,
					}
				}

				return Result{
					Value:   fmt.Sprintf("Tool %s executed successfully: %v", toolName, result),
					Success: true,
				}
			},
		}

		agent.Functions = append(agent.Functions, mcpToolFunc)
	}

	sc.agents[name] = agent
	return agent
}

// RegisterFunction registers a custom function for a specific agent
func (sc *swarmClient) RegisterFunction(agentName string, fn AgentFunction) error {
	// Store in global functions map
	sc.functions[fn.Name] = fn

	// Add to the specific agent's function list
	if agent, exists := sc.agents[agentName]; exists {
		agent.Functions = append(agent.Functions, fn)
	} else {
		return fmt.Errorf("agent %s not found", agentName)
	}

	return nil
}

// Run executes the swarm with the given agent and messages
func (sc *swarmClient) Run(agent *Agent, messages []Message, contextVars map[string]interface{}) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), sc.config.Timeout)
	defer cancel()
	return sc.RunWithContext(ctx, agent, messages, contextVars)
}

// RunWithContext executes the swarm with context
func (sc *swarmClient) RunWithContext(ctx context.Context, agent *Agent, messages []Message, contextVars map[string]interface{}) *Response {
	startTime := time.Now()

	if contextVars == nil {
		contextVars = make(map[string]interface{})
	}

	execCtx := &ExecutionContext{
		SessionID:     fmt.Sprintf("session_%d", time.Now().UnixNano()),
		CurrentAgent:  agent,
		MessageCount:  len(messages),
		ToolCallCount: 0,
		HandoffCount:  0,
		StartTime:     startTime,
		Memory:        sc.memory,
		Debug:         sc.config.Debug,
	}

	response := &Response{
		Messages:       make([]Message, len(messages)),
		Agent:          agent,
		ContextVars:    contextVars,
		TotalTurns:     0,
		ToolCallsCount: 0,
		HandoffsCount:  0,
		Success:        false,
	}

	// Copy input messages
	copy(response.Messages, messages)

	currentAgent := agent
	turnCount := 0

	// Main execution loop
	for turnCount < sc.config.MaxTurns {
		select {
		case <-ctx.Done():
			response.Error = ctx.Err()
			response.ExecutionTime = time.Since(startTime)
			return response
		default:
		}

		turnCount++
		execCtx.CurrentAgent = currentAgent

		// Simulate LLM processing with system prompt
		systemMessage := Message{
			Role:    "system",
			Content: currentAgent.Instructions,
		}

		// Add context variables to system prompt if available
		if len(contextVars) > 0 {
			contextStr := sc.buildContextString(contextVars)
			systemMessage.Content += fmt.Sprintf("\n\nContext: %s", contextStr)
		}

		// Process the current conversation turn
		turnResult := sc.processTurn(execCtx, currentAgent, response.Messages, contextVars)

		if turnResult.Error != nil {
			response.Error = turnResult.Error
			response.ExecutionTime = time.Since(startTime)
			return response
		}

		// Add assistant response
		if turnResult.ResponseMessage != nil {
			response.Messages = append(response.Messages, *turnResult.ResponseMessage)
		}

		// Handle agent handoff
		if turnResult.Agent != nil && turnResult.Agent != currentAgent {
			execCtx.HandoffCount++
			response.HandoffsCount++
			currentAgent = turnResult.Agent
			response.Agent = currentAgent

			// Add handoff message
			handoffMsg := Message{
				Role:    "assistant",
				Content: fmt.Sprintf("Transferring to %s", currentAgent.Name),
			}
			response.Messages = append(response.Messages, handoffMsg)
		}

		// Update context variables
		if turnResult.ContextVars != nil {
			for k, v := range turnResult.ContextVars {
				contextVars[k] = v
			}
		}

		// If no further actions needed, break
		if turnResult.Agent == nil && turnResult.Value != "" && !strings.Contains(turnResult.Value, "transfer_to_") {
			break
		}
	}

	response.TotalTurns = turnCount
	response.ContextVars = contextVars
	response.ExecutionTime = time.Since(startTime)
	response.Success = response.Error == nil

	return response
}

// processTurn handles a single conversation turn with LLM reasoning
func (sc *swarmClient) processTurn(ctx *ExecutionContext, agent *Agent, messages []Message, contextVars map[string]interface{}) Result {
	// Get the last user message
	lastMessage := ""
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastMessage = messages[i].Content
			break
		}
	}

	// Use LLM for intelligent reasoning - no fallback
	if !sc.HasLLM() {
		return Result{
			Error:   fmt.Errorf("no LLM model configured - swarm requires LLM for intelligent reasoning"),
			Success: false,
		}
	}

	return sc.processWithLLM(ctx, agent, lastMessage, messages, contextVars)
}

// processWithLLM uses LLM for intelligent agent reasoning and decision-making
func (sc *swarmClient) processWithLLM(ctx *ExecutionContext, agent *Agent, message string, messages []Message, contextVars map[string]interface{}) Result {
	// Create comprehensive system prompt for the agent
	systemPrompt := sc.buildAgentSystemPrompt(agent, contextVars)

	// Create conversation history
	conversationHistory := sc.buildConversationHistory(messages)

	// Build LLM prompt for agent reasoning
	prompt := fmt.Sprintf(`%s

CONVERSATION HISTORY:
%s

USER MESSAGE: %s

Please analyze this message and decide what action to take. You can:
1. Use one of your available tools
2. Transfer to another agent (if appropriate)
3. Respond directly

Your available tools: %s
Available agents for handoff: %s

Respond with a JSON object containing your decision:
{
  "action": "tool_use|handoff|respond",
  "reasoning": "explain your decision",
  "tool_name": "tool to use (if action=tool_use)",
  "tool_args": {"param": "value"},
  "handoff_agent": "agent name (if action=handoff)",
  "response": "direct response (if action=respond)"
}`, systemPrompt, conversationHistory, message, sc.getAvailableToolsForAgent(agent), sc.getAvailableAgentsForHandoff(agent))

	// Call LLM for reasoning - use agent-specific model if available
	llmResponse, err := sc.callAgentLLM(agent, prompt, ctx.SessionID)
	if err != nil {
		return Result{
			Error:   fmt.Errorf("LLM reasoning failed: %w", err),
			Success: false,
		}
	}

	// Parse LLM response and execute decision
	return sc.executeLLMDecision(ctx, agent, llmResponse, contextVars)
}

// callLLM makes a call to the configured LLM model
func (sc *swarmClient) callLLM(prompt string, sessionID string) (string, error) {
	if sc.modelFunc == nil {
		return "", fmt.Errorf("no LLM model configured")
	}

	ctx := mcp.ContextInput{
		ContextID: sessionID,
		Inputs: map[string]interface{}{
			"query": prompt,
		},
	}

	req := mcp.MCPRequest{
		SessionID: sessionID,
		Model:     sc.modelName,
		Stream:    false,
	}

	response, err := sc.modelFunc(ctx, req, sc.memory, nil)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return response, nil
}

// callAgentLLM makes a call to the agent-specific LLM model or fallback to swarm model
func (sc *swarmClient) callAgentLLM(agent *Agent, prompt string, sessionID string) (string, error) {
	// Try agent-specific model first
	modelFunc := agent.ModelFunc
	modelName := agent.Model

	// Fallback to swarm-level model if agent doesn't have one
	if modelFunc == nil {
		modelFunc = sc.modelFunc
		modelName = sc.modelName
	}

	// If still no model, return error
	if modelFunc == nil {
		return "", fmt.Errorf("no LLM configured for agent %s or swarm", agent.Name)
	}

	// Create context for LLM call
	ctx := mcp.ContextInput{
		ContextID: sessionID,
		Inputs: map[string]interface{}{
			"query": prompt,
		},
	}

	// Create request
	req := mcp.MCPRequest{
		SessionID: sessionID,
		Model:     modelName,
		Stream:    false,
	}

	// Call the model function
	response, err := modelFunc(ctx, req, sc.memory, nil)
	if err != nil {
		return "", fmt.Errorf("agent LLM call failed: %w", err)
	}

	return response, nil
}

// buildAgentSystemPrompt creates a comprehensive system prompt for the agent
func (sc *swarmClient) buildAgentSystemPrompt(agent *Agent, contextVars map[string]interface{}) string {
	prompt := fmt.Sprintf("You are %s.\n\n%s", agent.Name, agent.Instructions)

	if len(contextVars) > 0 {
		prompt += "\n\nCurrent context:\n"
		for k, v := range contextVars {
			prompt += fmt.Sprintf("- %s: %v\n", k, v)
		}
	}

	return prompt
}

// buildConversationHistory formats conversation history for LLM context
func (sc *swarmClient) buildConversationHistory(messages []Message) string {
	if len(messages) == 0 {
		return "No previous conversation."
	}

	history := ""
	// Limit to last 5 messages to prevent prompt overflow
	start := len(messages) - 5
	if start < 0 {
		start = 0
	}

	for _, msg := range messages[start:] {
		history += fmt.Sprintf("%s: %s\n", strings.ToUpper(msg.Role), msg.Content)
	}

	return strings.TrimSpace(history)
}

// getAvailableToolsForAgent returns list of tools available to the agent
func (sc *swarmClient) getAvailableToolsForAgent(agent *Agent) string {
	if len(agent.Functions) == 0 {
		return "No tools available"
	}

	tools := make([]string, 0, len(agent.Functions))
	for _, fn := range agent.Functions {
		if !strings.HasPrefix(fn.Name, "transfer_to_") {
			tools = append(tools, fmt.Sprintf("%s: %s", fn.Name, fn.Description))
		}
	}

	if len(tools) == 0 {
		return "No tools available"
	}

	return strings.Join(tools, ", ")
}

// getAvailableAgentsForHandoff returns list of agents available for handoff
func (sc *swarmClient) getAvailableAgentsForHandoff(currentAgent *Agent) string {
	agents := make([]string, 0, len(sc.agents))
	for name, agent := range sc.agents {
		if agent.Name != currentAgent.Name {
			agents = append(agents, name)
		}
	}

	if len(agents) == 0 {
		return "No other agents available"
	}

	return strings.Join(agents, ", ")
}

// executeLLMDecision parses and executes the LLM's decision
func (sc *swarmClient) executeLLMDecision(ctx *ExecutionContext, agent *Agent, llmResponse string, contextVars map[string]interface{}) Result {
	// Try to extract JSON from LLM response
	decision, err := sc.parseLLMDecision(llmResponse)
	if err != nil {
		return Result{
			Error:   fmt.Errorf("failed to parse LLM decision: %w", err),
			Success: false,
		}
	}

	switch decision.Action {
	case "tool_use":
		return sc.executeLLMToolUse(ctx, agent, decision, contextVars)
	case "handoff":
		return sc.executeLLMHandoff(agent, decision, contextVars)
	case "respond":
		return Result{
			Value:   decision.Response,
			Success: true,
			ResponseMessage: &Message{
				Role:    "assistant",
				Content: decision.Response,
			},
		}
	default:
		return Result{
			Error:   fmt.Errorf("unknown LLM action: %s", decision.Action),
			Success: false,
		}
	}
}

// LLMDecision represents the parsed decision from LLM
type LLMDecision struct {
	Action       string                 `json:"action"`
	Reasoning    string                 `json:"reasoning"`
	ToolName     string                 `json:"tool_name"`
	ToolArgs     map[string]interface{} `json:"tool_args"`
	HandoffAgent string                 `json:"handoff_agent"`
	Response     string                 `json:"response"`
}

// parseLLMDecision attempts to parse JSON decision from LLM response
func (sc *swarmClient) parseLLMDecision(response string) (*LLMDecision, error) {
	// Try to find JSON in the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		// If no JSON found, treat as direct response
		return &LLMDecision{
			Action:   "respond",
			Response: response,
		}, nil
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var decision LLMDecision
	if err := sc.parseJSON(jsonStr, &decision); err != nil {
		// If JSON parsing fails, treat as direct response
		return &LLMDecision{
			Action:   "respond",
			Response: response,
		}, nil
	}

	return &decision, nil
}

// parseJSON is a simple JSON parser for LLM responses
func (sc *swarmClient) parseJSON(jsonStr string, decision *LLMDecision) error {
	// Simple parsing logic - in production, use proper JSON parser
	// This is simplified for the demo

	if strings.Contains(jsonStr, `"action": "tool_use"`) {
		decision.Action = "tool_use"

		// Extract tool name
		if start := strings.Index(jsonStr, `"tool_name": "`); start != -1 {
			start += len(`"tool_name": "`)
			if end := strings.Index(jsonStr[start:], `"`); end != -1 {
				decision.ToolName = jsonStr[start : start+end]
			}
		}

		// Simple tool args extraction - in practice, use proper JSON
		decision.ToolArgs = map[string]interface{}{}

	} else if strings.Contains(jsonStr, `"action": "handoff"`) {
		decision.Action = "handoff"

		// Extract handoff agent
		if start := strings.Index(jsonStr, `"handoff_agent": "`); start != -1 {
			start += len(`"handoff_agent": "`)
			if end := strings.Index(jsonStr[start:], `"`); end != -1 {
				decision.HandoffAgent = jsonStr[start : start+end]
			}
		}

	} else {
		decision.Action = "respond"

		// Extract response
		if start := strings.Index(jsonStr, `"response": "`); start != -1 {
			start += len(`"response": "`)
			if end := strings.Index(jsonStr[start:], `"`); end != -1 {
				decision.Response = jsonStr[start : start+end]
			}
		}
	}

	return nil
}

// executeLLMToolUse executes a tool use decision from LLM
func (sc *swarmClient) executeLLMToolUse(ctx *ExecutionContext, agent *Agent, decision *LLMDecision, contextVars map[string]interface{}) Result {
	// Find the tool function
	for _, fn := range agent.Functions {
		if fn.Name == decision.ToolName {
			if sc.config.Debug {
				sc.logger.Debug("Executing LLM-selected tool", "tool", decision.ToolName, "reasoning", decision.Reasoning)
			}

			ctx.ToolCallCount++
			result := fn.Function(decision.ToolArgs, contextVars)

			if result.Success {
				result.ResponseMessage = &Message{
					Role:    "assistant",
					Content: fmt.Sprintf("I used %s: %s", decision.ToolName, result.Value),
				}
			}

			return result
		}
	}

	return Result{
		Error:   fmt.Errorf("tool %s not found", decision.ToolName),
		Success: false,
	}
}

// executeLLMHandoff executes a handoff decision from LLM
func (sc *swarmClient) executeLLMHandoff(agent *Agent, decision *LLMDecision, contextVars map[string]interface{}) Result {
	// Find the target agent
	if targetAgent, exists := sc.agents[decision.HandoffAgent]; exists {
		if sc.config.Debug {
			sc.logger.Info("LLM-directed agent handoff", "from", agent.Name, "to", targetAgent.Name, "reasoning", decision.Reasoning)
		}

		return Result{
			Value:   fmt.Sprintf("Transferring to %s", targetAgent.Name),
			Agent:   targetAgent,
			Success: true,
			ResponseMessage: &Message{
				Role:    "assistant",
				Content: fmt.Sprintf("Transferring to %s", targetAgent.Name),
			},
		}
	}

	return Result{
		Error:   fmt.Errorf("agent %s not found for handoff", decision.HandoffAgent),
		Success: false,
	}
}

// detectAndExecuteFunction detects if a function should be called based on message content
func (sc *swarmClient) detectAndExecuteFunction(ctx *ExecutionContext, agent *Agent, message string, contextVars map[string]interface{}) Result {
	messageLower := strings.ToLower(message)

	// Enhanced transfer detection based on message content and agent type
	if agent.Name == "Coordinator" {
		// Coordinator should delegate based on task type
		if strings.Contains(messageLower, "article") || strings.Contains(messageLower, "write") ||
			strings.Contains(messageLower, "content") || strings.Contains(messageLower, "research") {
			return sc.executeTransferFunction(agent, "transfer_to_content_creator", contextVars)
		}

		if strings.Contains(messageLower, "analyz") || strings.Contains(messageLower, "data") ||
			strings.Contains(messageLower, "dataset") || strings.Contains(messageLower, "report") {
			return sc.executeTransferFunction(agent, "transfer_to_data_analyst", contextVars)
		}

		if strings.Contains(messageLower, "retrieve") || strings.Contains(messageLower, "memory") ||
			strings.Contains(messageLower, "context") || strings.Contains(messageLower, "summarize") {
			return sc.executeTransferFunction(agent, "transfer_to_memory_manager", contextVars)
		}
	}

	// Check for tool calls based on agent capabilities and message content
	if len(agent.Functions) > 0 {
		// Look for MCP tool calls
		for _, fn := range agent.Functions {
			if sc.shouldExecuteFunction(fn.Name, messageLower, agent.Name) {
				args := sc.extractArgumentsFromMessage(fn.Name, message)

				if sc.config.Debug {
					sc.logger.Debug("Executing function", "function", fn.Name, "args", args)
				}

				ctx.ToolCallCount++
				result := fn.Function(args, contextVars)

				if result.Success {
					result.ResponseMessage = &Message{
						Role:    "assistant",
						Content: fmt.Sprintf("I executed %s: %s", fn.Name, result.Value),
					}
				}

				return result
			}
		}
	}

	return Result{Success: false}
}

// executeTransferFunction executes a transfer function by name
func (sc *swarmClient) executeTransferFunction(agent *Agent, functionName string, contextVars map[string]interface{}) Result {
	for _, fn := range agent.Functions {
		if fn.Name == functionName {
			result := fn.Function(map[string]interface{}{}, contextVars)
			if sc.config.Debug {
				sc.logger.Info("Agent handoff executed", "from", agent.Name, "to", result.Agent.Name)
			}
			return result
		}
	}
	return Result{Success: false}
}

// shouldExecuteFunction determines if a function should be executed based on context
func (sc *swarmClient) shouldExecuteFunction(functionName, messageLower, agentName string) bool {
	switch agentName {
	case "ContentCreator":
		if functionName == "research_topic" && (strings.Contains(messageLower, "research") || strings.Contains(messageLower, "article")) {
			return true
		}
		if functionName == "write_article" && (strings.Contains(messageLower, "write") || strings.Contains(messageLower, "article")) {
			return true
		}
		if functionName == "store_context" && strings.Contains(messageLower, "remember") {
			return true
		}
	case "DataAnalyst":
		if functionName == "analyze_data" && (strings.Contains(messageLower, "analyz") || strings.Contains(messageLower, "data")) {
			return true
		}
		if functionName == "generate_report" && strings.Contains(messageLower, "report") {
			return true
		}
		if functionName == "store_context" && strings.Contains(messageLower, "remember") {
			return true
		}
	case "MemoryManager":
		if functionName == "store_context" && (strings.Contains(messageLower, "store") || strings.Contains(messageLower, "remember")) {
			return true
		}
		if functionName == "retrieve_context" && (strings.Contains(messageLower, "retrieve") || strings.Contains(messageLower, "recall")) {
			return true
		}
	}

	return false
}

// extractArgumentsFromMessage extracts arguments from natural language message
func (sc *swarmClient) extractArgumentsFromMessage(toolName, message string) map[string]interface{} {
	args := make(map[string]interface{})

	switch toolName {
	case "research_topic":
		// Extract topic from message about articles or research
		if strings.Contains(strings.ToLower(message), "artificial intelligence") || strings.Contains(strings.ToLower(message), "ai") {
			if strings.Contains(strings.ToLower(message), "healthcare") {
				args["topic"] = "artificial intelligence in healthcare"
			} else {
				args["topic"] = "artificial intelligence"
			}
		} else {
			// Default topic extraction
			words := strings.Fields(message)
			if len(words) > 3 {
				args["topic"] = strings.Join(words[len(words)-3:], " ")
			} else {
				args["topic"] = "general research topic"
			}
		}
	case "write_article":
		// Extract title and topic for article writing
		if strings.Contains(strings.ToLower(message), "ai") && strings.Contains(strings.ToLower(message), "healthcare") {
			args["title"] = "AI in Healthcare: A Comprehensive Guide"
			args["topic"] = "artificial intelligence in healthcare"
		} else {
			args["title"] = "Generated Article"
			args["topic"] = "general topic"
		}
	case "analyze_data":
		// Extract dataset information
		if strings.Contains(strings.ToLower(message), "customer") {
			args["dataset"] = "customer_behavior_q4_2023.csv"
		} else {
			args["dataset"] = "data.csv"
		}
	case "generate_report":
		// Extract findings for report generation
		args["findings"] = "Analysis complete with key insights and findings"
	case "store_context", "retrieve_context":
		// Extract key-value pairs for memory operations
		if strings.Contains(strings.ToLower(message), "project") {
			args["key"] = "project_status"
			args["value"] = "completed"
		} else {
			args["key"] = "general_context"
			args["value"] = "stored information"
		}
	case "create_task":
		// Extract task and assignee information
		args["task"] = "Generated task from user request"
		args["assignee"] = "specialist"
	default:
		// For other tools, try to extract text content
		words := strings.Fields(message)
		if len(words) > 0 {
			args["text"] = strings.Join(words, " ")
		}
	}

	return args
}

// generateResponse generates a response based on agent instructions
func (sc *swarmClient) generateResponse(agent *Agent, message string, contextVars map[string]interface{}) string {
	// Simple response generation based on agent type
	agentName := strings.ToLower(agent.Name)

	switch {
	case strings.Contains(agentName, "coordinator"):
		return sc.handleCoordinatorResponse(message, contextVars)
	case strings.Contains(agentName, "content"):
		return "I can help you with text processing tasks like converting case, counting words, and more."
	case strings.Contains(agentName, "data"):
		return "I can help you analyze data, generate UUIDs, timestamps, and encode information."
	case strings.Contains(agentName, "memory"):
		return "I can help you store and retrieve information from memory."
	default:
		return fmt.Sprintf("Hello! I'm %s. How can I help you?", agent.Name)
	}
}

// handleCoordinatorResponse handles responses from the coordinator agent
func (sc *swarmClient) handleCoordinatorResponse(message string, contextVars map[string]interface{}) string {
	messageLower := strings.ToLower(message)

	if strings.Contains(messageLower, "text") || strings.Contains(messageLower, "convert") || strings.Contains(messageLower, "case") {
		return "I'll transfer you to the Content Creator for text processing tasks."
	} else if strings.Contains(messageLower, "uuid") || strings.Contains(messageLower, "data") || strings.Contains(messageLower, "analyze") {
		return "I'll transfer you to the Data Analyst for data processing tasks."
	} else if strings.Contains(messageLower, "remember") || strings.Contains(messageLower, "store") || strings.Contains(messageLower, "memory") {
		return "I'll transfer you to the Memory Manager for information storage tasks."
	}

	return "I can help route your request to the appropriate specialist. What would you like to do?"
}

// buildContextString builds a string representation of context variables
func (sc *swarmClient) buildContextString(contextVars map[string]interface{}) string {
	if len(contextVars) == 0 {
		return ""
	}

	parts := make([]string, 0, len(contextVars))
	for k, v := range contextVars {
		parts = append(parts, fmt.Sprintf("%s: %v", k, v))
	}

	return strings.Join(parts, ", ")
}

// GetAvailableTools returns the list of available MCP tools
func (sc *swarmClient) GetAvailableTools() []string {
	return sc.toolRegistry.GetRegisteredTools()
}

// GetMemory returns the MCP memory instance
func (sc *swarmClient) GetMemory() *mcp.Memory {
	return sc.memory
}

// SetModel sets the LLM model function for agent reasoning
func (sc *swarmClient) SetModel(modelFunc mcp.ModelFunc, modelName string) {
	sc.modelFunc = modelFunc
	sc.modelName = modelName
}

// HasLLM returns true if LLM is configured
func (sc *swarmClient) HasLLM() bool {
	return sc.modelFunc != nil
}

// CreateHandoffFunction creates a function that transfers to another agent
func CreateHandoffFunction(name string, targetAgent *Agent) AgentFunction {
	return AgentFunction{
		Name:        fmt.Sprintf("transfer_to_%s", name),
		Description: fmt.Sprintf("Transfer the conversation to %s", targetAgent.Name),
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Function: func(args map[string]interface{}, contextVars map[string]interface{}) Result {
			return Result{
				Value:   fmt.Sprintf("Transferring to %s", targetAgent.Name),
				Agent:   targetAgent,
				Success: true,
				ResponseMessage: &Message{
					Role:    "assistant",
					Content: fmt.Sprintf("Transferring to %s", targetAgent.Name),
				},
			}
		},
	}
}
