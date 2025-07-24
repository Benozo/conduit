package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/benozo/neuron/src/agents/rag"
	"github.com/benozo/neuron/src/llm/providers"
	"github.com/benozo/neuron/src/vectordb"
	vectordbProviders "github.com/benozo/neuron/src/vectordb/providers"
)

// VectorRAGRetriever implements RAG retrieval using vector databases
type VectorRAGRetriever struct {
	*rag.Retriever
	ragStore   vectordb.RAGStore
	collection string
}

// NewVectorRAGRetriever creates a new vector-based RAG retriever
func NewVectorRAGRetriever(source string, ragStore vectordb.RAGStore, collection string) *VectorRAGRetriever {
	baseRetriever := rag.NewRetriever(source)
	return &VectorRAGRetriever{
		Retriever:  baseRetriever,
		ragStore:   ragStore,
		collection: collection,
	}
}

// Retrieve overrides the base Retrieve method to use vector search
func (vrr *VectorRAGRetriever) Retrieve(ctx context.Context, query string) (string, error) {
	if vrr.Source == "" {
		return "", fmt.Errorf("no source specified for retrieval")
	}

	fmt.Printf("🔍 Performing vector search for: %s\n", query)

	// Use the RAG store for semantic search
	retrievedContext, err := vrr.ragStore.RetrieveForRAG(ctx, vrr.collection, query, 3)
	if err != nil {
		return "", fmt.Errorf("vector retrieval failed: %w", err)
	}

	return retrievedContext, nil
}

// AddKnowledge adds knowledge to the vector database
func (vrr *VectorRAGRetriever) AddKnowledge(ctx context.Context, documents []string, metadata []map[string]interface{}) error {
	fmt.Printf("📚 Adding %d knowledge documents to vector database...\n", len(documents))

	for i, doc := range documents {
		var meta map[string]interface{}
		if i < len(metadata) {
			meta = metadata[i]
		} else {
			meta = make(map[string]interface{})
		}

		// Set default metadata
		meta["added_at"] = fmt.Sprintf("%d", len(documents))
		meta["doc_index"] = i

		_, err := vrr.ragStore.AddDocument(ctx, vrr.collection, doc, vectordb.DocumentTypeText, meta)
		if err != nil {
			return fmt.Errorf("failed to add document %d: %w", i, err)
		}
	}

	fmt.Printf("✅ Successfully added %d documents to collection '%s'\n", len(documents), vrr.collection)
	return nil
}

// GetStats returns statistics about the knowledge base
func (vrr *VectorRAGRetriever) GetStats(ctx context.Context) (map[string]interface{}, error) {
	collectionInfo, err := vrr.ragStore.GetCollectionInfo(ctx, vrr.collection)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection info: %w", err)
	}

	return map[string]interface{}{
		"collection":     vrr.collection,
		"document_count": collectionInfo.DocumentCount,
		"dimension":      collectionInfo.Dimension,
		"index_type":     collectionInfo.IndexType,
		"metric_type":    collectionInfo.MetricType,
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

// GenerateContent overrides the base GenerateContent method to use Ollama
func (org *OllamaRAGGenerator) GenerateContent(inputData string) (string, error) {
	if inputData == "" {
		return "", fmt.Errorf("input data cannot be empty")
	}

	prompt := fmt.Sprintf(`Based on the following retrieved information from our knowledge base, generate a comprehensive and well-structured response:

Retrieved Context:
%s

Please provide a clear, informative, and well-organized answer that incorporates the retrieved information. Focus on accuracy and completeness while maintaining readability.`, inputData)

	response, err := org.ollamaProvider.GenerateResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate content with Ollama: %w", err)
	}

	fmt.Printf("✅ Generated content using Ollama (%d characters)\n", len(response))
	return response, nil
}

// GetModelInfo returns information about the Ollama model
func (org *OllamaRAGGenerator) GetModelInfo() string {
	info := org.ollamaProvider.GetModelInfo()
	return fmt.Sprintf("Model: %s | Provider: %s", info.Name, info.Provider)
}

