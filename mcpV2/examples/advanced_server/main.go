// Advanced MCP Server Example
//
// This example demonstrates the advanced features of the Go MCP SDK:
// - Middleware integration for logging and authentication
// - Resource management with file system access
// - Prompt management for dynamic content generation
// - Error handling and validation
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/benozo/neuron-mcp/middleware"
	"github.com/benozo/neuron-mcp/protocol"
	"github.com/benozo/neuron-mcp/server"
	"github.com/benozo/neuron-mcp/transport"
)

func main() {
	// Create logging middleware
	loggingMiddleware := middleware.LoggingMiddleware(&simpleLogger{})

	// Create metrics middleware
	metricsMiddleware := middleware.MetricsMiddleware(&simpleMetrics{})

	// Create server with advanced configuration
	opts := &server.ServerOptions{
		Info: protocol.Implementation{
			Name:    "advanced-example-server",
			Version: "1.0.0",
		},
		Capabilities: protocol.ServerCapabilities{
			Tools: &protocol.ToolsCapability{
				ListChanged: false,
			},
			Resources: &protocol.ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
			Prompts: &protocol.PromptsCapability{
				ListChanged: false,
			},
		},
		Logger: &simpleLogger{},
		// Add middleware chain
		Middleware: []middleware.Middleware{
			loggingMiddleware,
			metricsMiddleware,
		},
	}

	srv := server.NewServer(opts)

	// Register tools, resources, and prompts
	if err := registerTools(srv); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	if err := registerResources(srv); err != nil {
		log.Fatalf("Failed to register resources: %v", err)
	}

	if err := registerPrompts(srv); err != nil {
		log.Fatalf("Failed to register prompts: %v", err)
	}

	// Create STDIO transport
	transport := transport.NewStdioTransport(&transport.StdioOptions{
		TransportOptions: &transport.TransportOptions{
			Debug:  true,
			Logger: &simpleLogger{},
		},
	})

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Start the server
	log.Println("Starting advanced MCP server...")
	if err := srv.Serve(ctx, transport); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}

// registerTools adds various tools to the server
func registerTools(srv *server.Server) error {
	// Advanced text processing tool
	textTool := &protocol.Tool{
		Name:        "advanced_text_transform",
		Description: "Advanced text transformation with multiple operations",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"text": {
					Type:        "string",
					Description: "The text to transform",
				},
				"operations": {
					Type:        "array",
					Description: "List of operations to perform",
					Items: &protocol.JSONSchema{
						Type: "string",
						Enum: []interface{}{"uppercase", "lowercase", "reverse", "sort_words", "count_chars"},
					},
				},
			},
			Required: []string{"text", "operations"},
		},
	}

	if err := srv.RegisterTool(textTool, handleAdvancedTextTransform); err != nil {
		return fmt.Errorf("failed to register advanced_text_transform tool: %w", err)
	}

	// File operations tool
	fileTool := &protocol.Tool{
		Name:        "file_operations",
		Description: "Perform file system operations",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"operation": {
					Type:        "string",
					Description: "File operation to perform",
					Enum:        []interface{}{"list", "stat", "mkdir"},
				},
				"path": {
					Type:        "string",
					Description: "File path for the operation",
				},
			},
			Required: []string{"operation", "path"},
		},
	}

	if err := srv.RegisterTool(fileTool, handleFileOperations); err != nil {
		return fmt.Errorf("failed to register file_operations tool: %w", err)
	}

	return nil
}

// registerResources adds file system resources to the server
func registerResources(srv *server.Server) error {
	// Configuration file resource
	configResource := &protocol.Resource{
		URI:         "file:///config/app.json",
		Name:        "Application Configuration",
		Description: "Application configuration file",
		MimeType:    "application/json",
	}

	if err := srv.RegisterResource(configResource, handleConfigResource); err != nil {
		return fmt.Errorf("failed to register config resource: %w", err)
	}

	// Documentation resource
	docsResource := &protocol.Resource{
		URI:         "file:///docs/README.md",
		Name:        "Documentation",
		Description: "Application documentation",
		MimeType:    "text/markdown",
	}

	if err := srv.RegisterResource(docsResource, handleDocsResource); err != nil {
		return fmt.Errorf("failed to register docs resource: %w", err)
	}

	return nil
}

