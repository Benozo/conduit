package main

import (
	"log"
	"strings"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Example: Custom model integration patterns

	config := conduit.DefaultConfig()
	config.Port = 8083
	config.Mode = mcp.ModeHTTP

	server := conduit.NewServer(config)

	// Register standard tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)

	// Example 1: Echo model (for testing)
	echoModel := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := ctx.Inputs["query"].(string)
		response := "Echo: " + query

		// Simulate streaming
		if onToken != nil {
			words := strings.Split(response, " ")
			for _, word := range words {
				onToken(ctx.ContextID, word+" ")
			}
		}

		return response, nil
	}

	// Example 2: Mock AI model with predefined responses
	mockAIModel := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := strings.ToLower(ctx.Inputs["query"].(string))

		var response string
		switch {
		case strings.Contains(query, "hello"):
			response = "Hello! How can I help you today?"
		case strings.Contains(query, "weather"):
			response = "I'm sorry, I don't have access to real-time weather data, but you can check your local weather service."
		case strings.Contains(query, "time"):
			response = "I can help you with time-related queries. Try using the timestamp tool!"
		default:
			response = "I understand you're asking about: " + query + ". How can I assist you further?"
		}

		// Simulate streaming
		if onToken != nil {
			for i, char := range response {
				if i%5 == 0 { // Send every 5 characters
					onToken(ctx.ContextID, string(char))
				}
			}
		}

		return response, nil
	}

	// Example 3: Use Ollama model (requires Ollama running)
	ollamaModel := conduit.CreateOllamaModel("http://localhost:11434")

	// Switch between models based on configuration or runtime logic
	modelType := "mock" // Change to "echo", "ollama", etc.

	switch modelType {
	case "echo":
		server.SetModel(echoModel)
	case "mock":
		server.SetModel(mockAIModel)
	case "ollama":
		server.SetModel(ollamaModel)
	default:
		server.SetModel(echoModel) // Default fallback
	}

	log.Printf("Starting model integration server with %s model on port %d...", modelType, config.Port)
	log.Fatal(server.Start())
}
