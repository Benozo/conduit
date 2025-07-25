package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/benozo/neuron-mcp/protocol"
	"github.com/gorilla/websocket"
)

// WebSocketTransport implements WebSocket transport for MCP
type WebSocketTransport struct {
	conn      *websocket.Conn
	url       string
	headers   http.Header
	readChan  chan *protocol.JSONRPCMessage
	writeChan chan *protocol.JSONRPCMessage
	closeChan chan struct{}
	connected bool
	mu        sync.RWMutex
}

// NewWebSocketTransport creates a new WebSocket transport
func NewWebSocketTransport(wsURL string) *WebSocketTransport {
	return &WebSocketTransport{
		url:       wsURL,
		headers:   make(http.Header),
		readChan:  make(chan *protocol.JSONRPCMessage, 100),
		writeChan: make(chan *protocol.JSONRPCMessage, 100),
		closeChan: make(chan struct{}),
	}
}

// SetHeader sets a custom header for the WebSocket connection
func (t *WebSocketTransport) SetHeader(key, value string) {
	t.headers.Set(key, value)
}

// Connect establishes the WebSocket connection
func (t *WebSocketTransport) Connect(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.connected {
		return nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, t.url, t.headers)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %w", err)
	}

	t.conn = conn
	t.connected = true

	// Start read and write goroutines
	go t.readLoop()
	go t.writeLoop()

	return nil
}

// Send sends a message via WebSocket
func (t *WebSocketTransport) Send(ctx context.Context, message *protocol.JSONRPCMessage) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if !t.connected {
		return fmt.Errorf("websocket not connected")
	}

	select {
	case t.writeChan <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-t.closeChan:
		return fmt.Errorf("websocket closed")
	}
}

// Receive gets the next message from the WebSocket
func (t *WebSocketTransport) Receive(ctx context.Context) (*protocol.JSONRPCMessage, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if !t.connected {
		return nil, fmt.Errorf("websocket not connected")
	}

	select {
	case message := <-t.readChan:
		return message, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-t.closeChan:
		return nil, fmt.Errorf("websocket closed")
	}
}

// Close closes the WebSocket connection
func (t *WebSocketTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.connected {
		return nil
	}

	t.connected = false
	close(t.closeChan)

	if t.conn != nil {
		// Send close frame
		t.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		return t.conn.Close()
	}

	return nil
}

// IsConnected returns true if the WebSocket is connected
func (t *WebSocketTransport) IsConnected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.connected
}

// readLoop continuously reads messages from the WebSocket
func (t *WebSocketTransport) readLoop() {
	defer func() {
		t.mu.Lock()
		t.connected = false
		t.mu.Unlock()
	}()

	for {
		select {
		case <-t.closeChan:
			return
		default:
		}

		messageType, data, err := t.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Log unexpected close
			}
			return
		}

		if messageType != websocket.TextMessage {
			continue // Skip non-text messages
		}

		var msg protocol.JSONRPCMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue // Skip invalid JSON
		}

		select {
		case t.readChan <- &msg:
		case <-t.closeChan:
			return
		default:
			// Buffer full, drop message
		}
	}
}

// writeLoop continuously writes messages to the WebSocket
func (t *WebSocketTransport) writeLoop() {
	ticker := time.NewTicker(54 * time.Second) // Ping every 54 seconds
	defer ticker.Stop()

	for {
		select {
		case <-t.closeChan:
			return
		case message := <-t.writeChan:
			data, err := json.Marshal(message)
			if err != nil {
				return
			}
			t.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := t.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			t.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := t.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WebSocketServer provides WebSocket server functionality
type WebSocketServer struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]chan []byte
	handler  MessageHandler
	mu       sync.RWMutex
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(handler MessageHandler) *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		clients: make(map[*websocket.Conn]chan []byte),
		handler: handler,
	}
}

// HandleWebSocket handles WebSocket upgrade requests
func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	clientChan := make(chan []byte, 100)
	s.mu.Lock()
	s.clients[conn] = clientChan
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		close(clientChan)
		conn.Close()
	}()

	// Start write goroutine
	go s.writeToClient(conn, clientChan)

	// Read messages from client
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Log unexpected close
			}
			break
		}

		if messageType == websocket.TextMessage {
			response, err := s.handler(data)
			if err == nil && response != nil {
				select {
				case clientChan <- response:
				default:
					// Client buffer full
				}
			}
		}
	}
}

// BroadcastToAllClients sends a message to all connected WebSocket clients
func (s *WebSocketServer) BroadcastToAllClients(message []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, clientChan := range s.clients {
		select {
		case clientChan <- message:
		default:
			// Client buffer full, skip
		}
	}
}

// writeToClient handles writing messages to a specific client
func (s *WebSocketServer) writeToClient(conn *websocket.Conn, clientChan chan []byte) {
	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-clientChan:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
