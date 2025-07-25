// Advanced MCP Client Example
//
// This example demonstrates how to use the Go MCP client SDK to:
// - Connect to an MCP server
// - Call advanced tools with complex parameters
// - Access resources with different content types
// - Use dynamic prompts with arguments
// - Handle errors and timeouts properly
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/client"
	"github.com/modelcontextprotocol/go-sdk/protocol"
	"github.com/modelcontextprotocol/go-sdk/transport"
)

func main() {
	// Create client with advanced configuration
	opts := &client.ClientOptions{
		Timeout:        30 * time.Second,
		ConnectTimeout: 10 * time.Second,
		ClientInfo: protocol.Implementation{
			Name:    "advanced-example-client",
			Version: "1.0.0",
		},
		Logger: &simpleLogger{},
	}

	// For this example, we'll simulate a connection to a server
	// In practice, you would connect to a real server via STDIO, HTTP, or WebSocket
	transport := transport.NewStdioTransport(&transport.StdioOptions{
		TransportOptions: &transport.TransportOptions{
			Debug:  true,
			Logger: &simpleLogger{},
		},
	})

	mcpClient := client.NewClient(transport, opts)

	// Connect to server
	ctx := context.Background()
	capabilities := protocol.ClientCapabilities{
		Experimental: map[string]interface{}{},
		Roots: &protocol.RootsCapability{
			ListChanged: false,
		},
		Sampling: &protocol.SamplingCapability{},
	}

	if err := mcpClient.Connect(ctx, capabilities); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer mcpClient.Close()

	fmt.Println("✓ Connected to MCP server")
	fmt.Printf("✓ Server: %s %s\n", mcpClient.GetServerInfo().Name, mcpClient.GetServerInfo().Version)

	// Demonstrate tool usage
	if err := demonstrateTools(ctx, mcpClient); err != nil {
		log.Printf("Tool demonstration failed: %v", err)
	}

	// Demonstrate resource access
	if err := demonstrateResources(ctx, mcpClient); err != nil {
		log.Printf("Resource demonstration failed: %v", err)
	}

	// Demonstrate prompt usage
	if err := demonstratePrompts(ctx, mcpClient); err != nil {
		log.Printf("Prompt demonstration failed: %v", err)
	}

	fmt.Println("✓ All demonstrations completed")
}