// registerPrompts adds dynamic prompts to the server
func registerPrompts(srv *server.Server) error {
	// Code review prompt
	reviewPrompt := &protocol.Prompt{
		Name:        "code_review",
		Description: "Generate a code review prompt with context",
		Arguments: []protocol.PromptArgument{
			{
				Name:        "language",
				Description: "Programming language",
				Required:    true,
			},
			{
				Name:        "complexity",
				Description: "Code complexity level",
				Required:    false,
			},
		},
	}

	if err := srv.RegisterPrompt(reviewPrompt, handleCodeReviewPrompt); err != nil {
		return fmt.Errorf("failed to register code_review prompt: %w", err)
	}

	// Documentation prompt
	docPrompt := &protocol.Prompt{
		Name:        "documentation",
		Description: "Generate documentation prompts",
		Arguments: []protocol.PromptArgument{
			{
				Name:        "type",
				Description: "Documentation type (api, guide, reference)",
				Required:    true,
			},
			{
				Name:        "audience",
				Description: "Target audience",
				Required:    false,
			},
		},
	}

	if err := srv.RegisterPrompt(docPrompt, handleDocumentationPrompt); err != nil {
		return fmt.Errorf("failed to register documentation prompt: %w", err)
	}

	return nil
}

// Tool handlers

func handleAdvancedTextTransform(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
	text, ok := params["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text parameter must be a string")
	}

	operationsInterface, ok := params["operations"]
	if !ok {
		return nil, fmt.Errorf("operations parameter is required")
	}

	var operations []string
	switch ops := operationsInterface.(type) {
	case []interface{}:
		operations = make([]string, len(ops))
		for i, op := range ops {
			if opStr, ok := op.(string); ok {
				operations[i] = opStr
			} else {
				return nil, fmt.Errorf("all operations must be strings")
			}
		}
	default:
		return nil, fmt.Errorf("operations must be an array")
	}

	result := text
	for _, op := range operations {
		switch op {
		case "uppercase":
			result = strings.ToUpper(result)
		case "lowercase":
			result = strings.ToLower(result)
		case "reverse":
			runes := []rune(result)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			result = string(runes)
		case "sort_words":
			words := strings.Fields(result)
			// Simple bubble sort for demo
			for i := 0; i < len(words)-1; i++ {
				for j := 0; j < len(words)-i-1; j++ {
					if words[j] > words[j+1] {
						words[j], words[j+1] = words[j+1], words[j]
					}
				}
			}
			result = strings.Join(words, " ")
		case "count_chars":
			result = fmt.Sprintf("Character count: %d\nOriginal: %s", len(result), result)
		default:
			return nil, fmt.Errorf("unknown operation: %s", op)
		}
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			{
				Type: "text",
				Text: result,
			},
		},
		IsError: false,
	}, nil
}

func handleFileOperations(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter must be a string")
	}

	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter must be a string")
	}

	var result string
	var err error

	switch operation {
	case "list":
		entries, dirErr := os.ReadDir(path)
		if dirErr != nil {
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					{
						Type: "text",
						Text: fmt.Sprintf("Error listing directory: %v", dirErr),
					},
				},
				IsError: true,
			}, nil
		}

		var fileList []string
		for _, entry := range entries {
			if entry.IsDir() {
				fileList = append(fileList, entry.Name()+"/")
			} else {
				fileList = append(fileList, entry.Name())
			}
		}
		result = strings.Join(fileList, "\n")

	case "stat":
		info, statErr := os.Stat(path)
		if statErr != nil {
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					{
						Type: "text",
						Text: fmt.Sprintf("Error getting file info: %v", statErr),
					},
				},
				IsError: true,
			}, nil
		}

		result = fmt.Sprintf("Name: %s\nSize: %d bytes\nMode: %s\nModTime: %s\nIsDir: %t",
			info.Name(), info.Size(), info.Mode(), info.ModTime(), info.IsDir())

	case "mkdir":
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					{
						Type: "text",
						Text: fmt.Sprintf("Error creating directory: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
		result = fmt.Sprintf("Directory created: %s", path)

	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			{
				Type: "text",
				Text: result,
			},
		},
		IsError: false,
	}, nil
}

// Resource handlers

