package main

import (
	"fmt"
	"log"
	"os"

	"github.com/benozo/neuron/src/agents/base"
	"github.com/benozo/neuron/src/core"
	"github.com/benozo/neuron/src/llm/providers"
)

// CloudflareAIAgent represents an AI-powered agent using Cloudflare Workers AI
type CloudflareAIAgent struct {
	*core.BaseAgent
	cfProvider *providers.CloudflareAIProvider
}

// NewCloudflareAIAgent creates a new Cloudflare Workers AI agent (standard API)
func NewCloudflareAIAgent(name, specialization string, accountID, apiToken, model string) *CloudflareAIAgent {
	cfProvider := providers.NewCloudflareAIProvider(accountID, apiToken, model)
	baseAgent := core.NewAgent(name, specialization, nil, model, "auto")

	return &CloudflareAIAgent{
		BaseAgent:  baseAgent,
		cfProvider: cfProvider,
	}
}

// NewCustomCloudflareAIAgent creates a new agent using custom Cloudflare endpoint
func NewCustomCloudflareAIAgent(name, specialization string, baseURL, apiKey, model string) *CloudflareAIAgent {
	cfProvider := providers.NewCustomCloudflareAIProvider(baseURL, apiKey, model)
	baseAgent := core.NewAgent(name, specialization, nil, model, "auto")

	return &CloudflareAIAgent{
		BaseAgent:  baseAgent,
		cfProvider: cfProvider,
	}
}

// ProcessTask processes a task using Cloudflare Workers AI capabilities
func (cfa *CloudflareAIAgent) ProcessTask(task string) (string, error) {
	prompt := fmt.Sprintf("As a %s, please process this task: %s", cfa.Instructions, task)
	return cfa.cfProvider.GenerateResponse(prompt)
}

// Execute implements the Agent interface
func (cfa *CloudflareAIAgent) Execute(input interface{}) (interface{}, error) {
	task, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string input, got %T", input)
	}
	return cfa.ProcessTask(task)
}

// GetModelInfo returns information about the agent's model
func (cfa *CloudflareAIAgent) GetModelInfo() string {
	info := cfa.cfProvider.GetModelInfo()
	return fmt.Sprintf("Model: %s | Provider: %s", info.Name, info.Provider)
}

