# ConduitMCP RAG Pipeline Architecture

## Complete RAG Chat System Flow

```mermaid
graph TB
    %% User Interface Layer
    User[👤 User Input] --> Terminal[🖥️ Interactive Terminal]
    Terminal --> InputParser{📝 Parse Input}
    
    %% Special Commands Branch
    InputParser -->|Special Commands| SpecialCmd[🔧 Special Commands]
    SpecialCmd --> Help[/help - Show Help]
    SpecialCmd --> Stats[/stats - KB Stats]
    SpecialCmd --> Search[/search - Direct Search]
    SpecialCmd --> Quit[/quit - Exit]
    
    %% Main RAG Pipeline Branch
    InputParser -->|Natural Language Query| TaskCreator[📋 Create Task]
    
    %% Task Processing
    TaskCreator --> AgentManager[🤖 LLM Agent Manager]
    AgentManager --> LLMAnalysis[🧠 LLM Analysis & Planning]
    
    %% LLM Planning
    LLMAnalysis --> PlanSteps{📊 Plan Tool Usage}
    PlanSteps --> SemanticSearchPlan[🔍 Plan: Semantic Search]
    PlanSteps --> KnowledgeQueryPlan[❓ Plan: Knowledge Query]
    
    %% Tool Execution Layer
    SemanticSearchPlan --> ToolRegistry[🛠️ MCP Tool Registry]
    KnowledgeQueryPlan --> ToolRegistry
    
    ToolRegistry --> SemanticSearchTool[🔍 Semantic Search Tool]
    ToolRegistry --> KnowledgeQueryTool[❓ Knowledge Query Tool]
    ToolRegistry --> OtherTools[🔧 Other MCP Tools]
    
    %% RAG Engine Components
    SemanticSearchTool --> RAGEngine[🧠 RAG Engine]
    KnowledgeQueryTool --> RAGEngine
    
    RAGEngine --> EmbeddingProvider[🔢 Ollama Embeddings]
    RAGEngine --> VectorDB[(📊 PostgreSQL + pgvector)]
    RAGEngine --> Chunker[📄 Text Chunker]
    
    %% Knowledge Base
    VectorDB --> IndexedDocs[📚 Indexed Documents]
    IndexedDocs --> EmployeeHandbook[📖 Employee Handbook 2024]
    IndexedDocs --> RemoteWorkPolicy[🏠 Remote Work Policy]
    IndexedDocs --> SecurityGuidelines[🔐 Data Security Guidelines]
    IndexedDocs --> CustomerOnboarding[👥 Customer Onboarding Process]
    IndexedDocs --> ExpensePolicy[💰 Expense Reimbursement Policy]
    
    %% Search Processing
    EmbeddingProvider --> OllamaHost[🦙 Ollama Host<br/>192.168.10.10:11434]
    OllamaHost --> EmbeddingModel[🔢 nomic-embed-text:latest]
    
    %% Vector Search
    VectorDB --> SimilaritySearch[📐 Cosine Similarity Search]
    SimilaritySearch --> RankedResults[📊 Ranked Results<br/>Score + Content + Source]
    
    %% Response Generation
    RankedResults --> ResponseFormatter[🎨 Response Formatter]
    ResponseFormatter --> SourceGrouping[📋 Group by Source Document]
    SourceGrouping --> PrettyOutput[✨ Formatted Output]
    
    %% LLM Integration for Knowledge Query
    KnowledgeQueryTool --> LLMGeneration[🤖 LLM Answer Generation]
    LLMGeneration --> OllamaLLM[🦙 Ollama Llama 3.2]
    OllamaLLM --> GeneratedAnswer[📝 AI-Generated Answer]
    
    %% Final Output
    PrettyOutput --> FinalResponse[📤 Final Response]
    GeneratedAnswer --> FinalResponse
    FinalResponse --> Terminal
    
    %% Memory System
    AgentManager --> Memory[🧠 Agent Memory]
    Memory --> ContextStorage[💾 Session Context]
    Memory --> RAGEngineRef[🔗 RAG Engine Reference]
    
    %% Configuration
    Config[⚙️ Configuration] --> OllamaConfig[🦙 Ollama Settings]
    Config --> RAGConfig[🔧 RAG Settings]
    Config --> AgentConfig[🤖 Agent Settings]
    
    %% Styling
    classDef userLayer fill:#e1f5fe
    classDef llmLayer fill:#f3e5f5
    classDef ragLayer fill:#e8f5e8
    classDef dataLayer fill:#fff3e0
    classDef toolLayer fill:#fce4ec
    
    class User,Terminal,InputParser userLayer
    class AgentManager,LLMAnalysis,OllamaLLM,LLMGeneration llmLayer
    class RAGEngine,EmbeddingProvider,SimilaritySearch,Chunker ragLayer
    class VectorDB,IndexedDocs,EmployeeHandbook,RemoteWorkPolicy,SecurityGuidelines,CustomerOnboarding,ExpensePolicy dataLayer
    class ToolRegistry,SemanticSearchTool,KnowledgeQueryTool,OtherTools toolLayer
```

## Data Flow Sequence

