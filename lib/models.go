package conduit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/benozo/conduit/mcp"
)

// OllamaRequest represents a request to Ollama
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaChunk represents a streaming response chunk from Ollama
type OllamaChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// CreateOllamaModel creates an Ollama model function
func CreateOllamaModel(ollamaURL string) mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		payload := OllamaRequest{
			Model:  req.Model,
			Prompt: query,
			Stream: false,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(context.Background(), "POST", ollamaURL+"/api/generate", bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			return "", fmt.Errorf("failed to call Ollama: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("Ollama returned status %d", resp.StatusCode)
		}

		var result strings.Builder
		decoder := json.NewDecoder(resp.Body)

		for {
			var chunk OllamaChunk
			if err := decoder.Decode(&chunk); err != nil {
				if err == io.EOF {
					break
				}
				return "", fmt.Errorf("failed to decode chunk: %w", err)
			}

			if chunk.Response != "" {
				if onToken != nil {
					onToken(ctx.ContextID, chunk.Response)
				}
				result.WriteString(chunk.Response)
			}

			if chunk.Done {
				break
			}
		}

		return result.String(), nil
	}
}

// OllamaChatRequest represents a chat request to Ollama with tool support
type OllamaChatRequest struct {
	Model    string                  `json:"model"`
	Messages []OllamaChatMessage     `json:"messages"`
	Stream   bool                    `json:"stream"`
	Tools    []OllamaToolDescription `json:"tools,omitempty"`
}

// OllamaChatMessage represents a chat message
type OllamaChatMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content"`
	Name      string           `json:"name,omitempty"` // Tool name for tool messages
	ToolCalls []OllamaToolCall `json:"tool_calls,omitempty"`
}

// OllamaToolDescription represents a tool description for Ollama
type OllamaToolDescription struct {
	Type     string            `json:"type"`
	Function OllamaFunctionDef `json:"function"`
}

// OllamaFunctionDef represents a function definition
type OllamaFunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// OllamaChatChunk represents a streaming chat response from Ollama
type OllamaChatChunk struct {
	Message OllamaChatMessage `json:"message"`
	Done    bool              `json:"done"`
}

// OllamaToolCall represents a tool call from Ollama
type OllamaToolCall struct {
	Function OllamaFunctionCall `json:"function"`
}

// OllamaFunctionCall represents a function call
type OllamaFunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// CreateOllamaToolAwareModel creates an Ollama model function with tool support
func CreateOllamaToolAwareModel(ollamaURL string, tools *mcp.ToolRegistry) mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		// Use provided model or default
		model := req.Model
		if model == "" {
			model = "llama3.2" // Default model
		}

		log.Printf("🔍 CreateOllamaToolAwareModel called with query: %s", query)
		log.Printf("🔧 Ollama URL: %s", ollamaURL)
		log.Printf("🔧 Model: %s", model)

		// First try with tool-aware chat API
		result, err := tryOllamaWithTools(ollamaURL, query, model, tools, memory, onToken, ctx.ContextID)
		if err == nil && result != "" {
			log.Printf("✅ Ollama tool-aware request succeeded")
			return result, nil
		}

		log.Printf("⚠️ Ollama tool-aware failed, trying prompt-based approach: %v", err)

		// Fallback: Use prompt engineering to simulate tool calling
		return tryOllamaWithPromptTools(ollamaURL, query, model, tools, memory, onToken, ctx.ContextID)
	}
}

