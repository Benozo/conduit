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
	fmt.Println("🧠 LLM-Powered Advanced Agent Swarm Workflow Patterns Demo")
	fmt.Println("=========================================================")
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

	// Register custom memory tools
	registerCustomMemoryTools(server)

	// Register workflow-specific tools
	registerWorkflowTools(server)

	// Create Ollama model function
	ollamaModel := conduit.CreateOllamaModel(ollamaURL)

	// Create swarm client with LLM integration (no fallback)
	swarmClient := swarm.NewSwarmClientWithLLM(server, swarm.DefaultSwarmConfig(), ollamaModel, modelName)

	// Create specialized workflow agents
	orchestrator := createOrchestrator(swarmClient)
	dataProcessor := createDataProcessor(swarmClient)
	analyst := createWorkflowAnalyst(swarmClient)
	reporter := createReporter(swarmClient)
	qualityController := createQualityController(swarmClient)

	// Set up complex agent handoffs
	setupWorkflowHandoffs(swarmClient, orchestrator, dataProcessor, analyst, reporter, qualityController)

	fmt.Println("🎯 LLM-Powered Advanced Workflow Swarm Created:")
	fmt.Printf("   🎼 %s - Orchestrates complex multi-step workflows\n", orchestrator.Name)
	fmt.Printf("   🔄 %s - Handles data extraction, transformation, and loading\n", dataProcessor.Name)
	fmt.Printf("   📈 %s - Performs trend analysis and generates insights\n", analyst.Name)
	fmt.Printf("   📄 %s - Creates comprehensive reports and notifications\n", reporter.Name)
	fmt.Printf("   ✅ %s - Ensures quality and handles exceptions\n", qualityController.Name)
	fmt.Println()

	// Run advanced workflow scenarios
	runAdvancedWorkflowScenarios(swarmClient, orchestrator)
}

// registerCustomMemoryTools registers custom memory tools to avoid conflicts
func registerCustomMemoryTools(server *conduit.EnhancedServer) {
	server.GetToolRegistry().Register("store_context", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key, ok := params["key"].(string)
		if !ok {
			return nil, fmt.Errorf("key parameter is required")
		}

		value, ok := params["value"]
		if !ok {
			return nil, fmt.Errorf("value parameter is required")
		}

		// Store in memory
		memory.Set(key, value)

		return fmt.Sprintf("Stored context: %s", key), nil
	})

	server.GetToolRegistry().Register("retrieve_context", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key, ok := params["key"].(string)
		if !ok {
			return nil, fmt.Errorf("key parameter is required")
		}

		value := memory.Get(key)
		if value == nil {
			return fmt.Sprintf("No context found for key: %s", key), nil
		}

		return fmt.Sprintf("Retrieved context for %s: %v", key, value), nil
	})
}

