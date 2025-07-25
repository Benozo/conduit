// Mock WebSocket Server for Testing
//
// This is a simple WebSocket server that implements basic MCP protocol
// for testing the WebSocket client. Run this in a separate terminal.
//
// Usage: go run mock_server.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/benozo/neuron-mcp/protocol"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for testing
	},
}

func main() {
	http.HandleFunc("/mcp", handleWebSocket)

	log.Println("Mock WebSocket MCP server starting on :8082")
	log.Println("WebSocket endpoint: ws://localhost:8082/mcp")
	log.Println("Use this server to test the WebSocket client example")

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("New WebSocket connection from %s", r.RemoteAddr)

	for {
		var msg protocol.JSONRPCMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("Received: %s", msg.Method)

		response := handleMCPMessage(&msg)
		if response != nil {
			err = conn.WriteJSON(response)
			if err != nil {
				log.Printf("Write error: %v", err)
				break
			}
		}
	}

	log.Printf("WebSocket connection closed")
}

func handleMCPMessage(msg *protocol.JSONRPCMessage) *protocol.JSONRPCMessage {
	switch msg.Method {
	case "initialize":
		return &protocol.JSONRPCMessage{
			Version: "2.0",
			ID:      msg.ID,
			Result: protocol.InitializeResult{
				ProtocolVersion: "2025-03-26",
				Capabilities: protocol.ServerCapabilities{
					Tools: &protocol.ToolsCapability{
						ListChanged: true,
					},
				},
				ServerInfo: protocol.Implementation{
					Name:    "mock-websocket-server",
					Version: "1.0.0",
				},
			},
		}

	case "tools/list":
		tools := []protocol.Tool{
			{
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
			},
			{
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
			},
			{
				Name:        "long_task",
				Description: "Simulate a long-running task with progress updates",
				InputSchema: protocol.JSONSchema{
					Type: "object",
					Properties: map[string]*protocol.JSONSchema{
						"duration": {
							Type:        "number",
							Description: "Task duration in seconds",
						},
						"steps": {
							Type:        "number",
							Description: "Number of progress steps",
						},
					},
					Required: []string{"duration", "steps"},
				},
			},
		}

		return &protocol.JSONRPCMessage{
			Version: "2.0",
			ID:      msg.ID,
			Result: protocol.ListToolsResult{
				Tools: tools,
			},
		}

	case "tools/call":
		var req protocol.ToolCallRequest
		if data, err := json.Marshal(msg.Params); err == nil {
			json.Unmarshal(data, &req)
		}

		// Convert Arguments to map[string]interface{}
		args, ok := req.Arguments.(map[string]interface{})
		if !ok {
			return &protocol.JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Error: &protocol.RPCError{
					Code:    protocol.InvalidParams,
					Message: "Invalid arguments format",
				},
			}
		}

		switch req.Name {
		case "echo":
			message, ok := args["message"].(string)
			if !ok {
				return &protocol.JSONRPCMessage{
					Version: "2.0",
					ID:      msg.ID,
					Error: &protocol.RPCError{
						Code:    protocol.InvalidParams,
						Message: "Invalid message parameter",
					},
				}
			}
			return &protocol.JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Result: protocol.CallToolResult{
					Content: []protocol.Content{{
						Type: "text",
						Text: "Echo at " + time.Now().Format(time.RFC3339) + ": " + message,
					}},
				},
			}

		case "text_transform":
			text, textOk := args["text"].(string)
			operation, opOk := args["operation"].(string)

			if !textOk || !opOk {
				return &protocol.JSONRPCMessage{
					Version: "2.0",
					ID:      msg.ID,
					Error: &protocol.RPCError{
						Code:    protocol.InvalidParams,
						Message: "Invalid text or operation parameter",
					},
				}
			}

			var result string
			switch operation {
			case "uppercase":
				result = "TRANSFORMED: " + text
			case "lowercase":
				result = "transformed: " + text
			case "reverse":
				runes := []rune(text)
				for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
					runes[i], runes[j] = runes[j], runes[i]
				}
				result = string(runes)
			default:
				result = "Unknown operation: " + operation
			}

			return &protocol.JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Result: protocol.CallToolResult{
					Content: []protocol.Content{{
						Type: "text",
						Text: result,
					}},
				},
			}

		case "long_task":
			// In a real implementation, this would start a background task
			// and send progress notifications. For this mock, we just return success.
			return &protocol.JSONRPCMessage{
				Version: "2.0",
				ID:      msg.ID,
				Result: protocol.CallToolResult{
					Content: []protocol.Content{{
						Type: "text",
						Text: "Long task completed (mock implementation)",
					}},
				},
			}
		}

		return &protocol.JSONRPCMessage{
			Version: "2.0",
			ID:      msg.ID,
			Error: &protocol.RPCError{
				Code:    protocol.MethodNotFound,
				Message: "Tool not found: " + req.Name,
			},
		}

	default:
		return &protocol.JSONRPCMessage{
			Version: "2.0",
			ID:      msg.ID,
			Error: &protocol.RPCError{
				Code:    protocol.MethodNotFound,
				Message: "Method not found: " + msg.Method,
			},
		}
	}
}
