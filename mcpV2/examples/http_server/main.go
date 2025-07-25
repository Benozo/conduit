// HTTP Server Example
//
// This example demonstrates setting up an MCP server with HTTP/SSE transport
// and middleware for logging, authentication, and rate limiting.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/middleware"
	"github.com/modelcontextprotocol/go-sdk/protocol"
	"github.com/modelcontextprotocol/go-sdk/server"
	"github.com/modelcontextprotocol/go-sdk/transport"
)

// SimpleLogger implements the middleware.Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

// SimpleMetrics implements the middleware.Metrics interface
type SimpleMetrics struct{}

func (m *SimpleMetrics) IncrementCounter(name string, tags map[string]string) {
	log.Printf("[METRICS] Counter %s: %v", name, tags)
}

func (m *SimpleMetrics) RecordDuration(name string, duration time.Duration, tags map[string]string) {
	log.Printf("[METRICS] Duration %s: %v (%v)", name, duration, tags)
}

func (m *SimpleMetrics) RecordGauge(name string, value float64, tags map[string]string) {
	log.Printf("[METRICS] Gauge %s: %f (%v)", name, value, tags)
}

// SimpleRateLimiter implements basic rate limiting
type SimpleRateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewSimpleRateLimiter(limit int, window time.Duration) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (r *SimpleRateLimiter) Allow(key string) bool {
	now := time.Now()

	// Clean old requests
	if timestamps, exists := r.requests[key]; exists {
		var valid []time.Time
		for _, t := range timestamps {
			if now.Sub(t) < r.window {
				valid = append(valid, t)
			}
		}
		r.requests[key] = valid
	}

	// Check if we're under the limit
	if len(r.requests[key]) >= r.limit {
		return false
	}

	// Add current request
	r.requests[key] = append(r.requests[key], now)
	return true
}

func main() {
	// Create server with HTTP transport
	srv := server.NewServer(&server.ServerOptions{
		Info: protocol.Implementation{
			Name:    "http-mcp-server",
			Version: "1.0.0",
		},
		Capabilities: protocol.ServerCapabilities{
			Tools: &protocol.ToolsCapability{
				ListChanged: true,
			},
		},
	})

	// Set up middleware chain (TODO: Implement middleware integration in server)
	logger := &SimpleLogger{}
	metrics := &SimpleMetrics{}
	rateLimiter := NewSimpleRateLimiter(10, time.Minute) // 10 requests per minute

	middlewareChain := middleware.NewChain(
		middleware.LoggingMiddleware(logger),
		middleware.MetricsMiddleware(metrics),
		middleware.RateLimitMiddleware(rateLimiter),
		middleware.ValidationMiddleware(),
		middleware.ErrorHandlingMiddleware(logger),
	)

	// Apply middleware to server (TODO: Implement this)
	// srv.SetMiddleware(middlewareChain)
	_ = middlewareChain // Suppress unused variable warning

	// Register tools
	if err := registerTools(srv); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// Create HTTP transport server with proper message handler
	httpServer := transport.NewHTTPServer(":8081", func(data []byte) ([]byte, error) {
		return handleMCPMessage(srv, data)
	})

	// Start HTTP server
	log.Println("Starting HTTP MCP server on :8081")
	log.Println("Available endpoints:")
	log.Println("  POST /mcp - JSON-RPC requests")
	log.Println("  GET  /mcp/events - Server-Sent Events")

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// registerTools registers tools with the server
func registerTools(srv *server.Server) error {
	// Text transformation tool
	textTool := &protocol.Tool{
		Name:        "text_transform",
		Description: "Transform text using various operations",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"text": {
					Type:        "string",
					Description: "The text to transform",
				},
				"operation": {
					Type:        "string",
					Description: "Transform operation: uppercase, lowercase, reverse",
					Enum:        []interface{}{"uppercase", "lowercase", "reverse"},
				},
			},
			Required: []string{"text", "operation"},
		},
	}

	if err := srv.RegisterTool(textTool, handleTextTransform); err != nil {
		return fmt.Errorf("failed to register text_transform tool: %w", err)
	}

	// Echo tool for testing
	echoTool := &protocol.Tool{
		Name:        "echo",
		Description: "Echo back the input with a timestamp",
		InputSchema: protocol.JSONSchema{
			Type: "object",
			Properties: map[string]*protocol.JSONSchema{
				"message": {
					Type:        "string",
					Description: "Message to echo back",
				},
			},
			Required: []string{"message"},
		},
	}

	if err := srv.RegisterTool(echoTool, handleEcho); err != nil {
		return fmt.Errorf("failed to register echo tool: %w", err)
	}

	return nil
}

