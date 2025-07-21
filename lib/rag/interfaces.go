package rag

import (
	"context"
	"time"
)

// Document represents a document in the knowledge base
type Document struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	SourcePath  string                 `json:"source_path"`
	ContentType string                 `json:"content_type"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// DocumentChunk represents a chunk of a document with embeddings
type DocumentChunk struct {
	ID         string                 `json:"id"`
	DocumentID string                 `json:"document_id"`
	Index      int                    `json:"index"`
	Content    string                 `json:"content"`
	Embedding  []float32              `json:"embedding"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  time.Time              `json:"created_at"`
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Chunk      DocumentChunk `json:"chunk"`
	Document   Document      `json:"document"`
	Score      float64       `json:"score"`
	Highlights []string      `json:"highlights"`
}

// Source represents a source reference for RAG responses
type Source struct {
	DocumentID    string  `json:"document_id"`
	DocumentTitle string  `json:"document_title"`
	ChunkContent  string  `json:"chunk_content"`
	Score         float64 `json:"score"`
	PageNumber    int     `json:"page_number,omitempty"`
	Section       string  `json:"section,omitempty"`
}

// RAGResponse represents a complete RAG query response
type RAGResponse struct {
	Answer     string    `json:"answer"`
	Sources    []Source  `json:"sources"`
	Question   string    `json:"question"`
	Confidence float64   `json:"confidence"`
	Timestamp  time.Time `json:"timestamp"`
}

// TextChunk represents a chunk of text with metadata
type TextChunk struct {
	Content  string                 `json:"content"`
	Index    int                    `json:"index"`
	Metadata map[string]interface{} `json:"metadata"`
}

// VectorDB interface for vector database operations
type VectorDB interface {
	// Document operations
	StoreDocument(ctx context.Context, doc Document) error
	GetDocument(ctx context.Context, id string) (*Document, error)
	DeleteDocument(ctx context.Context, id string) error
	ListDocuments(ctx context.Context, limit, offset int) ([]Document, error)

	// Chunk operations
	StoreChunks(ctx context.Context, chunks []DocumentChunk) error
	GetChunk(ctx context.Context, id string) (*DocumentChunk, error)
	GetDocumentChunks(ctx context.Context, documentID string) ([]DocumentChunk, error)

	// Search operations
	SearchSimilar(ctx context.Context, embedding []float32, limit int, filters map[string]interface{}) ([]SearchResult, error)
	SearchByText(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]SearchResult, error)

	// Management operations
	CreateIndex(ctx context.Context, indexType string) error
	DropIndex(ctx context.Context, indexType string) error
	GetStats(ctx context.Context) (map[string]interface{}, error)

	// Health check
	Ping(ctx context.Context) error
	Close() error
}

// EmbeddingProvider interface for generating embeddings
type EmbeddingProvider interface {
	// Single embedding
	Embed(ctx context.Context, text string) ([]float32, error)

	// Batch embeddings
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)

	// Provider info
	GetDimensions() int
	GetModel() string
	GetProvider() string

	// Health check
	Ping(ctx context.Context) error
}

// DocumentProcessor interface for processing documents
type DocumentProcessor interface {
	// Process a document file
	ProcessFile(ctx context.Context, filePath string, metadata map[string]interface{}) (*Document, error)

	// Process raw content
	ProcessContent(ctx context.Context, content, title, contentType string, metadata map[string]interface{}) (*Document, error)

	// Extract text from various formats
	ExtractText(ctx context.Context, filePath string) (string, error)

	// Get supported content types
	GetSupportedTypes() []string
}

// TextChunker interface for chunking text
type TextChunker interface {
	// Chunk text into smaller pieces
	ChunkText(ctx context.Context, text string) ([]TextChunk, error)

	// Get chunking strategy
	GetStrategy() string

	// Configure chunking parameters
	Configure(size, overlap int, strategy string) error
}

