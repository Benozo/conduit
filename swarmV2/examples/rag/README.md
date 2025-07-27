# SwarmV2 RAG Examples

This directory contains Retrieval-Augmented Generation (RAG) examples **integrated with the SwarmV2 framework** using different vector databases.

## Overview

RAG (Retrieval-Augmented Generation) is a technique that combines information retrieval with language generation to create more accurate and informed AI responses. These examples demonstrate how to build RAG systems using **SwarmV2 agents** and popular vector databases.

## SwarmV2 Integration Features

### 🤖 Agent-Based Architecture
- **RAG Agents**: Specialized agents for document retrieval and response generation
- **Coordinator**: Centralized task delegation and workflow management
- **Multi-Agent Coordination**: Multiple agents working together on complex RAG tasks

### 🔄 Task Execution Patterns
- **rag_query:**: Full RAG pipeline (retrieve + generate)
- **search:**: Vector similarity search only
- **add_document:**: Knowledge base management
- **bulk_search:**: Batch processing operations

## Available Examples

### 🔍 [Weaviate RAG Agent](./weaviate/)
- **Technology**: Weaviate vector database + SwarmV2 agents
- **API**: GraphQL + REST through SwarmV2 coordination
- **UI**: Weaviate Console (http://localhost:8081)
- **Features**: Agent-based semantic search, task delegation, coordinated workflows

### 🔍 [Milvus RAG Agent](./milvus/)
- **Technology**: Milvus vector database + SwarmV2 high-performance agents
- **API**: REST API v2 through SwarmV2 coordination
- **UI**: Attu Admin UI (http://localhost:3000)
- **Features**: High-performance agent coordination, distributed search, bulk operations

## Infrastructure Setup

All examples use the same Docker Compose infrastructure. From the swarmV2 root directory:

```bash
# Start all vector databases and UIs
docker-compose up -d

# Or start specific services
docker-compose up -d weaviate weaviate-console  # For Weaviate
docker-compose up -d etcd minio milvus attu      # For Milvus
```

### Infrastructure Components

| Service | Port | Purpose | UI Access |
|---------|------|---------|-----------|
| **Weaviate** | 8080 | Vector database | [Console](http://localhost:8081) |
| **Milvus** | 19530 | Vector database | [Attu](http://localhost:3000) |
| **PostgreSQL** | 5432 | pgvector support | [pgAdmin](http://localhost:5050) |
| **Redis** | 6379 | Caching | [Commander](http://localhost:8082) |
| **MinIO** | 9000 | Object storage | [Console](http://localhost:9001) |

## Quick Start

### 1. Start Infrastructure
```bash
# From swarmV2 root
docker-compose up -d
```

### 2. Choose Your Vector Database

#### Weaviate Example
```bash
cd examples/rag/weaviate
go run main.go
```

#### Milvus Example  
```bash
cd examples/rag/milvus
go run main.go
```

### 3. Access Web UIs

- **Weaviate Console**: http://localhost:8081
- **Milvus Attu**: http://localhost:3000
- **pgAdmin**: http://localhost:5050 (admin@swarmv2.com / admin123)
- **Redis Commander**: http://localhost:8082

## SwarmV2 Agent Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Query    │───▶│   SwarmV2       │───▶│   RAG Agent     │
│                 │    │   Coordinator   │    │   (Weaviate/    │
└─────────────────┘    └─────────────────┘    │    Milvus)      │
                                │              └─────────────────┘
                                ▼                       │
┌─────────────────┐    ┌─────────────────┐              ▼
│   Final         │◀───│   Response      │    ┌─────────────────┐
│   Response      │    │   Coordination  │◀───│   Vector        │
└─────────────────┘    └─────────────────┘    │   Database      │
                                              └─────────────────┘
```

## SwarmV2 Task Execution Flow

1. **User submits query** to SwarmV2 Coordinator
2. **Coordinator delegates** task to appropriate RAG Agent
3. **RAG Agent executes** specialized RAG operations:
   - Document retrieval through vector search
   - Context building and relevance scoring
   - Response generation with confidence metrics
4. **Coordinator returns** final response to user

## Agent Capabilities

### Weaviate RAG Agent
- `document_retrieval`: GraphQL-based document search
- `semantic_search`: Advanced semantic similarity
- `context_generation`: Intelligent context building
- `rag_query`: Full RAG pipeline execution
- `knowledge_base_management`: Schema and collection management

### Milvus RAG Agent  
- `high_performance_search`: Sub-10ms vector search
- `vector_similarity`: Advanced similarity algorithms
- `distributed_rag`: Scalable distributed processing
- `scalable_indexing`: Billion-scale vector indexing
- `multi_vector_types`: Multiple vector type support
- `advanced_filtering`: Complex query filtering

## Vector Database Comparison

| Feature | Weaviate | Milvus | PostgreSQL+pgvector |
|---------|----------|--------|-------------------|
| **Performance** | High | Very High | Medium |
| **Scalability** | Good | Excellent | Limited |
| **API** | GraphQL + REST | REST | SQL |
| **Ease of Use** | Easy | Medium | Easy |
| **Features** | Knowledge Graph | Advanced Indexing | SQL Integration |
| **Best For** | Semantic Apps | Large Scale | Existing SQL Apps |

## Integration with SwarmV2

These RAG examples can be integrated with SwarmV2 agents:

```go
// Create a RAG-powered agent
type RAGAgent struct {
    name        string
    ragSystem   *WeaviateRAGSystem  // or MilvusRAGSystem
    llmProvider *providers.CloudflareAIProvider
}

func (r *RAGAgent) Execute(task string) (string, error) {
    // 1. Search for relevant documents
    docs, err := r.ragSystem.SearchSimilar(context.Background(), task, 3)
    if err != nil {
        return "", err
    }
    
    // 2. Build context and generate response
    context := buildContext(docs)
    prompt := buildPrompt(context, task)
    
    return r.llmProvider.GenerateResponse(prompt)
}
```

## Production Considerations

### 1. Embedding Models
Replace simple hash embeddings with production models:

```bash
# OpenAI
export OPENAI_API_KEY=your_key

# Cohere  
export COHERE_API_KEY=your_key

# Hugging Face
export HF_API_KEY=your_key
```

### 2. Performance Optimization
- Use proper vector indexes (HNSW, IVF)
- Implement connection pooling
- Add caching layer with Redis
- Monitor resource usage

### 3. Security
- Enable authentication on vector databases
- Use HTTPS for production deployments
- Implement rate limiting
- Add input sanitization

## Environment Variables

Common environment variables for all examples:

```bash
# Vector Database URLs
export WEAVIATE_URL=http://localhost:8080
export MILVUS_URL=http://localhost:19530
export POSTGRES_URL=postgresql://postgres:password@localhost:5432/vectordb

# LLM Integration
export CLOUDFLARE_CUSTOM_URL=https://intent.moanalabs.com
export CLOUDFLARE_CUSTOM_API_KEY=your_key
export CLOUDFLARE_MODEL=@cf/meta/llama-4-scout-17b-16e-instruct

# Redis Caching
export REDIS_URL=redis://localhost:6379
```

## Monitoring and Maintenance

### Health Checks
```bash
# Check all services
docker-compose ps

# Health endpoints
curl http://localhost:8080/v1/.well-known/ready  # Weaviate
curl http://localhost:19530/healthz              # Milvus
```

### Logs
```bash
# View logs
docker-compose logs -f weaviate
docker-compose logs -f milvus
```

### Backup
```bash
# Backup data volumes
docker run --rm -v swarmv2_weaviate:/data -v $(pwd):/backup alpine tar czf /backup/weaviate-backup.tar.gz /data
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**: Check if ports are already in use
2. **Memory Issues**: Increase Docker memory allocation
3. **Connection Errors**: Verify services are running and healthy
4. **Permission Issues**: Check Docker volume permissions

### Debug Commands
```bash
# Check container status
docker-compose ps

# View service logs
docker-compose logs [service-name]

# Test connectivity
curl -f http://localhost:8080/v1/.well-known/ready
curl -f http://localhost:19530/healthz
```

## Next Steps

1. **Experiment**: Try both Weaviate and Milvus examples
2. **Integrate**: Combine with LLM providers for full RAG
3. **Scale**: Move to production-ready embeddings
4. **Deploy**: Use the examples as foundation for your applications

## Related Resources

- **[SwarmV2 Documentation](../../README.md)**
- **[Cloudflare AI Example](../cloudflare_ai/)**
- **[Docker Compose Setup](../../docker-compose.yml)**
- **[Infrastructure Guide](../../docs/rag-infrastructure.md)**

Happy building with RAG and SwarmV2! 🚀
