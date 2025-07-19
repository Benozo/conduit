package main

import (
	"log"
	"os"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Get OpenAI API key and model from environment
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("❌ OPENAI_API_KEY environment variable is required")
	}

	openaiURL := os.Getenv("OPENAI_API_URL")
	if openaiURL == "" {
		openaiURL = "https://api.openai.com"
	}

	modelName := os.Getenv("OPENAI_MODEL")
	if modelName == "" {
		modelName = "gpt-40-mini"
	}

	// Create configuration
	config := conduit.DefaultConfig()
	config.Port = 9090
	config.Mode = mcp.ModeHTTP
	config.EnableLogging = true

	log.Printf("🧠 Using OpenAI at: %s", openaiURL)
	log.Printf("📦 Using model: %s", modelName)

	// Create server
	server := conduit.NewServer(config)

	// Register tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Add OpenAI-specific diagnostic tool
	server.RegisterTool("model_info", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		return map[string]interface{}{
			"openai_url": openaiURL,
			"model":      modelName,
			"status":     "connected",
		}, nil
	})

	server.RegisterTool("chat_history", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		history := memory.Get("chat_history")
		if history == nil {
			history = []map[string]string{}
		}

		if message, ok := params["message"]; ok {
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

	// Create OpenAI model with tool support
	openaiModel := conduit.CreateOpenAIToolAwareModel(openaiKey, openaiURL, server.GetToolRegistry())

	// Set the OpenAI model
	server.SetModel(openaiModel)

	log.Printf("🚀 Starting OpenAI-powered MCP server on port %d...", config.Port)
	log.Printf("Try these endpoints:")
	log.Printf("  GET  http://localhost:%d/health", config.Port)
	log.Printf("  GET  http://localhost:%d/schema", config.Port)
	log.Printf("  POST http://localhost:%d/tool", config.Port)
	log.Printf("  POST http://localhost:%d/chat", config.Port)
	log.Printf("  POST http://localhost:%d/mcp", config.Port)
	log.Printf("  POST http://localhost:%d/react", config.Port)
	log.Printf("")
	log.Printf("Environment variables:")
	log.Printf("  OPENAI_API_KEY=%s", openaiKey)
	log.Printf("  OPENAI_API_URL=%s", openaiURL)
	log.Printf("  OPENAI_MODEL=%s", modelName)

	log.Fatal(server.Start())
}
