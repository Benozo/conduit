package workflows

import (
	"context"
	"fmt"

	"github.com/benozo/neuron/src/agents/rag"
)

type RagWorkflow struct {
	Retriever *rag.Retriever
	Generator *rag.Generator
	Evaluator *rag.Evaluator
}

func NewRagWorkflow(retriever *rag.Retriever, generator *rag.Generator, evaluator *rag.Evaluator) *RagWorkflow {
	return &RagWorkflow{
		Retriever: retriever,
		Generator: generator,
		Evaluator: evaluator,
	}
}

func (w *RagWorkflow) Execute(query string) (string, error) {
	// Step 1: Retrieve information
	ctx := context.Background()
	retrievedData, err := w.Retriever.Retrieve(ctx, query)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err)
	}

	// Step 2: Generate content based on retrieved data
	generatedContent, err := w.Generator.GenerateContent(retrievedData)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	// Step 3: Evaluate the generated content
	isValid, err := w.Evaluator.Evaluate(generatedContent)
	if err != nil {
		return "", fmt.Errorf("evaluation failed: %w", err)
	}

	if !isValid {
		return "", fmt.Errorf("generated content is not valid")
	}

	return generatedContent, nil
}
