# ConduitMCP RAG API Server

A REST API server that provides RAG (Retrieval Augmented Generation) capabilities with chat and document management endpoints.

## Features

🌐 **REST API**: HTTP endpoints for chat and document management  
🧠 **Multi-provider Support**: Works with OpenAI and Ollama embeddings  
📊 **Vector Search**: PostgreSQL with pgvector for semantic similarity  
🔧 **MCP Tools**: Integrated tool calling capabilities  
⚡ **Optimized Performance**: Efficient indexing and resource management  
🛡️ **Production Ready**: Health checks, CORS support, graceful shutdown

## Quick Start

### Prerequisites

1. **PostgreSQL with pgvector**:
   ```bash
   # From project root
   docker compose up -d
   ```

2. **Ollama** (for local embeddings):
   ```bash
   # Install and pull embedding model
   ollama pull nomic-embed-text:latest
   ollama pull llama3.2
   ```

### Running the Server

```bash
# Default configuration (Ollama + PostgreSQL)
go run main.go

# Custom port
PORT=8090 go run main.go

# OpenAI configuration
RAG_PROVIDER=openai OPENAI_API_KEY=your_key go run main.go
```

## API Endpoints

### 🏥 Health Check
```bash
GET /health
# Response: {"status": "healthy", "message": "RAG system is operational"}
```

### 📊 Knowledge Base Stats
```bash
GET /stats
# Response: {
#   "document_count": 8,
#   "chunk_count": 64,
#   "embedding_model": "nomic-embed-text:latest",
#   "embedding_dimensions": 768
# }
```

### 💬 Chat with Knowledge Base
```bash
POST /chat
Content-Type: application/json

{
  "message": "What is ConduitMCP?",
  "limit": 5  # optional, default: 5
}

# Response: {
#   "response": "Based on the knowledge base: ...",
#   "sources": [...],
#   "response_time": "117ms"
# }
```

### 📄 Add Documents
```bash
POST /documents
Content-Type: application/json

{
  "content": "# Document Title\nDocument content here...",
  "title": "My Document",
  "type": "text/plain",  # optional
  "metadata": {          # optional
    "category": "docs",
    "version": "1.0"
  }
}

# Response: {
#   "document_id": "uuid-here",
#   "message": "Document 'My Document' successfully indexed"
# }
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RAG_PROVIDER` | `ollama` | Embedding provider: `openai` or `ollama` |
| `OPENAI_API_KEY` | - | Required for OpenAI provider |
| `OLLAMA_HOST` | `localhost` | Ollama server host |
| `OLLAMA_MODEL` | `nomic-embed-text:latest` | Ollama embedding model |
| `POSTGRES_HOST` | `localhost` | PostgreSQL host |
| `POSTGRES_USER` | `conduit` | PostgreSQL user |
| `POSTGRES_PASSWORD` | `conduit_password` | PostgreSQL password |
| `POSTGRES_DB` | `conduit_rag` | PostgreSQL database |

## Available MCP Tools

The example registers these tools for programmatic access:

- **`index_document`**: Index documents from files or content
- **`semantic_search`**: Search documents by semantic similarity
- **`knowledge_query`**: Ask questions and get AI-generated answers
- **`list_documents`**: List all indexed documents
- **`get_document`**: Retrieve a specific document
- **`delete_document`**: Remove a document from the index
- **`get_document_chunks`**: View document chunks
- **`get_rag_stats`**: Get system statistics

## What the Example Does

1. **Initialization**: Sets up database, embeddings, chunker, and RAG engine
2. **Health Check**: Verifies all components are working
3. **Sample Document**: Creates example documentation if database is empty
4. **Semantic Search**: Demonstrates similarity search capabilities
5. **RAG Query**: Shows AI-powered question answering
6. **Statistics**: Displays system metrics

## Example Queries

Try these semantic searches:
- `"embedding providers and their features"`
- `"How does vector search work?"`
- `"production deployment features"`
- `"difference between OpenAI and Ollama"`

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure PostgreSQL is running: `docker compose ps`
   - Check connection settings

2. **Ollama Connection Failed**
   - Verify Ollama is running: `ollama list`
   - Check if `nomic-embed-text:latest` model is available
   - Verify host and port settings

