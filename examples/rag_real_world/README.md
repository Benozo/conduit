# TechCorp Knowledge Management - Real World RAG Example

## 🏢 Overview

This example demonstrates a production-ready RAG (Retrieval Augmented Generation) system for a fictional company "TechCorp". It showcases how ConduitMCP can be used to build an enterprise knowledge management system that allows employees to query company policies, procedures, and documentation using natural language.

## 📋 Scenario

**TechCorp Knowledge Management System**

You are implementing a knowledge base for TechCorp, a technology company that needs to make company information easily accessible to employees. The system indexes various company documents and allows employees to ask questions about:

- HR policies and procedures
- Remote work guidelines  
- Security protocols and data protection
- Customer onboarding processes
- Expense reimbursement policies

## 🎯 Business Use Cases

The example demonstrates realistic business scenarios:

1. **HR Manager**: "What are the steps for onboarding a new software engineer?"
2. **Project Manager**: "What is our company policy on remote work and flexible hours?"
3. **Developer**: "What security practices should I follow when handling customer data?"
4. **Sales Team**: "How do we onboard new enterprise customers?"
5. **Finance**: "What expenses can I claim and what's the approval process?"

## 🚀 Features Demonstrated

### Core RAG Capabilities
- **Document Indexing**: Automatically chunks and embeds company documents
- **Semantic Search**: Find relevant information using natural language queries
- **AI-Powered Answers**: Generate contextual responses based on company knowledge
- **Source Attribution**: Track which documents were used to generate answers

### Advanced Features
- **Filtered Search**: Search within specific departments, categories, or time periods
- **Multiple Embedding Providers**: Support for both OpenAI and Ollama embeddings
- **Metadata-Rich Storage**: Store additional context like department, category, year
- **Statistics and Analytics**: Track knowledge base usage and content metrics

### Production-Ready Features
- **Error Handling**: Robust error handling and graceful degradation
- **Timeouts**: Proper timeout handling for all operations
- **Health Checks**: Embedding provider connectivity validation
- **Memory Optimization**: Efficient chunk size and overlap strategies

## 🛠️ Prerequisites

### 1. Database Setup
Start PostgreSQL with pgvector using Docker Compose:

```bash
# From the root directory
docker-compose up -d postgres adminer
```

Wait for the database to be ready (check logs):
```bash
docker-compose logs postgres
```

### 2. Embedding Provider

**Option A: Ollama (Recommended for local development)**
```bash
# Install and start Ollama
ollama serve

# Pull the embedding model
ollama pull nomic-embed-text
```

**Option B: OpenAI**
```bash
export OPENAI_API_KEY="your-api-key-here"
```

## 🏃 Running the Example

### Quick Start (Ollama)
```bash
cd examples/rag_real_world
go run main.go
```

### Using OpenAI
```bash
cd examples/rag_real_world
export RAG_PROVIDER=openai
export OPENAI_API_KEY="your-api-key-here"
go run main.go
```

### Custom Configuration
```bash
# Custom Ollama host
export OLLAMA_HOST="http://custom-host:11434"

# Custom database connection
export RAG_DB_HOST="custom-db-host"
export RAG_DB_PORT="5432"
export RAG_DB_NAME="rag_db"
export RAG_DB_USER="postgres"
export RAG_DB_PASSWORD="password"

go run main.go
```

## 📊 Example Output

```
🏢 ConduitMCP RAG - Real World Business Example
===============================================
Scenario: TechCorp Knowledge Management System
Indexing: Company policies, procedures, and documentation

📊 Using OLLAMA embeddings
🔧 Chunk size: 800 chars, Overlap: 150 chars

🚀 Initializing TechCorp Knowledge Base...
📚 Current knowledge base: 0 documents

📝 Indexing TechCorp Company Documents...
   📄 Indexing: Employee Handbook 2024 (1/5)
   ✅ Successfully indexed Employee Handbook 2024
   📄 Indexing: Remote Work Policy (2/5)
   ✅ Successfully indexed Remote Work Policy
   ...

🔍 Real-World Business Query Examples
=====================================

1. 👤 HR Manager - New Employee Onboarding
   ❓ "What are the steps for onboarding a new software engineer?"
   🔍 Found 3 relevant documents
      1. Score: 0.847 | ## Onboarding Process for New Employees ### Week 1: Get...
      2. Score: 0.792 | ### Software Engineers Onboarding - Complete security...
      3. Score: 0.734 | ### Required Training - Company culture and values...
   🤖 AI Answer (Confidence: 0.89):
      For onboarding a new software engineer at TechCorp, follow this structured 
      process: Week 1 begins with Day 1 HR paperwork and equipment setup, Days 2-3 
      include IT setup with laptop, accounts, and security training, and Days 4-5 
      cover department introduction and role-specific training...
   📄 Sources: 2 documents referenced

2. 👤 Project Manager - Remote Work Policy
   ❓ "What is our company policy on remote work and flexible hours?"
   🔍 Found 3 relevant documents
   ...
```