// VectorRAGWorkflow represents a RAG workflow with vector database integration
type VectorRAGWorkflow struct {
	Retriever *VectorRAGRetriever
	Generator *OllamaRAGGenerator
	Evaluator *rag.Evaluator
}

// NewVectorRAGWorkflow creates a new vector-based RAG workflow
func NewVectorRAGWorkflow(retriever *VectorRAGRetriever, generator *OllamaRAGGenerator, evaluator *rag.Evaluator) *VectorRAGWorkflow {
	return &VectorRAGWorkflow{
		Retriever: retriever,
		Generator: generator,
		Evaluator: evaluator,
	}
}

// Execute runs the RAG workflow with vector search
func (w *VectorRAGWorkflow) Execute(query string) (string, error) {
	ctx := context.Background()

	// Step 1: Vector-based retrieval
	fmt.Println("🔍 Step 1: Vector-based semantic retrieval...")
	retrievedData, err := w.Retriever.Retrieve(ctx, query)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err)
	}
	fmt.Printf("✅ Retrieved context (%d characters)\n", len(retrievedData))

	// Step 2: AI-powered generation
	fmt.Println("🤖 Step 2: AI-powered content generation...")
	generatedContent, err := w.Generator.GenerateContent(retrievedData)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	// Step 3: Content evaluation
	fmt.Println("📋 Step 3: Content quality evaluation...")
	isValid, err := w.Evaluator.Evaluate(generatedContent)
	if err != nil {
		return "", fmt.Errorf("evaluation failed: %w", err)
	}

	if !isValid {
		return "", fmt.Errorf("generated content did not pass quality evaluation")
	}
	fmt.Println("✅ Content evaluation passed!")

	return generatedContent, nil
}

