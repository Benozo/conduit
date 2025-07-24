package main

import (
	"fmt"
	"log"

	"github.com/benozo/neuron/src/agents/base"
	"github.com/benozo/neuron/src/core"
	"github.com/benozo/neuron/src/llm/providers"
)

// AISpecialistAgent represents an AI-powered specialist using Ollama
type AISpecialistAgent struct {
	*core.BaseAgent
	ollamaProvider *providers.OllamaProvider
}

// NewAISpecialistAgent creates a new AI specialist agent
func NewAISpecialistAgent(name, specialization string, ollamaURL, model string) *AISpecialistAgent {
	ollamaProvider := providers.NewOllamaProvider(ollamaURL, model)
	baseAgent := core.NewAgent(name, specialization, nil, model, "auto")

	return &AISpecialistAgent{
		BaseAgent:      baseAgent,
		ollamaProvider: ollamaProvider,
	}
}

// ProcessTask processes a task using AI capabilities
func (asa *AISpecialistAgent) ProcessTask(task string) (string, error) {
	prompt := fmt.Sprintf("As a %s, please process this task: %s", asa.Instructions, task)
	return asa.ollamaProvider.GenerateResponse(prompt)
}

// Execute implements the Agent interface
func (asa *AISpecialistAgent) Execute(input interface{}) (interface{}, error) {
	task, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string input, got %T", input)
	}
	return asa.ProcessTask(task)
}

// GetModelInfo returns information about the agent's model
func (asa *AISpecialistAgent) GetModelInfo() string {
	info := asa.ollamaProvider.GetModelInfo()
	return fmt.Sprintf("Model: %s | Provider: %s", info.Name, info.Provider)
}

// Coordinator Demo Example with Ollama Integration
// This example demonstrates how to use the coordinator to manage both traditional and AI-powered agents.
func main() {
	fmt.Println("=== Enhanced Coordinator Demo with Ollama ===")

	// Configuration
	ollamaURL := "http://192.168.10.10:11434"
	model := "llama3.2"

	// Create a coordinator
	coordinator := base.NewCoordinator()

	// Create traditional specialist agents
	traditionalAnalyst := base.NewSpecialist("TraditionalAnalyst", "Analyze and process data using traditional methods")
	reportGenerator := base.NewSpecialist("ReportGenerator", "Generate detailed reports")
	qualityController := base.NewSpecialist("QualityController", "Ensure quality standards")

	// Create AI-powered specialist
	aiAdvisor := NewAISpecialistAgent(
		"AIAdvisor",
		"intelligent advisor providing AI-powered insights and recommendations",
		ollamaURL, model)

	// Test AI connection
	fmt.Printf("🔍 Testing AI advisor connection to Ollama at %s...\n", ollamaURL)
	if err := aiAdvisor.ollamaProvider.Ping(); err != nil {
		log.Printf("❌ Failed to connect to Ollama: %v", err)
		log.Println("Continuing with traditional agents only...")
	} else {
		fmt.Println("✅ AI advisor connected to Ollama!")
	}
	// Register traditional agents with the coordinator
	fmt.Println("\n🏗️  Registering agents with coordinator...")
	coordinator.RegisterAgent(traditionalAnalyst)
	coordinator.RegisterAgent(reportGenerator)
	coordinator.RegisterAgent(qualityController)

	// Display coordinator status
	fmt.Printf("\n📊 Coordinator Status: %s\n", coordinator.GetStatus())
	fmt.Printf("📋 Registered Traditional Agents: %d\n", len(coordinator.GetAgents()))
	fmt.Printf("🤖 AI-Powered Agents: 1 (managed separately)\n")

	// List all traditional agents
	fmt.Println("\n📋 Registered Traditional Agents:")
	for name, agent := range coordinator.GetAgents() {
		fmt.Printf("  - %s: %s\n", name, agent.GetName())
	}

	// Show AI agent info
	fmt.Printf("\n🤖 AI-Powered Agent:\n")
	fmt.Printf("  - %s: %s\n", aiAdvisor.GetName(), aiAdvisor.GetModelInfo())

	// Test AI agent with a task
	if aiAdvisor.ollamaProvider != nil {
		fmt.Println("\n🧪 Testing AI advisor with a coordination task...")
		task := "Recommend the best approach for coordinating a team of data analysts, report generators, and quality controllers for a large data processing project."

		fmt.Printf("🤖 %s processing coordination task...\n", aiAdvisor.GetName())
		result, err := aiAdvisor.ProcessTask(task)
		if err != nil {
			log.Printf("❌ AI advisor task failed: %v", err)
		} else {
			fmt.Printf("💡 AI Recommendation (first 300 chars):\n%s...\n", result[:min(300, len(result))])
		}
	}

	// Get metrics
	metrics := coordinator.GetMetrics()
	fmt.Printf("\n📈 Coordinator Metrics:\n")
	fmt.Printf("  - Traditional Agents Registered: %d\n", metrics.AgentsRegistered)
	fmt.Printf("  - Workflows Executed: %d\n", metrics.WorkflowsExecuted)
	fmt.Printf("  - Uptime: %v\n", metrics.Uptime)

	fmt.Println("\n🎉 Enhanced coordinator demo completed successfully!")
	fmt.Println("Demonstrated coordination of both traditional and AI-powered agents!")
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
