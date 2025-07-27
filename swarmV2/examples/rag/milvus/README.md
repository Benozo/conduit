# Milvus RAG Example with SwarmV2

This example demonstrates a complete **Retrieval-Augmented Generation (RAG)** system using:
- **SwarmV2** agent architecture (retriever + generator pattern)
- **Milvus** vector database for knowledge storage using official Go SDK
- **Ollama** LLM for response generation
- **Automatic document setup** with 128-dimensional embeddings

## 🏗️ Project Structure

```
examples/rag/milvus/
├── main.go                 # Complete SwarmV2 RAG implementation with document setup
├── test_connection.go      # Milvus connection testing utility
├── test.sh                 # Test script
└── README.md              # This file
```

## 🚀 Quick Start

### 1. Prerequisites

- **Milvus** running on `localhost:19530`
- **Ollama** running on `192.168.10.10:11434` (configurable)
- **Go 1.19+** with modules enabled

### 2. Run the Example

The example includes automatic document setup - no manual data insertion needed!

```bash
# From SwarmV2 root directory
cd /home/engineone/Projects/AI/ConduitMCP/swarmV2
go run examples/rag/milvus/main.go
```

Or use the test script:

```bash
cd examples/rag/milvus
./test.sh
```

### 3. Test Milvus Connection (Optional)

```bash
cd examples/rag/milvus
go run test_connection.go
```

## 🔧 Configuration

Environment variables:

- `MILVUS_URL`: Milvus connection URL (default: `http://localhost:19530`)
- `OLLAMA_URL`: Ollama API URL (default: `http://192.168.10.10:11434`)
- `OLLAMA_MODEL`: Ollama model name (default: `llama3.2`)

## 📊 Automatic Data Setup

This example includes **automatic document setup** that:

1. **Creates Collection**: Automatically creates the `knowledge_base` collection with proper schema
2. **Adds Knowledge Documents**: Inserts 8 AI-related knowledge documents with embeddings
3. **Verifies Data**: Ensures documents are accessible before running SwarmV2 agents
4. **Uses Official SDK**: Leverages the official Milvus Go SDK for all operations

### Built-in Knowledge Documents

The example automatically adds documents covering:
- **Artificial Intelligence**: Core AI concepts and applications
- **Machine Learning**: ML algorithms and pattern recognition
- **Deep Learning**: Neural networks and complex data processing
- **Natural Language Processing**: Text understanding and generation
- **Computer Vision**: Image analysis and interpretation
- **Vector Databases**: Milvus and similarity search capabilities
- **Robotics**: AI-powered autonomous systems
- **Data Science**: Insights and knowledge extraction techniques

### Data Schema

The `knowledge_base` collection uses:
- **ID**: VarChar primary key (UUID-based)
- **Vector**: 128-dimensional FloatVector embeddings
- **Content**: Text content with AI knowledge documents
- **Metadata**: Document type, topic, and indexing information

## 🧪 Example Queries

The example tests these query types:

- "What is artificial intelligence?"
- "Explain machine learning and deep learning"
- "How does Milvus work for vector search?"
- "What are the applications of computer vision?"
- "Tell me about natural language processing"

## 🏛️ SwarmV2 Architecture

### Components

1. **MilvusRAGRetriever**
   - Performs vector search using official Milvus Go SDK
   - Retrieves relevant knowledge context from 128-dimensional vectors
   - Supports metadata filtering and automatic document setup

2. **OllamaRAGGenerator**
   - Generates responses using Ollama LLM
   - Context-aware prompting with retrieved knowledge
   - Handles empty/invalid queries gracefully

3. **Document Setup System**
   - Automatically creates collections with compatible schema
   - Verifies document accessibility before SwarmV2 execution
   - Error handling for missing or inaccessible data

### Design Pattern

```go
// Document Setup (automatic)
err := setupAndVerifyDocuments(ctx, retriever, ragStore, collection)

// Retrieval phase
context, err := retriever.Retrieve(ctx, query)

// Generation phase  
response, err := generator.GenerateRAGResponse(query, context)
```

## 🛠️ Development

### Build and Run

```bash
# From SwarmV2 root
go build -o milvus_rag examples/rag/milvus/main.go
./milvus_rag
```

### Test Connection

```bash
# Test Milvus connectivity and collection schema
cd examples/rag/milvus
go run test_connection.go
```

### Run Test Suite

```bash
# Run comprehensive test checks
./test.sh
```

### Debug

Check these if issues occur:

1. **Milvus Connection**: `curl http://localhost:19530/healthz`
2. **Ollama Status**: `curl http://192.168.10.10:11434/api/version`
3. **Collection Schema**: Use `test_connection.go` to verify collection structure
4. **Go SDK Version**: Ensure Milvus Go SDK v2 is properly installed

## 📈 Expected Output

When working correctly, you should see:

```
=== Milvus RAG Example with SwarmV2 ===
🚀 Initializing Milvus RAG system...
💡 This example includes automatic document setup!
✅ Successfully connected to Milvus!
✅ Ollama connection successful!
📋 Setting up and verifying documents in Milvus collection 'knowledge_base'...
📝 Adding 8 documents using Go Milvus provider...
✅ Added document 1 with ID: <uuid>
...
✅ Successfully added 8 documents to collection
✅ Document verification successful!
🔍 Query 1: What is artificial intelligence?
📄 Retrieved context (XXX characters)
🤖 Generated Response:
[Detailed AI explanation based on knowledge documents...]
```

## 🔗 Related Files

- `../../../src/vectordb/providers/milvus_sdk.go`: Official Milvus Go SDK provider
- `../../../src/agents/rag/`: SwarmV2 RAG agent implementations
- `test_connection.go`: Milvus connectivity testing utility

## 💡 Tips

1. **Automatic setup** - No need to run Python scripts or manual data insertion
2. **Check test_connection.go** to verify Milvus connectivity and collection schema
3. **Use environment variables** for different configurations
4. **Monitor logs** for connection, document setup, and retrieval details
5. **Official SDK** - Uses the official Milvus Go SDK for optimal performance

This example provides a **complete, self-contained** RAG system ready for production use! 🎯
