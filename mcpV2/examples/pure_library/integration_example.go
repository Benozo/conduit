//go:build integration
// +build integration

// Integration Example
//
// This shows how to integrate the pure library into your own application
// Run with: go run -tags=integration integration_example.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/benozo/neuron-mcp/library"
	"github.com/benozo/neuron-mcp/protocol"
)

// MyApp demonstrates embedding MCP library functionality
type MyApp struct {
	registry *library.ComponentRegistry
}

// NewMyApp creates an application with embedded MCP functionality
func NewMyApp() *MyApp {
	app := &MyApp{
		registry: library.NewComponentRegistry(),
	}

	// Register custom tools
	app.setupTools()

	return app
}

// setupTools registers application-specific tools
func (app *MyApp) setupTools() {
	tools := app.registry.Tools()

	// Register a custom business logic tool
	tools.Register("format_user", func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
		name, ok := params["name"].(string)
		if !ok {
			return nil, fmt.Errorf("name must be a string")
		}

		email, ok := params["email"].(string)
		if !ok {
			return nil, fmt.Errorf("email must be a string")
		}

		formatted := fmt.Sprintf("User: %s <%s>", name, email)

		return &protocol.ToolResult{
			Content: []protocol.Content{{
				Type: "text",
				Text: formatted,
			}},
		}, nil
	})

	// Register a data processing tool
	tools.Register("process_data", func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
		data, ok := params["data"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("data must be an array")
		}

		// Process the data (example: count items)
		count := len(data)
		result := fmt.Sprintf("Processed %d items", count)

		return &protocol.ToolResult{
			Content: []protocol.Content{{
				Type: "text",
				Text: result,
			}},
		}, nil
	})
}

// ProcessUser demonstrates using tools for business logic
func (app *MyApp) ProcessUser(name, email string) (string, error) {
	result, err := app.registry.Tools().Call(context.Background(), "format_user", map[string]interface{}{
		"name":  name,
		"email": email,
	})
	if err != nil {
		return "", fmt.Errorf("failed to format user: %w", err)
	}

	return result.Content[0].Text, nil
}

// StoreUserData demonstrates using memory for caching
func (app *MyApp) StoreUserData(userID string, data map[string]interface{}) error {
	return app.registry.Memory().Set(fmt.Sprintf("user:%s", userID), data)
}

// GetUserData retrieves cached user data
func (app *MyApp) GetUserData(userID string) (map[string]interface{}, error) {
	value, err := app.registry.Memory().Get(fmt.Sprintf("user:%s", userID))
	if err != nil {
		return nil, err
	}

	userData, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data format")
	}

	return userData, nil
}

// GetStats returns application statistics
func (app *MyApp) GetStats() (map[string]interface{}, error) {
	memStats, err := app.registry.Memory().Stats()
	if err != nil {
		return nil, err
	}

	tools := app.registry.Tools().List()

	return map[string]interface{}{
		"tools_count":     len(tools),
		"memory_keys":     memStats.ActiveKeys,
		"memory_backend":  memStats.Backend,
		"available_tools": tools,
	}, nil
}

func main() {
	runIntegrationExample()
}

func runIntegrationExample() {
	// Create application with embedded MCP
	app := NewMyApp()

	// Demonstrate usage
	fmt.Println("=== Integration Example ===")

	// 1. Use tool for business logic
	formatted, err := app.ProcessUser("Alice Smith", "alice@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Formatted user: %s\n", formatted)

	// 2. Store and retrieve data
	userData := map[string]interface{}{
		"name":       "Alice Smith",
		"email":      "alice@example.com",
		"created_at": "2025-07-25",
		"active":     true,
	}

	err = app.StoreUserData("alice", userData)
	if err != nil {
		log.Fatal(err)
	}

	retrieved, err := app.GetUserData("alice")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved data: %+v\n", retrieved)

	// 3. Get application statistics
	stats, err := app.GetStats()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("App stats: %+v\n", stats)

	fmt.Println("Integration example completed!")
}
