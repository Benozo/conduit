# multi_llm_swarm

## 🧠 What It Does

This example demonstrates the **Multi-LLM Agent Swarm** architecture where each agent uses its own specialized LLM provider and model. It showcases optimal model selection for different tasks while managing costs and performance across multiple AI providers.

## ⚙️ Requirements

**Local Models (Optional):**
- **Ollama** - For fast, local agents
- **Models** - `llama3.2`, `qwen2.5` pulled locally

**Cloud APIs (Optional):**
- **OpenAI API Key** - For premium reasoning tasks
- **DeepInfra API Key** - For specialized code generation
- **Go 1.21+** - For running the swarm

## 🚀 How to Run

```bash
# 1. Setup local models (optional)
ollama serve
ollama pull llama3.2    # For coordinator 
ollama pull qwen2.5     # For content creation

# 2. Set API keys (optional)
export OPENAI_API_KEY="sk-your-openai-key"
export DEEPINFRA_API_KEY="your-deepinfra-key"

# 3. Run the multi-LLM swarm demo
go run main.go
```

## 🔍 Agent Architecture

| Agent | Provider | Model | Purpose | Cost |
|-------|----------|-------|---------|------|
| **Coordinator** | Ollama | `llama3.2` | Fast task routing | Free |
| **ContentCreator** | Ollama | `qwen2.5` | Content generation | Free |
| **DataAnalyst** | OpenAI | `gpt-4` | Complex reasoning | Premium |
| **CodeGenerator** | DeepInfra | `Qwen Coder` | Code specialization | Mid-tier |

## 💡 Sample Output

```bash
🤖 Multi-LLM Agent Swarm Demo
============================

🔧 Creating Specialized Agents:
✅ Coordinator (Ollama llama3.2) - Task routing
✅ ContentCreator (Ollama qwen2.5) - Text processing  
✅ DataAnalyst (OpenAI gpt-4) - Complex analysis
✅ CodeGenerator (DeepInfra Qwen Coder) - Code generation

🎯 Task: "Analyze user data and generate Python code for visualization"

📍 Coordinator: Routing to DataAnalyst for analysis...
🧠 DataAnalyst (GPT-4): Analyzing patterns in user engagement data
📍 Coordinator: Routing to CodeGenerator for implementation...  
💻 CodeGenerator (Qwen Coder): Generating optimized Python visualization code

✅ Result: Complete analysis + production-ready Python code
```

## 🧪 Demo Scenarios

### 1. Cost-Optimized Processing
```
Simple text formatting → Ollama (free)
Complex data analysis → OpenAI (premium)
Code generation → DeepInfra (specialized)
```

### 2. Performance-Tuned Pipeline
```
Fast routing → Local llama3.2
Content creation → Local qwen2.5  
Critical reasoning → Cloud GPT-4
```

### 3. Provider Redundancy
```
Primary: Ollama models (always available)
Fallback: Cloud APIs (when local unavailable)
Specialized: Task-specific cloud models
```
go run main.go

# Full demo with real models (requires setup above)
go run main.go
```

## 🏗️ Multi-LLM Architecture

### Agent Creation with Individual Models

```go
// Coordinator with fast local model for routing
coordinator := swarmClient.CreateAgentWithModel("coordinator",
    "Route tasks to appropriate agents", []string{},
    &swarm.ModelConfig{
        Provider:    "ollama",
        Model:       "llama3.2", 
        URL:         "http://localhost:11434",
        Temperature: 0.7,
    })

// Analyst with premium model for complex reasoning
analyst := swarmClient.CreateAgentWithModel("data_analyst", 
    "Perform complex analysis", []string{"analyze", "report"},
    &swarm.ModelConfig{
        Provider:    "openai",
        Model:       "gpt-4",
        APIKey:      os.Getenv("OPENAI_API_KEY"),
        Temperature: 0.3,
    })

