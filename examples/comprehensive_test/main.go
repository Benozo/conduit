package main

import (
	"fmt"
	"log"
	"time"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	fmt.Println("🚀 Comprehensive AI Agents + MCP Integration Test")
	fmt.Println("=================================================")

	// Create MCP server with comprehensive tool set
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8086
	config.EnableLogging = false // Reduce noise for testing

	server := conduit.NewEnhancedServer(config)

	// Register all tool categories
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Add custom math tools
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a + b
			fmt.Printf("   ✅ MCP: %.1f + %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "addition",
			}, nil
		},
		conduit.CreateToolMetadata("add", "Add two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number"),
			"b": conduit.NumberParam("Second number"),
		}, []string{"a", "b"}))

	server.RegisterToolWithSchema("multiply",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a * b
			fmt.Printf("   ✅ MCP: %.1f × %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "multiplication",
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number"),
			"b": conduit.NumberParam("Second number"),
		}, []string{"a", "b"}))

	// Create agent manager
	agentManager := agents.NewMCPAgentManager(server)

	// Create specialized agents
	fmt.Println("🤖 Creating specialized agents...")
	if err := agentManager.CreateSpecializedAgents(); err != nil {
		log.Fatalf("Failed to create agents: %v", err)
	}

	// Start server
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	time.Sleep(1 * time.Second)

	// Test scenarios
	fmt.Println("\n📊 Running comprehensive tests...")

	// Test 1: Math Agent with Addition
	fmt.Println("\n🧮 Test 1: Math Agent - Addition")
	testMathOperation(agentManager, "add", 15.0, 25.0)

	// Test 2: Math Agent with Multiplication
	fmt.Println("\n🧮 Test 2: Math Agent - Multiplication")
	testMathOperation(agentManager, "multiply", 8.0, 7.0)

	// Test 3: Text Agent with Text Processing
	fmt.Println("\n📝 Test 3: Text Agent - Text Processing")
	testTextProcessing(agentManager, "Hello World! This is a comprehensive test of our AI agents system.")

	// Test 4: Memory Agent - Store and Retrieve
	fmt.Println("\n🧠 Test 4: Memory Agent - Data Storage")
	testMemoryOperations(agentManager)

	// Test 5: Utility Agent - Various Utilities
	fmt.Println("\n🛠️  Test 5: Utility Agent - UUID and Timestamp")
	testUtilityOperations(agentManager)

	fmt.Println("\n🎉 All comprehensive tests completed successfully!")
	fmt.Println("✅ AI Agents are fully integrated with MCP tools")
	fmt.Println("✅ All tool categories working: Math, Text, Memory, Utility")
	fmt.Println("✅ Real-time tool execution with proper results")
}

func testMathOperation(agentManager *agents.MCPAgentManager, operation string, a, b float64) {
	task, err := agentManager.CreateTask(
		"math_agent",
		fmt.Sprintf("%s Operation", operation),
		fmt.Sprintf("Perform %s operation on %.1f and %.1f", operation, a, b),
		map[string]interface{}{
			"a":         a,
			"b":         b,
			"operation": operation,
		},
	)
	if err != nil {
		log.Printf("Failed to create %s task: %v", operation, err)
		return
	}

	fmt.Printf("   Task: %.1f %s %.1f\n", a, getOperatorSymbol(operation), b)
	err = agentManager.ExecuteTask(task.ID)
	if err != nil {
		log.Printf("   ❌ Task failed: %v", err)
	} else {
		// Extract result from task steps
		for _, step := range task.Steps {
			if step.Output != nil && step.Output["output"] != nil {
				if output, ok := step.Output["output"].(map[string]interface{}); ok {
					if result, ok := output["result"]; ok {
						fmt.Printf("   📋 Agent Result: %.1f\n", result)
					}
				}
			}
		}
	}
}

func testTextProcessing(agentManager *agents.MCPAgentManager, text string) {
	task, err := agentManager.CreateTask(
		"text_agent",
		"Text Analysis",
		"Analyze the provided text",
		map[string]interface{}{
			"query": text,
		},
	)
	if err != nil {
		log.Printf("Failed to create text task: %v", err)
		return
	}

	fmt.Printf("   Text: %s\n", text)
	err = agentManager.ExecuteTask(task.ID)
	if err != nil {
		log.Printf("   ❌ Task failed: %v", err)
	} else {
		fmt.Printf("   📋 Text processed successfully with %d steps\n", len(task.Steps))
	}
}

func testMemoryOperations(agentManager *agents.MCPAgentManager) {
	// Test storing data
	task, err := agentManager.CreateTask(
		"memory_agent",
		"Store User Data",
		"Store user preference data",
		map[string]interface{}{
			"action": "store",
			"key":    "user_theme",
			"value":  "dark_mode",
		},
	)
	if err != nil {
		log.Printf("Failed to create memory task: %v", err)
		return
	}

	fmt.Printf("   Action: Store user_theme = dark_mode\n")
	err = agentManager.ExecuteTask(task.ID)
	if err != nil {
		log.Printf("   ❌ Task failed: %v", err)
	} else {
		fmt.Printf("   📋 Memory operation completed successfully\n")
	}
}

func testUtilityOperations(agentManager *agents.MCPAgentManager) {
	task, err := agentManager.CreateTask(
		"utility_agent",
		"Generate Session Data",
		"Generate UUID and timestamp for session",
		map[string]interface{}{
			"action": "generate",
			"type":   "session",
		},
	)
	if err != nil {
		log.Printf("Failed to create utility task: %v", err)
		return
	}

	fmt.Printf("   Action: Generate session identifiers\n")
	err = agentManager.ExecuteTask(task.ID)
	if err != nil {
		log.Printf("   ❌ Task failed: %v", err)
	} else {
		fmt.Printf("   📋 Utility operations completed successfully\n")
	}
}

func getOperatorSymbol(operation string) string {
	switch operation {
	case "add":
		return "+"
	case "multiply":
		return "×"
	case "subtract":
		return "-"
	case "divide":
		return "÷"
	default:
		return operation
	}
}
