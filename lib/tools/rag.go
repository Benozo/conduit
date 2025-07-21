package tools

import (
	"github.com/benozo/conduit/lib/rag/tools"
)

// RegisterRAGTools registers all RAG-related tools with the server
func RegisterRAGTools(server ToolRegistrar) {
	// Document management tools
	server.RegisterTool("index_document", tools.IndexDocumentFunc)
	server.RegisterTool("delete_document", tools.DeleteDocumentFunc)
	server.RegisterTool("list_documents", tools.ListDocumentsFunc)
	server.RegisterTool("get_document", tools.GetDocumentFunc)
	server.RegisterTool("get_document_chunks", tools.GetDocumentChunksFunc)

	// Search and retrieval tools
	server.RegisterTool("semantic_search", tools.SemanticSearchFunc)
	server.RegisterTool("knowledge_query", tools.KnowledgeQueryFunc)

	// System tools
	server.RegisterTool("get_rag_stats", tools.GetRAGStatsFunc)
}

// RegisterRAGToolsWithSchema registers RAG tools with enhanced schemas
func RegisterRAGToolsWithSchema(server interface {
	RegisterToolWithSchema(string, interface{}, interface{})
}) {
	// Get tool metadata
	toolMetadata := tools.GetRAGToolMetadata()

	// Register each tool with its schema
	for _, meta := range toolMetadata {
		switch meta.Name {
		case "index_document":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.IndexDocumentFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		case "semantic_search":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.SemanticSearchFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		case "knowledge_query":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.KnowledgeQueryFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		case "list_documents":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.ListDocumentsFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		case "get_document":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.GetDocumentFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		case "delete_document":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.DeleteDocumentFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		case "get_document_chunks":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.GetDocumentChunksFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		case "get_rag_stats":
			server.RegisterToolWithSchema(
				meta.Name,
				tools.GetRAGStatsFunc,
				createToolMetadata(meta.Name, meta.Description, meta.Parameters, meta.Required),
			)
		}
	}
}

// Helper function to create tool metadata (would use the framework's actual metadata creation function)
func createToolMetadata(name, description string, parameters map[string]interface{}, required []string) interface{} {
	// This would use the actual metadata creation function from the framework
	// For now, return a simple map structure
	return map[string]interface{}{
		"name":        name,
		"description": description,
		"parameters":  parameters,
		"required":    required,
	}
}
