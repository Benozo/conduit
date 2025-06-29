package conduit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
			Stream: true,
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
