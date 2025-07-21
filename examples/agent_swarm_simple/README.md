# agent_swarm_simple

## 🧠 What It Does

This example demonstrates a **LLM-powered agent swarm** with intelligent task routing and processing capabilities. It showcases how multiple AI agents can coordinate and hand off tasks using Ollama for local LLM reasoning.

## ⚙️ Requirements

- **Ollama running** - Local or remote Ollama server 
- **Compatible model** - `llama3.2`, `qwen2.5`, or similar
- **Go 1.21+** - For building and running
- **4GB+ RAM** - For the AI model

## 🚀 How to Run

```bash
# 1. Start Ollama (if not already running)
ollama serve

# 2. Pull a model
ollama pull llama3.2

# 3. Configure connection (optional)
export OLLAMA_URL="http://localhost:11434"    # Default
export OLLAMA_MODEL="llama3.2"                # Default

# 4. Run the swarm demo
go run main.go
```

## 🔍 Agents Used

- **Router** — Analyzes requests and routes to appropriate specialists using LLM reasoning
- **TextProcessor** — Handles text processing, formatting, and manipulation tasks  
- **Analyst** — Performs data analysis and generates insights

## 💡 Sample Output

```bash
� Starting LLM-Powered Simple Agent Swarm Demo

=== Scenario 1: Text Processing Task ===
� User: "Convert 'Agent Swarm Integration' to uppercase and remember it"

🧠 Router Agent (LLM Analysis):
- Task involves text transformation (uppercase)
- Memory storage required
- Routing to: TextProcessor

🔧 TextProcessor Agent:
- Using tool: uppercase
- Result: "AGENT SWARM INTEGRATION"
- Using tool: remember
- Stored in memory successfully

✅ Final Result: Text converted and stored
```

## 🧪 Test Scenarios

The demo runs four scenarios automatically:

### 1. Text Processing Task
```
Convert 'Agent Swarm Integration' to uppercase and remember it
→ Router → TextProcessor → uppercase + remember tools
```

### 2. Text Analysis Task  
```
Count words in 'AI Agent Coordination Systems' and analyze structure
→ Router → Analyst → word_count + analysis tools
```

### 3. Data Analysis Task
```
Analyze user engagement data and store insights
→ Router → Analyst → data analysis + memory storage
```

### 4. Multi-Step Processing
```
Process customer feedback and generate summary report
→ Router → TextProcessor → Analyst → Combined processing
```

## Expected Output

The demo will show:
- LLM-powered agent reasoning and routing decisions
- Tool usage based on natural language understanding
- Agent handoffs and collaboration
- Execution metrics (turns, tool calls, handoffs)
- Transparent reasoning for all decisions

## Error Handling

**Important**: This example uses pure LLM reasoning with **no rule-based fallback**. If Ollama is unavailable or the LLM fails, the example will error out as intended.

## Customization

You can customize this example by:
- Adding new specialized agents for your domain
- Implementing custom tools for specific tasks
- Modifying agent instructions for different behaviors
- Experimenting with different Ollama models

## Related Examples

- **agent_swarm**: Full-featured LLM-powered multi-agent workflow
- **agent_swarm_workflows**: Advanced workflow patterns with LLM orchestration
- **agent_swarm_llm**: Comprehensive LLM agent swarm demonstration

## LLM Integration Details

This example showcases:
- Natural language task analysis and routing
- Context-aware agent selection
- LLM-powered tool selection and usage
- Intelligent handoff decisions
- No rule-based logic - pure LLM reasoning

Perfect for understanding basic LLM-powered agent coordination!
