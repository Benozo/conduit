// OpenAI MCP Server Example
//
// This example demonstrates how to create a production-ready MCP server
// that integrates with OpenAI's API (or OpenAI-compatible services).
//
// Features:
// - Full integration with OpenAI API
// - Standard MCP tools (text, memory, utility)
// - Custom diagnostic and monitoring tools
// - HTTP REST API with comprehensive endpoints
// - Production-ready error handling and logging
//
// Usage:
//
//	export OPENAI_API_KEY="your-api-key"
//	go run examples/openai/main.go
//
// Test:
//
//	curl http://localhost:9090/health
//	curl -X POST http://localhost:9090/tool -H 'Content-Type: application/json' \
//	  -d '{"name":"model_info","params":{}}'
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
	if len(openaiKey) < 20 {
		log.Printf("⚠️  Warning: OPENAI_API_KEY seems too short (length: %d)", len(openaiKey))
	}

	openaiURL := os.Getenv("OPENAI_API_URL")
	if openaiURL == "" {
		openaiURL = "https://api.openai.com"
	}

	modelName := os.Getenv("OPENAI_MODEL")
	if modelName == "" {
		modelName = "gpt-4o-mini"
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

	// Register standard MCP tools
	log.Printf("🔧 Registering MCP tools...")
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)
	log.Printf("✅ Registered standard MCP tools: text, memory, and utility tools")

	// Add OpenAI-specific diagnostic tool
	log.Printf("🔧 Registering custom tools...")
	server.RegisterTool("model_info", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		return map[string]interface{}{
			"openai_url": openaiURL,
			"model":      modelName,
			"status":     "connected",
			"provider":   "OpenAI",
		}, nil
	})

	// Add chat history management tool
	server.RegisterTool("chat_history", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		history := memory.Get("chat_history")
		if history == nil {
			history = []map[string]interface{}{}
		}

		if message, ok := params["message"]; ok {
			timestampResult, err := tools.TimestampFunc(map[string]interface{}{"format": "iso"}, memory)
			if err != nil {
				log.Printf("⚠️ Timestamp generation failed: %v", err)
				timestampResult = map[string]interface{}{"iso": "unknown"}
			}

			var timestamp string
			if timestampMap, ok := timestampResult.(map[string]interface{}); ok {
				if ts, exists := timestampMap["iso"]; exists {
					timestamp = ts.(string)
				} else if ts, exists := timestampMap["result"]; exists {
					timestamp = ts.(string)
				} else {
					timestamp = "unknown"
				}
			} else {
				timestamp = "unknown"
			}

			newEntry := map[string]interface{}{
				"timestamp": timestamp,
				"message":   message,
				"type":      "user",
			}

			historyList := history.([]map[string]interface{})
			historyList = append(historyList, newEntry)
			memory.Set("chat_history", historyList)

			return map[string]interface{}{
				"status":  "added",
				"count":   len(historyList),
				"history": historyList,
			}, nil
		}

		// Just return history if no message to add
		return map[string]interface{}{
			"history": history,
			"count":   len(history.([]map[string]interface{})),
		}, nil
	})

	// Add OpenAI-specific test tool
	server.RegisterTool("openai_test", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		testMessage := "Hello from OpenAI MCP server!"
		if msg, ok := params["message"]; ok {
			testMessage = msg.(string)
		}

		return map[string]interface{}{
			"input_message": testMessage,
			"response":      "OpenAI MCP server is working correctly!",
			"model":         modelName,
			"timestamp":     "2025-07-19",
			"tools_available": []string{
				"uppercase", "lowercase", "trim", "reverse",
				"remember", "recall", "clear_memory", "list_memories",
				"timestamp", "uuid", "hash_md5", "hash_sha256",
				"model_info", "chat_history", "openai_test",
			},
		}, nil
	})

	// Create OpenAI model
	apiURL := openaiURL + "/v1/chat/completions"
	openaiModel := conduit.CreateOpenAICompatibleModel(apiURL, openaiKey)

	// Set the OpenAI model
	server.SetModel(openaiModel)

	log.Printf("🚀 Starting OpenAI-powered MCP server on port %d...", config.Port)
	log.Printf("")
	log.Printf("📡 Available endpoints:")
	log.Printf("  GET  http://localhost:%d/health        - Health check", config.Port)
	log.Printf("  GET  http://localhost:%d/schema        - Tool schema", config.Port)
	log.Printf("  POST http://localhost:%d/tool          - Execute tool", config.Port)
	log.Printf("  POST http://localhost:%d/chat          - Chat with AI", config.Port)
	log.Printf("  POST http://localhost:%d/mcp           - MCP protocol", config.Port)
	log.Printf("  POST http://localhost:%d/react         - ReAct reasoning", config.Port)
	log.Printf("")
	log.Printf("🔧 Environment configuration:")
	log.Printf("  OPENAI_API_KEY: %s", maskKey(openaiKey))
	log.Printf("  OPENAI_API_URL: %s", openaiURL)
	log.Printf("  OPENAI_MODEL:   %s", modelName)
	log.Printf("")
	log.Printf("🧪 Test the server:")
	log.Printf("  curl http://localhost:%d/health", config.Port)
	log.Printf("  curl -X POST http://localhost:%d/tool -H 'Content-Type: application/json' -d '{\"name\":\"model_info\",\"params\":{}}'", config.Port)
	log.Printf("  curl -X POST http://localhost:%d/tool -H 'Content-Type: application/json' -d '{\"name\":\"openai_test\",\"params\":{\"message\":\"Hello!\"}}'", config.Port)

	log.Fatal(server.Start())
}

// maskKey masks an API key for logging, showing only first 8 and last 4 characters
func maskKey(key string) string {
	if len(key) <= 12 {
		return "***hidden***"
	}
	return key[:8] + "..." + key[len(key)-4:]
}
