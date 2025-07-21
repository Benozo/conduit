package rag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TextProcessor handles plain text files
type TextProcessor struct{}

func (p *TextProcessor) ProcessFile(ctx context.Context, filePath string, metadata map[string]interface{}) (*Document, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract metadata
	title := filepath.Base(filePath)
	if titleMeta, exists := metadata["title"]; exists {
		if titleStr, ok := titleMeta.(string); ok && titleStr != "" {
			title = titleStr
		}
	}

	// Create document
	doc := &Document{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     string(content),
		SourcePath:  filePath,
		ContentType: ".txt",
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return doc, nil
}

func (p *TextProcessor) ProcessContent(ctx context.Context, content, title, contentType string, metadata map[string]interface{}) (*Document, error) {
	if title == "" {
		title = "Untitled Document"
	}

	doc := &Document{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     content,
		SourcePath:  "",
		ContentType: contentType,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return doc, nil
}

func (p *TextProcessor) ExtractText(ctx context.Context, filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

func (p *TextProcessor) GetSupportedTypes() []string {
	return []string{".txt", "text/plain"}
}

// MarkdownProcessor handles markdown files
type MarkdownProcessor struct{}

func (p *MarkdownProcessor) ProcessFile(ctx context.Context, filePath string, metadata map[string]interface{}) (*Document, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	// Extract title from markdown (first # heading or filename)
	title := extractMarkdownTitle(contentStr)
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	// Override with provided title if available
	if titleMeta, exists := metadata["title"]; exists {
		if titleStr, ok := titleMeta.(string); ok && titleStr != "" {
			title = titleStr
		}
	}

	// Create document
	doc := &Document{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     contentStr,
		SourcePath:  filePath,
		ContentType: ".md",
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return doc, nil
}

func (p *MarkdownProcessor) ProcessContent(ctx context.Context, content, title, contentType string, metadata map[string]interface{}) (*Document, error) {
	// Extract title from markdown if not provided
	if title == "" {
		title = extractMarkdownTitle(content)
		if title == "" {
			title = "Untitled Markdown Document"
		}
	}

	doc := &Document{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     content,
		SourcePath:  "",
		ContentType: contentType,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return doc, nil
}

func (p *MarkdownProcessor) ExtractText(ctx context.Context, filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// For now, just return raw markdown
	// In a full implementation, this would strip markdown formatting
	return string(content), nil
}

func (p *MarkdownProcessor) GetSupportedTypes() []string {
	return []string{".md", ".markdown", "text/markdown"}
}

// Helper function to extract title from markdown content
func extractMarkdownTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}
