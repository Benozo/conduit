package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/protocol"
)

// StdioTransport implements Transport for stdio communication
type StdioTransport struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	scanner *bufio.Scanner
	encoder *json.Encoder

	closed bool
	mu     sync.RWMutex

	options *TransportOptions
}

// StdioOptions configures STDIO transport
type StdioOptions struct {
	*TransportOptions

	// Custom stdin/stdout/stderr (for testing)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewStdioTransport creates a new STDIO transport
func NewStdioTransport(opts *StdioOptions) *StdioTransport {
	if opts == nil {
		opts = &StdioOptions{
			TransportOptions: DefaultTransportOptions(),
		}
	}
	if opts.TransportOptions == nil {
		opts.TransportOptions = DefaultTransportOptions()
	}

	// Use provided streams or default to os.Stdin/Stdout/Stderr
	stdin := opts.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	return &StdioTransport{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		scanner: bufio.NewScanner(stdin),
		encoder: json.NewEncoder(stdout),
		options: opts.TransportOptions,
	}
}

// Send implements Transport.Send
func (t *StdioTransport) Send(ctx context.Context, msg *protocol.JSONRPCMessage) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.closed {
		return ErrTransportClosed
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Encode and send message as single line JSON
	if err := t.encoder.Encode(msg); err != nil {
		if t.options.Logger != nil {
			t.options.Logger.Error("Failed to encode message: %v", err)
		}
		return fmt.Errorf("failed to encode message: %w", err)
	}

	if t.options.Debug && t.options.Logger != nil {
		t.options.Logger.Debug("Sent message: %s %v", msg.Method, msg.ID)
	}

	return nil
}

// Receive implements Transport.Receive
func (t *StdioTransport) Receive(ctx context.Context) (*protocol.JSONRPCMessage, error) {
	t.mu.RLock()
	closed := t.closed
	t.mu.RUnlock()

	if closed {
		return nil, ErrTransportClosed
	}

	// Use a channel to make scanner.Scan() cancellable
	msgChan := make(chan *protocol.JSONRPCMessage, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(msgChan)
		defer close(errChan)

		for {
			if !t.scanner.Scan() {
				if err := t.scanner.Err(); err != nil {
					errChan <- fmt.Errorf("scanner error: %w", err)
					return
				}
				// EOF reached
				errChan <- io.EOF
				return
			}

			line := t.scanner.Text()
			if line == "" {
				continue // Skip empty lines
			}

			var msg protocol.JSONRPCMessage
			if err := json.Unmarshal([]byte(line), &msg); err != nil {
				if t.options.Logger != nil {
					t.options.Logger.Warn("Invalid JSON received: %s", line)
				}
				// Continue reading instead of returning error for invalid JSON
				continue
			}

			if t.options.Debug && t.options.Logger != nil {
				t.options.Logger.Debug("Received message: %s %v", msg.Method, msg.ID)
			}

			msgChan <- &msg
			return
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-msgChan:
		if msg == nil {
			// Channel was closed without sending a message
			return nil, io.EOF
		}
		return msg, nil
	case err := <-errChan:
		if err == nil {
			return nil, io.EOF
		}
		return nil, err
	}
}

// Close implements Transport.Close
func (t *StdioTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	if t.options.Debug && t.options.Logger != nil {
		t.options.Logger.Debug("STDIO transport closed")
	}

	return nil
}

// IsConnected implements Transport.IsConnected
func (t *StdioTransport) IsConnected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return !t.closed
}
