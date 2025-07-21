-- Database schema for RAG with pgvector
-- Requires PostgreSQL 15+ with pgvector extension

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
);

-- Document chunks table with vector embeddings
CREATE TABLE IF NOT EXISTS document_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector, -- Dynamic dimension vector (auto-sized based on first insert)
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chunks_content_not_empty CHECK (length(content) > 0),
    CONSTRAINT chunks_index_positive CHECK (chunk_index >= 0),
    UNIQUE(document_id, chunk_index)
);

-- Indexes for optimal vector search performance
-- Using IVFFlat for good balance of speed and recall
CREATE INDEX IF NOT EXISTS document_chunks_embedding_cosine_idx 
    ON document_chunks USING ivfflat (embedding vector_cosine_ops) 
    WITH (lists = 100);

-- Additional indexes for metadata filtering and general queries
CREATE INDEX IF NOT EXISTS document_chunks_document_id_idx ON document_chunks(document_id);
CREATE INDEX IF NOT EXISTS document_chunks_metadata_idx ON document_chunks USING gin(metadata);
CREATE INDEX IF NOT EXISTS documents_metadata_idx ON documents USING gin(metadata);
CREATE INDEX IF NOT EXISTS documents_source_path_idx ON documents(source_path);
CREATE INDEX IF NOT EXISTS documents_content_type_idx ON documents(content_type);
CREATE INDEX IF NOT EXISTS documents_created_at_idx ON documents(created_at DESC);

-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at for documents
DROP TRIGGER IF EXISTS update_documents_updated_at ON documents;
CREATE TRIGGER update_documents_updated_at 
    BEFORE UPDATE ON documents 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to calculate vector similarity (for testing/debugging)
CREATE OR REPLACE FUNCTION vector_similarity(embedding1 vector, embedding2 vector)
RETURNS FLOAT AS $$
BEGIN
    RETURN 1 - (embedding1 <=> embedding2);
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- View for document statistics
CREATE OR REPLACE VIEW document_stats AS
SELECT 
    d.id,
    d.title,
    d.source_path,
    d.content_type,
    d.created_at,
    COUNT(dc.id) as chunk_count,
    AVG(LENGTH(dc.content)) as avg_chunk_length,
    SUM(LENGTH(dc.content)) as total_content_length
FROM documents d
LEFT JOIN document_chunks dc ON d.id = dc.document_id
GROUP BY d.id, d.title, d.source_path, d.content_type, d.created_at;

-- Index for the view
CREATE INDEX IF NOT EXISTS document_chunks_content_length_idx ON document_chunks((LENGTH(content)));
