// LangChain Go + MCP Integration Example (with Ollama)
//
// This example demonstrates seamless integration between LangChain Go agents
// and our MCP (Model Context Protocol) tools using Ollama for local LLM inference.
//
// Features:
// - LangChain Go agent reasoning with Ollama models
// - MCP tools wrapped as LangChain tools
// - Natural language task execution
// - HTML page generation with Tailwind CSS
// - Memory persistence across tool executions
//
// Prerequisites:
//  1. Ollama installed and running (ollama serve)
//  2. A model pulled (e.g., ollama pull llama3.2)
//
// Usage:
//
//	# Default (llama3.2 model, localhost:11434)
//	go run examples/langchain_mcp_integration/main.go
//
//	# Custom configuration
//	export OLLAMA_URL="http://localhost:11434"
//	export OLLAMA_MODEL="llama3.2"
//	go run examples/langchain_mcp_integration/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/tools"

	conduit "github.com/benozo/conduit/lib"
	conduitTools "github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

// MCPTool wraps our MCP tools to work with LangChain Go
type MCPTool struct {
	name        string
	description string
	mcpServer   *conduit.EnhancedServer
	toolName    string
}

// Name returns the tool name
func (t MCPTool) Name() string {
	return t.name
}

// Description returns the tool description
func (t MCPTool) Description() string {
	return t.description
} // Call executes the MCP tool
func (t MCPTool) Call(ctx context.Context, input string) (string, error) {
	// Parse input as parameters (simple key=value format for demo)
	params := make(map[string]interface{})

	// Clean up input
	input = strings.TrimSpace(input)
	if input == "" {
		// For tools that don't need input, use reasonable defaults
		if t.toolName == "list_memories" || t.toolName == "clear_memory" {
			// These tools don't need parameters
		} else {
			return "", fmt.Errorf("empty input provided for tool %s", t.toolName)
		}
	}

	// For text-based tools, use the input directly as "text" parameter
	if strings.Contains(t.toolName, "text") ||
		t.toolName == "uppercase" || t.toolName == "lowercase" ||
		t.toolName == "trim" || t.toolName == "reverse" {
		params["text"] = input
	} else if t.toolName == "remember" {
		// For memory tools, try to parse key=value or use simpler format
		if strings.Contains(input, "=") {
			parts := strings.SplitN(input, "=", 2)
			if len(parts) == 2 {
				params["key"] = strings.TrimSpace(parts[0])
				params["value"] = strings.TrimSpace(parts[1])
			}
		} else {
			// Try to extract key and value from natural language
			words := strings.Fields(input)
			if len(words) >= 2 {
				params["key"] = words[0]
				params["value"] = strings.Join(words[1:], " ")
			} else {
				params["key"] = "data"
				params["value"] = input
			}
		}
	} else if t.toolName == "recall" {
		// For recall, the input is the key
		if input == "" {
			params["key"] = "data" // default key
		} else {
			params["key"] = input
		}
	} else {
		params["text"] = input
	}

	// Get the tool function from MCP server
	toolRegistry := t.mcpServer.GetToolRegistry()
	if toolRegistry == nil {
		return "", fmt.Errorf("tool registry not available")
	}

	// Execute the tool
	result, err := toolRegistry.Call(t.toolName, params, t.mcpServer.GetMemory())
	if err != nil {
		return "", fmt.Errorf("MCP tool execution failed: %w", err)
	}

	// Convert result to string
	if resultMap, ok := result.(map[string]interface{}); ok {
		if output, exists := resultMap["result"]; exists {
			return fmt.Sprintf("%v", output), nil
		}
		if output, exists := resultMap["output"]; exists {
			return fmt.Sprintf("%v", output), nil
		}
		return fmt.Sprintf("%v", result), nil
	}

	return fmt.Sprintf("%v", result), nil
}

// MCPHTMLTool specifically for HTML creation
type MCPHTMLTool struct {
	mcpServer *conduit.EnhancedServer
	outputDir string
}

func (t MCPHTMLTool) Name() string {
	return "create_html_page"
}

func (t MCPHTMLTool) Description() string {
	return "Create HTML landing pages with Tailwind CSS. Input format: filename|content"
}

