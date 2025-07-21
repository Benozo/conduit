package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/benozo/conduit/lib/rag"
	"github.com/benozo/conduit/mcp"
)

// RAGToolRegistry holds references to RAG system components
type RAGToolRegistry struct {
	engine rag.RAGEngine
	config *rag.RAGConfig
}

// NewRAGToolRegistry creates a new RAG tool registry
func NewRAGToolRegistry(engine rag.RAGEngine, config *rag.RAGConfig) *RAGToolRegistry {
	return &RAGToolRegistry{
		engine: engine,
		config: config,
	}
}

// Global RAG engine instance (to be set during initialization)
var globalRAGEngine rag.RAGEngine

// SetRAGEngine sets the global RAG engine instance
func SetRAGEngine(engine rag.RAGEngine) {
	globalRAGEngine = engine
}

// IndexDocumentFunc indexes a document into the knowledge base
var IndexDocumentFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	// Extract parameters
	filePath, ok := params["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path is required and must be a string")
	}

	title, ok := params["title"].(string)
	if !ok {
		title = filepath.Base(filePath)
	}

	metadata := make(map[string]interface{})
	if metaParam, exists := params["metadata"]; exists {
		if metaMap, ok := metaParam.(map[string]interface{}); ok {
			metadata = metaMap
		}
	}

	// Add timestamp and source info to metadata
	metadata["indexed_at"] = time.Now().Unix()
	metadata["source_path"] = filePath
	metadata["title"] = title

	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Index the document
	doc, err := globalRAGEngine.IndexDocument(ctx, filePath, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to index document: %w", err)
	}

	return map[string]interface{}{
		"status":      "indexed",
		"document_id": doc.ID,
		"title":       doc.Title,
		"file_path":   filePath,
		"indexed_at":  metadata["indexed_at"],
		"size_bytes":  len(doc.Content),
	}, nil
}

// SemanticSearchFunc performs vector similarity search
var SemanticSearchFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	// Extract parameters
	query, ok := params["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query is required and must be a string")
	}

	limit := 10 // default
	if limitParam, exists := params["limit"]; exists {
		if limitFloat, ok := limitParam.(float64); ok {
			limit = int(limitFloat)
		}
	}

	filters := make(map[string]interface{})
	if filterParam, exists := params["filters"]; exists {
		if filterMap, ok := filterParam.(map[string]interface{}); ok {
			filters = filterMap
		}
	}

	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform search
	results, err := globalRAGEngine.Search(ctx, query, limit, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}

	// Format results for response
	searchResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		searchResults[i] = map[string]interface{}{
			"score":          result.Score,
			"content":        result.Chunk.Content,
			"document_id":    result.Document.ID,
			"document_title": result.Document.Title,
			"chunk_index":    result.Chunk.Index,
			"metadata":       result.Chunk.Metadata,
		}
	}

	return map[string]interface{}{
		"query":   query,
		"results": searchResults,
		"count":   len(results),
	}, nil
}

// KnowledgeQueryFunc performs RAG-powered question answering
var KnowledgeQueryFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	// Extract parameters
	question, ok := params["question"].(string)
	if !ok {
		return nil, fmt.Errorf("question is required and must be a string")
	}

	maxSources := 5 // default
	if maxParam, exists := params["max_sources"]; exists {
		if maxFloat, ok := maxParam.(float64); ok {
			maxSources = int(maxFloat)
		}
	}

	filters := make(map[string]interface{})
	if filterParam, exists := params["filters"]; exists {
		if filterMap, ok := filterParam.(map[string]interface{}); ok {
			filters = filterMap
		}
	}

	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform RAG query
	result, err := globalRAGEngine.Query(ctx, question, maxSources, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to perform RAG query: %w", err)
	}

	// Format sources for response
	sources := make([]map[string]interface{}, len(result.Sources))
	for i, source := range result.Sources {
		sources[i] = map[string]interface{}{
			"document_id":    source.DocumentID,
			"document_title": source.DocumentTitle,
			"chunk_content":  source.ChunkContent,
			"score":          source.Score,
			"section":        source.Section,
		}
	}

	return map[string]interface{}{
		"question":   question,
		"answer":     result.Answer,
		"sources":    sources,
		"confidence": result.Confidence,
		"timestamp":  result.Timestamp.Unix(),
	}, nil
}

// ListDocumentsFunc lists documents in the knowledge base
var ListDocumentsFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	limit := 50 // default
	if limitParam, exists := params["limit"]; exists {
		if limitFloat, ok := limitParam.(float64); ok {
			limit = int(limitFloat)
		}
	}

	offset := 0 // default
	if offsetParam, exists := params["offset"]; exists {
		if offsetFloat, ok := offsetParam.(float64); ok {
			offset = int(offsetFloat)
		}
	}

	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List documents
	documents, err := globalRAGEngine.ListDocuments(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	// Format documents for response
	docList := make([]map[string]interface{}, len(documents))
	for i, doc := range documents {
		docList[i] = map[string]interface{}{
			"id":           doc.ID,
			"title":        doc.Title,
			"source_path":  doc.SourcePath,
			"content_type": doc.ContentType,
			"created_at":   doc.CreatedAt.Unix(),
			"updated_at":   doc.UpdatedAt.Unix(),
			"metadata":     doc.Metadata,
			"size_bytes":   len(doc.Content),
		}
	}

	return map[string]interface{}{
		"documents": docList,
		"count":     len(documents),
		"limit":     limit,
		"offset":    offset,
	}, nil
}

