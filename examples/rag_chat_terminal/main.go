package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
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

func main() {
	fmt.Println("🤖 ConduitMCP RAG + LLM Interactive Chat Terminal")
	fmt.Println("================================================")
	fmt.Println("💬 Chat with TechCorp's Knowledge Base using Ollama + RAG")
	fmt.Println("🔧 Includes MCP tools for enhanced functionality")
	fmt.Println("")

	// Configuration
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "192.168.10.10:11434"
	}

	llmModel := os.Getenv("OLLAMA_LLM_MODEL")
	if llmModel == "" {
		llmModel = "llama3.2"
	}

	embeddingModel := os.Getenv("OLLAMA_EMBEDDING_MODEL")
	if embeddingModel == "" {
		embeddingModel = "nomic-embed-text:latest"
	}

	fmt.Printf("🦙 Ollama Host: %s\n", ollamaHost)
	fmt.Printf("🧠 LLM Model: %s\n", llmModel)
	fmt.Printf("🔢 Embedding Model: %s\n", embeddingModel)
	fmt.Println("")

	ctx := context.Background()

	// Initialize MCP server in library mode
	fmt.Println("🔧 Initializing MCP Tools in Library Mode...")
	config := conduit.DefaultConfig()
	config.EnableLogging = false // Quiet mode for chat
	mcpServer := conduit.NewEnhancedServer(config)

	// Register all MCP tools
	tools.RegisterTextTools(mcpServer)
	tools.RegisterMemoryTools(mcpServer)
	tools.RegisterUtilityTools(mcpServer)

	// Register RAG-specific tools
	registerRAGTools(mcpServer)

	fmt.Printf("✅ Registered MCP tools\n")

	// Initialize RAG system
	fmt.Println("📚 Initializing RAG Knowledge Base...")
	ragConfig := rag.DefaultOllamaRAGConfig()
	// Extract just the host part without port since NewOllamaEmbeddings adds :11434
	hostParts := strings.Split(ollamaHost, ":")
	ragConfig.Embeddings.Host = hostParts[0]
	ragConfig.Embeddings.Model = embeddingModel

	vectorDB, err := database.NewPgVectorDB(ragConfig.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer vectorDB.Close()

	embeddingProvider := embeddings.NewOllamaEmbeddings(
		ragConfig.Embeddings.Host,
		ragConfig.Embeddings.Model,
		ragConfig.Embeddings.Dimensions,
		ragConfig.Embeddings.Timeout,
	)

	// Test embedding connection
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	if err := embeddingProvider.Ping(pingCtx); err != nil {
		cancel()
		log.Fatalf("Embedding provider connection failed: %v", err)
	}
	cancel()

	chunker := processors.NewTextChunker(
		processors.Paragraph,
		ragConfig.Chunking.Size,
		ragConfig.Chunking.Overlap,
	)

	ragEngine := rag.NewRAGEngine(ragConfig, vectorDB, embeddingProvider, chunker)

	// Check knowledge base
	stats, err := ragEngine.GetStats(ctx)
	if err != nil {
		log.Printf("Warning: Could not get stats: %v", err)
		stats = make(map[string]interface{})
	}

	existingDocs := getIntValue(stats, "document_count")
	fmt.Printf("📊 Knowledge Base: %d documents, %d chunks\n",
		existingDocs, getIntValue(stats, "chunk_count"))

	// Index documents if needed
	if existingDocs == 0 {
		fmt.Println("📝 Indexing TechCorp knowledge base...")
		if err := indexTechCorpDocuments(ctx, ragEngine); err != nil {
			log.Printf("Warning: Failed to index documents: %v", err)
		}
	}

	// Initialize LLM agent with RAG capabilities
	fmt.Println("🧠 Creating RAG-enabled LLM Agent...")
	var ollamaURL string
	if !strings.HasPrefix(ollamaHost, "http") {
		ollamaURL = "http://" + ollamaHost
	} else {
		ollamaURL = ollamaHost
	}
	ollamaModel := conduit.CreateOllamaModel(ollamaURL)

	// Create specialized agent manager
	agentManager := agents.NewLLMAgentManager(mcpServer, ollamaModel, llmModel)

	// Create RAG-enabled assistant agent
	ragAgent, err := agentManager.CreateLLMAgent(
		"rag_assistant",
		"TechCorp Knowledge Assistant",
		"An intelligent assistant that can search company knowledge and use various tools",
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

	if err != nil {
		log.Fatalf("Failed to create RAG agent: %v", err)
	}

	// Set up RAG context in memory
	ragAgent.Memory.Set("rag_engine", ragEngine)
	ragAgent.Memory.Set("context", "TechCorp employee using knowledge base")

	fmt.Println("✅ RAG Chat System Ready!")
	fmt.Println("")

	// Display help
	displayHelp()

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start interactive chat loop
	scanner := bufio.NewScanner(os.Stdin)
	sessionID := fmt.Sprintf("chat_%d", time.Now().Unix())

	for {
		fmt.Print("\n💬 You: ")

		// Check for shutdown signal
		select {
		case <-sigChan:
			fmt.Println("\n👋 Goodbye! Chat session ended.")
			return
		default:
		}

		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}

		// Handle special commands
		if handleSpecialCommands(userInput, ragEngine, agentManager) {
			continue
		}

		// Process user query with RAG agent
		fmt.Print("🤖 Assistant: ")

		startTime := time.Now()

		// Create a task for the user query
		task, err := agentManager.CreateTask(
			"rag_assistant",
			"User Query",
			"Process user query with RAG and tool capabilities",
			map[string]interface{}{
				"user_query": userInput,
				"session_id": sessionID,
				"timestamp":  time.Now().Format(time.RFC3339),
			},
		)

		if err != nil {
			fmt.Printf("❌ Error creating task: %v\n", err)
			continue
		}

		// Execute with LLM reasoning
		err = agentManager.ExecuteTaskWithLLM(task.ID)

		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			// Use semantic search to provide comprehensive answers
			fmt.Printf("� Searching knowledge base...\n\n")

			memory := ragAgent.Memory
			toolRegistry := mcpServer.GetToolRegistry()

			// Use semantic search to get relevant content
			params := map[string]interface{}{
				"query": userInput,
				"limit": 5,
			}

			result, err := toolRegistry.Call("semantic_search", params, memory)
			if err != nil {
				fmt.Printf("❌ Search error: %v\n", err)
			} else if resultMap, ok := result.(map[string]interface{}); ok {
				if results, ok := resultMap["results"].([]map[string]interface{}); ok {
					if len(results) > 0 {
						fmt.Printf("📚 **Found relevant information:**\n\n") // Group results by source document
						sourceMap := make(map[string][]string)
						for _, resultItem := range results {
							if content, ok := resultItem["content"].(string); ok {
								if source, ok := resultItem["source"].(string); ok {
									sourceMap[source] = append(sourceMap[source], content)
								}
							}
						}

						// Display organized results
						for source, contents := range sourceMap {
							fmt.Printf("📄 **%s:**\n", source)
							for i, content := range contents {
								// Clean up and format content
								cleanContent := strings.TrimSpace(content)
								if strings.HasPrefix(cleanContent, "#") {
									// This is a header, format it nicely
									fmt.Printf("   %s\n", cleanContent)
								} else {
									// This is content, indent it
									fmt.Printf("   • %s\n", cleanContent)
								}

								if i < len(contents)-1 {
									fmt.Printf("\n")
								}
							}
							fmt.Printf("\n")
						}
					} else {
						fmt.Printf("❌ No relevant information found in the knowledge base.\n")
					}
				} else {
					fmt.Printf("❌ No search results found.\n")
				}
			}

			fmt.Printf("\n⏱️  Response time: %v\n", duration.Round(time.Millisecond))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

func createRAGSystemPrompt() string {
	return `You are TechCorp's intelligent knowledge assistant. You have access to:

1. **Company Knowledge Base**: Search using semantic_search and knowledge_query tools
2. **MCP Tools**: Various text processing, memory, and utility tools
3. **Memory System**: Remember context and user preferences

Your capabilities:
- Answer questions about TechCorp policies, procedures, and guidelines
- Search the knowledge base for relevant information
- Use tools to process text, generate IDs, encode data, etc.
- Remember important information across conversations
- Provide helpful, accurate, and contextual responses

Guidelines:
- Always search the knowledge base first for company-related questions
- Use knowledge_query tool with the exact user question as the "question" parameter
- Use semantic_search tool with relevant keywords as the "query" parameter
- Cite sources when referencing company documents
- Use tools when they can enhance your response
- Be conversational but professional
- Ask clarifying questions if needed
- Remember user context and preferences

CRITICAL TOOL PARAMETER REQUIREMENTS - FOLLOW EXACTLY:
1. knowledge_query: MUST use parameter name "question" (NOT "query")
2. semantic_search: MUST use parameter name "query" (NOT "question")
3. remember: MUST use parameters "key" and "value"
4. recall: MUST use parameter "key"

TOOL USAGE EXAMPLES - COPY EXACTLY:

For answering user questions, use knowledge_query:
{
  "name": "Get Answer",
  "description": "Get answer to user question",
  "tool": "knowledge_query",
  "input": {
    "question": "What is our remote work policy?"
  }
}

For searching documents, use semantic_search:
{
  "name": "Search Documents",
  "description": "Search for relevant documents",
  "tool": "semantic_search",
  "input": {
    "query": "remote work policy"
  }
}

Available tools:
- knowledge_query: Get AI-generated answers - parameter: "question" 
- semantic_search: Search documents - parameter: "query"
- remember: Store info - parameters: "key", "value"
- recall: Get stored info - parameter: "key"
- Text tools: uppercase, lowercase, word_count
- Utility tools: timestamp, uuid, base64_encode/decode

RESPOND WITH VALID JSON ONLY. NO TEXT OUTSIDE JSON.`
}

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
			Description: "Search the TechCorp knowledge base for relevant information",
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
			Description: "Ask questions and get AI-generated answers from the TechCorp knowledge base",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"question": map[string]interface{}{
						"type":        "string",
						"description": "Question to ask about TechCorp policies, procedures, or information",
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
			Description: "Get information about documents in the TechCorp knowledge base",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		})
}

