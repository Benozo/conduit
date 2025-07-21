# RAG + LLM Interactive Chat Terminal

🤖 **Interactive terminal chat combining RAG (Retrieval Augmented Generation) with Ollama LLM and MCP tools**

This example demonstrates a complete conversational AI system that can:
- Chat with TechCorp's knowledge base using semantic search
- Generate AI-powered answers using Ollama Llama 3.2
- Use MCP tools for text processing, memory, and utilities
- Provide an interactive terminal-based chat experience

## Features

### 🧠 **Intelligent Conversational AI**
- **Llama 3.2 Integration**: Uses Ollama for natural language understanding and generation
- **RAG-Enhanced Responses**: Searches company documents before generating answers
- **Context-Aware**: Maintains conversation context and user preferences
- **Tool Integration**: Access to 15+ MCP tools through natural language

### 📚 **Knowledge Base Integration**
- **Semantic Search**: Find relevant information using vector similarity
- **Document Indexing**: Automatically indexes TechCorp business documents
- **Source Attribution**: Cites sources and confidence scores
- **Filtered Search**: Search by department, category, or document type

### 🛠️ **MCP Tools Library Mode**
- **Text Processing**: uppercase, lowercase, word_count, reverse, trim
- **Memory System**: remember, recall, clear_memory, list_memories
- **Utilities**: timestamp, uuid, hash functions, base64 encoding
- **RAG Tools**: semantic_search, knowledge_query, list_documents

### 💬 **Interactive Terminal Interface**
- **Real-time Chat**: Type questions and get immediate responses
- **Special Commands**: `/help`, `/stats`, `/search`, `/quit`
- **Response Timing**: Shows how long each response takes
- **Graceful Shutdown**: Ctrl+C handling

## Prerequisites

### 1. **PostgreSQL + pgvector**
```bash
# Start the database (from project root)
docker compose up -d postgres
```

### 2. **Ollama Setup**
```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama server
ollama serve

# Pull required models
ollama pull llama3.2                    # For LLM reasoning
ollama pull nomic-embed-text:latest     # For embeddings
```

### 3. **Verify Setup**
```bash
# Check Ollama is running
curl http://localhost:11434/api/tags

# Check database
docker exec -it conduit-postgres psql -U conduit -d conduit_rag -c "SELECT 1;"
```

## Running the Example

### Default Configuration
```bash
cd /home/engineone/Projects/AI/ConduitMCP
go run examples/rag_chat_terminal/main.go
```

### Custom Configuration
```bash
# Set custom Ollama host
export OLLAMA_HOST="192.168.10.10:11434"

# Set custom models
export OLLAMA_LLM_MODEL="llama3.2"
export OLLAMA_EMBEDDING_MODEL="nomic-embed-text:latest"

# Run with custom settings
go run examples/rag_chat_terminal/main.go
```

## Example Conversations

### 1. **Company Policy Questions**
```
💬 You: What is our remote work policy?

🤖 Assistant: Based on TechCorp's Remote Work Policy, here are the key points:

**Eligibility:**
- Employees who have completed 6 months of employment
- Roles that don't require physical presence
- Approval from direct manager required

**Remote Work Options:**
1. Hybrid Remote: 2-3 days per week from home
2. Fully Remote: Permanent remote work arrangement
3. Temporary Remote: Short-term arrangements

The application process involves completing a Remote Work Request Form, discussing with your manager, HR review, IT equipment assessment, and a 30-day trial period.

⏱️ Response time: 2.3s
```

### 2. **Tool Integration**
```
💬 You: Remember that I work in the engineering team

🤖 Assistant: I've remembered that you work in the engineering team. This will help me provide more relevant information about engineering-specific policies and procedures.

⏱️ Response time: 1.1s
```

### 3. **Complex Queries**
```
💬 You: What are the onboarding steps for a new software engineer and can you generate a UUID for tracking?

🤖 Assistant: Here are the onboarding steps for new software engineers:

**Week 1: Getting Started**
1. Day 1: Complete HR paperwork and receive equipment
2. Day 2-3: IT setup including laptop, accounts, and security training
3. Day 4-5: Department introduction and role-specific training

**Software Engineer Specific:**
- Complete security awareness training
- Set up development environment
- Review coding standards and practices
- Assign mentor for first 30 days

I've also generated a tracking UUID for you: `f47ac10b-58cc-4372-a567-0e02b2c3d479`

⏱️ Response time: 3.7s
```

