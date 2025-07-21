# LLM-Powered Agent Swarm Example

This example demonstrates the **LLM-powered Agent Swarm** with **Ollama integration**, where agents use real LLM reasoning for intelligent decision-making, task routing, and tool selection.

## 🌟 What This Example Shows

### **🧠 LLM-Powered Agent Reasoning**
- Agents use Ollama LLM for intelligent task analysis
- Smart routing decisions based on content understanding
- Context-aware tool selection and execution
- Natural language decision explanations

### **🎯 Intelligent Task Coordination**
- **Coordinator**: Routes tasks using LLM analysis 
- **ContentCreator**: Handles writing and research with LLM reasoning
- **DataAnalyst**: Performs analysis with intelligent insights
- **MemoryManager**: Manages information with context understanding

### **🔧 Advanced LLM Features**
- System prompts specialized for each agent role
- Conversation history context for better decisions
- JSON-structured decision making
- Fallback to rule-based logic for reliability

## 🚀 Quick Start

### Prerequisites

1. **Ollama installed and running**:
   ```bash
   # Install Ollama
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Start Ollama server
   ollama serve
   
   # Pull a model
   ollama pull llama3.2
   ```

2. **Verify Ollama is working**:
   ```bash
   curl http://localhost:11434/api/tags
   ```

### Running the Example

```bash
# Default configuration (llama3.2, localhost:11434)
go run examples/agent_swarm_llm/main.go

# Custom configuration
export OLLAMA_URL="http://192.168.10.10:11434"
export OLLAMA_MODEL="llama3.1"
go run examples/agent_swarm_llm/main.go
```

## 🎭 Demo Scenarios

The example runs 4 intelligent scenarios:

### 1. **Content Creation Request**
```
"I need to write an article about artificial intelligence in healthcare. 
Can you help me research and create this content?"
```
- **LLM Analysis**: Identifies content creation need
- **Smart Routing**: Routes to ContentCreator agent
- **Tool Selection**: Uses research_topic and write_article tools
- **Reasoning**: Explains why ContentCreator was chosen

### 2. **Data Analysis Request**
```
"I have a customer behavior dataset that needs analysis. 
Please analyze the data and generate a comprehensive report."
```
- **LLM Analysis**: Recognizes data analysis requirements
- **Smart Routing**: Routes to DataAnalyst agent
- **Tool Usage**: Executes analyze_data and generate_report
- **Insights**: Provides analytical reasoning

### 3. **Memory Management Request**
```
"Please remember that our Q4 project deadline is December 15th 
and we're currently 75% complete."
```
- **LLM Analysis**: Identifies information storage need
- **Smart Routing**: Routes to MemoryManager agent
- **Context Storage**: Uses store_context for persistence
- **Organization**: Structures information intelligently

### 4. **Complex Multi-Agent Request**
```
"I need to research AI trends, analyze market data, and remember 
the key insights for our strategic planning meeting."
```
- **LLM Coordination**: Orchestrates multi-step workflow
- **Agent Handoffs**: Coordinates between multiple agents
- **Task Decomposition**: Breaks complex request into steps
- **Integration**: Combines results intelligently

## 🏗️ LLM Integration Architecture

### **System Prompt Engineering**
Each agent receives a specialized system prompt:
```go
systemPrompt := fmt.Sprintf(`You are %s.

%s

Your available tools: %s
Available agents for handoff: %s

Analyze the user message and decide your action...`, 
    agent.Name, agent.Instructions, tools, agents)
```

### **Decision Structure**
LLM responds with structured JSON:
```json
{
  "action": "tool_use|handoff|respond",
  "reasoning": "explanation of decision",
  "tool_name": "specific tool to use",
  "tool_args": {"param": "value"},
  "handoff_agent": "target agent name",
  "response": "direct response text"
}
```

### **Conversation Context**
- Maintains conversation history
- Includes context variables
- Preserves agent state across turns
- Enables contextual decision-making

## 🔄 LLM vs Rule-Based Comparison

| Feature | Rule-Based (Previous) | LLM-Powered (This Example) |
|---------|----------------------|----------------------------|
| **Decision Making** | `if contains("article")` | LLM analyzes full context |
| **Tool Selection** | Pattern matching | Semantic understanding |
| **Agent Routing** | Keyword triggers | Intent recognition |
| **Adaptability** | Fixed rules | Context-aware reasoning |
| **Explanations** | None | Reasoning provided |
| **Complexity** | Simple patterns | Natural language understanding |

## 📊 Expected Output