func handleConfigResource(ctx context.Context, uri string) (*protocol.ReadResourceResponse, error) {
	// Mock configuration data
	config := map[string]interface{}{
		"app_name":        "Advanced MCP Server",
		"version":         "1.0.0",
		"debug":           true,
		"max_connections": 100,
		"features": []string{
			"middleware",
			"resources",
			"prompts",
			"tools",
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	return &protocol.ReadResourceResponse{
		Contents: []protocol.ResourceContent{
			{
				URI:      uri,
				MimeType: "application/json",
				Text:     string(data),
			},
		},
	}, nil
}

func handleDocsResource(ctx context.Context, uri string) (*protocol.ReadResourceResponse, error) {
	// Mock documentation
	docs := `# Advanced MCP Server

This is an advanced example of an MCP server that demonstrates:

## Features

- **Middleware Integration**: Request/response logging and metrics
- **Resource Management**: File system and configuration access
- **Prompt Management**: Dynamic prompt generation
- **Advanced Tools**: Complex text processing and file operations

## Usage

Connect to this server using any MCP client to access:

1. Tools for text transformation and file operations
2. Resources for configuration and documentation
3. Prompts for code review and documentation generation

## Examples

### Text Transform Tool

` + "```json" + `
{
  "tool": "advanced_text_transform",
  "parameters": {
    "text": "hello world",
    "operations": ["uppercase", "reverse"]
  }
}
` + "```" + `

### File Operations Tool

` + "```json" + `
{
  "tool": "file_operations",
  "parameters": {
    "operation": "list",
    "path": "/tmp"
  }
}
` + "```" + `
`

	return &protocol.ReadResourceResponse{
		Contents: []protocol.ResourceContent{
			{
				URI:      uri,
				MimeType: "text/markdown",
				Text:     docs,
			},
		},
	}, nil
}

// Prompt handlers

func handleCodeReviewPrompt(ctx context.Context, name string, args map[string]interface{}) (*protocol.GetPromptResponse, error) {
	language, ok := args["language"].(string)
	if !ok {
		return nil, fmt.Errorf("language argument is required")
	}

	complexity, _ := args["complexity"].(string)
	if complexity == "" {
		complexity = "medium"
	}

	prompt := fmt.Sprintf(`You are a senior software engineer conducting a code review for %s code.

Context:
- Programming Language: %s
- Code Complexity: %s
- Review Focus: Best practices, security, performance, maintainability

Please review the following code and provide:
1. Overall assessment
2. Specific issues and suggestions
3. Positive aspects worth highlighting
4. Recommendations for improvement

Code to review:`, language, language, complexity)

	return &protocol.GetPromptResponse{
		Description: fmt.Sprintf("Code review prompt for %s (complexity: %s)", language, complexity),
		Messages: []protocol.PromptMessage{
			{
				Role: protocol.RoleUser,
				Content: protocol.Content{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}, nil
}

func handleDocumentationPrompt(ctx context.Context, name string, args map[string]interface{}) (*protocol.GetPromptResponse, error) {
	docType, ok := args["type"].(string)
	if !ok {
		return nil, fmt.Errorf("type argument is required")
	}

	audience, _ := args["audience"].(string)
	if audience == "" {
		audience = "developers"
	}

	var prompt string
	switch docType {
	case "api":
		prompt = fmt.Sprintf(`Create comprehensive API documentation for %s.

Include:
- Endpoint descriptions
- Parameter specifications
- Response formats
- Example requests/responses
- Error codes and handling
- Authentication requirements

Target audience: %s`, audience, audience)

	case "guide":
		prompt = fmt.Sprintf(`Create a step-by-step guide for %s.

Include:
- Clear objectives
- Prerequisites
- Detailed steps with examples
- Troubleshooting section
- Best practices
- Additional resources

Target audience: %s`, audience, audience)

	case "reference":
		prompt = fmt.Sprintf(`Create reference documentation for %s.

Include:
- Complete function/method listings
- Parameter descriptions
- Return value specifications
- Usage examples
- Cross-references
- Version compatibility notes

Target audience: %s`, audience, audience)

	default:
		return nil, fmt.Errorf("unknown documentation type: %s", docType)
	}

	return &protocol.GetPromptResponse{
		Description: fmt.Sprintf("%s documentation prompt for %s", docType, audience),
		Messages: []protocol.PromptMessage{
			{
				Role: protocol.RoleUser,
				Content: protocol.Content{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}, nil
}

// Supporting types

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

type simpleMetrics struct{}

func (m *simpleMetrics) IncrementCounter(name string, tags map[string]string) {
	log.Printf("[METRICS] Counter %s incremented with tags: %v", name, tags)
}

func (m *simpleMetrics) RecordDuration(name string, duration time.Duration, tags map[string]string) {
	log.Printf("[METRICS] Duration %s: %v with tags: %v", name, duration, tags)
}

func (m *simpleMetrics) RecordGauge(name string, value float64, tags map[string]string) {
	log.Printf("[METRICS] Gauge %s: %f with tags: %v", name, value, tags)
}
