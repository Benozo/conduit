package vectordb

import (
	"context"
	"fmt"
)

// SimpleDocumentProcessor implements basic document processing
type SimpleDocumentProcessor struct {
	supportedTypes []DocumentType
}

// NewSimpleDocumentProcessor creates a new document processor
func NewSimpleDocumentProcessor() *SimpleDocumentProcessor {
	return &SimpleDocumentProcessor{
		supportedTypes: []DocumentType{
			DocumentTypeText,
			DocumentTypeJSON,
			DocumentTypePDF, // Basic support
		},
	}
}

// ProcessDocument processes a single document into chunks
func (sdp *SimpleDocumentProcessor) ProcessDocument(ctx context.Context, document Document) ([]Document, error) {
	// For text documents, chunk the content
	if document.Type == DocumentTypeText {
		chunks, err := sdp.ChunkText(document.Content, 512, 50) // 512 chars with 50 char overlap
		if err != nil {
			return nil, fmt.Errorf("failed to chunk text: %w", err)
		}

		var processedDocs []Document
		for i, chunk := range chunks {
			chunkDoc := Document{
				ID:      fmt.Sprintf("%s_chunk_%d", document.ID, i),
				Content: chunk,
				Type:    DocumentTypeText,
				Metadata: map[string]interface{}{
					"parent_id":   document.ID,
					"chunk_index": i,
					"chunk_size":  len(chunk),
				},
			}

			// Copy original metadata
			for k, v := range document.Metadata {
				chunkDoc.Metadata[k] = v
			}

			processedDocs = append(processedDocs, chunkDoc)
		}

		return processedDocs, nil
	}

	// For other types, return as-is for now
	return []Document{document}, nil
}

// ProcessDocuments processes multiple documents
func (sdp *SimpleDocumentProcessor) ProcessDocuments(ctx context.Context, documents []Document) ([]Document, error) {
	var allProcessed []Document

	for _, doc := range documents {
		processed, err := sdp.ProcessDocument(ctx, doc)
		if err != nil {
			return nil, fmt.Errorf("failed to process document %s: %w", doc.ID, err)
		}
		allProcessed = append(allProcessed, processed...)
	}

	return allProcessed, nil
}

// ExtractText extracts text from different document formats
func (sdp *SimpleDocumentProcessor) ExtractText(ctx context.Context, data []byte, docType DocumentType) (string, error) {
	switch docType {
	case DocumentTypeText:
		return string(data), nil
	case DocumentTypeJSON:
		// Simple JSON to text conversion
		return string(data), nil
	case DocumentTypePDF:
		// In a real implementation, you'd use a PDF library
		return "PDF content extraction not implemented - use a library like github.com/ledongthuc/pdf", nil
	default:
		return "", fmt.Errorf("unsupported document type: %s", docType)
	}
}

// ChunkText splits text into overlapping chunks
func (sdp *SimpleDocumentProcessor) ChunkText(text string, chunkSize int, overlap int) ([]string, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunk size must be positive")
	}

	if overlap >= chunkSize {
		return nil, fmt.Errorf("overlap must be less than chunk size")
	}

	var chunks []string
	textLen := len(text)

	if textLen <= chunkSize {
		return []string{text}, nil
	}

	start := 0
	for start < textLen {
		end := start + chunkSize
		if end > textLen {
			end = textLen
		}

		chunk := text[start:end]
		chunks = append(chunks, chunk)

		// Move start position accounting for overlap
		start = end - overlap
		if start >= textLen {
			break
		}
	}

	return chunks, nil
}

// GetSupportedTypes returns supported document types
func (sdp *SimpleDocumentProcessor) GetSupportedTypes() []DocumentType {
	return sdp.supportedTypes
}

// SimpleEmbeddingProvider implements basic embedding generation (mock)
type SimpleEmbeddingProvider struct {
	dimension int
	modelName string
}

// NewSimpleEmbeddingProvider creates a new embedding provider
func NewSimpleEmbeddingProvider(dimension int, modelName string) *SimpleEmbeddingProvider {
	return &SimpleEmbeddingProvider{
		dimension: dimension,
		modelName: modelName,
	}
}

// GenerateTextEmbedding generates embedding for text (mock implementation)
func (sep *SimpleEmbeddingProvider) GenerateTextEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Mock embedding generation - in reality, you'd call an embedding API
	embedding := make([]float32, sep.dimension)

	// Simple hash-based mock embedding
	hash := simpleHash(text)
	for i := range embedding {
		embedding[i] = float32((hash+i)%100) / 100.0
	}

	return embedding, nil
}

// GenerateTextEmbeddings generates embeddings for multiple texts
func (sep *SimpleEmbeddingProvider) GenerateTextEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	var embeddings [][]float32

	for _, text := range texts {
		embedding, err := sep.GenerateTextEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text: %w", err)
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

// GenerateImageEmbedding generates embedding for image data
func (sep *SimpleEmbeddingProvider) GenerateImageEmbedding(ctx context.Context, imageData []byte) ([]float32, error) {
	return nil, fmt.Errorf("image embedding not implemented")
}

// GenerateImageEmbeddings generates embeddings for multiple images
func (sep *SimpleEmbeddingProvider) GenerateImageEmbeddings(ctx context.Context, images [][]byte) ([][]float32, error) {
	return nil, fmt.Errorf("image embeddings not implemented")
}

// GenerateMultimodalEmbedding generates embedding for text and image
func (sep *SimpleEmbeddingProvider) GenerateMultimodalEmbedding(ctx context.Context, text string, imageData []byte) ([]float32, error) {
	return nil, fmt.Errorf("multimodal embedding not implemented")
}

// GetDimension returns the embedding dimension
func (sep *SimpleEmbeddingProvider) GetDimension() int {
	return sep.dimension
}

// GetModelName returns the model name
func (sep *SimpleEmbeddingProvider) GetModelName() string {
	return sep.modelName
}

// GetProviderName returns the provider name
func (sep *SimpleEmbeddingProvider) GetProviderName() string {
	return "SimpleEmbeddingProvider"
}

// Helper function for simple hash
func simpleHash(s string) int {
	hash := 0
	for _, c := range s {
		hash = hash*31 + int(c)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}
