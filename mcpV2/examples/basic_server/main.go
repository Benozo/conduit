// Basic MCP Server Example
//
// This example demonstrates how to create a simple MCP server using the Go SDK.
// The server provides a few basic tools like text transformation and math operations.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/benozo/neuron-mcp/protocol"
	"github.com/benozo/neuron-mcp/server"
	"github.com/benozo/neuron-mcp/transport"
)

func main() {
	// Create server with basic configuration
	opts := &server.ServerOptions{
		Info: protocol.Implementation{
			Name:    "example-server",
			Version: "1.0.0",
		},
		Capabilities: protocol.ServerCapabilities{
			Tools: &protocol.ToolsCapability{
				ListChanged: false,
			},
		},
		Logger: &simpleLogger{},
	}

	srv := server.NewServer(opts)

	// Register tools
	if err := registerTools(srv); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// Create STDIO transport
	transport := transport.NewStdioTransport(&transport.StdioOptions{
		TransportOptions: &transport.TransportOptions{
			Debug:  true,
			Logger: &simpleLogger{},
		},
	})

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...", sig)
		cancel()
	}()

	// Start serving
	log.Println("Starting MCP server...")
	if err := srv.Serve(ctx, transport); err != nil && err != context.Canceled {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server shutdown complete")
}

// registerTools registers example tools with the server
func registerTools(srv *server.Server) error {
	// Text transformation tool
	textTool := &protocol.Tool{
		Name:        "text_transform",
		Description: "Transform text using various operations",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"text": {
					Type:        "string",
					Description: "The text to transform",
				},
				"operation": {
					Type:        "string",
					Description: "Transform operation: uppercase, lowercase, reverse",
					Enum:        []interface{}{"uppercase", "lowercase", "reverse"},
				},
			},
			Required: []string{"text", "operation"},
		},
	}

	if err := srv.RegisterTool(textTool, handleTextTransform); err != nil {
		return fmt.Errorf("failed to register text_transform tool: %w", err)
	}

	// Math calculator tool
	mathTool := &protocol.Tool{
		Name:        "calculator",
		Description: "Perform basic mathematical operations",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"operation": {
					Type:        "string",
					Description: "Math operation: add, subtract, multiply, divide",
					Enum:        []interface{}{"add", "subtract", "multiply", "divide"},
				},
				"a": {
					Type:        "number",
					Description: "First number",
				},
				"b": {
					Type:        "number",
					Description: "Second number",
				},
			},
			Required: []string{"operation", "a", "b"},
		},
	}

	if err := srv.RegisterTool(mathTool, handleCalculator); err != nil {
		return fmt.Errorf("failed to register calculator tool: %w", err)
	}

	// Echo tool for testing
	echoTool := &protocol.Tool{
		Name:        "echo",
		Description: "Echo back the provided message",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"message": {
					Type:        "string",
					Description: "Message to echo back",
				},
			},
			Required: []string{"message"},
		},
	}

	if err := srv.RegisterTool(echoTool, handleEcho); err != nil {
		return fmt.Errorf("failed to register echo tool: %w", err)
	}

	return nil
}

// handleTextTransform implements the text transformation tool
func handleTextTransform(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
	text, ok := params["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text parameter must be a string")
	}

	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter must be a string")
	}

	var result string
	switch operation {
	case "uppercase":
		result = strings.ToUpper(text)
	case "lowercase":
		result = strings.ToLower(text)
	case "reverse":
		result = reverseString(text)
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: result,
		}},
	}, nil
}

// handleCalculator implements the calculator tool
func handleCalculator(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter must be a string")
	}

	// Handle different number types
	a, err := getNumberParam(params, "a")
	if err != nil {
		return nil, err
	}

	b, err := getNumberParam(params, "b")
	if err != nil {
		return nil, err
	}

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
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: fmt.Sprintf("%.2f", result),
		}},
	}, nil
}

// handleEcho implements the echo tool
func handleEcho(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
	message, ok := params["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter must be a string")
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: message,
		}},
	}, nil
}

// Helper functions

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// getNumberParam extracts a number parameter handling different JSON number types
func getNumberParam(params map[string]interface{}, key string) (float64, error) {
	value, exists := params[key]
	if !exists {
		return 0, fmt.Errorf("%s parameter is required", key)
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		// Try to parse as number
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, nil
		}
		return 0, fmt.Errorf("%s parameter must be a number, got string: %s", key, v)
	default:
		return 0, fmt.Errorf("%s parameter must be a number, got %T", key, value)
	}
}

// simpleLogger implements a basic logger
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