// VectorSearcher interface for search operations
type VectorSearcher interface {
	// Search by text query
	Search(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]SearchResult, error)

	// Search by embedding vector
	SearchByVector(ctx context.Context, embedding []float32, limit int, filters map[string]interface{}) ([]SearchResult, error)

	// Hybrid search (vector + keyword)
	HybridSearch(ctx context.Context, query string, limit int, vectorWeight, keywordWeight float64, filters map[string]interface{}) ([]SearchResult, error)

	// Get search statistics
	GetSearchStats() map[string]interface{}
}

// RAGEngine interface for the complete RAG system
type RAGEngine interface {
	// Index operations
	IndexDocument(ctx context.Context, filePath string, metadata map[string]interface{}) (*Document, error)
	IndexContent(ctx context.Context, content, title, contentType string, metadata map[string]interface{}) (*Document, error)
	DeleteDocument(ctx context.Context, documentID string) error

	// Query operations
	Query(ctx context.Context, question string, maxSources int, filters map[string]interface{}) (*RAGResponse, error)
	Search(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]SearchResult, error)

	// Management operations
	ListDocuments(ctx context.Context, limit, offset int) ([]Document, error)
	GetDocument(ctx context.Context, documentID string) (*Document, error)
	GetDocumentChunks(ctx context.Context, documentID string) ([]DocumentChunk, error)

	// System operations
	GetStats(ctx context.Context) (map[string]interface{}, error)
	HealthCheck(ctx context.Context) error

	// Configuration
	UpdateConfig(config *RAGConfig) error
	GetConfig() *RAGConfig
}

// LLMProvider interface for language model integration
type LLMProvider interface {
	// Generate response with context
	GenerateWithContext(ctx context.Context, question string, context []string, systemPrompt string) (string, error)

	// Get model information
	GetModel() string
	GetProvider() string

	// Configuration
	SetTemperature(temp float64)
	SetMaxTokens(tokens int)
}

// RAGAgent interface for intelligent RAG operations
type RAGAgent interface {
	// Intelligent query processing
	ProcessQuery(ctx context.Context, question string, options QueryOptions) (*RAGResponse, error)

	// Query analysis and reformulation
	AnalyzeQuery(ctx context.Context, question string) (*QueryAnalysis, error)
	ReformulateQuery(ctx context.Context, originalQuery string, context []string) (string, error)

	// Context management
	BuildContext(ctx context.Context, chunks []DocumentChunk, maxTokens int) ([]string, error)

	// Response generation
	GenerateAnswer(ctx context.Context, question string, context []string) (string, float64, error)

	// Source attribution
	AttributeSources(ctx context.Context, answer string, chunks []DocumentChunk) ([]Source, error)
}

// QueryOptions for advanced query processing
type QueryOptions struct {
	MaxSources      int                    `json:"max_sources"`
	MinConfidence   float64                `json:"min_confidence"`
	Filters         map[string]interface{} `json:"filters"`
	IncludeMetadata bool                   `json:"include_metadata"`
	Rerank          bool                   `json:"rerank"`
	Temperature     float64                `json:"temperature"`
}

// QueryAnalysis represents the analysis of a user query
type QueryAnalysis struct {
	Intent          string                 `json:"intent"`
	Entities        []string               `json:"entities"`
	Keywords        []string               `json:"keywords"`
	Complexity      string                 `json:"complexity"` // "simple", "medium", "complex"
	RequiredSources int                    `json:"required_sources"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// IndexStats represents indexing statistics
type IndexStats struct {
	TotalDocuments int       `json:"total_documents"`
	TotalChunks    int       `json:"total_chunks"`
	IndexSize      int64     `json:"index_size_bytes"`
	LastUpdated    time.Time `json:"last_updated"`
	VectorCount    int       `json:"vector_count"`
	AvgChunkSize   float64   `json:"avg_chunk_size"`
}

// SearchStats represents search performance statistics
type SearchStats struct {
	TotalQueries    int64         `json:"total_queries"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	SuccessRate     float64       `json:"success_rate"`
	LastQuery       time.Time     `json:"last_query"`
	PopularFilters  []string      `json:"popular_filters"`
}
