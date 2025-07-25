package transport

import (
	"testing"
	"time"
)

func TestStdioTransport(t *testing.T) {
	transport := NewStdioTransport(nil)

	if transport == nil {
		t.Fatal("NewStdioTransport returned nil")
	}

	// STDIO transport is considered connected by default since STDIO is always available
	if !transport.IsConnected() {
		t.Error("STDIO transport should be connected by default")
	}
}

func TestWebSocketTransport(t *testing.T) {
	transport := NewWebSocketTransport("ws://localhost:8080/test")

	if transport == nil {
		t.Fatal("NewWebSocketTransport returned nil")
	}

	// Test IsConnected initially false
	if transport.IsConnected() {
		t.Error("Transport should not be connected initially")
	}

	// Test header setting
	transport.SetHeader("Authorization", "Bearer token123")
	if transport.headers.Get("Authorization") != "Bearer token123" {
		t.Error("Header not set correctly")
	}
}

func TestHTTPTransport(t *testing.T) {
	transport := NewHTTPTransport("http://localhost:8080", "http://localhost:8080/events")

	if transport == nil {
		t.Fatal("NewHTTPTransport returned nil")
	}

	// Test timeout setting
	transport.SetTimeout(30 * time.Second)
	if transport.client.Timeout != 30*time.Second {
		t.Error("Timeout not set correctly")
	}

	// Test header setting
	transport.SetHeader("User-Agent", "test-client")
	if transport.headers["User-Agent"] != "test-client" {
		t.Error("Header not set correctly")
	}
}

func TestTransportErrorCodes(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrTransportClosed", ErrTransportClosed},
		{"ErrInvalidMessage", ErrInvalidMessage},
		{"ErrTimeout", ErrTimeout},
		{"ErrNotConnected", ErrNotConnected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() == "" {
				t.Errorf("Error %s should have a message", tt.name)
			}
		})
	}
}

func TestDefaultTransportOptions(t *testing.T) {
	opts := DefaultTransportOptions()

	if opts.BufferSize <= 0 {
		t.Error("Default buffer size should be positive")
	}

	if opts.Debug {
		t.Error("Debug should be false by default")
	}
}
