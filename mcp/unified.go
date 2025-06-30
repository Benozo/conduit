package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// ServerMode defines the server operating mode
type ServerMode int

const (
	// ModeStdio runs the server in stdio mode (for Copilot)
	ModeStdio ServerMode = iota
	// ModeHTTP runs the server in HTTP mode (for web apps)
	ModeHTTP
	// ModeBoth runs both stdio and HTTP servers
	ModeBoth
)

// UnifiedServer supports both stdio and HTTP protocols
type UnifiedServer struct {
	tools       *ToolRegistry
	memory      *Memory
	processor   *MCPProcessor
	stdioServer *StdioServer
	httpServer  *http.Server
	mode        ServerMode
	port        string
}

// NewUnifiedServer creates a new unified MCP server
func NewUnifiedServer(model ModelFunc, tools *ToolRegistry) *UnifiedServer {
	memory := NewMemory()
	processor := NewProcessor(model, tools)
	stdioServer := NewStdioServer(tools, memory)

	return &UnifiedServer{
		tools:       tools,
		memory:      memory,
		processor:   processor,
		stdioServer: stdioServer,
		mode:        ModeBoth,
		port:        ":8080",
	}
}

// NewUnifiedServerWithSchemaProvider creates a unified server with enhanced schema support
func NewUnifiedServerWithSchemaProvider(model ModelFunc, tools *ToolRegistry, schemaProvider EnhancedSchemaProvider) *UnifiedServer {
	memory := NewMemory()
	processor := NewProcessor(model, tools)
	stdioServer := NewStdioServerWithSchemaProvider(tools, memory, schemaProvider)

	return &UnifiedServer{
		tools:       tools,
		memory:      memory,
		processor:   processor,
		stdioServer: stdioServer,
		mode:        ModeBoth,
		port:        ":8080",
	}
}

// SetMode sets the server operating mode
func (s *UnifiedServer) SetMode(mode ServerMode) {
	s.mode = mode
}

// SetPort sets the HTTP server port
func (s *UnifiedServer) SetPort(port string) {
	s.port = port
}

// Run starts the server in the configured mode
func (s *UnifiedServer) Run() error {
	switch s.mode {
	case ModeStdio:
		return s.runStdio()
	case ModeHTTP:
		return s.runHTTP()
	case ModeBoth:
		return s.runBoth()
	default:
		return fmt.Errorf("unsupported server mode: %d", s.mode)
	}
}

// runStdio runs only the stdio server
func (s *UnifiedServer) runStdio() error {
	log.Println("Starting MCP server in stdio mode...")
	return s.stdioServer.Run()
}

// runHTTP runs only the HTTP server
func (s *UnifiedServer) runHTTP() error {
	log.Printf("Starting MCP server in HTTP mode on %s...", s.port)
	s.setupHTTPRoutes()
	return s.httpServer.ListenAndServe()
}

// runBoth runs both servers (stdio in background, HTTP in foreground)
func (s *UnifiedServer) runBoth() error {
	// Check if we're being called from stdio (no TTY)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// No TTY, run stdio mode
		log.Println("Detected stdio mode (no TTY)")
		return s.runStdio()
	}

	// TTY detected, run HTTP mode
	log.Println("Detected TTY, running HTTP mode")
	return s.runHTTP()
}

// setupHTTPRoutes sets up HTTP routes for the existing endpoints
func (s *UnifiedServer) setupHTTPRoutes() {
	mux := http.NewServeMux()

	// MCP endpoint (SSE)
	mux.HandleFunc("/mcp", s.handleMCPHTTP)

	// ReAct endpoint
	mux.HandleFunc("/react", s.handleReActHTTP)

	// Schema endpoint
	mux.HandleFunc("/schema", s.handleSchemaHTTP)

	// Health check
	mux.HandleFunc("/health", s.handleHealthHTTP)

	s.httpServer = &http.Server{
		Addr:    s.port,
		Handler: mux,
	}
}

// handleMCPHTTP handles the SSE MCP endpoint
func (s *UnifiedServer) handleMCPHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	s.processor.EnableStreaming(func(ctxID, token string) {
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ctxID, token)
		flusher.Flush()
	})

	result, err := s.processor.Run(req)
	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
		return
	}

	out, _ := json.Marshal(result)
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", out)
	flusher.Flush()
}

// handleReActHTTP handles the ReAct demonstration endpoint
func (s *UnifiedServer) handleReActHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Set initial memory
	s.processor.Memory.Set("latest", "hello world")

	thoughts := []string{
		"transform to uppercase",
		"no action",
	}

	steps, err := ReActAgent(thoughts, s.tools, s.processor.Memory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(steps)
}

// handleSchemaHTTP handles the schema endpoint
func (s *UnifiedServer) handleSchemaHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	tools := s.stdioServer.getToolSchemas()
	response := map[string]interface{}{"tools": tools}
	json.NewEncoder(w).Encode(response)
}

// handleHealthHTTP handles health checks
func (s *UnifiedServer) handleHealthHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"status":    "healthy",
		"server":    "conduit-unified",
		"protocols": []string{"stdio", "http"},
	}
	json.NewEncoder(w).Encode(response)
}

// Shutdown gracefully shuts down the server
func (s *UnifiedServer) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// RegisterTool registers a tool with the server
func (s *UnifiedServer) RegisterTool(name string, fn ToolFunc) {
	s.tools.Register(name, fn)
}

// GetMemory returns the server's memory instance
func (s *UnifiedServer) GetMemory() *Memory {
	return s.memory
}
