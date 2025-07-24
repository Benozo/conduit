# Enhanced Coordinator Demo with Ollama Integration

This example demonstrates how to use the Coordinator agent to manage both traditional specialist agents and AI-powered agents using Ollama.

## What it demonstrates

- **Hybrid Coordination**: Managing both traditional and AI-powered agents
- **Coordinator**: Central agent that manages and coordinates other agents
- **Traditional Specialists**: Specialized agents with specific tasks
- **AI-Powered Advisor**: Intelligent agent using Ollama for advanced insights
- **Agent Registration**: How to register traditional agents with a coordinator
- **AI Integration**: How to integrate Ollama-powered agents alongside traditional ones
- **Metrics and Monitoring**: How to track coordinator performance

## Prerequisites

- Ollama server at `192.168.10.10:11434`
- `llama3.2` model available
- Traditional framework components

```bash
# Pull the required model
ollama pull llama3.2

# Verify model is available
ollama list
```

## How to run

```bash
cd /home/engineone/Projects/AI/ConduitMCP/swarmV2/examples/coordinator_demo
go run main.go
```

## Expected output

The example will:
1. Create a coordinator instance
2. Create multiple traditional specialist agents with different roles
3. Create an AI-powered advisor agent using Ollama
4. Test connection to the Ollama server
5. Register traditional specialists with the coordinator
6. Display coordinator status and metrics for both types of agents
7. Test the AI advisor with a coordination task
8. Show performance metrics and agent information

## Architecture

```
Enhanced Coordinator System
├── Traditional Coordinator
│   ├── TraditionalAnalyst
│   ├── ReportGenerator
│   └── QualityController
└── AI-Powered Agent (separate)
    └── AIAdvisor (using llama3.2)
```

## Code structure

- Uses agents from `src/agents/base/` for traditional coordination
- Uses `src/llm/providers/ollama.go` for AI integration
- Demonstrates hybrid agent management
- Shows both traditional and AI agent coordination patterns
- Displays comprehensive metrics and status information
