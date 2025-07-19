package main

import (
	"log"
	"os"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Create server configuration
	config := conduit.DefaultConfig()

	// Parse command line arguments to set mode
	mode := mcp.ModeBoth // Default to both protocols
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--stdio":
			mode = mcp.ModeStdio
		case "--http":
			mode = mcp.ModeHTTP
		case "--both":
			mode = mcp.ModeBoth
		case "--agents":
			// Demo AI Agents functionality
			demoAgents()
			return
		}
	}
	config.Mode = mode

	// Create the server
	// server := conduit.NewServer(config)
	server := conduit.NewEnhancedServer(config)

	// Register all tool packages
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Use Custom tools
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			return map[string]interface{}{
				"result":    a + b,
				"operation": "addition",
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("add", "Add two numbers together", map[string]interface{}{
			"a": conduit.NumberParam("First number to add"),
			"b": conduit.NumberParam("Second number to add"),
		}, []string{"a", "b"}))

	// Register another custom tool for curl

	server.RegisterToolWithSchema("multiply",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			return map[string]interface{}{
				"result":    a * b,
				"operation": "multiplication",
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number to multiply"),
			"b": conduit.NumberParam("Second number to multiply"),
		}, []string{"a", "b"}))

	// Start the server
	log.Printf("Starting Conduit server (mode: %v)...", mode)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

// demoAgents demonstrates the AI Agents functionality
func demoAgents() {
	log.Println("🤖 AI Agents Demo Mode")
	log.Println("======================")

	// Create MCP server for agents
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8080
	config.EnableLogging = true

	server := conduit.NewEnhancedServer(config)

	// Register tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Add math tools
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a + b
			log.Printf("🧮 Math: %.1f + %.1f = %.1f", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "addition",
			}, nil
		},
		conduit.CreateToolMetadata("add", "Add two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number"),
			"b": conduit.NumberParam("Second number"),
		}, []string{"a", "b"}))

	// Create agent manager
	agentManager := agents.NewMCPAgentManager(server)

	// Create specialized agents
	log.Println("📝 Creating AI agents...")
	if err := agentManager.CreateSpecializedAgents(); err != nil {
		log.Fatalf("Failed to create agents: %v", err)
	}

	// Start server
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait and create a sample task
	log.Println("⏳ Starting server...")
	// In a real scenario, you would wait for proper startup
	// For demo purposes, we'll show the concept

	log.Println("✅ AI Agents ready!")
	log.Println("🔗 Server running on http://localhost:8080")
	log.Println("📚 See examples/ai_agents/ for complete usage examples")
	log.Println("📖 See agents/README.md for full documentation")

	// Keep running
	select {}
}
