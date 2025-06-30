package main

import (
	"log"
	"os"

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
