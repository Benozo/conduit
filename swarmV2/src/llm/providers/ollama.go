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

// Ensure OllamaProvider implements LanguageModelProvider interface
var _ llm.LanguageModelProvider = (*OllamaProvider)(nil)

// OllamaProvider represents an Ollama language model provider.
type OllamaProvider struct {
	BaseURL string // Ollama server URL
	Model   string // Current model name
	Client  *http.Client
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents a response from Ollama API
type OllamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// NewOllamaProvider initializes a new OllamaProvider instance.
func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	return &OllamaProvider{
		BaseURL: baseURL,
		Model:   model,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate generates a response based on the provided prompt using Ollama.
func (op *OllamaProvider) Generate(prompt string) (string, error) {
	return op.GenerateResponse(prompt)
}

// GenerateResponse implements the LanguageModelProvider interface.
func (op *OllamaProvider) GenerateResponse(prompt string) (string, error) {
	return op.CallOllama(context.Background(), prompt)
}

// SetModel sets the model name for the Ollama provider.
func (op *OllamaProvider) SetModel(model string) {
	op.Model = model
}

// CallOllama makes a request to the Ollama API
func (op *OllamaProvider) CallOllama(ctx context.Context, prompt string) (string, error) {
	// Prepare the request
	reqBody := OllamaRequest{
		Model:  op.Model,
		Prompt: prompt,
		Stream: false, // We want a single response, not streaming
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/generate", op.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := op.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, resp.Status)
	}

	// Parse the response
	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}

// GetModelInfo returns information about the Ollama model.
func (op *OllamaProvider) GetModelInfo() llm.ModelInfo {
	return llm.ModelInfo{
		Name:        op.Model,
		Version:     "latest",
		Provider:    "Ollama",
		MaxTokens:   4096, // Default, varies by model
		Description: fmt.Sprintf("Ollama model %s running at %s", op.Model, op.BaseURL),
	}
}

// Ping checks if the Ollama server is reachable
func (op *OllamaProvider) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", op.BaseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create ping request: %w", err)
	}

	resp, err := op.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to ping Ollama server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama server responded with status %d", resp.StatusCode)
	}

	return nil
}
