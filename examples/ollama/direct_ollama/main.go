package main

import (
	"fmt"
	"log"
	"os"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/mcp"
)

// This file demonstrates direct usage of the Ollama model functions
// Run with: go run main.go

func main() {
	fmt.Println("=== Direct Ollama Integration Example ===")
	fmt.Println()

	// Get Ollama configuration
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3.2"
	}

	fmt.Printf("Ollama URL: %s\n", ollamaURL)
	fmt.Printf("Model: %s\n", modelName)
	fmt.Println()

	// Create memory for context
	memory := mcp.NewMemory()

	// Create Ollama model function directly
	ollamaModel := conduit.CreateOllamaModel(ollamaURL)

	// Test the model with a simple query
	ctx := mcp.ContextInput{
		ContextID: "direct-test",
		Inputs: map[string]interface{}{
			"query": "Hello! Please explain what you are in one sentence.",
		},
	}

	req := mcp.MCPRequest{
		Model: modelName,
	}

	fmt.Println("=== Testing Direct Ollama Call ===")
	fmt.Println("Query: Hello! Please explain what you are in one sentence.")
	fmt.Println("Response:")

	// Call the model directly
	response, err := ollamaModel(ctx, req, memory, func(contextID string, token string) {
		fmt.Print(token) // Stream tokens as they arrive
	})

	if err != nil {
		log.Printf("Error calling Ollama: %v", err)
		fmt.Println("\nNote: Make sure Ollama is running and the model is available.")
		fmt.Printf("Try: ollama serve (in one terminal) and ollama pull %s (in another)\n", modelName)
		return
	}

	fmt.Println("\n")
	fmt.Printf("Complete response length: %d characters\n", len(response))
	fmt.Println()

	// Test with a more complex query using memory
	memory.Set("previous_query", "Hello! Please explain what you are in one sentence.")
	memory.Set("user_name", "Developer")

	ctx2 := mcp.ContextInput{
		ContextID: "direct-test-2",
		Inputs: map[string]interface{}{
			"query": "Now tell me about Go programming language in 2 sentences.",
		},
	}

	fmt.Println("=== Testing with Memory Context ===")
	fmt.Println("Query: Now tell me about Go programming language in 2 sentences.")
	fmt.Println("Response:")

	response2, err := ollamaModel(ctx2, req, memory, func(contextID string, token string) {
		fmt.Print(token)
	})

	if err != nil {
		log.Printf("Error in second call: %v", err)
		return
	}

	fmt.Println("\n")
	fmt.Printf("Complete response length: %d characters\n", len(response2))

	// Show memory state
	fmt.Println("\n=== Memory State ===")
	if prev := memory.Get("previous_query"); prev != nil {
		fmt.Printf("previous_query: %v\n", prev)
	}
	if name := memory.Get("user_name"); name != nil {
		fmt.Printf("user_name: %v\n", name)
	}

	fmt.Println("\n=== Direct Ollama Usage Complete ===")
	fmt.Println("This demonstrates using Ollama models directly without a server.")
}
