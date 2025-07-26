package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/benozo/neuron/src/llm"
)

// Ensure CloudflareAIProvider implements LanguageModelProvider interface
var _ llm.LanguageModelProvider = (*CloudflareAIProvider)(nil)

// CloudflareAIProvider represents a Cloudflare Workers AI provider.
type CloudflareAIProvider struct {
	AccountID string // Cloudflare Account ID (for standard API)
	APIToken  string // Cloudflare API Token (for standard API)
	APIKey    string // Custom X-API-Key for custom endpoints
	BaseURL   string // Custom base URL for API endpoint
	Model     string // Current model name
	Client    *http.Client
	UseCustom bool // Whether to use custom endpoint or standard Cloudflare API
}

// CloudflareAIRequest represents a request to Cloudflare Workers AI API
type CloudflareAIRequest struct {
	Model       string              `json:"model"`
	Messages    []CloudflareMessage `json:"messages"`
	Stream      bool                `json:"stream,omitempty"`
	Temperature float64             `json:"temperature,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
	Tools       []CloudflareTool    `json:"tools,omitempty"`
}

// CloudflareTool represents a tool/function definition
type CloudflareTool struct {
	Type     string             `json:"type"`
	Function CloudflareFunction `json:"function"`
}

// CloudflareFunction represents a function definition
type CloudflareFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

// CloudflareMessage represents a message in the Cloudflare AI format
type CloudflareMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CloudflareAIResponse represents a response from Cloudflare Workers AI
type CloudflareAIResponse struct {
	Result    CloudflareResult  `json:"result,omitempty"`
	Response  string            `json:"response,omitempty"` // For custom endpoint
	Success   bool              `json:"success"`
	Errors    []CloudflareError `json:"errors,omitempty"`
	ToolCalls []interface{}     `json:"tool_calls,omitempty"`
	Usage     interface{}       `json:"usage,omitempty"`
}

// CloudflareResult contains the AI response
type CloudflareResult struct {
	Response string `json:"response"`
}

// CloudflareError represents an error from Cloudflare API
type CloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewCloudflareAIProvider creates a new Cloudflare Workers AI provider.
func NewCloudflareAIProvider(accountID, apiToken, model string) *CloudflareAIProvider {
	return &CloudflareAIProvider{
		AccountID: accountID,
		APIToken:  apiToken,
		Model:     model,
		UseCustom: false,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewCustomCloudflareAIProvider creates a new provider for custom Cloudflare endpoint.
func NewCustomCloudflareAIProvider(baseURL, apiKey, model string) *CloudflareAIProvider {
	return &CloudflareAIProvider{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		Model:     model,
		UseCustom: true,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateResponse generates a response using Cloudflare Workers AI.
func (c *CloudflareAIProvider) GenerateResponse(prompt string) (string, error) {
	var url string

	if c.UseCustom {
		// Use custom endpoint
		url = c.BaseURL + "/v1/chat/completions"
	} else {
		// Use standard Cloudflare API
		url = fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", c.AccountID, c.Model)
	}

	// Prepare the request
	cfRequest := CloudflareAIRequest{
		Model: c.Model,
		Messages: []CloudflareMessage{
			{
				Role:    "assitant", // Use exact typo from your working curl
				Content: "You are a helpful assistant",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream:      true, // Set to true to match your curl example
		Temperature: 0.7,
		MaxTokens:   150,
		Tools: []CloudflareTool{
			{
				Type: "function",
				Function: CloudflareFunction{
					Name:        "getWeather",
					Description: "Get current weather in a city",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"city": map[string]interface{}{
								"type":        "string",
								"description": "The city name",
							},
						},
						"required": []string{"city"},
					},
				},
			},
		},
	}

	requestBody, err := json.Marshal(cfRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "insomnia/11.3.0") // Match your working curl

	if c.UseCustom {
		// Use X-API-Key for custom endpoint
		req.Header.Set("X-API-Key", c.APIKey)
	} else {
		// Use Bearer token for standard Cloudflare API
		req.Header.Set("Authorization", "Bearer "+c.APIToken)
	}

	// Send the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse the response
	var cfResponse CloudflareAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Handle custom endpoint response format
	if c.UseCustom {
		if cfResponse.Response != "" {
			return cfResponse.Response, nil
		}
		// If response field is empty, check if there's an error
		if len(cfResponse.Errors) > 0 {
			return "", fmt.Errorf("custom endpoint error: %v", cfResponse.Errors)
		}
		return "", fmt.Errorf("no response received from custom endpoint, raw response: %+v", cfResponse)
	}

	// Handle standard Cloudflare API response format
	if !cfResponse.Success {
		if len(cfResponse.Errors) > 0 {
			return "", fmt.Errorf("cloudflare AI error: %s", cfResponse.Errors[0].Message)
		}
		return "", fmt.Errorf("cloudflare AI request failed")
	}

	return cfResponse.Result.Response, nil
}

// SetModel sets the model for the provider.
func (c *CloudflareAIProvider) SetModel(model string) {
	c.Model = model
}

// GetModelInfo returns information about the current model.
func (c *CloudflareAIProvider) GetModelInfo() llm.ModelInfo {
	provider := "Cloudflare Workers AI"
	if c.UseCustom {
		provider = "Custom Cloudflare AI"
	}
	return llm.ModelInfo{
		Name:        c.Model,
		Provider:    provider,
		Description: fmt.Sprintf("%s model: %s", provider, c.Model),
	}
}

// Ping tests the connection to Cloudflare Workers AI.
func (c *CloudflareAIProvider) Ping() error {
	// Use a simple prompt to test connectivity
	_, err := c.GenerateResponse("Hello")
	return err
}

// GetAvailableModels returns a list of available Cloudflare AI models.
func (c *CloudflareAIProvider) GetAvailableModels() []string {
	// These are some popular Cloudflare Workers AI models
	return []string{
		"@cf/meta/llama-3.1-8b-instruct",
		"@cf/meta/llama-3.1-70b-instruct",
		"@cf/meta/llama-3-8b-instruct",
		"@cf/meta/llama-4-scout-17b-16e-instruct", // Added based on your example
		"@cf/mistral/mistral-7b-instruct-v0.1",
		"@cf/microsoft/phi-2",
		"@cf/qwen/qwen1.5-7b-chat-awq",
		"@cf/google/gemma-7b-it",
		"@cf/openchat/openchat-3.5-0106",
		"@cf/meta/llama-2-7b-chat-int8",
		"@cf/thebloke/neural-chat-7b-v3-1-awq",
	}
}