func handleSpecialCommands(input string, ragEngine rag.RAGEngine, agentManager *agents.LLMAgentManager) bool {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "/help", "/h":
		displayHelp()
		return true
	case "/quit", "/exit", "/q":
		fmt.Println("👋 Goodbye! Chat session ended.")
		os.Exit(0)
		return true
	case "/stats":
		displayStats(ragEngine)
		return true
	case "/clear":
		fmt.Print("\033[2J\033[H") // Clear screen
		fmt.Println("💬 Chat history cleared.")
		return true
	case "/tasks":
		displayTasks(agentManager)
		return true
	default:
		if strings.HasPrefix(strings.ToLower(input), "/search ") {
			query := strings.TrimSpace(input[8:])
			performDirectSearch(ragEngine, query)
			return true
		}
	}
	return false
}

func displayHelp() {
	fmt.Println("🆘 Available Commands:")
	fmt.Println("   /help or /h     - Show this help")
	fmt.Println("   /quit or /q     - Exit chat")
	fmt.Println("   /stats          - Show knowledge base statistics")
	fmt.Println("   /search <query> - Direct search (no AI generation)")
	fmt.Println("   /clear          - Clear screen")
	fmt.Println("   /tasks          - Show recent tasks")
	fmt.Println("")
	fmt.Println("💡 Examples:")
	fmt.Println("   'What is our remote work policy?'")
	fmt.Println("   'How do I submit expenses?'")
	fmt.Println("   'What are the onboarding steps for developers?'")
	fmt.Println("   'Remember that I work in the engineering team'")
	fmt.Println("   'Generate a UUID for my project'")
}

