package main

import (
	"log"
	"os"

	"github.com/benozo/conduit/mcp"
)

func main() {
	// Pure MCP usage - no wrapper library needed

	// Create tools registry
	tools := mcp.NewToolRegistry()

	// Register some basic tools directly
	tools.Register("uppercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		return map[string]string{"result": text}, nil
	})

	tools.Register("echo", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		return map[string]string{"echo": text}, nil
	})

	// Simple model function
	model := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := ctx.Inputs["query"].(string)
		return "Response to: " + query, nil
	}

	// Create and start server
	server := mcp.NewUnifiedServer(model, tools)

	// Set mode from command line or default to both
	mode := mcp.ModeBoth
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--stdio":
			mode = mcp.ModeStdio
		case "--http":
			mode = mcp.ModeHTTP
		}
	}
	server.SetMode(mode)

	log.Printf("Starting pure MCP server (mode: %v)...", mode)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
