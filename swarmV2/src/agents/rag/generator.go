package rag

import (
	"fmt"
	"strings"
)

// Generator represents the agent responsible for generating content based on retrieved information.
type Generator struct {
	// ModelName specifies the language model to be used for content generation.
	ModelName string
}

// NewGenerator creates a new instance of the Generator agent with the specified model name.
func NewGenerator(modelName string) *Generator {
	return &Generator{
		ModelName: modelName,
	}
}

// GenerateContent takes the input data and generates content based on it.
func (g *Generator) GenerateContent(inputData string) (string, error) {
	if strings.TrimSpace(inputData) == "" {
		return "", fmt.Errorf("input data cannot be empty")
	}

	// Simulate content generation logic based on the input data.
	generatedContent := fmt.Sprintf("Generated content based on: %s using model: %s", inputData, g.ModelName)
	return generatedContent, nil
}