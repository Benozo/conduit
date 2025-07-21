package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/benozo/conduit/agents"
	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/rag"
	"github.com/benozo/conduit/lib/rag/database"
	"github.com/benozo/conduit/lib/rag/embeddings"
	"github.com/benozo/conduit/lib/rag/processors"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

// Global variables for server components
var (
	ragEngine    rag.RAGEngine
	ragAgent     *agents.LLMAgent
	agentManager *agents.LLMAgentManager
	mcpServer    *conduit.EnhancedServer
)

// Request/Response types
type ChatRequest struct {
	Message string `json:"message"`
	Limit   int    `json:"limit,omitempty"`
}

type ChatResponse struct {
	Response     string                   `json:"response"`
	Sources      []map[string]interface{} `json:"sources"`
	ResponseTime string                   `json:"response_time"`
	Error        string                   `json:"error,omitempty"`
}

type DocumentRequest struct {
	Content  string                 `json:"content"`
	Title    string                 `json:"title"`
	Type     string                 `json:"type,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type DocumentResponse struct {
	DocumentID string `json:"document_id"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
}

type StatsResponse struct {
	DocumentCount       int    `json:"document_count"`
	ChunkCount          int    `json:"chunk_count"`
	EmbeddingModel      string `json:"embedding_model"`
	EmbeddingDimensions int    `json:"embedding_dimensions"`
	Error               string `json:"error,omitempty"`
}

