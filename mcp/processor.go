package mcp

import "fmt"

type ModelFunc func(ctx ContextInput, req MCPRequest, memory *Memory, onToken StreamCallback) (string, error)

type MCPProcessor struct {
	Model        ModelFunc
	Tools        *ToolRegistry
	Memory       *Memory
	StreamTokens bool
	OnToken      StreamCallback
}

func NewProcessor(model ModelFunc, tools *ToolRegistry) *MCPProcessor {
	return &MCPProcessor{
		Model:        model,
		Tools:        tools,
		Memory:       NewMemory(),
		StreamTokens: false,
		OnToken:      nil,
	}
}

func (p *MCPProcessor) EnableStreaming(cb StreamCallback) {
	p.StreamTokens = true
	p.OnToken = cb
}

func (p *MCPProcessor) Run(req MCPRequest) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	for _, ctx := range req.Contexts {
		var out interface{}
		var err error

		if req.ToolChoice != nil {
			out, err = p.Tools.Call(req.ToolChoice.Name, req.ToolChoice.Parameters, p.Memory)
		} else {
			onToken := func(contextID, token string) {}
			if p.StreamTokens && p.OnToken != nil {
				onToken = func(contextID, token string) {
					p.OnToken(contextID, token)
				}
			}
			out, err = p.Model(ctx, req, p.Memory, onToken)
		}

		if err != nil {
			return nil, fmt.Errorf("context %s error: %w", ctx.ContextID, err)
		}
		results[ctx.ContextID] = out
	}

	return results, nil
}