package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	fmt.Println("🧠 AI Agents with Mock LLM - Library Mode Demo")
	fmt.Println("===============================================")
	fmt.Println("This demo shows LLM-powered agents in library mode")
	fmt.Println("(No HTTP server - pure library usage with mock LLM)")

	// Create MCP core for library mode (no server)
	config := conduit.DefaultConfig()
	config.EnableLogging = false // Disable server logging for library mode

	server := conduit.NewEnhancedServer(config)

	// Register tools for library use
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

	// Create LLM-powered agent manager (library mode)
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

	fmt.Printf("✅ Created: %s (Library Mode)\n", agent.Name)

	// Demonstrate different types of LLM-powered reasoning in library mode
	fmt.Println("\n🧪 LLM-Powered Agent Library Demonstrations")
	fmt.Println("==========================================")

	// Demo 1: Mathematical reasoning
	fmt.Println("\n🧮 Demo 1: Mathematical Problem Solving")
	fmt.Println("Problem: Calculate the total cost of 3 items at $15 each")

	mathTask, err := llmAgentManager.CreateTask(
		"intelligent_agent",
		"Calculate Total Cost",
		"Calculate total cost for multiple items",
		map[string]interface{}{
			"problem": "Calculate the total cost of 3 items at $15 each",
			"items":   3.0,
			"price":   15.0,
		},
	)
	if err != nil {
		log.Printf("Failed to create math task: %v", err)
	} else {
		fmt.Println("🧠 LLM analyzing the problem...")
		if err := llmAgentManager.ExecuteTaskWithLLM(mathTask.ID); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Println("✅ Problem solved successfully!")
			printTaskSteps(mathTask)
		}
	}

	// Demo 2: Text analysis reasoning
	fmt.Println("\n📝 Demo 2: Text Analysis with Context")
	fmt.Println("Text: 'Artificial Intelligence is transforming technology'")

	textTask, err := llmAgentManager.CreateTask(
		"intelligent_agent",
		"Analyze Important Text",
		"Analyze text and extract key insights",
		map[string]interface{}{
			"text": "Artificial Intelligence is transforming technology",
			"task": "analyze_and_store",
		},
	)
	if err != nil {
		log.Printf("Failed to create text task: %v", err)
	} else {
		fmt.Println("🧠 LLM analyzing the text...")
		if err := llmAgentManager.ExecuteTaskWithLLM(textTask.ID); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Println("✅ Text analyzed successfully!")
			printTaskSteps(textTask)
		}
	}

	// Demo 3: Simple task execution (library style)
	fmt.Println("\n⚡ Demo 3: Simple Task Execution")
	fmt.Println("Creating and executing a custom calculation task")

	simpleTask, err := llmAgentManager.CreateTask(
		"intelligent_agent",
		"Simple Calculation",
		"Calculate 7 × 8 using available tools",
		map[string]interface{}{
			"operation": "Calculate 7 × 8 and store the result",
			"a":         7.0,
			"b":         8.0,
		},
	)
	if err != nil {
		fmt.Printf("❌ Task creation failed: %v\n", err)
	} else {
		fmt.Println("🧠 LLM planning the calculation...")
		if err := llmAgentManager.ExecuteTaskWithLLM(simpleTask.ID); err != nil {
			fmt.Printf("❌ Execution failed: %v\n", err)
		} else {
			fmt.Printf("✅ Calculation completed successfully!\n")
			printTaskSteps(simpleTask)
		}
	}

	// Show available tools and agent status
	fmt.Println("\n📊 Library Mode Status")
	fmt.Println("=====================")

	availableTools := llmAgentManager.GetAvailableTools()
	fmt.Printf("🔧 Available Tools: %v\n", availableTools)

	agents := llmAgentManager.ListAgents()
	fmt.Printf("🤖 Active Agents: %d\n", len(agents))
	for _, agent := range agents {
		fmt.Printf("  - %s: %s\n", agent.ID, agent.Name)
	}

	tasks := llmAgentManager.ListTasks()
	fmt.Printf("📋 Total Tasks: %d\n", len(tasks))

	fmt.Println("\n🎓 Library Mode Benefits Demonstrated:")
	fmt.Println("====================================")
	fmt.Println("✅ No HTTP server required - pure Go library")
	fmt.Println("✅ Direct function calls for maximum performance")
	fmt.Println("✅ LLM-powered intelligent reasoning and planning")
	fmt.Println("✅ Agents can be embedded in any Go application")
	fmt.Println("✅ Memory and tool state managed internally")
	fmt.Println("✅ Easy integration with existing Go codebases")

	fmt.Println("\n🔧 Integration Examples:")
	fmt.Println("=======================")
	fmt.Println("// Embed in your Go app:")
	fmt.Println("manager := agents.NewLLMAgentManager(server, llm, \"model\")")
	fmt.Println("agent, _ := manager.CreateLLMAgent(\"my_agent\", ...)")
	fmt.Println("result, _ := manager.ExecuteAgentAction(\"my_agent\", \"task\", params)")

	fmt.Println("\n✨ Library mode demo completed successfully!")
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

		if strings.Contains(query, "7") && strings.Contains(query, "8") {
			// Direct calculation response
			return `{
  "analysis": "I need to calculate 7 × 8 and store the result as requested.",
  "steps": [
    {
      "name": "calculate_product",
      "description": "Multiply 7 by 8",
      "tool": "multiply",
      "input": {"a": 7.0, "b": 8.0}
    },
    {
      "name": "store_calculation",
      "description": "Store the calculation result",
      "tool": "remember",
      "input": {"key": "direct_calc", "value": "7 × 8 = 56"}
    }
  ],
  "reasoning": "Direct multiplication requested, I'll use the multiply tool and store the result."
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
