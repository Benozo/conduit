package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	fmt.Println("🤖 AI Agents Example using MCP")
	fmt.Println("===============================")

	// Create MCP server configuration
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8084
	config.EnableLogging = true

	// Create the MCP server
	server := conduit.NewEnhancedServer(config)

	// Register all available tools
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Register custom calculation tools
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

	// Create the AI Agent Manager
	agentManager := agents.NewMCPAgentManager(server)

	// Create specialized agents
	fmt.Println("📝 Creating specialized agents...")
	if err := agentManager.CreateSpecializedAgents(); err != nil {
		log.Fatalf("Failed to create agents: %v", err)
	}

	// List all created agents
	fmt.Println("\n🤖 Available Agents:")
	for _, agent := range agentManager.ListAgents() {
		fmt.Printf("  - %s (%s): %s\n", agent.Name, agent.ID, agent.Description)
	}

	// Example 1: Math calculation task
	fmt.Println("\n🧮 Example 1: Mathematical Calculation")
	mathTask, err := agentManager.CreateTaskForAgent("math_agent", agents.TaskTypeMath, map[string]interface{}{
		"a":         25.0,
		"b":         15.0,
		"operation": "add",
	})
	if err != nil {
		log.Fatalf("Failed to create math task: %v", err)
	}

	fmt.Printf("Created task: %s\n", mathTask.Title)

	// Start MCP server in background
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(2 * time.Second)

	// Execute the math task
	fmt.Println("Executing math task...")
	if err := agentManager.ExecuteTask(mathTask.ID); err != nil {
		log.Printf("Math task execution failed: %v", err)
	} else {
		fmt.Printf("✅ Math task completed successfully!\n")
		printTaskResult(mathTask)
	}

	// Example 2: Text processing task
	fmt.Println("\n📝 Example 2: Text Processing")
	textTask, err := agentManager.CreateTaskForAgent("text_agent", agents.TaskTypeTextProcessing, map[string]interface{}{
		"query": "Hello World! This is a sample text for processing and analysis.",
	})
	if err != nil {
		log.Fatalf("Failed to create text task: %v", err)
	}

	fmt.Printf("Created task: %s\n", textTask.Title)
	fmt.Println("Executing text processing task...")
	if err := agentManager.ExecuteTask(textTask.ID); err != nil {
		log.Printf("Text task execution failed: %v", err)
	} else {
		fmt.Printf("✅ Text task completed successfully!\n")
		printTaskResult(textTask)
	}

	// Example 3: Memory management task
	fmt.Println("\n🧠 Example 3: Memory Management")
	memoryTask, err := agentManager.CreateTaskForAgent("memory_agent", agents.TaskTypeMemoryManagement, map[string]interface{}{
		"action": "store",
		"key":    "user_preference",
		"value":  "dark_mode_enabled",
	})
	if err != nil {
		log.Fatalf("Failed to create memory task: %v", err)
	}

	fmt.Printf("Created task: %s\n", memoryTask.Title)
	fmt.Println("Executing memory management task...")
	if err := agentManager.ExecuteTask(memoryTask.ID); err != nil {
		log.Printf("Memory task execution failed: %v", err)
	} else {
		fmt.Printf("✅ Memory task completed successfully!\n")
		printTaskResult(memoryTask)
	}

	// Example 4: General purpose task
	fmt.Println("\n🌟 Example 4: General Purpose Agent")
	generalTask, err := agentManager.CreateTaskForAgent("general_agent", agents.TaskTypeGeneral, map[string]interface{}{
		"query": "Generate a UUID and timestamp for a new session",
	})
	if err != nil {
		log.Fatalf("Failed to create general task: %v", err)
	}

	fmt.Printf("Created task: %s\n", generalTask.Title)
	fmt.Println("Executing general purpose task...")
	if err := agentManager.ExecuteTask(generalTask.ID); err != nil {
		log.Printf("General task execution failed: %v", err)
	} else {
		fmt.Printf("✅ General task completed successfully!\n")
		printTaskResult(generalTask)
	}

	// Show agent templates
	fmt.Println("\n📋 Available Agent Templates:")
	templates := agents.GetAgentTemplates()
	for _, template := range templates {
		fmt.Printf("  - %s (%s): %s\n", template.Name, template.Type, template.Description)
	}

	// Create a custom agent from template
	fmt.Println("\n🔧 Creating Custom Agent from Template")
	customAgent, err := agentManager.CreateAgentFromTemplate(templates[0], "custom_math_agent")
	if err != nil {
		log.Printf("Failed to create custom agent: %v", err)
	} else {
		fmt.Printf("✅ Created custom agent: %s\n", customAgent.Name)
	}

	// Show final status
	fmt.Println("\n📊 Final Status:")
	fmt.Printf("Total Agents: %d\n", len(agentManager.ListAgents()))
	fmt.Printf("Total Tasks: %d\n", len(agentManager.ListTasks()))

	// Show available tools
	fmt.Println("\n🛠️  Available Tools:")
	availableTools := agentManager.GetAvailableTools()
	for _, tool := range availableTools {
		fmt.Printf("  - %s\n", tool)
	}

	fmt.Println("\n🎉 AI Agents example completed!")
	fmt.Println("The MCP server is running on http://localhost:8084")
	fmt.Println("Press Ctrl+C to exit...")

	// Keep the program running if not running in example mode
	if len(os.Args) > 1 && os.Args[1] == "--interactive" {
		select {} // Block forever
	}
}

// printTaskResult prints the results of a completed task
func printTaskResult(task *agents.Task) {
	fmt.Printf("Task ID: %s\n", task.ID)
	fmt.Printf("Status: %s\n", task.Status)
	fmt.Printf("Progress: %.1f%%\n", task.Progress*100)
	fmt.Printf("Steps executed: %d\n", len(task.Steps))

	if task.Error != "" {
		fmt.Printf("❌ Error: %s\n", task.Error)
	}

	for i, step := range task.Steps {
		fmt.Printf("  Step %d: %s (%s)\n", i+1, step.Name, step.Status)
		if step.Output != nil {
			if result, ok := step.Output["result"]; ok {
				fmt.Printf("    Result: %v\n", result)
			}
			if output, ok := step.Output["output"]; ok {
				fmt.Printf("    Output: %v\n", output)
			}
		}
	}
}