// tryOllamaWithTools attempts to use Ollama's native tool calling
func tryOllamaWithTools(ollamaURL, query, model string, tools *mcp.ToolRegistry, memory *mcp.Memory, onToken mcp.StreamCallback, contextID string) (string, error) {
	// Convert tools to Ollama format
	var ollamaTools []OllamaToolDescription
	if tools != nil {
		log.Printf("🔧 Available tools: %v", tools.GetRegisteredTools())
		for _, name := range tools.GetRegisteredTools() {
			// Use the same pattern as the working curl example
			// All tools get "Tool: name" description and text parameter by default
			description := fmt.Sprintf("Tool: %s", name)
			parameters := map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "Input text",
					},
				},
				"required": []string{"text"},
			}

			// Special cases for tools that need different parameters
			switch name {
			case "remember":
				description = "Store a value in memory"
				parameters = map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"key": map[string]interface{}{
							"type":        "string",
							"description": "Memory key",
						},
						"value": map[string]interface{}{
							"type":        "string",
							"description": "Value to store",
						},
					},
					"required": []string{"key", "value"},
				}
			case "recall":
				description = "Retrieve a value from memory"
				parameters = map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"key": map[string]interface{}{
							"type":        "string",
							"description": "Memory key",
						},
					},
					"required": []string{"key"},
				}
			case "forget":
				description = "Remove a value from memory"
				parameters = map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"key": map[string]interface{}{
							"type":        "string",
							"description": "Memory key",
						},
					},
					"required": []string{"key"},
				}
			case "replace":
				description = "Replace text patterns"
				parameters = map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type":        "string",
							"description": "Input text",
						},
						"old": map[string]interface{}{
							"type":        "string",
							"description": "Text to replace",
						},
						"new": map[string]interface{}{
							"type":        "string",
							"description": "Replacement text",
						},
					},
					"required": []string{"text", "old", "new"},
				}
			case "add":
				description = "Add two numbers together"
				parameters = map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"a": map[string]interface{}{
							"type":        "number",
							"description": "First number to add",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second number to add",
						},
					},
					"required": []string{"a", "b"},
				}
			case "timestamp":
				description = "Get current timestamp"
				parameters = map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"format": map[string]interface{}{
							"type":        "string",
							"description": "Timestamp format (iso, unix, readable)",
						},
					},
					// timestamp doesn't require any parameters, making format optional
				}
			case "uuid":
				description = "Generate a UUID"
				// Keep the text parameter as in working example for consistency
			}

			ollamaTools = append(ollamaTools, OllamaToolDescription{
				Type: "function",
				Function: OllamaFunctionDef{
					Name:        name,
					Description: description,
					Parameters:  parameters,
				},
			})
		}
		log.Printf("🔧 Converted %d tools to Ollama format", len(ollamaTools))
	}

	payload := OllamaChatRequest{
		Model: model,
		Messages: []OllamaChatMessage{
			{
				Role:    "user",
				Content: query,
			},
		},
		Stream: false, // Try non-streaming first
		Tools:  ollamaTools,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	log.Printf("🚀 Sending tool-aware request to Ollama")
	log.Printf("📤 Payload: %s", string(body))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", ollamaURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create chat request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama chat API: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("📡 Response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama chat API returned status %d", resp.StatusCode)
	}

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("📦 Raw response: %s", string(respBody))

	var chatResp OllamaChatChunk
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("🔧 Decoded response: Message.Content='%s', ToolCalls=%d", chatResp.Message.Content, len(chatResp.Message.ToolCalls))

	// Check if we got tool calls
	if len(chatResp.Message.ToolCalls) > 0 {
		log.Printf("🔧 Ollama requested %d tool calls!", len(chatResp.Message.ToolCalls))

		// Execute tools and collect results
		var toolMessages []OllamaChatMessage

		// Start with the original user message
		toolMessages = append(toolMessages, OllamaChatMessage{
			Role:    "user",
			Content: query,
		})

		// Add the assistant's message with tool calls
		toolMessages = append(toolMessages, OllamaChatMessage{
			Role:      "assistant",
			Content:   chatResp.Message.Content,
			ToolCalls: chatResp.Message.ToolCalls,
		})

		// Execute each tool and create tool result messages
		for i, toolCall := range chatResp.Message.ToolCalls {
			log.Printf("🔧 Tool call %d: %s with args %+v", i+1, toolCall.Function.Name, toolCall.Function.Arguments)

			var toolResult interface{}
			var toolErr error

			if tools != nil {
				toolResult, toolErr = tools.Call(toolCall.Function.Name, toolCall.Function.Arguments, memory)
			}

			// Create tool result message
			var resultContent string
			if toolErr != nil {
				log.Printf("❌ Tool %s failed: %v", toolCall.Function.Name, toolErr)
				resultContent = fmt.Sprintf("Error: %v", toolErr)
			} else {
				log.Printf("✅ Tool %s succeeded: %v", toolCall.Function.Name, toolResult)
				resultContent = fmt.Sprintf("%v", toolResult)
			}

			// Add tool result message
			toolMessages = append(toolMessages, OllamaChatMessage{
				Role:    "tool",
				Content: resultContent,
				Name:    toolCall.Function.Name,
			})
		}

		// Send tool results back to Ollama for final response
		log.Printf("🔄 Sending tool results back to Ollama for final response...")
		return sendToolResultsToOllama(ollamaURL, model, query, toolMessages, onToken, contextID)
	}

	// If no tool calls but we have content, return it
	if chatResp.Message.Content != "" {
		log.Printf("💬 Got message content without tool calls: %s", chatResp.Message.Content)
		return chatResp.Message.Content, nil
	}

	// No tool calls and no content - this is the failure case
	return "", fmt.Errorf("ollama returned empty response")
}

