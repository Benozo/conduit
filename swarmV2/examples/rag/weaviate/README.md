# Weaviate RAG Example

This example demonstrates how to build a Retrieval-Augmented Generation (RAG) system using Weaviate vector database.

## Features

- **Document Storage**: Store documents with vector embeddings in Weaviate
- **Similarity Search**: Find relevant documents using vector similarity
- **GraphQL API**: Interact with Weaviate using GraphQL queries
- **HTTP Integration**: Pure HTTP-based implementation (no external dependencies)
- **Simple Embeddings**: Basic hash-based embeddings for demonstration

## Prerequisites

1. **Weaviate running**: Ensure Weaviate is running via Docker Compose
2. **Go 1.21+**: Required for running the example

## Quick Start

### 1. Start the Infrastructure

From the swarmV2 root directory:
```bash
# Start Weaviate and other services
docker-compose up -d weaviate weaviate-console

# Check if Weaviate is running
curl http://localhost:8080/v1/.well-known/ready
```

### 2. Run the Example

```bash
# Navigate to the weaviate example
cd examples/rag/weaviate

# Run the demo
go run main.go
```

### 3. Environment Variables (Optional)

```bash
export WEAVIATE_URL=http://localhost:8080
go run main.go
```

## What the Example Does

1. **Schema Initialization**: Creates a Weaviate class for documents
2. **Document Ingestion**: Adds sample documents with embeddings
3. **Similarity Search**: Demonstrates vector-based document retrieval
4. **RAG Workflow**: Shows the complete retrieval process

## Sample Output

```
=== Weaviate RAG Demo ===
🌐 Weaviate URL: http://localhost:8080
📚 Class Name: Document

🏗️ Initializing Weaviate schema...
✅ Created schema for class: Document

📚 Adding sample documents...
✅ Added document: Introduction to Vector Databases
✅ Added document: RAG Systems Architecture
✅ Added document: Weaviate Features and Capabilities
✅ Added document: Building AI-Powered Applications

🔍 Similarity Search Demo:

📝 Query 1: vector database technology
Found 2 similar documents:
  1. Introduction to Vector Databases (Category: Technology)
     Content preview: Vector databases are specialized databases designed to store and query high-dimensional vectors...
  2. Weaviate Features and Capabilities (Category: Database)
     Content preview: Weaviate is an open-source vector database that supports semantic search, automatic classification...
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Documents     │───▶│   Weaviate      │───▶│   Search        │
│   + Embeddings  │    │   Vector DB     │    │   Results       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │    │   GraphQL       │    │   JSON          │
│   Integration   │    │   Queries       │    │   Response      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Key Components

### 1. Document Structure
```go
type Document struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Category  string    `json:"category"`
    Timestamp time.Time `json:"timestamp"`
}
```

### 2. Weaviate Schema
- **Class**: Document
- **Properties**: title, content, category, timestamp
- **Vectorizer**: none (manual embeddings)

### 3. Search Process
1. Generate embedding for query
2. Send GraphQL nearVector query
3. Parse results and return documents

## Weaviate Console

Access the Weaviate Console at: **http://localhost:8081**

Features:
- Browse collections and objects
- Execute GraphQL queries
- Visualize vector data
- Monitor system health

## Production Considerations

1. **Real Embeddings**: Replace simple hash embeddings with proper models:
   - OpenAI text-embedding-ada-002
   - Cohere embed-english-v3.0
   - Hugging Face sentence-transformers

2. **Authentication**: Add API key authentication for production

3. **Indexing**: Configure appropriate vector indexes for scale:
   - HNSW for high recall
   - IVF for memory efficiency

4. **Monitoring**: Implement proper logging and metrics

## Integration with LLM

For full RAG functionality, combine with an LLM provider:

```go
// Add LLM integration for generation
func (w *WeaviateRAGSystem) Query(ctx context.Context, question string) (string, error) {
    // 1. Search for relevant documents
    docs, err := w.SearchSimilar(ctx, question, 3)
    if err != nil {
        return "", err
    }
    
    // 2. Build context
    context := buildContext(docs)
    
    // 3. Generate response using LLM
    return llmProvider.GenerateResponse(buildPrompt(context, question))
}
```

## Related Examples

- **[Milvus RAG](../milvus/)** - Similar example using Milvus
- **[Cloudflare AI](../../cloudflare_ai/)** - LLM integration example

## Troubleshooting

1. **Connection Issues**: Ensure Weaviate is running on port 8080
2. **Schema Errors**: Check if class already exists
3. **Search Problems**: Verify embeddings are being generated

This example provides a solid foundation for building production RAG systems with Weaviate!
