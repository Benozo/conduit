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
	fmt.Println("🧠 LLM-Powered Agent Swarm with Ollama")
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

	// Register custom swarm-specific tools
	registerSwarmTools(server)

	// Create Ollama model function
	ollamaModel := conduit.CreateOllamaModel(ollamaURL)

	// Create swarm client with LLM integration
	swarmClient := swarm.NewSwarmClientWithLLM(server, swarm.DefaultSwarmConfig(), ollamaModel, modelName)

	// Create agents with specialized roles
	coordinator := createCoordinator(swarmClient)
	contentCreator := createContentCreator(swarmClient)
	dataAnalyst := createDataAnalyst(swarmClient)
	memoryManager := createMemoryManager(swarmClient)

	// Set up agent handoff functions
	setupAgentHandoffs(swarmClient, coordinator, contentCreator, dataAnalyst, memoryManager)

	fmt.Println("🎯 Agent Swarm Created with LLM Intelligence:")
	fmt.Printf("   📋 %s - Routes tasks to appropriate specialists\n", coordinator.Name)
	fmt.Printf("   ✍️  %s - Handles content creation and text processing\n", contentCreator.Name)
	fmt.Printf("   📊 %s - Performs data analysis and reporting\n", dataAnalyst.Name)
	fmt.Printf("   🧠 %s - Manages information storage and retrieval\n", memoryManager.Name)
	fmt.Println()

	// Demo scenarios with LLM reasoning
	runLLMDemoScenarios(swarmClient, coordinator)
}

func registerSwarmTools(server *conduit.EnhancedServer) {
	// Content creation tools
	server.GetToolRegistry().Register("write_article", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		title := "Generated Article"
		topic := "general topic"

		if t, ok := params["title"].(string); ok {
			title = t
		}
		if tp, ok := params["topic"].(string); ok {
			topic = tp
		}

		return fmt.Sprintf("📝 Article '%s' about %s has been written and saved.", title, topic), nil
	})

	server.GetToolRegistry().Register("research_topic", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		topic := "general research"
		if t, ok := params["topic"].(string); ok {
			topic = t
		}

		return fmt.Sprintf("🔍 Research completed for topic: %s. Found 5 relevant sources and key insights.", topic), nil
	})

	// Data analysis tools
	server.GetToolRegistry().Register("analyze_data", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		dataset := "data.csv"
		if d, ok := params["dataset"].(string); ok {
			dataset = d
		}

		return fmt.Sprintf("📈 Analysis complete for %s: Found 3 key insights, 2 anomalies, and 5 recommendations.", dataset), nil
	})

	server.GetToolRegistry().Register("generate_report", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		findings := "Analysis complete"
		if f, ok := params["findings"].(string); ok {
			findings = f
		}

		return fmt.Sprintf("📊 Report generated based on: %s", findings), nil
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

func createCoordinator(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"Coordinator",
		`You are the Coordinator agent responsible for routing tasks to appropriate specialists.

Your role:
- Analyze incoming requests and determine the best agent to handle them
- Route content creation tasks to ContentCreator
- Route data analysis tasks to DataAnalyst  
- Route memory/storage tasks to MemoryManager
- Provide guidance when users are unsure what they need

CRITICAL HANDOFF RULES:
- ALWAYS transfer to ONE AGENT AT A TIME - never multiple agents
- Use exact function names: "transfer_to_content_creator", "transfer_to_data_analyst", "transfer_to_memory_manager"
- NEVER try to handoff to multiple agents or use empty names
- Example of WRONG usage: transfer to multiple agents or empty string
- Example of CORRECT usage: Call "transfer_to_content_creator" for content tasks

When deciding on handoffs, consider:
- Content creation: articles, writing, text processing, research → use "transfer_to_content_creator"
- Data analysis: datasets, reports, statistics, analysis → use "transfer_to_data_analyst"
- Memory management: storing information, remembering context → use "transfer_to_memory_manager"

Always explain your reasoning when transferring to another agent.`,
		[]string{})
}

func createContentCreator(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"ContentCreator",
		`You are the ContentCreator agent specialized in content creation and text processing.

Your capabilities:
- Writing articles and content
- Researching topics
- Text processing and manipulation
- Content optimization
- Information synthesis

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "write_article", "research_topic", "store_context", etc.
- Wait for each tool response before calling another tool

Use your tools to complete content creation tasks efficiently. When you need to store important information for later use, remember to use memory tools.`,
		[]string{"write_article", "research_topic", "uppercase", "lowercase", "word_count", "store_context"})
}

func createDataAnalyst(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"DataAnalyst",
		`You are the DataAnalyst agent specialized in data analysis and reporting.

Your capabilities:
- Analyzing datasets and finding insights
- Generating comprehensive reports
- Statistical analysis
- Data visualization recommendations
- Trend identification

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "analyze_data", "generate_report", "store_context", etc.
- Wait for each tool response before calling another tool

Use your analytical tools to process data and generate actionable insights. Store important findings in memory for future reference.`,
		[]string{"analyze_data", "generate_report", "uuid_generate", "timestamp", "encode_base64", "store_context"})
}

