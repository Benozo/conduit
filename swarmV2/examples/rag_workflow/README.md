# Enhanced RAG Workflow with Ollama Integration

This example demonstrates the RAG (Retrieve, Augment, Generate) workflow pattern with AI-powered content generation using Ollama.

## What it demonstrates

- **Retriever Agent**: Retrieves relevant information from a knowledge base
- **Ollama Generator**: AI-powered content generation using llama3.2 model
- **Evaluator Agent**: Evaluates the quality of generated content
- **Custom RAG Workflow**: Orchestrates the interaction with real AI integration
- **Step-by-step Processing**: Clear visualization of each RAG stage

## Prerequisites

- Ollama server at `192.168.10.10:11434`
- `llama3.2` model available
- Traditional framework components

```bash
# Pull the required model
ollama pull llama3.2

# Verify model is available
ollama list
```

## How to run

```bash
cd /home/engineone/Projects/AI/ConduitMCP/swarmV2/examples/rag_workflow
go run main.go
```

## Expected output

The example will:
1. Test connection to Ollama server
2. Initialize the three RAG agents (including Ollama-powered generator)
3. Create a custom RAG workflow with AI integration
4. Execute the workflow with a comprehensive query about machine learning
5. Show step-by-step processing:
   - Information retrieval from knowledge base
   - AI-powered content generation using Ollama
   - Quality evaluation of generated content
6. Display the final AI-generated content and workflow summary

## Architecture

```
Enhanced RAG Workflow
├── 🔍 Retriever (Traditional)
│   └── Retrieves from knowledge base
├── 🤖 OllamaGenerator (AI-Powered)
│   └── llama3.2 model for content generation
└── 📋 Evaluator (Traditional)
    └── Quality assessment and validation
```

## Code structure

- Uses agents from `src/agents/rag/` for retrieval and evaluation
- Uses `src/llm/providers/ollama.go` for AI-powered generation
- Custom RAG workflow implementation for Ollama integration
- Real-time AI content generation with comprehensive prompting
- Demonstrates error handling and result processing
