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
	fmt.Println("🧪 Testing Real MCP Tool Integration")
	fmt.Println("====================================")

	// Create MCP server
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8085
	config.EnableLogging = true

	server := conduit.NewEnhancedServer(config)

	// Register tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Register our custom math tools
	server.RegisterToolWithSchema("add",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			a := params["a"].(float64)
			b := params["b"].(float64)
			result := a + b
			fmt.Printf("🔧 MCP Tool 'add' called: %.1f + %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
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
			result := a * b
			fmt.Printf("🔧 MCP Tool 'multiply' called: %.1f × %.1f = %.1f\n", a, b, result)
			return map[string]interface{}{
				"result":    result,
				"operation": "multiplication",
				"operands":  []float64{a, b},
			}, nil
		},
		conduit.CreateToolMetadata("multiply", "Multiply two numbers", map[string]interface{}{
			"a": conduit.NumberParam("First number to multiply"),
			"b": conduit.NumberParam("Second number to multiply"),
		}, []string{"a", "b"}))

	// Create agent manager
	agentManager := agents.NewMCPAgentManager(server)

	// Create a simple math agent
	agent, err := agentManager.CreateAgent(
		"test_math_agent",
		"Test Math Agent",
		"A test agent for verifying MCP tool integration",
		"You are a mathematical assistant that uses MCP tools.",
		[]string{"add", "multiply"},
		agents.DefaultAgentConfig(),
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	fmt.Printf("✅ Created agent: %s\n", agent.Name)

	// Start server in background
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test 1: Addition
	fmt.Println("\n🧪 Test 1: Testing Addition")
	addTask, err := agentManager.CreateTask(
		"test_math_agent",
		"Addition Test",
		"Test addition operation",
		map[string]interface{}{
			"a":         10.0,
			"b":         5.0,
			"operation": "add",
		},
	)
	if err != nil {
		log.Fatalf("Failed to create addition task: %v", err)
	}

	err = agentManager.ExecuteTask(addTask.ID)
	if err != nil {
		log.Printf("Addition task failed: %v", err)
	} else {
		fmt.Println("✅ Addition task completed!")
		for _, step := range addTask.Steps {
			if step.Name == "perform_calculation" && step.Output != nil {
				if output, ok := step.Output["output"]; ok {
					if result, ok := output.(map[string]interface{})["result"]; ok {
						fmt.Printf("   Result: %.1f\n", result)
					}
				}
			}
		}
	}

	// Test 2: Multiplication
	fmt.Println("\n🧪 Test 2: Testing Multiplication")
	multiplyTask, err := agentManager.CreateTask(
		"test_math_agent",
		"Multiplication Test",
		"Test multiplication operation",
		map[string]interface{}{
			"a":         7.0,
			"b":         8.0,
			"operation": "multiply",
		},
	)
	if err != nil {
		log.Fatalf("Failed to create multiplication task: %v", err)
	}

	err = agentManager.ExecuteTask(multiplyTask.ID)
	if err != nil {
		log.Printf("Multiplication task failed: %v", err)
	} else {
		fmt.Println("✅ Multiplication task completed!")
		for _, step := range multiplyTask.Steps {
			if step.Name == "perform_calculation" && step.Output != nil {
				if output, ok := step.Output["output"]; ok {
					if result, ok := output.(map[string]interface{})["result"]; ok {
						fmt.Printf("   Result: %.1f\n", result)
					}
				}
			}
		}
	}

	// Test 3: Test actual MCP tool calling directly
	fmt.Println("\n🧪 Test 3: Direct MCP Tool Calling")
	directResult, err := agentManager.ExecuteToolWithMCP("add", map[string]interface{}{
		"a": 100.0,
		"b": 200.0,
	}, agent.Memory)
	if err != nil {
		log.Printf("Direct tool call failed: %v", err)
	} else {
		fmt.Printf("✅ Direct MCP tool call successful: %v\n", directResult)
	}

	fmt.Println("\n🎉 All tests completed!")
	fmt.Println("The agents are successfully using real MCP tools!")
}
