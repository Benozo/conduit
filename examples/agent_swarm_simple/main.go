package main

import (
	"fmt"
	"os"
	"strings"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
	"github.com/benozo/conduit/swarm"
)

func main() {
	fmt.Println("🧠 LLM-Powered Simple Agent Swarm Demo")
	fmt.Println("=====================================")
	fmt.Println()

	// Get Ollama configuration
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://192.168.10.10:11434"
	}

	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3.2"
	}

	fmt.Printf("🦙 Ollama URL: %s\n", ollamaURL)
	fmt.Printf("🤖 Model: %s\n", modelName)
	fmt.Println()

	// Create MCP server
	config := conduit.DefaultConfig()
	config.EnableLogging = false
	server := conduit.NewEnhancedServer(config)

	// Register comprehensive tool set (skip memory tools to avoid conflicts)
	tools.RegisterTextTools(server)
	// tools.RegisterMemoryTools(server) // Skip to avoid conflicts with custom store_context
	tools.RegisterUtilityTools(server)

	// Register custom simple tools
	registerSimpleTools(server)

	// Create Ollama model function
	ollamaModel := conduit.CreateOllamaModel(ollamaURL)

	// Create swarm client with LLM integration (no fallback)
	swarmClient := swarm.NewSwarmClientWithLLM(server, swarm.DefaultSwarmConfig(), ollamaModel, modelName)

	// Create simple agent setup
	router := createRouter(swarmClient)
	textProcessor := createTextProcessor(swarmClient)
	analyst := createAnalyst(swarmClient)

	// Set up agent handoffs
	setupSimpleHandoffs(swarmClient, router, textProcessor, analyst)

	fmt.Println("🎯 LLM-Powered Simple Agent Swarm Created:")
	fmt.Printf("   🚦 %s - Routes tasks to appropriate specialists\n", router.Name)
	fmt.Printf("   📝 %s - Handles text processing and formatting\n", textProcessor.Name)
	fmt.Printf("   📊 %s - Performs analysis and insights\n", analyst.Name)
	fmt.Println()

	// Run simple demo scenarios
	runSimpleDemoScenarios(swarmClient, router)
}

func registerSimpleTools(server *conduit.EnhancedServer) {
	// Enhanced text processing tool
	server.GetToolRegistry().Register("process_text", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := "sample text"
		operation := "uppercase"

		if t, ok := params["text"].(string); ok {
			text = t
		}
		if op, ok := params["operation"].(string); ok {
			operation = op
		}

		switch operation {
		case "uppercase":
			return fmt.Sprintf("📝 Processed text (uppercase): %s", strings.ToUpper(text)), nil
		case "lowercase":
			return fmt.Sprintf("📝 Processed text (lowercase): %s", strings.ToLower(text)), nil
		case "count_words":
			words := len(strings.Fields(text))
			return fmt.Sprintf("📊 Word count for text: %d words", words), nil
		case "reverse":
			runes := []rune(text)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return fmt.Sprintf("📝 Processed text (reversed): %s", string(runes)), nil
		default:
			return fmt.Sprintf("❓ Unknown operation: %s. Available: uppercase, lowercase, count_words, reverse", operation), nil
		}
	})

	// Enhanced analysis tool
	server.GetToolRegistry().Register("analyze", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		data := "sample data"
		if d, ok := params["data"].(string); ok {
			data = d
		}

		return fmt.Sprintf("📈 Analysis complete for: %s - Found patterns, trends, and key insights", data), nil
	})

	// Memory management tools (custom implementation)
	server.GetToolRegistry().Register("store_context", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := "context"
		value := "stored information"

		if k, ok := params["key"].(string); ok {
			key = k
		}
		if v, ok := params["value"].(string); ok {
			value = v
		}

		memory.Set(key, value)
		return fmt.Sprintf("🧠 Stored context: %s = %s", key, value), nil
	})

	server.GetToolRegistry().Register("retrieve_context", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := "context"
		if k, ok := params["key"].(string); ok {
			key = k
		}

		value := memory.Get(key)
		if value == nil {
			return fmt.Sprintf("🔍 No context found for key: %s", key), nil
		}
		return fmt.Sprintf("🧠 Retrieved context: %s = %s", key, value), nil
	})
}

