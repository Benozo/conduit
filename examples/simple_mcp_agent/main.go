package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	conduit "github.com/benozo/conduit/lib"
	conduitTools "github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

// SimpleMCPAgent demonstrates MCP tool usage without LangChain dependency
type SimpleMCPAgent struct {
	mcpServer *conduit.EnhancedServer
	memory    *mcp.Memory
}

func NewSimpleMCPAgent() *SimpleMCPAgent {
	// Initialize MCP server
	config := conduit.DefaultConfig()
	config.EnableLogging = true
	mcpServer := conduit.NewEnhancedServer(config)

	// Register MCP tools
	conduitTools.RegisterTextTools(mcpServer)
	conduitTools.RegisterMemoryTools(mcpServer)
	conduitTools.RegisterUtilityTools(mcpServer)

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

			outputDir := "./generated_pages"
			os.MkdirAll(outputDir, 0755)
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

	return &SimpleMCPAgent{
		mcpServer: mcpServer,
		memory:    mcpServer.GetMemory(),
	}
}

func (a *SimpleMCPAgent) ExecuteTool(toolName string, params map[string]interface{}) (interface{}, error) {
	toolRegistry := a.mcpServer.GetToolRegistry()
	return toolRegistry.Call(toolName, params, a.memory)
}

func (a *SimpleMCPAgent) RunWorkflow(ctx context.Context) error {
	fmt.Println("🤖 Simple MCP Agent Workflow Demo")
	fmt.Println("==================================")

	// Step 1: Text processing
	fmt.Println("\n📝 Step 1: Text Processing")
	result1, err := a.ExecuteTool("uppercase", map[string]interface{}{
		"text": "hello world from mcp",
	})
	if err != nil {
		return fmt.Errorf("uppercase failed: %w", err)
	}
	fmt.Printf("Uppercase result: %v\n", result1)

	// Step 2: Store in memory
	fmt.Println("\n💾 Step 2: Memory Storage")
	_, err = a.ExecuteTool("remember", map[string]interface{}{
		"key":   "greeting",
		"value": result1,
	})
	if err != nil {
		return fmt.Errorf("remember failed: %w", err)
	}
	fmt.Println("Stored greeting in memory")

	// Step 3: Recall from memory
	fmt.Println("\n🧠 Step 3: Memory Recall")
	result3, err := a.ExecuteTool("recall", map[string]interface{}{
		"key": "greeting",
	})
	if err != nil {
		return fmt.Errorf("recall failed: %w", err)
	}
	fmt.Printf("Recalled from memory: %v\n", result3)

	// Step 4: Create HTML page
	fmt.Println("\n🎨 Step 4: HTML Creation")
	htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MCP Demo</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gradient-to-br from-blue-50 to-purple-50 min-h-screen">
    <div class="container mx-auto px-4 py-16">
        <div class="max-w-4xl mx-auto text-center">
            <h1 class="text-5xl font-bold text-gray-900 mb-6">MCP Tools Integration</h1>
            <p class="text-xl text-gray-600 mb-8">This page was created using Model Context Protocol tools!</p>
            
            <div class="bg-white rounded-lg shadow-lg p-8 mb-8">
                <h2 class="text-2xl font-semibold text-gray-800 mb-4">Workflow Results</h2>
                <div class="space-y-4">
                    <div class="p-4 bg-blue-50 rounded">
                        <span class="font-medium">Text Processing:</span> ` + fmt.Sprintf("%v", result1) + `
                    </div>
                    <div class="p-4 bg-green-50 rounded">
                        <span class="font-medium">Memory Storage:</span> Stored as 'greeting'
                    </div>
                    <div class="p-4 bg-purple-50 rounded">
                        <span class="font-medium">Memory Recall:</span> ` + fmt.Sprintf("%v", result3) + `
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div class="bg-white p-6 rounded-lg shadow-md">
                    <h3 class="text-lg font-semibold mb-2">🔧 MCP Tools</h3>
                    <p class="text-gray-600">Modular, composable functions</p>
                </div>
                <div class="bg-white p-6 rounded-lg shadow-md">
                    <h3 class="text-lg font-semibold mb-2">💾 Memory</h3>
                    <p class="text-gray-600">Persistent data storage</p>
                </div>
                <div class="bg-white p-6 rounded-lg shadow-md">
                    <h3 class="text-lg font-semibold mb-2">🎨 HTML Gen</h3>
                    <p class="text-gray-600">Dynamic page creation</p>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

	_, err = a.ExecuteTool("create_html_page", map[string]interface{}{
		"filename": "mcp_demo",
		"content":  htmlContent,
	})
	if err != nil {
		return fmt.Errorf("HTML creation failed: %w", err)
	}

	// Step 5: Generate timestamp
	fmt.Println("\n⏰ Step 5: Timestamp Generation")
	result5, err := a.ExecuteTool("timestamp", map[string]interface{}{
		"format": "iso",
	})
	if err != nil {
		return fmt.Errorf("timestamp failed: %w", err)
	}
	fmt.Printf("Current timestamp: %v\n", result5)

	return nil
}

func main() {
	// Create agent
	agent := NewSimpleMCPAgent()

	// Run the workflow
	if err := agent.RunWorkflow(context.Background()); err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	fmt.Println("\n✅ MCP Agent workflow completed successfully!")
	fmt.Println("🌐 Check generated HTML at: ./generated_pages/mcp_demo.html")

	// Show available tools
	fmt.Println("\n🔧 Available MCP Tools:")
	toolNames := []string{
		"uppercase", "lowercase", "trim", "reverse",
		"remember", "recall", "clear_memory", "list_memories",
		"timestamp", "uuid", "hash_md5", "hash_sha256",
		"base64_encode", "base64_decode", "url_encode", "url_decode",
		"create_html_page",
	}

	for i, tool := range toolNames {
		if i%4 == 0 {
			fmt.Println()
		}
		fmt.Printf("  %-18s", tool)
	}
	fmt.Println()
}
