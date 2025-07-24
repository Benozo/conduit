package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/benozo/neuron/src/agents/react"
	"github.com/benozo/neuron/src/core"
	"github.com/benozo/neuron/src/llm/providers"
	"github.com/benozo/neuron/src/workflows"
)

// React Workflow Example with Ollama Integration
// This example showcases the interaction between React agents for decision-making and action execution,
// enhanced with Ollama for AI-powered reasoning.
func main() {
	fmt.Println("=== React Loop Workflow Demo with Ollama ===")

	// Setup Ollama provider
	fmt.Println("🔍 Setting up Ollama integration...")
	ollamaHost := "http://192.168.10.10:11434"
	ollamaModel := "llama3.2"

	ollama := providers.NewOllamaProvider(ollamaHost, ollamaModel)

	// Test Ollama connection
	fmt.Printf("🔍 Testing connection to Ollama at %s...\n", ollamaHost)
	ctx := context.Background()
	testPrompt := "Hello"
	if _, err := ollama.GenerateResponse(testPrompt); err != nil {
		log.Printf("⚠️  Ollama connection failed: %v", err)
		log.Println("💡 Continuing with simulated responses...")
		ollama = nil // Will trigger fallback behavior
	} else {
		fmt.Printf("✅ Successfully connected to Ollama!\n")
		fmt.Printf("🤖 Using model: %s\n", ollamaModel)
	}
	fmt.Println()

	// Initialize React agents with Ollama support
	reasoner := react.NewReasoner("AI-Reasoner")
	actor := react.NewActor("AI-Actor", "Execute decisions based on AI reasoning", reasoner)
	observer := react.NewObserver("AI-Observer")

	// Create a coordinator
	coordinator := core.NewAgent("ReactCoordinator", "Coordinate React agents", nil, "gpt-4", "auto")

	// Create React workflow
	reactWorkflow := workflows.NewReactWorkflow(coordinator, reasoner, actor, observer)

	fmt.Println("🚀 Starting AI-enhanced React loop workflow...")
	fmt.Printf("Coordinator: %s\n", coordinator.Name)
	fmt.Printf("Reasoner: %s (AI-powered)\n", reasoner.GetName())
	fmt.Printf("Actor: %s (AI-guided)\n", actor.Name)
	fmt.Printf("Observer: %s (AI-monitoring)\n", observer.GetName())
	fmt.Println()

	// Demonstrate multiple React cycles with different scenarios
	scenarios := []string{
		"A user reports slow application performance during peak hours",
		"The system detected unusual data patterns in user behavior",
		"A new feature request requires integration with external APIs",
	}

	for i, scenario := range scenarios {
		fmt.Printf("🧪 Scenario %d: %s\n", i+1, scenario)
		fmt.Println(strings.Repeat("-", 60))

		// Enhanced execution with Ollama reasoning
		if ollama != nil {
			executeWithOllama(ctx, reactWorkflow, ollama, scenario)
		} else {
			reactWorkflow.Execute()
		}

		fmt.Println()
	}

	fmt.Println("🎉 AI-enhanced React loop workflow completed successfully!")
}

// executeWithOllama runs the React workflow with AI-enhanced reasoning using Ollama
func executeWithOllama(ctx context.Context, workflow *workflows.ReactWorkflow, ollama *providers.OllamaProvider, scenario string) {
	fmt.Println("🧠 AI Reasoning Phase...")

	// Enhanced reasoning with Ollama
	reasoningPrompt := fmt.Sprintf(`You are an expert system analyst. Analyze this situation and provide a structured response:

Situation: %s

Please provide:
1. Analysis of the situation
2. Potential root causes
3. Recommended actions
4. Risk assessment

Format your response as a clear, actionable decision.`, scenario)

	decision, err := ollama.GenerateResponse(reasoningPrompt)
	if err != nil {
		fmt.Printf("❌ AI reasoning failed: %v\n", err)
		decision = "Fallback: Investigate the situation manually and gather more data"
	} else {
		fmt.Printf("✅ AI Analysis completed\n")
		fmt.Printf("📋 Decision: %s\n", truncateString(decision, 200))
	}

	fmt.Println("\n⚡ Action Phase...")

	// Enhanced action planning with Ollama
	actionPrompt := fmt.Sprintf(`Based on this analysis and decision:

%s

Create a specific action plan with:
1. Immediate steps to take
2. Timeline for implementation
3. Success metrics
4. Monitoring approach

Provide a concise, actionable plan.`, decision)

	actionPlan, err := ollama.GenerateResponse(actionPrompt)
	if err != nil {
		fmt.Printf("❌ Action planning failed: %v\n", err)
		actionPlan = "Execute standard procedures for this type of situation"
	} else {
		fmt.Printf("✅ Action plan generated\n")
		fmt.Printf("🎯 Plan: %s\n", truncateString(actionPlan, 200))
	}

	fmt.Println("\n👁️  Observation Phase...")

	// Enhanced observation with Ollama
	observationPrompt := fmt.Sprintf(`You are monitoring the execution of this action plan:

%s

Provide an assessment of:
1. Implementation progress
2. Potential issues or risks
3. Effectiveness indicators
4. Recommendations for adjustments

Give a brief monitoring report.`, actionPlan)

	observation, err := ollama.GenerateResponse(observationPrompt)
	if err != nil {
		fmt.Printf("❌ Observation analysis failed: %v\n", err)
		observation = "Continue monitoring standard metrics and alert if anomalies detected"
	} else {
		fmt.Printf("✅ Monitoring assessment completed\n")
		fmt.Printf("📊 Observation: %s\n", truncateString(observation, 200))
	}

	fmt.Println("🔄 React cycle completed for this scenario")
}

// truncateString truncates a string to the specified length and adds "..." if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
