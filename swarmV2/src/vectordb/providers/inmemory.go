package providers

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/benozo/neuron/src/vectordb"
)

// InMemoryProvider implements VectorDB interface with in-memory storage for demos
type InMemoryProvider struct {
	collections map[string]*InMemoryCollection
	connected   bool
}

// InMemoryCollection represents a collection stored in memory
type InMemoryCollection struct {
	Name       string
	Dimension  int
	Documents  map[string]vectordb.Document
	Vectors    map[string]vectordb.Vector
	IndexType  string
	MetricType string
	CreatedAt  time.Time
}

// NewInMemoryProvider creates a new in-memory vector database provider
func NewInMemoryProvider() *InMemoryProvider {
	return &InMemoryProvider{
		collections: make(map[string]*InMemoryCollection),
	}
}

// Connect establishes connection (always succeeds for in-memory)
func (imp *InMemoryProvider) Connect(ctx context.Context) error {
	log.Printf("Connecting to In-Memory Vector Database")
	imp.connected = true
	return nil
}

// Disconnect closes the connection
func (imp *InMemoryProvider) Disconnect(ctx context.Context) error {
	imp.connected = false
	return nil
}

// Ping tests the connection
func (imp *InMemoryProvider) Ping(ctx context.Context) error {
	if !imp.connected {
		return fmt.Errorf("not connected to in-memory database")
	}
	return nil
}

// CreateCollection creates a new in-memory collection
func (imp *InMemoryProvider) CreateCollection(ctx context.Context, name string, dimension int, options map[string]interface{}) error {
	if !imp.connected {
		return fmt.Errorf("not connected to in-memory database")
	}

	indexType := "flat"
	metricType := "cosine"

	if it, ok := options["index_type"].(string); ok {
		indexType = it
	}
	if mt, ok := options["metric_type"].(string); ok {
		metricType = mt
	}

	imp.collections[name] = &InMemoryCollection{
		Name:       name,
		Dimension:  dimension,
		Documents:  make(map[string]vectordb.Document),
		Vectors:    make(map[string]vectordb.Vector),
		IndexType:  indexType,
		MetricType: metricType,
		CreatedAt:  time.Now(),
	}

	log.Printf("Created in-memory collection '%s' with dimension %d", name, dimension)
	return nil
}

// DeleteCollection removes the collection
func (imp *InMemoryProvider) DeleteCollection(ctx context.Context, name string) error {
	if !imp.connected {
		return fmt.Errorf("not connected to in-memory database")
	}

	delete(imp.collections, name)
	log.Printf("Deleted in-memory collection '%s'", name)
	return nil
}

// ListCollections returns list of collections
func (imp *InMemoryProvider) ListCollections(ctx context.Context) ([]vectordb.CollectionInfo, error) {
	if !imp.connected {
		return nil, fmt.Errorf("not connected to in-memory database")
	}

	var collections []vectordb.CollectionInfo
	for _, col := range imp.collections {
		collections = append(collections, vectordb.CollectionInfo{
			Name:          col.Name,
			Dimension:     col.Dimension,
			IndexType:     col.IndexType,
			MetricType:    col.MetricType,
			DocumentCount: int64(len(col.Documents)),
			CreatedAt:     col.CreatedAt,
		})
	}

	return collections, nil
}

// GetCollectionInfo returns collection metadata
func (imp *InMemoryProvider) GetCollectionInfo(ctx context.Context, name string) (*vectordb.CollectionInfo, error) {
	if !imp.connected {
		return nil, fmt.Errorf("not connected to in-memory database")
	}

	col, exists := imp.collections[name]
	if !exists {
		return nil, fmt.Errorf("collection '%s' not found", name)
	}

	return &vectordb.CollectionInfo{
		Name:          col.Name,
		Dimension:     col.Dimension,
		IndexType:     col.IndexType,
		MetricType:    col.MetricType,
		DocumentCount: int64(len(col.Documents)),
		CreatedAt:     col.CreatedAt,
	}, nil
}

// AddDocuments stores documents in memory
func (imp *InMemoryProvider) AddDocuments(ctx context.Context, collection string, documents []vectordb.Document) error {
	if !imp.connected {
		return fmt.Errorf("not connected to in-memory database")
	}

	col, exists := imp.collections[collection]
	if !exists {
		return fmt.Errorf("collection '%s' not found", collection)
	}

	for _, doc := range documents {
		col.Documents[doc.ID] = doc
	}

	log.Printf("Added %d documents to in-memory collection '%s'", len(documents), collection)
	return nil
}