func createRouter(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"Router",
		`You are the Router agent responsible for analyzing requests and routing them to appropriate specialists.

Your role:
- Analyze incoming user requests to understand their intent
- Route text processing tasks to TextProcessor agent
- Route analysis tasks to Analyst agent
- You are NOT responsible for doing the actual work - always delegate to specialists

CRITICAL TOOL USAGE RULES:
- ONLY use "store_context" for routing decisions, not for the main task
- Do NOT perform text processing yourself - always route to TextProcessor
- Do NOT perform analysis yourself - always route to Analyst

CRITICAL HANDOFF RULES:
- ALWAYS transfer to ONE AGENT AT A TIME - never multiple agents
- Use exact function names: "transfer_to_text_processor", "transfer_to_analyst"
- NEVER try to call agent names as tools (e.g., don't call "TextProcessor" as a tool)
- Example of WRONG usage: calling "TextProcessor" as a tool or doing work yourself
- Example of CORRECT usage: Call "transfer_to_text_processor" function for text tasks

ROUTING DECISIONS - YOU MUST ROUTE:
- Text processing tasks (uppercase, lowercase, word count, reverse) → ALWAYS use "transfer_to_text_processor"
- Analysis tasks (data analysis, insights, patterns) → ALWAYS use "transfer_to_analyst"
- ANY text manipulation request → ALWAYS use "transfer_to_text_processor"
- Mixed requests: If request has BOTH text processing AND analysis, route to TextProcessor FIRST

IMPORTANT: For requests like "count words and analyze", route to TextProcessor first for the word counting, then that agent can coordinate with Analyst if needed.

Always explain your routing decision and then immediately transfer to the appropriate specialist.`,
		[]string{"store_context"})
}

func createTextProcessor(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"TextProcessor",
		`You are the TextProcessor agent specialized in text processing and manipulation.

Your capabilities:
- Text formatting (uppercase, lowercase, reverse)
- Word counting and text analysis
- Text transformation and processing
- Memory storage for processed results

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "process_text", "store_context", "uppercase", "lowercase", "word_count"
- Wait for each tool response before calling another tool

Your workflow:
1. First call the appropriate processing tool (process_text, uppercase, etc.)
2. Wait for response, then call "store_context" if needed to save results
3. Provide clear feedback on what was processed

Use your tools efficiently to handle all text processing tasks.`,
		[]string{"process_text", "store_context", "uppercase", "lowercase", "word_count"})
}

func createAnalyst(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"Analyst",
		`You are the Analyst agent specialized in data analysis and insights generation.

Your capabilities:
- Data analysis and pattern recognition
- Insight generation and trend identification
- Information synthesis and reporting
- Memory storage for analytical findings

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "analyze", "store_context", "uuid_generate", "timestamp"
- Wait for each tool response before calling another tool

Your workflow:
1. First call "analyze" to examine the data or information
2. Wait for response, then call other tools like "uuid_generate" or "timestamp" if needed
3. Finally call "store_context" to save important findings

Focus on providing valuable insights and clear recommendations.`,
		[]string{"analyze", "store_context", "uuid_generate", "timestamp"})
}

func setupSimpleHandoffs(client swarm.SwarmClient, router, textProcessor, analyst *swarm.Agent) {
	fmt.Println("🔗 Setting up simple agent handoff functions...")

	// Router can transfer to specialists
	client.RegisterFunction(router.Name, swarm.CreateHandoffFunction("text_processor", textProcessor))
	client.RegisterFunction(router.Name, swarm.CreateHandoffFunction("analyst", analyst))

	// Specialists can transfer back to router or to each other
	client.RegisterFunction(textProcessor.Name, swarm.CreateHandoffFunction("router", router))
	client.RegisterFunction(textProcessor.Name, swarm.CreateHandoffFunction("analyst", analyst))

	client.RegisterFunction(analyst.Name, swarm.CreateHandoffFunction("router", router))
	client.RegisterFunction(analyst.Name, swarm.CreateHandoffFunction("text_processor", textProcessor))

	fmt.Println("✅ Simple agent handoff functions registered")
}

