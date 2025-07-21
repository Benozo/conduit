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
	fmt.Println("🧠 LLM-Powered Agent Swarm - Multi-Agent Content & Analysis Workflow")
	fmt.Println("====================================================================")
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

	// Register custom swarm-specific tools first
	registerSwarmTools(server)

	// Register comprehensive tool set (skip memory tools since we have custom ones)
	tools.RegisterTextTools(server)
	// tools.RegisterMemoryTools(server) // Skip to avoid conflicts with custom store_context/retrieve_context
	tools.RegisterUtilityTools(server)

	// Create Ollama model function
	ollamaModel := conduit.CreateOllamaModel(ollamaURL)

	// Create swarm client with LLM integration (no fallback)
	swarmClient := swarm.NewSwarmClientWithLLM(server, swarm.DefaultSwarmConfig(), ollamaModel, modelName)

	// Debug: Show that tools are registered
	fmt.Println("🔧 Tool registration complete")

	// Create agents with specialized roles
	coordinator := createCoordinator(swarmClient)
	contentCreator := createContentCreator(swarmClient)
	dataAnalyst := createDataAnalyst(swarmClient)
	memoryManager := createMemoryManager(swarmClient)

	// Set up agent handoff functions
	setupAgentHandoffs(swarmClient, coordinator, contentCreator, dataAnalyst, memoryManager)

	fmt.Println("🎯 LLM-Powered Agent Swarm Created:")
	fmt.Printf("   📋 %s - Coordinates projects and delegates to specialists\n", coordinator.Name)
	fmt.Printf("   ✍️  %s - Handles content creation and research\n", contentCreator.Name)
	fmt.Printf("   📊 %s - Performs data analysis and reporting\n", dataAnalyst.Name)
	fmt.Printf("   🧠 %s - Manages shared information and context\n", memoryManager.Name)
	fmt.Println()

	// Run comprehensive demo scenarios
	runComprehensiveDemoScenarios(swarmClient, coordinator)
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

	// Task coordination tools
	server.GetToolRegistry().Register("create_task", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		task := "New task"
		assignee := "Team member"

		if t, ok := params["task"].(string); ok {
			task = t
		}
		if a, ok := params["assignee"].(string); ok {
			assignee = a
		}

		return fmt.Sprintf("📋 Task created: '%s' assigned to %s", task, assignee), nil
	})

	// Memory management tools (explicitly register these)
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
		`You are the Coordinator agent responsible for project management and task delegation.

Your role:
- Analyze complex projects and break them down into manageable tasks
- Route tasks to appropriate specialist agents based on their expertise
- Coordinate multi-agent workflows and ensure project completion
- Store project context and track progress

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names only: "store_context" or "create_task"
- Wait for tool response before calling another tool
- Example of WRONG usage: "store_context|create_task" 
- Example of CORRECT usage: First call "store_context", then call "create_task"

CRITICAL HANDOFF RULES:
- ALWAYS transfer to ONE AGENT AT A TIME - never multiple agents
- NEVER try to handoff to multiple agents like "DataAnalyst, MemoryManager"
- Use exact function names: "transfer_to_content_creator", "transfer_to_data_analyst", "transfer_to_memory_manager"
- Example of WRONG usage: transfer to "DataAnalyst, MemoryManager"
- Example of CORRECT usage: First transfer to "transfer_to_data_analyst", then that agent can transfer to "transfer_to_memory_manager"

Available tools (call them individually):
- store_context: Store important project information and context
- create_task: Create new tasks and assign them to team members

Available handoff functions for delegation:
- transfer_to_content_creator: For articles, research, writing tasks
- transfer_to_data_analyst: For data analysis, insights, reports  
- transfer_to_memory_manager: For information storage/retrieval

When you receive a complex project:
1. First call "store_context" to save project details
2. Then call "create_task" to break down the work
3. Finally delegate by calling the appropriate transfer function
4. For content creation projects → call "transfer_to_content_creator"
5. For data analysis projects → call "transfer_to_data_analyst"

Always delegate complex work to specialist agents after initial coordination.`,
		[]string{"store_context", "create_task"})
}