func main() {
	fmt.Println("🚀 ConduitMCP RAG API Server")
	fmt.Println("============================")
	fmt.Println("🌐 REST API with chat and document endpoints")

	// Load configuration
	provider := strings.ToLower(os.Getenv("RAG_PROVIDER"))
	if provider == "" {
		provider = "ollama"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}

	var config *rag.RAGConfig
	switch provider {
	case "openai":
		config = rag.DefaultRAGConfig()
		if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
			config.Embeddings.APIKey = openaiKey
		}
		if config.Embeddings.APIKey == "" {
			log.Fatal("❌ OPENAI_API_KEY environment variable is required for OpenAI provider")
		}
	case "ollama":
		config = rag.DefaultOllamaRAGConfig()
		if ollamaHost := os.Getenv("OLLAMA_HOST"); ollamaHost != "" {
			hostParts := strings.Split(ollamaHost, ":")
			config.Embeddings.Host = hostParts[0]
		}
		if ollamaModel := os.Getenv("OLLAMA_MODEL"); ollamaModel != "" {
			config.Embeddings.Model = ollamaModel
		}
	default:
		log.Fatalf("❌ Unsupported provider: %s", provider)
	}

	// Override database settings
	if dbHost := os.Getenv("POSTGRES_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbUser := os.Getenv("POSTGRES_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("POSTGRES_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbName := os.Getenv("POSTGRES_DB"); dbName != "" {
		config.Database.Name = dbName
	}

	fmt.Printf("📊 Configuration:\n")
	fmt.Printf("  Provider: %s\n", strings.ToUpper(provider))
	fmt.Printf("  Port: %s\n", port)
	fmt.Printf("  Database: %s:%d/%s\n", config.Database.Host, config.Database.Port, config.Database.Name)

	ctx := context.Background()

	// Initialize RAG components
	fmt.Println("\n🔧 Initializing RAG components...")
	if err := initializeRAGSystem(ctx, config, provider); err != nil {
		log.Fatalf("Failed to initialize RAG system: %v", err)
	}

	// Set up HTTP routes
	fmt.Println("\n🌐 Setting up HTTP endpoints...")
	setupRoutes()

	// Start server
	fmt.Printf("\n🎉 RAG API Server starting on port %s\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  POST /chat         - Chat with knowledge base")
	fmt.Println("  POST /documents    - Add documents to knowledge base")
	fmt.Println("  GET  /stats        - Get knowledge base statistics")
	fmt.Println("  GET  /health       - Health check")
	fmt.Println("")

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	fmt.Println("\n👋 Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("✅ Server stopped gracefully")
}

func initializeRAGSystem(ctx context.Context, config *rag.RAGConfig, provider string) error {
	// 1. Initialize PostgreSQL with pgvector
	fmt.Print("  📊 Connecting to PostgreSQL with pgvector... ")
	vectorDB, err := database.NewPgVectorDB(config.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	fmt.Println("✅")

	// 2. Initialize embedding provider
	fmt.Printf("  🧠 Initializing %s embeddings... ", strings.ToUpper(provider))
	var embeddingProvider rag.EmbeddingProvider

	if provider == "openai" {
		embeddingProvider = embeddings.NewOpenAIEmbeddings(
			config.Embeddings.APIKey,
			config.Embeddings.Model,
			config.Embeddings.Dimensions,
			config.Embeddings.Timeout,
		)
	} else {
		embeddingProvider = embeddings.NewOllamaEmbeddings(
			config.Embeddings.Host,
			config.Embeddings.Model,
			config.Embeddings.Dimensions,
			config.Embeddings.Timeout,
		)
	}

	// Test embeddings connection with short timeout
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	if err := embeddingProvider.Ping(pingCtx); err != nil {
		cancel()
		return fmt.Errorf("failed to connect to %s: %v", provider, err)
	}
	cancel()
	fmt.Println("✅")

	// 3. Initialize text chunker with optimized settings for performance
	fmt.Print("  📝 Initializing text chunker... ")
	var strategy processors.ChunkingStrategy
	switch config.Chunking.Strategy {
	case "semantic":
		strategy = processors.Semantic
	case "paragraph":
		strategy = processors.Paragraph
	case "sentence":
		strategy = processors.Sentence
	default:
		strategy = processors.FixedSize
	}

	// Use very small chunk sizes for minimal CPU load
	chunkSize := 256   // Much smaller chunks
	chunkOverlap := 20 // Minimal overlap

	chunker := processors.NewTextChunker(
		strategy,
		chunkSize,
		chunkOverlap,
	)
	fmt.Println("✅")

	// 4. Create RAG engine
	fmt.Print("  🎯 Creating RAG engine... ")
	ragEngine = rag.NewRAGEngine(config, vectorDB, embeddingProvider, chunker)
	fmt.Println("✅")

	// 5. Create MCP server with RAG tools
	fmt.Print("  🛠️  Setting up MCP server with tools... ")
	serverConfig := conduit.DefaultConfig()
	serverConfig.EnableLogging = false
	mcpServer = conduit.NewEnhancedServer(serverConfig)

	// Register standard MCP tools
	tools.RegisterTextTools(mcpServer)
	tools.RegisterMemoryTools(mcpServer)
	tools.RegisterUtilityTools(mcpServer)

	// Register RAG-specific tools
	registerRAGTools(mcpServer)
	fmt.Println("✅")

	// 6. Initialize LLM agent
	fmt.Print("  🤖 Setting up LLM agent... ")
	var ollamaURL string
	ollamaHost := config.Embeddings.Host
	if !strings.HasPrefix(ollamaHost, "http") {
		ollamaURL = "http://" + ollamaHost + ":11434"
	} else {
		ollamaURL = ollamaHost
	}

	llmModel := os.Getenv("OLLAMA_LLM_MODEL")
	if llmModel == "" {
		llmModel = "llama3.2"
	}

	ollamaModel := conduit.CreateOllamaModel(ollamaURL)
	agentManager = agents.NewLLMAgentManager(mcpServer, ollamaModel, llmModel)

	// Create RAG-enabled assistant agent
	var err2 error
	ragAgent, err2 = agentManager.CreateLLMAgent(
		"rag_assistant",
		"Knowledge Assistant",
		"An intelligent assistant that can search knowledge base and answer questions",
		createRAGSystemPrompt(),
		[]string{
			// RAG tools
			"semantic_search", "knowledge_query", "list_documents",
			// Standard MCP tools
			"uppercase", "lowercase", "word_count", "remember", "recall",
			"timestamp", "uuid", "base64_encode", "base64_decode",
		},
		&agents.AgentConfig{
			MaxTokens:     2000,
			Temperature:   0.3,
			EnableMemory:  true,
			EnableLogging: false,
		},
	)

	if err2 != nil {
		return fmt.Errorf("failed to create RAG agent: %v", err2)
	}

	// Set up RAG context in memory
	ragAgent.Memory.Set("rag_engine", ragEngine)
	ragAgent.Memory.Set("context", "API user chatting with knowledge base")
	fmt.Println("✅")

	// 7. Check and populate knowledge base (with skip option)
	fmt.Println("  📚 Checking knowledge base...")
	stats, err := ragEngine.GetStats(ctx)
	if err != nil {
		log.Printf("Warning: Could not get stats: %v", err)
		stats = make(map[string]interface{})
	}

	existingDocs := getIntValue(stats, "document_count")
	fmt.Printf("  📊 Current: %d documents, %d chunks\n",
		existingDocs, getIntValue(stats, "chunk_count"))

	// Only index if explicitly requested via environment variable
	skipIndexing := os.Getenv("SKIP_INDEXING")
	if existingDocs == 0 && skipIndexing != "true" {
		fmt.Println("  📝 Indexing minimal sample documents...")
		fmt.Println("  💡 Set SKIP_INDEXING=true to skip this step")
		if err := indexSampleDocuments(ctx, ragEngine, provider); err != nil {
			log.Printf("Warning: Failed to index documents: %v", err)
		}
	} else if skipIndexing == "true" {
		fmt.Println("  ⏭️  Skipping indexing (SKIP_INDEXING=true)")
	} else {
		fmt.Println("  ✅ Using existing documents")
	}

	return nil
}

func setupRoutes() {
	// Enable CORS for all origins
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":   "ConduitMCP RAG API Server",
			"version":   "1.0.0",
			"endpoints": "/chat, /documents, /stats, /health",
		})
	})

	// Chat endpoint
	http.HandleFunc("/chat", handleChat)

	// Documents endpoint
	http.HandleFunc("/documents", handleDocuments)

	// Stats endpoint
	http.HandleFunc("/stats", handleStats)

	// Health check endpoint
	http.HandleFunc("/health", handleHealth)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ChatResponse{Error: "Method not allowed"})
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{Error: "Invalid JSON"})
		return
	}

	if req.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{Error: "Message is required"})
		return
	}

	if req.Limit == 0 {
		req.Limit = 5
	}

	startTime := time.Now()

	// Use semantic search to get relevant content
	memory := ragAgent.Memory
	toolRegistry := mcpServer.GetToolRegistry()

	params := map[string]interface{}{
		"query": req.Message,
		"limit": req.Limit,
	}

	result, err := toolRegistry.Call("semantic_search", params, memory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{Error: fmt.Sprintf("Search error: %v", err)})
		return
	}

	var sources []map[string]interface{}
	var response strings.Builder

	if resultMap, ok := result.(map[string]interface{}); ok {
		if results, ok := resultMap["results"].([]map[string]interface{}); ok {
			if len(results) > 0 {
				response.WriteString("Based on the knowledge base:\n\n")

				// Group results by source document
				sourceMap := make(map[string][]string)
				for _, resultItem := range results {
					if content, ok := resultItem["content"].(string); ok {
						if source, ok := resultItem["source"].(string); ok {
							sourceMap[source] = append(sourceMap[source], content)
							sources = append(sources, map[string]interface{}{
								"source":  source,
								"content": content,
								"score":   resultItem["score"],
							})
						}
					}
				}

				// Format response
				for source, contents := range sourceMap {
					response.WriteString(fmt.Sprintf("**%s:**\n", source))
					for _, content := range contents {
						cleanContent := strings.TrimSpace(content)
						response.WriteString(fmt.Sprintf("• %s\n", cleanContent))
					}
					response.WriteString("\n")
				}
			} else {
				response.WriteString("No relevant information found in the knowledge base.")
			}
		}
	}

	if response.Len() == 0 {
		response.WriteString("I couldn't find relevant information to answer your question.")
	}

	duration := time.Since(startTime)

	json.NewEncoder(w).Encode(ChatResponse{
		Response:     response.String(),
		Sources:      sources,
		ResponseTime: duration.Round(time.Millisecond).String(),
	})
}

func handleDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(DocumentResponse{Error: "Method not allowed"})
		return
	}

	var req DocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DocumentResponse{Error: "Invalid JSON"})
		return
	}

	if req.Content == "" || req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DocumentResponse{Error: "Content and title are required"})
		return
	}

	if req.Type == "" {
		req.Type = "text/plain"
	}

	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}

	// Add default metadata
	req.Metadata["indexed_at"] = time.Now().Format(time.RFC3339)
	req.Metadata["api_upload"] = true

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	doc, err := ragEngine.IndexContent(ctx, req.Content, req.Title, req.Type, req.Metadata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DocumentResponse{Error: fmt.Sprintf("Failed to index document: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(DocumentResponse{
		DocumentID: doc.ID,
		Message:    fmt.Sprintf("Document '%s' successfully indexed", req.Title),
	})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(StatsResponse{Error: "Method not allowed"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats, err := ragEngine.GetStats(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(StatsResponse{Error: fmt.Sprintf("Failed to get stats: %v", err)})
		return
	}

	json.NewEncoder(w).Encode(StatsResponse{
		DocumentCount:       getIntValue(stats, "document_count"),
		ChunkCount:          getIntValue(stats, "chunk_count"),
		EmbeddingModel:      getStringValue(stats, "embedding_model"),
		EmbeddingDimensions: getIntValue(stats, "embedding_dimensions"),
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ragEngine.HealthCheck(ctx)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"message": "RAG system is operational",
	})
}

// Helper functions
func getIntValue(m map[string]interface{}, key string) int {
	if val, exists := m[key]; exists {
		if intVal, ok := val.(int); ok {
			return intVal
		}
		if floatVal, ok := val.(float64); ok {
			return int(floatVal)
		}
	}
	return 0
}

func getStringValue(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Missing helper functions implementation

func registerRAGTools(server *conduit.EnhancedServer) {
	// Semantic search tool
	server.RegisterToolWithSchema("semantic_search",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			query, ok := params["query"].(string)
			if !ok {
				return nil, fmt.Errorf("query parameter required")
			}

			limit := 3
			if l, ok := params["limit"].(float64); ok {
				limit = int(l)
			}

			// Get RAG engine from memory
			ragEngineInterface := memory.Get("rag_engine")
			if ragEngineInterface == nil {
				return nil, fmt.Errorf("RAG engine not available")
			}

			ragEngine, ok := ragEngineInterface.(rag.RAGEngine)
			if !ok {
				return nil, fmt.Errorf("invalid RAG engine type")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			results, err := ragEngine.Search(ctx, query, limit, nil)
			if err != nil {
				return nil, fmt.Errorf("search failed: %w", err)
			}

			// Format results
			var formattedResults []map[string]interface{}
			for _, result := range results {
				formattedResults = append(formattedResults, map[string]interface{}{
					"score":   result.Score,
					"content": result.Chunk.Content,
					"source":  result.Document.Title,
				})
			}

			return map[string]interface{}{
				"query":   query,
				"results": formattedResults,
				"count":   len(results),
			}, nil
		}, conduit.ToolMetadata{
			Name:        "semantic_search",
			Description: "Search the knowledge base for relevant information",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query for finding relevant documents",
					},
					"limit": map[string]interface{}{
						"type":        "number",
						"description": "Maximum number of results to return (default: 3)",
					},
				},
				"required": []string{"query"},
			},
		})

	// Knowledge query tool (RAG with AI generation)
	server.RegisterToolWithSchema("knowledge_query",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			question, ok := params["question"].(string)
			if !ok {
				return nil, fmt.Errorf("question parameter required")
			}

			// Get RAG engine from memory
			ragEngineInterface := memory.Get("rag_engine")
			if ragEngineInterface == nil {
				return nil, fmt.Errorf("RAG engine not available")
			}

			ragEngine, ok := ragEngineInterface.(rag.RAGEngine)
			if !ok {
				return nil, fmt.Errorf("invalid RAG engine type")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			response, err := ragEngine.Query(ctx, question, 5, nil)
			if err != nil {
				return nil, fmt.Errorf("RAG query failed: %w", err)
			}

			return map[string]interface{}{
				"question":   question,
				"answer":     response.Answer,
				"confidence": response.Confidence,
				"sources":    len(response.Sources),
			}, nil
		}, conduit.ToolMetadata{
			Name:        "knowledge_query",
			Description: "Ask questions and get AI-generated answers from the knowledge base",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"question": map[string]interface{}{
						"type":        "string",
						"description": "Question to ask about policies, procedures, or information",
					},
				},
				"required": []string{"question"},
			},
		})

	// List documents tool
	server.RegisterToolWithSchema("list_documents",
		func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
			// Get RAG engine from memory
			ragEngineInterface := memory.Get("rag_engine")
			if ragEngineInterface == nil {
				return nil, fmt.Errorf("RAG engine not available")
			}

			ragEngine, ok := ragEngineInterface.(rag.RAGEngine)
			if !ok {
				return nil, fmt.Errorf("invalid RAG engine type")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			stats, err := ragEngine.GetStats(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get stats: %w", err)
			}

			return map[string]interface{}{
				"document_count":       getIntValue(stats, "document_count"),
				"chunk_count":          getIntValue(stats, "chunk_count"),
				"embedding_model":      getStringValue(stats, "embedding_model"),
				"embedding_dimensions": getIntValue(stats, "embedding_dimensions"),
			}, nil
		}, conduit.ToolMetadata{
			Name:        "list_documents",
			Description: "Get information about documents in the knowledge base",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		})
}

