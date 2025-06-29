package main

import (
	"log"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Example: Using Conduit as an embedded library

	// Create custom configuration
	config := conduit.DefaultConfig()
	config.Port = 8081
	config.Mode = mcp.ModeHTTP
	config.EnableLogging = true

	// Create the server
	server := conduit.NewServer(config)

	// Register standard tool packages
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Register custom tools
	server.RegisterTool("custom_greeting", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		name := "World"
		if nameParam, ok := params["name"]; ok {
			name = nameParam.(string)
		}
		return map[string]string{
			"greeting": "Hello, " + name + "!",
		}, nil
	})

	// Set a custom model (optional)
	server.SetModel(func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		// Simple echo model for testing
		query := ctx.Inputs["query"].(string)
		return "Echo: " + query, nil
	})

	// Start the server
	log.Printf("Starting embedded conduit server on port %d...", config.Port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
