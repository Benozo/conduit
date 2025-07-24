package providers

import (
	"context"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

type AnthropicProvider struct {
	apiKey string
}

func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	return &AnthropicProvider{apiKey: apiKey}
}

func (p *AnthropicProvider) CallModel(ctx context.Context, prompt string) (string, error) {
	url := "https://api.anthropic.com/v1/complete"
	body, err := json.Marshal(map[string]interface{}{
		"prompt": prompt,
		"model":  "claude-v1",
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var response struct {
		Completion string `json:"completion"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Completion, nil
}