func (t MCPHTMLTool) Call(ctx context.Context, input string) (string, error) {
	// Parse input as filename|content or try to extract from natural language
	var filename, content string

	if strings.Contains(input, "|") {
		// Direct format: filename|content
		parts := strings.SplitN(input, "|", 2)
		filename = strings.TrimSpace(parts[0])
		content = strings.TrimSpace(parts[1])
	} else {
		// Try to extract from natural language request
		lines := strings.Split(input, "\n")

		// Look for filename in the first few lines
		for i, line := range lines {
			if i > 5 { // Don't search too far
				break
			}
			if strings.Contains(strings.ToLower(line), "page") ||
				strings.Contains(strings.ToLower(line), "file") ||
				strings.Contains(strings.ToLower(line), "html") {
				// Extract potential filename
				words := strings.Fields(line)
				for _, word := range words {
					if len(word) > 2 && !strings.Contains(word, " ") {
						filename = word
						break
					}
				}
				break
			}
		}

		// If no filename found, use default
		if filename == "" {
			filename = "generated-page"
		}

		// Look for HTML content (starts with <!DOCTYPE or <html)
		contentStart := -1
		for i, line := range lines {
			if strings.Contains(line, "<!DOCTYPE") || strings.Contains(line, "<html") {
				contentStart = i
				break
			}
		}

		if contentStart >= 0 {
			content = strings.Join(lines[contentStart:], "\n")
		} else {
			// Generate simple HTML content
			content = fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50 p-8">
    <div class="max-w-4xl mx-auto">
        <h1 class="text-4xl font-bold text-blue-900 mb-4">Generated Page</h1>
        <p class="text-lg text-gray-700">%s</p>
    </div>
</body>
</html>`, filename, input)
		}
	}

	params := map[string]interface{}{
		"filename": filename,
		"content":  content,
	}

	// Execute the tool
	toolRegistry := t.mcpServer.GetToolRegistry()
	_, err := toolRegistry.Call("create_html_page", params, t.mcpServer.GetMemory())
	if err != nil {
		return "", fmt.Errorf("HTML creation failed: %w", err)
	}

	return fmt.Sprintf("Created HTML file: %s", filename), nil
}

func main() {
	// Initialize MCP server
	log.Printf("🚀 Starting LangChain Go + MCP Integration with Ollama...")

	// Check Ollama configuration
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://192.168.10.10:11434"
	}

	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3.2"
	}

	log.Printf("🦙 Ollama URL: %s", ollamaURL)
	log.Printf("📦 Model: %s", modelName)
	log.Printf("💡 Tip: Make sure Ollama is running (ollama serve) and model is pulled (ollama pull %s)", modelName)
	config := conduit.DefaultConfig()
	config.EnableLogging = true
	mcpServer := conduit.NewEnhancedServer(config)

	// Register MCP tools
	conduitTools.RegisterTextTools(mcpServer)
	conduitTools.RegisterMemoryTools(mcpServer)
	conduitTools.RegisterUtilityTools(mcpServer)

	// Create output directory for HTML
	outputDir := "./generated_pages"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Register HTML creation tool
	mcpServer.RegisterToolWithSchema("create_html_page",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			filename, ok := params["filename"].(string)
			if !ok {
				return nil, fmt.Errorf("filename parameter required")
			}

			content, ok := params["content"].(string)
			if !ok {
				return nil, fmt.Errorf("content parameter required")
			}

			if !strings.HasSuffix(filename, ".html") {
				filename += ".html"
			}

			filePath := fmt.Sprintf("%s/%s", outputDir, filename)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return nil, fmt.Errorf("failed to write file: %w", err)
			}

			fmt.Printf("📄 Created HTML: %s\n", filename)
			return map[string]interface{}{
				"filepath": filePath,
				"filename": filename,
				"status":   "created",
			}, nil
		},
		conduit.CreateToolMetadata("create_html_page", "Create HTML landing page", map[string]interface{}{
			"filename": conduit.StringParam("HTML filename"),
			"content":  conduit.StringParam("Complete HTML content"),
		}, []string{"filename", "content"}))

	// Create LangChain LLM with Ollama
	log.Printf("🦙 Using Ollama at: %s", ollamaURL)
	log.Printf("📦 Using model: %s", modelName)

	llm, err := ollama.New(
		ollama.WithServerURL(ollamaURL),
		ollama.WithModel(modelName),
	)
	if err != nil {
		log.Fatalf("Failed to create Ollama LLM: %v", err)
	}

	// Create MCP tools wrapped for LangChain
	mcpTools := []tools.Tool{
		MCPTool{
			name:        "uppercase",
			description: "Convert text to uppercase. Provide the text to convert.",
			mcpServer:   mcpServer,
			toolName:    "uppercase",
		},
		MCPTool{
			name:        "lowercase",
			description: "Convert text to lowercase. Provide the text to convert.",
			mcpServer:   mcpServer,
			toolName:    "lowercase",
		},
		MCPTool{
			name:        "remember",
			description: "Store data in memory. Use format: key=value or provide key and value separately.",
			mcpServer:   mcpServer,
			toolName:    "remember",
		},
		MCPTool{
			name:        "recall",
			description: "Retrieve data from memory. Provide the key to look up.",
			mcpServer:   mcpServer,
			toolName:    "recall",
		},
		MCPHTMLTool{
			mcpServer: mcpServer,
			outputDir: outputDir,
		},
		tools.Calculator{},
	}

	// Create LangChain agent with increased iterations for complex tasks
	agent := agents.NewOneShotAgent(
		llm,
		mcpTools,
		agents.WithMaxIterations(15), // Increased from 10 to handle more complex reasoning
	)

	executor := agents.NewExecutor(agent)

	// Test the integrated system
	fmt.Println("🤖 LangChain Go + MCP Integration Demo")
	fmt.Println("=====================================")

	// Test 1: Simple text processing
	fmt.Println("\n📝 Test 1: Text Processing")
	question1 := "Use the uppercase tool to convert 'hello world' to uppercase"
	answer1, err := chains.Run(context.Background(), executor, question1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Q: %s\nA: %s\n", question1, answer1)
	}

	// Test 2: Memory operation with simple format
	fmt.Println("\n🧮 Test 2: Memory Operation")
	question2 := "Use the remember tool to store the value 42 with key answer"
	answer2, err := chains.Run(context.Background(), executor, question2)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Q: %s\nA: %s\n", question2, answer2)
	}

	// Test 3: HTML creation with simple description
	fmt.Println("\n🎨 Test 3: HTML Creation")
	question3 := "Use the create_html_page tool to create a simple demo page"
	answer3, err := chains.Run(context.Background(), executor, question3)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Q: %s\nA: %s\n", question3, answer3)
	}

	fmt.Println("\n✅ Integration demo completed!")
	fmt.Printf("🌐 Check generated files in: %s\n", outputDir)
}