func displayStats(ragEngine rag.RAGEngine) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats, err := ragEngine.GetStats(ctx)
	if err != nil {
		fmt.Printf("❌ Error getting stats: %v\n", err)
		return
	}

	fmt.Println("📊 Knowledge Base Statistics:")
	fmt.Printf("   📚 Documents: %d\n", getIntValue(stats, "document_count"))
	fmt.Printf("   📄 Chunks: %d\n", getIntValue(stats, "chunk_count"))
	fmt.Printf("   🧠 Model: %s\n", getStringValue(stats, "embedding_model"))
	fmt.Printf("   📏 Dimensions: %d\n", getIntValue(stats, "embedding_dimensions"))
}

func performDirectSearch(ragEngine rag.RAGEngine, query string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("🔍 Searching for: %s\n", query)
	results, err := ragEngine.Search(ctx, query, 3, nil)
	if err != nil {
		fmt.Printf("❌ Search failed: %v\n", err)
		return
	}

	fmt.Printf("📊 Found %d results:\n", len(results))
	for i, result := range results {
		fmt.Printf("   %d. [%.3f] %s\n", i+1, result.Score,
			truncateString(result.Chunk.Content, 80))
		fmt.Printf("      Source: %s\n", result.Document.Title)
	}
}

func displayTasks(agentManager *agents.LLMAgentManager) {
	// Since GetAllTasks is not available, just show a placeholder
	fmt.Println("📋 Task tracking not implemented in this example")
}