```mermaid
sequenceDiagram
    participant U as 👤 User
    participant T as 🖥️ Terminal
    participant AM as 🤖 Agent Manager
    participant LLM as 🦙 Llama 3.2
    participant TR as 🛠️ Tool Registry
    participant RE as 🧠 RAG Engine
    participant VDB as 📊 Vector DB
    participant OE as 🔢 Ollama Embeddings
    
    U->>T: "What is our remote work policy?"
    T->>AM: Create Task with user query
    AM->>LLM: Analyze task & plan tool usage
    LLM->>AM: JSON plan: [semantic_search, knowledge_query]
    
    AM->>TR: Execute semantic_search tool
    TR->>RE: Search(query="remote work policy", limit=5)
    RE->>OE: Generate embedding for query
    OE->>RE: Query embedding vector
    RE->>VDB: Similarity search with embedding
    VDB->>RE: Top 5 relevant chunks with scores
    RE->>TR: Formatted search results
    TR->>AM: Tool execution complete
    
    AM->>TR: Execute knowledge_query tool
    TR->>RE: Query(question="What is our remote work policy?")
    RE->>VDB: Retrieve relevant context
    VDB->>RE: Context chunks
    RE->>LLM: Generate answer with context
    LLM->>RE: AI-generated answer
    RE->>TR: Complete answer with sources
    TR->>AM: Tool execution complete
    
    AM->>T: Task completed
    T->>T: Format results by source document
    T->>U: 📚 Organized response with sources
    
    Note over U,OE: Response Time: ~3 seconds
```

## Architecture Components

```mermaid
graph LR
    subgraph "🖥️ Interface Layer"
        CLI[Command Line Interface]
        Help[Help System]
        Stats[Statistics Display]
    end
    
    subgraph "🤖 Agent Layer"
        LLMAgent[LLM Agent]
        TaskManager[Task Manager]
        Memory[Memory System]
        SystemPrompt[System Prompt]
    end
    
    subgraph "🛠️ Tool Layer"
        MCPTools[MCP Tool Registry]
        SemanticSearch[semantic_search]
        KnowledgeQuery[knowledge_query]
        ListDocs[list_documents]
        TextTools[Text Processing Tools]
        UtilityTools[Utility Tools]
    end
    
    subgraph "🧠 RAG Layer"
        RAGEngine[RAG Engine]
        TextChunker[Text Chunker]
        EmbeddingProvider[Embedding Provider]
        QueryProcessor[Query Processor]
    end
    
    subgraph "💾 Data Layer"
        VectorDB[(PostgreSQL + pgvector)]
        Documents[Document Store]
        Embeddings[Vector Embeddings]
        Metadata[Document Metadata]
    end
    
    subgraph "🦙 Ollama Services"
        LlamaModel[Llama 3.2 LLM]
        NomicEmbeddings[nomic-embed-text]
        OllamaAPI[Ollama API Server]
    end
    
    CLI --> LLMAgent
    LLMAgent --> MCPTools
    MCPTools --> RAGEngine
    RAGEngine --> VectorDB
    RAGEngine --> EmbeddingProvider
    EmbeddingProvider --> OllamaAPI
    LLMAgent --> OllamaAPI
    
    classDef interface fill:#e1f5fe
    classDef agent fill:#f3e5f5
    classDef tool fill:#fce4ec
    classDef rag fill:#e8f5e8
    classDef data fill:#fff3e0
    classDef ollama fill:#ffeb3b
    
    class CLI,Help,Stats interface
    class LLMAgent,TaskManager,Memory,SystemPrompt agent
    class MCPTools,SemanticSearch,KnowledgeQuery,ListDocs,TextTools,UtilityTools tool
    class RAGEngine,TextChunker,EmbeddingProvider,QueryProcessor rag
    class VectorDB,Documents,Embeddings,Metadata data
    class LlamaModel,NomicEmbeddings,OllamaAPI ollama
```

## Key Features

### 🎯 **Core Capabilities**
- **Interactive Chat**: Real-time conversation with TechCorp knowledge base
- **Semantic Search**: Vector-based similarity search across 110 document chunks
- **AI-Powered Responses**: LLM-generated answers with source citations
- **Tool Integration**: 12+ MCP tools for enhanced functionality
- **Memory System**: Context preservation across conversation

### 📊 **Performance Metrics**
- **Knowledge Base**: 5 documents, 110 chunks indexed
- **Response Time**: ~3 seconds average
- **Embedding Model**: nomic-embed-text (384 dimensions)
- **LLM Model**: Llama 3.2 via Ollama
- **Vector DB**: PostgreSQL with pgvector extension

### 🔧 **Technical Stack**
- **Language**: Go
- **Vector Database**: PostgreSQL + pgvector
- **Embeddings**: Ollama nomic-embed-text
- **LLM**: Ollama Llama 3.2
- **Tools Framework**: Model Context Protocol (MCP)
- **Agent System**: Custom LLM-powered agents

This RAG pipeline successfully combines retrieval-augmented generation with intelligent agent orchestration to provide accurate, contextual responses about company policies and procedures.
