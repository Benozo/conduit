package main

import (
	"fmt"
	"log"

	"github.com/benozo/neuron/src/agents/base"
	"github.com/benozo/neuron/src/core"
	"github.com/benozo/neuron/src/llm/providers"
)

// AIAssistantAgent represents an AI-powered agent using Ollama
type AIAssistantAgent struct {
	*core.BaseAgent
	ollamaProvider *providers.OllamaProvider
}

// NewAIAssistantAgent creates a new AI assistant agent
func NewAIAssistantAgent(name, specialization string, ollamaURL, model string) *AIAssistantAgent {
	ollamaProvider := providers.NewOllamaProvider(ollamaURL, model)
	baseAgent := core.NewAgent(name, specialization, nil, model, "auto")

	return &AIAssistantAgent{
		BaseAgent:      baseAgent,
		ollamaProvider: ollamaProvider,
	}
}

// Consult asks the AI assistant for advice or analysis
func (aia *AIAssistantAgent) Consult(query string) (string, error) {
	prompt := fmt.Sprintf("As a %s, please help with this request: %s", aia.Instructions, query)
	return aia.ollamaProvider.GenerateResponse(prompt)
}

// Execute implements the Agent interface
func (aia *AIAssistantAgent) Execute(input interface{}) (interface{}, error) {
	query, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string input, got %T", input)
	}
	return aia.Consult(query)
}

// GetModelInfo returns information about the agent's model
func (aia *AIAssistantAgent) GetModelInfo() string {
	info := aia.ollamaProvider.GetModelInfo()
	return fmt.Sprintf("Model: %s | Provider: %s", info.Name, info.Provider)
}

func main() {
	fmt.Println("=== Multi-Agent Ollama System Demo ===")

	// Configuration
	ollamaURL := "http://192.168.10.10:11434"
	analystModel := "llama3.2"
	writerModel := "tinyLlama"

	// Create coordinator
	coordinator := base.NewCoordinator()

	// Create AI-powered specialists with different models
	aiAnalyst := NewAIAssistantAgent(
		"AIDataAnalyst",
		"expert data analyst specializing in pattern recognition and insights",
		ollamaURL, analystModel)

	aiWriter := NewAIAssistantAgent(
		"AIContentWriter",
		"professional content writer creating clear and engaging text",
		ollamaURL, writerModel)

	// Create traditional agents
	projectManager := base.NewSpecialist("ProjectManager", "Coordinate and manage project tasks")
	qualityController := base.NewSpecialist("QualityController", "Ensure deliverable quality")

	// Register traditional agents with coordinator
	fmt.Println("🏗️  Setting up multi-agent system...")
	coordinator.RegisterAgent(projectManager)
	coordinator.RegisterAgent(qualityController)

	fmt.Printf("✅ Coordinator managing %d traditional agents\n", len(coordinator.GetAgents()))
	fmt.Printf("🤖 AI Analyst using model: %s\n", analystModel)
	fmt.Printf("🤖 AI Writer using model: %s\n", writerModel)

	// Test connection to AI agents
	fmt.Printf("\n🔍 Testing AI agent connections...\n")
	fmt.Printf("Testing %s with %s...\n", aiAnalyst.GetName(), analystModel)
	if err := aiAnalyst.ollamaProvider.Ping(); err != nil {
		log.Printf("❌ Failed to connect to Ollama for analyst: %v", err)
		return
	}
	fmt.Printf("✅ %s connected!\n", aiAnalyst.GetName())

	fmt.Printf("Testing %s with %s...\n", aiWriter.GetName(), writerModel)
	if err := aiWriter.ollamaProvider.Ping(); err != nil {
		log.Printf("❌ Failed to connect to Ollama for writer: %v", err)
		return
	}
	fmt.Printf("✅ %s connected!\n", aiWriter.GetName())
	fmt.Println("✅ All AI agents connected to Ollama!")

	// Simulate a simple workflow with AI assistance
	fmt.Println("\n📋 Simulating workflow with AI assistance...")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	// Step 1: AI Data Analysis
	fmt.Println("\n📊 Step 1: AI Data Analysis")
	analysisQuery := "Analyze key trends in mobile app user engagement. What are the top 3 factors affecting user retention?"

	fmt.Printf("🤖 %s (%s) processing...\n", aiAnalyst.GetName(), analystModel)
	analysis, err := aiAnalyst.Consult(analysisQuery)
	if err != nil {
		log.Printf("❌ Analysis failed: %v", err)
		return
	}
	fmt.Printf("📈 Analysis Result (first 200 chars):\n%s...\n", analysis[:min(200, len(analysis))])

	// Step 2: Content Creation
	fmt.Println("\n✍️  Step 2: AI Content Creation")
	contentQuery := "Create a brief executive summary about mobile app user retention strategies."

	fmt.Printf("🤖 %s (%s) creating content...\n", aiWriter.GetName(), writerModel)
	content, err := aiWriter.Consult(contentQuery)
	if err != nil {
		log.Printf("❌ Content creation failed: %v", err)
		return
	}
	fmt.Printf("📝 Content Summary (first 200 chars):\n%s...\n", content[:min(200, len(content))])

	// Show final status
	fmt.Println("\n📊 System Status:")
	fmt.Printf("Traditional Agents: %d active\n", len(coordinator.GetAgents()))
	fmt.Printf("AI Agents: 2 active and responding\n")
	fmt.Printf("  - %s using %s\n", aiAnalyst.GetName(), analystModel)
	fmt.Printf("  - %s using %s\n", aiWriter.GetName(), writerModel)

	fmt.Println("\n🎉 Multi-agent system demo completed!")
	fmt.Println("Successfully demonstrated multi-model AI-traditional agent collaboration!")
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
