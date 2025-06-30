package conduit

import (
	"fmt"
	"log"

	"github.com/benozo/conduit/mcp"
)

// ToolMetadata contains schema information for custom tools
type ToolMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// EnhancedServer extends the base Conduit server with metadata support
type EnhancedServer struct {
	*Server
	toolMetadata map[string]ToolMetadata
}

// NewEnhancedServer creates a new enhanced server with metadata support
func NewEnhancedServer(config *Config) *EnhancedServer {
	return &EnhancedServer{
		Server:       NewServer(config),
		toolMetadata: make(map[string]ToolMetadata),
	}
}

// RegisterToolWithSchema registers a tool with full schema metadata
func (es *EnhancedServer) RegisterToolWithSchema(name string, tool mcp.ToolFunc, metadata ToolMetadata) {
	// Register the tool with the base server
	es.Server.RegisterTool(name, tool)

	// Store metadata for schema generation
	es.toolMetadata[name] = metadata

	fmt.Printf("✓ Registered tool '%s': %s\n", name, metadata.Description)
}

// GetToolMetadata returns all stored tool metadata (implements EnhancedSchemaProvider)
func (es *EnhancedServer) GetToolMetadata() map[string]interface{} {
	result := make(map[string]interface{})
	for name, metadata := range es.toolMetadata {
		result[name] = map[string]interface{}{
			"name":        metadata.Name,
			"description": metadata.Description,
			"inputSchema": metadata.InputSchema,
		}
	}
	return result
}

// GetToolSchema returns the schema for a specific tool (implements EnhancedSchemaProvider)
func (es *EnhancedServer) GetToolSchema(toolName string) (interface{}, bool) {
	metadata, exists := es.toolMetadata[toolName]
	if !exists {
		return nil, false
	}

	return map[string]interface{}{
		"name":        metadata.Name,
		"description": metadata.Description,
		"inputSchema": metadata.InputSchema,
	}, true
}

// GetCustomToolCount returns the number of custom tools with metadata
func (es *EnhancedServer) GetCustomToolCount() int {
	return len(es.toolMetadata)
}

// ListCustomTools returns a list of custom tool names and descriptions
func (es *EnhancedServer) ListCustomTools() []map[string]string {
	var tools []map[string]string
	for name, metadata := range es.toolMetadata {
		tools = append(tools, map[string]string{
			"name":        name,
			"description": metadata.Description,
		})
	}
	return tools
}

// Start starts the enhanced server with schema provider support
func (es *EnhancedServer) Start() error {
	// Configure the model if not set
	if es.Server.model == nil {
		es.Server.model = CreateOllamaModel(es.Server.config.OllamaURL)
	}

	// Create unified server with enhanced schema support
	es.Server.unified = mcp.NewUnifiedServerWithSchemaProvider(es.Server.model, es.Server.tools, es)
	es.Server.unified.SetMode(es.Server.config.Mode)

	if es.Server.config.Mode == mcp.ModeHTTP || es.Server.config.Mode == mcp.ModeBoth {
		es.Server.unified.SetPort(fmt.Sprintf(":%d", es.Server.config.Port))
	}

	log.Printf("Starting enhanced server with %d custom tools...", len(es.toolMetadata))
	return es.Server.unified.Run()
}

// Helper function to create number parameter schema
func NumberParam(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "number",
		"description": description,
	}
}

// Helper function to create string parameter schema
func StringParam(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
	}
}

// Helper function to create array parameter schema
func ArrayParam(description string, itemType string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "array",
		"description": description,
		"items":       map[string]interface{}{"type": itemType},
	}
}

// Helper function to create boolean parameter schema
func BoolParam(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "boolean",
		"description": description,
	}
}

// Helper function to create enum parameter schema
func EnumParam(description string, values []string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
		"enum":        values,
	}
}

// Helper function to create object schema
func CreateObjectSchema(properties map[string]interface{}, required []string) map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
}

// CreateToolMetadata is a convenience function to create ToolMetadata
func CreateToolMetadata(name, description string, properties map[string]interface{}, required []string) ToolMetadata {
	return ToolMetadata{
		Name:        name,
		Description: description,
		InputSchema: CreateObjectSchema(properties, required),
	}
}
