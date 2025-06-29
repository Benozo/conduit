package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Example: ReAct (Reasoning + Acting) Agent with conduit

	// Create configuration
	config := conduit.DefaultConfig()
	config.Port = 8085
	config.Mode = mcp.ModeHTTP
	config.EnableLogging = true

	// Create server
	server := conduit.NewServer(config)

	// Register comprehensive tool set for the agent
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Add ReAct-specific tools
	server.RegisterTool("analyze_sentiment", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
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
	})

	server.RegisterTool("math_calculate", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
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
	})

	server.RegisterTool("web_search_mock", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
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
	})

	server.RegisterTool("decision_maker", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		options := params["options"].([]interface{})
		criteria := params["criteria"].(string)

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
	})

	// Create a ReAct model that uses the MCP package's ReAct agent
	reactModel := func(ctx mcp.ContextInput, req mcp.MCPRequest, memory *mcp.Memory, onToken mcp.StreamCallback) (string, error) {
		query := ctx.Inputs["query"].(string)

		// Create tool registry for direct MCP usage
		toolRegistry := mcp.NewToolRegistry()

		// Register basic tools in the MCP registry
		toolRegistry.Register("uppercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			text := params["text"].(string)
			return strings.ToUpper(text), nil
		})

		toolRegistry.Register("analyze_sentiment", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			text := params["text"].(string)
			// Simple sentiment analysis
			if strings.Contains(strings.ToLower(text), "good") || strings.Contains(strings.ToLower(text), "great") {
				return map[string]interface{}{"sentiment": "positive", "confidence": 0.8}, nil
			} else if strings.Contains(strings.ToLower(text), "bad") || strings.Contains(strings.ToLower(text), "terrible") {
				return map[string]interface{}{"sentiment": "negative", "confidence": 0.8}, nil
			}
			return map[string]interface{}{"sentiment": "neutral", "confidence": 0.5}, nil
		})

		// Store the query in memory for ReAct processing
		memory.Set("latest", query)

		// Define thoughts for ReAct processing
		thoughts := []string{
			"analyze the input text",
			"transform to uppercase",
			"store results",
		}

		// Use MCP's ReAct agent
		steps, err := mcp.ReActAgent(thoughts, toolRegistry, memory)
		if err != nil {
			return "Error in ReAct processing: " + err.Error(), err
		}

		// Build response from ReAct steps
		var response strings.Builder
		response.WriteString("ReAct Agent Processing:\n\n")

		for i, step := range steps {
			response.WriteString(fmt.Sprintf("Step %d:\n", i+1))
			response.WriteString(fmt.Sprintf("Thought: %s\n", step.Thought))
			response.WriteString(fmt.Sprintf("Action: %s\n", step.Action))
			if len(step.Params) > 0 {
				response.WriteString(fmt.Sprintf("Params: %v\n", step.Params))
			}
			response.WriteString(fmt.Sprintf("Observation: %s\n\n", step.Observed))
		}

		// Add final result
		finalResult := memory.Get("latest_result")
		if finalResult != nil {
			response.WriteString(fmt.Sprintf("Final Result: %v\n", finalResult))
		}

		// Simulate streaming if callback provided
		responseText := response.String()
		if onToken != nil {
			words := strings.Split(responseText, " ")
			for _, word := range words {
				onToken(ctx.ContextID, word+" ")
			}
		}

		return responseText, nil
	}

	// Set the model
	server.SetModel(reactModel)

	log.Printf("Starting ReAct Agent MCP server on port %d...", config.Port)
	log.Printf("")
	log.Printf("ReAct Pattern: Reasoning + Acting")
	log.Printf("Available tools for the agent:")
	log.Printf("  • Text processing tools (uppercase, lowercase, etc.)")
	log.Printf("  • Memory tools (remember, recall, etc.)")
	log.Printf("  • Utility tools (timestamp, uuid, etc.)")
	log.Printf("  • Analysis tools (analyze_sentiment)")
	log.Printf("  • Calculation tools (math_calculate)")
	log.Printf("  • Search tools (web_search_mock)")
	log.Printf("  • Decision tools (decision_maker)")
	log.Printf("")
	log.Printf("Try these endpoints:")
	log.Printf("  GET  http://localhost:%d/health", config.Port)
	log.Printf("  GET  http://localhost:%d/schema", config.Port)
	log.Printf("  POST http://localhost:%d/react", config.Port)
	log.Printf("")
	log.Printf("Example ReAct request:")
	log.Printf(`  curl -X POST http://localhost:%d/react \`, config.Port)
	log.Printf(`    -H "Content-Type: application/json" \`)
	log.Printf(`    -d '{"thoughts": "I need to analyze some text and remember the results"}'`)

	// Use stdio mode if requested
	if len(os.Args) > 1 && os.Args[1] == "--stdio" {
		config.Mode = mcp.ModeStdio
		log.Printf("\nSwitching to stdio mode for VS Code Copilot integration...")
		server = conduit.NewServer(config) // Recreate server with new config
	}

	log.Fatal(server.Start())
}
