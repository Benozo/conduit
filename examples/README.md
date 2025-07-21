# Conduit Examples

This folder contains working examples for ConduitMCP including tool calling, LLM integration, agent frameworks, and protocol demonstrations. Each example is self-contained with clear setup instructions and sample outputs.

## ⚡ Quick Start

### 🎯 Try It Now (30 seconds)

```bash
# 1. Clone and enter examples
git clone https://github.com/benozo/conduit && cd conduit/examples

# 2. Pick your adventure:
cd pure_library && go run main.go     # 📚 Library demo (no deps)
cd stdio_example && go build .       # 🔌 MCP client ready  
cd ollama && go run main.go          # 🤖 Local AI (needs Ollama)
cd openai && OPENAI_API_KEY=sk-... go run main.go  # ☁️ Cloud AI
```

### 🎯 Choose Your Path

- **🆕 New to Conduit?** → [`pure_library`](./pure_library) • [`stdio_example`](./stdio_example)
- **🤖 Want Local AI?** → [`ollama`](./ollama) • [`agents_ollama`](./agents_ollama)  
- **☁️ Need Cloud AI?** → [`openai`](./openai) • [`agents_deepinfra`](./agents_deepinfra)
- **🐝 Multi-Agent Systems?** → [`agent_swarm_simple`](./agent_swarm_simple) • [`multi_llm_swarm`](./multi_llm_swarm)
- **🌐 Web Integration?** → [`sse_example`](./sse_example) • [`pure_library_web`](./pure_library_web)

## 📋 Example Gallery

