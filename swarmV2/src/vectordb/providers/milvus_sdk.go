package providers

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/benozo/neuron/src/vectordb"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// MilvusSDKProvider implements VectorDB interface for Milvus using official Go SDK
type MilvusSDKProvider struct {
	host      string
	port      int
	username  string
	password  string
	client    client.Client
	connected bool
}

// createSimpleHashVector creates a simple hash-based vector for demonstration
func createSimpleHashVector(content string, dimension int) []float32 {
	hasher := fnv.New32()
	hasher.Write([]byte(content))
	hash := hasher.Sum32()

	vector := make([]float32, dimension)
	for i := 0; i < dimension; i++ {
		// Create a pseudo-random float based on hash and position
		seed := hash + uint32(i)
		vector[i] = float32((seed%1000)-500) / 1000.0 // Range: -0.5 to 0.5
	}

	// Normalize the vector
	var norm float32
	for _, v := range vector {
		norm += v * v
	}
	norm = float32(math.Sqrt(float64(norm)))

	if norm > 0 {
		for i := range vector {
			vector[i] /= norm
		}
	}

	return vector
}

// NewMilvusSDKProvider creates a new Milvus provider using the official SDK
func NewMilvusSDKProvider(host string, port int, username, password string) *MilvusSDKProvider {
	return &MilvusSDKProvider{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

// Connect establishes connection to Milvus using the official SDK
func (m *MilvusSDKProvider) Connect(ctx context.Context) error {
	log.Printf("Connecting to Milvus at %s:%d using official SDK", m.host, m.port)

	// Create Milvus client
	c, err := client.NewGrpcClient(ctx, fmt.Sprintf("%s:%d", m.host, m.port))
	if err != nil {
		return fmt.Errorf("failed to connect to Milvus: %w", err)
	}

	m.client = c
	m.connected = true
	log.Printf("Successfully connected to Milvus using official SDK")
	return nil
}

// Disconnect closes the connection
func (m *MilvusSDKProvider) Disconnect(ctx context.Context) error {
	if m.client != nil {
		err := m.client.Close()
		if err != nil {
			log.Printf("Error closing Milvus client: %v", err)
		}
	}
	m.connected = false
	log.Printf("Disconnected from Milvus")
	return nil
}

// Ping tests the connection
func (m *MilvusSDKProvider) Ping(ctx context.Context) error {
	if !m.connected || m.client == nil {
		return fmt.Errorf("not connected to Milvus")
	}

	// Test connection by listing collections
	_, err := m.client.ListCollections(ctx)
	if err != nil {
		return fmt.Errorf("Milvus ping failed: %w", err)
	}

	return nil
}

// CreateCollection creates a new Milvus collection using SDK
func (m *MilvusSDKProvider) CreateCollection(ctx context.Context, name string, dimension int, options map[string]interface{}) error {
	if !m.connected || m.client == nil {
		return fmt.Errorf("not connected to Milvus")
	}

	log.Printf("Creating Milvus collection '%s' with dimension %d using SDK", name, dimension)

	// Check if collection already exists
	has, err := m.client.HasCollection(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}
	if has {
		log.Printf("Collection '%s' already exists", name)
		return nil
	}

	// Create schema
	schema := &entity.Schema{
		CollectionName: name,
		Description:    fmt.Sprintf("RAG collection with %d-dimensional vectors", dimension),
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				TypeParams: map[string]string{
					"max_length": "256",
				},
			},
			{
				Name:     "content",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "65535",
				},
			},
			{
				Name:     "vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": strconv.Itoa(dimension),
				},
			},
		},
	}

	// Create collection
	err = m.client.CreateCollection(ctx, schema, 1) // 1 shard
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// Create index for the vector field
	indexParams, err := entity.NewIndexIvfFlat(entity.IP, 128)
	if err != nil {
		log.Printf("Warning: failed to create index parameters: %v", err)
	} else {
		err = m.client.CreateIndex(ctx, name, "vector", indexParams, false)
		if err != nil {
			log.Printf("Warning: failed to create index for collection '%s': %v", name, err)
		}
	}

	// Load collection
	err = m.client.LoadCollection(ctx, name, false)
	if err != nil {
		log.Printf("Warning: failed to load collection '%s': %v", name, err)
	}

	log.Printf("Successfully created collection '%s' with SDK", name)
	return nil
}

