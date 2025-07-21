# LLM-Powered Multi-Agent Content & Analysis Workflow

This example demonstrates a comprehensive **LLM-powered** agent swarm for content creation and data analysis workflows using **Ollama integration**.

## Overview

This is the main agent swarm example showcasing advanced multi-agent coordination for real-world business workflows. The swarm consists of four specialized agents:

- **Coordinator**: Project manager that orchestrates complex workflows using LLM reasoning
- **ContentCreator**: Handles content creation, research, and writing tasks
- **DataAnalyst**: Performs data analysis, insights generation, and reporting
- **MemoryManager**: Manages shared information and context across agents

## Key Features

- 🧠 **LLM-Powered Coordination**: Project orchestration using Ollama LLM reasoning
- 🎯 **Intelligent Task Delegation**: Context-aware routing based on task requirements
- 🔧 **Smart Tool Selection**: LLM-driven tool usage and workflow optimization
- 🔄 **Multi-Agent Collaboration**: Seamless handoffs and knowledge sharing
- 💭 **Transparent Reasoning**: All decisions include LLM explanations
- 📚 **Persistent Memory**: Cross-conversation context and knowledge management

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
cd examples/agent_swarm

# Run the example
go run main.go
```

## Key Features Demonstrated

### Multi-Agent Handoffs
```go
// Agents can transfer execution to other specialists
transferToContentCreator := swarm.CreateHandoffFunction(
    "transfer_to_content_creator",
````markdown
## Demo Scenarios

The example runs four comprehensive scenarios:

1. **Content Creation Project**: Comprehensive article writing with research coordination
2. **Data Analysis Project**: Customer behavior analysis with reporting workflow
3. **Cross-Agent Memory Sharing**: Knowledge management and status reporting
4. **Complex Multi-Phase Project**: AI product launch with multi-agent coordination

## Expected Output

The demo will show:
- LLM-powered project breakdown and task delegation
- Multi-agent workflow coordination and collaboration
- Context-aware tool usage and decision making
- Agent handoffs with transparent reasoning
- Persistent memory and knowledge sharing
- Comprehensive execution metrics and timing

## Error Handling

**Important**: This example uses pure LLM reasoning with **no rule-based fallback**. If Ollama is unavailable or the LLM fails, the example will error out as intended.

## Advanced Capabilities

### Project Coordination
- LLM-powered project analysis and breakdown
- Intelligent task delegation based on agent specializations
- Dynamic workflow adaptation based on requirements

### Content Workflows
- Research and content creation coordination
- Multi-step writing processes with quality assurance
- Information synthesis and presentation

### Data Analysis Workflows  
- Comprehensive data analysis with insights generation
- Multi-stage reporting with stakeholder communication
- Cross-project knowledge management

### Memory Management
- Persistent context storage and retrieval
- Cross-agent knowledge sharing
- Project continuity and status tracking

## Customization

You can extend this example by:
- Adding specialized agents for your domain (legal, medical, financial, etc.)
- Implementing custom tools for specific workflows
- Creating industry-specific project templates
- Adding external system integrations
- Scaling to larger agent swarms

## Related Examples

- **agent_swarm_simple**: Basic LLM-powered agent routing and coordination
- **agent_swarm_workflows**: Advanced workflow patterns and orchestration
- **agent_swarm_llm**: Comprehensive LLM integration demonstration

## LLM Integration Details

This example showcases:
- Complex project reasoning and coordination
- Multi-agent workflow orchestration
- Context-aware decision making and tool selection
- Intelligent handoff strategies
- Persistent memory and knowledge management
- No rule-based logic - pure LLM coordination

Perfect for understanding production-scale LLM-powered multi-agent systems!
```
