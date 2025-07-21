package processors

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/benozo/conduit/lib/rag"
)

// ChunkingStrategy defines different text chunking approaches
type ChunkingStrategy string

const (
	FixedSize ChunkingStrategy = "fixed"
	Semantic  ChunkingStrategy = "semantic"
	Paragraph ChunkingStrategy = "paragraph"
	Sentence  ChunkingStrategy = "sentence"
)

// TextChunker implements text chunking with various strategies
type TextChunker struct {
	strategy ChunkingStrategy
	size     int
	overlap  int
}

// NewTextChunker creates a new text chunker
func NewTextChunker(strategy ChunkingStrategy, size, overlap int) *TextChunker {
	return &TextChunker{
		strategy: strategy,
		size:     size,
		overlap:  overlap,
	}
}

// ChunkText splits text into chunks based on the configured strategy
func (tc *TextChunker) ChunkText(ctx context.Context, text string) ([]rag.TextChunk, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	switch tc.strategy {
	case FixedSize:
		return tc.chunkByFixedSize(text), nil
	case Semantic:
		return tc.chunkBySemantic(text), nil
	case Paragraph:
		return tc.chunkByParagraph(text), nil
	case Sentence:
		return tc.chunkBySentence(text), nil
	default:
		return nil, fmt.Errorf("unsupported chunking strategy: %s", tc.strategy)
	}
}

// GetStrategy returns the current chunking strategy
func (tc *TextChunker) GetStrategy() string {
	return string(tc.strategy)
}

// Configure updates chunking parameters
func (tc *TextChunker) Configure(size, overlap int, strategy string) error {
	if size <= 0 {
		return fmt.Errorf("chunk size must be positive")
	}
	if overlap < 0 {
		return fmt.Errorf("overlap cannot be negative")
	}
	if overlap >= size {
		return fmt.Errorf("overlap must be less than chunk size")
	}

	tc.size = size
	tc.overlap = overlap
	tc.strategy = ChunkingStrategy(strategy)

	return nil
}

// chunkByFixedSize splits text into fixed-size chunks with overlap
func (tc *TextChunker) chunkByFixedSize(text string) []rag.TextChunk {
	if len(text) <= tc.size {
		return []rag.TextChunk{{
			Content: text,
			Index:   0,
			Metadata: map[string]interface{}{
				"strategy": "fixed",
				"size":     len(text),
			},
		}}
	}

	var chunks []rag.TextChunk
	index := 0
	start := 0

	for start < len(text) {
		end := start + tc.size
		if end > len(text) {
			end = len(text)
		}

		// Try to break at word boundary
		if end < len(text) {
			end = tc.findWordBoundary(text, end)
		}

		chunk := rag.TextChunk{
			Content: text[start:end],
			Index:   index,
			Metadata: map[string]interface{}{
				"strategy":      "fixed",
				"size":          end - start,
				"start_pos":     start,
				"end_pos":       end,
				"word_boundary": end < len(text),
			},
		}

		chunks = append(chunks, chunk)

		// Move start position with overlap
		start = end - tc.overlap
		if start <= 0 {
			start = end
		}
		index++
	}

	return chunks
}

// chunkBySemantic splits text by semantic boundaries (paragraphs and sentences)
func (tc *TextChunker) chunkBySemantic(text string) []rag.TextChunk {
	// Split by paragraphs first
	paragraphs := strings.Split(text, "\n\n")
	var chunks []rag.TextChunk
	index := 0
	currentChunk := ""
	startPos := 0

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// If adding this paragraph would exceed size, save current chunk
		if len(currentChunk)+len(paragraph)+2 > tc.size && currentChunk != "" {
			chunk := rag.TextChunk{
				Content: strings.TrimSpace(currentChunk),
				Index:   index,
				Metadata: map[string]interface{}{
					"strategy":   "semantic",
					"size":       len(currentChunk),
					"start_pos":  startPos,
					"boundaries": "paragraph",
				},
			}
			chunks = append(chunks, chunk)

			// Start new chunk with overlap
			if tc.overlap > 0 {
				overlapText := tc.getLastWords(currentChunk, tc.overlap)
				currentChunk = overlapText + " " + paragraph
			} else {
				currentChunk = paragraph
			}
			startPos += len(chunk.Content) - tc.overlap
			index++
		} else {
			if currentChunk != "" {
				currentChunk += "\n\n" + paragraph
			} else {
				currentChunk = paragraph
			}
		}
	}

	// Add remaining chunk
	if currentChunk != "" {
		chunk := rag.TextChunk{
			Content: strings.TrimSpace(currentChunk),
			Index:   index,
			Metadata: map[string]interface{}{
				"strategy":   "semantic",
				"size":       len(currentChunk),
				"start_pos":  startPos,
				"boundaries": "paragraph",
			},
		}
		chunks = append(chunks, chunk)
	}

	return chunks
}

// chunkByParagraph splits text by paragraph boundaries
func (tc *TextChunker) chunkByParagraph(text string) []rag.TextChunk {
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(text, -1)
	var chunks []rag.TextChunk
	index := 0

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// If paragraph is too long, split it further
		if len(paragraph) > tc.size {
			subChunks := tc.splitLongParagraph(paragraph, index)
			chunks = append(chunks, subChunks...)
			index += len(subChunks)
		} else {
			chunk := rag.TextChunk{
				Content: paragraph,
				Index:   index,
				Metadata: map[string]interface{}{
					"strategy":   "paragraph",
					"size":       len(paragraph),
					"boundaries": "paragraph",
				},
			}
			chunks = append(chunks, chunk)
			index++
		}
	}

	return chunks
}

