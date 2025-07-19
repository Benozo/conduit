package conduit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/benozo/conduit/mcp"
)

type OpenAIMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
	Name       string           `json:"name,omitempty"`
	ToolCalls  []OpenAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

type OpenAIToolCall struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Function OpenAIFunctionCall `json:"function"`
}

type OpenAIFunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type OpenAIFunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Type        string                 `json:"type,omitempty"` // "function" for OpenAI
}

type OpenAIRequest struct {
	Model      string              `json:"model"`
	Messages   []OpenAIMessage     `json:"messages"`
	Stream     bool                `json:"stream"`
	Functions  []OpenAIFunctionDef `json:"functions,omitempty"`
	ToolChoice interface{}         `json:"tool_choice,omitempty"` // "auto", "none", or map
	Tools      []OpenAIToolWrapper `json:"tools,omitempty"`       // Deprecated, use Functions
	Type       string              `json:"type,omitempty"`        // "chat.completions"
}

type OpenAIToolWrapper struct {
	Type     string            `json:"type"` // "function"
	Function OpenAIFunctionDef `json:"function"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
}

func CreateOpenAIToolAwareModel(apiKey, openaiURL string, tools *mcp.ToolRegistry) mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])
		model := req.Model
		if model == "" {
			model = "gpt-4o-mini"
		}

		log.Printf("🧠 OpenAI ToolAwareModel query: %s", query)

		messages := []OpenAIMessage{
			{Role: "user", Content: query},
		}

		functions := []OpenAIFunctionDef{}
		if tools != nil {
			for _, name := range tools.GetRegisteredTools() {
				fn := OpenAIFunctionDef{
					Name:        name,
					Description: fmt.Sprintf("Tool: %s", name),
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"text": map[string]interface{}{
								"type":        "string",
								"description": "Input text",
							},
						},
						"required": []string{"text"},
					},
					Type: "function",
				}

				switch name {
				case "remember":
					fn.Description = "Store a value in memory"
					fn.Parameters = map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"key":   map[string]interface{}{"type": "string", "description": "Memory key"},
							"value": map[string]interface{}{"type": "string", "description": "Value to store"},
						},
						"required": []string{"key", "value"},
					}
				case "recall", "forget":
					fn.Description = fmt.Sprintf("%s a memory value", name)
					fn.Parameters = map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"key": map[string]interface{}{"type": "string", "description": "Memory key"},
						},
						"required": []string{"key"},
					}
				case "add":
					fn.Description = "Add two numbers"
					fn.Parameters = map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"a": map[string]interface{}{"type": "number"},
							"b": map[string]interface{}{"type": "number"},
						},
						"required": []string{"a", "b"},
					}
				case "timestamp":
					fn.Description = "Get current timestamp"
					fn.Parameters = map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"format": map[string]interface{}{"type": "string"},
						},
					}
				}

				functions = append(functions, fn)
			}
		}

		payload := OpenAIRequest{
			Model:      model,
			Messages:   messages,
			Stream:     false,
			Tools:      buildToolWrappers(functions),
			ToolChoice: "auto",
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal: %w", err)
		}

		reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpReq, err := http.NewRequestWithContext(reqCtx, "POST", openaiURL+"/v1/chat/completions", bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("http request error: %w", err)
		}
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return "", fmt.Errorf("OpenAI call failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("OpenAI error %d: %s", resp.StatusCode, string(b))
		}

		var res OpenAIResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return "", fmt.Errorf("decode error: %w", err)
		}

		if len(res.Choices) == 0 {
			log.Printf("❌ OpenAI returned 0 choices in response")
			return "", fmt.Errorf("OpenAI returned no choices")
		}

		msg := res.Choices[0].Message

		if len(msg.ToolCalls) > 0 {
			log.Printf("🛠 OpenAI tool calls: %+v", msg.ToolCalls)

			var followup []OpenAIMessage

			for _, call := range msg.ToolCalls {
				log.Printf("🛠 Tool call received: %s", call.Function.Name)

				var args map[string]interface{}
				var raw string
				if err := json.Unmarshal(call.Function.Arguments, &raw); err == nil {
					if err := json.Unmarshal([]byte(raw), &args); err != nil {
						log.Printf("❌ Failed to parse tool arguments (wrapped): %v", err)
						continue
					}
				} else {
					if err := json.Unmarshal(call.Function.Arguments, &args); err != nil {
						log.Printf("❌ Failed to parse tool arguments (direct): %v", err)
						continue
					}
				}

				result := "unknown"
				if tools != nil {
					r, err := tools.Call(call.Function.Name, args, memory)
					if err != nil {
						result = fmt.Sprintf("error: %v", err)
					} else {
						result = fmt.Sprintf("%v", r)
					}
				}

				followup = append(followup, OpenAIMessage{
					Role:       "tool",
					Name:       call.Function.Name,
					Content:    result,
					ToolCallID: call.ID,
				})
			}

			assistantWithToolCalls := OpenAIMessage{
				Role:      "assistant",
				ToolCalls: msg.ToolCalls,
			}
			followupMessages := append(messages, assistantWithToolCalls)
			followupMessages = append(followupMessages, followup...)

			followupReq := OpenAIRequest{
				Model:      model,
				Messages:   followupMessages,
				Stream:     false,
				Tools:      buildToolWrappers(functions),
				ToolChoice: "none",
			}

			followupBody, err := json.Marshal(followupReq)
			if err != nil {
				return "", fmt.Errorf("failed to marshal follow-up: %w", err)
			}

			followupHTTPReq, err := http.NewRequestWithContext(reqCtx, "POST", openaiURL+"/v1/chat/completions", bytes.NewReader(followupBody))
			if err != nil {
				return "", fmt.Errorf("follow-up http request error: %w", err)
			}
			followupHTTPReq.Header.Set("Authorization", "Bearer "+apiKey)
			followupHTTPReq.Header.Set("Content-Type", "application/json")

			followupResp, err := http.DefaultClient.Do(followupHTTPReq)
			if err != nil {
				return "", fmt.Errorf("OpenAI follow-up error: %w", err)
			}
			defer followupResp.Body.Close()

			if followupResp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(followupResp.Body)
				log.Printf("❌ OpenAI follow-up error %d: %s", followupResp.StatusCode, string(b))
				return "", fmt.Errorf("OpenAI follow-up error %d: %s", followupResp.StatusCode, string(b))
			}

			// Read and log the raw response for debugging
			respBody, err := io.ReadAll(followupResp.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read follow-up response: %w", err)
			}
			log.Printf("📦 OpenAI follow-up raw response: %s", string(respBody))

			var followupRes OpenAIResponse
			if err := json.Unmarshal(respBody, &followupRes); err != nil {
				return "", fmt.Errorf("decode error (follow-up): %w", err)
			}
			if len(followupRes.Choices) == 0 {
				log.Printf("❌ OpenAI returned %d choices in follow-up response", len(followupRes.Choices))
				return "", fmt.Errorf("OpenAI returned no choices after tool call")
			}

			msg = followupRes.Choices[0].Message
		}

		content := msg.Content
		if onToken != nil {
			onToken(ctx.ContextID, content)
		}
		return content, nil
	}
}

func buildToolWrappers(funcDefs []OpenAIFunctionDef) []OpenAIToolWrapper {
	var wrappers []OpenAIToolWrapper
	for _, fn := range funcDefs {
		wrappers = append(wrappers, OpenAIToolWrapper{
			Type:     "function",
			Function: fn,
		})
	}
	return wrappers
}