```
🧠 LLM-Powered Agent Swarm with Ollama
=====================================

🦙 Ollama URL: http://localhost:11434
🤖 Model: llama3.2

🎯 Agent Swarm Created with LLM Intelligence:
   📋 Coordinator - Routes tasks to appropriate specialists
   ✍️  ContentCreator - Handles content creation and text processing
   📊 DataAnalyst - Performs data analysis and reporting
   🧠 MemoryManager - Manages information storage and retrieval

🚀 Running LLM-Powered Demo Scenarios:
=====================================

📝 Scenario 1: Content Creation Request
📄 Description: Tests LLM routing to ContentCreator and tool usage
💬 Request: I need to write an article about artificial intelligence...
🔄 LLM Processing...
📊 Response:
✅ Success! Turns: 3, Tool calls: 2, Handoffs: 1
   🤖 Transferring to ContentCreator
   🤖 I used research_topic: 🔍 Research completed for topic: artificial intelligence in healthcare...
   🤖 I used write_article: 📝 Article 'AI in Healthcare Guide' about artificial intelligence in healthcare...
🎯 Final Agent: ContentCreator
⏱️  Execution Time: 2.1s
```

## 🔧 Configuration Options

### Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `OLLAMA_URL` | `http://localhost:11434` | Ollama server URL |
| `OLLAMA_MODEL` | `llama3.2` | Model for agent reasoning |

### Recommended Models
- **Fast**: `llama3.2:1b` (Quick responses)
- **Balanced**: `llama3.2` (Good quality)
- **High Quality**: `llama3.1:8b` (Best reasoning)
- **Code-focused**: `codellama:7b` (Technical tasks)

## 🛠️ Troubleshooting

### Ollama Connection Issues
```bash
# Check Ollama status
curl http://localhost:11434/api/tags

# Restart if needed
ollama serve

# Verify model is available
ollama list
ollama pull llama3.2
```

### LLM Response Issues
- If LLM responses are inconsistent, try a different model
- Increase model size for better reasoning (e.g., `llama3.1:8b`)
- Check Ollama logs for model loading issues

### Performance Optimization
```bash
# Use faster model for development
export OLLAMA_MODEL="llama3.2:1b"

# Use larger model for production
export OLLAMA_MODEL="llama3.1:8b"
```

## 🚀 Integration in Your Applications

### Basic Integration
```go
// Create Ollama model
ollamaModel := conduit.CreateOllamaModel("http://localhost:11434")

// Create LLM-powered swarm
swarmClient := swarm.NewSwarmClientWithLLM(
    mcpServer, 
    swarm.DefaultSwarmConfig(), 
    ollamaModel, 
    "llama3.2",
)

// Create intelligent agents
agent := swarmClient.CreateAgent(
    "MyAgent",
    "Specialized instructions for your domain...",
    []string{"tool1", "tool2"},
)

// Run with LLM reasoning
response := swarmClient.Run(agent, messages, contextVars)
```

### Custom Agent Prompts
```go
agent := swarmClient.CreateAgent(
    "CustomerSupportAgent",
    `You are a customer support specialist with expertise in:
    - Product troubleshooting
    - Account management  
    - Escalation procedures
    
    Always be helpful, empathetic, and solution-focused.
    Use tools to resolve issues efficiently.`,
    []string{"lookup_account", "create_ticket", "escalate"},
)
```

## 📈 Performance Characteristics

- **Startup Time**: 2-3 seconds (model loading)
- **Response Time**: 1-5 seconds per turn (model dependent)
- **Memory Usage**: Depends on Ollama model size
- **Accuracy**: High context understanding and task routing

## 🔮 Advanced Use Cases

1. **Multi-Domain Agents**: Specialized agents for different business domains
2. **Workflow Orchestration**: Complex multi-step business processes
3. **Customer Service**: Intelligent routing and escalation
4. **Content Pipeline**: Research → Writing → Review → Publishing
5. **Data Processing**: Collection → Analysis → Reporting → Storage

## 🤝 Related Examples

- [`examples/agent_swarm/`](../agent_swarm/) - Rule-based agent swarm
- [`examples/agent_swarm_simple/`](../agent_swarm_simple/) - Basic swarm concepts
- [`examples/agent_swarm_workflows/`](../agent_swarm_workflows/) - Advanced workflow patterns
- [`examples/ollama/`](../ollama/) - Basic Ollama integration
- [`examples/agents_ollama/`](../agents_ollama/) - LLM agents with tools

---

**🎉 Result**: A production-ready, LLM-powered agent swarm that provides intelligent task routing, context-aware decision making, and natural language reasoning - all running locally with Ollama!