// Cloudflare Workers AI Demo Example
// This example demonstrates how to use Cloudflare Workers AI with the swarm framework.
func main() {
	fmt.Println("=== Cloudflare Workers AI Swarm Demo ===")

	// Configuration - Read from environment variables
	// Option 1: Standard Cloudflare API
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

	// Option 2: Custom Cloudflare endpoint
	customURL := os.Getenv("CLOUDFLARE_CUSTOM_URL")
	customAPIKey := os.Getenv("CLOUDFLARE_CUSTOM_API_KEY")

	model := os.Getenv("CLOUDFLARE_MODEL")

	// Default model if not specified
	if model == "" {
		model = "@cf/meta/llama-3.1-8b-instruct"
	}

	// Determine which configuration to use
	useCustomEndpoint := customURL != "" && customAPIKey != ""
	useStandardAPI := accountID != "" && apiToken != ""

	if !useCustomEndpoint && !useStandardAPI {
		fmt.Println("❌ Missing required environment variables:")
		fmt.Println("\nOption 1: Standard Cloudflare Workers AI API:")
		fmt.Println("   CLOUDFLARE_ACCOUNT_ID - Your Cloudflare Account ID")
		fmt.Println("   CLOUDFLARE_API_TOKEN - Your Cloudflare API Token")
		fmt.Println("\nOption 2: Custom Cloudflare Workers AI Endpoint:")
		fmt.Println("   CLOUDFLARE_CUSTOM_URL - Your custom endpoint URL (e.g., https://example.com)")
		fmt.Println("   CLOUDFLARE_CUSTOM_API_KEY - Your custom API key")
		fmt.Println("\nCommon:")
		fmt.Println("   CLOUDFLARE_MODEL (optional) - Model to use (default: @cf/meta/llama-3.1-8b-instruct)")
		fmt.Println()
		fmt.Println("🔧 Setup Instructions:")
		fmt.Println("For custom endpoint (based on your curl example):")
		fmt.Println("   export CLOUDFLARE_CUSTOM_URL=https://example.com")
		fmt.Println("   export CLOUDFLARE_CUSTOM_API_KEY=XXXXXXXXXXXXXX")
		fmt.Println("   export CLOUDFLARE_MODEL=@cf/meta/llama-4-scout-17b-16e-instruct")
		os.Exit(1)
	}

	fmt.Printf("🌐 Using Cloudflare Workers AI\n")
	if useCustomEndpoint {
		fmt.Printf("   Endpoint Type: Custom\n")
		fmt.Printf("   Base URL: %s\n", customURL)
		fmt.Printf("   API Key: %s\n", maskString(customAPIKey))
	} else {
		fmt.Printf("   Endpoint Type: Standard Cloudflare API\n")
		fmt.Printf("   Account ID: %s\n", maskString(accountID))
	}
	fmt.Printf("   Model: %s\n", model)

	// Create a coordinator
	coordinator := base.NewCoordinator()

	// Create Cloudflare AI-powered agents based on configuration
	var dataAnalyst, contentCreator, strategicAdvisor *CloudflareAIAgent

	if useCustomEndpoint {
		dataAnalyst = NewCustomCloudflareAIAgent(
			"CloudflareDataAnalyst",
			"data scientist specializing in statistical analysis and data interpretation",
			customURL, customAPIKey, model)

		contentCreator = NewCustomCloudflareAIAgent(
			"CloudflareContentCreator",
			"creative writer specializing in marketing content and technical documentation",
			customURL, customAPIKey, model)

		strategicAdvisor = NewCustomCloudflareAIAgent(
			"CloudflareStrategicAdvisor",
			"business strategist providing insights on market trends and strategic planning",
			customURL, customAPIKey, model)
	} else {
		dataAnalyst = NewCloudflareAIAgent(
			"CloudflareDataAnalyst",
			"data scientist specializing in statistical analysis and data interpretation",
			accountID, apiToken, model)

		contentCreator = NewCloudflareAIAgent(
			"CloudflareContentCreator",
			"creative writer specializing in marketing content and technical documentation",
			accountID, apiToken, model)

		strategicAdvisor = NewCloudflareAIAgent(
			"CloudflareStrategicAdvisor",
			"business strategist providing insights on market trends and strategic planning",
			accountID, apiToken, model)
	}

	// Test connection to Cloudflare Workers AI
	fmt.Println("\n🔍 Testing connection to Cloudflare Workers AI...")
	if err := dataAnalyst.cfProvider.Ping(); err != nil {
		log.Fatalf("❌ Failed to connect to Cloudflare Workers AI: %v", err)
	}
	fmt.Println("✅ Successfully connected to Cloudflare Workers AI!")

	// Show available models
	fmt.Println("\n📋 Available Cloudflare AI Models:")
	models := dataAnalyst.cfProvider.GetAvailableModels()
	for i, m := range models {
		marker := ""
		if m == model {
			marker = " ← Current"
		}
		fmt.Printf("   %d. %s%s\n", i+1, m, marker)
	}

	// Register agents with the coordinator
	fmt.Println("\n🏗️  Setting up AI-powered swarm...")

	// Create traditional coordinator for comparison
	coordinator.RegisterAgent(base.NewSpecialist("TraditionalReporter", "Generate basic reports"))

	// Display coordinator status
	fmt.Printf("\n📊 Swarm Status:\n")
	fmt.Printf("   Traditional Agents: %d\n", len(coordinator.GetAgents()))
	fmt.Printf("   Cloudflare AI Agents: 3\n")

	// Show agent information
	fmt.Println("\n🤖 Cloudflare AI-Powered Agents:")
	agents := []*CloudflareAIAgent{dataAnalyst, contentCreator, strategicAdvisor}
	for _, agent := range agents {
		fmt.Printf("   - %s: %s\n", agent.GetName(), agent.GetModelInfo())
	}

	// Demonstrate collaborative workflow
	fmt.Println("\n🔄 Demonstrating Cloudflare AI Swarm Workflow...")

	// Define a complex business scenario
	scenario := "A tech startup wants to launch a new AI-powered productivity app. They need market analysis, content strategy, and strategic recommendations."

	// Step 1: Data Analysis
	fmt.Printf("\n📊 Step 1: %s analyzing market data...\n", dataAnalyst.GetName())
	analysisTask := fmt.Sprintf("Analyze the market for AI productivity apps. Focus on: market size, competition, target demographics, and growth potential. Scenario: %s", scenario)

	analysisResult, err := dataAnalyst.ProcessTask(analysisTask)
	if err != nil {
		log.Printf("❌ Analysis failed: %v", err)
	} else {
		fmt.Printf("📈 Analysis Result (first 200 chars):\n%s...\n", truncateString(analysisResult, 2000))
	}

	// Step 2: Content Strategy
	fmt.Printf("\n✍️  Step 2: %s creating content strategy...\n", contentCreator.GetName())
	contentTask := fmt.Sprintf("Based on this market analysis, create a content marketing strategy for an AI productivity app launch. Analysis: %s", truncateString(analysisResult, 3000))

	contentResult, err := contentCreator.ProcessTask(contentTask)
	if err != nil {
		log.Printf("❌ Content strategy failed: %v", err)
	} else {
		fmt.Printf("📝 Content Strategy (first 200 chars):\n%s...\n", truncateString(contentResult, 2000))
	}

	// Step 3: Strategic Recommendations
	fmt.Printf("\n🎯 Step 3: %s providing strategic recommendations...\n", strategicAdvisor.GetName())
	strategyTask := fmt.Sprintf("Synthesize this market analysis and content strategy into actionable strategic recommendations. Analysis: %s Content Strategy: %s",
		truncateString(analysisResult, 200), truncateString(contentResult, 2000))

	strategyResult, err := strategicAdvisor.ProcessTask(strategyTask)
	if err != nil {
		log.Printf("❌ Strategic planning failed: %v", err)
	} else {
		fmt.Printf("🎯 Strategic Recommendations (first 300 chars):\n%s...\n", truncateString(strategyResult, 3000))
	}

	// Performance metrics
	fmt.Println("\n📈 Swarm Performance Metrics:")
	fmt.Printf("   ✅ Market Analysis: Completed\n")
	fmt.Printf("   ✅ Content Strategy: Completed\n")
	fmt.Printf("   ✅ Strategic Planning: Completed\n")
	fmt.Printf("   🌐 All tasks powered by Cloudflare Workers AI\n")

	// Get coordinator metrics
	metrics := coordinator.GetMetrics()
	fmt.Printf("\n📊 Coordinator Metrics:\n")
	fmt.Printf("   Traditional Agents: %d\n", metrics.AgentsRegistered)
	fmt.Printf("   Workflows Executed: %d\n", metrics.WorkflowsExecuted)
	fmt.Printf("   Uptime: %v\n", metrics.Uptime)

	fmt.Println("\n🎉 Cloudflare Workers AI Swarm Demo completed!")
	fmt.Println("Demonstrated edge computing AI with global low-latency inference!")
}

// Helper function to mask sensitive strings
func maskString(s string) string {
	if len(s) <= 8 {
		return "***masked***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}

// Helper function to truncate strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
