# SwarmV2 Examples

This directory contains working examples that demonstrate various patterns and capabilities of the SwarmV2 agent framework.

## Available Examples

### 1. RAG Workflow (`rag_workflow/`)
Demonstrates the Retrieve-Augment-Generate pattern with:
- Retriever, Evaluator agents and AI-powered Generator
- Real AI content generation using Ollama llama3.2
- Complete RAG workflow execution with step-by-step processing
- Error handling and result processing

### 2. React Workflow (`react_workflow/`)
Shows the AI-enhanced Reason-Act-Observe pattern with:
- **AI-Powered Reasoner**: Uses Ollama for intelligent situation analysis
- **AI-Guided Actor**: Generates detailed action plans using AI reasoning
- **AI-Enhanced Observer**: Provides intelligent monitoring and assessment
- **Multiple Scenarios**: Demonstrates React cycles across different problem domains
- **Real-Time AI Integration**: Complete integration with Ollama llama3.2
- **Fallback Handling**: Graceful degradation when AI is unavailable

### 3. Coordinator Demo (`coordinator_demo/`)
Illustrates hybrid agent coordination with:
- Coordinator managing multiple traditional specialists
- AI-powered advisor agent using Ollama
- Agent registration and status tracking
- Metrics collection and reporting
- Real AI integration for coordination insights

### 4. Cloudflare Workers AI (`cloudflare_ai/`)
Demonstrates edge computing AI with Cloudflare Workers AI:
- **Global Edge Network**: AI processing on 200+ worldwide locations
- **Multiple AI Models**: Support for Llama, Mistral, Phi, Gemma models
- **Collaborative Workflow**: Multi-agent business scenario (analysis → content → strategy)
- **Cost-Effective**: Pay-per-use serverless AI without infrastructure
- **Low Latency**: Edge computing for faster AI responses
- **Easy Integration**: Simple API integration with robust error handling

### 5. Ollama Agent (`ollama_agent/`)
Shows real LLM integration with Ollama:
- Connection to Ollama server at 192.168.10.10
- Using llama3.2 model for AI responses
- Agent-LLM hybrid architecture
- Error handling and connection testing

### 6. Multi-Agent Ollama (`multi_agent_ollama/`)
Demonstrates hybrid multi-agent system with:
- Traditional specialist agents
- AI-powered assistant agents using Ollama
- Coordinated workflow between different agent types
- Real-world problem-solving scenario

### 7. Vector RAG Demo (`vector_rag_demo/`)
Showcases comprehensive vector-enhanced RAG with:
- Vector database integration (Milvus, PgVector, Weaviate, Pinecone)
- Semantic search and document retrieval
- Knowledge base management and document processing
- AI-powered content generation using Ollama
- End-to-end RAG pipeline with vector embeddings

## How to run examples

Each example is a standalone Go program with its own `main.go` file:

```bash
# Run RAG workflow example
cd rag_workflow && go run main.go

# Run React workflow example  
cd react_workflow && go run main.go

# Run Coordinator demo
cd coordinator_demo && go run main.go

# Run Cloudflare Workers AI example
cd cloudflare_ai && go run main.go

# Run Custom workflow example
cd custom_workflow && go run main.go

# Run Ollama agent example
cd ollama_agent && go run main.go

# Run Multi-agent Ollama example
cd multi_agent_ollama && go run main.go

# Run Vector RAG demo
cd vector_rag_demo && go run main.go
```

## Prerequisites

Most examples work out of the box, but the Ollama agent requires:
- Ollama server running at `192.168.10.10:11434`
- llama3.2 model pulled and available

```bash
# Install and setup Ollama
curl -fsSL https://ollama.ai/install.sh | sh
ollama serve
ollama pull llama3.2
```

## Framework Components Used

- **Core**: Agent interfaces, Registry, Swarm management
- **Base Agents**: Coordinator, Specialist, Workflow
- **React Agents**: Reasoner, Actor, Observer
- **RAG Agents**: Retriever, Generator, Evaluator
- **Workflows**: React, RAG, and Custom patterns

## Architecture Highlights

- **Interface-driven design**: All agents implement common interfaces
- **Type safety**: Compile-time verification of agent interactions
- **Modularity**: Each component can be used independently
- **Extensibility**: Easy to create new agent types and workflows
- **Thread safety**: Concurrent execution support with proper synchronization

Each example includes detailed comments and error handling to help understand the framework usage patterns.
