// Pure Library Example
//
// This example demonstrates how to use MCP components as a pure library
// without any client-server infrastructure. This provides maximum performance
// and flexibility for embedding MCP functionality directly into applications.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/benozo/neuron-mcp/library"
	"github.com/benozo/neuron-mcp/protocol"
)

func main() {
	// Create a component registry for pure library usage
	registry := library.NewComponentRegistry()

	// Register tools directly
	if err := registerLibraryTools(registry); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// Demonstrate tool usage
	demonstrateTools(registry)

	// Demonstrate memory usage
	demonstrateMemory(registry)

	fmt.Println("Pure library example completed!")
}

// registerLibraryTools registers tools with the library registry
func registerLibraryTools(registry *library.ComponentRegistry) error {
	tools := registry.Tools()

	// Register text transformation tool
	if err := tools.Register("text_transform", handleTextTransform); err != nil {
		return fmt.Errorf("failed to register text_transform: %w", err)
	}

	// Register calculator tool with explicit schema
	calculatorSchema := &protocol.JSONSchema{
		Type: "object",
		Properties: map[string]*protocol.JSONSchema{
			"operation": {
				Type:        "string",
				Description: "Math operation: add, subtract, multiply, divide",
				Enum:        []interface{}{"add", "subtract", "multiply", "divide"},
			},
			"a": {
				Type:        "number",
				Description: "First number",
			},
			"b": {
				Type:        "number",
				Description: "Second number",
			},
		},
		Required: []string{"operation", "a", "b"},
	}

	if err := tools.RegisterWithSchema("calculator", handleCalculator, calculatorSchema); err != nil {
		return fmt.Errorf("failed to register calculator: %w", err)
	}

	// Register greeting tool
	if err := tools.Register("greet", handleGreeting); err != nil {
		return fmt.Errorf("failed to register greet: %w", err)
	}

	return nil
}

// demonstrateTools shows how to use tools in library mode
func demonstrateTools(registry *library.ComponentRegistry) {
	ctx := context.Background()
	tools := registry.Tools()

	fmt.Println("=== Tool Demonstration ===")

	// List available tools
	toolNames := tools.List()
	fmt.Printf("Available tools: %v\n", toolNames)

	// Test text transformation tool
	fmt.Println("\n1. Text Transformation:")
	result, err := tools.Call(ctx, "text_transform", map[string]interface{}{
		"text":      "Hello, World!",
		"operation": "uppercase",
	})
	if err != nil {
		log.Printf("Error calling text_transform: %v", err)
	} else {
		fmt.Printf("   Input: 'Hello, World!' -> Output: '%s'\n", result.Content[0].Text)
	}

	// Test calculator tool
	fmt.Println("\n2. Calculator:")
	result, err = tools.Call(ctx, "calculator", map[string]interface{}{
		"operation": "multiply",
		"a":         15.5,
		"b":         2.0,
	})
	if err != nil {
		log.Printf("Error calling calculator: %v", err)
	} else {
		fmt.Printf("   15.5 * 2.0 = %s\n", result.Content[0].Text)
	}

	// Test greeting tool
	fmt.Println("\n3. Greeting:")
	result, err = tools.Call(ctx, "greet", map[string]interface{}{
		"name": "Alice",
		"time": "morning",
	})
	if err != nil {
		log.Printf("Error calling greet: %v", err)
	} else {
		fmt.Printf("   %s\n", result.Content[0].Text)
	}

	// Test error handling
	fmt.Println("\n4. Error Handling:")
	_, err = tools.Call(ctx, "calculator", map[string]interface{}{
		"operation": "divide",
		"a":         10.0,
		"b":         0.0,
	})
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}
}

// demonstrateMemory shows how to use memory in library mode
func demonstrateMemory(registry *library.ComponentRegistry) {
	memory := registry.Memory()

	fmt.Println("\n=== Memory Demonstration ===")

	// Store some values
	memory.Set("user:1", map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
		"age":   30,
	})

	memory.Set("user:2", map[string]interface{}{
		"name":  "Bob",
		"email": "bob@example.com",
		"age":   25,
	})

	memory.Set("config:theme", "dark")

	// Retrieve values
	if user1, err := memory.Get("user:1"); err == nil {
		fmt.Printf("Retrieved user:1: %+v\n", user1)
	}

	// List all keys
	if keys, err := memory.List(); err == nil {
		fmt.Printf("All keys: %v\n", keys)
	}

	// Get memory statistics
	if stats, err := memory.Stats(); err == nil {
		fmt.Printf("Memory stats: %d active keys, backend: %s\n", stats.ActiveKeys, stats.Backend)
	}

	// Delete a key
	memory.Delete("config:theme")
	if keys, err := memory.List(); err == nil {
		fmt.Printf("Keys after deletion: %v\n", keys)
	}
}

// Tool implementations

// handleTextTransform implements text transformation
func handleTextTransform(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
	text, ok := params["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text parameter must be a string")
	}

	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter must be a string")
	}

	var result string
	switch operation {
	case "uppercase":
		result = strings.ToUpper(text)
	case "lowercase":
		result = strings.ToLower(text)
	case "reverse":
		result = reverseString(text)
	case "title":
		result = strings.Title(strings.ToLower(text))
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	return &protocol.ToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: result,
		}},
	}, nil
}

// handleCalculator implements mathematical operations
func handleCalculator(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter must be a string")
	}

	a, ok := params["a"].(float64)
	if !ok {
		return nil, fmt.Errorf("parameter 'a' must be a number")
	}

	b, ok := params["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("parameter 'b' must be a number")
	}

	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	return &protocol.ToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: fmt.Sprintf("%.2f", result),
		}},
	}, nil
}

// handleGreeting implements a greeting tool
func handleGreeting(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name parameter must be a string")
	}

	time, ok := params["time"].(string)
	if !ok {
		time = "day" // default
	}

	var greeting string
	switch time {
	case "morning":
		greeting = fmt.Sprintf("Good morning, %s!", name)
	case "afternoon":
		greeting = fmt.Sprintf("Good afternoon, %s!", name)
	case "evening":
		greeting = fmt.Sprintf("Good evening, %s!", name)
	default:
		greeting = fmt.Sprintf("Hello, %s!", name)
	}

	return &protocol.ToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: greeting,
		}},
	}, nil
}

// Helper functions

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