// tryOllamaWithPromptTools uses prompt engineering to simulate tool calling
func tryOllamaWithPromptTools(ollamaURL, query, model string, tools *mcp.ToolRegistry, memory *mcp.Memory, onToken mcp.StreamCallback, contextID string) (string, error) {
	log.Printf("🔧 Using prompt-based tool calling approach")

	// Create a prompt that instructs the model to use tools
	var toolList strings.Builder
	if tools != nil {
		toolList.WriteString("Available tools:\n")
		for _, name := range tools.GetRegisteredTools() {
			switch name {
			case "uuid":
				toolList.WriteString("- uuid(): Generate a unique identifier\n")
			case "timestamp":
				toolList.WriteString("- timestamp(format): Get current time (formats: iso, unix, readable)\n")
			case "uppercase":
				toolList.WriteString("- uppercase(text): Convert text to uppercase\n")
			case "lowercase":
				toolList.WriteString("- lowercase(text): Convert text to lowercase\n")
			case "base64_encode":
				toolList.WriteString("- base64_encode(text): Encode text to base64\n")
			case "hash_sha256":
				toolList.WriteString("- hash_sha256(text): Generate SHA256 hash\n")
			default:
				toolList.WriteString(fmt.Sprintf("- %s(text): %s tool\n", name, name))
			}
		}
	}

	prompt := fmt.Sprintf(`You are an AI assistant with access to tools. When the user asks for something, you should:
1. Identify which tools are needed
2. Call the tools using the format: TOOL_CALL:tool_name:parameters
3. Provide a helpful response

%s

User request: %s

Please identify which tools you need and call them using the TOOL_CALL format, then provide a summary.`, toolList.String(), query)

	// Use the simple generate API with our enhanced prompt
	payload := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("� Sending prompt-based request to Ollama")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", ollamaURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaChunk
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	response := ollamaResp.Response
	log.Printf("💬 Ollama response: %s", response)

	// Parse the response for tool calls
	return parseAndExecuteToolCalls(response, tools, memory), nil
}

// parseAndExecuteToolCalls parses the LLM response for tool calls and executes them
func parseAndExecuteToolCalls(response string, tools *mcp.ToolRegistry, memory *mcp.Memory) string {
	if tools == nil {
		return response
	}

	var result strings.Builder
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "TOOL_CALL:") {
			// Parse tool call: TOOL_CALL:tool_name:parameters
			parts := strings.SplitN(line, ":", 3)
			if len(parts) >= 2 {
				toolName := parts[1]
				var params map[string]interface{}

				if len(parts) == 3 {
					params = map[string]interface{}{"text": parts[2]}
				} else {
					params = map[string]interface{}{}
				}

				log.Printf("🔧 Executing parsed tool call: %s with params %+v", toolName, params)

				toolResult, err := tools.Call(toolName, params, memory)
				if err != nil {
					log.Printf("❌ Tool %s failed: %v", toolName, err)
					result.WriteString(fmt.Sprintf("[Tool %s error: %v]\n", toolName, err))
				} else {
					log.Printf("✅ Tool %s succeeded: %v", toolName, toolResult)
					result.WriteString(fmt.Sprintf("[Tool %s result: %v]\n", toolName, toolResult))
				}
			}
		} else {
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}

