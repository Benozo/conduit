package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/benozo/conduit/mcp"
)

// This file demonstrates direct usage of the MCP package for ReAct
// Run with: go run main.go

func main() {
	fmt.Println("=== Direct MCP Package ReAct Example ===")
	fmt.Println()

	// Create memory and tool registry
	memory := mcp.NewMemory()
	toolRegistry := mcp.NewToolRegistry()

	// Register some basic tools
	toolRegistry.Register("uppercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		result := strings.ToUpper(text)
		fmt.Printf("Tool 'uppercase' called with: %s -> %s\n", text, result)
		return result, nil
	})

	toolRegistry.Register("lowercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		result := strings.ToLower(text)
		fmt.Printf("Tool 'lowercase' called with: %s -> %s\n", text, result)
		return result, nil
	})

	toolRegistry.Register("word_count", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		count := len(strings.Fields(text))
		result := map[string]interface{}{
			"text":       text,
			"word_count": count,
		}
		fmt.Printf("Tool 'word_count' called with: %s -> %d words\n", text, count)
		return result, nil
	})

	// Set up initial data
	inputText := "Hello World from conduit"
	memory.Set("latest", inputText)
	fmt.Printf("Initial input: %s\n", inputText)
	fmt.Println()

	// Define ReAct thoughts/actions
	thoughts := []string{
		"transform to uppercase",
		"count words in result",
		"store final result",
	}

	fmt.Println("=== ReAct Agent Processing ===")

	// Execute ReAct agent
	steps, err := mcp.ReActAgent(thoughts, toolRegistry, memory)
	if err != nil {
		log.Fatalf("ReAct agent error: %v", err)
	}

	// Display results
	for i, step := range steps {
		fmt.Printf("\nStep %d:\n", i+1)
		fmt.Printf("  Thought: %s\n", step.Thought)
		fmt.Printf("  Action: %s\n", step.Action)
		if len(step.Params) > 0 {
			fmt.Printf("  Params: %v\n", step.Params)
		}
		fmt.Printf("  Observation: %s\n", step.Observed)
	}

	// Show final memory state
	fmt.Println("\n=== Final Memory State ===")
	if latest := memory.Get("latest"); latest != nil {
		fmt.Printf("latest: %v\n", latest)
	}
	if result := memory.Get("latest_result"); result != nil {
		fmt.Printf("latest_result: %v\n", result)
	}

	fmt.Println("\n=== Direct MCP Usage Complete ===")
}