| Example | Type | Description | Features | Setup Time |
|---------|------|-------------|----------|------------|
| **🔌 Protocol & Integration** |
| [`stdio_example`](./stdio_example) | ![MCP](https://img.shields.io/badge/MCP-Compatible-blue) | MCP stdio server for VS Code, Cline, etc. | stdio protocol, tool calling | ⚡ 1 min |
| [`sse_example`](./sse_example) | ![Web](https://img.shields.io/badge/Web-HTTP/SSE-green) | HTTP/SSE server for web applications | HTTP API, real-time streaming | ⚡ 2 min |
| [`pure_mcp`](./pure_mcp) | ![Core](https://img.shields.io/badge/Core-MCP-purple) | Raw MCP protocol implementation | Pure MCP, no server wrapper | ⚡ 1 min |
| **📚 Library Usage** |
| [`pure_library`](./pure_library) | ![Library](https://img.shields.io/badge/Library-Go-cyan) | Use Conduit as Go library | Library integration, custom tools | ⚡ 1 min |
| [`pure_library_cli`](./pure_library_cli) | ![CLI](https://img.shields.io/badge/CLI-Terminal-orange) | CLI tool with MCP components | Command-line interface | ⚡ 2 min |
| [`pure_library_web`](./pure_library_web) | ![Web](https://img.shields.io/badge/Web-Custom-green) | Custom web server with MCP | Web server, custom endpoints | ⚡ 3 min |
| [`embedded`](./embedded) | ![Embedded](https://img.shields.io/badge/Embedded-App-yellow) | Embed Conduit in existing apps | Application integration | ⚡ 3 min |
| **🤖 LLM Integration** |
| [`ollama`](./ollama) | ![Local](https://img.shields.io/badge/Local-Ollama-red) | Local LLM with tool calling | Ollama integration, auto tool selection | 🔥 5 min |
| [`openai`](./openai) | ![Cloud](https://img.shields.io/badge/Cloud-OpenAI-blue) | OpenAI GPT with tools | OpenAI API, cloud LLM | 🔥 3 min |
| [`model_integration`](./model_integration) | ![Custom](https://img.shields.io/badge/Custom-Model-purple) | Custom model patterns | Custom LLM integration | 🚀 10 min |
| **🛠️ Tool Development** |
| [`custom_tools`](./custom_tools) | ![Tools](https://img.shields.io/badge/Tools-Enhanced-green) | Enhanced tool registration | Rich schemas, validation | ⚡ 3 min |
| [`builtin_tools_test`](./builtin_tools_test) | ![Test](https://img.shields.io/badge/Test-Tools-gray) | Test all built-in tools | Tool testing, validation | ⚡ 1 min |
| **🤖 AI Agents** |
| [`ai_agents`](./ai_agents) | ![Agents](https://img.shields.io/badge/Agents-Framework-purple) | AI Agents with task management | Agent framework, task execution | 🚀 10 min |
| [`agents_test`](./agents_test) | ![Test](https://img.shields.io/badge/Test-Agents-gray) | Basic agent functionality | Agent testing | ⚡ 2 min |
| [`agents_ollama`](./agents_ollama) | ![Local+AI](https://img.shields.io/badge/Local-Agents-red) | Agents with Ollama LLM | AI agents, local LLM | 🔥 5 min |
| [`agents_deepinfra`](./agents_deepinfra) | ![Cloud+AI](https://img.shields.io/badge/Cloud-Agents-blue) | Agents with DeepInfra | AI agents, cloud inference | 🔥 5 min |
| [`agents_library_mode`](./agents_library_mode) | ![Library+AI](https://img.shields.io/badge/Library-Agents-cyan) | Library-mode agent usage | Agent patterns | ⚡ 3 min |
| [`agents_mock_llm`](./agents_mock_llm) | ![Mock](https://img.shields.io/badge/Mock-Testing-gray) | Mock LLM for testing | Testing, development | ⚡ 2 min |
| [`agents_vue_builder`](./agents_vue_builder) | ![Code](https://img.shields.io/badge/Code-Vue.js-green) | Vue.js app builder agent | Code generation, Vue.js | 🚀 15 min |
| **🐝 Agent Swarms** |
| [`agent_swarm_simple`](./agent_swarm_simple) | ![Swarm](https://img.shields.io/badge/Swarm-Simple-orange) | Simple swarm demo | Basic swarm patterns | 🔥 5 min |
| [`agent_swarm`](./agent_swarm) | ![Swarm](https://img.shields.io/badge/Swarm-Advanced-orange) | Basic agent coordination | Multi-agent, handoffs | 🚀 10 min |
| [`agent_swarm_llm`](./agent_swarm_llm) | ![Swarm+LLM](https://img.shields.io/badge/Swarm-LLM-red) | LLM-powered agent swarm | LLM coordination, Ollama | 🚀 15 min |
| [`agent_swarm_workflows`](./agent_swarm_workflows) | ![Workflows](https://img.shields.io/badge/Workflows-DAG-purple) | Advanced workflow patterns | DAG, Supervisor, Pipeline | 🚀 20 min |
| [`multi_llm_swarm`](./multi_llm_swarm) | ![Multi-LLM](https://img.shields.io/badge/Multi-LLM-Enterprise-gold) | Multi-LLM agent architecture | Multiple LLM providers | 🚀 20 min |
| **📖 RAG & Advanced** |
| [`rag`](./rag) | ![RAG](https://img.shields.io/badge/RAG-Documents-brown) | RAG with document processing | Document analysis, embeddings | 🚀 15 min |
| [`rag_chat_terminal`](./rag_chat_terminal) | ![RAG+CLI](https://img.shields.io/badge/RAG-Terminal-brown) | Terminal RAG chat interface | CLI RAG, interactive | 🚀 10 min |
| [`rag_real_world`](./rag_real_world) | ![RAG+Prod](https://img.shields.io/badge/RAG-Production-brown) | Production RAG patterns | Real-world RAG | 🚀 30 min |
| [`langchain_mcp_integration`](./langchain_mcp_integration) | ![LangChain](https://img.shields.io/badge/LangChain-Bridge-brown) | LangChain integration | LangChain + MCP | 🚀 15 min |

**Setup Time Legend:**
- ⚡ = 1-3 min (no external deps)  
- 🔥 = 5-10 min (needs API keys or local services)
- 🚀 = 15-30 min (complex setup or multiple services)

## 🎯 Use Case Guide

### I want to...

**Create an MCP server for VS Code/Cline:**
→ [`stdio_example`](./stdio_example)

**Build a web application with real-time features:**
→ [`sse_example`](./sse_example) → [`pure_library_web`](./pure_library_web)

**Integrate local AI with tool calling:**
→ [`ollama`](./ollama) → [`agents_ollama`](./agents_ollama)

**Create autonomous AI agents:**
→ [`ai_agents`](./ai_agents) → [`agent_swarm`](./agent_swarm)

**Build multi-agent workflows:**
→ [`agent_swarm_simple`](./agent_swarm_simple) → [`agent_swarm_workflows`](./agent_swarm_workflows)

**Process documents with AI:**
→ [`rag`](./rag) → [`rag_chat_terminal`](./rag_chat_terminal)

**Test different LLM providers:**
→ [`multi_llm_swarm`](./multi_llm_swarm)

**Develop custom tools:**
→ [`custom_tools`](./custom_tools) → [`builtin_tools_test`](./builtin_tools_test)

## 📖 Documentation Standards

Each example follows our standardized documentation format for consistency and usability:

### ✅ **Required Elements**
- **🧠 What It Does** - Clear purpose and use case
- **⚙️ Requirements** - Prerequisites and dependencies  
- **🚀 How to Run** - Step-by-step setup commands
- **✅ Sample Output** - Realistic terminal output examples

### 🎯 **Quality Standards** 
- Copy-pasteable commands that actually work
- Real terminal output (not pseudo-code)
- Troubleshooting for common issues
- Cross-references to related examples

### 📝 **Template Available**
New examples should use our standardized template: [`_README_TEMPLATE.md`](./_README_TEMPLATE.md)

### 📊 **Documentation Status**

| Category | Examples | ✅ Complete | 🔄 In Progress | ❌ Missing |
|----------|----------|-------------|---------------|------------|
| **Protocol** | 3 | 2 | 1 | 0 |
| **Library** | 4 | 2 | 1 | 1 |
| **LLM** | 3 | 2 | 1 | 0 |
| **Agents** | 7 | 3 | 2 | 2 |
| **Swarms** | 5 | 3 | 1 | 1 |
| **RAG** | 4 | 1 | 1 | 2 |

**Target: 100% documentation coverage with sample output and troubleshooting**

## 🎥 Visual Previews

We're working on adding visual demonstrations to each example:

### 🎬 **Coming Soon**
- Terminal recordings (asciinema) for complex workflows
- Screenshots for web-based examples
- GIF demos for agent interactions
- Interactive browser demos

### 📝 **Current Status**
- **Sample Output**: ✅ Text-based examples in all major READMEs
- **Screenshots**: 🔄 In progress for web examples
- **GIFs**: 🔄 Planned for agent and swarm examples
- **Interactive Demos**: 💡 Future enhancement

**Help Wanted**: Contribute visual content via PRs!

## 🤝 Contributing

To add a new example:

1. **Create folder**: `examples/my_example/`
2. **Follow template**: Copy structure from [`stdio_example`](./stdio_example)
3. **Update index**: Add entry to this README table
4. **Test**: Ensure `go run main.go` works
5. **Document**: Include setup, usage, and sample output

## 🐛 Troubleshooting

**Common Issues:**

- **Import errors**: Run `go mod tidy` in the example directory
- **Connection refused**: Check prerequisites (Ollama, API keys, etc.)
- **Permission denied**: Ensure executable permissions on scripts

**Getting Help:**

- Check individual example READMEs for specific troubleshooting
- Review main [Conduit documentation](../README.md)
- Open an issue on GitHub with example name and error details

---

**Total Examples: 33** | **Getting Started: 5 min** | **Full Tour: 2 hours**