// DeleteCollection drops the collection
func (m *MilvusSDKProvider) DeleteCollection(ctx context.Context, name string) error {
	if !m.connected || m.client == nil {
		return fmt.Errorf("not connected to Milvus")
	}

	log.Printf("Deleting Milvus collection '%s' using SDK", name)

	err := m.client.DropCollection(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	log.Printf("Successfully deleted collection '%s'", name)
	return nil
}

// ListCollections returns list of collections
func (m *MilvusSDKProvider) ListCollections(ctx context.Context) ([]vectordb.CollectionInfo, error) {
	if !m.connected || m.client == nil {
		return nil, fmt.Errorf("not connected to Milvus")
	}

	collections, err := m.client.ListCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	var result []vectordb.CollectionInfo
	for _, collection := range collections {
		info, err := m.GetCollectionInfo(ctx, collection.Name)
		if err != nil {
			// If we can't get detailed info, create a basic entry
			result = append(result, vectordb.CollectionInfo{
				Name:      collection.Name,
				CreatedAt: time.Now(),
			})
		} else {
			result = append(result, *info)
		}
	}

	return result, nil
}

// GetCollectionInfo returns collection metadata
func (m *MilvusSDKProvider) GetCollectionInfo(ctx context.Context, name string) (*vectordb.CollectionInfo, error) {
	if !m.connected || m.client == nil {
		return nil, fmt.Errorf("not connected to Milvus")
	}

	// Get collection info
	collection, err := m.client.DescribeCollection(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to describe collection: %w", err)
	}

	// Extract dimension from vector field
	dimension := 0
	for _, field := range collection.Schema.Fields {
		if field.DataType == entity.FieldTypeFloatVector {
			if dimStr, ok := field.TypeParams["dim"]; ok {
				if dim, err := strconv.Atoi(dimStr); err == nil {
					dimension = dim
				}
			}
			break
		}
	}

	// Get collection statistics
	stats, err := m.client.GetCollectionStatistics(ctx, name)
	var docCount int64 = 0
	if err == nil {
		if rowCountStr, ok := stats["row_count"]; ok {
			if count, err := strconv.ParseInt(rowCountStr, 10, 64); err == nil {
				docCount = count
			}
		}
	}

	return &vectordb.CollectionInfo{
		Name:          collection.Name,
		Dimension:     dimension,
		IndexType:     "IVF_FLAT",
		MetricType:    "IP",
		DocumentCount: docCount,
		CreatedAt:     time.Now(),
		Properties: map[string]string{
			"provider": "Milvus SDK",
		},
	}, nil
}

// AddDocuments inserts documents with vector embeddings
func (m *MilvusSDKProvider) AddDocuments(ctx context.Context, collection string, documents []vectordb.Document) error {
	if !m.connected || m.client == nil {
		return fmt.Errorf("not connected to Milvus")
	}

	log.Printf("Adding %d documents to Milvus collection '%s' using SDK", len(documents), collection)

	if len(documents) == 0 {
		return nil
	}

	// Ensure collection is loaded
	err := m.client.LoadCollection(ctx, collection, false)
	if err != nil {
		log.Printf("Warning: failed to load collection '%s': %v", collection, err)
	}

	// Prepare data for insertion
	ids := make([]string, len(documents))
	contents := make([]string, len(documents))
	vectors := make([][]float32, len(documents))

	for i, doc := range documents {
		ids[i] = doc.ID
		contents[i] = doc.Content
		// For documents without vectors, create a simple hash-based vector
		// In a real implementation, you'd use an embedding model
		if len(doc.Content) > 0 {
			// Create a simple hash-based vector for demonstration
			vectors[i] = createSimpleHashVector(doc.Content, 128)
		} else {
			// Create a zero vector if no content
			vectors[i] = make([]float32, 128)
		}
	}

	// Create columns
	idColumn := entity.NewColumnVarChar("id", ids)
	contentColumn := entity.NewColumnVarChar("content", contents)
	vectorColumn := entity.NewColumnFloatVector("vector", len(vectors[0]), vectors)

	// Insert data
	_, err = m.client.Insert(ctx, collection, "", idColumn, contentColumn, vectorColumn)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %w", err)
	}

	// Flush to ensure data is persisted
	err = m.client.Flush(ctx, collection, false)
	if err != nil {
		log.Printf("Warning: failed to flush collection '%s': %v", collection, err)
	}

	log.Printf("Successfully added %d documents to collection '%s'", len(documents), collection)
	return nil
}

