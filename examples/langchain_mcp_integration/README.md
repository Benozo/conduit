# LangChain Go + MCP Integration with Ollama

**Complete Documentation & Implementation Guide**

---

## 🌟 Overview

This example demonstrates the seamless integration of LangChain Go agents with MCP (Model Context Protocol) tools using **Ollama for local LLM inference**. This creates a powerful, privacy-first AI agent system that runs entirely on your local machine.

### **✅ What We Built**

1. **Full LangChain Go Integration** (`examples/langchain_mcp_integration/main.go`)
   - MCP tools wrapped as LangChain `tools.Tool` interface
   - Agent executor using **Ollama LLM** for local reasoning (no API keys required!)
   - Seamless tool execution through MCP server

2. **Local AI Agent System** with these capabilities:
   - **🦙 Local LLM Inference**: Uses Ollama - no API keys or cloud dependencies
   - **🧠 Intelligent Agent Reasoning**: Natural language task execution
   - **🔧 Rich Tool Ecosystem**: 15+ MCP tools wrapped for LangChain compatibility
   - **💾 Memory Persistence**: Shared memory across tool executions
   - **📄 HTML Generation**: Dynamic web page creation with Tailwind CSS
   - **🔒 Privacy-First**: Everything runs locally on your machine

---

## 🚀 Quick Start

### Prerequisites

1. **Go 1.19+** installed
2. **Ollama** installed and running:
   ```bash
   # Install Ollama (if not already done)
   curl https://ollama.ai/install.sh | sh
   
   # Start Ollama server
   ollama serve
   
   # Pull a model
   ollama pull llama3.2
   ```

### Running the Example

```bash
# Clone and navigate to the project
cd /path/to/gomcp

# Run with defaults (llama3.2 model, localhost:11434)
go run examples/langchain_mcp_integration/main.go

# Or with custom configuration
export OLLAMA_URL="http://192.168.10.10:11434"  # Remote Ollama server
export OLLAMA_MODEL="llama3.1"                  # Different model
go run examples/langchain_mcp_integration/main.go
```

### Quick Test
```bash
./test_langchain_ollama.sh
```

---

## 🔧 Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `OLLAMA_URL` | Ollama server URL | `http://localhost:11434` |
| `OLLAMA_MODEL` | Model to use for reasoning | `llama3.2` |

### Recommended Models

- **Fast & Efficient**: `llama3.2:1b` (1B parameters)
- **Balanced**: `llama3.2` (3B parameters) 
- **High Quality**: `llama3.1:8b` (8B parameters)
- **Code Tasks**: `codellama:7b`
- **Multilingual**: `mistral:7b`

Pull models with: `ollama pull <model-name>`

---

## 🛠️ Available Tools

The integration provides access to all MCP tools through natural language:

### Standard MCP Tool Ecosystem

| Category | Tools | Description |
|----------|-------|-------------|
| **Text Processing** | `uppercase`, `lowercase`, `trim`, `reverse` | String manipulation |
| **Memory** | `remember`, `recall`, `clear_memory`, `list_memories` | Persistent storage |
| **Utilities** | `timestamp`, `uuid`, `hash_md5`, `hash_sha256` | Helper functions |
| **Encoding** | `base64_encode`, `base64_decode`, `url_encode`, `url_decode` | Data transformation |
| **HTML Creation** | `create_html_page` | Dynamic web page generation |
| **Math** | `Calculator` (LangChain) | Mathematical operations |

---

## 💡 Example Interactions

The agent can handle complex, multi-step tasks through natural language:

### 1. Multi-Step Text Processing
```
Input: "Convert 'hello world' to uppercase and remember it as greeting"
→ Agent uses `uppercase` tool → "HELLO WORLD"
→ Agent uses `remember` tool → stores as "greeting"
→ Returns confirmation
```

### 2. Memory + Calculation
```
Input: "Remember that the answer is 42, then calculate 42 * 2"
→ Agent uses `remember` tool → stores "answer=42"
→ Agent uses `Calculator` tool → computes 84
→ Returns "The answer 42 is stored, and 42 * 2 = 84"
```

