package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Enhanced ReAct (Reasoning + Acting) Agent with Ollama Integration

	// Create configuration with Ollama support
	config := conduit.DefaultConfig()
	config.Port = 8085
	config.Mode = mcp.ModeHTTP
	config.EnableLogging = true

	// Configure Ollama integration
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://192.168.10.10:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.2"
	}

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
	// Create enhanced server for rich tool schemas
	server := conduit.NewEnhancedServer(config)

	// Register comprehensive tool set for the agent
	tools.RegisterTextTools(server.Server)
	tools.RegisterMemoryTools(server.Server)
	tools.RegisterUtilityTools(server.Server)

	// Add ReAct-specific tools with enhanced schemas
	server.RegisterToolWithSchema("analyze_sentiment",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			text := params["text"].(string)

			// Simple sentiment analysis (in real implementation, use ML model)
			positive_words := []string{"good", "great", "excellent", "amazing", "wonderful", "happy", "love"}
			negative_words := []string{"bad", "terrible", "awful", "hate", "sad", "angry", "disappointed"}

			pos_count := 0
			neg_count := 0

			text_lower := strings.ToLower(text)

			for _, word := range positive_words {
				if strings.Contains(text_lower, word) {
					pos_count++
				}
			}

			for _, word := range negative_words {
				if strings.Contains(text_lower, word) {
					neg_count++
				}
			}

			sentiment := "neutral"
			confidence := 0.5

			if pos_count > neg_count {
				sentiment = "positive"
				confidence = 0.7
			} else if neg_count > pos_count {
				sentiment = "negative"
				confidence = 0.7
			}

			return map[string]interface{}{
				"sentiment":           sentiment,
				"confidence":          confidence,
				"positive_indicators": pos_count,
				"negative_indicators": neg_count,
			}, nil
		},
		conduit.CreateToolMetadata("analyze_sentiment", "Analyze the sentiment of text", map[string]interface{}{
			"text": conduit.StringParam("Text to analyze for sentiment"),
		}, []string{"text"}),
	)

	server.RegisterToolWithSchema("math_calculate",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			// Handle both proper schema and text-only schema from Ollama
			if textParam, hasText := params["text"]; hasText {
				// Parse mathematical expression from text
				text := textParam.(string)
				return parseAndCalculate(text), nil
			}

			// Use proper parameters
			operation := params["operation"].(string)
			a := params["a"].(float64)
			b := params["b"].(float64)

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
					return map[string]interface{}{"error": "division by zero"}, nil
				}
				result = a / b
			default:
				return map[string]interface{}{"error": "unsupported operation"}, nil
			}

			return map[string]interface{}{
				"result":    result,
				"operation": operation,
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("math_calculate", "Perform mathematical calculations", map[string]interface{}{
			"operation": conduit.EnumParam("Mathematical operation", []string{"add", "subtract", "multiply", "divide"}),
			"a":         conduit.NumberParam("First number"),
			"b":         conduit.NumberParam("Second number"),
		}, []string{"operation", "a", "b"}),
	)

	server.RegisterToolWithSchema("web_search_mock",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			query := params["query"].(string)

			// Mock web search results (in real implementation, use actual search API)
			mockResults := []map[string]interface{}{
				{
					"title":   "Mock Result 1 for: " + query,
					"url":     "https://example.com/result1",
					"snippet": "This is a mock search result for demonstration purposes.",
				},
				{
					"title":   "Mock Result 2 for: " + query,
					"url":     "https://example.com/result2",
					"snippet": "Another mock result showing how ReAct can use web search.",
				},
			}

			return map[string]interface{}{
				"query":   query,
				"results": mockResults,
				"count":   len(mockResults),
			}, nil
		},
		conduit.CreateToolMetadata("web_search_mock", "Search the web (mock implementation)", map[string]interface{}{
			"query": conduit.StringParam("Search query"),
		}, []string{"query"}),
	)

	server.RegisterToolWithSchema("decision_maker",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			optionsParam, ok := params["options"]
			if !ok || optionsParam == nil {
				return map[string]interface{}{"error": "options parameter is required"}, nil
			}
			options, ok := optionsParam.([]any)
			if !ok {
				return map[string]interface{}{"error": "options must be an array"}, nil
			}

			criteriaParam, ok := params["criteria"]
			if !ok || criteriaParam == nil {
				return map[string]interface{}{"error": "criteria parameter is required"}, nil
			}
			criteria, ok := criteriaParam.(string)
			if !ok {
				return map[string]interface{}{"error": "criteria must be a string"}, nil
			}

			// Simple decision making logic
			decision := "Unable to decide"
			reasoning := "No clear criteria provided"

			if len(options) > 0 {
				// For demo, just pick the first option
				decision = options[0].(string)
				reasoning = "Selected based on criteria: " + criteria
			}

			return map[string]interface{}{
				"decision":  decision,
				"reasoning": reasoning,
				"options":   options,
				"criteria":  criteria,
			}, nil
		},
		conduit.CreateToolMetadata("decision_maker", "Make decisions based on options and criteria", map[string]interface{}{
			"options":  conduit.ArrayParam("Available options", "string"),
			"criteria": conduit.StringParam("Decision criteria"),
		}, []string{"options", "criteria"}),
	)

	// Set up intelligent model with fallback options
	// Try Ollama first, fallback to smart mock if unavailable
	// ollamaModel := conduit.CreateOllamaToolAwareModel(ollamaURL, server.Server.GetToolRegistry())
	openaiModel := conduit.CreateOpenAIToolAwareModel(openaiKey, openaiURL, server.GetToolRegistry())

	// Create a ReAct-enhanced model that combines Ollama with ReAct patterns
	reactModel := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := fmt.Sprintf("%v", ctx.Inputs["query"])
		log.Printf("🧠 ReAct Agent processing: %s", query)

		// Try Ollama or openai first for sophisticated reasoning
		// result, err := ollamaModel(ctx, req, memory, onToken)
		result, err := openaiModel(ctx, req, memory, onToken)

		if err == nil && result != "" {
			log.Printf("✅ Ollama ReAct processing successful")
			return result, nil
		}

		log.Printf("⚠️ Ollama unavailable, using ReAct fallback reasoning")

		// Fallback: Manual ReAct pattern implementation
		return processWithReActPattern(query, server.Server.GetToolRegistry(), memory, onToken, ctx.ContextID)
	}

	server.SetModel(reactModel)

	log.Printf("Starting Enhanced ReAct Agent with Ollama Integration on port %d...", config.Port)
	log.Printf("")
	log.Printf("🧠 ReAct Pattern: Reasoning + Acting with LLM Intelligence")
	log.Printf("🤖 Ollama Model: %s at %s", model, ollamaURL)
	log.Printf("🔧 Available tools for intelligent agent:")
	log.Printf("  • 31+ Conduit Tools: Text, Memory, Utility operations")
	log.Printf("  • Enhanced Analysis: analyze_sentiment")
	log.Printf("  • Mathematics: math_calculate (add, subtract, multiply, divide)")
	log.Printf("  • Web Search: web_search_mock")
	log.Printf("  • Decision Making: decision_maker")
	log.Printf("")
	log.Printf("🌐 HTTP Endpoints:")
	log.Printf("  GET  http://localhost:%d/health", config.Port)
	log.Printf("  GET  http://localhost:%d/schema", config.Port)
	log.Printf("  POST http://localhost:%d/tool", config.Port)
	log.Printf("  POST http://localhost:%d/chat", config.Port)
	log.Printf("  POST http://localhost:%d/react", config.Port)
	log.Printf("")
	log.Printf("💬 Natural Language ReAct Examples:")
	log.Printf(`  # Sentiment analysis with reasoning`)
	log.Printf(`  curl -X POST http://localhost:%d/chat \`, config.Port)
	log.Printf(`    -H "Content-Type: application/json" \`)
	log.Printf(`    -d '{"message": "analyze the sentiment of this text: I love this product, it works amazingly well!"}'`)
	log.Printf("")
	log.Printf(`  # Mathematical reasoning`)
	log.Printf(`  curl -X POST http://localhost:%d/chat \`, config.Port)
	log.Printf(`    -H "Content-Type: application/json" \`)
	log.Printf(`    -d '{"message": "calculate 25 multiplied by 4, then add 100 to the result"}'`)
	log.Printf("")
	log.Printf(`  # Multi-step reasoning with memory`)
	log.Printf(`  curl -X POST http://localhost:%d/chat \`, config.Port)
	log.Printf(`    -H "Content-Type: application/json" \`)
	log.Printf(`    -d '{"message": "remember that my favorite number is 42, then convert it to uppercase text"}'`)

	// Use stdio mode if requested (for MCP clients)
	if len(os.Args) > 1 && os.Args[1] == "--stdio" {
		config.Mode = mcp.ModeStdio
		log.Printf("\n🔄 Switching to stdio mode for VS Code Copilot integration...")
		// For stdio mode, we don't need the Ollama model, use a simpler approach
		simpleServer := conduit.NewEnhancedServer(config)
		tools.RegisterTextTools(simpleServer.Server)
		tools.RegisterMemoryTools(simpleServer.Server)
		tools.RegisterUtilityTools(simpleServer.Server)
		log.Fatal(simpleServer.Start())
	}

	log.Fatal(server.Start())
}

// processWithReActPattern implements manual ReAct reasoning when Ollama is unavailable
func processWithReActPattern(query string, tools *mcp.ToolRegistry, memory *mcp.Memory, onToken mcp.StreamCallback, contextID string) (string, error) {
	var response strings.Builder
	response.WriteString("🧠 ReAct Agent Reasoning:\n\n")

	// Step 1: Thought - Analyze the query
	thought := analyzeQuery(query)
	response.WriteString(fmt.Sprintf("💭 **Thought**: %s\n\n", thought))

	// Step 2: Action - Determine what tools to use
	actions := determineActions(query, tools.GetRegisteredTools())
	response.WriteString("🔧 **Planned Actions**:\n")
	for i, action := range actions {
		response.WriteString(fmt.Sprintf("  %d. %s\n", i+1, action))
	}
	response.WriteString("\n")

	// Step 3: Act - Execute the actions
	response.WriteString("⚡ **Execution**:\n")
	for i, action := range actions {
		result, err := executeAction(action, query, tools, memory)
		if err != nil {
			response.WriteString(fmt.Sprintf("  %d. ❌ %s: Error - %v\n", i+1, action, err))
		} else {
			response.WriteString(fmt.Sprintf("  %d. ✅ %s: %v\n", i+1, action, result))
		}
	}

	// Step 4: Observe - Summarize results
	response.WriteString("\n📋 **Final Observation**: ")
	observation := generateObservation(query, actions, memory)
	response.WriteString(observation)

	// Simulate streaming if callback provided
	responseText := response.String()
	if onToken != nil {
		words := strings.Split(responseText, " ")
		for _, word := range words {
			onToken(contextID, word+" ")
		}
	}

	return responseText, nil
}

// analyzeQuery provides reasoning about what the query is asking for
func analyzeQuery(query string) string {
	query = strings.ToLower(query)

	if strings.Contains(query, "sentiment") || strings.Contains(query, "feel") || strings.Contains(query, "emotion") {
		return "The user wants sentiment analysis. I should use analyze_sentiment tool."
	}
	if strings.Contains(query, "calculate") || strings.Contains(query, "math") || strings.Contains(query, "add") || strings.Contains(query, "multiply") {
		return "The user wants mathematical computation. I should use math_calculate tool."
	}
	if strings.Contains(query, "remember") || strings.Contains(query, "store") || strings.Contains(query, "save") {
		return "The user wants to store information in memory. I should use remember tool."
	}
	if strings.Contains(query, "recall") || strings.Contains(query, "retrieve") || strings.Contains(query, "get") {
		return "The user wants to retrieve stored information. I should use recall tool."
	}
	if strings.Contains(query, "uppercase") || strings.Contains(query, "capital") || strings.Contains(query, "upper") {
		return "The user wants text transformation to uppercase. I should use uppercase tool."
	}
	if strings.Contains(query, "search") || strings.Contains(query, "find") || strings.Contains(query, "look") {
		return "The user wants to search for information. I should use web_search_mock tool."
	}
	if strings.Contains(query, "decide") || strings.Contains(query, "choose") || strings.Contains(query, "pick") {
		return "The user wants decision making help. I should use decision_maker tool."
	}

	return "I should analyze this query and determine the most appropriate tools to use based on the user's request."
}

// determineActions selects which tools to use based on the query
func determineActions(query string, availableTools []string) []string {
	query = strings.ToLower(query)
	actions := []string{}

	// Multi-step reasoning: look for multiple intents
	if strings.Contains(query, "sentiment") || strings.Contains(query, "analyze") {
		if contains(availableTools, "analyze_sentiment") {
			actions = append(actions, "analyze_sentiment")
		}
	}

	if strings.Contains(query, "remember") || strings.Contains(query, "store") {
		if contains(availableTools, "remember") {
			actions = append(actions, "remember")
		}
	}

	if strings.Contains(query, "calculate") || strings.Contains(query, "math") {
		if contains(availableTools, "math_calculate") {
			actions = append(actions, "math_calculate")
		}
	}

	if strings.Contains(query, "uppercase") || strings.Contains(query, "upper") {
		if contains(availableTools, "uppercase") {
			actions = append(actions, "uppercase")
		}
	}

	if strings.Contains(query, "search") {
		if contains(availableTools, "web_search_mock") {
			actions = append(actions, "web_search_mock")
		}
	}

	// If no specific actions determined, use general text processing
	if len(actions) == 0 {
		if contains(availableTools, "uppercase") {
			actions = append(actions, "uppercase")
		}
	}

	return actions
}

// executeAction runs a specific tool action
func executeAction(action, query string, tools *mcp.ToolRegistry, memory *mcp.Memory) (interface{}, error) {
	params := make(map[string]interface{})

	switch action {
	case "analyze_sentiment":
		// Extract text for sentiment analysis
		text := extractTextFromQuery(query, "sentiment")
		params["text"] = text

	case "math_calculate":
		// Extract math operation
		params["operation"] = "add" // Default for demo
		params["a"] = 10.0
		params["b"] = 5.0

	case "remember":
		// Extract key-value for memory
		params["key"] = "react_query"
		params["value"] = query

	case "recall":
		params["key"] = "react_query"

	case "uppercase":
		text := extractTextFromQuery(query, "uppercase")
		params["text"] = text

	case "web_search_mock":
		searchQuery := extractTextFromQuery(query, "search")
		params["query"] = searchQuery

	default:
		params["text"] = query
	}

	return tools.Call(action, params, memory)
}

// extractTextFromQuery extracts relevant text for different operations
func extractTextFromQuery(query, operation string) string {
	// Simple text extraction - in a real implementation, this would be more sophisticated
	switch operation {
	case "sentiment":
		// Look for text after "sentiment of" or similar phrases
		if idx := strings.Index(strings.ToLower(query), "sentiment of"); idx != -1 {
			return strings.TrimSpace(query[idx+12:])
		}
		if idx := strings.Index(strings.ToLower(query), "analyze"); idx != -1 {
			return strings.TrimSpace(query[idx+7:])
		}
	case "uppercase":
		// Look for text to convert
		if idx := strings.Index(strings.ToLower(query), "uppercase"); idx != -1 {
			return strings.TrimSpace(query[:idx])
		}
	case "search":
		// Look for search terms
		if idx := strings.Index(strings.ToLower(query), "search for"); idx != -1 {
			return strings.TrimSpace(query[idx+10:])
		}
	}

	return query
}

// generateObservation creates a summary of what was accomplished
func generateObservation(query string, actions []string, memory *mcp.Memory) string {
	if len(actions) == 0 {
		return "No actions were executed."
	}

	if len(actions) == 1 {
		return fmt.Sprintf("Successfully executed %s based on the user's request.", actions[0])
	}

	return fmt.Sprintf("Successfully executed %d actions (%s) to address the user's multi-step request.",
		len(actions), strings.Join(actions, ", "))
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// parseAndCalculate parses simple mathematical expressions from text
func parseAndCalculate(text string) map[string]interface{} {
	text = strings.ToLower(strings.TrimSpace(text))

	// Simple pattern matching for common mathematical operations
	if strings.Contains(text, "multiply") || strings.Contains(text, "*") {
		// Extract numbers for multiplication
		if numbers := extractTwoNumbers(text); len(numbers) == 2 {
			result := numbers[0] * numbers[1]
			return map[string]interface{}{
				"result":    result,
				"operation": "multiply",
				"operands":  numbers,
				"text":      text,
			}
		}
	}

	if strings.Contains(text, "add") || strings.Contains(text, "+") {
		if numbers := extractTwoNumbers(text); len(numbers) == 2 {
			result := numbers[0] + numbers[1]
			return map[string]interface{}{
				"result":    result,
				"operation": "add",
				"operands":  numbers,
				"text":      text,
			}
		}
	}

	if strings.Contains(text, "subtract") || strings.Contains(text, "-") {
		if numbers := extractTwoNumbers(text); len(numbers) == 2 {
			result := numbers[0] - numbers[1]
			return map[string]interface{}{
				"result":    result,
				"operation": "subtract",
				"operands":  numbers,
				"text":      text,
			}
		}
	}

	if strings.Contains(text, "divide") || strings.Contains(text, "/") {
		if numbers := extractTwoNumbers(text); len(numbers) == 2 {
			if numbers[1] == 0 {
				return map[string]interface{}{"error": "division by zero", "text": text}
			}
			result := numbers[0] / numbers[1]
			return map[string]interface{}{
				"result":    result,
				"operation": "divide",
				"operands":  numbers,
				"text":      text,
			}
		}
	}

	return map[string]interface{}{
		"error": "Could not parse mathematical expression",
		"text":  text,
	}
}

// extractTwoNumbers extracts two numbers from text
func extractTwoNumbers(text string) []float64 {
	// Simple regex to find numbers (including decimals)
	numbers := []float64{}

	// Split by common separators and look for numbers
	words := strings.FieldsFunc(text, func(c rune) bool {
		return c == ' ' || c == ',' || c == '*' || c == '+' || c == '-' || c == '/' || c == '(' || c == ')'
	})

	for _, word := range words {
		if num, err := strconv.ParseFloat(word, 64); err == nil {
			numbers = append(numbers, num)
			if len(numbers) >= 2 {
				break
			}
		}
	}

	return numbers
}
