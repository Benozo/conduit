package main

import (
	"log"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/mcp"
)

func main() {
	log.Println("=== Enhanced Custom Tools Example ===")

	config := conduit.DefaultConfig()
	config.Port = 8082
	config.Mode = mcp.ModeHTTP

	// Use enhanced server for custom tool schemas
	server := conduit.NewEnhancedServer(config)

	// Math calculation tools with enhanced schemas
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			return map[string]interface{}{
				"result":    a + b,
				"operation": "addition",
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("add", "Add two numbers together", map[string]interface{}{
			"a": conduit.NumberParam("First number to add"),
			"b": conduit.NumberParam("Second number to add"),
		}, []string{"a", "b"}))

	server.RegisterToolWithSchema("multiply",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			return map[string]interface{}{
				"result":    a * b,
				"operation": "multiplication",
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number to multiply"),
			"b": conduit.NumberParam("Second number to multiply"),
		}, []string{"a", "b"}))

	// Data processing tools with enhanced schemas
	server.RegisterToolWithSchema("filter_array",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			arr := params["array"].([]interface{})
			threshold := params["threshold"].(float64)

			var filtered []interface{}
			for _, item := range arr {
				if val, ok := item.(float64); ok && val > threshold {
					filtered = append(filtered, item)
				}
			}

			return map[string]interface{}{
				"filtered":       filtered,
				"original_count": len(arr),
				"filtered_count": len(filtered),
				"threshold":      threshold,
			}, nil
		},
		conduit.CreateToolMetadata("filter_array", "Filter array elements above a threshold value", map[string]interface{}{
			"array":     conduit.ArrayParam("Array of numbers to filter", "number"),
			"threshold": conduit.NumberParam("Minimum value to include in results"),
		}, []string{"array", "threshold"}))

	// Advanced calculation tool
	server.RegisterToolWithSchema("calculate",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			operation := params["operation"].(string)
			a := params["a"].(float64)
			b := params["b"].(float64)

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
					return map[string]interface{}{
						"error":     "Division by zero",
						"operation": operation,
						"operands":  []float64{a, b},
					}, nil
				}
				result = a / b
			case "power":
				result = 1
				for i := 0; i < int(b); i++ {
					result *= a
				}
			default:
				return map[string]interface{}{
					"error":     "Unknown operation: " + operation,
					"supported": []string{"add", "subtract", "multiply", "divide", "power"},
				}, nil
			}

			return map[string]interface{}{
				"result":    result,
				"operation": operation,
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("calculate", "Perform mathematical operations on two numbers", map[string]interface{}{
			"operation": conduit.EnumParam("Mathematical operation to perform", []string{"add", "subtract", "multiply", "divide", "power"}),
			"a":         conduit.NumberParam("First operand"),
			"b":         conduit.NumberParam("Second operand"),
		}, []string{"operation", "a", "b"}))

	// File system tools (mock implementation)
	server.RegisterToolWithSchema("list_files",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			path := params["path"].(string)
			// Mock implementation - in production, implement proper security checks
			return map[string]interface{}{
				"files": []string{"example1.txt", "example2.txt", "config.json"},
				"path":  path,
				"count": 3,
				"note":  "This is a mock implementation for demonstration",
			}, nil
		},
		conduit.CreateToolMetadata("list_files", "List files in a directory (mock implementation)", map[string]interface{}{
			"path": conduit.StringParam("Directory path to list"),
		}, []string{"path"}))

	log.Printf("Starting enhanced custom tools server on port %d...", config.Port)
	log.Printf("Custom tools registered: %d", server.GetCustomToolCount())

	// List custom tools
	for _, tool := range server.ListCustomTools() {
		log.Printf("  - %s: %s", tool["name"], tool["description"])
	}

	log.Fatal(server.Start())
}