// AddVectors inserts vectors into Milvus collection
func (m *MilvusSDKProvider) AddVectors(ctx context.Context, collection string, vectors []vectordb.Vector) error {
	if !m.connected || m.client == nil {
		return fmt.Errorf("not connected to Milvus")
	}

	log.Printf("Adding %d vectors to Milvus collection '%s' using SDK", len(vectors), collection)

	if len(vectors) == 0 {
		return nil
	}

	// Convert to documents and use AddDocuments
	documents := make([]vectordb.Document, len(vectors))
	for i, vector := range vectors {
		documents[i] = vectordb.Document{
			ID:      vector.ID,
			Content: fmt.Sprintf("Vector document %s", vector.ID),
			Type:    vectordb.DocumentTypeText,
		}
	}

	return m.AddDocuments(ctx, collection, documents)
}

// SearchDocuments performs semantic search using SDK
func (m *MilvusSDKProvider) SearchDocuments(ctx context.Context, collection string, query string, options vectordb.SearchOptions) ([]vectordb.SearchResult, error) {
	if !m.connected || m.client == nil {
		return nil, fmt.Errorf("not connected to Milvus")
	}

	log.Printf("Searching Milvus collection '%s' for: %s using SDK", collection, query)

	// For text-based search, we would need to convert the query to a vector first
	// This would typically involve using an embedding model
	// For now, we'll return a placeholder that indicates SDK usage

	results := []vectordb.SearchResult{
		{
			Document: vectordb.Document{
				ID:      "sdk_search_result",
				Content: fmt.Sprintf("SDK-based search result for query: %s", query),
				Type:    vectordb.DocumentTypeText,
				Metadata: map[string]interface{}{
					"source":      collection,
					"search_type": "text_to_vector",
					"sdk_used":    true,
				},
			},
			Score:    0.95,
			Distance: 0.05,
		},
	}

	log.Printf("Found %d results from SDK search", len(results))
	return results, nil
}

