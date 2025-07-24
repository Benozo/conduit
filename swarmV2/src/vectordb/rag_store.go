package vectordb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SimpleRAGStore implements RAGStore interface
type SimpleRAGStore struct {
	vectorDB  VectorDB
	embedder  EmbeddingProvider
	processor DocumentProcessor
	connected bool
}

// NewSimpleRAGStore creates a new RAG store
func NewSimpleRAGStore(vectorDB VectorDB, embedder EmbeddingProvider, processor DocumentProcessor) *SimpleRAGStore {
	return &SimpleRAGStore{
		vectorDB:  vectorDB,
		embedder:  embedder,
		processor: processor,
	}
}

// Connect establishes connection to the vector database
func (srs *SimpleRAGStore) Connect(ctx context.Context) error {
	err := srs.vectorDB.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to vector database: %w", err)
	}
	srs.connected = true
	return nil
}

// Disconnect closes the connection
func (srs *SimpleRAGStore) Disconnect(ctx context.Context) error {
	srs.connected = false
	return srs.vectorDB.Disconnect(ctx)
}

// Ping tests the connection
func (srs *SimpleRAGStore) Ping(ctx context.Context) error {
	if !srs.connected {
		return fmt.Errorf("not connected to RAG store")
	}
	return srs.vectorDB.Ping(ctx)
}

// AddDocument adds a document with automatic processing and embedding
func (srs *SimpleRAGStore) AddDocument(ctx context.Context, collection string, content string, docType DocumentType, metadata map[string]interface{}) (string, error) {
	if !srs.connected {
		return "", fmt.Errorf("not connected to RAG store")
	}

	// Create document ID
	docID := fmt.Sprintf("doc_%d", time.Now().UnixNano())

	// Create document
	document := Document{
		ID:       docID,
		Content:  content,
		Type:     docType,
		Metadata: metadata,
	}

	// Process document (chunking, etc.)
	processedDocs, err := srs.processor.ProcessDocument(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to process document: %w", err)
	}

	// Generate embeddings and add to vector database
	for _, doc := range processedDocs {
		embedding, err := srs.embedder.GenerateTextEmbedding(ctx, doc.Content)
		if err != nil {
			return "", fmt.Errorf("failed to generate embedding: %w", err)
		}

		// Add metadata about embedding
		doc.Metadata["embedding_model"] = srs.embedder.GetModelName()
		doc.Metadata["embedding_dimension"] = srs.embedder.GetDimension()
		doc.Metadata["processed_at"] = time.Now().UTC()

		// Create vector
		vector := Vector{
			ID:       doc.ID,
			Values:   embedding,
			Metadata: doc.Metadata,
		}

		// Add document and vector
		err = srs.vectorDB.AddDocuments(ctx, collection, []Document{doc})
		if err != nil {
			return "", fmt.Errorf("failed to add document: %w", err)
		}

		err = srs.vectorDB.AddVectors(ctx, collection, []Vector{vector})
		if err != nil {
			return "", fmt.Errorf("failed to add vector: %w", err)
		}
	}

	return docID, nil
}

// AddDocumentFromFile adds a document from a file
func (srs *SimpleRAGStore) AddDocumentFromFile(ctx context.Context, collection string, filePath string, metadata map[string]interface{}) (string, error) {
	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Determine document type from file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	var docType DocumentType
	switch ext {
	case ".txt", ".md":
		docType = DocumentTypeText
	case ".json":
		docType = DocumentTypeJSON
	case ".pdf":
		docType = DocumentTypePDF
	default:
		docType = DocumentTypeText // Default to text
	}

	// Extract text content
	textContent, err := srs.processor.ExtractText(ctx, content, docType)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from file: %w", err)
	}

	// Add file metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["source_file"] = filePath
	metadata["file_extension"] = ext
	metadata["file_size"] = len(content)

	return srs.AddDocument(ctx, collection, textContent, docType, metadata)
}

// AddDocumentsFromDirectory adds all documents from a directory
func (srs *SimpleRAGStore) AddDocumentsFromDirectory(ctx context.Context, collection string, dirPath string, recursive bool) ([]string, error) {
	var docIDs []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			if !recursive && path != dirPath {
				return filepath.SkipDir
			}
			return nil
		}

		// Add file
		docID, err := srs.AddDocumentFromFile(ctx, collection, path, map[string]interface{}{
			"directory":     dirPath,
			"relative_path": strings.TrimPrefix(path, dirPath),
		})
		if err != nil {
			return fmt.Errorf("failed to add file %s: %w", path, err)
		}

		docIDs = append(docIDs, docID)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", dirPath, err)
	}

	return docIDs, nil
}

