package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	fmt.Println("🧠 AI Agents with Ollama LLM Integration")
	fmt.Println("========================================")

	// Get Ollama configuration from environment
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434" // Default Ollama URL
	}

	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3.2" // Default model
	}

	fmt.Printf("🔗 Ollama URL: %s\n", ollamaURL)
	fmt.Printf("🤖 Model: %s\n", modelName)

	// Create MCP server configuration
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8088
	config.OllamaURL = ollamaURL
	config.EnableLogging = true

	// Create the MCP server with Ollama integration
	server := conduit.NewEnhancedServer(config)

	// Register all available tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Register math tools for agents
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a + b
			fmt.Printf("   🧮 Executing: %.1f + %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "addition",
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("add", "Add two numbers together", map[string]interface{}{
			"a": conduit.NumberParam("First number to add"),
			"b": conduit.NumberParam("Second number to add"),
		}, []string{"a", "b"}))

	server.RegisterToolWithSchema("multiply",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a * b
			fmt.Printf("   🧮 Executing: %.1f × %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "multiplication",
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number to multiply"),
			"b": conduit.NumberParam("Second number to multiply"),
		}, []string{"a", "b"}))

	// Create Ollama model function
	ollamaModel := conduit.CreateOllamaModel(ollamaURL)

	// Create LLM-powered agent manager
	fmt.Println("🚀 Creating LLM-powered agent manager...")
	llmAgentManager := agents.NewLLMAgentManager(server, ollamaModel, modelName)

	// Create intelligent agents with LLM reasoning
	fmt.Println("🧠 Creating intelligent agents...")

	// Math reasoning agent
	mathAgent, err := llmAgentManager.CreateLLMAgent(
		"llm_math_agent",
		"LLM Math Reasoner",
		"An intelligent agent that uses LLM reasoning for mathematical problem solving",
		`You are a mathematical reasoning assistant. You can analyze mathematical problems, break them down into steps, and use available tools (add, multiply) to solve them. Always explain your reasoning and show your work step by step.`,
		[]string{"add", "multiply"},
		&agents.AgentConfig{
			MaxTokens:     1000,
			Temperature:   0.1, // Low temperature for precise math
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create math agent: %v", err)
	}

	// Text analysis agent
	textAgent, err := llmAgentManager.CreateLLMAgent(
		"llm_text_agent",
		"LLM Text Analyst",
		"An intelligent agent that uses LLM reasoning for text analysis and processing",
		`You are a text analysis expert. You can analyze text content, extract insights, and use text processing tools effectively. You understand context and can make intelligent decisions about which text operations to perform.`,
		[]string{"word_count", "char_count", "uppercase", "lowercase", "title_case", "reverse"},
		&agents.AgentConfig{
			MaxTokens:     1200,
			Temperature:   0.3, // Moderate temperature for creative text analysis
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create text agent: %v", err)
	}

	// Data management agent
	dataAgent, err := llmAgentManager.CreateLLMAgent(
		"llm_data_agent",
		"LLM Data Manager",
		"An intelligent agent that uses LLM reasoning for data management and organization",
		`You are a data management specialist. You can organize, store, and retrieve information intelligently using memory and utility tools. You make smart decisions about data structure and storage strategies.`,
		[]string{"remember", "recall", "forget", "uuid", "timestamp", "base64_encode"},
		&agents.AgentConfig{
			MaxTokens:     800,
			Temperature:   0.2, // Low temperature for consistent data handling
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create data agent: %v", err)
	}

	fmt.Printf("✅ Created %s\n", mathAgent.Name)
	fmt.Printf("✅ Created %s\n", textAgent.Name)
	fmt.Printf("✅ Created %s\n", dataAgent.Name)

	// Start the MCP server
	go func() {
		fmt.Println("🔄 Starting MCP server with Ollama integration...")
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(3 * time.Second)

	// Test LLM-powered agents
	fmt.Println("\n🧪 Testing LLM-Powered Agents")
	fmt.Println("==============================")

	// Test 1: Complex Math Problem with LLM Reasoning
	fmt.Println("\n🧮 Test 1: Complex Math Problem")
	fmt.Println("Problem: Calculate the area of a rectangle with length 15 and width 8")

	mathTask, err := llmAgentManager.CreateTask(
		"llm_math_agent",
		"Rectangle Area Calculation",
		"Calculate the area of a rectangle using length and width",
		map[string]interface{}{
			"problem":   "Calculate the area of a rectangle with length 15 and width 8",
			"length":    15.0,
			"width":     8.0,
			"operation": "area_calculation",
		},
	)
	if err != nil {
		log.Printf("Failed to create math task: %v", err)
	} else {
		fmt.Println("🧠 LLM is reasoning about the problem...")
		if err := llmAgentManager.ExecuteTaskWithLLM(mathTask.ID); err != nil {
			fmt.Printf("❌ LLM task failed: %v\n", err)
		} else {
			fmt.Printf("✅ LLM math task completed!\n")
			printLLMTaskResult(mathTask)
		}
	}

	// Test 2: Intelligent Text Analysis
	fmt.Println("\n📝 Test 2: Intelligent Text Analysis")
	text := "The quick brown fox jumps over the lazy dog. This sentence contains every letter of the alphabet!"
	fmt.Printf("Text: %s\n", text)

	textTask, err := llmAgentManager.CreateTask(
		"llm_text_agent",
		"Comprehensive Text Analysis",
		"Analyze text and provide comprehensive insights",
		map[string]interface{}{
			"text":          text,
			"analysis_type": "comprehensive",
			"include_stats": true,
		},
	)
	if err != nil {
		log.Printf("Failed to create text task: %v", err)
	} else {
		fmt.Println("🧠 LLM is analyzing the text...")
		if err := llmAgentManager.ExecuteTaskWithLLM(textTask.ID); err != nil {
			fmt.Printf("❌ LLM text task failed: %v\n", err)
		} else {
			fmt.Printf("✅ LLM text analysis completed!\n")
			printLLMTaskResult(textTask)
		}
	}

	// Test 3: Smart Data Management
	fmt.Println("\n🗄️  Test 3: Smart Data Management")
	fmt.Println("Task: Organize user session data intelligently")

	dataTask, err := llmAgentManager.CreateTask(
		"llm_data_agent",
		"User Session Management",
		"Create and organize user session data with proper structure",
		map[string]interface{}{
			"user_id":                "user_12345",
			"session_type":           "login",
			"create_structured_data": true,
		},
	)
	if err != nil {
		log.Printf("Failed to create data task: %v", err)
	} else {
		fmt.Println("🧠 LLM is organizing the data...")
		if err := llmAgentManager.ExecuteTaskWithLLM(dataTask.ID); err != nil {
			fmt.Printf("❌ LLM data task failed: %v\n", err)
		} else {
			fmt.Printf("✅ LLM data management completed!\n")
			printLLMTaskResult(dataTask)
		}
	}

	// Comparison with non-LLM agents
	fmt.Println("\n🔄 Comparison: Non-LLM vs LLM Agents")
	fmt.Println("====================================")

	// Create traditional agent for comparison
	standardAgentManager := agents.NewMCPAgentManager(server)
	_, _ = standardAgentManager.CreateAgent(
		"standard_math",
		"Standard Math Agent",
		"Traditional rule-based math agent",
		"Basic math agent",
		[]string{"add", "multiply"},
		agents.DefaultAgentConfig(),
	)

	// Test with same problem
	standardTask, _ := standardAgentManager.CreateTask(
		"standard_math",
		"Standard Area Calculation",
		"Basic area calculation",
		map[string]interface{}{
			"a":         15.0,
			"b":         8.0,
			"operation": "multiply",
		},
	)

	fmt.Println("🤖 Standard agent executing...")
	standardAgentManager.ExecuteTask(standardTask.ID)

	fmt.Println("\n📊 Results Comparison:")
	fmt.Printf("Standard Agent: %d steps, rule-based planning\n", len(standardTask.Steps))
	fmt.Printf("LLM Agent: %d steps, intelligent reasoning\n", len(mathTask.Steps))

	fmt.Println("\n🎉 LLM-Powered AI Agents Demo Completed!")
	fmt.Println("Key Benefits of LLM Integration:")
	fmt.Println("✅ Intelligent problem analysis and reasoning")
	fmt.Println("✅ Context-aware decision making")
	fmt.Println("✅ Natural language understanding")
	fmt.Println("✅ Adaptive error recovery")
	fmt.Println("✅ Complex multi-step planning")

	fmt.Printf("\n🔗 Server running on http://localhost:%d\n", config.Port)

	// Keep running if interactive mode
	if len(os.Args) > 1 && os.Args[1] == "--interactive" {
		fmt.Println("Press Ctrl+C to exit...")
		select {}
	}
}

// printLLMTaskResult prints detailed results of LLM-powered task execution
func printLLMTaskResult(task *agents.Task) {
	fmt.Printf("  📋 Task: %s\n", task.Title)
	fmt.Printf("  📊 Status: %s (%.1f%% complete)\n", task.Status, task.Progress*100)
	fmt.Printf("  🔢 Steps: %d executed\n", len(task.Steps))

	for i, step := range task.Steps {
		fmt.Printf("    Step %d: %s (%s)\n", i+1, step.Name, step.Status)

		if step.Name == "llm_reasoning" && step.Output != nil {
			if analysis, ok := step.Output["llm_analysis"].(string); ok {
				// Extract key insights from LLM analysis
				if len(analysis) > 200 {
					fmt.Printf("      🧠 LLM Analysis: %s...\n", analysis[:200])
				} else {
					fmt.Printf("      🧠 LLM Analysis: %s\n", analysis)
				}
			}
		}

		if step.Output != nil {
			if result, ok := step.Output["result"]; ok {
				fmt.Printf("      ✅ Result: %v\n", result)
			}
			if output, ok := step.Output["output"]; ok {
				if outputMap, ok := output.(map[string]interface{}); ok {
					if result, ok := outputMap["result"]; ok {
						fmt.Printf("      ✅ Tool Result: %v\n", result)
					}
				}
			}
		}
	}
}