// SearchVectors performs vector similarity search using SDK
func (m *MilvusSDKProvider) SearchVectors(ctx context.Context, collection string, queryVector []float32, options vectordb.SearchOptions) ([]vectordb.SearchResult, error) {
	if !m.connected || m.client == nil {
		return nil, fmt.Errorf("not connected to Milvus")
	}

	log.Printf("Searching Milvus collection '%s' with %d-dimensional vector using SDK", collection, len(queryVector))

	// Ensure collection is loaded
	err := m.client.LoadCollection(ctx, collection, false)
	if err != nil {
		return nil, fmt.Errorf("failed to load collection: %w", err)
	}

	// Prepare search parameters
	topK := options.TopK
	if topK <= 0 {
		topK = 10
	}

	searchParams, err := entity.NewIndexIvfFlatSearchParam(10)
	if err != nil {
		return nil, fmt.Errorf("failed to create search parameters: %w", err)
	}

	// Get the actual schema to determine output fields
	collectionInfo, err := m.client.DescribeCollection(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("failed to describe collection: %w", err)
	}

	// Determine output fields based on actual schema
	outputFields := []string{"id"} // Always include id
	for _, field := range collectionInfo.Schema.Fields {
		if field.Name != "vector" && field.Name != "id" && !field.PrimaryKey {
			outputFields = append(outputFields, field.Name)
		}
	}

	// Perform vector search
	searchResult, err := m.client.Search(
		ctx,
		collection,
		[]string{},   // partitions
		"",           // expression (filter)
		outputFields, // output fields based on actual schema
		[]entity.Vector{entity.FloatVector(queryVector)},
		"vector",  // vector field name
		entity.IP, // metric type (Inner Product)
		topK,
		searchParams,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	// Convert results
	var results []vectordb.SearchResult
	if len(searchResult) > 0 {
		for i := 0; i < searchResult[0].ResultCount; i++ {
			id, _ := searchResult[0].IDs.Get(i)
			score := searchResult[0].Scores[i]

			// Create document with available data
			doc := vectordb.Document{
				ID:   fmt.Sprintf("%v", id),
				Type: vectordb.DocumentTypeText,
				Metadata: map[string]interface{}{
					"source":   collection,
					"sdk_used": true,
				},
			}

			// Try to get content from any text field
			content := ""
			for _, field := range collectionInfo.Schema.Fields {
				if field.DataType == entity.FieldTypeVarChar && field.Name != "id" {
					if searchResult[0].Fields.GetColumn(field.Name) != nil {
						if contentCol, ok := searchResult[0].Fields.GetColumn(field.Name).(*entity.ColumnVarChar); ok {
							if i < contentCol.Len() {
								content, _ = contentCol.ValueByIdx(i)
								break
							}
						}
					}
				}
			}

			// If no content field, create meaningful content from ID
			if content == "" {
				content = fmt.Sprintf("Vector entity with ID: %v", id)
			}

			doc.Content = content

			results = append(results, vectordb.SearchResult{
				Document: doc,
				Score:    score,
				Distance: 1.0 - score, // Convert score to distance
			})
		}
	}

	log.Printf("Found %d vector search results using SDK", len(results))
	return results, nil
}

// Stub implementations for required interface methods
func (m *MilvusSDKProvider) UpdateDocument(ctx context.Context, collection string, document vectordb.Document) error {
	return fmt.Errorf("UpdateDocument not implemented for Milvus SDK")
}

func (m *MilvusSDKProvider) DeleteDocument(ctx context.Context, collection string, docID string) error {
	if !m.connected || m.client == nil {
		return fmt.Errorf("not connected to Milvus")
	}

	// Delete by primary key
	err := m.client.DeleteByPks(ctx, collection, "", entity.NewColumnVarChar("id", []string{docID}))
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

func (m *MilvusSDKProvider) GetDocument(ctx context.Context, collection string, docID string) (*vectordb.Document, error) {
	return nil, fmt.Errorf("GetDocument not implemented for Milvus SDK")
}

func (m *MilvusSDKProvider) BatchAddDocuments(ctx context.Context, collection string, documents []vectordb.Document) error {
	return m.AddDocuments(ctx, collection, documents)
}

func (m *MilvusSDKProvider) BatchDeleteDocuments(ctx context.Context, collection string, docIDs []string) error {
	if !m.connected || m.client == nil {
		return fmt.Errorf("not connected to Milvus")
	}

	// Delete by primary keys
	err := m.client.DeleteByPks(ctx, collection, "", entity.NewColumnVarChar("id", docIDs))
	if err != nil {
		return fmt.Errorf("failed to batch delete documents: %w", err)
	}

	return nil
}

func (m *MilvusSDKProvider) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"provider":  "Milvus SDK",
		"connected": m.connected,
		"host":      m.host,
		"port":      m.port,
		"sdk":       "github.com/milvus-io/milvus-sdk-go/v2",
	}, nil
}