// Code specialist with domain-specific model
coder := swarmClient.CreateAgentWithModel("code_generator",
    "Generate and review code", []string{"format", "validate"}, 
    &swarm.ModelConfig{
        Provider:    "deepinfra", 
        Model:       "Qwen/Qwen2.5-Coder-32B-Instruct",
        APIKey:      os.Getenv("DEEPINFRA_API_KEY"),
        Temperature: 0.1,
    })
```

### Backward Compatibility

```go
// Existing code continues to work unchanged
swarm := conduit.NewSwarmClient(mcpServer)
agent := swarm.CreateAgent("coordinator", "Route tasks", []string{})
// Uses swarm-level LLM (current behavior)
```

## 📊 Model Selection Strategy

| Agent Type | Provider | Model | Why? |
|------------|----------|-------|------|
| **Coordinator** | Ollama | llama3.2 | Fast routing decisions, always available |
| **ContentCreator** | Ollama | qwen2.5 | Better multilingual and content generation |
| **DataAnalyst** | OpenAI | GPT-4 | Premium reasoning for complex analysis |
| **CodeGenerator** | DeepInfra | Qwen Coder | Specialized for code generation tasks |

## 🎯 Usage Patterns

### 1. Cost-Optimized Setup
- **Simple tasks**: Local Ollama models (free)
- **Complex tasks**: Cloud models (pay per use)
- **Routing**: Always local for speed and reliability

### 2. Performance-Optimized Setup  
- **Speed critical**: Small fast models (llama3.2:1b)
- **Quality critical**: Large models (GPT-4, Claude)
- **Specialized**: Domain-specific models (Codellama, Qwen Coder)

### 3. Reliability Setup
- **Primary**: Preferred models for each agent type
- **Fallback**: Swarm-level model as backup
- **Redundancy**: Multiple providers for critical agents

## 🔧 Configuration Options

```go
type ModelConfig struct {
    Provider     string  // "ollama", "openai", "deepinfra"
    Model        string  // "llama3.2", "gpt-4", "qwen2.5" 
    URL          string  // Custom endpoints
    APIKey       string  // API authentication
    Temperature  float64 // Creativity level (0.0-1.0)
    MaxTokens    int     // Response length limit
    TopK         int     // Sampling diversity
}
```

## 📈 Expected Output

```
🚀 Multi-LLM Agent Swarm Demo
=============================

🤖 Creating agents with different LLM providers:
   📋 coordinator - Ollama llama3.2 (fast routing)
   ✍️  content_creator - Ollama qwen2.5 (optimized for content)
   📊 data_analyst - OpenAI GPT-4 (premium reasoning)
   💻 code_generator - DeepInfra Qwen Coder (code specialist)

🎯 Multi-LLM Demo Scenarios:

📝 Scenario 1: Text Processing
📄 Coordinator (llama3.2) routes to ContentCreator (qwen2.5)
💬 Request: Convert 'Hello Multi-LLM World' to snake_case format
🔄 Processing with ollama llama3.2...
✅ Would route to appropriate specialist agent
...
```

## 🚀 Integration in Your Applications

### Basic Multi-LLM Setup

```go
// Create swarm with no default LLM
swarmClient := swarm.NewSwarmClient(server, nil)

// Create agents with specific models
routerAgent := swarmClient.CreateAgentWithModel("router", 
    instructions, tools, &swarm.ModelConfig{
    Provider: "ollama", Model: "llama3.2", // Fast & free
})

reasoningAgent := swarmClient.CreateAgentWithModel("reasoner",
    instructions, tools, &swarm.ModelConfig{
    Provider: "openai", Model: "gpt-4", // Powerful & accurate
})

// Use normally - agents automatically use their configured models
response := swarmClient.Run(routerAgent, messages, contextVars)
```

## 🔗 Related Examples

- **agent_swarm_llm**: Single-LLM swarm with Ollama
- **agent_swarm**: Rule-based multi-agent coordination
- **ollama**: Basic Ollama LLM integration
- **model_integration**: Custom model integration patterns

---

**🎉 Result**: A flexible, cost-efficient, and high-performance multi-agent system that uses the right model for each task!