func registerWorkflowTools(server *conduit.EnhancedServer) {
	// Data processing pipeline tools
	server.GetToolRegistry().Register("extract_data", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		source := "default_source"
		if s, ok := params["source"].(string); ok {
			source = s
		}

		return fmt.Sprintf("📥 Extracted data from %s: [records: 1000, quality: high, format: structured]", source), nil
	})

	server.GetToolRegistry().Register("transform_data", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		data := "extracted_data"
		if d, ok := params["data"].(string); ok {
			data = d
		}

		return fmt.Sprintf("🔄 Transformed %s: normalized, cleaned, validated, enriched", data), nil
	})

	server.GetToolRegistry().Register("validate_data", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		data := "transformed_data"
		if d, ok := params["data"].(string); ok {
			data = d
		}

		// Store quality result in memory
		memory.Set("data_quality", "high")
		return fmt.Sprintf("✅ Validated %s: 98%% quality score, ready for analysis", data), nil
	})

	server.GetToolRegistry().Register("load_data", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		data := "validated_data"
		destination := "data_warehouse"

		if d, ok := params["data"].(string); ok {
			data = d
		}
		if dest, ok := params["destination"].(string); ok {
			destination = dest
		}

		return fmt.Sprintf("📤 Loaded %s to %s successfully", data, destination), nil
	})

	// Analysis tools
	server.GetToolRegistry().Register("analyze_trends", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		data := "processed_data"
		if d, ok := params["data"].(string); ok {
			data = d
		}

		return fmt.Sprintf("📈 Trend analysis complete for %s: upward trend detected (15%% growth), seasonal patterns identified", data), nil
	})

	server.GetToolRegistry().Register("generate_insights", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		analysis := "trend_analysis"
		if a, ok := params["analysis"].(string); ok {
			analysis = a
		}

		insights := "Key insights: 1) Customer engagement up 15%, 2) Peak activity on weekends, 3) Mobile usage growing, 4) Geographic expansion opportunities, 5) Product recommendation engine needed"
		memory.Set("insights", insights)
		return fmt.Sprintf("💡 Generated insights from %s: %s", analysis, insights), nil
	})

	// Reporting tools
	server.GetToolRegistry().Register("create_report", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		insights := "analysis_insights"
		reportType := "standard"

		if i, ok := params["insights"].(string); ok {
			insights = i
		}
		if rt, ok := params["report_type"].(string); ok {
			reportType = rt
		}

		return fmt.Sprintf("📊 Created %s report with %s: Executive summary, detailed findings, recommendations, and action items", reportType, insights), nil
	})

	server.GetToolRegistry().Register("send_notification", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		report := "generated_report"
		recipients := "stakeholders"

		if r, ok := params["report"].(string); ok {
			report = r
		}
		if rec, ok := params["recipients"].(string); ok {
			recipients = rec
		}

		return fmt.Sprintf("📧 Sent notification to %s: %s is ready for review", recipients, report), nil
	})

	// Quality control tools
	server.GetToolRegistry().Register("check_quality", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		data := "input_data"
		if d, ok := params["data"].(string); ok {
			data = d
		}

		quality := "high"
		memory.Set("data_quality", quality)
		return fmt.Sprintf("🔍 Quality check for %s: PASSED (%s quality) - meets all standards", data, quality), nil
	})

	server.GetToolRegistry().Register("emergency_process", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		issue := "quality_issue"
		if i, ok := params["issue"].(string); ok {
			issue = i
		}

		return fmt.Sprintf("🚨 Emergency processing for %s: issue resolved, fallback procedures activated", issue), nil
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

func createOrchestrator(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"Orchestrator",
		`You are the Orchestrator agent responsible for managing complex multi-step workflows.

Your role:
- Analyze complex workflow requirements and design execution plans
- Coordinate between multiple agents to ensure smooth workflow execution
- Monitor workflow progress and handle exceptions
- Make decisions about workflow routing and sequencing

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "store_context", "retrieve_context", "timestamp"
- Use exact handoff function names: "data_processor", "analyst", "reporter", "quality_controller"
- Wait for each tool response before calling another tool

Workflow coordination strategy:
1. First call "store_context" to save workflow plan and context
2. For ETL/data processing tasks: call "data_processor" handoff function
3. For analysis/insights tasks: call "analyst" handoff function  
4. For reporting/presentation tasks: call "reporter" handoff function
5. For quality control tasks: call "quality_controller" handoff function
6. Use "retrieve_context" and "timestamp" tools as needed

IMPORTANT: You should primarily delegate tasks to specialized agents rather than trying to handle everything yourself. Always use the handoff functions to route work to the appropriate agent.`,
		[]string{"store_context", "retrieve_context", "timestamp", "data_processor", "analyst", "reporter", "quality_controller"})
}

func createDataProcessor(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"DataProcessor",
		`You are the DataProcessor agent specialized in ETL (Extract, Transform, Load) operations.

Your capabilities:
- Data extraction from various sources
- Data transformation and normalization
- Data validation and quality assurance
- Data loading to target systems

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "extract_data", "transform_data", "validate_data", "load_data", "store_context", "encode_base64"
- Wait for each tool response before calling another tool

Your workflow:
1. First call "extract_data" from specified sources
2. Then call "transform_data" to process the extracted data
3. Next call "validate_data" to check quality and completeness
4. Then call "load_data" to target destinations
5. Finally call "store_context" to save processing status and metrics

Focus on data quality and efficient processing. Always validate before proceeding to next steps.`,
		[]string{"extract_data", "transform_data", "validate_data", "load_data", "store_context", "encode_base64"})
}

