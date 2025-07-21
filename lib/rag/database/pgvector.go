package database

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"

	"github.com/benozo/conduit/lib/rag"
)

// PgVectorDB implements VectorDB interface using PostgreSQL with pgvector
type PgVectorDB struct {
	db     *sql.DB
	config rag.DatabaseConfig
}

// NewPgVectorDB creates a new PostgreSQL vector database connection
func NewPgVectorDB(config rag.DatabaseConfig) (*PgVectorDB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	pgvector := &PgVectorDB{
		db:     db,
		config: config,
	}

	// Initialize database schema
	if err := pgvector.initializeSchema(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return pgvector, nil
}

// initializeSchema creates necessary tables and indexes with dynamic vector dimensions
func (p *PgVectorDB) initializeSchema(ctx context.Context) error {
	// First create documents table and extensions
	documentsSQL := `
		-- Enable pgvector extension
		CREATE EXTENSION IF NOT EXISTS vector;
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		-- Documents table to store document metadata
		CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			source_path TEXT,
			content_type TEXT DEFAULT 'text/plain',
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			
			-- Indexes for performance
			CONSTRAINT documents_title_not_empty CHECK (length(title) > 0),
			CONSTRAINT documents_content_not_empty CHECK (length(content) > 0)
		);`

	if _, err := p.db.ExecContext(ctx, documentsSQL); err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	// Create document_chunks table without vector column first
	chunksSQL := `
		CREATE TABLE IF NOT EXISTS document_chunks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
			chunk_index INTEGER NOT NULL,
			content TEXT NOT NULL,
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			
			-- Constraints
			CONSTRAINT chunks_content_not_empty CHECK (length(content) > 0),
			CONSTRAINT chunks_index_positive CHECK (chunk_index >= 0),
			UNIQUE(document_id, chunk_index)
		);`

	if _, err := p.db.ExecContext(ctx, chunksSQL); err != nil {
		return fmt.Errorf("failed to create chunks table: %w", err)
	}

	// Check if embedding column exists, if not add it (this allows for dynamic sizing)
	checkColumnSQL := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = 'document_chunks' AND column_name = 'embedding';`

	var columnName string
	err := p.db.QueryRowContext(ctx, checkColumnSQL).Scan(&columnName)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check embedding column: %w", err)
	}

	// If embedding column doesn't exist, we'll add it when first vector is inserted
	if columnName == "" {
		log.Printf("Embedding column will be created dynamically on first vector insert")
	}

	// Create indexes and triggers
	indexesSQL := `
		-- Indexes for performance
		CREATE INDEX IF NOT EXISTS idx_documents_created_at ON documents(created_at);
		CREATE INDEX IF NOT EXISTS idx_documents_content_type ON documents(content_type);
		CREATE INDEX IF NOT EXISTS idx_chunks_document_id ON document_chunks(document_id);
		CREATE INDEX IF NOT EXISTS idx_chunks_created_at ON document_chunks(created_at);

		-- Full-text search indexes
		CREATE INDEX IF NOT EXISTS idx_documents_content_fts ON documents USING gin(to_tsvector('english', content));
		CREATE INDEX IF NOT EXISTS idx_chunks_content_fts ON document_chunks USING gin(to_tsvector('english', content));

		-- Function to update timestamps
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ language 'plpgsql';

		-- Triggers for automatic timestamp updates (skip if exists)
		DROP TRIGGER IF EXISTS update_documents_updated_at ON documents;
		CREATE TRIGGER update_documents_updated_at
			BEFORE UPDATE ON documents
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`

	if _, err := p.db.ExecContext(ctx, indexesSQL); err != nil {
		return fmt.Errorf("failed to create indexes and triggers: %w", err)
	}

	log.Printf("Database schema initialized successfully")
	return nil
}

// StoreDocument stores a document in the database
func (p *PgVectorDB) StoreDocument(ctx context.Context, doc rag.Document) error {
	query := `
		INSERT INTO documents (id, title, content, source_path, content_type, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			source_path = EXCLUDED.source_path,
			content_type = EXCLUDED.content_type,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at`

	metadataJSON, err := json.Marshal(doc.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query,
		doc.ID, doc.Title, doc.Content, doc.SourcePath, doc.ContentType,
		metadataJSON, doc.CreatedAt, doc.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	return nil
}

// GetDocument retrieves a document by ID
func (p *PgVectorDB) GetDocument(ctx context.Context, id string) (*rag.Document, error) {
	query := `SELECT id, title, content, source_path, content_type, metadata, created_at, updated_at 
			  FROM documents WHERE id = $1`

	var doc rag.Document
	var metadataJSON []byte

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&doc.ID, &doc.Title, &doc.Content, &doc.SourcePath, &doc.ContentType,
		&metadataJSON, &doc.CreatedAt, &doc.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if err := jsonUnmarshal(metadataJSON, &doc.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &doc, nil
}

// DeleteDocument removes a document and its chunks
func (p *PgVectorDB) DeleteDocument(ctx context.Context, id string) error {
	query := `DELETE FROM documents WHERE id = $1`

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found: %s", id)
	}

	return nil
}

// ListDocuments returns a paginated list of documents
func (p *PgVectorDB) ListDocuments(ctx context.Context, limit, offset int) ([]rag.Document, error) {
	query := `SELECT id, title, content, source_path, content_type, metadata, created_at, updated_at 
			  FROM documents ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	var documents []rag.Document
	for rows.Next() {
		var doc rag.Document
		var metadataJSON []byte

		err := rows.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.SourcePath, &doc.ContentType,
			&metadataJSON, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		if err := jsonUnmarshal(metadataJSON, &doc.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		documents = append(documents, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}

	return documents, nil
}

// StoreChunks stores document chunks with embeddings
func (p *PgVectorDB) StoreChunks(ctx context.Context, chunks []rag.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	// Check if embedding column exists and create it if needed
	if err := p.ensureEmbeddingColumn(ctx, chunks[0]); err != nil {
		return fmt.Errorf("failed to ensure embedding column: %w", err)
	}

	// Prepare batch insert
	valueStrings := make([]string, 0, len(chunks))
	valueArgs := make([]interface{}, 0, len(chunks)*6)

	for i, chunk := range chunks {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))

		metadataJSON, err := jsonMarshal(chunk.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal chunk metadata: %w", err)
		}

		valueArgs = append(valueArgs,
			chunk.DocumentID, chunk.Index, chunk.Content,
			pgvector.NewVector(chunk.Embedding), metadataJSON, chunk.CreatedAt)
	}

	query := fmt.Sprintf(`
		INSERT INTO document_chunks (document_id, chunk_index, content, embedding, metadata, created_at)
		VALUES %s
		ON CONFLICT (document_id, chunk_index) DO UPDATE SET
			content = EXCLUDED.content,
			embedding = EXCLUDED.embedding,
			metadata = EXCLUDED.metadata`,
		strings.Join(valueStrings, ","))

	_, err := p.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to store chunks: %w", err)
	}

	return nil
}

// GetChunk retrieves a document chunk by ID
func (p *PgVectorDB) GetChunk(ctx context.Context, id string) (*rag.DocumentChunk, error) {
	query := `SELECT id, document_id, chunk_index, content, embedding, metadata, created_at 
			  FROM document_chunks WHERE id = $1`

	var chunk rag.DocumentChunk
	var embedding pgvector.Vector
	var metadataJSON []byte

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&chunk.ID, &chunk.DocumentID, &chunk.Index, &chunk.Content,
		&embedding, &metadataJSON, &chunk.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chunk not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get chunk: %w", err)
	}

	chunk.Embedding = embedding.Slice()

	if err := jsonUnmarshal(metadataJSON, &chunk.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &chunk, nil
}

// GetDocumentChunks retrieves all chunks for a document
func (p *PgVectorDB) GetDocumentChunks(ctx context.Context, documentID string) ([]rag.DocumentChunk, error) {
	query := `SELECT id, document_id, chunk_index, content, embedding, metadata, created_at 
			  FROM document_chunks WHERE document_id = $1 ORDER BY chunk_index`

	rows, err := p.db.QueryContext(ctx, query, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document chunks: %w", err)
	}
	defer rows.Close()

	var chunks []rag.DocumentChunk
	for rows.Next() {
		var chunk rag.DocumentChunk
		var embedding pgvector.Vector
		var metadataJSON []byte

		err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.Index, &chunk.Content,
			&embedding, &metadataJSON, &chunk.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chunk: %w", err)
		}

		chunk.Embedding = embedding.Slice()

		if err := jsonUnmarshal(metadataJSON, &chunk.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		chunks = append(chunks, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chunks: %w", err)
	}

	return chunks, nil
}

// SearchSimilar performs vector similarity search
func (p *PgVectorDB) SearchSimilar(ctx context.Context, embedding []float32, limit int, filters map[string]interface{}) ([]rag.SearchResult, error) {
	// Build the query with optional filters
	query := `
		SELECT 
			c.id, c.document_id, c.chunk_index, c.content, c.embedding, c.metadata, c.created_at,
			d.id, d.title, d.content, d.source_path, d.content_type, d.metadata, d.created_at, d.updated_at,
			1 - (c.embedding <=> $1) as similarity_score
		FROM document_chunks c
		JOIN documents d ON c.document_id = d.id`

	args := []interface{}{pgvector.NewVector(embedding)}
	whereConditions := []string{}
	argIndex := 2

	// Add metadata filters
	for key, value := range filters {
		switch key {
		case "document_id":
			whereConditions = append(whereConditions, fmt.Sprintf("c.document_id = $%d", argIndex))
			args = append(args, value)
			argIndex++
		case "content_type":
			whereConditions = append(whereConditions, fmt.Sprintf("d.content_type = $%d", argIndex))
			args = append(args, value)
			argIndex++
		default:
			// Handle JSONB metadata filters
			whereConditions = append(whereConditions, fmt.Sprintf("(c.metadata ->> '%s' = $%d OR d.metadata ->> '%s' = $%d)", key, argIndex, key, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	query += " ORDER BY c.embedding <=> $1 LIMIT $" + fmt.Sprintf("%d", argIndex)
	args = append(args, limit)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar: %w", err)
	}
	defer rows.Close()

	var results []rag.SearchResult
	for rows.Next() {
		var result rag.SearchResult
		var chunkEmbedding pgvector.Vector
		var chunkMetadataJSON, docMetadataJSON []byte

		err := rows.Scan(
			&result.Chunk.ID, &result.Chunk.DocumentID, &result.Chunk.Index,
			&result.Chunk.Content, &chunkEmbedding, &chunkMetadataJSON, &result.Chunk.CreatedAt,
			&result.Document.ID, &result.Document.Title, &result.Document.Content,
			&result.Document.SourcePath, &result.Document.ContentType, &docMetadataJSON,
			&result.Document.CreatedAt, &result.Document.UpdatedAt,
			&result.Score)

		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}

		result.Chunk.Embedding = chunkEmbedding.Slice()

		if err := jsonUnmarshal(chunkMetadataJSON, &result.Chunk.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chunk metadata: %w", err)
		}

		if err := jsonUnmarshal(docMetadataJSON, &result.Document.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal document metadata: %w", err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return results, nil
}

// SearchByText performs text-based search (placeholder for hybrid search)
func (p *PgVectorDB) SearchByText(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]rag.SearchResult, error) {
	// For now, this is a simple text search. In a full implementation,
	// this would involve generating embeddings for the query and calling SearchSimilar
	// This is a placeholder that searches by content similarity

	sqlQuery := `
		SELECT 
			c.id, c.document_id, c.chunk_index, c.content, c.embedding, c.metadata, c.created_at,
			d.id, d.title, d.content, d.source_path, d.content_type, d.metadata, d.created_at, d.updated_at,
			similarity(c.content, $1) as text_score
		FROM document_chunks c
		JOIN documents d ON c.document_id = d.id
		WHERE c.content ILIKE $2
		ORDER BY similarity(c.content, $1) DESC
		LIMIT $3`

	likePattern := "%" + query + "%"
	rows, err := p.db.QueryContext(ctx, sqlQuery, query, likePattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search by text: %w", err)
	}
	defer rows.Close()

	var results []rag.SearchResult
	for rows.Next() {
		var result rag.SearchResult
		var chunkEmbedding pgvector.Vector
		var chunkMetadataJSON, docMetadataJSON []byte

		err := rows.Scan(
			&result.Chunk.ID, &result.Chunk.DocumentID, &result.Chunk.Index,
			&result.Chunk.Content, &chunkEmbedding, &chunkMetadataJSON, &result.Chunk.CreatedAt,
			&result.Document.ID, &result.Document.Title, &result.Document.Content,
			&result.Document.SourcePath, &result.Document.ContentType, &docMetadataJSON,
			&result.Document.CreatedAt, &result.Document.UpdatedAt,
			&result.Score)

		if err != nil {
			return nil, fmt.Errorf("failed to scan text search result: %w", err)
		}

		result.Chunk.Embedding = chunkEmbedding.Slice()

		if err := jsonUnmarshal(chunkMetadataJSON, &result.Chunk.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chunk metadata: %w", err)
		}

		if err := jsonUnmarshal(docMetadataJSON, &result.Document.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal document metadata: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

// CreateIndex creates a vector index
func (p *PgVectorDB) CreateIndex(ctx context.Context, indexType string) error {
	var query string

	switch indexType {
	case "ivfflat_cosine":
		query = `CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chunks_embedding_cosine_new 
				 ON document_chunks USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100)`
	case "ivfflat_l2":
		query = `CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chunks_embedding_l2_new 
				 ON document_chunks USING ivfflat (embedding vector_l2_ops) WITH (lists = 100)`
	case "hnsw_cosine":
		query = `CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chunks_embedding_hnsw_cosine 
				 ON document_chunks USING hnsw (embedding vector_cosine_ops)`
	default:
		return fmt.Errorf("unsupported index type: %s", indexType)
	}

	_, err := p.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

// DropIndex drops a vector index
func (p *PgVectorDB) DropIndex(ctx context.Context, indexType string) error {
	var indexName string

	switch indexType {
	case "ivfflat_cosine":
		indexName = "idx_chunks_embedding_cosine"
	case "ivfflat_l2":
		indexName = "idx_chunks_embedding_l2"
	case "hnsw_cosine":
		indexName = "idx_chunks_embedding_hnsw_cosine"
	default:
		return fmt.Errorf("unsupported index type: %s", indexType)
	}

	query := fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)
	_, err := p.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop index: %w", err)
	}

	return nil
}

// GetStats returns database statistics
func (p *PgVectorDB) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Document count
	var docCount int
	err := p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM documents").Scan(&docCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get document count: %w", err)
	}
	stats["document_count"] = docCount

	// Chunk count
	var chunkCount int
	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM document_chunks").Scan(&chunkCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunk count: %w", err)
	}
	stats["chunk_count"] = chunkCount

	// Database size
	var dbSize int64
	err = p.db.QueryRowContext(ctx, "SELECT pg_database_size(current_database())").Scan(&dbSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get database size: %w", err)
	}
	stats["database_size_bytes"] = dbSize

	// Average chunk length
	var avgChunkLength float64
	err = p.db.QueryRowContext(ctx, "SELECT AVG(LENGTH(content)) FROM document_chunks").Scan(&avgChunkLength)
	if err != nil {
		return nil, fmt.Errorf("failed to get average chunk length: %w", err)
	}
	stats["avg_chunk_length"] = avgChunkLength

	return stats, nil
}

// Ping checks database connectivity
func (p *PgVectorDB) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// Close closes the database connection
func (p *PgVectorDB) Close() error {
	return p.db.Close()
}

// Helper functions for JSON marshaling/unmarshaling
func jsonMarshal(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// GenerateID generates a new UUID for documents or chunks
func GenerateID() string {
	return uuid.New().String()
}

// ensureEmbeddingColumn creates the embedding column if it doesn't exist
func (p *PgVectorDB) ensureEmbeddingColumn(ctx context.Context, sampleChunk rag.DocumentChunk) error {
	// Check if embedding column exists
	checkColumnSQL := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = 'document_chunks' AND column_name = 'embedding';`

	var columnName string
	err := p.db.QueryRowContext(ctx, checkColumnSQL).Scan(&columnName)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check embedding column: %w", err)
	}

	// If column doesn't exist, create it with the correct dimensions
	if columnName == "" {
		dimensions := len(sampleChunk.Embedding)
		if dimensions == 0 {
			return fmt.Errorf("cannot create embedding column with zero dimensions")
		}

		alterSQL := fmt.Sprintf(`
			ALTER TABLE document_chunks 
			ADD COLUMN embedding vector(%d);`, dimensions)

		if _, err := p.db.ExecContext(ctx, alterSQL); err != nil {
			return fmt.Errorf("failed to add embedding column: %w", err)
		}

		// Create the vector search index
		indexSQL := `
			CREATE INDEX IF NOT EXISTS idx_chunks_embedding_cosine 
			ON document_chunks USING ivfflat (embedding vector_cosine_ops) 
			WITH (lists = 100);`

		if _, err := p.db.ExecContext(ctx, indexSQL); err != nil {
			log.Printf("Warning: failed to create vector index: %v", err)
		}

		log.Printf("Created embedding column with %d dimensions", dimensions)
	}

	return nil
}
