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
	fmt.Println("🧠 AI Agents Library Mode - No Server Required")
	fmt.Println("==============================================")
	fmt.Println("Using agents directly as a library without HTTP server")

	// Create a minimal configuration for library mode
	config := conduit.DefaultConfig()
	config.EnableLogging = false // Disable server logging

	// Create server instance for tool registry (but don't start it)
	server := conduit.NewEnhancedServer(config)

	// Register all available tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Add custom math tools
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a + b
			fmt.Printf("   🧮 Add: %.1f + %.1f = %.1f\n", a, b, result)
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
			fmt.Printf("   🧮 Multiply: %.1f × %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "multiplication",
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number"),
			"b": conduit.NumberParam("Second number"),
		}, []string{"a", "b"}))

	// Create mock LLM for demonstration
	mockLLM := createIntelligentMockLLM()

	// Create LLM-powered agent manager (library mode - no server)
	fmt.Println("📚 Creating agent manager in library mode...")
	llmAgentManager := agents.NewLLMAgentManager(server, mockLLM, "mock-llm-v1")

	// Create intelligent agents
	fmt.Println("🤖 Creating intelligent agents...")

	mathAgent, err := llmAgentManager.CreateLLMAgent(
		"math_genius",
		"Math Genius",
		"Expert at mathematical problem solving with step-by-step reasoning",
		`You are a mathematical genius. You excel at breaking down complex problems into simple steps and using the right tools to solve them efficiently.`,
		[]string{"add", "multiply"},
		&agents.AgentConfig{
			MaxTokens:     500,
			Temperature:   0.1,
			EnableMemory:  true,
			EnableLogging: false,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create math agent: %v", err)
	}

	textAgent, err := llmAgentManager.CreateLLMAgent(
		"text_analyst",
		"Text Analyst",
		"Expert at analyzing and processing text with intelligent insights",
		`You are a text analysis expert. You can understand context, extract key information, and make intelligent decisions about text processing.`,
		[]string{"word_count", "char_count", "uppercase", "lowercase", "remember", "uuid"},
		&agents.AgentConfig{
			MaxTokens:     600,
			Temperature:   0.3,
			EnableMemory:  true,
			EnableLogging: false,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create text agent: %v", err)
	}

	fmt.Printf("✅ Created: %s\n", mathAgent.Name)
	fmt.Printf("✅ Created: %s\n", textAgent.Name)

	// Note: No server.Start() - we're using pure library mode!
	fmt.Println("🚀 Agents ready! (No server required)")

	// Demo 1: Math problem solving
	fmt.Println("\n🧮 Demo 1: Math Problem Solving")
	fmt.Println("Problem: A bakery sells 12 dozen cookies at $3 per dozen. What's the total revenue?")

	mathTask, _ := llmAgentManager.CreateTask(
		"math_genius",
		"Bakery Revenue Calculation",
		"Calculate total revenue from cookie sales",
		map[string]interface{}{
			"problem":         "A bakery sells 12 dozen cookies at $3 per dozen. What's the total revenue?",
			"dozens":          12.0,
			"price_per_dozen": 3.0,
		},
	)

	fmt.Println("🧠 Math Genius is analyzing the problem...")
	if err := llmAgentManager.ExecuteTaskWithLLM(mathTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Problem solved!")
		printLibraryTaskResult(mathTask)
	}

	// Demo 2: Text analysis
	fmt.Println("\n📝 Demo 2: Text Analysis")
	sampleText := "The future of artificial intelligence lies in creating systems that can reason, learn, and adapt intelligently."
	fmt.Printf("Text: %s\n", sampleText)

	textTask, _ := llmAgentManager.CreateTask(
		"text_analyst",
		"AI Text Analysis",
		"Analyze text about AI and extract insights",
		map[string]interface{}{
			"text":          sampleText,
			"analysis_goal": "extract_key_concepts_and_metrics",
		},
	)

	fmt.Println("🧠 Text Analyst is analyzing the content...")
	if err := llmAgentManager.ExecuteTaskWithLLM(textTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Analysis complete!")
		printLibraryTaskResult(textTask)
	}

	// Demo 3: Complex multi-step problem
	fmt.Println("\n🔢 Demo 3: Complex Multi-Step Problem")
	fmt.Println("Problem: Calculate compound interest: $1000 principal, 5% rate, compounded twice")

	compoundTask, _ := llmAgentManager.CreateTask(
		"math_genius",
		"Compound Interest Calculation",
		"Calculate compound interest step by step",
		map[string]interface{}{
			"principal":    1000.0,
			"rate":         0.05,
			"compounds":    2.0,
			"problem_type": "compound_interest",
		},
	)

	fmt.Println("🧠 Math Genius is working on compound interest...")
	if err := llmAgentManager.ExecuteTaskWithLLM(compoundTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Compound interest calculated!")
		printLibraryTaskResult(compoundTask)
	}

	// Show library mode benefits
	fmt.Println("\n🎯 Library Mode Benefits:")
	fmt.Println("========================")
	fmt.Println("✅ No HTTP server required")
	fmt.Println("✅ Direct function calls")
	fmt.Println("✅ Faster execution")
	fmt.Println("✅ Lower memory footprint")
	fmt.Println("✅ Easy integration into existing applications")
	fmt.Println("✅ Full access to all MCP tools")

	// Show total usage
	fmt.Printf("\n📊 Session Summary:\n")
	fmt.Printf("Total Agents: %d\n", len(llmAgentManager.ListAgents()))
	fmt.Printf("Total Tasks: %d\n", len(llmAgentManager.ListTasks()))
	fmt.Printf("Available Tools: %d\n", len(llmAgentManager.GetAvailableTools()))

	fmt.Println("\n🔧 Real LLM Integration Options:")
	fmt.Println("===============================")
	fmt.Println("Replace the mock LLM with real AI:")
	fmt.Println("")
	fmt.Println("🌊 DeepInfra (Recommended):")
	fmt.Println("  export DEEPINFRA_TOKEN=your_token")
	fmt.Println("  model := conduit.CreateDeepInfraModel(token)")
	fmt.Println("  manager := agents.NewLLMAgentManager(server, model, \"meta-llama/Meta-Llama-3.1-8B-Instruct\")")
	fmt.Println("")
	fmt.Println("🦙 Local Ollama:")
	fmt.Println("  export OLLAMA_URL=http://localhost:11434")
	fmt.Println("  model := conduit.CreateOllamaModel(url)")
	fmt.Println("  manager := agents.NewLLMAgentManager(server, model, \"llama3.2\")")
	fmt.Println("")
	fmt.Println("📁 Example Files:")
	fmt.Println("  examples/agents_deepinfra/    - DeepInfra integration")
	fmt.Println("  examples/agents_ollama/       - Local Ollama integration")

	fmt.Println("\n🎉 Library mode demonstration completed!")
	fmt.Println("All processing done locally without any server!")
}

// createIntelligentMockLLM creates a more sophisticated mock LLM
func createIntelligentMockLLM() mcp.ModelFunc {
	return func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])

		// Bakery revenue calculation
		if strings.Contains(query, "bakery") && strings.Contains(query, "dozen") {
			return `{
  "analysis": "This is a revenue calculation problem. I need to multiply 12 dozen by $3 per dozen to get total revenue.",
  "steps": [
    {
      "name": "calculate_revenue",
      "description": "Multiply dozens sold by price per dozen",
      "tool": "multiply",
      "input": {"a": 12.0, "b": 3.0}
    },
    {
      "name": "store_calculation",
      "description": "Remember this business calculation",
      "tool": "remember",
      "input": {"key": "bakery_revenue", "value": "12 dozen × $3 = $36 total revenue"}
    }
  ],
  "reasoning": "I identified this as a simple multiplication problem: 12 dozen × $3/dozen = total revenue. I'll use the multiply tool and store the result."
}`, nil
		}

		// AI text analysis
		if strings.Contains(query, "artificial intelligence") && strings.Contains(query, "future") {
			return `{
  "analysis": "This text discusses AI future concepts. I should analyze word count, extract key themes, and store insights about AI development.",
  "steps": [
    {
      "name": "count_words",
      "description": "Get word count and text metrics",
      "tool": "word_count",
      "input": {"text": "The future of artificial intelligence lies in creating systems that can reason, learn, and adapt intelligently."}
    },
    {
      "name": "store_ai_insight",
      "description": "Store key AI insight for future reference",
      "tool": "remember",
      "input": {"key": "ai_future_concept", "value": "AI future focuses on reasoning, learning, and adaptation"}
    },
    {
      "name": "generate_analysis_id",
      "description": "Create unique ID for this analysis session",
      "tool": "uuid",
      "input": {}
    }
  ],
  "reasoning": "I analyzed the text about AI's future, focusing on the key concepts of reasoning, learning, and adaptation. I'll count words for metrics and store the insight."
}`, nil
		}

		// Compound interest calculation
		if strings.Contains(query, "compound interest") || strings.Contains(query, "principal") {
			return `{
  "analysis": "This is a compound interest problem. I need to calculate it step by step: first multiply principal by (1 + rate), then multiply again for the second compounding.",
  "steps": [
    {
      "name": "first_compounding",
      "description": "Calculate first compounding: principal × (1 + rate) = 1000 × 1.05",
      "tool": "multiply",
      "input": {"a": 1000.0, "b": 1.05}
    },
    {
      "name": "second_compounding", 
      "description": "Calculate second compounding: result × (1 + rate) again",
      "tool": "multiply",
      "input": {"a": 1050.0, "b": 1.05}
    },
    {
      "name": "store_compound_result",
      "description": "Store the compound interest calculation",
      "tool": "remember",
      "input": {"key": "compound_interest", "value": "1000 compounded twice at 5% = 1102.50"}
    }
  ],
  "reasoning": "Compound interest requires multiple multiplication steps. First: 1000 × 1.05 = 1050, then 1050 × 1.05 = 1102.50. I'll break this into sequential multiplications."
}`, nil
		}

		// Default response
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

// printLibraryTaskResult prints task results in library mode
func printLibraryTaskResult(task *agents.Task) {
	fmt.Printf("  📋 Task: %s\n", task.Title)
	fmt.Printf("  📊 Status: %s (%.0f%% complete)\n", task.Status, task.Progress*100)

	for i, step := range task.Steps {
		fmt.Printf("    Step %d: %s (%s)\n", i+1, step.Name, step.Status)

		if step.Name == "llm_reasoning" && step.Output != nil {
			if analysis, ok := step.Output["llm_analysis"].(string); ok {
				// Show first part of LLM analysis
				lines := strings.Split(analysis, "\n")
				if len(lines) > 0 {
					fmt.Printf("      🧠 LLM: Analyzed and planned execution\n")
				}
			}
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