// SemanticSearch performs semantic search using embeddings
func (srs *SimpleRAGStore) SemanticSearch(ctx context.Context, collection string, query string, options SearchOptions) ([]SearchResult, error) {
	if !srs.connected {
		return nil, fmt.Errorf("not connected to RAG store")
	}

	// Generate query embedding
	queryEmbedding, err := srs.embedder.GenerateTextEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search vectors
	return srs.vectorDB.SearchVectors(ctx, collection, queryEmbedding, options)
}

// HybridSearch combines semantic and keyword search
func (srs *SimpleRAGStore) HybridSearch(ctx context.Context, collection string, query string, options SearchOptions) ([]SearchResult, error) {
	// For now, just use semantic search
	// In a real implementation, you'd combine with keyword search
	return srs.SemanticSearch(ctx, collection, query, options)
}

// RetrieveForRAG retrieves relevant documents for RAG applications
func (srs *SimpleRAGStore) RetrieveForRAG(ctx context.Context, collection string, query string, topK int) (string, error) {
	options := SearchOptions{
		TopK:           topK,
		IncludeContent: true,
		IncludeVector:  false,
	}

	results, err := srs.SemanticSearch(ctx, collection, query, options)
	if err != nil {
		return "", fmt.Errorf("failed to search for RAG: %w", err)
	}

	// Combine results into a single context string
	var contexts []string
	for i, result := range results {
		context := fmt.Sprintf("Document %d (Score: %.3f):\n%s", i+1, result.Score, result.Document.Content)
		contexts = append(contexts, context)
	}

	return strings.Join(contexts, "\n\n"), nil
}

// RetrieveWithContext retrieves documents with additional context window
func (srs *SimpleRAGStore) RetrieveWithContext(ctx context.Context, collection string, query string, topK int, contextWindow int) (string, error) {
	// For now, just use basic retrieval
	// In a real implementation, you'd expand context by retrieving neighboring chunks
	return srs.RetrieveForRAG(ctx, collection, query, topK)
}

// Delegate other methods to the underlying vector database
func (srs *SimpleRAGStore) CreateCollection(ctx context.Context, name string, dimension int, options map[string]interface{}) error {
	return srs.vectorDB.CreateCollection(ctx, name, dimension, options)
}

func (srs *SimpleRAGStore) DeleteCollection(ctx context.Context, name string) error {
	return srs.vectorDB.DeleteCollection(ctx, name)
}

func (srs *SimpleRAGStore) ListCollections(ctx context.Context) ([]CollectionInfo, error) {
	return srs.vectorDB.ListCollections(ctx)
}

func (srs *SimpleRAGStore) GetCollectionInfo(ctx context.Context, name string) (*CollectionInfo, error) {
	return srs.vectorDB.GetCollectionInfo(ctx, name)
}

func (srs *SimpleRAGStore) AddDocuments(ctx context.Context, collection string, documents []Document) error {
	return srs.vectorDB.AddDocuments(ctx, collection, documents)
}

func (srs *SimpleRAGStore) UpdateDocument(ctx context.Context, collection string, document Document) error {
	return srs.vectorDB.UpdateDocument(ctx, collection, document)
}

func (srs *SimpleRAGStore) DeleteDocument(ctx context.Context, collection string, docID string) error {
	return srs.vectorDB.DeleteDocument(ctx, collection, docID)
}

func (srs *SimpleRAGStore) GetDocument(ctx context.Context, collection string, docID string) (*Document, error) {
	return srs.vectorDB.GetDocument(ctx, collection, docID)
}

func (srs *SimpleRAGStore) AddVectors(ctx context.Context, collection string, vectors []Vector) error {
	return srs.vectorDB.AddVectors(ctx, collection, vectors)
}

func (srs *SimpleRAGStore) SearchVectors(ctx context.Context, collection string, queryVector []float32, options SearchOptions) ([]SearchResult, error) {
	return srs.vectorDB.SearchVectors(ctx, collection, queryVector, options)
}

func (srs *SimpleRAGStore) SearchDocuments(ctx context.Context, collection string, query string, options SearchOptions) ([]SearchResult, error) {
	return srs.vectorDB.SearchDocuments(ctx, collection, query, options)
}

func (srs *SimpleRAGStore) BatchAddDocuments(ctx context.Context, collection string, documents []Document) error {
	return srs.vectorDB.BatchAddDocuments(ctx, collection, documents)
}

func (srs *SimpleRAGStore) BatchDeleteDocuments(ctx context.Context, collection string, docIDs []string) error {
	return srs.vectorDB.BatchDeleteDocuments(ctx, collection, docIDs)
}

func (srs *SimpleRAGStore) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return srs.vectorDB.GetStats(ctx)
}
