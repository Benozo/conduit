package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OllamaEmbeddings implements EmbeddingProvider for Ollama
type OllamaEmbeddings struct {
	client     *http.Client
	baseURL    string
	model      string
	dimensions int
	timeout    time.Duration
}

// OllamaEmbeddingRequest represents the request to Ollama embedding API
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse represents the response from Ollama embedding API
type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// NewOllamaEmbeddings creates a new Ollama embeddings provider
func NewOllamaEmbeddings(host, model string, dimensions int, timeout time.Duration) *OllamaEmbeddings {
	if host == "" {
		host = "localhost"
	}

	baseURL := fmt.Sprintf("http://%s:11434", host)

	return &OllamaEmbeddings{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL:    baseURL,
		model:      model,
		dimensions: dimensions,
		timeout:    timeout,
	}
}

// Embed generates embedding for a single text
func (o *OllamaEmbeddings) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Create context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	// Prepare request
	reqBody := OllamaEmbeddingRequest{
		Model:  o.model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/embeddings", o.baseURL)
	req, err := http.NewRequestWithContext(ctxWithTimeout, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	// Parse response
	var embeddingResp OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(embeddingResp.Embedding) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return embeddingResp.Embedding, nil
}

// EmbedBatch generates embeddings for multiple texts
// Note: Ollama doesn't support batch requests, so we'll process them sequentially
func (o *OllamaEmbeddings) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts cannot be empty")
	}

	results := make([][]float32, len(texts))

	for i, text := range texts {
		if text == "" {
			// Skip empty texts but maintain array position
			results[i] = nil
			continue
		}

		embedding, err := o.Embed(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to embed text %d: %w", i, err)
		}

		results[i] = embedding
	}

	return results, nil
}

// GetDimensions returns the embedding dimensions
func (o *OllamaEmbeddings) GetDimensions() int {
	return o.dimensions
}

// GetModel returns the model name
func (o *OllamaEmbeddings) GetModel() string {
	return o.model
}

// GetProvider returns the provider name
func (o *OllamaEmbeddings) GetProvider() string {
	return "ollama"
}

// Ping checks if the Ollama API is accessible
func (o *OllamaEmbeddings) Ping(ctx context.Context) error {
	// Test with a simple embedding request
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// First check if Ollama is running
	url := fmt.Sprintf("%s/api/tags", o.baseURL)
	req, err := http.NewRequestWithContext(ctxWithTimeout, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create ping request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("Ollama API ping failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	// Test embedding with a simple text
	_, err = o.Embed(ctxWithTimeout, "test")
	if err != nil {
		return fmt.Errorf("Ollama embedding test failed: %w", err)
	}

	return nil
}

// GetAvailableModels retrieves the list of available models from Ollama
func (o *OllamaEmbeddings) GetAvailableModels(ctx context.Context) ([]string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/api/tags", o.baseURL)
	req, err := http.NewRequestWithContext(ctxWithTimeout, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var tagsResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]string, len(tagsResp.Models))
	for i, model := range tagsResp.Models {
		models[i] = model.Name
	}

	return models, nil
}

// PullModel pulls a model from Ollama if it's not already available
func (o *OllamaEmbeddings) PullModel(ctx context.Context, model string) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute) // Pulling can take time
	defer cancel()

	reqBody := map[string]string{
		"name": model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/pull", o.baseURL)
	req, err := http.NewRequestWithContext(ctxWithTimeout, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pull model, status: %d", resp.StatusCode)
	}

	return nil
}
