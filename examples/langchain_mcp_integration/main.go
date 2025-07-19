package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
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
}

// Call executes the MCP tool
func (t MCPTool) Call(ctx context.Context, input string) (string, error) {
	// Parse input as parameters (simple key=value format for demo)
	params := make(map[string]interface{})

	// For text-based tools, use the input directly as "text" parameter
	if strings.Contains(t.toolName, "text") ||
		t.toolName == "uppercase" || t.toolName == "lowercase" ||
		t.toolName == "trim" || t.toolName == "reverse" {
		params["text"] = input
	} else if t.toolName == "remember" {
		// For memory tools, try to parse key=value
		parts := strings.SplitN(input, "=", 2)
		if len(parts) == 2 {
			params["key"] = strings.TrimSpace(parts[0])
			params["value"] = strings.TrimSpace(parts[1])
		} else {
			params["text"] = input // fallback
		}
	} else if t.toolName == "recall" {
		params["key"] = input
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
	// Parse input as filename|content
	parts := strings.SplitN(input, "|", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("input must be in format: filename|content")
	}

	filename := strings.TrimSpace(parts[0])
	content := strings.TrimSpace(parts[1])

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
	// Check for required environment variables
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Initialize MCP server
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

	// Create LangChain LLM
	llm, err := openai.New()
	if err != nil {
		log.Fatalf("Failed to create OpenAI LLM: %v", err)
	}

	// Create MCP tools wrapped for LangChain
	mcpTools := []tools.Tool{
		MCPTool{
			name:        "uppercase",
			description: "Convert text to uppercase",
			mcpServer:   mcpServer,
			toolName:    "uppercase",
		},
		MCPTool{
			name:        "lowercase",
			description: "Convert text to lowercase",
			mcpServer:   mcpServer,
			toolName:    "lowercase",
		},
		MCPTool{
			name:        "remember",
			description: "Store data in memory. Use format: key=value",
			mcpServer:   mcpServer,
			toolName:    "remember",
		},
		MCPTool{
			name:        "recall",
			description: "Retrieve data from memory. Provide the key.",
			mcpServer:   mcpServer,
			toolName:    "recall",
		},
		MCPHTMLTool{
			mcpServer: mcpServer,
			outputDir: outputDir,
		},
		tools.Calculator{},
	}

	// Create LangChain agent
	agent := agents.NewOneShotAgent(
		llm,
		mcpTools,
		agents.WithMaxIterations(5),
	)

	executor := agents.NewExecutor(agent)

	// Test the integrated system
	fmt.Println("🤖 LangChain Go + MCP Integration Demo")
	fmt.Println("=====================================")

	// Test 1: Basic text processing
	fmt.Println("\n📝 Test 1: Text Processing")
	question1 := "Convert 'hello world' to uppercase and remember it as greeting"
	answer1, err := chains.Run(context.Background(), executor, question1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Q: %s\nA: %s\n", question1, answer1)
	}

	// Test 2: Memory and calculation
	fmt.Println("\n🧮 Test 2: Memory + Calculation")
	question2 := "Remember that the answer is 42, then calculate 42 * 2"
	answer2, err := chains.Run(context.Background(), executor, question2)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Q: %s\nA: %s\n", question2, answer2)
	}

	// Test 3: HTML creation
	fmt.Println("\n🎨 Test 3: HTML Creation")
	question3 := `Create an HTML page called "demo" with this content: <!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LangChain + MCP Demo</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-blue-50 p-8">
    <div class="max-w-4xl mx-auto">
        <h1 class="text-4xl font-bold text-blue-900 mb-4">LangChain Go + MCP Integration</h1>
        <p class="text-lg text-gray-700">This page was created by LangChain Go using MCP tools!</p>
    </div>
</body>
</html>`
	answer3, err := chains.Run(context.Background(), executor, question3)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Q: Create HTML demo page\nA: %s\n", answer3)
	}

	fmt.Println("\n✅ Integration demo completed!")
	fmt.Printf("🌐 Check generated files in: %s\n", outputDir)
}
