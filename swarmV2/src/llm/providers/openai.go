package providers

import (
	"context"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

type OpenAIProvider struct {
	APIKey string
}

type OpenAIRequest struct {
	Model string `json:"model"`
	Prompt string `json:"prompt"`
	MaxTokens int `json:"max_tokens"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{APIKey: apiKey}
}

func (o *OpenAIProvider) CallModel(ctx context.Context, prompt string, model string, maxTokens int) (string, error) {
	requestBody := OpenAIRequest{
		Model: model,
		Prompt: prompt,
		MaxTokens: maxTokens,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/engines/"+model+"/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API returned non-200 status: %s", resp.Status)
	}

	var openAIResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI API")
	}

	return openAIResponse.Choices[0].Text, nil
}