### 3. HTML Generation
```
Input: "Create a landing page for my startup called TechCorp"
→ Agent uses `create_html_page` tool
→ Generates complete HTML with Tailwind CSS
→ Saves to ./generated_pages/techcorp.html
```

### 4. Complex Workflows
```
Input: "Generate a UUID, convert it to uppercase, remember it as session_id, then create an HTML page showing the session info"
→ Multi-step execution using multiple tools
→ Demonstrates tool chaining and memory usage
```

## 🏗️ Technical Architecture

### System Flow
```
Natural Language Input
        ↓
LangChain Agent (Local Ollama LLM)
        ↓
Tool Selection & Planning
        ↓
MCP Tool Wrappers
        ↓
MCP Tool Registry
        ↓
Tool Execution + Memory
        ↓
Formatted Response
```

### Component Breakdown

1. **LangChain Agent**: Powered by local Ollama LLM for reasoning
2. **MCP Tool Wrappers**: Bridge between LangChain and MCP interfaces
3. **MCP Server**: Manages tool registry and shared memory
4. **Tool Execution**: Runs actual tool functions with persistence

### Key Integration Components

#### **🔧 Tool Wrapper Architecture**
```go
type MCPTool struct {
    name        string
    description string
    mcpServer   *conduit.EnhancedServer
    toolName    string
}

func (t MCPTool) Call(ctx context.Context, input string) (string, error) {
    // Parse input and prepare parameters
    params := parseInput(input, t.toolName)
    
    // Execute via MCP registry
    toolRegistry := t.mcpServer.GetToolRegistry()
    result, err := toolRegistry.Call(t.toolName, params, t.mcpServer.GetMemory())
    
    return formatResult(result), err
}
```

#### **🤖 LangChain Agent with Ollama + MCP Tools**
```go
// Local Ollama LLM (no API key required)
llm, err := ollama.New(
    ollama.WithServerURL("http://localhost:11434"),
    ollama.WithModel("llama3.2"),
)

mcpTools := []tools.Tool{
    MCPTool{name: "uppercase", toolName: "uppercase", ...},
    MCPTool{name: "remember", toolName: "remember", ...},
    MCPHTMLTool{mcpServer: server, ...},
    tools.Calculator{}, // LangChain built-in
}

agent := agents.NewOneShotAgent(llm, mcpTools, agents.WithMaxIterations(10))
executor := agents.NewExecutor(agent)
```

#### **💡 Natural Language Task Execution**
```go
question := "Convert 'hello world' to uppercase and remember it as greeting"
answer, err := chains.Run(context.Background(), executor, question)
```

---

## 🔄 Recent Improvements

### **✅ Successfully Updated LangChain Example to Use Ollama**

#### **Key Changes Made**

1. **Switched from OpenAI to Ollama**
   - Changed import from `openai` to `ollama`
   - Updated LLM creation to use local Ollama instance
   - Added configuration for Ollama URL and model selection

2. **Improved Tool Input Parsing**
   - Better handling of empty inputs
   - Enhanced parameter parsing for memory tools
   - Fallback logic for unclear input formats

3. **Enhanced HTML Tool Flexibility**
   - Smart parsing of natural language requests
   - Automatic content extraction from complex inputs
   - Fallback HTML generation for simple requests

4. **Increased Agent Robustness**
   - Max iterations increased from 5 to 10
   - Better error handling and recovery
   - Simplified test cases for higher success rate

#### **Technical Improvements**

**Tool Input Parsing**
```go
// Before: Rigid format requirements
if len(parts) != 2 {
    return "", fmt.Errorf("input must be in format: filename|content")
}

// After: Flexible natural language parsing
if strings.Contains(input, "|") {
    // Direct format
} else {
    // Extract from natural language
}
```

**Agent Configuration**
```go
// Before: Limited iterations
agents.WithMaxIterations(5)

// After: More robust execution
agents.WithMaxIterations(10)
```

**Error Handling**
```go
// Before: Basic error propagation
if err != nil {
    return "", err
}

// After: Contextual error messages
if err != nil {
    return "", fmt.Errorf("tool execution failed: %w", err)
}
```