## 🎯 Advanced Features

### Filtered Search Examples

The example demonstrates sophisticated filtering capabilities:

```go
// Search only HR documents
filters := map[string]interface{}{"department": "HR"}
results, err := ragEngine.Search(ctx, "vacation policy", 5, filters)

// Search security-related documents
filters := map[string]interface{}{"category": "security"}
results, err := ragEngine.Search(ctx, "data protection", 5, filters)

// Search recent documents
filters := map[string]interface{}{"year": 2024}
results, err := ragEngine.Search(ctx, "company policies", 5, filters)
```

### Custom Document Metadata

Each document includes rich metadata for enhanced search:

```go
metadata := map[string]interface{}{
    "category":   "HR",
    "department": "Human Resources", 
    "year":       2024,
    "type":       "policy",
    "indexed_at": "2024-01-15T10:30:00Z",
}
```

## 📈 Business Impact

This example demonstrates how RAG systems provide tangible business value:

- **Instant Access**: Employees get immediate answers to policy questions
- **Consistency**: Everyone gets the same, accurate information from official sources
- **Efficiency**: Reduced time spent searching through documents and manuals
- **Onboarding**: New employees can quickly find answers to common questions
- **Compliance**: Better adherence to policies through improved accessibility
- **Scalability**: Knowledge scales automatically as new documents are added

## 🔧 Customization

### Adding Your Own Documents

Replace the sample documents with your actual company content:

```go
documents := []CompanyDocument{
    {
        Title:    "Your Company Policy",
        Content:  "Your policy content here...",
        Category: "your_category",
        Department: "Your Department",
        Year:     2024,
        Type:     "policy",
    },
    // Add more documents...
}
```

### Custom Chunking Strategy

Optimize chunking for your document types:

```go
// For legal documents (larger chunks)
config.Chunking.Size = 1200
config.Chunking.Overlap = 200

// For technical documentation (smaller chunks)
config.Chunking.Size = 600
config.Chunking.Overlap = 100

// For FAQ-style content
config.Chunking.Strategy = "sentence"
```

### Integration with Your Systems

The example can be extended to integrate with:

- **Document Management Systems**: Automatically index new documents
- **HR Systems**: Pull employee-specific information
- **Slack/Teams Bots**: Provide answers in chat platforms
- **Web Portals**: Build employee self-service portals
- **Analytics**: Track usage patterns and popular queries

## 🚀 Next Steps

1. **Index Real Documents**: Replace sample content with your actual company documents
2. **Automated Updates**: Set up processes to automatically index new/updated documents
3. **User Training**: Train employees on how to effectively query the knowledge base
4. **Analytics**: Implement usage tracking to understand popular queries and gaps
5. **Integration**: Connect with existing tools like Slack, SharePoint, or your intranet
6. **Advanced Search**: Add support for boolean queries, date ranges, and complex filters
7. **Multilingual Support**: Extend for multiple languages if needed
8. **Access Controls**: Implement role-based access to sensitive documents

## 📚 Related Examples

- **Basic RAG**: `examples/rag/` - Simple RAG implementation and concepts
- **OpenAI Integration**: `examples/openai/` - OpenAI-specific features
- **Ollama Integration**: `examples/ollama/` - Local Ollama setup and usage

## 🤝 Contributing

To improve this example:

1. Add more realistic business scenarios
2. Enhance document content with real-world complexity
3. Implement additional filtering and search capabilities
4. Add performance benchmarks and optimization tips
5. Create integration examples with popular business tools

## 📖 Documentation

For more information about the underlying RAG system:

- [RAG System Architecture](../../ROADMAP_RAG_PGVECTOR.md)
- [Database Setup Guide](../../DOCKER_SETUP.md)
- [API Documentation](../../lib/rag/)
