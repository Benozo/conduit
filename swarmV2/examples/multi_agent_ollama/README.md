# Multi-Agent Ollama System

This example demonstrates a hybrid system combining traditional framework agents with AI-powered agents using Ollama.

## What it demonstrates

- **Hybrid Architecture**: Traditional agents + AI-powered agents working together
- **Multi-Model AI**: Different AI agents using different language models
- **AI Specialists**: Multiple AI agents with different specializations and models
- **Workflow Coordination**: Using coordinator to manage traditional agents
- **Real AI Integration**: Live interaction with llama3.2 and gemma2:latest for specialized tasks
- **Multi-step Processing**: Complex workflows with AI assistance

## How to run

```bash
cd /home/engineone/Projects/AI/ConduitMCP/swarmV2/examples/multi_agent_ollama
go run main.go
```

## System Architecture

```
Coordinator
├── Traditional Agents
│   ├── ProjectManager (coordination)
│   └── QualityController (quality assurance)
└── AI-Powered Agents (separate)
    ├── AIDataAnalyst (data analysis via llama3.2)
    └── AIContentWriter (content creation via gemma2:latest)
```

## Expected Workflow

1. **Setup**: Initialize coordinator with traditional agents
2. **AI Connection**: Verify connection to Ollama server
3. **Data Analysis**: AI analyst processes business queries
4. **Content Creation**: AI writer creates summaries and reports
5. **Status Report**: Show system status and agent activity

## Prerequisites

- Ollama server at `192.168.10.10:11434`
- `llama3.2` model available for data analysis
- `gemma2:latest` model available for content creation
- All traditional framework components

```bash
# Pull both required models
ollama pull llama3.2
ollama pull gemma2:latest

# Verify models are available
ollama list
```

## Key Features

- **Multi-Model Support**: Each AI agent uses a different language model optimized for its task
- **Specialization**: Each AI agent has a specific role and prompt context
- **Model Diversity**: Demonstrates llama3.2 for analysis and gemma2:latest for content creation
- **Error Handling**: Robust error handling for AI connections
- **Performance**: Shows response times and connection status for each model
- **Scalability**: Easy to add more AI specialists or traditional agents

This example shows how to build sophisticated multi-agent systems that combine the reliability of traditional agents with the intelligence of modern LLMs!
