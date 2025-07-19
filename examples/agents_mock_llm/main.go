package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	fmt.Println("🧠 AI Agents with Mock LLM Integration Demo")
	fmt.Println("===========================================")
	fmt.Println("This demo shows how agents integrate with LLMs")
	fmt.Println("(Using mock LLM for demonstration purposes)")

	// Create MCP server
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8089
	config.EnableLogging = true

	server := conduit.NewEnhancedServer(config)

	// Register tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Add math tools
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a + b
			fmt.Printf("   🧮 Tool executed: %.1f + %.1f = %.1f\n", a, b, result)
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
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a * b
			fmt.Printf("   🧮 Tool executed: %.1f × %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "multiplication",
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number"),
			"b": conduit.NumberParam("Second number"),
		}, []string{"a", "b"}))

	// Create a mock LLM that demonstrates intelligent reasoning
	mockLLM := createMockLLM()

	// Create LLM-powered agent manager
	llmAgentManager := agents.NewLLMAgentManager(server, mockLLM, "mock-llm-v1")

	// Create an intelligent agent
	agent, err := llmAgentManager.CreateLLMAgent(
		"intelligent_agent",
		"Intelligent Problem Solver",
		"An agent that uses LLM reasoning to solve problems intelligently",
		`You are an intelligent problem-solving assistant. You analyze problems, reason about them step by step, and use the appropriate tools to solve them effectively.`,
		[]string{"add", "multiply", "word_count", "remember", "uuid"},
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

	fmt.Printf("✅ Created: %s\n", agent.Name)

	// Start server
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)

	// Demonstrate different types of LLM-powered reasoning
	fmt.Println("\n🧪 LLM-Powered Agent Demonstrations")
	fmt.Println("===================================")

	// Demo 1: Mathematical reasoning
	fmt.Println("\n🧮 Demo 1: Mathematical Problem Solving")
	fmt.Println("Problem: Calculate the total cost of 3 items at $15 each")

	mathTask, _ := llmAgentManager.CreateTask(
		"intelligent_agent",
		"Calculate Total Cost",
		"Calculate total cost for multiple items",
		map[string]interface{}{
			"problem": "Calculate the total cost of 3 items at $15 each",
			"items":   3.0,
			"price":   15.0,
		},
	)

	fmt.Println("🧠 LLM analyzing the problem...")
	if err := llmAgentManager.ExecuteTaskWithLLM(mathTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Problem solved successfully!")
		printTaskSteps(mathTask)
	}

	// Demo 2: Text analysis reasoning
	fmt.Println("\n📝 Demo 2: Text Analysis with Context")
	fmt.Println("Text: 'Artificial Intelligence is transforming technology'")

	textTask, _ := llmAgentManager.CreateTask(
		"intelligent_agent",
		"Analyze Important Text",
		"Analyze text and extract key insights",
		map[string]interface{}{
			"text": "Artificial Intelligence is transforming technology",
			"task": "analyze_and_store",
		},
	)

	fmt.Println("🧠 LLM analyzing the text...")
	if err := llmAgentManager.ExecuteTaskWithLLM(textTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Text analyzed successfully!")
		printTaskSteps(textTask)
	}

	fmt.Println("\n🎓 Key LLM Integration Concepts Demonstrated:")
	fmt.Println("==============================================")
	fmt.Println("✅ LLM analyzes tasks and creates intelligent plans")
	fmt.Println("✅ Agents use LLM reasoning to choose appropriate tools")
	fmt.Println("✅ Context-aware decision making based on task requirements")
	fmt.Println("✅ Multi-step problem decomposition and execution")
	fmt.Println("✅ Error recovery and adaptive planning")

	fmt.Println("\n🔧 How to Use with Real Ollama:")
	fmt.Println("===============================")
	fmt.Println("1. Install Ollama: https://ollama.ai")
	fmt.Println("2. Pull a model: ollama pull llama3.2")
	fmt.Println("3. Set environment variables:")
	fmt.Println("   export OLLAMA_URL=http://localhost:11434")
	fmt.Println("   export OLLAMA_MODEL=llama3.2")
	fmt.Println("4. Run: go run examples/agents_ollama/main.go")

	fmt.Printf("\n🔗 Mock server running on http://localhost:%d\n", config.Port)
}

// createMockLLM creates a mock LLM that demonstrates intelligent reasoning
func createMockLLM() mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		// Simulate LLM reasoning based on the query content
		if strings.Contains(query, "Calculate") && strings.Contains(query, "total cost") {
			// Mathematical reasoning response
			return `{
  "analysis": "This is a multiplication problem. I need to calculate 3 items × $15 each. I should use the multiply tool to get the result.",
  "steps": [
    {
      "name": "calculate_total",
      "description": "Multiply number of items by price per item",
      "tool": "multiply",
      "input": {"a": 3.0, "b": 15.0}
    },
    {
      "name": "store_result",
      "description": "Remember the calculation result",
      "tool": "remember",
      "input": {"key": "last_calculation", "value": "3 items × $15 = $45"}
    }
  ],
  "reasoning": "I identified this as a multiplication problem and planned to use the multiply tool followed by storing the result for reference."
}`, nil
		}

		if strings.Contains(query, "Analyze") && strings.Contains(query, "text") {
			// Text analysis reasoning response
			return `{
  "analysis": "This text is about AI and technology. I should analyze its content and store it for future reference since it seems important.",
  "steps": [
    {
      "name": "count_words",
      "description": "Count words in the text to understand its length",
      "tool": "word_count",
      "input": {"text": "Artificial Intelligence is transforming technology"}
    },
    {
      "name": "store_content",
      "description": "Store this important text about AI",
      "tool": "remember",
      "input": {"key": "ai_insight", "value": "Artificial Intelligence is transforming technology"}
    },
    {
      "name": "generate_id",
      "description": "Generate a unique ID for this analysis session",
      "tool": "uuid",
      "input": {}
    }
  ],
  "reasoning": "I analyzed the text content, decided to count words for metrics, store the important content, and generate a session ID for tracking."
}`, nil
		}

		// Default reasoning response
		return `{
  "analysis": "I need to analyze this task and determine the best approach using available tools.",
  "steps": [
    {
      "name": "general_analysis",
      "description": "Perform general analysis of the task",
      "tool": "uuid",
      "input": {}
    }
  ],
  "reasoning": "This appears to be a general task that requires basic analysis."
}`, nil
	}
}

// printTaskSteps prints the execution steps with LLM reasoning
func printTaskSteps(task *agents.Task) {
	fmt.Printf("  📋 Task: %s\n", task.Title)
	fmt.Printf("  📊 Status: %s\n", task.Status)

	for i, step := range task.Steps {
		fmt.Printf("    Step %d: %s\n", i+1, step.Name)

		if step.Name == "llm_reasoning" && step.Output != nil {
			fmt.Printf("      🧠 LLM provided intelligent analysis and planning\n")
		}

		if step.Output != nil && step.Status == "completed" {
			if output, ok := step.Output["output"]; ok {
				if outputMap, ok := output.(map[string]interface{}); ok {
					if result, ok := outputMap["result"]; ok {
						fmt.Printf("      ✅ Result: %v\n", result)
					}
				}
			}
		}
	}
}
