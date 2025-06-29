package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/benozo/conduit/mcp"
)

// Simple CLI using MCP components as pure library
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: cli-tool <tool_name> <json_params>")
		fmt.Println("Example: cli-tool uppercase '{\"text\":\"hello world\"}'")
		os.Exit(1)
	}

	// Initialize MCP components
	memory := mcp.NewMemory()
	tools := mcp.NewToolRegistry()

	// Register tools
	tools.Register("uppercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		return map[string]string{"result": strings.ToUpper(text)}, nil
	})

	tools.Register("lowercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		return map[string]string{"result": strings.ToLower(text)}, nil
	})

	tools.Register("reverse", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		runes := []rune(text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return map[string]string{"result": string(runes)}, nil
	})

	tools.Register("store", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := params["key"].(string)
		value := params["value"]
		memory.Set(key, value)
		return map[string]string{"status": "stored " + key}, nil
	})

	tools.Register("recall", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := params["key"].(string)
		value := memory.Get(key)
		if value == nil {
			return map[string]string{"error": "not found"}, nil
		}
		return map[string]interface{}{"value": value}, nil
	})

	// Parse command line arguments
	toolName := os.Args[1]
	paramsJSON := os.Args[2]

	// Parse parameters
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		log.Fatalf("Invalid JSON parameters: %v", err)
	}

	// Execute tool using MCP
	result, err := tools.Call(toolName, params, memory)
	if err != nil {
		log.Fatalf("Tool execution failed: %v", err)
	}

	// Output result
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}