## Available Commands

| Command | Description |
|---------|-------------|
| `/help` or `/h` | Show help and examples |
| `/quit` or `/q` | Exit the chat |
| `/stats` | Show knowledge base statistics |
| `/search <query>` | Direct semantic search (no AI) |
| `/clear` | Clear the terminal screen |
| `/tasks` | Show recent task information |

## Architecture

```
User Input
    ↓
LLM Agent (Llama 3.2)
    ↓
Task Planning & Tool Selection
    ↓
┌─────────────────┬─────────────────┐
│   RAG System    │   MCP Tools     │
│   ┌─────────┐   │   ┌─────────┐   │
│   │ Search  │   │   │  Text   │   │
│   │ Query   │   │   │ Memory  │   │
│   │ Index   │   │   │ Utility │   │
│   └─────────┘   │   └─────────┘   │
└─────────────────┴─────────────────┘
    ↓
Response Generation
    ↓
Terminal Output
```

## Configuration Options

### Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `OLLAMA_HOST` | `192.168.10.10:11434` | Ollama server host:port |
| `OLLAMA_LLM_MODEL` | `llama3.2` | Model for chat and reasoning |
| `OLLAMA_EMBEDDING_MODEL` | `nomic-embed-text:latest` | Model for embeddings |

### Recommended Models
- **Fast**: `llama3.2:1b` (1B parameters, quick responses)
- **Balanced**: `llama3.2` (3B parameters, good quality)
- **High Quality**: `llama3.1:8b` (8B parameters, best results)
- **Code-focused**: `codellama:7b` (for technical queries)

## Troubleshooting

### Connection Issues
```bash
# Check Ollama
curl http://localhost:11434/api/tags

# Check database
docker ps | grep postgres

# Test embeddings
curl -X POST http://localhost:11434/api/embeddings \
  -H "Content-Type: application/json" \
  -d '{"model": "nomic-embed-text:latest", "prompt": "test"}'
```

### Performance Issues
```bash
# Use smaller model for faster responses
export OLLAMA_LLM_MODEL="llama3.2:1b"

# Check system resources
ollama ps
docker stats conduit-postgres
```

### Model Issues
```bash
# List available models
ollama list

# Pull missing models
ollama pull llama3.2
ollama pull nomic-embed-text:latest

# Remove unused models
ollama rm old-model-name
```

## Advanced Usage

### Custom Knowledge Base
```go
// Add your own documents
documents := []struct {
    title   string
    content string
}{
    {"My Company Policy", "Custom policy content..."},
    {"Technical Documentation", "API documentation..."},
}
```

### Custom Tools
```go
// Register custom MCP tools
server.RegisterToolWithSchema("custom_tool",
    func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
        // Your custom logic
        return result, nil
    }, conduit.ToolMetadata{
        Name: "custom_tool",
        Description: "Description of your tool",
        // ... schema definition
    })
```

### Integration with Your App
```go
// Use the chat system in library mode
import "github.com/benozo/conduit/examples/rag_chat_terminal"

// Create and configure the system
chatSystem := NewRAGChatSystem(config)
response := chatSystem.ProcessQuery("Your question")
```

## Performance Metrics

- **Average Response Time**: 1-4 seconds
- **Knowledge Base**: 5 documents, 110 chunks
- **Vector Dimensions**: 768 (Ollama embeddings)
- **Memory Usage**: ~200MB (with Llama 3.2)
- **Concurrent Users**: Single-user terminal interface

## Next Steps

1. **Scale Knowledge Base**: Add your company's actual documents
2. **Custom Tools**: Integrate domain-specific tools
3. **Web Interface**: Build a web UI for multi-user access
4. **API Integration**: Connect to external systems
5. **Advanced RAG**: Implement hierarchical search, re-ranking

---

🎉 **Result**: A complete interactive chat system that combines the power of local LLMs, semantic search, and tool integration for intelligent business knowledge assistance!