// GetDocumentFunc retrieves a specific document
var GetDocumentFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	documentID, ok := params["document_id"].(string)
	if !ok {
		return nil, fmt.Errorf("document_id is required and must be a string")
	}

	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get document
	doc, err := globalRAGEngine.GetDocument(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	return map[string]interface{}{
		"id":           doc.ID,
		"title":        doc.Title,
		"content":      doc.Content,
		"source_path":  doc.SourcePath,
		"content_type": doc.ContentType,
		"created_at":   doc.CreatedAt.Unix(),
		"updated_at":   doc.UpdatedAt.Unix(),
		"metadata":     doc.Metadata,
		"size_bytes":   len(doc.Content),
	}, nil
}

// DeleteDocumentFunc removes a document from the knowledge base
var DeleteDocumentFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	documentID, ok := params["document_id"].(string)
	if !ok {
		return nil, fmt.Errorf("document_id is required and must be a string")
	}

	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete document
	err := globalRAGEngine.DeleteDocument(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}

	return map[string]interface{}{
		"status":      "deleted",
		"document_id": documentID,
		"deleted_at":  time.Now().Unix(),
	}, nil
}

// GetDocumentChunksFunc retrieves chunks for a specific document
var GetDocumentChunksFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	documentID, ok := params["document_id"].(string)
	if !ok {
		return nil, fmt.Errorf("document_id is required and must be a string")
	}

	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get document chunks
	chunks, err := globalRAGEngine.GetDocumentChunks(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document chunks: %w", err)
	}

	// Format chunks for response
	chunkList := make([]map[string]interface{}, len(chunks))
	for i, chunk := range chunks {
		chunkList[i] = map[string]interface{}{
			"id":          chunk.ID,
			"document_id": chunk.DocumentID,
			"index":       chunk.Index,
			"content":     chunk.Content,
			"metadata":    chunk.Metadata,
			"created_at":  chunk.CreatedAt.Unix(),
			"size_bytes":  len(chunk.Content),
		}
	}

	return map[string]interface{}{
		"chunks":      chunkList,
		"document_id": documentID,
		"count":       len(chunks),
	}, nil
}

// GetRAGStatsFunc returns statistics about the RAG system
var GetRAGStatsFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	// Get RAG engine
	if globalRAGEngine == nil {
		return nil, fmt.Errorf("RAG engine not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get statistics
	stats, err := globalRAGEngine.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get RAG stats: %w", err)
	}

	return map[string]interface{}{
		"stats":        stats,
		"retrieved_at": time.Now().Unix(),
	}, nil
}

// RAGToolMetadata provides metadata for RAG tools
type RAGToolMetadata struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Required    []string
}

// GetRAGToolMetadata returns metadata for all RAG tools
func GetRAGToolMetadata() []RAGToolMetadata {
	return []RAGToolMetadata{
		{
			Name:        "index_document",
			Description: "Index a document into the knowledge base for semantic search and retrieval",
			Parameters: map[string]interface{}{
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the document file to index",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Title for the document (optional, defaults to filename)",
				},
				"metadata": map[string]interface{}{
					"type":        "object",
					"description": "Additional metadata for the document",
				},
			},
			Required: []string{"file_path"},
		},
		{
			Name:        "semantic_search",
			Description: "Search for documents using semantic similarity",
			Parameters: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query text",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of results to return (default: 10)",
				},
				"filters": map[string]interface{}{
					"type":        "object",
					"description": "Metadata filters to apply to search",
				},
			},
			Required: []string{"query"},
		},
		{
			Name:        "knowledge_query",
			Description: "Ask a question and get an AI-generated answer based on indexed knowledge",
			Parameters: map[string]interface{}{
				"question": map[string]interface{}{
					"type":        "string",
					"description": "Question to ask about the knowledge base",
				},
				"max_sources": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of source documents to use (default: 5)",
				},
				"filters": map[string]interface{}{
					"type":        "object",
					"description": "Metadata filters to apply when selecting sources",
				},
			},
			Required: []string{"question"},
		},
		{
			Name:        "list_documents",
			Description: "List documents in the knowledge base",
			Parameters: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of documents to return (default: 50)",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Number of documents to skip (default: 0)",
				},
			},
			Required: []string{},
		},
		{
			Name:        "get_document",
			Description: "Retrieve a specific document by ID",
			Parameters: map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the document to retrieve",
				},
			},
			Required: []string{"document_id"},
		},
		{
			Name:        "delete_document",
			Description: "Remove a document from the knowledge base",
			Parameters: map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the document to delete",
				},
			},
			Required: []string{"document_id"},
		},
		{
			Name:        "get_document_chunks",
			Description: "Retrieve all chunks for a specific document",
			Parameters: map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the document to get chunks for",
				},
			},
			Required: []string{"document_id"},
		},
		{
			Name:        "get_rag_stats",
			Description: "Get statistics about the RAG system",
			Parameters:  map[string]interface{}{},
			Required:    []string{},
		},
	}
}
