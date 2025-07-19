package main

import (
	"bufio"
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

	// Read existing HTML tool
	server.RegisterToolWithSchema("read_existing_html",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			filename, ok := params["filename"].(string)
			if !ok {
				return nil, fmt.Errorf("filename parameter required")
			}

			if !strings.HasSuffix(filename, ".html") {
				filename += ".html"
			}

			filePath := filepath.Join(outputDir, filename)

			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %v", err)
			}

			return map[string]interface{}{
				"content":  string(content),
				"filename": filename,
				"status":   "read",
			}, nil
		},
		conduit.CreateToolMetadata("read_existing_html", "Read existing HTML file", map[string]interface{}{
			"filename": conduit.StringParam("HTML filename to read"),
		}, []string{"filename"}))

	// Create DeepInfra model
	deepInfraModel := conduit.CreateDeepInfraModel(bearerToken)

	// Create agent manager
	llmAgentManager := agents.NewLLMAgentManager(server, deepInfraModel, modelName)

	// Create HTML builder agent
	_, err := llmAgentManager.CreateLLMAgent(
		"html_builder",
		"HTML + Tailwind Builder",
		"Creates and modifies HTML landing pages with Tailwind CSS",
		`You are an expert HTML + Tailwind CSS developer that creates and modifies landing pages.

TASK: Create or modify HTML landing pages based on user requirements.

TOOLS AVAILABLE:
- create_html_page(filename, content): Create HTML file with complete content
- read_existing_html(filename): Read existing HTML file content
- remember(key, value): Store data for later use
- recall(key): Retrieve stored data

IMPORTANT: Respond ONLY with valid JSON. No explanations, no thinking tags, no markdown.

For NEW pages, create complete HTML with Tailwind CSS.
For MODIFICATIONS, read existing HTML first, then create updated version.

Response format (JSON ONLY):
{
  "analysis": "Brief analysis of requirements",
  "steps": [
    {
      "name": "step_name",
      "description": "What this step does",
      "tool": "tool_name",
      "input": {
        "param": "value"
      }
    }
  ],
  "reasoning": "Design approach explanation"
}`,
		[]string{"create_html_page", "read_existing_html", "remember", "recall"},
		&agents.AgentConfig{
			MaxTokens:     3000,
			Temperature:   0.7,
			EnableMemory:  true,
			EnableLogging: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Interactive loop for user amendments
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("🎨 HTML Landing Page Generator with Amendments")
	fmt.Println("============================================")

	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Create new landing page")
		fmt.Println("2. Modify existing page")
		fmt.Println("3. List existing pages")
		fmt.Println("4. Exit")
		fmt.Print("\nChoose option (1-4): ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			createNewPage(llmAgentManager, scanner)
		case "2":
			modifyExistingPage(llmAgentManager, scanner, outputDir)
		case "3":
			listExistingPages(outputDir)
		case "4":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid option. Please choose 1-4.")
		}
	}
}

func createNewPage(llmAgentManager *agents.LLMAgentManager, scanner *bufio.Scanner) {
	fmt.Print("\nCompany name: ")
	if !scanner.Scan() {
		return
	}
	company := scanner.Text()

	fmt.Print("Description: ")
	if !scanner.Scan() {
		return
	}
	description := scanner.Text()

	fmt.Print("Key features (comma-separated): ")
	if !scanner.Scan() {
		return
	}
	featuresStr := scanner.Text()
	features := strings.Split(featuresStr, ",")
	for i := range features {
		features[i] = strings.TrimSpace(features[i])
	}

	fmt.Print("Filename (without .html): ")
	if !scanner.Scan() {
		return
	}
	filename := scanner.Text()

	// Create task
	task, _ := llmAgentManager.CreateTask(
		"html_builder",
		fmt.Sprintf("Landing page for %s", company),
		"Create new landing page",
		map[string]interface{}{
			"company":     company,
			"description": description,
			"features":    features,
			"filename":    filename,
			"action":      "create_new",
		},
	)

	fmt.Printf("\n🧠 Creating landing page for %s...\n", company)
	start := time.Now()
	if err := llmAgentManager.ExecuteTaskWithLLM(task.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		elapsed := time.Since(start)
		fmt.Printf("✅ Created in %v\n", elapsed)
		fmt.Printf("📁 File saved as: %s.html\n", filename)
	}
}

func modifyExistingPage(llmAgentManager *agents.LLMAgentManager, scanner *bufio.Scanner, outputDir string) {
	// List available files
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Printf("❌ Error reading directory: %v\n", err)
		return
	}

	htmlFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".html") {
			htmlFiles = append(htmlFiles, file.Name())
		}
	}

	if len(htmlFiles) == 0 {
		fmt.Println("❌ No HTML files found. Create a new page first.")
		return
	}

	fmt.Println("\nAvailable HTML files:")
	for i, file := range htmlFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	fmt.Print("\nChoose file number: ")
	if !scanner.Scan() {
		return
	}

	var fileIndex int
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &fileIndex); err != nil || fileIndex < 1 || fileIndex > len(htmlFiles) {
		fmt.Println("❌ Invalid file number")
		return
	}

	selectedFile := htmlFiles[fileIndex-1]
	filename := strings.TrimSuffix(selectedFile, ".html")

	fmt.Print("\nWhat changes would you like to make? ")
	if !scanner.Scan() {
		return
	}
	amendments := scanner.Text()

	// Create modification task
	task, _ := llmAgentManager.CreateTask(
		"html_builder",
		fmt.Sprintf("Modify %s", selectedFile),
		"Modify existing landing page",
		map[string]interface{}{
			"filename":   filename,
			"amendments": amendments,
			"action":     "modify_existing",
		},
	)

	fmt.Printf("\n🔧 Modifying %s...\n", selectedFile)
	start := time.Now()
	if err := llmAgentManager.ExecuteTaskWithLLM(task.ID); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		elapsed := time.Since(start)
		fmt.Printf("✅ Modified in %v\n", elapsed)
		fmt.Printf("📁 Updated file: %s\n", selectedFile)
	}
}

func listExistingPages(outputDir string) {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Printf("❌ Error reading directory: %v\n", err)
		return
	}

	htmlFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".html") {
			htmlFiles = append(htmlFiles, file.Name())
		}
	}

	if len(htmlFiles) == 0 {
		fmt.Println("📁 No HTML files found.")
		return
	}

	fmt.Println("\n📁 Existing HTML files:")
	for _, file := range htmlFiles {
		fmt.Printf("  • %s\n", file)

		// Show file preview URL
		absPath, _ := filepath.Abs(filepath.Join(outputDir, file))
		fmt.Printf("    🌐 file://%s\n", absPath)
	}
}
