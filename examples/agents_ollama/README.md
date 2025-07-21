# agents_ollama

## 🧠 What It Does

An MCP-compatible AI agent that uses [Ollama](https://ollama.com) to respond to user input with intelligent tool usage. The agent can perform math calculations, search operations, memory management, and text processing using natural language understanding.

## ⚙️ Requirements

- **Ollama installed** - Download from [ollama.com](https://ollama.com)
- **AI model pulled** - `ollama pull llama3.2` (recommended)
- **Go 1.21+** - For building and running
- **4GB+ RAM** - For running the AI model

## 🚀 How to Run

```bash
# 1. Install and start Ollama
ollama serve

# 2. Pull a compatible model
ollama pull llama3.2

# 3. Set environment variables (optional)
export OLLAMA_HOST=http://localhost:11434    # Default
export OLLAMA_MODEL=llama3.2                 # Default

# 4. Run the agent
go run main.go
```

## 🧪 Sample Prompts

Try these natural language requests with the agent:

- `"What is 23 * 12?"`
- `"Search for the capital of France"`
- `"Remember that my name is Alice"`
- `"Convert 'Hello World' to uppercase"`
- `"What's the square root of 144?"`

## ✅ Sample Output

```bash
🤖 Starting Ollama-powered AI Agent...
🔌 Connected to Ollama at http://localhost:11434
🧠 Using model: llama3.2

> What is 23 * 12?

🧠 Agent Analysis: User wants to calculate 23 * 12
🔧 Tool Selection: math_multiply
⚡ Executing: math_multiply(a=23, b=12)
✅ Result: 276

Agent: I calculated 23 × 12 for you. The answer is 276.

> Search for the capital of France

🧠 Agent Analysis: User wants to search for geographical information
🔧 Tool Selection: search
⚡ Executing: search(query="capital of France")
✅ Result: Paris is the capital and largest city of France

Agent: I found that Paris is the capital of France. It's also the largest city in the country.

> Remember that my name is Alice

🧠 Agent Analysis: User wants to store personal information
🔧 Tool Selection: remember
⚡ Executing: remember(key="user_name", value="Alice")
✅ Result: Information stored successfully

Agent: Got it! I'll remember that your name is Alice.
```

## 🔧 Tools Available

The agent has access to these MCP tools:

- **Math Tools**: `add`, `multiply`, `divide`, `subtract`
- **Search Tools**: `search`, `web_lookup` 
- **Memory Tools**: `remember`, `recall`, `forget`, `list_memories`
- **Text Tools**: `uppercase`, `lowercase`, `reverse`, `word_count`
- **Utility Tools**: `timestamp`, `uuid`, `hash_md5`

## ⚙️ Configuration Options

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `OLLAMA_HOST` | `http://localhost:11434` | Ollama server URL |
| `OLLAMA_MODEL` | `llama3.2` | Model to use for reasoning |
| `AGENT_PORT` | `8080` | HTTP server port (if applicable) |

## 🎯 Key Features

- ✅ **Natural Language Processing**: Understands conversational requests
- ✅ **Intelligent Tool Selection**: Automatically chooses the right tools
- ✅ **Context Awareness**: Maintains conversation context and memory
- ✅ **Local AI**: Privacy-focused local model execution
- ✅ **Real-time Responses**: Fast response times with local models
- ✅ **Error Handling**: Graceful handling of tool failures and retries

## 🔍 How It Works

1. **User Input** → Natural language request
2. **AI Analysis** → Ollama model analyzes intent and available tools  
3. **Tool Selection** → Agent chooses appropriate MCP tools
4. **Execution** → Tools execute with extracted parameters
5. **Response** → Natural language response with results

## ⚠️ Troubleshooting

**Ollama Connection Issues:**
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Start Ollama if not running
ollama serve
```

**Model Not Found:**
```bash
# List available models
ollama list

# Pull required model
ollama pull llama3.2
```

**Memory Issues:**
- Ensure at least 4GB RAM available
- Try smaller models: `ollama pull llama3.2:1b`
- Close other applications if needed

## 📚 Related Examples

- [`ollama/`](../ollama) - Basic Ollama integration without agents
- [`agents_deepinfra/`](../agents_deepinfra) - Cloud-based agent alternative
- [`agent_swarm_simple/`](../agent_swarm_simple) - Multi-agent coordination
- [`pure_library/`](../pure_library) - Library-only usage without agents

## 🚀 Next Steps

After trying this example:

1. Experiment with different models: `qwen2.5`, `codellama`, `mistral`
2. Try the multi-agent swarm: [`agent_swarm_llm/`](../agent_swarm_llm)
3. Add custom tools: [`custom_tools/`](../custom_tools)
4. Deploy in production: [`openai/`](../openai) for cloud reliability
