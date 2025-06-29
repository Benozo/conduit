package main

import (
	"log"
	"os"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Example: Ollama integration with conduit

	// Get Ollama URL from environment or use default
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	// Get model name from environment or use default
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3.2" // Default model
	}

	// Create configuration
	config := conduit.DefaultConfig()
	config.Port = 8084
	config.Mode = mcp.ModeHTTP
	config.OllamaURL = ollamaURL
	config.EnableLogging = true

	log.Printf("Using Ollama at: %s", ollamaURL)
	log.Printf("Using model: %s", modelName)

	// Create server
	server := conduit.NewServer(config)

	// Register comprehensive tool set
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Add some Ollama-specific tools
	server.RegisterTool("model_info", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		return map[string]interface{}{
			"ollama_url": ollamaURL,
			"model":      modelName,
			"status":     "connected",
		}, nil
	})

	server.RegisterTool("chat_history", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		// Simple chat history implementation
		history := memory.Get("chat_history")
		if history == nil {
			history = []map[string]string{}
		}

		if message, ok := params["message"]; ok {
			// Get timestamp
			timestampResult, _ := tools.TimestampFunc(map[string]interface{}{"format": "iso"}, memory)
			timestamp := timestampResult.(map[string]interface{})["timestamp"].(string)

			newEntry := map[string]string{
				"timestamp": timestamp,
				"message":   message.(string),
			}

			historyList := history.([]map[string]string)
			historyList = append(historyList, newEntry)
			memory.Set("chat_history", historyList)

			return map[string]interface{}{
				"status":  "added",
				"history": historyList,
			}, nil
		}

		return map[string]interface{}{
			"history": history,
		}, nil
	})

	// Create Ollama model with custom configuration
	ollamaModel := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		// Use the library's Ollama model with our configuration
		model := conduit.CreateOllamaModel(ollamaURL)

		// Enhance the context with model information
		enhancedCtx := ctx
		if enhancedCtx.Inputs == nil {
			enhancedCtx.Inputs = make(map[string]interface{})
		}

		// Add model name to the request if not specified
		enhancedReq := req
		if enhancedReq.Model == "" {
			enhancedReq.Model = modelName
		}

		// Store conversation in memory
		if query, ok := ctx.Inputs["query"]; ok {
			memory.Set("last_query", query)
		}

		return model(enhancedCtx, enhancedReq, memory, onToken)
	}

	// Set the Ollama model
	server.SetModel(ollamaModel)

	log.Printf("Starting Ollama-powered MCP server on port %d...", config.Port)
	log.Printf("Try these endpoints:")
	log.Printf("  GET  http://localhost:%d/health", config.Port)
	log.Printf("  GET  http://localhost:%d/schema", config.Port)
	log.Printf("  POST http://localhost:%d/mcp", config.Port)
	log.Printf("  POST http://localhost:%d/react", config.Port)
	log.Printf("")
	log.Printf("Environment variables:")
	log.Printf("  OLLAMA_URL=%s (Ollama server URL)", ollamaURL)
	log.Printf("  OLLAMA_MODEL=%s (Model to use)", modelName)

	log.Fatal(server.Start())
}