func createRAGSystemPrompt() string {
	return `You are an intelligent knowledge assistant. You have access to:

1. **Knowledge Base**: Search using semantic_search and knowledge_query tools
2. **MCP Tools**: Various text processing, memory, and utility tools
3. **Memory System**: Remember context and user preferences

Your capabilities:
- Answer questions about policies, procedures, and guidelines
- Search the knowledge base for relevant information
- Use tools to process text, generate IDs, encode data, etc.
- Remember important information across conversations
- Provide helpful, accurate, and contextual responses

Guidelines:
- Always search the knowledge base first for relevant questions
- Use knowledge_query tool with the exact user question as the "question" parameter
- Use semantic_search tool with relevant keywords as the "query" parameter
- Cite sources when referencing documents
- Use tools when they can enhance your response
- Be conversational but professional
- Ask clarifying questions if needed
- Remember user context and preferences

CRITICAL TOOL PARAMETER REQUIREMENTS - FOLLOW EXACTLY:
1. knowledge_query: MUST use parameter name "question" (NOT "query")
2. semantic_search: MUST use parameter name "query" (NOT "question")
3. remember: MUST use parameters "key" and "value"
4. recall: MUST use parameter "key"

Available tools:
- knowledge_query: Get AI-generated answers - parameter: "question" 
- semantic_search: Search documents - parameter: "query"
- remember: Store info - parameters: "key", "value"
- recall: Get stored info - parameter: "key"
- Text tools: uppercase, lowercase, word_count
- Utility tools: timestamp, uuid, base64_encode/decode

RESPOND WITH VALID JSON ONLY. NO TEXT OUTSIDE JSON.`
}

