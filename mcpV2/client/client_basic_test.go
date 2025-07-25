package client

import (
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/transport"
)

func TestNewClient(t *testing.T) {
	// Use STDIO transport for testing
	stdio := transport.NewStdioTransport(nil)

	client := NewClient(stdio, &ClientOptions{
		Timeout: 5 * time.Second,
	})

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.options.Timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", client.options.Timeout)
	}
}

func TestClientOptions(t *testing.T) {
	opts := DefaultClientOptions()

	if opts.Timeout <= 0 {
		t.Error("Default timeout should be positive")
	}

	if opts.ConnectTimeout <= 0 {
		t.Error("Default connect timeout should be positive")
	}
}

func TestClient_IsConnected(t *testing.T) {
	stdio := transport.NewStdioTransport(nil)
	client := NewClient(stdio, nil)

	// Initially not connected
	if client.IsConnected() {
		t.Error("Client should not be connected initially")
	}
}