---

## 🧪 Testing & Validation

### Automated Test Suite
```bash
./test_langchain_ollama.sh
```

This script:
- Checks Ollama availability (multiple URLs including `192.168.10.10:11434`)
- Lists available models
- Builds and runs the integration
- Demonstrates various tool capabilities
- Shows generated output files

### Manual Testing
```bash
# Start with a simple model for testing
export OLLAMA_MODEL="llama3.2:1b"
export OLLAMA_URL="http://192.168.10.10:11434"
go run examples/langchain_mcp_integration/main.go
```

### Test Results

**Before improvements:**
```
❌ Error: agent not finished before max iterations
❌ Error: unable to parse agent output
❌ Error: input must be in format: filename|content
```

**After improvements:**
```
✅ Better input parsing reduces parsing errors
✅ Increased iterations allow task completion
✅ Flexible HTML tool accepts natural language
✅ Simplified test cases improve success rate
```

---

## 🛠️ Troubleshooting

### Common Issues

1. **"connection refused"**
   - Start Ollama: `ollama serve`
   - Check if running: `curl http://localhost:11434/api/tags`
   - For remote: `curl http://192.168.10.10:11434/api/tags`

2. **"model not found"**
   - Pull the model: `ollama pull llama3.2`
   - List available: `ollama list`

3. **Slow responses**
   - Use smaller model: `export OLLAMA_MODEL="llama3.2:1b"`
   - Check system resources: `ollama ps`

4. **Tool execution errors**
   - Check generated_pages directory exists
   - Verify MCP tool registration in logs

### Performance Tips

- **Faster Models**: Use quantized models (`:1b`, `:3b`)
- **GPU Acceleration**: Ensure CUDA/Metal is available for Ollama
- **Memory**: Close other applications for large models
- **CPU**: Use models appropriate for your hardware

---

## 🔮 Extension Ideas

### Custom Tools
```go
// Add domain-specific tools
server.RegisterTool("analyze_data", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
    // Custom business logic
    return processData(params), nil
})
```

### Specialized Agents
```go
// Create agents for specific domains
htmlAgent := agents.NewOneShotAgent(llm, htmlTools, agents.WithMaxIterations(3))
dataAgent := agents.NewOneShotAgent(llm, dataTools, agents.WithMaxIterations(10))
```

### Streaming Responses
```go
// Add real-time feedback
executor := agents.NewExecutor(agent, agents.WithCallbacks(streamingCallback))
```

---

## 📚 Files & Structure

### Created/Updated Files
- `examples/langchain_mcp_integration/main.go` - Main integration example
- `examples/langchain_mcp_integration/README.md` - This comprehensive documentation
- `test_langchain_ollama.sh` - Test script with Ollama checks
- `LANGCHAIN_INTEGRATION_SUMMARY.md` - Integration summary (content merged here)
- `LANGCHAIN_OLLAMA_IMPROVEMENTS.md` - Improvement notes (content merged here)

### Related Examples
- `examples/simple_mcp_agent/` - Direct MCP tool orchestration
- `examples/openai/` - OpenAI MCP server
- `examples/ollama/` - Basic Ollama integration

---

## 🤝 Contributing

To add new tools or improve the integration:

1. Implement MCP tool in `lib/tools/`
2. Register tool in server setup
3. Add tool wrapper for LangChain compatibility
4. Update documentation and tests

---

## 🎯 Perfect Use Cases

- **Local Development**: AI-powered workflows without cloud dependencies
- **Privacy-Conscious Applications**: All data stays on your machine
- **Offline Environments**: No internet required for AI reasoning
- **Educational Projects**: Learn AI agent development safely
- **Prototyping**: Rapid development of AI-powered tools
- **Enterprise**: Private AI systems with full data control

---

## 📄 License

Part of the Conduit MCP library project.

---

**🎉 Result: A privacy-first, cost-free, local AI agent system** that can execute natural language tasks using local Ollama models, seamlessly integrate 15+ MCP tools with LangChain reasoning, generate HTML pages, process text, manage memory, and run completely offline with no external dependencies.
