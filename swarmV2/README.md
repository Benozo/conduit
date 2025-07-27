# Neuron Swarm Framework

The Neuron Swarm Framework is a next-generation, production-ready multi-agent system designed for complex AI workflows. It provides a unified architecture for coordinating multiple AI agents, integrating various language models, and implementing advanced patterns like RAG (Retrieve-Augment-Generate) and ReAct (Reasoning and Acting).

## 🚀 Key Features

- **🤖 Multi-Agent Orchestration**: Coordinate multiple AI agents with different specializations and capabilities
- **🧠 Multi-LLM Support**: Seamlessly integrate OpenAI, Anthropic, Ollama, and other language model providers
- **📚 Vector Database Integration**: Comprehensive vector database support (pgvector, Milvus, Weaviate, Pinecone, in-memory)
- **🔍 Advanced RAG Workflows**: Built-in document processing, embedding, and semantic search capabilities
- **⚡ Flexible Workflow Engine**: Support for RAG, ReAct, and custom workflow patterns
- **🏗️ Type-Safe Architecture**: Strongly typed interfaces ensuring reliability and maintainability
- **🔧 Production Ready**: Comprehensive error handling, logging, and monitoring capabilities

## 🛠️ Prerequisites

- **Go 1.19+**: Modern Go version with generics support
- **Ollama** (optional): For local LLM inference
- **Vector Database** (optional): For production RAG workflows
- **API Keys**: For external LLM providers (OpenAI, Anthropic)

## 📦 Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/benozo/neuron.git
   cd neuron/swarmV2
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Verify installation**:
   ```bash
   ./status_check.sh
   ```

## 🚀 Quick Start

### Basic Agent System

```go
package main

import (
    "context"
    "log"
    
    "github.com/benozo/neuron/src/core"
    "github.com/benozo/neuron/src/agents/base"
    "github.com/benozo/neuron/src/llm/providers"
)

func main() {
    // Initialize the swarm
    swarm := core.NewSwarm()
    
    // Create coordinator
    coordinator := base.NewCoordinator("MainCoordinator")
    swarm.RegisterAgent(coordinator)
    
    // Add AI-powered agent with Ollama
    ollama := providers.NewOllamaProvider("http://localhost:11434", "llama3.2")
    aiAgent := base.NewSpecialist("AIAnalyst", "Data analysis specialist")
    aiAgent.SetLLMProvider(ollama)
    
    swarm.RegisterAgent(aiAgent)
    
    // Execute workflow
    ctx := context.Background()
    result, err := coordinator.ProcessTask(ctx, "Analyze user behavior patterns")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Result: %s", result)
}
```

### Vector-Enhanced RAG Workflow

```go
package main

import (
    "context"
    "log"
    
    "github.com/benozo/neuron/src/vectordb"
    "github.com/benozo/neuron/src/vectordb/providers"
    "github.com/benozo/neuron/src/workflows"
    "github.com/benozo/neuron/src/llm/providers"
)

func main() {
    ctx := context.Background()
    
    // Setup vector database
    vectorDB := providers.NewInMemoryProvider()
    vectorDB.Connect(ctx)
    
    // Create RAG store
    ragStore := vectordb.NewRAGStore(vectorDB)
    
    // Setup LLM
    ollama := providers.NewOllamaProvider("http://localhost:11434", "llama3.2")
    
    // Create RAG workflow
    ragWorkflow := workflows.NewRAGWorkflow(ragStore, ollama)
    
    // Add documents to knowledge base
    documents := []vectordb.Document{
        {
            ID:      "doc1",
            Content: "Machine learning is a subset of artificial intelligence...",
            Type:    vectordb.DocumentTypeText,
        },
    }
    
    ragStore.AddDocuments(ctx, "knowledge_base", documents)
    
    // Execute RAG query
    result, err := ragWorkflow.Execute(ctx, "What is machine learning?")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Generated response: %s", result)
}
```

    coordinator := base.NewCoordinator("Coordinator")
    specialist := base.NewSpecialist("Specialist")

    swarm.AddAgent(coordinator)
    swarm.AddAgent(specialist)

    workflow := workflows.NewRAGWorkflow(swarm)
    workflow.Execute()
}
```

## 📚 Examples

The framework includes comprehensive examples demonstrating various use cases:

### Available Examples

1. **📊 Coordinator Demo** (`examples/coordinator_demo/`)
   - Multi-agent coordination
   - Task delegation and management

2. **🔍 RAG Workflow** (`examples/rag_workflow/`)
   - Document retrieval and generation
   - Ollama integration for content creation

3. **🧠 ReAct Workflow** (`examples/react_workflow/`)
   - Reasoning and acting patterns
   - Step-by-step problem solving

4. **⚙️ Custom Workflow** (`examples/custom_workflow/`)
   - Building custom agent workflows
   - Extensibility demonstration

5. **🤖 Ollama Agent** (`examples/ollama_agent/`)
   - Local LLM integration
   - Standalone AI agent

6. **🚀 Multi-Agent Ollama** (`examples/multi_agent_ollama/`)
   - Multiple AI agents with different models
   - Collaborative AI workflows

7. **🗄️ Vector RAG Demo** (`examples/vector_rag_demo/`)
   - End-to-end vector database RAG
   - Semantic search and generation

8. **☁️ Cloudflare Workers AI** (`examples/cloudflare_ai/`)
   - Edge computing AI with Cloudflare Workers
   - Global low-latency inference

## 🐳 RAG Infrastructure (Docker Compose)

SwarmV2 includes a comprehensive Docker Compose setup for RAG infrastructure:

```bash
# Start all RAG services (Milvus, Weaviate, PostgreSQL+pgvector, Ollama, UIs)
./rag-infrastructure.sh start

