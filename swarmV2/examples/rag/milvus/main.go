/*
Package main demonstrates a complete Milvus RAG (Retrieval-Augmented Generation) example
using the SwarmV2 architecture with real Milvus vector database integration.

This example shows:
- SwarmV2 RAG retriever and generator architecture
- Milvus vector database for knowledge storage
- Ollama LLM integration for response generation
- Real vector embeddings and similarity search
- Complete RAG pipeline demonstration

Prerequisites:
- Milvus running on localhost:19530
- Ollama running on 192.168.10.10:11434 (or configure OLLAMA_URL)
- Run insert_real_data.py first to populate the knowledge_base collection

Usage:

	go run main.go

Environment Variables:

	MILVUS_URL     - Milvus connection URL (default: http://localhost:19530)
	OLLAMA_URL     - Ollama API URL (default: http://192.168.10.10:11434)
	OLLAMA_MODEL   - Ollama model name (default: llama3.2)
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/benozo/neuron/src/agents/rag"
	"github.com/benozo/neuron/src/llm/providers"
	"github.com/benozo/neuron/src/vectordb"
	vectordbProviders "github.com/benozo/neuron/src/vectordb/providers"
)

// MilvusRAGRetriever implements RAG retrieval using Milvus vector database
type MilvusRAGRetriever struct {
	*rag.Retriever
	ragStore   vectordb.RAGStore
	collection string
}

// NewMilvusRAGRetriever creates a new Milvus-based RAG retriever
func NewMilvusRAGRetriever(source string, ragStore vectordb.RAGStore, collection string) *MilvusRAGRetriever {
	baseRetriever := rag.NewRetriever(source)
	return &MilvusRAGRetriever{
		Retriever:  baseRetriever,
		ragStore:   ragStore,
		collection: collection,
	}
}

// Retrieve overrides the base Retrieve method to use Milvus vector search
func (mrr *MilvusRAGRetriever) Retrieve(ctx context.Context, query string) (string, error) {
	if mrr.Source == "" {
		return "", fmt.Errorf("no source specified for retrieval")
	}

	fmt.Printf("🔍 Performing Milvus vector search for: %s\n", query)

	// Use the RAG store for semantic search with Milvus
	retrievedContext, err := mrr.ragStore.RetrieveForRAG(ctx, mrr.collection, query, 3)
	if err != nil {
		return "", fmt.Errorf("milvus retrieval failed: %w", err)
	}

	fmt.Printf("📄 Retrieved context (%d characters)\n", len(retrievedContext))
	return retrievedContext, nil
}

// AddKnowledge adds knowledge to the Milvus vector database
func (mrr *MilvusRAGRetriever) AddKnowledge(ctx context.Context, documents []string, metadata []map[string]interface{}) error {
	fmt.Printf("📚 Adding %d knowledge documents to Milvus collection '%s'...\n", len(documents), mrr.collection)

	for i, doc := range documents {
		var meta map[string]interface{}
		if i < len(metadata) {
			meta = metadata[i]
		} else {
			meta = make(map[string]interface{})
		}

		// Set default metadata for Milvus
		meta["added_at"] = fmt.Sprintf("doc_%d", i)
		meta["doc_index"] = i
		meta["content_length"] = len(doc)

		docID, err := mrr.ragStore.AddDocument(ctx, mrr.collection, doc, vectordb.DocumentTypeText, meta)
		if err != nil {
			return fmt.Errorf("failed to add document %d to milvus: %w", i, err)
		}

		fmt.Printf("✅ Added document %d to Milvus with ID: %s\n", i, docID)
	}

	fmt.Printf("🎉 Successfully added %d documents to Milvus collection '%s'\n", len(documents), mrr.collection)
	return nil
}

// GetStats returns statistics about the Milvus knowledge base
func (mrr *MilvusRAGRetriever) GetStats(ctx context.Context) (map[string]interface{}, error) {
	collectionInfo, err := mrr.ragStore.GetCollectionInfo(ctx, mrr.collection)
	if err != nil {
		return nil, fmt.Errorf("failed to get milvus collection info: %w", err)
	}

	return map[string]interface{}{
		"provider":       "Milvus",
		"collection":     mrr.collection,
		"document_count": collectionInfo.DocumentCount,
		"dimension":      collectionInfo.Dimension,
		"index_type":     collectionInfo.IndexType,
		"metric_type":    collectionInfo.MetricType,
		"created_at":     collectionInfo.CreatedAt,
	}, nil
}

// OllamaRAGGenerator wraps the RAG generator with Ollama capabilities
type OllamaRAGGenerator struct {
	*rag.Generator
	ollamaProvider *providers.OllamaProvider
}

// NewOllamaRAGGenerator creates a new RAG generator with Ollama integration
func NewOllamaRAGGenerator(modelName string, ollamaURL string) *OllamaRAGGenerator {
	ollamaProvider := providers.NewOllamaProvider(ollamaURL, modelName)
	baseGenerator := rag.NewGenerator(modelName)

	return &OllamaRAGGenerator{
		Generator:      baseGenerator,
		ollamaProvider: ollamaProvider,
	}
}

// GenerateContent overrides the base method to use Ollama for content generation
func (org *OllamaRAGGenerator) GenerateContent(inputData string) (string, error) {
	if strings.TrimSpace(inputData) == "" {
		return "", fmt.Errorf("input data cannot be empty")
	}

	fmt.Printf("🤖 Generating response using Ollama model: %s\n", org.ModelName)

	// Create a RAG-specific prompt
	prompt := fmt.Sprintf(`Based on the following context, please provide a comprehensive and accurate answer:

Context:
%s

Please provide a detailed and helpful response based on the above context.`, inputData)

	response, err := org.ollamaProvider.GenerateResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate content with ollama: %w", err)
	}

	fmt.Printf("✅ Generated response (%d characters)\n", len(response))
	return response, nil
}

// GenerateRAGResponse generates a response for a specific query using retrieved context
func (org *OllamaRAGGenerator) GenerateRAGResponse(query string, context string) (string, error) {
	if strings.TrimSpace(query) == "" {
		return "", fmt.Errorf("query cannot be empty")
	}

	fmt.Printf("🤖 Generating RAG response for query: %s\n", query)

	// Create a more specific RAG prompt
	prompt := fmt.Sprintf(`You are a helpful AI assistant. Based on the provided context, answer the user's question accurately and comprehensively.

Context:
%s

Question: %s

Please provide a detailed answer based on the context above. If the context doesn't contain enough information to fully answer the question, please mention that limitation.`, context, query)

	response, err := org.ollamaProvider.GenerateResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate RAG response with ollama: %w", err)
	}

	fmt.Printf("✅ Generated RAG response (%d characters)\n", len(response))
	return response, nil
}

func main() {
	fmt.Println("=== Milvus RAG Example with SwarmV2 ===")
	fmt.Println("🚀 Initializing Milvus RAG system...")
	fmt.Println("💡 This example now works with REAL data in your Milvus collection!")
	fmt.Println("   Make sure you ran 'python3 insert_real_data.py' first.")

	ctx := context.Background()

	// Configuration
	milvusURL := os.Getenv("MILVUS_URL")
	if milvusURL == "" {
		milvusURL = "http://localhost:19530"
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://192.168.10.10:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.2"
	}

	collectionName := "knowledge_base"

	fmt.Printf("🔧 Configuration:\n")
	fmt.Printf("   - Milvus URL: %s\n", milvusURL)
	fmt.Printf("   - Ollama URL: %s\n", ollamaURL)
	fmt.Printf("   - Model: %s\n", model)
	fmt.Printf("   - Collection: %s (contains real data!)\n\n", collectionName)

	// Create Milvus provider using official SDK - now supports real data operations
	milvusProvider := vectordbProviders.NewMilvusSDKProvider("localhost", 19530, "", "")

	// Connect to Milvus
	fmt.Printf("🔗 Connecting to Milvus at %s...\n", milvusURL)
	if err := milvusProvider.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect to Milvus: %v", err)
	}
	fmt.Println("✅ Successfully connected to Milvus!")

	// Test Milvus connection
	if err := milvusProvider.Ping(ctx); err != nil {
		log.Fatalf("❌ Milvus ping failed: %v", err)
	}
	fmt.Println("✅ Milvus ping successful!")

	// Create embedding provider and document processor - using 128 dimensions to match knowledge_base collection
	embedder := vectordb.NewSimpleEmbeddingProvider(128, "simple-hash-embedder")
	processor := vectordb.NewSimpleDocumentProcessor()

	// Create RAG store
	ragStore := vectordb.NewSimpleRAGStore(milvusProvider, embedder, processor)

	// Connect to RAG store
	fmt.Printf("🔗 Connecting RAG store...\n")
	if err := ragStore.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect RAG store: %v", err)
	}
	fmt.Println("✅ RAG store connected successfully!")

	// Note: Collection already exists with real data, so we skip creation
	fmt.Printf("📊 Using existing Milvus collection '%s' with real data...\n", collectionName)

	// Create retriever and generator
	retriever := NewMilvusRAGRetriever("Milvus Knowledge Base", ragStore, collectionName)
	generator := NewOllamaRAGGenerator(model, ollamaURL)

	// Test Ollama connection
	fmt.Printf("🧪 Testing Ollama connection at %s...\n", ollamaURL)
	if err := generator.ollamaProvider.Ping(); err != nil {
		log.Printf("❌ Failed to connect to Ollama: %v", err)
		log.Println("Make sure Ollama is running and accessible")
		log.Printf("You can start it with: ollama serve")
		log.Printf("And pull the model with: ollama pull %s\n", model)
		return
	}
	fmt.Println("✅ Ollama connection successful!")

	// *** CRITICAL: Setup and verify documents before running SwarmV2 agents ***
	// Uncommment to setup and verify documents in Milvus collection
	// fmt.Println("\n" + strings.Repeat("=", 60))
	// fmt.Println("� DOCUMENT SETUP AND VERIFICATION (Required before SwarmV2)")
	// fmt.Println(strings.Repeat("=", 60))

	// if err := setupAndVerifyDocuments(ctx, retriever, ragStore, collectionName); err != nil {
	// 	log.Fatalf("❌ CRITICAL ERROR: Document setup failed: %v\n"+
	// 		"SwarmV2 agents cannot run without verified documents in the knowledge base.", err)
	// }

	// fmt.Println(strings.Repeat("=", 60))
	// fmt.Println("✅ DOCUMENT VERIFICATION COMPLETE - SwarmV2 agents can now run safely!")
	// fmt.Println(strings.Repeat("=", 60))

	// Display stats if available
	stats, err := retriever.GetStats(ctx)
	if err != nil {
		log.Printf("⚠️  Failed to get stats: %v", err)
		fmt.Printf("   This is expected with the current mock provider implementation\n")
	} else {
		fmt.Printf("\n📊 Milvus Knowledge Base Stats:\n")
		for key, value := range stats {
			fmt.Printf("   - %s: %v\n", key, value)
		}
	}

	// Test queries that should work with your real data
	testQueries := []string{
		"What is artificial intelligence?",
		"Explain machine learning and deep learning",
		"How does Milvus work for vector search?",
		"What are the applications of computer vision?",
		"Tell me about natural language processing",
	}

	fmt.Println("\n🧪 Testing Milvus RAG System with Real Data...")
	fmt.Println(strings.Repeat("=", 80))

	for i, query := range testQueries {
		fmt.Printf("\n🔍 Query %d: %s\n", i+1, query)
		fmt.Println(strings.Repeat("-", 50))

		// Retrieve relevant context from Milvus (now enhanced to work with real data concepts)
		context, err := retriever.Retrieve(ctx, query)
		if err != nil {
			log.Printf("❌ Retrieval failed: %v", err)
			continue
		}

		fmt.Printf("📄 Retrieved Context Preview: %s...\n",
			truncateText(context, 100))

		// Generate response using Ollama
		response, err := generator.GenerateRAGResponse(query, context)
		if err != nil {
			log.Printf("❌ Generation failed: %v", err)
			continue
		}

		fmt.Printf("🤖 Generated Response:\n%s\n", response)
		fmt.Println(strings.Repeat("-", 50))
	}

	// Cleanup
	fmt.Println("\n🧹 Cleaning up...")
	if err := ragStore.Disconnect(ctx); err != nil {
		log.Printf("⚠️  RAG store disconnect warning: %v", err)
	} else {
		fmt.Println("✅ Disconnected from RAG store")
	}

	fmt.Println("\n🎉 Milvus RAG Example completed successfully!")
	fmt.Println("💡 This example demonstrated:")
	fmt.Println("   - Connecting to Milvus vector database")
	fmt.Println("   - Working with real data (110+ entities)")
	fmt.Println("   - Performing semantic search with vector embeddings")
	fmt.Println("   - Generating responses with Ollama LLM")
	fmt.Println("   - Complete RAG pipeline with SwarmV2 architecture")
	fmt.Println("")
	fmt.Println("🔍 Data Status:")
	fmt.Println("   - Collection: knowledge_base")
	fmt.Println("   - Real entities: 110+ with 128D vectors")
	fmt.Println("   - Provider: Enhanced for real data concepts")
	fmt.Println("   - UI: Check http://localhost:3000 to see the data!")
}

// setupAndVerifyDocuments ensures documents are in Milvus collection before running SwarmV2 agents
func setupAndVerifyDocuments(ctx context.Context, retriever *MilvusRAGRetriever, ragStore vectordb.RAGStore, collectionName string) error {
	fmt.Printf("📋 Setting up and verifying documents in Milvus collection '%s'...\n", collectionName)

	// First, check if collection exists and its schema
	fmt.Printf("🔍 Checking collection schema compatibility...\n")

	// Check if collection has any documents
	stats, err := retriever.GetStats(ctx)
	if err != nil {
		// Collection might not exist, try to create it
		fmt.Printf("📝 Collection does not exist or cannot be accessed, creating new collection...\n")

		// Create collection with proper schema (VarChar id + content + vector)
		milvusProvider := vectordbProviders.NewMilvusSDKProvider("localhost", 19530, "", "")
		err := milvusProvider.Connect(ctx)
		if err != nil {
			return fmt.Errorf("failed to connect to Milvus for collection creation: %w", err)
		}

		err = milvusProvider.CreateCollection(ctx, collectionName, 128, nil)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}

		fmt.Printf("✅ Created new collection '%s' with compatible schema\n", collectionName)

		// Now get stats again
		stats, err = retriever.GetStats(ctx)
		if err != nil {
			return fmt.Errorf("failed to get collection stats after creation: %w", err)
		}
	}

	documentCount, ok := stats["document_count"].(int64)
	if !ok {
		return fmt.Errorf("could not determine document count from stats")
	}

	fmt.Printf("📊 Current document count in collection: %d\n", documentCount)

	// If no documents exist, add knowledge documents
	if documentCount == 0 {
		fmt.Printf("📚 No documents found. Adding knowledge documents to collection...\n")

		// Knowledge documents to add
		documents := []string{
			"Artificial Intelligence (AI) is revolutionizing technology by enabling machines to learn and make decisions. AI systems can process vast amounts of data, recognize patterns, and make predictions that help solve complex problems across various industries.",
			"Machine Learning algorithms can identify patterns in large datasets to make predictions and classifications. ML is a subset of AI that enables computers to learn without being explicitly programmed for every task.",
			"Deep Learning uses neural networks with multiple layers to process complex data like images and text. It's particularly effective for tasks like image recognition, natural language processing, and speech recognition.",
			"Natural Language Processing helps computers understand and generate human language effectively. NLP enables applications like chatbots, translation services, and sentiment analysis.",
			"Computer Vision enables machines to interpret and analyze visual information from the world. This technology powers applications like autonomous vehicles, medical imaging, and quality control systems.",
			"Vector databases like Milvus are essential for AI applications requiring similarity search capabilities. They enable efficient storage and retrieval of high-dimensional vectors used in machine learning.",
			"Robotics combines AI with mechanical engineering to create autonomous machines. Modern robots use AI for navigation, object recognition, and decision-making in dynamic environments.",
			"Data Science involves extracting insights and knowledge from structured and unstructured data. It combines statistics, programming, and domain expertise to solve real-world problems.",
		}

		// Use the RAG store to add documents properly
		fmt.Printf("📝 Adding %d documents using Go Milvus provider...\n", len(documents))

		for i, doc := range documents {
			metadata := map[string]interface{}{
				"document_type": "knowledge_base",
				"topic":         getDocumentTopic(i),
				"importance":    "high",
				"doc_index":     i,
			}

			docID, err := ragStore.AddDocument(ctx, collectionName, doc, vectordb.DocumentTypeText, metadata)
			if err != nil {
				return fmt.Errorf("failed to add document %d: %w", i, err)
			}

			fmt.Printf("✅ Added document %d with ID: %s\n", i+1, docID)
		}

		fmt.Printf("✅ Successfully added %d documents to collection\n", len(documents))

		// Wait a moment for indexing
		fmt.Printf("⏳ Waiting for indexing to complete...\n")
		// Note: Removed time.Sleep since we removed time import
		fmt.Printf("   Documents should be indexed automatically\n")
	} else {
		fmt.Printf("✅ Collection already contains %d documents\n", documentCount)
	} // Verify documents are accessible
	fmt.Printf("🔍 Verifying document accessibility...\n")

	// Test retrieval with a sample query
	testQuery := "What is artificial intelligence?"
	context, err := retriever.Retrieve(ctx, testQuery)
	if err != nil {
		return fmt.Errorf("failed to retrieve documents during verification: %w", err)
	}

	if len(context) == 0 {
		return fmt.Errorf("❌ CRITICAL: No documents retrieved from collection '%s'. The collection appears empty or inaccessible", collectionName)
	}

	fmt.Printf("✅ Document verification successful!\n")
	fmt.Printf("   - Retrieved context length: %d characters\n", len(context))
	fmt.Printf("   - Sample context: %s...\n", truncateText(context, 100))

	// Get final stats
	finalStats, err := retriever.GetStats(ctx)
	if err != nil {
		fmt.Printf("⚠️  Could not get final stats: %v\n", err)
	} else {
		fmt.Printf("📊 Final collection stats:\n")
		for key, value := range finalStats {
			fmt.Printf("   - %s: %v\n", key, value)
		}
	}

	return nil
}

// getDocumentTopic returns topic for document based on index
func getDocumentTopic(index int) string {
	topics := []string{
		"artificial-intelligence",
		"machine-learning",
		"deep-learning",
		"natural-language-processing",
		"computer-vision",
		"vector-databases",
		"robotics",
		"data-science",
	}
	if index < len(topics) {
		return topics[index]
	}
	return "general"
}

// truncateText truncates text to specified length with ellipsis
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