// CreateSimpleModel creates a simple echo model for testing
func CreateSimpleModel() mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])
		response := fmt.Sprintf("Echo: %s", query)

		if onToken != nil {
			onToken(ctx.ContextID, response)
		}

		return response, nil
	}
}

// CreateCustomModel creates a model from a custom function
func CreateCustomModel(modelFunc func(query string, memory *mcp.Memory) (string, error)) mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		response, err := modelFunc(query, memory)
		if err != nil {
			return "", err
		}

		if onToken != nil {
			onToken(ctx.ContextID, response)
		}

		return response, nil
	}
}

// sendToolResultsToOllama sends tool results back to Ollama and gets the final response
func sendToolResultsToOllama(ollamaURL, model, originalQuery string, messages []OllamaChatMessage, onToken mcp.StreamCallback, contextID string) (string, error) {
	// Create a new chat request with the conversation history including tool results
	payload := OllamaChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Tools:    nil, // Don't include tools in follow-up request to avoid infinite loops
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal follow-up request: %w", err)
	}

	log.Printf("🔄 Sending follow-up request to Ollama with tool results")
	log.Printf("📤 Follow-up payload: %s", string(body))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", ollamaURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create follow-up request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama follow-up API: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("📡 Follow-up response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama follow-up API returned status %d", resp.StatusCode)
	}

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read follow-up response: %w", err)
	}

	log.Printf("📦 Follow-up raw response: %s", string(respBody))

	var chatResp OllamaChatChunk
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to decode follow-up response: %w", err)
	}

	// Return the final response content
	if chatResp.Message.Content != "" {
		log.Printf("💬 Final response from Ollama: %s", chatResp.Message.Content)
		return chatResp.Message.Content, nil
	}

	return "", fmt.Errorf("ollama returned empty follow-up response")
}

// OpenAIRequest represents a request to OpenAI-compatible APIs (like DeepInfra)
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Stream      bool            `json:"stream"`
}

// OpenAIMessage represents a message in OpenAI format
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a response from OpenAI-compatible APIs
type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
}

// OpenAIChoice represents a choice in the response
type OpenAIChoice struct {
	Message OpenAIMessage `json:"message"`
}

// CreateOpenAICompatibleModel creates a model function for OpenAI-compatible APIs like DeepInfra
func CreateOpenAICompatibleModel(apiURL, bearerToken string) mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		// Use provided model or default
		model := req.Model
		if model == "" {
			model = "meta-llama/Meta-Llama-3.1-8B-Instruct" // Default DeepInfra model
		}

		log.Printf("🔍 CreateOpenAICompatibleModel called with query: %s", query)
		log.Printf("🔧 API URL: %s", apiURL)
		log.Printf("🔧 Model: %s", model)

		// Create OpenAI-compatible request
		payload := OpenAIRequest{
			Model: model,
			Messages: []OpenAIMessage{
				{
					Role:    "user",
					Content: query,
				},
			},
			Temperature: req.Temperature,
			MaxTokens:   1000, // Default max tokens
			Stream:      false,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(context.Background(), "POST", apiURL, bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers for OpenAI-compatible API
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+bearerToken)

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			return "", fmt.Errorf("failed to call API: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var result OpenAIResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", fmt.Errorf("failed to decode response: %w", err)
		}

		if len(result.Choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		response := result.Choices[0].Message.Content

		// Call token callback if provided
		if onToken != nil {
			onToken(ctx.ContextID, response)
		}

		log.Printf("✅ OpenAI-compatible API response received: %d characters", len(response))
		log.Printf("💬 Response content: %s", response) // Log first 100 chars for brevity
		return response, nil
	}
}

// CreateDeepInfraModel creates a model function specifically for DeepInfra
func CreateDeepInfraModel(bearerToken string) mcp.ModelFunc {
	return CreateOpenAICompatibleModel("https://api.deepinfra.com/v1/openai/chat/completions", bearerToken)
}

