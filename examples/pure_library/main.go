package main

import (
	"fmt"
	"log"

	"github.com/benozo/conduit/mcp"
)

func main() {
	// Pure library usage - no server, just MCP components

	// Create memory for state management
	memory := mcp.NewMemory()

	// Create tools registry for functionality
	tools := mcp.NewToolRegistry()

	// Register tools as library functions
	tools.Register("uppercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		return map[string]string{"result": text}, nil
	})

	tools.Register("remember", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := params["key"].(string)
		value := params["value"]
		memory.Set(key, value)
		return map[string]string{"status": "remembered " + key}, nil
	})

	tools.Register("recall", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := params["key"].(string)
		value := memory.Get(key)
		if value == nil {
			return map[string]string{"result": "not found"}, nil
		}
		return map[string]interface{}{"result": value}, nil
	})

	// Use the tools directly as library functions
	fmt.Println("=== Pure MCP Library Usage ===")

	// Example 1: Use uppercase tool
	result, err := tools.Call("uppercase", map[string]interface{}{
		"text": "hello world",
	}, memory)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Uppercase result: %v\n", result)

	// Example 2: Use memory functions
	_, err = tools.Call("remember", map[string]interface{}{
		"key":   "user_name",
		"value": "John Doe",
	}, memory)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Stored user name in memory")

	// Example 3: Recall from memory
	recalled, err := tools.Call("recall", map[string]interface{}{
		"key": "user_name",
	}, memory)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recalled from memory: %v\n", recalled)

	// Example 4: Direct memory access (no tools)
	memory.Set("direct_key", "direct_value")
	directValue := memory.Get("direct_key")
	fmt.Printf("Direct memory access: %v\n", directValue)

	// Example 5: Create a simple model function for processing
	model := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])
		// Simple echo processing
		return "Processed: " + query, nil
	}

	// Use model directly through processor
	processor := mcp.NewProcessor(model, tools)

	// Create a request to process
	request := mcp.MCPRequest{
		SessionID: "example-session",
		Model:     "default",
		Contexts: []mcp.ContextInput{{
			ContextID: "example",
			Inputs:    map[string]interface{}{"query": "What is the weather like?"},
		}},
	}

	response, err := processor.Run(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Model response: %v\n", response)

	fmt.Println("\n=== Library components used successfully ===")
	fmt.Println("Users can integrate these components into their own servers/applications")
}
