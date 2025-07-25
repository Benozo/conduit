package mcp

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/benozo/neuron-mcp/middleware"
	"github.com/benozo/neuron-mcp/protocol"
	"github.com/benozo/neuron-mcp/server"
)

// TestIntegration verifies that all components work together correctly
func TestIntegration(t *testing.T) {
	// Test all major components integrate correctly
	t.Run("MiddlewareIntegration", testMiddlewareIntegration)
	t.Run("ResourceManagement", testResourceManagement)
	t.Run("PromptManagement", testPromptManagement)
}

func testMiddlewareIntegration(t *testing.T) {
	// Create a server with middleware
	logger := &testLogger{}
	metrics := &testMetrics{}

	opts := &server.ServerOptions{
		Info: protocol.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		Capabilities: protocol.ServerCapabilities{
			Tools: &protocol.ToolsCapability{
				ListChanged: false,
			},
		},
		Logger: logger,
		Middleware: []middleware.Middleware{
			middleware.LoggingMiddleware(logger),
			middleware.MetricsMiddleware(metrics),
		},
	}

	srv := server.NewServer(opts)

	// Register a test tool
	tool := &protocol.Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"message": {
					Type:        "string",
					Description: "Test message",
				},
			},
			Required: []string{"message"},
		},
	}

	err := srv.RegisterTool(tool, func(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
		message := params["message"].(string)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				{
					Type: "text",
					Text: fmt.Sprintf("Echo: %s", message),
				},
			},
			IsError: false,
		}, nil
	})

	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	// Verify middleware was added
	chain := srv.GetMiddlewareChain()
	if chain == nil {
		t.Fatal("Middleware chain should not be nil")
	}

	// Test that we can add more middleware
	srv.AddMiddleware(func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			// Custom middleware logic
			return next(ctx, req)
		}
	})

	t.Log("✓ Middleware integration test passed")
}

func testResourceManagement(t *testing.T) {
	// Create server with resource capability
	opts := &server.ServerOptions{
		Info: protocol.Implementation{
			Name:    "resource-test-server",
			Version: "1.0.0",
		},
		Capabilities: protocol.ServerCapabilities{
			Resources: &protocol.ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
		},
		Logger: &testLogger{},
	}

	srv := server.NewServer(opts)

	// Register a test resource
	resource := &protocol.Resource{
		URI:         "test://example.txt",
		Name:        "Example Resource",
		Description: "A test resource",
		MimeType:    "text/plain",
	}

	err := srv.RegisterResource(resource, func(ctx context.Context, uri string) (*protocol.ReadResourceResponse, error) {
		return &protocol.ReadResourceResponse{
			Contents: []protocol.ResourceContent{
				{
					URI:      uri,
					MimeType: "text/plain",
					Text:     "This is test content",
				},
			},
		}, nil
	})

	if err != nil {
		t.Fatalf("Failed to register resource: %v", err)
	}

	// Test unregistering
	srv.UnregisterResource("test://example.txt")

	t.Log("✓ Resource management test passed")
}

func testPromptManagement(t *testing.T) {
	// Create server with prompt capability
	opts := &server.ServerOptions{
		Info: protocol.Implementation{
			Name:    "prompt-test-server",
			Version: "1.0.0",
		},
		Capabilities: protocol.ServerCapabilities{
			Prompts: &protocol.PromptsCapability{
				ListChanged: false,
			},
		},
		Logger: &testLogger{},
	}

	srv := server.NewServer(opts)

	// Register a test prompt
	prompt := &protocol.Prompt{
		Name:        "test_prompt",
		Description: "A test prompt",
		Arguments: []protocol.PromptArgument{
			{
				Name:        "subject",
				Description: "Subject for the prompt",
				Required:    true,
			},
		},
	}

	err := srv.RegisterPrompt(prompt, func(ctx context.Context, name string, args map[string]interface{}) (*protocol.GetPromptResponse, error) {
		subject := args["subject"].(string)
		return &protocol.GetPromptResponse{
			Description: fmt.Sprintf("Test prompt for %s", subject),
			Messages: []protocol.PromptMessage{
				{
					Role: protocol.RoleUser,
					Content: protocol.Content{
						Type: "text",
						Text: fmt.Sprintf("Please provide information about %s", subject),
					},
				},
			},
		}, nil
	})

	if err != nil {
		t.Fatalf("Failed to register prompt: %v", err)
	}

	// Test unregistering
	srv.UnregisterPrompt("test_prompt")

	t.Log("✓ Prompt management test passed")
}

// Test helper types
type testLogger struct{}

func (l *testLogger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] "+msg, args...)
}

func (l *testLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *testLogger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg, args...)
}

func (l *testLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

type testMetrics struct{}

func (m *testMetrics) IncrementCounter(name string, tags map[string]string) {
	log.Printf("[METRICS] Counter %s incremented", name)
}

func (m *testMetrics) RecordDuration(name string, duration time.Duration, tags map[string]string) {
	log.Printf("[METRICS] Duration %s: %v", name, duration)
}

func (m *testMetrics) RecordGauge(name string, value float64, tags map[string]string) {
	log.Printf("[METRICS] Gauge %s: %f", name, value)
}
