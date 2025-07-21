package rag

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RAGEngineImpl implements the RAGEngine interface
type RAGEngineImpl struct {
	config     *RAGConfig
	vectorDB   VectorDB
	embeddings EmbeddingProvider
	chunker    TextChunker
	processors map[string]DocumentProcessor
}

// NewRAGEngine creates a new RAG engine with all components
func NewRAGEngine(config *RAGConfig, vectorDB VectorDB, embeddings EmbeddingProvider, chunker TextChunker) *RAGEngineImpl {
	engine := &RAGEngineImpl{
		config:     config,
		vectorDB:   vectorDB,
		embeddings: embeddings,
		chunker:    chunker,
		processors: make(map[string]DocumentProcessor),
	}

	// Register default processors
	engine.registerDefaultProcessors()

	return engine
}

// registerDefaultProcessors registers built-in document processors
func (r *RAGEngineImpl) registerDefaultProcessors() {
	// Text processor
	r.processors[".txt"] = &TextProcessor{}
	r.processors[".md"] = &MarkdownProcessor{}
	// TODO: Add PDF, DOCX processors
}

// IndexDocument processes and indexes a document into the vector database
func (r *RAGEngineImpl) IndexDocument(ctx context.Context, filePath string, metadata map[string]interface{}) (*Document, error) {
	// Determine file type and get appropriate processor
	ext := strings.ToLower(filepath.Ext(filePath))
	processor, exists := r.processors[ext]
	if !exists {
		// Default to text processor
		processor = &TextProcessor{}
	}

	// Extract content from file
	doc, err := processor.ProcessFile(ctx, filePath, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to process file: %w", err)
	}

	// Store document in database
	if err := r.vectorDB.StoreDocument(ctx, *doc); err != nil {
		return nil, fmt.Errorf("failed to store document: %w", err)
	}

	// Chunk the document content
	chunks, err := r.chunker.ChunkText(ctx, doc.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to chunk text: %w", err)
	}

	// Process chunks and generate embeddings
	var documentChunks []DocumentChunk
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	// Generate embeddings for all chunks in batch
	embeddings, err := r.embeddings.EmbedBatch(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Create document chunks
	for i, chunk := range chunks {
		if embeddings[i] == nil {
			continue // Skip chunks that failed to embed
		}

		chunkMetadata := make(map[string]interface{})
		chunkMetadata["chunk_length"] = len(chunk.Content)
		chunkMetadata["chunk_words"] = len(strings.Fields(chunk.Content))

		documentChunk := DocumentChunk{
			ID:         uuid.New().String(),
			DocumentID: doc.ID,
			Index:      i,
			Content:    chunk.Content,
			Embedding:  embeddings[i],
			Metadata:   chunkMetadata,
			CreatedAt:  time.Now(),
		}

		documentChunks = append(documentChunks, documentChunk)
	}

	// Store chunks in database
	if err := r.vectorDB.StoreChunks(ctx, documentChunks); err != nil {
		return nil, fmt.Errorf("failed to store chunks: %w", err)
	}

	return doc, nil
}

// IndexContent processes and indexes content directly
func (r *RAGEngineImpl) IndexContent(ctx context.Context, content, title, contentType string, metadata map[string]interface{}) (*Document, error) {
	// Get appropriate processor
	processor, exists := r.processors[contentType]
	if !exists {
		processor = &TextProcessor{}
	}

	// Process content
	doc, err := processor.ProcessContent(ctx, content, title, contentType, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to process content: %w", err)
	}

	// Store document in database
	if err := r.vectorDB.StoreDocument(ctx, *doc); err != nil {
		return nil, fmt.Errorf("failed to store document: %w", err)
	}

	// Chunk and index the content
	chunks, err := r.chunker.ChunkText(ctx, doc.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to chunk text: %w", err)
	}

	// Process chunks similar to IndexDocument
	var documentChunks []DocumentChunk
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	embeddings, err := r.embeddings.EmbedBatch(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	for i, chunk := range chunks {
		if embeddings[i] == nil {
			continue
		}

		documentChunk := DocumentChunk{
			ID:         uuid.New().String(),
			DocumentID: doc.ID,
			Index:      i,
			Content:    chunk.Content,
			Embedding:  embeddings[i],
			Metadata:   chunk.Metadata,
			CreatedAt:  time.Now(),
		}

		documentChunks = append(documentChunks, documentChunk)
	}

	if err := r.vectorDB.StoreChunks(ctx, documentChunks); err != nil {
		return nil, fmt.Errorf("failed to store chunks: %w", err)
	}

	return doc, nil
}

// Search performs semantic search using vector similarity
func (r *RAGEngineImpl) Search(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]SearchResult, error) {
	// Generate embedding for query
	queryEmbedding, err := r.embeddings.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search in vector database
	results, err := r.vectorDB.SearchSimilar(ctx, queryEmbedding, limit, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar: %w", err)
	}

	return results, nil
}

// Query performs a RAG query with context retrieval and generation
func (r *RAGEngineImpl) Query(ctx context.Context, question string, maxSources int, filters map[string]interface{}) (*RAGResponse, error) {
	// Search for relevant context
	searchResults, err := r.Search(ctx, question, maxSources, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search for context: %w", err)
	}

	// Build sources from search results
	var sources []Source
	for _, result := range searchResults {
		source := Source{
			DocumentID:    result.Document.ID,
			DocumentTitle: result.Document.Title,
			ChunkContent:  result.Chunk.Content,
			Score:         result.Score,
		}

		// Add section info if available in metadata
		if section, exists := result.Chunk.Metadata["section"]; exists {
			if sectionStr, ok := section.(string); ok {
				source.Section = sectionStr
			}
		}

		sources = append(sources, source)
	}

	// For now, return a simple response without LLM generation
	// In a full implementation, this would call an LLM to generate the answer
	response := &RAGResponse{
		Answer:     fmt.Sprintf("Based on %d sources, here's what I found about: %s", len(sources), question),
		Sources:    sources,
		Question:   question,
		Confidence: calculateConfidence(searchResults),
		Timestamp:  time.Now(),
	}

	return response, nil
}

// GetDocument retrieves a document by ID
func (r *RAGEngineImpl) GetDocument(ctx context.Context, id string) (*Document, error) {
	return r.vectorDB.GetDocument(ctx, id)
}

// DeleteDocument removes a document and its chunks
func (r *RAGEngineImpl) DeleteDocument(ctx context.Context, id string) error {
	return r.vectorDB.DeleteDocument(ctx, id)
}

// ListDocuments returns a paginated list of documents
func (r *RAGEngineImpl) ListDocuments(ctx context.Context, limit, offset int) ([]Document, error) {
	return r.vectorDB.ListDocuments(ctx, limit, offset)
}

// GetDocumentChunks retrieves all chunks for a document
func (r *RAGEngineImpl) GetDocumentChunks(ctx context.Context, documentID string) ([]DocumentChunk, error) {
	return r.vectorDB.GetDocumentChunks(ctx, documentID)
}

// GetStats returns RAG system statistics
func (r *RAGEngineImpl) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats, err := r.vectorDB.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	// Add RAG-specific stats
	stats["embedding_model"] = r.embeddings.GetModel()
	stats["embedding_provider"] = r.embeddings.GetProvider()
	stats["embedding_dimensions"] = r.embeddings.GetDimensions()
	stats["chunk_size"] = r.config.Chunking.Size
	stats["chunk_overlap"] = r.config.Chunking.Overlap

	return stats, nil
}

// HealthCheck checks the health of all RAG components
func (r *RAGEngineImpl) HealthCheck(ctx context.Context) error {
	// Check vector database
	if err := r.vectorDB.Ping(ctx); err != nil {
		return fmt.Errorf("vector database health check failed: %w", err)
	}

	// Check embeddings provider
	if err := r.embeddings.Ping(ctx); err != nil {
		return fmt.Errorf("embeddings provider health check failed: %w", err)
	}

	return nil
}

// UpdateConfig updates the RAG configuration
func (r *RAGEngineImpl) UpdateConfig(config *RAGConfig) error {
	r.config = config

	// Update chunker configuration
	return r.chunker.Configure(config.Chunking.Size, config.Chunking.Overlap, config.Chunking.Strategy)
}

// GetConfig returns the current configuration
func (r *RAGEngineImpl) GetConfig() *RAGConfig {
	return r.config
}

// Helper function to calculate confidence score
func calculateConfidence(results []SearchResult) float64 {
	if len(results) == 0 {
		return 0.0
	}

	// Simple confidence calculation based on top result score
	// In a real implementation, this would be more sophisticated
	return results[0].Score
}