# View all access URLs
./rag-infrastructure.sh urls
```

**Included Services:**
- **Vector Databases**: Milvus, Weaviate, PostgreSQL+pgvector
- **Management UIs**: Attu (Milvus), Weaviate Console, pgAdmin
- **AI Services**: Ollama with WebUI
- **Supporting**: Redis, MinIO, etcd

**Access URLs:**
- Milvus UI (Attu): http://localhost:3000
- Weaviate Console: http://localhost:8081  
- pgAdmin: http://localhost:5050
- Ollama WebUI: http://localhost:8083
- Redis Commander: http://localhost:8082

For detailed setup instructions, see [RAG Infrastructure Guide](docs/RAG_INFRASTRUCTURE.md).

### Running Examples

```bash
# Run all examples
cd examples && ./run_all.sh

# Run specific example
cd examples/vector_rag_demo && go run main.go

# Build and run
cd examples/multi_agent_ollama && go build && ./multi_agent_ollama
```

## 🏗️ Architecture

### Core Components

```
SwarmV2/
├── src/
│   ├── core/           # Core swarm management
│   ├── agents/         # Agent implementations
│   │   ├── base/       # Basic agent types
│   │   ├── react/      # ReAct pattern agents
│   │   └── rag/        # RAG-specialized agents
│   ├── workflows/      # Workflow orchestration
│   ├── llm/           # LLM provider interface
│   │   └── providers/ # OpenAI, Anthropic, Ollama
│   └── vectordb/      # Vector database interface
│       └── providers/ # pgvector, Milvus, Weaviate, etc.
├── examples/          # Working demonstrations
├── config/           # Configuration templates
└── README.md
```

### Key Interfaces

- **Agent**: Base interface for all agents
- **LLMProvider**: Interface for language model integration
- **VectorDB**: Interface for vector database operations
- **Workflow**: Interface for workflow execution
- **RAGStore**: High-level RAG operations

## 🔧 Configuration

### Environment Variables

```bash
# Ollama Configuration
export OLLAMA_HOST="http://localhost:11434"

# OpenAI Configuration
export OPENAI_API_KEY="your-api-key"

# Anthropic Configuration
export ANTHROPIC_API_KEY="your-api-key"

# Vector Database Configuration (for production)
export PGVECTOR_HOST="localhost"
export PGVECTOR_PORT="5432"
export MILVUS_HOST="localhost"
export MILVUS_PORT="19530"
```

### Configuration Files

Example configuration in `config/swarm.yaml`:

```yaml
swarm:
  name: "production-swarm"
  max_agents: 100
  
llm:
  default_provider: "ollama"
  providers:
    ollama:
      host: "http://localhost:11434"
      default_model: "llama3.2"
    openai:
      model: "gpt-4"
      
vectordb:
  default_provider: "inmemory"
  providers:
    inmemory:
      dimension: 384
    pgvector:
      host: "localhost"
      port: 5432
      database: "vectordb"
```

## 🧪 Testing

### Validation Script

```bash
# Comprehensive framework validation
./status_check.sh

# Manual compilation test
go build ./src/...

# Test specific example
cd examples/vector_rag_demo && go test
```

### Unit Tests

```bash
# Run all tests
go test ./src/...

# Run with coverage
go test -cover ./src/...

# Benchmark tests
go test -bench=. ./src/...
```

## 📈 Performance

- **Concurrent Agent Processing**: High-performance goroutine-based agent execution
- **Efficient Vector Search**: Optimized embedding and search operations
- **Memory Management**: Smart caching and cleanup for long-running workflows
- **Streaming Support**: Real-time response streaming for interactive applications

## 🔐 Security

- **API Key Management**: Secure environment variable handling
- **Input Validation**: Comprehensive input sanitization
- **Error Handling**: Graceful error recovery and logging
- **Resource Limits**: Configurable timeouts and resource constraints

## 🚀 Production Deployment

### Docker Support

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o swarm ./src

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/swarm .
CMD ["./swarm"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: swarmv2
spec:
  replicas: 3
  selector:
    matchLabels:
      app: swarmv2
  template:
    metadata:
      labels:
        app: swarmv2
    spec:
      containers:
      - name: swarmv2
        image: swarmv2:latest
        ports:
        - containerPort: 8080
        env:
        - name: OLLAMA_HOST
          value: "http://ollama-service:11434"
```

## 🛠️ Development

### Adding New Agents

```go
type CustomAgent struct {
    base.BaseAgent
    specialization string
}

func (a *CustomAgent) ProcessTask(ctx context.Context, task string) (string, error) {
    // Custom agent logic
    return result, nil
}
```

### Adding New LLM Providers

```go
type CustomLLMProvider struct {
    apiKey string
    baseURL string
}

func (p *CustomLLMProvider) GenerateResponse(ctx context.Context, prompt string) (string, error) {
    // Custom LLM integration
    return response, nil
}
```

### Adding New Vector Databases

```go
type CustomVectorDB struct {
    connection interface{}
}

func (db *CustomVectorDB) SearchDocuments(ctx context.Context, collection string, query string, options vectordb.SearchOptions) ([]vectordb.SearchResult, error) {
    // Custom vector search implementation
    return results, nil
}
```

## 📝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- OpenAI for GPT models
- Anthropic for Claude models
- Ollama for local LLM inference
- Vector database communities (pgvector, Milvus, Weaviate, Pinecone)

## 📞 Support

- **Documentation**: See `examples/` directory for comprehensive usage examples
- **Issues**: GitHub Issues for bug reports and feature requests
- **Discussions**: GitHub Discussions for community support

---

**SwarmV2 Framework** - Orchestrating the future of multi-agent AI systems 🚀