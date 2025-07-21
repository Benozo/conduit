package rag

import (
	"fmt"
	"time"
)

// RAGConfig holds configuration for RAG system
type RAGConfig struct {
	// Database configuration
	Database DatabaseConfig `json:"database"`

	// Embedding configuration
	Embeddings EmbeddingConfig `json:"embeddings"`

	// Chunking configuration
	Chunking ChunkingConfig `json:"chunking"`

	// Search configuration
	Search SearchConfig `json:"search"`
}

// DatabaseConfig for PostgreSQL with pgvector
type DatabaseConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Name         string        `json:"name"`
	User         string        `json:"user"`
	Password     string        `json:"password"`
	SSLMode      string        `json:"ssl_mode"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

// EmbeddingConfig for embedding providers
type EmbeddingConfig struct {
	Provider   string        `json:"provider"` // "openai", "ollama"
	APIKey     string        `json:"api_key"`  // For OpenAI
	Host       string        `json:"host"`     // For Ollama (e.g., "192.168.10.10")
	Model      string        `json:"model"`    // "text-embedding-ada-002" or "nomic-embed-text:latest"
	Dimensions int           `json:"dimensions"`
	BatchSize  int           `json:"batch_size"`
	Timeout    time.Duration `json:"timeout"`
}

// ChunkingConfig for text chunking
type ChunkingConfig struct {
	Size     int    `json:"size"`     // Default chunk size in characters
	Overlap  int    `json:"overlap"`  // Overlap between chunks
	Strategy string `json:"strategy"` // "fixed", "semantic", "paragraph"
}

// SearchConfig for vector search
type SearchConfig struct {
	DefaultLimit int     `json:"default_limit"`
	MaxLimit     int     `json:"max_limit"`
	Threshold    float64 `json:"threshold"` // Similarity threshold
	Algorithm    string  `json:"algorithm"` // "cosine", "l2", "inner_product"
}

// DefaultRAGConfig returns default configuration with OpenAI
func DefaultRAGConfig() *RAGConfig {
	return &RAGConfig{
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         5432,
			Name:         "conduit_rag",
			User:         "conduit",
			Password:     "conduit_password",
			SSLMode:      "disable",
			MaxOpenConns: 25,
			MaxIdleConns: 5,
			MaxLifetime:  5 * time.Minute,
		},
		Embeddings: EmbeddingConfig{
			Provider:   "openai",
			Model:      "text-embedding-ada-002",
			Dimensions: 1536,
			BatchSize:  100,
			Timeout:    30 * time.Second,
		},
		Chunking: ChunkingConfig{
			Size:     1000,
			Overlap:  200,
			Strategy: "fixed",
		},
		Search: SearchConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
			Threshold:    0.7,
			Algorithm:    "cosine",
		},
	}
}

// DefaultOllamaRAGConfig returns default configuration with Ollama
func DefaultOllamaRAGConfig() *RAGConfig {
	config := DefaultRAGConfig()
	config.Embeddings = EmbeddingConfig{
		Provider:   "ollama",
		Host:       "192.168.10.10",
		Model:      "nomic-embed-text:latest",
		Dimensions: 768,              // Default for nomic-embed-text
		BatchSize:  10,               // Lower batch size for Ollama
		Timeout:    60 * time.Second, // Higher timeout for local processing
	}
	return config
}

// LoadConfigFromEnv loads configuration from environment variables
func (c *RAGConfig) LoadFromEnv() {
	// Implementation to load from environment variables
	// This would use os.Getenv() to override defaults
}

// Validate checks if configuration is valid
func (c *RAGConfig) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Embeddings.Provider == "" {
		return fmt.Errorf("embedding provider is required")
	}
	if c.Embeddings.Dimensions <= 0 {
		return fmt.Errorf("embedding dimensions must be positive")
	}
	if c.Chunking.Size <= 0 {
		return fmt.Errorf("chunk size must be positive")
	}
	if c.Search.DefaultLimit <= 0 {
		return fmt.Errorf("search default limit must be positive")
	}
	return nil
}
