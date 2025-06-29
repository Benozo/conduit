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
	server := conduit.NewServer(config)

	// Register all tool packages
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Start the server
	log.Printf("Starting Conduit server (mode: %v)...", mode)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