// Vector-Enhanced RAG Workflow Demo
func main() {
	fmt.Println("=== Vector-Enhanced RAG Workflow Demo ===")

	// Configuration
	ollamaURL := "http://192.168.10.10:11434"
	model := "llama3.2"
	collectionName := "knowledge_base"

	// 1. Set up vector database (using In-Memory for demo)
	fmt.Println("🗄️  Setting up vector database...")
	vectorDB := vectordbProviders.NewInMemoryProvider() // Use in-memory for working demo
	embedder := vectordb.NewSimpleEmbeddingProvider(384, "sentence-transformers/all-MiniLM-L6-v2")
	processor := vectordb.NewSimpleDocumentProcessor()

	// Create RAG store
	ragStore := vectordb.NewSimpleRAGStore(vectorDB, embedder, processor)

	// Connect to vector database
	ctx := context.Background()
	if err := ragStore.Connect(ctx); err != nil {
		log.Printf("❌ Failed to connect to vector database: %v", err)
		log.Println("Continuing with mock vector operations...")
	} else {
		fmt.Println("✅ Connected to vector database!")
	}

	// 2. Initialize vector-based RAG components
	fmt.Println("🔧 Initializing vector-enhanced RAG components...")

	vectorRetriever := NewVectorRAGRetriever("vector_knowledge_base", ragStore, collectionName)
	ollamaGenerator := NewOllamaRAGGenerator(model, ollamaURL)
	evaluator := rag.NewEvaluator("QualityChecker", "Evaluate content quality and relevance")

	// 3. Test Ollama connection
	fmt.Printf("🔍 Testing connection to Ollama at %s...\n", ollamaURL)
	if err := ollamaGenerator.ollamaProvider.Ping(); err != nil {
		log.Printf("❌ Failed to connect to Ollama: %v", err)
		log.Println("Make sure Ollama is running and accessible")
		return
	}
	fmt.Println("✅ Successfully connected to Ollama!")
	fmt.Printf("🤖 Using generator: %s\n", ollamaGenerator.GetModelInfo())

	// 4. Create collection and add knowledge
	fmt.Println("\n📚 Setting up knowledge base...")

	// Create collection
	err := ragStore.CreateCollection(ctx, collectionName, embedder.GetDimension(), map[string]interface{}{
		"metric_type": "cosine",
		"index_type":  "IVF_FLAT",
	})
	if err != nil {
		log.Printf("❌ Failed to create collection: %v", err)
	}

	// Add sample knowledge documents
	knowledgeDocs := []string{
		"Machine learning is a subset of artificial intelligence that focuses on algorithms that can learn from and make predictions or decisions based on data. Key principles include supervised learning, unsupervised learning, and reinforcement learning.",
		"Deep learning uses neural networks with multiple layers to model and understand complex patterns in data. It has revolutionized fields like computer vision, natural language processing, and speech recognition.",
		"Supervised learning algorithms learn from labeled training data to make predictions on new, unseen data. Common algorithms include linear regression, decision trees, random forests, and support vector machines.",
		"Unsupervised learning finds hidden patterns in data without labeled examples. Techniques include clustering (like k-means), dimensionality reduction (like PCA), and association rule mining.",
		"Reinforcement learning involves an agent learning to make decisions by interacting with an environment and receiving rewards or penalties. It's used in game playing, robotics, and autonomous systems.",
		"Feature engineering is the process of selecting, modifying, or creating new features from raw data to improve machine learning model performance. It requires domain expertise and understanding of the data.",
		"Cross-validation is a technique for assessing how well a machine learning model will generalize to unseen data. It involves splitting the data into training and validation sets multiple times.",
		"Overfitting occurs when a model learns the training data too well, including noise and outliers, leading to poor performance on new data. Regularization techniques help prevent overfitting.",
	}

	knowledgeMetadata := []map[string]interface{}{
		{"topic": "ML Fundamentals", "difficulty": "beginner", "category": "overview"},
		{"topic": "Deep Learning", "difficulty": "intermediate", "category": "neural_networks"},
		{"topic": "Supervised Learning", "difficulty": "beginner", "category": "algorithms"},
		{"topic": "Unsupervised Learning", "difficulty": "intermediate", "category": "algorithms"},
		{"topic": "Reinforcement Learning", "difficulty": "advanced", "category": "algorithms"},
		{"topic": "Feature Engineering", "difficulty": "intermediate", "category": "preprocessing"},
		{"topic": "Model Validation", "difficulty": "intermediate", "category": "evaluation"},
		{"topic": "Model Generalization", "difficulty": "intermediate", "category": "evaluation"},
	}

	err = vectorRetriever.AddKnowledge(ctx, knowledgeDocs, knowledgeMetadata)
	if err != nil {
		log.Printf("❌ Failed to add knowledge: %v", err)
	}

	// 5. Create and test vector RAG workflow
	fmt.Println("\n🔄 Creating vector-enhanced RAG workflow...")
	vectorWorkflow := NewVectorRAGWorkflow(vectorRetriever, ollamaGenerator, evaluator)

	// Test queries
	queries := []string{
		"What are the main types of machine learning and how do they differ?",
		"Explain overfitting and how to prevent it in machine learning models",
		"What is feature engineering and why is it important?",
	}

	fmt.Println("\n🧪 Testing vector-enhanced RAG workflow...")
	fmt.Println("=" + strings.Repeat("=", 70))

	for i, query := range queries {
		fmt.Printf("\n📝 Test Query %d: %s\n", i+1, query)
		fmt.Println(strings.Repeat("-", 70))

		result, err := vectorWorkflow.Execute(query)
		if err != nil {
			log.Printf("❌ Workflow error: %v\n", err)
			continue
		}

		fmt.Printf("📄 Generated Answer (first 300 chars):\n%s...\n", result[:min(300, len(result))])
	}

	// 6. Show knowledge base statistics
	fmt.Println("\n📊 Knowledge Base Statistics:")
	stats, err := vectorRetriever.GetStats(ctx)
	if err != nil {
		log.Printf("❌ Failed to get stats: %v", err)
	} else {
		for key, value := range stats {
			fmt.Printf("  - %s: %v\n", key, value)
		}
	}

	fmt.Println("\n🎉 Vector-enhanced RAG workflow demo completed!")
	fmt.Println("Successfully demonstrated:")
	fmt.Println("  ✅ Vector database integration")
	fmt.Println("  ✅ Semantic search and retrieval")
	fmt.Println("  ✅ AI-powered content generation")
	fmt.Println("  ✅ Knowledge base management")
	fmt.Println("  ✅ End-to-end RAG pipeline")
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
