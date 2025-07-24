package main

import (
	"context"
	"fmt"
	"log"

	"github.com/benozo/neuron/src/agents/rag"
	"github.com/benozo/neuron/src/llm/providers"
)

// OllamaRAGGenerator wraps the RAG generator with Ollama capabilities
type OllamaRAGGenerator struct {
	*rag.Generator
	ollamaProvider *providers.OllamaProvider
}

// NewOllamaRAGGenerator creates a new RAG generator with Ollama integration
func NewOllamaRAGGenerator(modelName string, ollamaURL string) *OllamaRAGGenerator {
	// Create the Ollama provider
	ollamaProvider := providers.NewOllamaProvider(ollamaURL, modelName)

	// Create the base RAG generator
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

	// Create a prompt for content generation based on retrieved data
	prompt := fmt.Sprintf(`Based on the following retrieved information, generate a comprehensive and well-structured response:

Retrieved Information:
%s

Please provide a clear, informative, and well-organized answer that incorporates the retrieved information.`, inputData)

	// Generate response using Ollama
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

// CustomRAGWorkflow represents a RAG workflow with Ollama integration
type CustomRAGWorkflow struct {
	Retriever       *rag.Retriever
	OllamaGenerator *OllamaRAGGenerator
	Evaluator       *rag.Evaluator
}

// NewCustomRAGWorkflow creates a new custom RAG workflow
func NewCustomRAGWorkflow(retriever *rag.Retriever, generator *OllamaRAGGenerator, evaluator *rag.Evaluator) *CustomRAGWorkflow {
	return &CustomRAGWorkflow{
		Retriever:       retriever,
		OllamaGenerator: generator,
		Evaluator:       evaluator,
	}
}

// Execute runs the RAG workflow with Ollama generation
func (w *CustomRAGWorkflow) Execute(query string) (string, error) {
	// Step 1: Retrieve information
	fmt.Println("🔍 Step 1: Retrieving relevant information...")
	ctx := context.Background()
	retrievedData, err := w.Retriever.Retrieve(ctx, query)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err)
	}
	fmt.Printf("✅ Retrieved information: %s\n", retrievedData)

	// Step 2: Generate content using Ollama
	fmt.Println("🤖 Step 2: Generating content with AI...")
	generatedContent, err := w.OllamaGenerator.GenerateContent(retrievedData)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	// Step 3: Evaluate the generated content
	fmt.Println("📋 Step 3: Evaluating content quality...")
	isValid, err := w.Evaluator.Evaluate(generatedContent)
	if err != nil {
		return "", fmt.Errorf("evaluation failed: %w", err)
	}

	if !isValid {
		return "", fmt.Errorf("generated content is not valid")
	}
	fmt.Println("✅ Content evaluation passed!")

	return generatedContent, nil
}

// Enhanced RAG Workflow Example with Ollama Integration
// This example demonstrates the interaction between RAG agents with AI-powered content generation.
func main() {
	fmt.Println("=== Enhanced RAG Workflow Demo with Ollama ===")

	// Configuration
	ollamaURL := "http://192.168.10.10:11434"
	model := "llama3.2"

	// Initialize RAG agents
	retriever := rag.NewRetriever("knowledge_base")
	ollamaGenerator := NewOllamaRAGGenerator(model, ollamaURL)
	evaluator := rag.NewEvaluator("QualityChecker", "Evaluate content quality and relevance")

	// Test Ollama connection
	fmt.Printf("🔍 Testing connection to Ollama at %s...\n", ollamaURL)
	if err := ollamaGenerator.ollamaProvider.Ping(); err != nil {
		log.Printf("❌ Failed to connect to Ollama: %v", err)
		log.Println("Make sure Ollama is running and accessible")
		return
	}
	fmt.Println("✅ Successfully connected to Ollama!")
	fmt.Printf("🤖 Using generator: %s\n\n", ollamaGenerator.GetModelInfo())

	// Create custom RAG workflow with Ollama generator
	ragWorkflow := NewCustomRAGWorkflow(retriever, ollamaGenerator, evaluator)

	// Execute workflow with a sample query
	query := "What are the main principles of machine learning and how do they apply in real-world applications?"
	fmt.Printf("📝 Query: %s\n\n", query)

	fmt.Println("🔄 Executing RAG workflow...")
	fmt.Println("=" + string(make([]byte, 60)) + "=")

	result, err := ragWorkflow.Execute(query)
	if err != nil {
		fmt.Printf("❌ RAG workflow error: %v\n", err)
		return
	}

	fmt.Printf("\n🎉 RAG workflow completed successfully!\n")
	fmt.Printf("📄 Generated content preview (first 300 chars):\n")
	resultStr := fmt.Sprintf("%v", result)
	fmt.Printf("%s...\n\n", resultStr[:min(300, len(resultStr))])

	fmt.Println("✅ Workflow Summary:")
	fmt.Println("  1. ✅ Retrieved relevant information from knowledge base")
	fmt.Println("  2. ✅ Generated AI-powered content using Ollama")
	fmt.Println("  3. ✅ Evaluated content quality and relevance")
	fmt.Printf("  4. ✅ Final result: %d characters of generated content\n", len(resultStr))
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
