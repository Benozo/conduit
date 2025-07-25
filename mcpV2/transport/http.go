package transport

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/protocol"
)

// HTTPTransport implements HTTP/SSE transport for MCP
type HTTPTransport struct {
	client    *http.Client
	serverURL string
	sseURL    string
	eventChan chan []byte
	closeChan chan struct{}
	headers   map[string]string
	connected bool
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(serverURL, sseURL string) *HTTPTransport {
	return &HTTPTransport{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		serverURL: serverURL,
		sseURL:    sseURL,
		eventChan: make(chan []byte, 100),
		closeChan: make(chan struct{}),
		headers:   make(map[string]string),
	}
}

// SetHeader sets a custom header for HTTP requests
func (t *HTTPTransport) SetHeader(key, value string) {
	t.headers[key] = value
}

// SetTimeout sets the HTTP client timeout
func (t *HTTPTransport) SetTimeout(timeout time.Duration) {
	t.client.Timeout = timeout
}

// Connect establishes the HTTP/SSE connection
func (t *HTTPTransport) Connect(ctx context.Context) error {
	if t.connected {
		return nil
	}

	// Start SSE connection for receiving messages
	if t.sseURL != "" {
		go t.startSSEConnection(ctx)
	}

	t.connected = true
	return nil
}

// Send sends a message via HTTP POST
func (t *HTTPTransport) Send(ctx context.Context, message []byte) error {
	if !t.connected {
		return fmt.Errorf("transport not connected")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.serverURL, bytes.NewReader(message))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	// For request-response pattern, read the response
	if resp.Header.Get("Content-Type") == "application/json" {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response: %w", err)
		}

		// Send response to event channel
		select {
		case t.eventChan <- responseBody:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// Receive gets the next message from the transport
func (t *HTTPTransport) Receive(ctx context.Context) ([]byte, error) {
	if !t.connected {
		return nil, fmt.Errorf("transport not connected")
	}

	select {
	case message := <-t.eventChan:
		return message, nil
	case <-t.closeChan:
		return nil, fmt.Errorf("transport closed")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Close closes the HTTP transport
func (t *HTTPTransport) Close() error {
	if !t.connected {
		return nil
	}

	t.connected = false
	close(t.closeChan)
	return nil
}

// startSSEConnection starts the Server-Sent Events connection
func (t *HTTPTransport) startSSEConnection(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.closeChan:
			return
		default:
		}

		err := t.connectSSE(ctx)
		if err != nil {
			// Retry after delay
			select {
			case <-ctx.Done():
				return
			case <-t.closeChan:
				return
			case <-time.After(5 * time.Second):
				continue
			}
		}
	}
}

// connectSSE establishes a single SSE connection
func (t *HTTPTransport) connectSSE(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", t.sseURL, nil)
	if err != nil {
		return fmt.Errorf("creating SSE request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("connecting to SSE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SSE connection failed with status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var eventData strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.closeChan:
			return nil
		default:
		}

		if line == "" {
			// Empty line indicates end of event
			if eventData.Len() > 0 {
				data := eventData.String()
				eventData.Reset()

				select {
				case t.eventChan <- []byte(data):
				case <-ctx.Done():
					return ctx.Err()
				case <-t.closeChan:
					return nil
				}
			}
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			eventData.WriteString(data)
		}
		// Ignore other SSE fields like "event:", "id:", "retry:"
	}

	return scanner.Err()
}

// HTTPServer provides HTTP server functionality for MCP
type HTTPServer struct {
	server     *http.Server
	sseClients map[string]chan []byte
	handler    MessageHandler
}

// MessageHandler handles incoming HTTP messages
type MessageHandler func([]byte) ([]byte, error)

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(addr string, handler MessageHandler) *HTTPServer {
	s := &HTTPServer{
		sseClients: make(map[string]chan []byte),
		handler:    handler,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", s.handleMCP)
	mux.HandleFunc("/mcp/events", s.handleSSE)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server
func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// BroadcastEvent sends an event to all SSE clients
func (s *HTTPServer) BroadcastEvent(message []byte) {
	for _, client := range s.sseClients {
		select {
		case client <- message:
		default:
			// Client buffer full, skip
		}
	}
}

// handleMCP handles JSON-RPC requests
func (s *HTTPServer) handleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	response, err := s.handler(body)
	if err != nil {
		// Create JSON-RPC error response
		var req protocol.JSONRPCMessage
		json.Unmarshal(body, &req)

		errorResp := protocol.JSONRPCMessage{
			Version: "2.0",
			ID:      req.ID,
			Error: &protocol.RPCError{
				Code:    protocol.InternalError,
				Message: err.Error(),
			},
		}

		response, _ = json.Marshal(errorResp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// handleSSE handles Server-Sent Events connections
func (s *HTTPServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Server-Sent Events not supported", http.StatusInternalServerError)
		return
	}

	clientID := fmt.Sprintf("%p", r)
	clientChan := make(chan []byte, 100)
	s.sseClients[clientID] = clientChan

	defer func() {
		delete(s.sseClients, clientID)
		close(clientChan)
	}()

	for {
		select {
		case message := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", string(message))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
