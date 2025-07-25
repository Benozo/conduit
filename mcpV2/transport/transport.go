// Package transport provides transport layer abstractions and implementations
// for the Model Context Protocol (MCP).
//
// This package defines the Transport interface and provides implementations
// for different communication channels including STDIO, HTTP/SSE, and WebSocket.
package transport

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/protocol"
)

// Common transport errors
var (
	ErrTransportClosed = errors.New("transport is closed")
	ErrInvalidMessage  = errors.New("invalid message")
	ErrTimeout         = errors.New("operation timed out")
	ErrNotConnected    = errors.New("not connected")
)

// Transport defines the interface for MCP message transport
type Transport interface {
	// Send sends a message through the transport
	Send(ctx context.Context, msg *protocol.JSONRPCMessage) error

	// Receive receives a message from the transport
	// This method should block until a message is available or the context is cancelled
	Receive(ctx context.Context) (*protocol.JSONRPCMessage, error)

	// Close closes the transport and releases any resources
	Close() error

	// IsConnected returns true if the transport is connected and ready to use
	IsConnected() bool
}

// TransportOptions configures transport behavior
type TransportOptions struct {
	// Buffer size for message queues
	BufferSize int

	// Enable debug logging
	Debug bool

	// Custom logger (optional)
	Logger Logger
}

// Logger interface for transport logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// DefaultTransportOptions returns default transport options
func DefaultTransportOptions() *TransportOptions {
	return &TransportOptions{
		BufferSize: 100,
		Debug:      false,
	}
}