3. **OpenAI Connection Failed**
   - Verify API key is correct
   - Check internet connectivity
   - Ensure you have API credits

4. **Memory Issues**
   - The system uses memory-efficient batch processing
   - Vector dimensions are automatically detected
   - Database schema adapts to embedding provider

### Performance Tips

- Use smaller chunk sizes for faster processing
- Ollama provides good performance for local development
- OpenAI offers higher quality embeddings for production
- Enable connection pooling for high-throughput scenarios

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Documents     │    │   Text Chunker   │    │   Embeddings    │
│                 │───▶│                  │───▶│   (OpenAI/      │
│   - Files       │    │   - Fixed Size   │    │    Ollama)      │
│   - Content     │    │   - Semantic     │    │                 │
│   - Metadata    │    │   - Paragraph    │    │                 │
└─────────────────┘    │   - Sentence     │    └─────────────────┘
                       └──────────────────┘              │
                                                         │
┌─────────────────┐    ┌──────────────────┐    ┌─────────▼─────────┐
│   MCP Tools     │    │   RAG Engine     │    │   PostgreSQL      │
│                 │───▶│                  │───▶│   + pgvector      │
│   - Search      │    │   - Query        │    │                   │
│   - Index       │    │   - Retrieve     │    │   - Documents     │
│   - Manage      │    │   - Generate     │    │   - Chunks        │
└─────────────────┘    └──────────────────┘    │   - Vectors       │
                                               └───────────────────┘
```

This example provides a complete foundation for building RAG-powered applications with ConduitMCP.

## 🎯 Use Cases & Applications

The ConduitMCP RAG pipeline enables a wide variety of intelligent applications. Here are practical use cases with implementation examples:

### 📚 Knowledge Management Systems

**Corporate Knowledge Base**
- Index company policies, procedures, and documentation
- Enable employees to ask natural language questions
- Automatically cite sources and maintain accuracy
- Example: "What's our remote work policy?" → Returns policy with source citations

**Technical Documentation Assistant**
- Index API docs, code repositories, and technical guides
- Help developers find relevant information quickly
- Provide context-aware code examples and best practices
- Example: "How do I implement OAuth?" → Returns relevant docs and code snippets

### 🏥 Healthcare & Medical Applications

**Medical Literature Search**
- Index research papers, clinical guidelines, and medical databases
- Help healthcare professionals find relevant studies
- Support evidence-based decision making
- Example: "Latest treatments for diabetes" → Returns recent research with confidence scores

**Patient Care Assistant**
- Index patient records, treatment protocols, and medical histories
- Assist healthcare providers with treatment recommendations
- Maintain HIPAA compliance with private deployments
- Example: "Treatment options for hypertension in elderly patients"

### 🏛️ Legal & Compliance

**Legal Document Analysis**
- Index contracts, regulations, and case law
- Help lawyers find relevant precedents and clauses
- Support contract review and compliance checking
- Example: "Data privacy clauses in EU contracts" → Returns relevant contract sections

**Regulatory Compliance**
- Index changing regulations and compliance requirements
- Help teams stay updated on regulatory changes
- Automate compliance checking and reporting
- Example: "GDPR requirements for data processing" → Returns specific articles and guidance

### 🎓 Education & Training

**Personalized Learning Assistant**
- Index educational content, textbooks, and course materials
- Provide personalized explanations and examples
- Support different learning styles and paces
- Example: "Explain quantum physics with simple examples" → Tailored explanations

**Corporate Training Platform**
- Index training materials, videos, and certification content
- Help employees find relevant training resources
- Track learning progress and recommend next steps
- Example: "Python programming basics for beginners" → Structured learning path

### 🛒 E-commerce & Customer Support

**Product Information Assistant**
- Index product catalogs, specifications, and reviews
- Help customers find products matching their needs
- Provide detailed comparisons and recommendations
- Example: "Laptop for video editing under $2000" → Filtered recommendations

**Intelligent Customer Support**
- Index support tickets, FAQs, and troubleshooting guides
- Provide instant answers to common questions
- Escalate complex issues to human agents with context
- Example: "WiFi connection problems" → Step-by-step troubleshooting

### 🏢 Business Intelligence

**Market Research Assistant**
- Index market reports, competitor analysis, and industry trends
- Help analysts find relevant insights quickly
- Support strategic decision making with data
- Example: "AI market trends 2024" → Latest market analysis and predictions

**Financial Analysis**
- Index financial reports, earnings calls, and market data
- Help analysts identify investment opportunities
- Provide context-aware financial insights
- Example: "Tech stock performance analysis" → Relevant financial data and trends

### 🔬 Research & Development

**Scientific Research Assistant**
- Index research papers, patents, and experimental data
- Help researchers find related work and avoid duplication
- Support literature reviews and hypothesis generation
- Example: "Machine learning applications in drug discovery" → Relevant research papers

**Innovation Management**
- Index patent databases, R&D projects, and innovation reports
- Help teams identify innovation opportunities
- Support IP strategy and competitive analysis
- Example: "Renewable energy patents 2023" → Patent landscape analysis

### 🏭 Manufacturing & Operations

**Equipment Maintenance Assistant**
- Index maintenance manuals, repair procedures, and troubleshooting guides
- Help technicians find relevant maintenance information
- Support predictive maintenance and downtime reduction
- Example: "Turbine vibration issues" → Diagnostic procedures and solutions

**Quality Control System**
- Index quality standards, inspection procedures, and defect databases
- Help quality teams identify and resolve issues
- Support continuous improvement initiatives
- Example: "Surface finish defects in automotive parts" → Root causes and solutions

### 🌐 Content Management

**Content Discovery Platform**
- Index articles, videos, podcasts, and multimedia content
- Help users discover relevant content based on interests
- Support content curation and recommendation engines
- Example: "Machine learning tutorials for beginners" → Curated learning resources

**Digital Asset Management**
- Index digital assets, brand guidelines, and creative resources
- Help teams find and reuse existing assets
- Maintain brand consistency across organizations
- Example: "Brand logos for social media" → Approved asset variations

## 🛠️ Implementation Patterns

### Basic RAG Pattern
```bash
# 1. Index documents
POST /documents
{"content": "...", "title": "...", "metadata": {...}}

