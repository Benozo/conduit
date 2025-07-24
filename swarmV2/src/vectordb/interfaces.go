package vectordb

import (
	"context"
	"time"
)

// Vector represents a high-dimensional vector with metadata
type Vector struct {
	ID       string                 `json:"id"`
	Values   []float32              `json:"values"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Document represents a document to be indexed
type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Type     DocumentType           `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

// DocumentType defines the type of document
type DocumentType string

const (
	DocumentTypeText  DocumentType = "text"
	DocumentTypeImage DocumentType = "image"
	DocumentTypePDF   DocumentType = "pdf"
	DocumentTypeJSON  DocumentType = "json"
)

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Document   Document `json:"document"`
	Score      float32  `json:"score"`
	Vector     Vector   `json:"vector"`
	Distance   float32  `json:"distance"`
	Highlights []string `json:"highlights,omitempty"`
}

// SearchOptions defines search parameters
type SearchOptions struct {
	TopK           int                    `json:"top_k"`
	ScoreThreshold float32                `json:"score_threshold,omitempty"`
	Filter         map[string]interface{} `json:"filter,omitempty"`
	IncludeVector  bool                   `json:"include_vector"`
	IncludeContent bool                   `json:"include_content"`
}

// CollectionInfo represents collection metadata
type CollectionInfo struct {
	Name          string            `json:"name"`
	Dimension     int               `json:"dimension"`
	IndexType     string            `json:"index_type"`
	MetricType    string            `json:"metric_type"`
	DocumentCount int64             `json:"document_count"`
	CreatedAt     time.Time         `json:"created_at"`
	Properties    map[string]string `json:"properties,omitempty"`
}

// VectorDB defines the interface for vector database operations
type VectorDB interface {
	// Collection management
	CreateCollection(ctx context.Context, name string, dimension int, options map[string]interface{}) error
	DeleteCollection(ctx context.Context, name string) error
	ListCollections(ctx context.Context) ([]CollectionInfo, error)
	GetCollectionInfo(ctx context.Context, name string) (*CollectionInfo, error)

	// Document operations
	AddDocuments(ctx context.Context, collection string, documents []Document) error
	UpdateDocument(ctx context.Context, collection string, document Document) error
	DeleteDocument(ctx context.Context, collection string, docID string) error
	GetDocument(ctx context.Context, collection string, docID string) (*Document, error)

	// Vector operations
	AddVectors(ctx context.Context, collection string, vectors []Vector) error
	SearchVectors(ctx context.Context, collection string, queryVector []float32, options SearchOptions) ([]SearchResult, error)
	SearchDocuments(ctx context.Context, collection string, query string, options SearchOptions) ([]SearchResult, error)

	// Batch operations
	BatchAddDocuments(ctx context.Context, collection string, documents []Document) error
	BatchDeleteDocuments(ctx context.Context, collection string, docIDs []string) error

	// Connection management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context) error
	GetStats(ctx context.Context) (map[string]interface{}, error)
}

// EmbeddingProvider defines the interface for generating embeddings
type EmbeddingProvider interface {
	// Text embeddings
	GenerateTextEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateTextEmbeddings(ctx context.Context, texts []string) ([][]float32, error)

	// Image embeddings
	GenerateImageEmbedding(ctx context.Context, imageData []byte) ([]float32, error)
	GenerateImageEmbeddings(ctx context.Context, images [][]byte) ([][]float32, error)

	// Multimodal embeddings
	GenerateMultimodalEmbedding(ctx context.Context, text string, imageData []byte) ([]float32, error)

	// Provider info
	GetDimension() int
	GetModelName() string
	GetProviderName() string
}

// DocumentProcessor defines the interface for processing different document types
type DocumentProcessor interface {
	// Process documents into chunks
	ProcessDocument(ctx context.Context, document Document) ([]Document, error)
	ProcessDocuments(ctx context.Context, documents []Document) ([]Document, error)

	// Extract text from different formats
	ExtractText(ctx context.Context, data []byte, docType DocumentType) (string, error)

	// Chunk text into smaller pieces
	ChunkText(text string, chunkSize int, overlap int) ([]string, error)

	// Supported document types
	GetSupportedTypes() []DocumentType
}

// RAGStore combines vector database and document processing capabilities
type RAGStore interface {
	VectorDB

	// High-level document management
	AddDocument(ctx context.Context, collection string, content string, docType DocumentType, metadata map[string]interface{}) (string, error)
	AddDocumentFromFile(ctx context.Context, collection string, filePath string, metadata map[string]interface{}) (string, error)
	AddDocumentsFromDirectory(ctx context.Context, collection string, dirPath string, recursive bool) ([]string, error)

	// Semantic search
	SemanticSearch(ctx context.Context, collection string, query string, options SearchOptions) ([]SearchResult, error)
	HybridSearch(ctx context.Context, collection string, query string, options SearchOptions) ([]SearchResult, error)

	// Document retrieval for RAG
	RetrieveForRAG(ctx context.Context, collection string, query string, topK int) (string, error)
	RetrieveWithContext(ctx context.Context, collection string, query string, topK int, contextWindow int) (string, error)
}