// CreateModelFunction creates a model function from configuration
func CreateModelFunction(config interface{}) (mcp.ModelFunc, error) {
	// Handle different config types
	switch cfg := config.(type) {
	case *ModelConfig:
		return CreateModelFunctionFromConfig(cfg)
	default:
		return nil, fmt.Errorf("unsupported config type: %T", config)
	}
}

// ModelConfig holds configuration for individual models (defined in swarm package)
// We need to import it to avoid circular dependency, so we'll define a local version
type ModelConfig struct {
	Provider    string  `json:"provider"`
	Model       string  `json:"model"`
	URL         string  `json:"url,omitempty"`
	APIKey      string  `json:"api_key,omitempty"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	TopK        int     `json:"top_k"`
}

// CreateModelFunctionFromConfig creates a model function from ModelConfig
func CreateModelFunctionFromConfig(config *ModelConfig) (mcp.ModelFunc, error) {
	if config == nil {
		return nil, fmt.Errorf("model config is nil")
	}

	switch strings.ToLower(config.Provider) {
	case "ollama":
		return CreateOllamaModelWithConfig(config), nil
	case "openai":
		return CreateOpenAIModelWithConfig(config), nil
	case "deepinfra":
		return CreateDeepInfraModelWithConfig(config), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}

// CreateOllamaModelWithConfig creates an Ollama model function with configuration
func CreateOllamaModelWithConfig(config *ModelConfig) mcp.ModelFunc {
	ollamaURL := config.URL
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		payload := OllamaRequest{
			Model:  config.Model,
			Prompt: query,
			Stream: false,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(context.Background(), "POST", ollamaURL+"/api/generate", bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 300 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			return "", fmt.Errorf("failed to call Ollama: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("Ollama returned status %d", resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response: %w", err)
		}

		var ollamaResp OllamaChunk
		if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
			return "", fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return ollamaResp.Response, nil
	}
}

// CreateOpenAIModelWithConfig creates an OpenAI model function with configuration
func CreateOpenAIModelWithConfig(config *ModelConfig) mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		// OpenAI API implementation
		payload := map[string]interface{}{
			"model": config.Model,
			"messages": []map[string]string{
				{"role": "user", "content": query},
			},
			"temperature": config.Temperature,
			"max_tokens":  config.MaxTokens,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal OpenAI request: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(context.Background(), "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("failed to create OpenAI request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+config.APIKey)

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			return "", fmt.Errorf("failed to call OpenAI: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("OpenAI returned status %d", resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read OpenAI response: %w", err)
		}

		var openaiResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(respBody, &openaiResp); err != nil {
			return "", fmt.Errorf("failed to unmarshal OpenAI response: %w", err)
		}

		if len(openaiResp.Choices) == 0 {
			return "", fmt.Errorf("no choices in OpenAI response")
		}

		return openaiResp.Choices[0].Message.Content, nil
	}
}

// CreateDeepInfraModelWithConfig creates a DeepInfra model function with configuration
func CreateDeepInfraModelWithConfig(config *ModelConfig) mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		// DeepInfra API implementation (similar to OpenAI format)
		payload := map[string]interface{}{
			"model": config.Model,
			"messages": []map[string]string{
				{"role": "user", "content": query},
			},
			"temperature": config.Temperature,
			"max_tokens":  config.MaxTokens,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal DeepInfra request: %w", err)
		}

		url := "https://api.deepinfra.com/v1/openai/chat/completions"
		if config.URL != "" {
			url = config.URL
		}

		httpReq, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("failed to create DeepInfra request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+config.APIKey)

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			return "", fmt.Errorf("failed to call DeepInfra: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("DeepInfra returned status %d", resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read DeepInfra response: %w", err)
		}

		var deepinfraResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(respBody, &deepinfraResp); err != nil {
			return "", fmt.Errorf("failed to unmarshal DeepInfra response: %w", err)
		}

		if len(deepinfraResp.Choices) == 0 {
			return "", fmt.Errorf("no choices in DeepInfra response")
		}

		return deepinfraResp.Choices[0].Message.Content, nil
	}
}
