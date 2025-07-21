# Conduit Examples

This folder contains working examples for ConduitMCP including tool calling, LLM integration, agent frameworks, and protocol demonstrations. Each example is self-contained with clear setup instructions and sample outputs.

## 🚀 Quick Start

Choose an example based on your use case:

- **New to Conduit?** → Start with [`stdio_example`](./stdio_example) or [`pure_library`](./pure_library)
- **Want LLM integration?** → Try [`ollama`](./ollama) or [`agents_ollama`](./agents_ollama)
- **Building web apps?** → Check out [`sse_example`](./sse_example) or [`pure_library_web`](./pure_library_web)
- **Need agent coordination?** → Explore [`agent_swarm`](./agent_swarm) or [`multi_llm_swarm`](./multi_llm_swarm)

## 📋 Complete Example Index

| Example | Description | Features | Complexity | Prerequisites |
|---------|-------------|----------|------------|---------------|
| **Core Protocol** |
| [`stdio_example`](./stdio_example) | MCP stdio server for VS Code, Cline, etc. | stdio protocol, tool calling | ⭐ | None |
| [`sse_example`](./sse_example) | HTTP/SSE server for web applications | HTTP API, real-time streaming | ⭐⭐ | None |
| [`pure_mcp`](./pure_mcp) | Raw MCP protocol implementation | Pure MCP, no server wrapper | ⭐ | None |
| **Library Usage** |
| [`pure_library`](./pure_library) | Use Conduit as Go library | Library integration, custom tools | ⭐ | None |
| [`pure_library_cli`](./pure_library_cli) | CLI tool with MCP components | Command-line interface | ⭐⭐ | None |
| [`pure_library_web`](./pure_library_web) | Custom web server with MCP | Web server, custom endpoints | ⭐⭐ | None |
| [`embedded`](./embedded) | Embed Conduit in existing apps | Application integration | ⭐⭐ | None |
| **LLM Integration** |
| [`ollama`](./ollama) | Local LLM with tool calling | Ollama integration, auto tool selection | ⭐⭐ | Ollama |
| [`openai`](./openai) | OpenAI GPT with tools | OpenAI API, cloud LLM | ⭐⭐ | OpenAI API key |
| [`model_integration`](./model_integration) | Custom model patterns | Custom LLM integration | ⭐⭐⭐ | Custom model |
| **Tool Development** |
| [`custom_tools`](./custom_tools) | Enhanced tool registration | Rich schemas, validation | ⭐⭐ | None |
| [`builtin_tools_test`](./builtin_tools_test) | Test all built-in tools | Tool testing, validation | ⭐ | None |
| **Agent Framework** |
| [`ai_agents`](./ai_agents) | AI Agents with task management | Agent framework, task execution | ⭐⭐⭐ | None |
| [`agents_test`](./agents_test) | Basic agent functionality | Agent testing | ⭐⭐ | None |
| [`agents_ollama`](./agents_ollama) | Agents with Ollama LLM | AI agents, local LLM | ⭐⭐⭐ | Ollama |
| [`agents_deepinfra`](./agents_deepinfra) | Agents with DeepInfra | AI agents, cloud inference | ⭐⭐⭐ | DeepInfra API |
| [`agents_library_mode`](./agents_library_mode) | Library-mode agent usage | Agent patterns | ⭐⭐ | None |
| [`agents_mock_llm`](./agents_mock_llm) | Mock LLM for testing | Testing, development | ⭐⭐ | None |
| [`agents_vue_builder`](./agents_vue_builder) | Vue.js app builder agent | Code generation, Vue.js | ⭐⭐⭐ | None |
| **Agent Swarm** |
| [`agent_swarm`](./agent_swarm) | Basic agent coordination | Multi-agent, handoffs | ⭐⭐⭐ | None |
| [`agent_swarm_llm`](./agent_swarm_llm) | LLM-powered agent swarm | LLM coordination, Ollama | ⭐⭐⭐⭐ | Ollama |
| [`agent_swarm_simple`](./agent_swarm_simple) | Simple swarm demo | Basic swarm patterns | ⭐⭐ | None |
| [`agent_swarm_workflows`](./agent_swarm_workflows) | Advanced workflow patterns | DAG, Supervisor, Pipeline | ⭐⭐⭐⭐ | None |
| [`multi_llm_swarm`](./multi_llm_swarm) | Multi-LLM agent architecture | Multiple LLM providers | ⭐⭐⭐⭐ | Multiple APIs |
| **RAG & Advanced** |
| [`rag`](./rag) | RAG with document processing | Document analysis, embeddings | ⭐⭐⭐ | Vector DB |
| [`rag_chat_terminal`](./rag_chat_terminal) | Terminal RAG chat interface | CLI RAG, interactive | ⭐⭐⭐ | Documents |
| [`rag_real_world`](./rag_real_world) | Production RAG patterns | Real-world RAG | ⭐⭐⭐⭐ | Vector DB |
| [`langchain_mcp_integration`](./langchain_mcp_integration) | LangChain integration | LangChain + MCP | ⭐⭐⭐ | LangChain |
| **Specialized** |
| [`simple_mcp_agent`](./simple_mcp_agent) | Minimal MCP agent | Basic agent pattern | ⭐ | None |
| [`comprehensive_test`](./comprehensive_test) | Full system testing | Integration testing | ⭐⭐ | None |
| [`agents_html_amender`](./agents_html_amender) | HTML processing agent | HTML manipulation | ⭐⭐ | None |

**Complexity Legend:**
- ⭐ = Beginner (5 min setup)
- ⭐⭐ = Intermediate (15 min setup)
- ⭐⭐⭐ = Advanced (30 min setup)
- ⭐⭐⭐⭐ = Expert (1+ hour setup)

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

## 🏃‍♂️ Run All Examples

```bash
# Run the quick setup script
./run_all.sh

# Or test specific categories
./run_all.sh --protocol     # stdio, sse, pure_mcp
./run_all.sh --llm          # ollama, openai, model_integration
./run_all.sh --agents       # ai_agents, agent_swarm, multi_llm_swarm
./run_all.sh --rag          # rag examples
```

## 📖 Documentation Standards

Each example follows this structure:

```
example_name/
├── README.md          # Standardized docs with setup & examples
├── main.go           # Primary executable
├── go.mod            # Dependencies
├── test_*.sh         # Test scripts (optional)
└── media/            # Screenshots/GIFs (optional)
```

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
