package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

// extractNumber safely extracts a number from an interface{} value
func extractNumber(value interface{}) (float64, error) {
	if value == nil {
		return 0, fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case string:
		// Try to parse string as number
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, nil
		}
		return 0, fmt.Errorf("cannot parse string '%s' as number", v)
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}

func main() {
	// Get DeepInfra configuration from environment
	bearerToken := os.Getenv("DEEPINFRA_TOKEN")
	if bearerToken == "" {
		fmt.Println("❌ DEEPINFRA_TOKEN environment variable is required")
		fmt.Println("   Set it with: export DEEPINFRA_TOKEN=your_token_here")
		return
	}

	modelName := os.Getenv("DEEPINFRA_MODEL")
	if modelName == "" {
		modelName = "meta-llama/Meta-Llama-3.1-8B-Instruct" // Default model
	}

	fmt.Printf("🤖 Model: %s\n", modelName)

	// Create MCP server for library mode
	config := conduit.DefaultConfig()
	config.EnableLogging = true // Keep logging to see LLM interactions

	server := conduit.NewEnhancedServer(config)

	// Register tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			// Robust number extraction
			a, err := extractNumber(params["a"])
			if err != nil {
				return nil, fmt.Errorf("invalid parameter 'a': %v", err)
			}
			b, err := extractNumber(params["b"])
			if err != nil {
				return nil, fmt.Errorf("invalid parameter 'b': %v", err)
			}

			result := a + b
			fmt.Printf("🧮 ADD: %.1f + %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "addition",
			}, nil
		},
		conduit.CreateToolMetadata("add", "Add two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number"),
			"b": conduit.NumberParam("Second number"),
		}, []string{"a", "b"}))

	server.RegisterToolWithSchema("multiply",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			// Robust number extraction
			a, err := extractNumber(params["a"])
			if err != nil {
				return nil, fmt.Errorf("invalid parameter 'a': %v", err)
			}
			b, err := extractNumber(params["b"])
			if err != nil {
				return nil, fmt.Errorf("invalid parameter 'b': %v", err)
			}

			result := a * b
			fmt.Printf("🧮 MULTIPLY: %.1f × %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "multiplication",
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number"),
			"b": conduit.NumberParam("Second number"),
		}, []string{"a", "b"}))

	// Create DeepInfra model function
	deepInfraModel := conduit.CreateDeepInfraModel(bearerToken)

	// Create LLM-powered agent manager
	llmAgentManager := agents.NewLLMAgentManager(server, deepInfraModel, modelName)

	// Create an intelligent agent
	_, err := llmAgentManager.CreateLLMAgent(
		"deepinfra_agent",
		"DeepInfra Problem Solver",
		"An agent powered by DeepInfra's LLM for intelligent reasoning and problem solving",
		`You are an intelligent assistant that can analyze problems, break them down into steps, and use available tools to solve them. When given a task, think step by step and create a JSON plan with the tools you need to use.

Available tools:
- add: Add two numbers (parameters: a, b as numbers)
- multiply: Multiply two numbers (parameters: a, b as numbers)
- word_count: Count words in text (parameter: text as string)
- remember: Store information in memory (parameters: key, value as strings)
- recall: Retrieve stored information (parameter: key as string)
- uuid: Generate unique IDs (no parameters)

IMPORTANT RULES:
1. When using math tools (add, multiply), always provide ACTUAL numeric values, NOT references to previous steps
2. If you need to use a result from a previous step, use the actual number that would result from that step
3. Plan all calculations so you can provide concrete numbers for each step

For multi-step calculations:
- Step 1: Calculate the first operation with concrete numbers
- Step 2: Calculate the next operation using the RESULT from step 1 (provide the actual number)

Format your response as JSON with:
{
  "analysis": "Your analysis of the problem",
  "steps": [
    {
      "name": "step_name",
      "description": "What this step does",
      "tool": "tool_name",
      "input": {"param": value}
    }
  ],
  "reasoning": "Your reasoning for this approach"
}

CORRECT example for multi-step math:
{
  "steps": [
    {"name": "calc_weekly", "tool": "multiply", "input": {"a": 40, "b": 50}},
    {"name": "calc_team", "tool": "multiply", "input": {"a": 5, "b": 2000}},
    {"name": "calc_monthly", "tool": "multiply", "input": {"a": 10000, "b": 4.3}}
  ]
}

WRONG example (DO NOT USE):
{"tool": "multiply", "input": {"a": 5, "b": "result of step 1"}}`,
		[]string{"add", "multiply", "word_count", "remember", "recall", "uuid"},
		&agents.AgentConfig{
			MaxTokens:     1000,
			Temperature:   0.3,
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Demo 1: Mathematical reasoning with real LLM
	mathTask, _ := llmAgentManager.CreateTask(
		"deepinfra_agent",
		"Restaurant Revenue Calculation",
		"Calculate hourly revenue from table service",
		map[string]interface{}{
			"problem":         "A restaurant serves 8 tables per hour. If each table pays $25 on average, how much revenue per hour?",
			"tables_per_hour": 8.0,
			"avg_payment":     25.0,
		},
	)

	fmt.Println("🧠 DeepInfra LLM is analyzing the problem...")
	start := time.Now()
	if err := llmAgentManager.ExecuteTaskWithLLM(mathTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		elapsed := time.Since(start)
		fmt.Printf("✅ Problem solved in %v!\n", elapsed)
	}

	// Demo 2: Text analysis with real LLM
	sampleText := "DeepInfra provides access to cutting-edge language models through a simple API, enabling developers to build intelligent applications with powerful AI capabilities."

	textTask, _ := llmAgentManager.CreateTask(
		"deepinfra_agent",
		"DeepInfra Text Analysis",
		"Analyze text about DeepInfra and extract insights",
		map[string]interface{}{
			"text": sampleText,
			"task": "analyze_and_store",
		},
	)

	fmt.Println("🧠 DeepInfra LLM is analyzing the text...")
	start = time.Now()
	if err := llmAgentManager.ExecuteTaskWithLLM(textTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		elapsed := time.Since(start)
		fmt.Printf("✅ Text analyzed in %v!\n", elapsed)
	}

	// Demo 3: Complex reasoning task
	complexTask, _ := llmAgentManager.CreateTask(
		"deepinfra_agent",
		"Team Cost Calculation",
		"Calculate monthly team cost with multiple steps",
		map[string]interface{}{
			"problem":         "A software team of 5 developers works 40 hours/week at $50/hour. What's the monthly cost?",
			"developers":      5.0,
			"hours_per_week":  40.0,
			"hourly_rate":     50.0,
			"weeks_per_month": 4.3, // Average
		},
	)

	fmt.Println("🧠 DeepInfra LLM is working on the complex problem...")
	start = time.Now()
	if err := llmAgentManager.ExecuteTaskWithLLM(complexTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		elapsed := time.Since(start)
		fmt.Printf("✅ Complex problem solved in %v!\n", elapsed)
	}

	fmt.Println("\n🎉 DeepInfra demonstration completed!")
}