func demonstrateTools(ctx context.Context, client *client.Client) error {
	fmt.Println("\n=== Tool Demonstrations ===")

	// List available tools
	tools, err := client.ListTools(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	fmt.Printf("Available tools: %d\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}

	// Test advanced text transform tool
	fmt.Println("\n1. Testing Advanced Text Transform:")
	textParams := map[string]interface{}{
		"text":       "Hello World Example",
		"operations": []string{"lowercase", "reverse", "count_chars"},
	}

	result, err := client.CallTool(ctx, "advanced_text_transform", textParams)
	if err != nil {
		return fmt.Errorf("text transform failed: %w", err)
	}

	fmt.Println("Result:")
	for _, content := range result.Content {
		fmt.Printf("  %s\n", content.Text)
	}

	// Test file operations tool
	fmt.Println("\n2. Testing File Operations:")
	fileParams := map[string]interface{}{
		"operation": "list",
		"path":      "/tmp",
	}

	result, err = client.CallTool(ctx, "file_operations", fileParams)
	if err != nil {
		return fmt.Errorf("file operations failed: %w", err)
	}

	fmt.Println("Directory listing:")
	for _, content := range result.Content {
		if result.IsError {
			fmt.Printf("  Error: %s\n", content.Text)
		} else {
			fmt.Printf("  %s\n", content.Text)
		}
	}

	// Test error handling with invalid parameters
	fmt.Println("\n3. Testing Error Handling:")
	invalidParams := map[string]interface{}{
		"text": 123, // Invalid type - should be string
	}

	result, err = client.CallTool(ctx, "advanced_text_transform", invalidParams)
	if err != nil {
		fmt.Printf("  Expected error: %v\n", err)
	} else if result.IsError {
		fmt.Printf("  Tool returned error: %s\n", result.Content[0].Text)
	}

	return nil
}

func demonstrateResources(ctx context.Context, client *client.Client) error {
	fmt.Println("\n=== Resource Demonstrations ===")

	// List available resources
	resources, err := client.ListResources(ctx)
	if err != nil {
		return fmt.Errorf("failed to list resources: %w", err)
	}

	fmt.Printf("Available resources: %d\n", len(resources))
	for _, resource := range resources {
		fmt.Printf("  - %s (%s): %s\n", resource.Name, resource.MimeType, resource.Description)
	}

	// Read configuration resource
	fmt.Println("\n1. Reading Configuration Resource:")
	configResponse, err := client.ReadResource(ctx, "file:///config/app.json")
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	for _, content := range configResponse.Contents {
		fmt.Printf("  URI: %s\n", content.URI)
		fmt.Printf("  Type: %s\n", content.MimeType)

		// Pretty print JSON if it's a JSON resource
		if content.MimeType == "application/json" {
			var jsonData interface{}
			if err := json.Unmarshal([]byte(content.Text), &jsonData); err == nil {
				prettyJSON, _ := json.MarshalIndent(jsonData, "  ", "  ")
				fmt.Printf("  Content:\n%s\n", string(prettyJSON))
			} else {
				fmt.Printf("  Content: %s\n", content.Text)
			}
		} else {
			fmt.Printf("  Content: %s\n", content.Text)
		}
	}

	// Read documentation resource
	fmt.Println("\n2. Reading Documentation Resource:")
	docsResponse, err := client.ReadResource(ctx, "file:///docs/README.md")
	if err != nil {
		return fmt.Errorf("failed to read docs: %w", err)
	}

	for _, content := range docsResponse.Contents {
		fmt.Printf("  URI: %s\n", content.URI)
		fmt.Printf("  Type: %s\n", content.MimeType)
		fmt.Printf("  Content (first 200 chars):\n  %s...\n",
			truncateString(content.Text, 200))
	}

	return nil
}

func demonstratePrompts(ctx context.Context, client *client.Client) error {
	fmt.Println("\n=== Prompt Demonstrations ===")

	// List available prompts
	prompts, err := client.ListPrompts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list prompts: %w", err)
	}

	fmt.Printf("Available prompts: %d\n", len(prompts))
	for _, prompt := range prompts {
		fmt.Printf("  - %s: %s\n", prompt.Name, prompt.Description)
		fmt.Println("    Arguments:")
		for _, arg := range prompt.Arguments {
			required := ""
			if arg.Required {
				required = " (required)"
			}
			fmt.Printf("      - %s: %s%s\n", arg.Name, arg.Description, required)
		}
	}

	// Get code review prompt
	fmt.Println("\n1. Getting Code Review Prompt:")
	reviewArgs := map[string]interface{}{
		"language":   "go",
		"complexity": "high",
	}

	reviewResponse, err := client.GetPrompt(ctx, "code_review", reviewArgs)
	if err != nil {
		return fmt.Errorf("failed to get code review prompt: %w", err)
	}

	fmt.Printf("  Description: %s\n", reviewResponse.Description)
	fmt.Printf("  Messages: %d\n", len(reviewResponse.Messages))
	for i, message := range reviewResponse.Messages {
		fmt.Printf("  Message %d (%s):\n", i+1, message.Role)
		fmt.Printf("    %s\n", truncateString(message.Content.Text, 150))
	}

	// Get documentation prompt
	fmt.Println("\n2. Getting Documentation Prompt:")
	docArgs := map[string]interface{}{
		"type":     "api",
		"audience": "developers",
	}

	docResponse, err := client.GetPrompt(ctx, "documentation", docArgs)
	if err != nil {
		return fmt.Errorf("failed to get documentation prompt: %w", err)
	}

	fmt.Printf("  Description: %s\n", docResponse.Description)
	fmt.Printf("  Messages: %d\n", len(docResponse.Messages))
	for i, message := range docResponse.Messages {
		fmt.Printf("  Message %d (%s):\n", i+1, message.Role)
		fmt.Printf("    %s\n", truncateString(message.Content.Text, 150))
	}

	// Test prompt with missing required arguments
	fmt.Println("\n3. Testing Error Handling with Missing Arguments:")
	_, err = client.GetPrompt(ctx, "code_review", map[string]interface{}{
		"complexity": "medium", // Missing required "language" argument
	})
	if err != nil {
		fmt.Printf("  Expected error: %v\n", err)
	}

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

type simpleLogger struct{}

func (l *simpleLogger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] "+msg, args...)
}

func (l *simpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *simpleLogger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg, args...)
}

func (l *simpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}