func runSimpleDemoScenarios(client swarm.SwarmClient, router *swarm.Agent) {
	fmt.Println("🚀 Running LLM-Powered Simple Demo Scenarios:")
	fmt.Println("==============================================")

	scenarios := []struct {
		name        string
		message     string
		description string
	}{
		{
			"Text Processing Task",
			"Convert 'Hello World' to uppercase and remember it for later use",
			"Tests LLM routing to TextProcessor and tool usage",
		},
		{
			"Text Analysis Task",
			"Please count the words in this sentence: 'Count the words in this sentence and analyze the text structure'",
			"Tests LLM text processing with word counting capabilities",
		},
		{
			"Data Analysis Task",
			"Analyze customer feedback data and remember the key insights for our team",
			"Tests LLM routing to Analyst and insight generation",
		},
		{
			"Multi-Step Processing",
			"Please reverse this text: 'AI is transforming business'",
			"Tests LLM routing to TextProcessor for text reversal",
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("\n📝 Scenario %d: %s\n", i+1, scenario.name)
		fmt.Printf("📄 Description: %s\n", scenario.description)
		fmt.Printf("💬 Request: %s\n", scenario.message)
		fmt.Println("🔄 LLM Processing...")

		// Create message
		messages := []swarm.Message{
			{Role: "user", Content: scenario.message},
		}

		// Run swarm with LLM reasoning (no fallback)
		response := client.Run(router, messages, map[string]interface{}{
			"scenario": scenario.name,
			"demo":     "simple_swarm",
		})

		fmt.Println("📊 Response:")
		if response.Success {
			fmt.Printf("✅ Success! Turns: %d, Tool calls: %d, Handoffs: %d\n",
				response.TotalTurns, response.ToolCallsCount, response.HandoffsCount)

			// Show conversation flow
			for j, msg := range response.Messages {
				if j > 0 { // Skip initial user message
					role := "🤖"
					if msg.Role == "user" {
						role = "👤"
					}
					fmt.Printf("   %s %s\n", role, msg.Content)
				}
			}

			fmt.Printf("🎯 Final Agent: %s\n", response.Agent.Name)
		} else {
			fmt.Printf("❌ Error: %v\n", response.Error)
		}

		fmt.Printf("⏱️  Execution Time: %v\n", response.ExecutionTime)
		fmt.Println(strings.Repeat("-", 80))
	}

	fmt.Println("\n🎉 LLM-Powered Simple Agent Swarm Demo Complete!")
	fmt.Println("\n✨ Key Features Demonstrated:")
	fmt.Println("   🧠 LLM-powered task routing and decision-making")
	fmt.Println("   🎯 Intelligent agent selection based on task type")
	fmt.Println("   🔧 Smart tool usage with natural language understanding")
	fmt.Println("   🔄 Context-aware agent handoffs")
	fmt.Println("   💭 Clear reasoning for all routing decisions")

	fmt.Println("\n🔧 Simple LLM Integration:")
	fmt.Println("   • Router analyzes requests and routes intelligently")
	fmt.Println("   • TextProcessor handles text manipulation tasks")
	fmt.Println("   • Analyst performs analysis and generates insights")
	fmt.Println("   • No rule-based fallback - pure LLM reasoning")

	fmt.Println("\n🚀 Next Steps:")
	fmt.Println("   • Try different text processing operations")
	fmt.Println("   • Experiment with various analysis requests")
	fmt.Println("   • Test complex multi-step workflows")
	fmt.Println("   • Add more specialized agents for your needs")
}