func createMemoryManager(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"MemoryManager",
		`You are the MemoryManager agent specialized in information storage and retrieval.

Your capabilities:
- Storing and organizing information
- Retrieving relevant context
- Managing shared knowledge across agents
- Information categorization
- Context synthesis

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "store_context", "retrieve_context"
- Wait for each tool response before calling another tool

You help maintain continuity across conversations and tasks by managing shared memory effectively.`,
		[]string{"store_context", "retrieve_context"})
}

func setupAgentHandoffs(client swarm.SwarmClient, coordinator, contentCreator, dataAnalyst, memoryManager *swarm.Agent) {
	// Coordinator can transfer to any specialist
	client.RegisterFunction(coordinator.Name, swarm.CreateHandoffFunction("content_creator", contentCreator))
	client.RegisterFunction(coordinator.Name, swarm.CreateHandoffFunction("data_analyst", dataAnalyst))
	client.RegisterFunction(coordinator.Name, swarm.CreateHandoffFunction("memory_manager", memoryManager))

	// Specialists can transfer back to coordinator or to each other when needed
	client.RegisterFunction(contentCreator.Name, swarm.CreateHandoffFunction("coordinator", coordinator))
	client.RegisterFunction(contentCreator.Name, swarm.CreateHandoffFunction("data_analyst", dataAnalyst))
	client.RegisterFunction(contentCreator.Name, swarm.CreateHandoffFunction("memory_manager", memoryManager))

	client.RegisterFunction(dataAnalyst.Name, swarm.CreateHandoffFunction("coordinator", coordinator))
	client.RegisterFunction(dataAnalyst.Name, swarm.CreateHandoffFunction("content_creator", contentCreator))
	client.RegisterFunction(dataAnalyst.Name, swarm.CreateHandoffFunction("memory_manager", memoryManager))

	client.RegisterFunction(memoryManager.Name, swarm.CreateHandoffFunction("coordinator", coordinator))
	client.RegisterFunction(memoryManager.Name, swarm.CreateHandoffFunction("content_creator", contentCreator))
	client.RegisterFunction(memoryManager.Name, swarm.CreateHandoffFunction("data_analyst", dataAnalyst))
}

func runLLMDemoScenarios(client swarm.SwarmClient, coordinator *swarm.Agent) {
	fmt.Println("🚀 Running LLM-Powered Demo Scenarios:")
	fmt.Println("=====================================")

	scenarios := []struct {
		name        string
		message     string
		description string
	}{
		{
			"Content Creation Request",
			"I need to write an article about artificial intelligence in healthcare. Can you help me research and create this content?",
			"Tests LLM routing to ContentCreator and tool usage",
		},
		{
			"Data Analysis Request",
			"I have a customer behavior dataset that needs analysis. Please analyze the data and generate a comprehensive report.",
			"Tests LLM routing to DataAnalyst and tool usage",
		},
		{
			"Memory Management Request",
			"Please remember that our Q4 project deadline is December 15th and we're currently 75% complete.",
			"Tests LLM routing to MemoryManager and context storage",
		},
		{
			"Complex Multi-Agent Request",
			"I need to research AI trends, analyze market data, and remember the key insights for our strategic planning meeting.",
			"Tests LLM coordination across multiple agents",
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

		// Run swarm with LLM reasoning
		response := client.Run(coordinator, messages, map[string]interface{}{
			"scenario": scenario.name,
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

	fmt.Println("\n🎉 LLM-Powered Agent Swarm Demo Complete!")
	fmt.Println("\n✨ Key Features Demonstrated:")
	fmt.Println("   🧠 LLM-powered agent reasoning and decision-making")
	fmt.Println("   🎯 Intelligent task routing based on content understanding")
	fmt.Println("   🔧 Smart tool selection using natural language analysis")
	fmt.Println("   🔄 Context-aware agent handoffs")
	fmt.Println("   💭 Reasoning explanations for transparency")
	fmt.Println("   📚 Conversation memory and continuity")

	fmt.Println("\n🔧 How LLM Integration Works:")
	fmt.Println("   1. Each agent receives a specialized system prompt")
	fmt.Println("   2. LLM analyzes user messages and conversation context")
	fmt.Println("   3. LLM decides between tool usage, agent handoffs, or direct responses")
	fmt.Println("   4. Decisions include reasoning for transparency")
	fmt.Println("   5. Fallback to rule-based logic if LLM fails")

	fmt.Println("\n🚀 Next Steps:")
	fmt.Println("   • Experiment with different Ollama models")
	fmt.Println("   • Add custom tools for your domain")
	fmt.Println("   • Create specialized agents for specific workflows")
	fmt.Println("   • Scale to larger agent swarms with complex coordination")
}