func indexTechCorpDocuments(ctx context.Context, ragEngine rag.RAGEngine) error {
	documents := []struct {
		title   string
		content string
	}{
		{"Employee Handbook 2024", getEmployeeHandbook()},
		{"Remote Work Policy", getRemoteWorkPolicy()},
		{"Data Security Guidelines", getDataSecurityGuidelines()},
		{"Customer Onboarding Process", getCustomerOnboarding()},
		{"Expense Reimbursement Policy", getExpensePolicy()},
	}

	for i, doc := range documents {
		fmt.Printf("   📄 Indexing: %s (%d/%d)\n", doc.title, i+1, len(documents))

		metadata := map[string]interface{}{
			"department": "TechCorp",
			"year":       2024,
			"indexed_at": time.Now().Format(time.RFC3339),
		}

		indexCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		_, err := ragEngine.IndexContent(indexCtx, doc.content, doc.title, "text/plain", metadata)
		cancel()

		if err != nil {
			return fmt.Errorf("failed to index %s: %w", doc.title, err)
		}
	}

	return nil
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

// Sample document content (simplified versions)
func getEmployeeHandbook() string {
	return `# TechCorp Employee Handbook 2024

## Onboarding Process for New Employees

### Week 1: Getting Started
1. **Day 1**: Complete HR paperwork and receive equipment
2. **Day 2-3**: IT setup including laptop, accounts, and security training
3. **Day 4-5**: Department introduction and role-specific training

### Software Engineers Onboarding
- Complete security awareness training
- Set up development environment
- Review coding standards and practices
- Assign mentor for first 30 days

## Work Schedule and Time Off

### Standard Work Hours
- Core hours: 9:00 AM - 5:00 PM
- Flexible start time: 7:00 AM - 10:00 AM

### Vacation Policy
- 25 days annual leave for full-time employees
- Must be approved by direct manager
- Minimum 2 weeks notice required for vacations longer than 5 days`
}

func getRemoteWorkPolicy() string {
	return `# TechCorp Remote Work Policy

## Eligibility
- Employees who have completed 6 months of employment
- Roles that don't require physical presence
- Approval from direct manager required

## Remote Work Options
1. **Hybrid Remote**: 2-3 days per week from home
2. **Fully Remote**: Permanent remote work arrangement
3. **Temporary Remote**: Short-term arrangements

## Application Process
1. Complete Remote Work Request Form
2. Discuss with direct manager
3. HR review and approval
4. IT equipment assessment
5. Trial period (30 days for new arrangements)`
}

func getDataSecurityGuidelines() string {
	return `# TechCorp Data Security Guidelines

## Customer Data Handling

### Collection and Processing
- Collect only necessary data for business purposes
- Obtain explicit consent for data processing
- Document legal basis for data processing

### Access Controls
- Role-based access to customer data
- Minimum necessary access principle
- Multi-factor authentication required

## Security Practices for Developers
- Regular security code reviews
- Use of approved libraries and frameworks
- Static code analysis tools`
}

func getCustomerOnboarding() string {
	return `# TechCorp Customer Onboarding Process

## Onboarding Phases

### Phase 1: Project Kickoff (Week 1)
- Welcome call with customer leadership
- Project team introductions
- Scope and timeline confirmation

### Phase 2: Technical Setup (Weeks 2-4)
- Environment provisioning
- Integration planning and setup
- Data migration planning

### Phase 3: Configuration & Training (Weeks 5-8)
- System configuration based on requirements
- User access and permissions setup
- Admin training sessions`
}

func getExpensePolicy() string {
	return `# TechCorp Expense Reimbursement Policy

## Approved Expense Categories

### Travel Expenses
- Flights: Economy class for domestic
- Hotels: Reasonable business hotels, up to $200/night
- Ground Transportation: Taxis, rideshares, public transit

### Approval Requirements
- $0-$500: Direct manager approval
- $501-$2,000: Department head approval
- $2,001+: VP and Finance approval

## Reimbursement Process
1. Submit expenses within 30 days
2. Use company expense management system
3. Include all required receipts and documentation`
}
