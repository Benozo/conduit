// WebSocket Client Example
//
// This example demonstrates connecting to an MCP server using WebSocket transport
// and calling tools with progress tracking.
package main

import (
	"context"
	"log"
	"time"

	"github.com/benozo/neuron-mcp/client"
	"github.com/benozo/neuron-mcp/protocol"
	"github.com/benozo/neuron-mcp/transport"
)

func main() {
	// Create WebSocket transport
	wsTransport := transport.NewWebSocketTransport("ws://localhost:8082/mcp")

	// Create client with progress tracking
	mcpClient := client.NewClient(wsTransport, &client.ClientOptions{
		Timeout:        30 * time.Second,
		ConnectTimeout: 10 * time.Second,
		ClientInfo: protocol.Implementation{
			Name:    "websocket-mcp-client",
			Version: "1.0.0",
		},
		ProgressHandler: handleProgress,
	})

	ctx := context.Background()

	// Connect to server
	log.Println("Connecting to WebSocket MCP server...")
	err := wsTransport.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	// Perform MCP handshake
	capabilities := protocol.ClientCapabilities{
		Experimental: map[string]interface{}{
			"progressNotifications": true,
		},
	}

	err = mcpClient.Connect(ctx, capabilities)
	if err != nil {
		log.Fatalf("Failed to connect to MCP server: %v", err)
	}

	log.Println("Connected to MCP server successfully!")

	// Demonstrate client functionality
	demonstrateToolListing(ctx, mcpClient)
	demonstrateToolCalling(ctx, mcpClient)
	demonstrateProgressTracking(ctx, mcpClient)

	// Clean up
	mcpClient.Close()
	wsTransport.Close()

	log.Println("WebSocket client example completed!")
}

// handleProgress handles progress notifications
func handleProgress(token string, progress float64, total int64) {
	percentage := progress * 100
	log.Printf("Progress [%s]: %.1f%% (%d/%d)", token, percentage, int64(progress*float64(total)), total)
}

// demonstrateToolListing shows how to list available tools
func demonstrateToolListing(ctx context.Context, client *client.Client) {
	log.Println("\n=== Tool Listing ===")

	tools, err := client.ListTools(ctx)
	if err != nil {
		log.Printf("Failed to list tools: %v", err)
		return
	}

	log.Printf("Available tools (%d):", len(tools))
	for _, tool := range tools {
		log.Printf("  - %s: %s", tool.Name, tool.Description)
	}
}

// demonstrateToolCalling shows how to call tools
func demonstrateToolCalling(ctx context.Context, client *client.Client) {
	log.Println("\n=== Tool Calling ===")

	// Call text transformation tool
	result, err := client.CallTool(ctx, "text_transform", map[string]interface{}{
		"text":      "Hello, WebSocket!",
		"operation": "uppercase",
	})

	if err != nil {
		log.Printf("Tool call failed: %v", err)
		return
	}

	log.Printf("Text transform result: %s", result.Content[0].Text)

	// Call echo tool
	result, err = client.CallTool(ctx, "echo", map[string]interface{}{
		"message": "Testing WebSocket transport",
	})

	if err != nil {
		log.Printf("Echo tool call failed: %v", err)
		return
	}

	log.Printf("Echo result: %s", result.Content[0].Text)
}

// demonstrateProgressTracking shows progress tracking functionality
func demonstrateProgressTracking(ctx context.Context, client *client.Client) {
	log.Println("\n=== Progress Tracking ===")

	// Create a context with progress tracking
	progressCtx := context.WithValue(ctx, "progressToken", "demo_task_001")

	// Simulate a long-running task
	result, err := client.CallTool(progressCtx, "long_task", map[string]interface{}{
		"duration": 5, // 5 seconds
		"steps":    10,
	})

	if err != nil {
		log.Printf("Long task failed: %v", err)
		return
	}

	log.Printf("Long task completed: %s", result.Content[0].Text)
}

// mockLongRunningTask simulates a long-running task with progress updates
func mockLongRunningTask(ctx context.Context, duration int, steps int) error {
	stepDuration := time.Duration(duration) * time.Second / time.Duration(steps)

	for i := 0; i < steps; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Simulate work
		time.Sleep(stepDuration)

		// Report progress
		progress := float64(i+1) / float64(steps)
		if progressToken := ctx.Value("progressToken"); progressToken != nil {
			handleProgress(progressToken.(string), progress, int64(steps))
		}
	}

	return nil
}