func createContentCreator(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"ContentCreator",
		`You are the ContentCreator agent specialized in content creation and research.

Your capabilities:
- Writing articles, blog posts, and other content
- Researching topics thoroughly and gathering information
- Text processing and content optimization
- Information synthesis and presentation

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "research_topic", "write_article", "store_context"
- Wait for each tool response before calling another tool

Your workflow:
1. First call ONLY "research_topic" to gather information
2. Wait for response, then call ONLY "write_article" to create content  
3. Finally call ONLY "store_context" to save important findings

Use your tools efficiently and ensure content quality.`,
		[]string{"research_topic", "write_article", "store_context"})
}

func createDataAnalyst(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"DataAnalyst",
		`You are the DataAnalyst agent specialized in data analysis and reporting.

Your capabilities:
- Analyzing datasets and extracting insights
- Generating comprehensive reports with findings
- Statistical analysis and pattern recognition
- Data visualization recommendations
- Trend identification and forecasting

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "analyze_data", "generate_report", "store_context"
- Wait for each tool response before calling another tool

Your workflow:
1. First call ONLY "analyze_data" to examine the dataset
2. Wait for response, then call ONLY "generate_report" with findings
3. Finally call ONLY "store_context" to save important insights

Focus on providing actionable insights and clear recommendations.`,
		[]string{"analyze_data", "generate_report", "store_context"})
}

func createMemoryManager(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"MemoryManager",
		`You are the MemoryManager agent responsible for information management across the swarm.

Your capabilities:
- Storing and organizing shared information
- Retrieving relevant context for other agents
- Managing project knowledge and continuity
- Providing summaries of stored information
- Cross-referencing related information

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "store_context", "retrieve_context"
- Wait for each tool response before calling another tool

Your role is crucial for maintaining context across conversations and ensuring knowledge sharing between agents.`,
		[]string{"store_context", "retrieve_context"})
}

func setupAgentHandoffs(client swarm.SwarmClient, coordinator, contentCreator, dataAnalyst, memoryManager *swarm.Agent) {
	fmt.Println("🔗 Setting up agent handoff functions...")

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

	fmt.Println("✅ Agent handoff functions registered")
}

func runComprehensiveDemoScenarios(client swarm.SwarmClient, coordinator *swarm.Agent) {
	fmt.Println("🚀 Running LLM-Powered Comprehensive Demo Scenarios:")
	fmt.Println("==================================================")

	scenarios := []struct {
		name        string
		message     string
		description string
	}{
		{
			"Simple Tool Test",
			"Please use the store_context tool to store a simple message",
			"Tests basic tool usage",
		},
		{
			"Agent Handoff Test",
			"Please use the transfer_to_content_creator function to transfer me to the ContentCreator agent",
			"Tests agent handoff functionality",
		},
		{
			"Content Creation Project",
			"I need a comprehensive article about artificial intelligence in healthcare. Please coordinate this project with proper research and content creation.",
			"Tests LLM project coordination and content creation workflow",
		},
		{
			"Data Analysis Project",
			"I have a customer behavior dataset from Q4 2023 that needs thorough analysis. Please coordinate the analysis and create a comprehensive report with insights.",
			"Tests LLM data analysis coordination and reporting workflow",
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
		response := client.Run(coordinator, messages, map[string]interface{}{
			"scenario":     scenario.name,
			"project_type": "coordinated_workflow",
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
					// Truncate long messages for readability
					content := msg.Content
					if len(content) > 200 {
						content = content[:200] + "..."
					}
					fmt.Printf("   %s %s\n", role, content)
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
	fmt.Println("   🧠 LLM-powered project coordination and task delegation")
	fmt.Println("   🎯 Intelligent agent selection based on task requirements")
	fmt.Println("   🔧 Smart tool usage with natural language reasoning")
	fmt.Println("   🔄 Context-aware agent handoffs and collaboration")
	fmt.Println("   💭 Transparent reasoning for all decisions")
	fmt.Println("   📚 Persistent memory and knowledge sharing")

	fmt.Println("\n🔧 LLM Integration Details:")
	fmt.Println("   • Each agent has specialized LLM instructions")
	fmt.Println("   • LLM analyzes context and selects appropriate actions")
	fmt.Println("   • Decisions include reasoning for transparency")
	fmt.Println("   • No rule-based fallback - pure LLM reasoning")
	fmt.Println("   • Ollama integration for local AI processing")

	fmt.Println("\n🚀 Next Steps:")
	fmt.Println("   • Customize agent roles for your specific domain")
	fmt.Println("   • Add specialized tools for your workflows")
	fmt.Println("   • Experiment with different Ollama models")
	fmt.Println("   • Scale to larger swarms with more specialized agents")
}
