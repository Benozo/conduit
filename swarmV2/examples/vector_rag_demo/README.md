# Vector-Enhanced RAG Workflow Demo

This example demonstrates a comprehensive RAG (Retrieve-Augment-Generate) workflow enhanced with vector database integration for semantic search and knowledge management.

## What it demonstrates

- **Vector Database Integration**: Support for multiple vector databases (Milvus, PgVector, Weaviate, Pinecone)
- **Semantic Search**: Vector-based retrieval using embeddings for better relevance
- **Knowledge Management**: Document processing, chunking, and indexing
- **AI-Powered Generation**: Ollama-based content generation with retrieved context
- **End-to-End Pipeline**: Complete RAG workflow from knowledge ingestion to answer generation

## Architecture

```
Vector-Enhanced RAG Pipeline
├── 📚 Knowledge Base
│   ├── Document Processing (chunking, metadata)
│   ├── Embedding Generation (text → vectors)
│   └── Vector Storage (Milvus/PgVector/Weaviate/Pinecone)
├── 🔍 Semantic Retrieval
│   ├── Query Embedding
│   ├── Vector Similarity Search
│   └── Context Assembly
├── 🤖 AI Generation
│   ├── Context-Aware Prompting
│   ├── Ollama llama3.2 Generation
│   └── Response Formatting
└── 📋 Quality Evaluation
    ├── Content Validation
    ├── Relevance Scoring
    └── Final Assessment
```

## Prerequisites

### Ollama Setup
- Ollama server at `192.168.10.10:11434`
- `llama3.2` model available

```bash
# Pull the required model
ollama pull llama3.2
```

### Vector Database (Optional)
Choose one of the following:

#### Milvus (Default in demo)
```bash
# Using Docker
docker run -d --name milvus-standalone \
  -p 19530:19530 -p 9091:9091 \
  milvusdb/milvus:latest standalone
```

#### PostgreSQL + pgvector
```bash
# Install pgvector extension
# CREATE EXTENSION vector;
```

#### Weaviate
```bash
# Using Docker
docker run -d --name weaviate \
  -p 8080:8080 \
  semitechnologies/weaviate:latest
```

#### Pinecone
- Sign up for Pinecone account
- Get API key and environment

## How to run

```bash
cd /home/engineone/Projects/AI/ConduitMCP/swarmV2/examples/vector_rag_demo
go run main.go
```

## Features Demonstrated

### 1. Vector Database Integration
- **Multiple Providers**: Support for 4 different vector databases
- **Unified Interface**: Common API across all providers
- **Connection Management**: Robust connection handling and testing

### 2. Document Processing
- **Text Chunking**: Intelligent document splitting with overlap
- **Metadata Management**: Rich metadata tracking and filtering
- **Multiple Formats**: Support for text, JSON, PDF (extensible)

### 3. Embedding Generation
- **Text Embeddings**: Convert text to high-dimensional vectors
- **Configurable Models**: Support for different embedding models
- **Batch Processing**: Efficient batch embedding generation

### 4. Semantic Search
- **Vector Similarity**: Find semantically similar content
- **Relevance Scoring**: Score-based result ranking
- **Context Assembly**: Combine multiple relevant chunks

### 5. AI-Powered Generation
- **Context-Aware Prompts**: Use retrieved context for generation
- **Model Integration**: Ollama llama3.2 for high-quality responses
- **Response Formatting**: Clean, structured output

### 6. Knowledge Base Management
- **Dynamic Addition**: Add documents at runtime
- **Statistics Tracking**: Monitor collection size and performance
- **Flexible Metadata**: Rich metadata for filtering and organization

## Sample Knowledge Base

The demo includes sample knowledge covering:
- Machine Learning Fundamentals
- Deep Learning Concepts
- Supervised/Unsupervised Learning
- Reinforcement Learning
- Feature Engineering
- Model Validation
- Overfitting Prevention

## Expected Output

The demo will:
1. **Setup**: Initialize vector database and components
2. **Knowledge Ingestion**: Add sample ML knowledge to the vector store
3. **Collection Creation**: Create a vector collection with proper indexing
4. **Query Processing**: Process multiple test queries through the pipeline
5. **Results Display**: Show generated answers with context
6. **Statistics**: Display knowledge base statistics and metrics

## Sample Queries

- "What are the main types of machine learning and how do they differ?"
- "Explain overfitting and how to prevent it in machine learning models"
- "What is feature engineering and why is it important?"

## Code Structure

- **Vector Database Interfaces**: `src/vectordb/interfaces.go`
- **Provider Implementations**: `src/vectordb/providers/`
- **Document Processing**: `src/vectordb/processor.go`
- **RAG Store**: `src/vectordb/rag_store.go`
- **Enhanced Retriever**: Vector-based semantic search
- **Ollama Generator**: AI-powered content generation
- **Workflow Orchestration**: End-to-end pipeline management

## Customization Options

### Vector Database Provider
Change the provider in main.go:
```go
// Milvus
vectorDB := vectordbProviders.NewMilvusProvider("localhost", 19530, "", "")

// PgVector
vectorDB := vectordbProviders.NewPgVectorProvider("localhost", 5432, "rag_db", "user", "pass")

// Weaviate
vectorDB := vectordbProviders.NewWeaviateProvider("http://localhost:8080", "api-key")

// Pinecone
vectorDB := vectordbProviders.NewPineconeProvider("api-key", "env", "index")
```

### Embedding Model
```go
embedder := vectordb.NewSimpleEmbeddingProvider(768, "your-model-name")
```

### Collection Configuration
```go
options := map[string]interface{}{
    "metric_type": "cosine",  // or "euclidean", "dot_product"
    "index_type":  "IVF_FLAT", // or "HNSW", "ANNOY"
}
```

This example showcases a production-ready approach to implementing RAG systems with vector databases and modern AI models.
