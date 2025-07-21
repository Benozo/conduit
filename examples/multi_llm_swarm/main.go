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
	fmt.Println("🚀 Multi-LLM Agent Swarm Demo")
	fmt.Println("=============================")
	fmt.Println()

	// Create MCP server
	config := conduit.DefaultConfig()
	config.EnableLogging = false
	server := conduit.NewEnhancedServer(config)

	// Register tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Create swarm client (no default LLM)
	swarmClient := swarm.NewSwarmClient(server, nil)

	// Create agents with different LLM providers
	fmt.Println("🤖 Creating agents with different LLM providers:")

	// Coordinator with Ollama llama3.2
	coordinator := swarmClient.CreateAgentWithModel("coordinator",
		"Route tasks to appropriate agents based on request type",
		[]string{},
		&conduit.ModelConfig{
			Provider:    "ollama",
			Model:       "llama3.2",
			URL:         "http://192.168.10.10:11434",
			Temperature: 0.7,
			MaxTokens:   1000,
		})
	fmt.Printf("   📋 %s - Ollama llama3.2 (fast routing)\n", coordinator.Name)

	// Content creator with Ollama qwen2.5 (better for content)
	contentCreator := swarmClient.CreateAgentWithModel("content_creator",
		"Handle content creation and text processing tasks",
		[]string{"uppercase", "lowercase", "snake_case", "camel_case", "word_count"},
		&conduit.ModelConfig{
			Provider:    "ollama",
			Model:       "qwen2.5",
			URL:         "http://192.168.10.10:11434",
			Temperature: 0.5,
			MaxTokens:   1500,
		})
	fmt.Printf("   ✍️  %s - Ollama qwen2.5 (optimized for content)\n", contentCreator.Name)

	// Data analyst with OpenAI GPT-4 (premium reasoning)
	dataAnalyst := swarmClient.CreateAgentWithModel("data_analyst",
		"Perform complex data analysis and generate insights",
		[]string{"word_count", "timestamp", "json_format"},
		&conduit.ModelConfig{
			Provider:    "openai",
			Model:       "gpt-4o",
			APIKey:      os.Getenv("OPENAI_API_KEY"),
			Temperature: 0.3,
			MaxTokens:   2000,
		})
	fmt.Printf("   📊 %s - OpenAI GPT-4 (premium reasoning)\n", dataAnalyst.Name)

	// Code generator with DeepInfra Qwen Coder
	codeGenerator := swarmClient.CreateAgentWithModel("code_generator",
		"Generate and review code efficiently",
		[]string{"base64_encode", "base64_decode", "json_format", "json_minify"},
		&conduit.ModelConfig{
			Provider:    "deepinfra",
			Model:       "meta-llama/Meta-Llama-3.1-8B-Instruct",
			APIKey:      os.Getenv("DEEPINFRA_API_KEY"),
			Temperature: 0.1,
			MaxTokens:   2000,
		})
	fmt.Printf("   💻 %s - DeepInfra Qwen Coder (code specialist)\n", codeGenerator.Name)

	fmt.Println()

	// Demo scenarios showing different models working together
	scenarios := []struct {
		name        string
		message     string
		startAgent  *swarm.Agent
		description string
	}{
		{
			"Text Processing",
			"Convert 'Hello Multi-LLM World' to snake_case format",
			contentCreator,
			"ContentCreator (qwen2.5) handles text processing directly",
		},
		{
			"Data Analysis",
			"Analyze this text and count words: 'Multi-agent systems with different LLMs are powerful'",
			dataAnalyst,
			"DataAnalyst (GPT-4) performs complex analysis",
		},
		{
			"Code Generation",
			"Generate a JSON object with user data and encode it in base64 format",
			codeGenerator,
			"CodeGenerator (Meta-Llama) creates optimized code",
		},
	}

	fmt.Println("🎯 Multi-LLM Demo Scenarios:")
	for i, scenario := range scenarios {
		fmt.Printf("\n📝 Scenario %d: %s\n", i+1, scenario.name)
		fmt.Printf("📄 %s\n", scenario.description)
		fmt.Printf("💬 Request: %s\n", scenario.message)

		// Check if models are available (in real usage)
		if !checkModelAvailability(scenario.startAgent) {
			fmt.Printf("⚠️  Skipping - %s model not available\n", getModelInfo(scenario.startAgent))
			continue
		}

		fmt.Printf("🔄 Processing with %s...\n", getModelInfo(scenario.startAgent)) // Actually call the LLM if available
		if scenario.startAgent.ModelFunc != nil {
			// Create proper MCP request structure with the actual message
			ctx := mcp.ContextInput{
				ContextID: "demo",
				Inputs:    map[string]interface{}{"query": scenario.message},
			}
			req := mcp.MCPRequest{
				SessionID:   "demo-session",
				Contexts:    []mcp.ContextInput{ctx},
				Model:       scenario.startAgent.ModelConfig.Model,
				Temperature: scenario.startAgent.ModelConfig.Temperature,
			}

			response, err := scenario.startAgent.ModelFunc(ctx, req, nil, nil)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("🤖 Response: %s\n", truncateResponse(response))
			}
		} else {
			fmt.Printf("✅ Would route to appropriate specialist agent (demo mode)\n")
		}
	}

	fmt.Println("\n🧪 Testing Per-Agent Model Selection:")

	// Test that each agent has its own model configuration
	agents := []*swarm.Agent{coordinator, contentCreator, dataAnalyst, codeGenerator}
	for _, agent := range agents {
		if agent.ModelConfig != nil {
			fmt.Printf("   🤖 %s: %s %s (temp: %.1f)\n",
				agent.Name,
				agent.ModelConfig.Provider,
				agent.ModelConfig.Model,
				agent.ModelConfig.Temperature)
		}
	}

	fmt.Println("✨ Multi-LLM Features Demonstrated:")
	fmt.Println("   🎯 Task-specific model selection")
	fmt.Println("   💰 Cost optimization (local vs cloud models)")
	fmt.Println("   ⚡ Performance optimization (speed vs quality)")
	fmt.Println("   🔄 Intelligent routing between different providers")
	fmt.Println("   🛡️  Provider redundancy and fallback")

	fmt.Println("\n🚀 To run with real models:")
	fmt.Println("   1. Start Ollama: ollama serve")
	fmt.Println("   2. Pull models: ollama pull llama3.2 && ollama pull qwen2.5")
	fmt.Println("   3. Set API keys: export OPENAI_API_KEY=... DEEPINFRA_API_KEY=...")
	fmt.Println("   4. Run: go run examples/multi_llm_swarm/main.go")
}

func checkModelAvailability(agent *swarm.Agent) bool {
	// In a real implementation, this would check if the model is available
	// For demo purposes, we'll check environment variables
	if agent.ModelConfig == nil {
		return false
	}

	switch agent.ModelConfig.Provider {
	case "ollama":
		return true // Assume Ollama is available for demo
	case "openai":
		return os.Getenv("OPENAI_API_KEY") != ""
	case "deepinfra":
		return os.Getenv("DEEPINFRA_API_KEY") != ""
	default:
		return false
	}
}

func getModelInfo(agent *swarm.Agent) string {
	if agent.ModelConfig == nil {
		return "unknown model"
	}
	return fmt.Sprintf("%s %s", agent.ModelConfig.Provider, agent.ModelConfig.Model)
}

func truncateResponse(response string) string {
	// Clean up the response and truncate if too long
	cleaned := strings.TrimSpace(response)
	if len(cleaned) > 100 {
		return cleaned[:97] + "..."
	}
	return cleaned
}