// chunkBySentence splits text by sentence boundaries
func (tc *TextChunker) chunkBySentence(text string) []rag.TextChunk {
	sentences := tc.splitIntoSentences(text)
	var chunks []rag.TextChunk
	index := 0
	currentChunk := ""

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// If adding this sentence would exceed size, save current chunk
		if len(currentChunk)+len(sentence)+1 > tc.size && currentChunk != "" {
			chunk := rag.TextChunk{
				Content: strings.TrimSpace(currentChunk),
				Index:   index,
				Metadata: map[string]interface{}{
					"strategy":   "sentence",
					"size":       len(currentChunk),
					"boundaries": "sentence",
				},
			}
			chunks = append(chunks, chunk)

			// Start new chunk with overlap
			if tc.overlap > 0 {
				overlapText := tc.getLastSentences(currentChunk, tc.overlap)
				currentChunk = overlapText + " " + sentence
			} else {
				currentChunk = sentence
			}
			index++
		} else {
			if currentChunk != "" {
				currentChunk += " " + sentence
			} else {
				currentChunk = sentence
			}
		}
	}

	// Add remaining chunk
	if currentChunk != "" {
		chunk := rag.TextChunk{
			Content: strings.TrimSpace(currentChunk),
			Index:   index,
			Metadata: map[string]interface{}{
				"strategy":   "sentence",
				"size":       len(currentChunk),
				"boundaries": "sentence",
			},
		}
		chunks = append(chunks, chunk)
	}

	return chunks
}

// Helper functions

// findWordBoundary finds the nearest word boundary before the given position
func (tc *TextChunker) findWordBoundary(text string, pos int) int {
	if pos >= len(text) {
		return len(text)
	}

	// Look backwards for a space or punctuation
	for i := pos - 1; i >= 0 && i > pos-50; i-- {
		if unicode.IsSpace(rune(text[i])) || unicode.IsPunct(rune(text[i])) {
			return i + 1
		}
	}

	// If no good boundary found, use original position
	return pos
}

// splitLongParagraph splits a paragraph that's too long into smaller chunks
func (tc *TextChunker) splitLongParagraph(paragraph string, startIndex int) []rag.TextChunk {
	// Use fixed-size chunking for long paragraphs
	tempChunker := &TextChunker{
		strategy: FixedSize,
		size:     tc.size,
		overlap:  tc.overlap,
	}

	chunks := tempChunker.chunkByFixedSize(paragraph)

	// Update metadata and indices
	for i := range chunks {
		chunks[i].Index = startIndex + i
		chunks[i].Metadata["original_strategy"] = "paragraph"
		chunks[i].Metadata["split_reason"] = "paragraph_too_long"
	}

	return chunks
}

// splitIntoSentences splits text into sentences using simple rules
func (tc *TextChunker) splitIntoSentences(text string) []string {
	// Simple sentence splitting regex
	sentenceEnders := regexp.MustCompile(`[.!?]+\s+`)
	sentences := sentenceEnders.Split(text, -1)

	// Clean up and filter empty sentences
	var result []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			result = append(result, sentence)
		}
	}

	return result
}

// getLastWords gets the last N characters worth of words for overlap
func (tc *TextChunker) getLastWords(text string, maxChars int) string {
	if len(text) <= maxChars {
		return text
	}

	// Find word boundary
	start := len(text) - maxChars
	for start < len(text) && !unicode.IsSpace(rune(text[start])) {
		start++
	}

	return strings.TrimSpace(text[start:])
}

// getLastSentences gets the last sentences that fit within maxChars
func (tc *TextChunker) getLastSentences(text string, maxChars int) string {
	sentences := tc.splitIntoSentences(text)
	if len(sentences) == 0 {
		return ""
	}

	result := ""
	for i := len(sentences) - 1; i >= 0; i-- {
		candidate := sentences[i]
		if i < len(sentences)-1 {
			candidate += " " + result
		}

		if len(candidate) <= maxChars {
			result = candidate
		} else {
			break
		}
	}

	return result
}

// ChunkingStats provides statistics about chunking results
type ChunkingStats struct {
	TotalChunks  int     `json:"total_chunks"`
	AvgChunkSize float64 `json:"avg_chunk_size"`
	MinChunkSize int     `json:"min_chunk_size"`
	MaxChunkSize int     `json:"max_chunk_size"`
	Strategy     string  `json:"strategy"`
	OverlapUsed  int     `json:"overlap_used"`
}

// GetChunkingStats calculates statistics for a set of chunks
func GetChunkingStats(chunks []rag.TextChunk, strategy string, overlap int) ChunkingStats {
	if len(chunks) == 0 {
		return ChunkingStats{Strategy: strategy, OverlapUsed: overlap}
	}

	totalSize := 0
	minSize := len(chunks[0].Content)
	maxSize := len(chunks[0].Content)

	for _, chunk := range chunks {
		size := len(chunk.Content)
		totalSize += size

		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
	}

	return ChunkingStats{
		TotalChunks:  len(chunks),
		AvgChunkSize: float64(totalSize) / float64(len(chunks)),
		MinChunkSize: minSize,
		MaxChunkSize: maxSize,
		Strategy:     strategy,
		OverlapUsed:  overlap,
	}
}
