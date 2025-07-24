package main

import (
	"fmt"
	"log"

	"github.com/benozo/neuron/src/core"
	"github.com/benozo/neuron/src/llm/providers"
)

// OllamaAgent represents an agent that uses Ollama for LLM capabilities
type OllamaAgent struct {
	*core.BaseAgent
	ollamaProvider *providers.OllamaProvider
}

// NewOllamaAgent creates a new agent with Ollama integration
func NewOllamaAgent(name, instructions string, ollamaURL, model string) *OllamaAgent {
	// Create the Ollama provider
	ollamaProvider := providers.NewOllamaProvider(ollamaURL, model)

	// Create the base agent
	baseAgent := core.NewAgent(name, instructions, nil, model, "auto")

	return &OllamaAgent{
		BaseAgent:      baseAgent,
		ollamaProvider: ollamaProvider,
	}
}

// ProcessPrompt sends a prompt to Ollama and returns the response
func (oa *OllamaAgent) ProcessPrompt(prompt string) (string, error) {
	fmt.Printf("🤖 %s is processing: %s\n", oa.GetName(), prompt)

	response, err := oa.ollamaProvider.GenerateResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to get response from Ollama: %w", err)
	}

	fmt.Printf("✅ %s generated response (%d characters)\n", oa.GetName(), len(response))
	return response, nil
}

// Execute implements the Agent interface with Ollama integration
func (oa *OllamaAgent) Execute(input interface{}) (interface{}, error) {
	prompt, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string input, got %T", input)
	}

	return oa.ProcessPrompt(prompt)
}

// GetModelInfo returns information about the underlying model
func (oa *OllamaAgent) GetModelInfo() string {
	info := oa.ollamaProvider.GetModelInfo()
	return fmt.Sprintf("Model: %s | Provider: %s | Description: %s",
		info.Name, info.Provider, info.Description)
}

func main() {
	fmt.Println("=== Ollama Agent Demo ===")

	// Configuration
	ollamaURL := "http://192.168.10.10:11434"
	model := "llama3.2"

	// Create the Ollama agent
	agent := NewOllamaAgent(
		"OllamaAssistant",
		"You are a helpful AI assistant powered by Ollama running llama3.2 model.",
		ollamaURL,
		model,
	)

	fmt.Printf("🚀 Created agent: %s\n", agent.GetName())
	fmt.Printf("📋 Instructions: %s\n", agent.Instructions)
	fmt.Printf("🔧 %s\n\n", agent.GetModelInfo())

	// Test connection to Ollama
	fmt.Printf("🔍 Testing connection to Ollama at %s...\n", ollamaURL)
	if err := agent.ollamaProvider.Ping(); err != nil {
		log.Printf("❌ Failed to connect to Ollama: %v", err)
		log.Println("Make sure Ollama is running at", ollamaURL)
		log.Println("You can start it with: ollama serve")
		log.Printf("And pull the model with: ollama pull %s\n", model)
		return
	}
	fmt.Println("✅ Successfully connected to Ollama!")

	// Test prompts
	prompts := []string{
		"What is artificial intelligence?",
		"Explain the concept of machine learning in simple terms.",
		"Write a short poem about robots.",
		"What are the main benefits of using large language models?",
	}

	fmt.Println("\n🧪 Testing agent with various prompts...")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	for i, prompt := range prompts {
		fmt.Printf("\n📝 Test %d: %s\n", i+1, prompt)
		fmt.Println("-" + string(make([]byte, 60)) + "-")

		response, err := agent.ProcessPrompt(prompt)
		if err != nil {
			log.Printf("❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("🤖 Response:\n%s\n", response)
	}

	fmt.Println("\n🎉 Ollama Agent demo completed successfully!")
	fmt.Println("The agent successfully communicated with Ollama and processed all prompts.")
}
