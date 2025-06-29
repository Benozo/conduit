package mcp

import "fmt"

type AgentStep struct {
	Thought  string
	Action   string
	Params   map[string]interface{}
	Observed string
}

func ReActAgent(thoughts []string, tools *ToolRegistry, memory *Memory) ([]AgentStep, error) {
	var steps []AgentStep
	for _, t := range thoughts {
		step := AgentStep{Thought: t}

		if t == "transform to uppercase" {
			step.Action = "uppercase"
			step.Params = map[string]interface{}{"text": memory.Get("latest")}
			res, err := tools.Call(step.Action, step.Params, memory)
			if err != nil {
				return steps, err
			}
			step.Observed = fmt.Sprintf("%v", res)
			memory.Set("latest_result", step.Observed)
		} else {
			step.Observed = "no action"
		}

		steps = append(steps, step)
	}

	return steps, nil
}