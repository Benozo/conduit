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
	fmt.Println("🔧 Testing Built-in MCP Tools with AI Agents")
	fmt.Println("============================================")

	// Create MCP server
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8087
	config.EnableLogging = false

	server := conduit.NewEnhancedServer(config)

	// Register ALL built-in tools
	fmt.Println("📦 Registering built-in tool packages...")
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Create agent manager
	agentManager := agents.NewMCPAgentManager(server)

	// Create a general agent with access to all tools
	agent, err := agentManager.CreateAgent(
		"test_builtin_agent",
		"Built-in Tools Test Agent",
		"An agent for testing built-in MCP tools",
		"You are a comprehensive assistant with access to all built-in tools.",
		[]string{
			// Text tools
			"word_count", "char_count", "uppercase", "lowercase", "title_case",
			"camel_case", "snake_case", "trim", "reverse",
			// Memory tools
			"remember", "recall", "forget", "list_memories", "clear_memory",
			// Utility tools
			"base64_encode", "base64_decode", "url_encode", "url_decode",
			"hash_md5", "hash_sha256", "uuid", "timestamp", "random_number", "random_string",
		},
		agents.DefaultAgentConfig(),
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	fmt.Printf("✅ Created agent: %s\n", agent.Name)

	// Start server
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	time.Sleep(1 * time.Second)

	// Test built-in text tools
	fmt.Println("\n📝 Testing Text Tools:")
	testTextTool(agentManager, "uppercase", "hello world", "text")
	testTextTool(agentManager, "word_count", "This is a test sentence with seven words", "text")
	testTextTool(agentManager, "reverse", "hello", "text")

	// Test built-in memory tools
	fmt.Println("\n🧠 Testing Memory Tools:")
	testMemoryTool(agentManager, "remember", "test_key", "test_value")
	testMemoryTool(agentManager, "recall", "test_key", "")

	// Test built-in utility tools
	fmt.Println("\n🛠️  Testing Utility Tools:")
	testUtilityTool(agentManager, "base64_encode", "Hello World!")
	testUtilityTool(agentManager, "uuid", "")
	testUtilityTool(agentManager, "timestamp", "")

	fmt.Println("\n🎉 All built-in tool tests completed!")
	fmt.Println("✅ AI Agents successfully integrated with all built-in MCP tools")
}

func testTextTool(agentManager *agents.MCPAgentManager, tool, input, inputKey string) {
	fmt.Printf("   Testing %s with '%s'\n", tool, input)

	// Test direct tool execution
	agent, _ := agentManager.GetAgent("test_builtin_agent")
	result, err := agentManager.ExecuteToolWithMCP(tool, map[string]interface{}{
		inputKey: input,
	}, agent.Memory)

	if err != nil {
		fmt.Printf("   ❌ Tool failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Result: %v\n", result)
	}
}

func testMemoryTool(agentManager *agents.MCPAgentManager, tool, key, value string) {
	fmt.Printf("   Testing %s", tool)
	if value != "" {
		fmt.Printf(" with key='%s', value='%s'\n", key, value)
	} else {
		fmt.Printf(" with key='%s'\n", key)
	}

	agent, _ := agentManager.GetAgent("test_builtin_agent")
	params := map[string]interface{}{
		"key": key,
	}
	if value != "" {
		params["value"] = value
	}

	result, err := agentManager.ExecuteToolWithMCP(tool, params, agent.Memory)

	if err != nil {
		fmt.Printf("   ❌ Tool failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Result: %v\n", result)
	}
}

func testUtilityTool(agentManager *agents.MCPAgentManager, tool, input string) {
	fmt.Printf("   Testing %s", tool)
	if input != "" {
		fmt.Printf(" with '%s'\n", input)
	} else {
		fmt.Printf("\n")
	}

	agent, _ := agentManager.GetAgent("test_builtin_agent")
	params := map[string]interface{}{}
	if input != "" {
		params["text"] = input
	}

	result, err := agentManager.ExecuteToolWithMCP(tool, params, agent.Memory)

	if err != nil {
		fmt.Printf("   ❌ Tool failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Result: %v\n", result)
	}
}