// AddVectors stores vectors in memory
func (imp *InMemoryProvider) AddVectors(ctx context.Context, collection string, vectors []vectordb.Vector) error {
	if !imp.connected {
		return fmt.Errorf("not connected to in-memory database")
	}

	col, exists := imp.collections[collection]
	if !exists {
		return fmt.Errorf("collection '%s' not found", collection)
	}

	for _, vector := range vectors {
		if len(vector.Values) != col.Dimension {
			return fmt.Errorf("vector dimension %d does not match collection dimension %d", len(vector.Values), col.Dimension)
		}
		col.Vectors[vector.ID] = vector
	}

	log.Printf("Added %d vectors to in-memory collection '%s'", len(vectors), collection)
	return nil
}

// SearchVectors performs similarity search on vectors
func (imp *InMemoryProvider) SearchVectors(ctx context.Context, collection string, queryVector []float32, options vectordb.SearchOptions) ([]vectordb.SearchResult, error) {
	if !imp.connected {
		return nil, fmt.Errorf("not connected to in-memory database")
	}

	col, exists := imp.collections[collection]
	if !exists {
		return nil, fmt.Errorf("collection '%s' not found", collection)
	}

	type scoreResult struct {
		vectorID string
		score    float32
		distance float32
	}

	var scores []scoreResult

	// Calculate similarity for each vector
	for vectorID, vector := range col.Vectors {
		similarity := cosineSimilarity(queryVector, vector.Values)
		distance := 1.0 - similarity

		// Apply score threshold if specified
		if options.ScoreThreshold > 0 && similarity < options.ScoreThreshold {
			continue
		}

		scores = append(scores, scoreResult{
			vectorID: vectorID,
			score:    similarity,
			distance: distance,
		})
	}

	// Sort by similarity (highest first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Limit results
	if options.TopK > 0 && len(scores) > options.TopK {
		scores = scores[:options.TopK]
	}

	// Build results
	var results []vectordb.SearchResult
	for _, score := range scores {
		vector := col.Vectors[score.vectorID]
		document, docExists := col.Documents[score.vectorID]

		result := vectordb.SearchResult{
			Score:    score.score,
			Distance: score.distance,
		}

		if options.IncludeVector {
			result.Vector = vector
		}

		if docExists && options.IncludeContent {
			result.Document = document
		}

		results = append(results, result)
	}

	return results, nil
}

// SearchDocuments performs text-based search (fallback to vector search)
func (imp *InMemoryProvider) SearchDocuments(ctx context.Context, collection string, query string, options vectordb.SearchOptions) ([]vectordb.SearchResult, error) {
	// For simplicity, this would need an embedding service to convert query to vector
	// For now, return empty results
	return []vectordb.SearchResult{}, nil
}

// Helper function to calculate cosine similarity
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}

// Stub implementations for other required methods
func (imp *InMemoryProvider) UpdateDocument(ctx context.Context, collection string, document vectordb.Document) error {
	return imp.AddDocuments(ctx, collection, []vectordb.Document{document})
}

func (imp *InMemoryProvider) DeleteDocument(ctx context.Context, collection string, docID string) error {
	if col, exists := imp.collections[collection]; exists {
		delete(col.Documents, docID)
		delete(col.Vectors, docID)
	}
	return nil
}

func (imp *InMemoryProvider) GetDocument(ctx context.Context, collection string, docID string) (*vectordb.Document, error) {
	if col, exists := imp.collections[collection]; exists {
		if doc, docExists := col.Documents[docID]; docExists {
			return &doc, nil
		}
	}
	return nil, fmt.Errorf("document '%s' not found", docID)
}

func (imp *InMemoryProvider) BatchAddDocuments(ctx context.Context, collection string, documents []vectordb.Document) error {
	return imp.AddDocuments(ctx, collection, documents)
}

func (imp *InMemoryProvider) BatchDeleteDocuments(ctx context.Context, collection string, docIDs []string) error {
	for _, docID := range docIDs {
		imp.DeleteDocument(ctx, collection, docID)
	}
	return nil
}

func (imp *InMemoryProvider) GetStats(ctx context.Context) (map[string]interface{}, error) {
	totalDocs := 0
	totalVectors := 0

	for _, col := range imp.collections {
		totalDocs += len(col.Documents)
		totalVectors += len(col.Vectors)
	}

	return map[string]interface{}{
		"provider":        "InMemory",
		"connected":       imp.connected,
		"collections":     len(imp.collections),
		"total_documents": totalDocs,
		"total_vectors":   totalVectors,
	}, nil
}
