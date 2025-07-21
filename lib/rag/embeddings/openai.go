package embeddings

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/benozo/conduit/lib/rag"
)

// OpenAIEmbeddings implements EmbeddingProvider for OpenAI
type OpenAIEmbeddings struct {
	client     *openai.Client
	model      openai.EmbeddingModel
	modelStr   string
	dimensions int
	timeout    time.Duration
}

// NewOpenAIEmbeddings creates a new OpenAI embeddings provider
func NewOpenAIEmbeddings(apiKey, model string, dimensions int, timeout time.Duration) *OpenAIEmbeddings {
	var embeddingModel openai.EmbeddingModel
	if err := embeddingModel.UnmarshalText([]byte(model)); err != nil {
		// If unmarshal fails, use the default AdaEmbeddingV2
		embeddingModel = openai.AdaEmbeddingV2
	}

	return &OpenAIEmbeddings{
		client:     openai.NewClient(apiKey),
		model:      embeddingModel,
		modelStr:   model,
		dimensions: dimensions,
		timeout:    timeout,
	}
}

// Embed generates embedding for a single text
func (o *OpenAIEmbeddings) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Create context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: o.model,
	}

	resp, err := o.client.CreateEmbeddings(ctxWithTimeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}

// EmbedBatch generates embeddings for multiple texts
func (o *OpenAIEmbeddings) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts cannot be empty")
	}

	// Filter out empty texts
	validTexts := make([]string, 0, len(texts))
	indexMap := make(map[int]int) // maps result index to original index

	for i, text := range texts {
		if text != "" {
			indexMap[len(validTexts)] = i
			validTexts = append(validTexts, text)
		}
	}

	if len(validTexts) == 0 {
		return nil, fmt.Errorf("no valid texts provided")
	}

	// Create context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	req := openai.EmbeddingRequest{
		Input: validTexts,
		Model: o.model,
	}

	resp, err := o.client.CreateEmbeddings(ctxWithTimeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch embeddings: %w", err)
	}

	if len(resp.Data) != len(validTexts) {
		return nil, fmt.Errorf("mismatch between input texts (%d) and embeddings (%d)",
			len(validTexts), len(resp.Data))
	}

	// Create result array with same length as input, filling in embeddings for valid texts
	results := make([][]float32, len(texts))
	for i, embedding := range resp.Data {
		originalIndex := indexMap[i]
		results[originalIndex] = embedding.Embedding
	}

	return results, nil
}

// GetDimensions returns the embedding dimensions
func (o *OpenAIEmbeddings) GetDimensions() int {
	return o.dimensions
}

// GetModel returns the model name
func (o *OpenAIEmbeddings) GetModel() string {
	return o.modelStr
}

// GetProvider returns the provider name
func (o *OpenAIEmbeddings) GetProvider() string {
	return "openai"
}

// Ping checks if the OpenAI API is accessible
func (o *OpenAIEmbeddings) Ping(ctx context.Context) error {
	// Test with a simple embedding request
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := openai.EmbeddingRequest{
		Input: []string{"test"},
		Model: o.model,
	}

	_, err := o.client.CreateEmbeddings(ctxWithTimeout, req)
	if err != nil {
		return fmt.Errorf("OpenAI API ping failed: %w", err)
	}

	return nil
}

// EmbeddingCache provides caching for embeddings
type EmbeddingCache struct {
	provider rag.EmbeddingProvider
	cache    map[string][]float32
	maxSize  int
}

// NewEmbeddingCache creates a new embedding cache
func NewEmbeddingCache(provider rag.EmbeddingProvider, maxSize int) *EmbeddingCache {
	return &EmbeddingCache{
		provider: provider,
		cache:    make(map[string][]float32),
		maxSize:  maxSize,
	}
}

// Embed with caching
func (c *EmbeddingCache) Embed(ctx context.Context, text string) ([]float32, error) {
	// Check cache first
	if embedding, exists := c.cache[text]; exists {
		return embedding, nil
	}

	// Generate embedding
	embedding, err := c.provider.Embed(ctx, text)
	if err != nil {
		return nil, err
	}

	// Cache the result (with simple size management)
	if len(c.cache) >= c.maxSize {
		// Simple eviction: remove first item
		for key := range c.cache {
			delete(c.cache, key)
			break
		}
	}
	c.cache[text] = embedding

	return embedding, nil
}

// EmbedBatch with caching
func (c *EmbeddingCache) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	results := make([][]float32, len(texts))
	uncachedTexts := make([]string, 0)
	uncachedIndices := make([]int, 0)

	// Check cache for each text
	for i, text := range texts {
		if embedding, exists := c.cache[text]; exists {
			results[i] = embedding
		} else {
			uncachedTexts = append(uncachedTexts, text)
			uncachedIndices = append(uncachedIndices, i)
		}
	}

	// Generate embeddings for uncached texts
	if len(uncachedTexts) > 0 {
		embeddings, err := c.provider.EmbedBatch(ctx, uncachedTexts)
		if err != nil {
			return nil, err
		}

		// Cache and assign results
		for i, embedding := range embeddings {
			if embedding != nil {
				text := uncachedTexts[i]
				originalIndex := uncachedIndices[i]

				// Cache management
				if len(c.cache) >= c.maxSize {
					for key := range c.cache {
						delete(c.cache, key)
						break
					}
				}
				c.cache[text] = embedding
				results[originalIndex] = embedding
			}
		}
	}

	return results, nil
}

// GetDimensions returns the embedding dimensions
func (c *EmbeddingCache) GetDimensions() int {
	return c.provider.GetDimensions()
}

// GetModel returns the model name
func (c *EmbeddingCache) GetModel() string {
	return c.provider.GetModel()
}

// GetProvider returns the provider name
func (c *EmbeddingCache) GetProvider() string {
	return c.provider.GetProvider() + "_cached"
}

// Ping checks the underlying provider
func (c *EmbeddingCache) Ping(ctx context.Context) error {
	return c.provider.Ping(ctx)
}

// ClearCache clears the embedding cache
func (c *EmbeddingCache) ClearCache() {
	c.cache = make(map[string][]float32)
}

// GetCacheStats returns cache statistics
func (c *EmbeddingCache) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"size":     len(c.cache),
		"max_size": c.maxSize,
		"usage":    float64(len(c.cache)) / float64(c.maxSize),
	}
}