// handleTextTransform implements the text transformation tool
func handleTextTransform(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
	text, ok := params["text"].(string)
	if !ok {
		return nil, fmt.Errorf("text parameter must be a string")
	}

	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter must be a string")
	}

	var result string
	switch operation {
	case "uppercase":
		result = strings.ToUpper(text)
	case "lowercase":
		result = strings.ToLower(text)
	case "reverse":
		result = reverseString(text)
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: result,
		}},
	}, nil
}

// handleEcho implements the echo tool
func handleEcho(ctx context.Context, params map[string]interface{}) (*protocol.CallToolResult, error) {
	message, ok := params["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter must be a string")
	}

	response := fmt.Sprintf("Echo at %s: %s", time.Now().Format(time.RFC3339), message)

	return &protocol.CallToolResult{
		Content: []protocol.Content{{
			Type: "text",
			Text: response,
		}},
	}, nil
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// handleMCPMessage processes incoming JSON-RPC messages through the MCP server
func handleMCPMessage(srv *server.Server, data []byte) ([]byte, error) {
	// Parse the incoming JSON-RPC message
	var msg protocol.JSONRPCMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		// Return JSON-RPC parse error
		errorResponse := protocol.NewJSONRPCError(
			nil, // ID is unknown due to parse error
			protocol.ParseError,
			"Parse error",
			nil,
		)
		return json.Marshal(errorResponse)
	}

	// Handle different methods
	ctx := context.Background()

	switch msg.Method {
	case "initialize":
		// Return initialization response for HTTP
		result := protocol.InitializeResult{
			ProtocolVersion: "2025-03-26",
			Capabilities: protocol.ServerCapabilities{
				Tools: &protocol.ToolsCapability{
					ListChanged: true,
				},
			},
			ServerInfo: protocol.Implementation{
				Name:    "http-mcp-server",
				Version: "1.0.0",
			},
		}
		response := protocol.NewJSONRPCResponse(msg.ID, result)
		return json.Marshal(response)

	case "tools/list":
		// Handle tools/list directly
		tools := srv.ListTools()
		result := protocol.ListToolsResult{
			Tools: make([]protocol.Tool, len(tools)),
		}
		for i, tool := range tools {
			result.Tools[i] = *tool
		}
		response := protocol.NewJSONRPCResponse(msg.ID, result)
		return json.Marshal(response)

	case "tools/call":
		// Handle tools/call directly
		var req protocol.ToolCallRequest
		if data, err := json.Marshal(msg.Params); err != nil {
			response := protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
			return json.Marshal(response)
		} else if err := json.Unmarshal(data, &req); err != nil {
			response := protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid parameters", nil)
			return json.Marshal(response)
		}

		// Get the tool handler
		if handler, exists := srv.GetToolHandler(req.Name); exists {
			// Convert arguments to map[string]interface{}
			var params map[string]interface{}
			if req.Arguments != nil {
				if data, err := json.Marshal(req.Arguments); err != nil {
					response := protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid arguments", nil)
					return json.Marshal(response)
				} else if err := json.Unmarshal(data, &params); err != nil {
					response := protocol.NewJSONRPCError(msg.ID, protocol.InvalidParams, "Invalid arguments", nil)
					return json.Marshal(response)
				}
			} else {
				params = make(map[string]interface{})
			}

			// Call the tool handler
			result, err := handler.Handler(ctx, params)
			if err != nil {
				response := protocol.NewJSONRPCError(
					msg.ID,
					protocol.InternalError,
					fmt.Sprintf("Tool call failed: %v", err),
					nil,
				)
				return json.Marshal(response)
			}
			response := protocol.NewJSONRPCResponse(msg.ID, result)
			return json.Marshal(response)
		}

		response := protocol.NewJSONRPCError(
			msg.ID,
			protocol.MethodNotFound,
			fmt.Sprintf("Tool not found: %s", req.Name),
			nil,
		)
		return json.Marshal(response)

	default:
		response := protocol.NewJSONRPCError(
			msg.ID,
			protocol.MethodNotFound,
			fmt.Sprintf("Method not found: %s", msg.Method),
			nil,
		)
		return json.Marshal(response)
	}
}

// tempInitialize sets the server as initialized for HTTP requests
func tempInitialize(srv *server.Server) {
	// This is a hack to make HTTP work - in a real implementation
	// this would be handled by a stateless wrapper or modified server
	initMsg := &protocol.JSONRPCMessage{
		Version: "2.0",
		ID:      json.RawMessage(`"init"`),
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2025-03-26",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "http-client",
				"version": "1.0.0",
			},
		},
	}
	srv.ProcessMessage(context.Background(), initMsg)
}
