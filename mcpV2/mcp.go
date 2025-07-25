// Package mcp provides the official Go SDK for the Model Context Protocol (MCP).
//
// This package implements both client and server functionality for MCP,
// supporting multiple transport layers and providing both traditional
// client-server patterns and pure library usage.
//
// The SDK supports:
//   - JSON-RPC 2.0 protocol implementation
//   - Multiple transport layers (stdio, HTTP/SSE, WebSocket)
//   - Tool registration and execution
//   - Resource management
//   - Memory backends with multiple storage options
//   - Pure library usage without server overhead
//   - Comprehensive middleware system
//   - Progress tracking and observability
//
// Quick Start - Pure Library Usage:
//
//	registry := library.NewComponentRegistry()
//	registry.Tools().Register("echo", func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
//	    return &protocol.ToolResult{
//	        Content: []protocol.Content{{Type: "text", Text: params["text"].(string)}},
//	    }, nil
//	})
//	result, err := registry.Tools().Call(context.Background(), "echo", map[string]interface{}{"text": "Hello"})
//
// Quick Start - Client Usage:
//
//	transport := transport.NewStdioTransport(nil)
//	client := client.NewClient(transport, nil)
//	err := client.Connect(ctx, protocol.ClientCapabilities{})
//	result, err := client.CallTool(ctx, "echo", map[string]interface{}{"text": "Hello"})
//
// Quick Start - Server Usage:
//
//	server := server.NewServer(nil)
//	tool := &protocol.Tool{Name: "echo", InputSchema: protocol.JSONSchema{Type: "object"}}
//	server.RegisterTool(tool, handlerFunc)
//	transport := transport.NewStdioTransport(nil)
//	server.Serve(ctx, transport)
//
// For more examples and detailed documentation, see the examples/ directory
// and the comprehensive README.md file.
package mcp

// Re-export commonly used types and functions for convenience
import (
	"github.com/modelcontextprotocol/go-sdk/client"
	"github.com/modelcontextprotocol/go-sdk/library"
	"github.com/modelcontextprotocol/go-sdk/protocol"
	"github.com/modelcontextprotocol/go-sdk/server"
	"github.com/modelcontextprotocol/go-sdk/transport"
)

// Protocol types and constants
type (
	JSONRPCMessage     = protocol.JSONRPCMessage
	Tool               = protocol.Tool
	ToolResult         = protocol.ToolResult
	Content            = protocol.Content
	JSONSchema         = protocol.JSONSchema
	Implementation     = protocol.Implementation
	ClientCapabilities = protocol.ClientCapabilities
	ServerCapabilities = protocol.ServerCapabilities
	RPCError           = protocol.RPCError
	Meta               = protocol.Meta
)

// Client types
type (
	Client        = client.Client
	ClientOptions = client.ClientOptions
)

// Server types
type (
	Server        = server.Server
	ServerOptions = server.ServerOptions
	ToolFunc      = server.ToolFunc
)

// Transport types
type (
	Transport      = transport.Transport
	StdioTransport = transport.StdioTransport
)

// Library types
type (
	ComponentRegistry = library.ComponentRegistry
	ToolRegistry      = library.ToolRegistry
	Memory            = library.Memory
)

// Constructor functions
var (
	// Protocol constructors
	NewJSONRPCRequest      = protocol.NewJSONRPCRequest
	NewJSONRPCNotification = protocol.NewJSONRPCNotification
	NewJSONRPCResponse     = protocol.NewJSONRPCResponse
	NewJSONRPCError        = protocol.NewJSONRPCError

	// Client constructors
	NewClient = client.NewClient

	// Server constructors
	NewServer = server.NewServer

	// Transport constructors
	NewStdioTransport = transport.NewStdioTransport

	// Library constructors
	NewComponentRegistry = library.NewComponentRegistry
	NewToolRegistry      = library.NewToolRegistry
	NewMemory            = library.NewMemory
)

// Constants
const (
	JSONRPCVersion     = protocol.JSONRPCVersion
	MCPProtocolVersion = protocol.MCPProtocolVersion
)

// Standard error codes
const (
	ParseError           = protocol.ParseError
	InvalidRequest       = protocol.InvalidRequest
	MethodNotFound       = protocol.MethodNotFound
	InvalidParams        = protocol.InvalidParams
	InternalError        = protocol.InternalError
	NotInitialized       = protocol.NotInitialized
	RequestFailed        = protocol.RequestFailed
	InvalidTool          = protocol.InvalidTool
	InvalidResource      = protocol.InvalidResource
	MethodDisabled       = protocol.MethodDisabled
	InvalidPrompt        = protocol.InvalidPrompt
	AuthenticationFailed = protocol.AuthenticationFailed
	PermissionDenied     = protocol.PermissionDenied
	RateLimitExceeded    = protocol.RateLimitExceeded
)
