package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	// Get DeepInfra configuration from environment
	bearerToken := os.Getenv("DEEPINFRA_TOKEN")
	if bearerToken == "" {
		fmt.Println("❌ DEEPINFRA_TOKEN environment variable is required")
		return
	}

	modelName := os.Getenv("DEEPINFRA_MODEL")
	if modelName == "" {
		// modelName = "meta-llama/Meta-Llama-3.1-8B-Instruct"
		modelName = "Qwen/Qwen3-14B"
	}

	// Create output directory
	outputDir := "./generated_pages"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create MCP server
	config := conduit.DefaultConfig()
	config.EnableLogging = true
	server := conduit.NewEnhancedServer(config)

	// Register tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// HTML page creation tool
	server.RegisterToolWithSchema("create_html_page",
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

			filePath := filepath.Join(outputDir, filename)

			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return nil, fmt.Errorf("failed to write file: %v", err)
			}

			fmt.Printf("📄 HTML: %s\n", filename)

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

	// Create DeepInfra model
	deepInfraModel := conduit.CreateDeepInfraModel(bearerToken)

	// Create agent manager
	llmAgentManager := agents.NewLLMAgentManager(server, deepInfraModel, modelName)

	// Create HTML builder agent
	_, err := llmAgentManager.CreateLLMAgent(
		"html_builder",
		"HTML + Tailwind Builder",
		"Creates HTML landing pages with Tailwind CSS",
		`You are an expert HTML + Tailwind CSS developer that creates landing pages.

TASK: Create a complete HTML landing page with the specified requirements.

TOOLS AVAILABLE:
- create_html_page(filename, content): Create HTML file with complete content

IMPORTANT: Respond ONLY with valid JSON. Do not use thinking tags or explanations.

Response format (JSON ONLY):
{
  "analysis": "Brief analysis of requirements",
  "steps": [
    {
      "name": "create_page",
      "description": "Create HTML landing page",
      "tool": "create_html_page",
      "input": {
        "filename": "landing.html",
        "content": "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>Page Title</title>\n    <script src=\"https://cdn.tailwindcss.com\"></script>\n</head>\n<body>\n    <!-- Complete page content here -->\n</body>\n</html>"
      }
    }
  ],
  "reasoning": "Design approach explanation"
}

The content field must contain the COMPLETE HTML document with Tailwind CSS.`,
		[]string{"create_html_page", "remember", "recall"},
		&agents.AgentConfig{
			MaxTokens:     2000,
			Temperature:   0.7,
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Create SaaS landing page
	saasTask, _ := llmAgentManager.CreateTask(
		"html_builder",
		"SaaS Analytics Platform",
		"Modern SaaS landing page",
		map[string]interface{}{
			"company":     "DataFlow Analytics",
			"description": "Analytics platform for data-driven decisions",
			"features":    []string{"Real-time data", "Beautiful dashboards", "Team collaboration"},
		},
	)

	fmt.Println("🧠 LLM creating SaaS page...")
	start := time.Now()
	if err := llmAgentManager.ExecuteTaskWithLLM(saasTask.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		elapsed := time.Since(start)
		fmt.Printf("✅ Created in %v\n", elapsed)
	}

	fmt.Println("🎉 Completed!")
}