func createWorkflowAnalyst(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"Analyst",
		`You are the Analyst agent specialized in advanced data analysis and trend identification.

Your capabilities:
- Trend analysis and pattern recognition
- Statistical analysis and forecasting
- Insight generation and recommendation development
- Performance metrics and KPI analysis

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "analyze_trends", "generate_insights", "store_context", "uuid_generate", "timestamp"
- Wait for each tool response before calling another tool

Your workflow:
1. First call "analyze_trends" to examine processed data for trends and patterns
2. Then call "generate_insights" to create actionable insights from analysis
3. Finally call "store_context" to save insights for reporting and future reference

Focus on providing valuable business insights and actionable recommendations.`,
		[]string{"analyze_trends", "generate_insights", "store_context", "uuid_generate", "timestamp"})
}

func createReporter(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"Reporter",
		`You are the Reporter agent specialized in creating comprehensive reports and communications.

Your capabilities:
- Report generation from insights and analysis
- Stakeholder communication and notifications
- Executive summary creation
- Multi-format report delivery

CRITICAL TOOL USAGE RULES:
- ALWAYS call tools ONE AT A TIME - never combine multiple tools
- NEVER use pipe symbols (|) or ampersands (&) in tool calls
- Use exact tool names: "create_report", "send_notification", "store_context", "uppercase", "lowercase"
- Wait for each tool response before calling another tool

Your workflow:
1. First call "create_report" to generate structured reports with findings
2. Then call "send_notification" to notify relevant stakeholders
3. Finally call "store_context" to save reporting status and results

Focus on clear communication and actionable reporting.`,
		[]string{"create_report", "send_notification", "store_context", "uppercase", "lowercase"})
}

func createQualityController(client swarm.SwarmClient) *swarm.Agent {
	return client.CreateAgent(
		"QualityController",
		`You are the QualityController agent responsible for ensuring workflow quality and handling exceptions.

Your capabilities:
- Quality assurance and validation
- Exception handling and error recovery
- Process monitoring and control
- Emergency response procedures

Your workflow:
1. Perform quality checks at each workflow stage
2. Validate data and process integrity
3. Handle exceptions and quality issues
4. Implement emergency procedures when needed

Focus on maintaining high quality standards and smooth workflow execution.`,
		[]string{"check_quality", "emergency_process", "store_context", "retrieve_context"})
}

func setupWorkflowHandoffs(client swarm.SwarmClient, orchestrator, dataProcessor, analyst, reporter, qualityController *swarm.Agent) {
	// Orchestrator can coordinate with all agents
	client.RegisterFunction(orchestrator.Name, swarm.CreateHandoffFunction("data_processor", dataProcessor))
	client.RegisterFunction(orchestrator.Name, swarm.CreateHandoffFunction("analyst", analyst))
	client.RegisterFunction(orchestrator.Name, swarm.CreateHandoffFunction("reporter", reporter))
	client.RegisterFunction(orchestrator.Name, swarm.CreateHandoffFunction("quality_controller", qualityController))

	// DataProcessor workflow handoffs
	client.RegisterFunction(dataProcessor.Name, swarm.CreateHandoffFunction("orchestrator", orchestrator))
	client.RegisterFunction(dataProcessor.Name, swarm.CreateHandoffFunction("analyst", analyst))
	client.RegisterFunction(dataProcessor.Name, swarm.CreateHandoffFunction("quality_controller", qualityController))

	// Analyst workflow handoffs
	client.RegisterFunction(analyst.Name, swarm.CreateHandoffFunction("orchestrator", orchestrator))
	client.RegisterFunction(analyst.Name, swarm.CreateHandoffFunction("reporter", reporter))
	client.RegisterFunction(analyst.Name, swarm.CreateHandoffFunction("quality_controller", qualityController))

	// Reporter workflow handoffs
	client.RegisterFunction(reporter.Name, swarm.CreateHandoffFunction("orchestrator", orchestrator))
	client.RegisterFunction(reporter.Name, swarm.CreateHandoffFunction("quality_controller", qualityController))

	// QualityController can coordinate back to any agent
	client.RegisterFunction(qualityController.Name, swarm.CreateHandoffFunction("orchestrator", orchestrator))
	client.RegisterFunction(qualityController.Name, swarm.CreateHandoffFunction("data_processor", dataProcessor))
	client.RegisterFunction(qualityController.Name, swarm.CreateHandoffFunction("analyst", analyst))
	client.RegisterFunction(qualityController.Name, swarm.CreateHandoffFunction("reporter", reporter))
}