func indexSampleDocuments(ctx context.Context, ragEngine rag.RAGEngine, provider string) error {
	// Ultra-lightweight documents to prevent CPU/memory issues
	documents := []struct {
		title   string
		content string
	}{
		{"Quick Start", getMinimalQuickStart()},
		{"API Guide", getMinimalAPIGuide()},
	}

	fmt.Printf("    📝 Indexing %d minimal documents...\n", len(documents))

	for i, doc := range documents {
		fmt.Printf("    📄 Processing: %s (%d/%d) - ", doc.title, i+1, len(documents))

		metadata := map[string]interface{}{
			"provider": provider,
			"type":     "docs",
			"minimal":  true,
		}

		// Very short timeout and careful resource management
		indexCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		_, err := ragEngine.IndexContent(indexCtx, doc.content, doc.title, "text/plain", metadata)
		cancel()

		if err != nil {
			fmt.Printf("❌ Skipped: %v\n", err)
			continue // Skip problematic documents
		}
		fmt.Printf("✅\n")

		// Longer pause to prevent system overload
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("    ✅ Minimal indexing completed\n")
	return nil
}

// Ultra-minimal document content to prevent CPU/memory overload
func getMinimalQuickStart() string {
	return `ConduitMCP is a Go RAG library with REST API endpoints for chat and documents.`
}

func getMinimalAPIGuide() string {
	return `API endpoints: POST /chat for questions, POST /documents for adding content, GET /stats for info.`
}
