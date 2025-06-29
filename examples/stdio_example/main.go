package main

import (
	"log"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	log.Println("=== Conduit Stdio MCP Example ===")
	log.Println("This server runs in stdio mode for MCP client integration")
	log.Println("Use this with VS Code Copilot, Cline, Claude Desktop, etc.")

	// Create configuration for stdio mode
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeStdio
	config.EnableLogging = true

	// Create server
	server := conduit.NewServer(config)

	// Register all available tool packages
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Register some custom tools specific to this example
	server.RegisterTool("stdio_demo", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		return map[string]interface{}{
			"result":      "Hello from stdio MCP server!",
			"mode":        "stdio",
			"description": "This tool demonstrates stdio MCP integration",
			"timestamp":   "2025-06-29",
		}, nil
	})

	server.RegisterTool("client_info", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		return map[string]interface{}{
			"result": "MCP Client Communication",
			"supported_clients": []string{
				"VS Code Copilot",
				"Cline",
				"Claude Desktop",
				"Continue.dev",
				"Cursor IDE",
				"Any MCP-compatible client",
			},
			"protocol": "stdio",
			"tools":    31,
		}, nil
	})

	log.Println("Starting stdio MCP server...")
	log.Println("Configure your MCP client with:")
	log.Println("  Command: /path/to/this/binary")
	log.Println("  Args: [\"--stdio\"]")
	log.Println("  Protocol: stdio")

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