func runAdvancedWorkflowScenarios(client swarm.SwarmClient, orchestrator *swarm.Agent) {
	fmt.Println("🚀 Running LLM-Powered Advanced Workflow Scenarios:")
	fmt.Println("===================================================")

	scenarios := []struct {
		name        string
		message     string
		description string
	}{
		{
			"ETL Pipeline Workflow",
			"Execute a complete ETL pipeline: extract customer data from CRM, transform and clean it, validate quality, and load to analytics warehouse",
			"Tests LLM orchestration of complex data processing pipeline",
		},
		{
			"Analytics and Reporting Workflow",
			"Analyze Q4 sales trends, generate insights about customer behavior, create executive report, and notify stakeholders",
			"Tests LLM coordination of analysis and reporting workflow",
		},
		{
			"Quality-Controlled Workflow",
			"Process user engagement data with quality checks at each step, handle any data quality issues, and ensure high-quality output",
			"Tests LLM quality control and exception handling",
		},
		{
			"End-to-End Business Intelligence",
			"Complete BI workflow: extract market data, process and analyze for trends, generate strategic insights, create board presentation, and distribute to leadership",
			"Tests LLM complex multi-agent workflow coordination",
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
		response := client.Run(orchestrator, messages, map[string]interface{}{
			"scenario":      scenario.name,
			"workflow_type": "advanced_patterns",
		})

		fmt.Println("📊 Response:")
		if response.Success {
			fmt.Printf("✅ Success! Turns: %d, Tool calls: %d, Handoffs: %d\n",
				response.TotalTurns, response.ToolCallsCount, response.HandoffsCount)

			// Show workflow execution summary
			for j, msg := range response.Messages {
				if j > 0 { // Skip initial user message
					role := "🤖"
					if msg.Role == "user" {
						role = "👤"
					}
					// Truncate long messages for readability
					content := msg.Content
					if len(content) > 150 {
						content = content[:150] + "..."
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

	fmt.Println("\n🎉 LLM-Powered Advanced Workflow Demo Complete!")
	fmt.Println("\n✨ Advanced Features Demonstrated:")
	fmt.Println("   🧠 LLM-powered workflow orchestration")
	fmt.Println("   🔄 Complex multi-step process coordination")
	fmt.Println("   🎯 Intelligent agent routing based on workflow state")
	fmt.Println("   🔧 Context-aware tool selection and execution")
	fmt.Println("   ✅ Quality control and exception handling")
	fmt.Println("   📊 End-to-end business process automation")

	fmt.Println("\n🔧 Advanced LLM Workflow Integration:")
	fmt.Println("   • Orchestrator designs and manages complex workflows")
	fmt.Println("   • DataProcessor handles ETL operations with quality checks")
	fmt.Println("   • Analyst performs advanced analytics and insights")
	fmt.Println("   • Reporter creates comprehensive deliverables")
	fmt.Println("   • QualityController ensures process integrity")
	fmt.Println("   • No rule-based fallback - pure LLM workflow reasoning")

	fmt.Println("\n🚀 Enterprise Applications:")
	fmt.Println("   • Data pipeline automation")
	fmt.Println("   • Business intelligence workflows")
	fmt.Println("   • Quality-controlled processing")
	fmt.Println("   • Multi-stakeholder reporting")
	fmt.Println("   • Exception handling and recovery")
	fmt.Println("   • Scalable workflow orchestration")
}
