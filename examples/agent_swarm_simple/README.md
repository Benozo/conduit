# LLM-Powered Simple Agent Swarm Demo

This example demonstrates a **LLM-powered** agent swarm with simple task routing and processing capabilities using **Ollama integration**.

## Overview

The simple agent swarm consists of three specialized agents:

- **Router**: Analyzes requests and routes them to appropriate specialists using LLM reasoning
- **TextProcessor**: Handles text processing, formatting, and manipulation tasks
- **Analyst**: Performs data analysis and generates insights

## Key Features

- 🧠 **LLM-Powered Reasoning**: All agent decisions powered by Ollama LLM (no rule-based fallback)
- 🎯 **Intelligent Routing**: Router agent uses natural language understanding to route tasks
- 🔧 **Smart Tool Selection**: Agents use LLM to decide which tools to use and how
- 🔄 **Context-Aware Handoffs**: Agents can transfer tasks based on requirements
- 💭 **Transparent Reasoning**: All decisions include LLM reasoning explanations

## Prerequisites

- **Ollama** running at `192.168.10.10:11434` with a compatible model (e.g., `llama3.2`)
- Go 1.21+ for building and running the example

## Configuration

The example uses the following Ollama configuration:
- **Host**: `192.168.10.10:11434` (configurable via `OLLAMA_URL` env var)
- **Model**: `llama3.2` (configurable via `OLLAMA_MODEL` env var)

## Running the Example

```bash
# Navigate to the example directory
cd examples/agent_swarm_simple

# Run the example
go run main.go
```

## Demo Scenarios

The example runs four demo scenarios:

1. **Text Processing Task**: Convert text to uppercase and store in memory
2. **Text Analysis Task**: Count words and analyze text structure
3. **Data Analysis Task**: Analyze data and store insights
4. **Multi-Step Processing**: Combine text processing with analysis

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