# 2. Query knowledge base
POST /chat
{"message": "Your question here"}
```

### Advanced Filtering
```bash
# Search with metadata filters
POST /chat
{
  "message": "Security protocols",
  "metadata_filter": {"department": "IT", "classification": "internal"}
}
```

### Batch Processing
```bash
# Upload multiple documents
for file in *.pdf; do
  curl -X POST /documents \
    -H "Content-Type: application/json" \
    -d "{\"content\": \"$(cat $file)\", \"title\": \"$file\"}"
done
```

### Integration Examples
```bash
# Slack Bot Integration
curl -X POST /chat \
  -H "Content-Type: application/json" \
  -d "{\"message\": \"${slack_message}\"}" | \
  jq -r '.response' | \
  slack_send_message

# Email Support Integration
incoming_email | extract_question | \
curl -X POST /chat \
  -H "Content-Type: application/json" \
  -d "{\"message\": \"$(cat -)\"}" | \
  jq -r '.response' | \
  send_email_response
```

## 💡 Best Practices by Use Case

### High-Volume Applications
- Use connection pooling for database connections
- Implement caching for frequently accessed content
- Consider horizontal scaling with multiple RAG instances
- Monitor response times and optimize chunk sizes

### Privacy-Sensitive Applications
- Use Ollama for local embeddings (no data sent to external APIs)
- Implement proper access controls and authentication
- Consider on-premises deployment with private networks
- Encrypt data at rest and in transit

### Real-Time Applications
- Pre-compute embeddings for static content
- Use async processing for document indexing
- Implement result caching for common queries
- Optimize database indexes for fast retrieval

### Multi-Language Applications
- Use multilingual embedding models
- Implement language detection and routing
- Consider separate indexes per language
- Translate queries if needed for better matching

## 🚀 Deployment Scenarios

### Development Environment
```bash
# Local Ollama + PostgreSQL
docker compose up -d
go run main.go
```

### Production Environment
```bash
# Scaled deployment with load balancer
docker-compose -f docker-compose.prod.yml up -d
```

### Cloud Deployment
```bash
# Kubernetes deployment
kubectl apply -f k8s/rag-deployment.yaml
```

### Edge Deployment
```bash
# Lightweight deployment for edge computing
RAG_PROVIDER=ollama OLLAMA_MODEL=all-minilm go run main.go
```
