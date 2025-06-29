package conduit

import (
	"fmt"
	"log"

	"github.com/benozo/conduit/mcp"
)

// Server represents an embeddable MCP server
type Server struct {
	tools   *mcp.ToolRegistry
	memory  *mcp.Memory
	model   mcp.ModelFunc
	config  *Config
	unified *mcp.UnifiedServer
}

// Config holds server configuration
type Config struct {
	Port          int               `json:"port"`
	OllamaURL     string            `json:"ollama_url"`
	Mode          mcp.ServerMode    `json:"mode"`
	Environment   map[string]string `json:"environment"`
	EnableCORS    bool              `json:"enable_cors"`
	EnableHTTPS   bool              `json:"enable_https"`
	CertFile      string            `json:"cert_file"`
	KeyFile       string            `json:"key_file"`
	EnableLogging bool              `json:"enable_logging"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:          8080,
		OllamaURL:     "http://localhost:11434",
		Mode:          mcp.ModeBoth,
		Environment:   make(map[string]string),
		EnableCORS:    true,
		EnableHTTPS:   false,
		EnableLogging: true,
	}
}

// NewServer creates a new embeddable MCP server
func NewServer(config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	tools := mcp.NewToolRegistry()
	memory := mcp.NewMemory()

	server := &Server{
		tools:  tools,
		memory: memory,
		config: config,
	}

	return server
}

// NewServerWithModel creates a new server with a custom model
func NewServerWithModel(config *Config, model mcp.ModelFunc) *Server {
	server := NewServer(config)
	server.model = model
	return server
}

// RegisterTool adds a tool to the server
func (s *Server) RegisterTool(name string, tool mcp.ToolFunc) {
	s.tools.Register(name, tool)
}

// SetModel sets a custom model function
func (s *Server) SetModel(model mcp.ModelFunc) {
	s.model = model
}

// GetMemory returns the server's memory instance
func (s *Server) GetMemory() *mcp.Memory {
	return s.memory
}

// GetToolRegistry returns the tool registry for advanced usage
func (s *Server) GetToolRegistry() *mcp.ToolRegistry {
	return s.tools
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *Config {
	return s.config
}

// Start starts the server with the configured mode
func (s *Server) Start() error {
	if s.model == nil {
		s.model = createDefaultOllamaModel(s.config.OllamaURL)
	}

	s.unified = mcp.NewUnifiedServer(s.model, s.tools)
	s.unified.SetMode(s.config.Mode)

	if s.config.Port != 8080 {
		s.unified.SetPort(fmt.Sprintf(":%d", s.config.Port))
	}

	if s.config.EnableLogging {
		log.Printf("Starting conduit server on port %d (mode: %v)", s.config.Port, s.config.Mode)
	}

	return s.unified.Run()
}

// StartWithMode starts the server with a specific mode, overriding config
func (s *Server) StartWithMode(mode mcp.ServerMode) error {
	originalMode := s.config.Mode
	s.config.Mode = mode
	err := s.Start()
	s.config.Mode = originalMode // Restore original config
	return err
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	if s.unified != nil {
		return s.unified.Shutdown(nil)
	}
	return nil
}

// createDefaultOllamaModel creates a default Ollama model function
func createDefaultOllamaModel(ollamaURL string) mcp.ModelFunc {
	return CreateOllamaModel(ollamaURL)
}
