package main

import (
	"log"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Example: Creating custom tools for specialized functionality

	config := conduit.DefaultConfig()
	config.Port = 8082
	config.Mode = mcp.ModeHTTP

	server := conduit.NewServer(config)

	// Math calculation tools
	server.RegisterTool("add", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		a := params["a"].(float64)
		b := params["b"].(float64)
		return map[string]float64{"result": a + b}, nil
	})

	server.RegisterTool("multiply", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		a := params["a"].(float64)
		b := params["b"].(float64)
		return map[string]float64{"result": a * b}, nil
	})

	// Data processing tools
	server.RegisterTool("filter_array", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		arr := params["array"].([]interface{})
		threshold := params["threshold"].(float64)

		var filtered []interface{}
		for _, item := range arr {
			if val, ok := item.(float64); ok && val > threshold {
				filtered = append(filtered, item)
			}
		}

		return map[string]interface{}{"filtered": filtered}, nil
	})

	// File system tools (example - be careful with security in production)
	server.RegisterTool("list_files", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		// This is just an example - implement proper security checks
		return map[string]interface{}{
			"files": []string{"example1.txt", "example2.txt"},
			"note":  "This is a mock implementation",
		}, nil
	})

	log.Printf("Starting custom tools server on port %d...", config.Port)
	log.Fatal(server.Start())
}
