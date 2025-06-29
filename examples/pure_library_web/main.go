package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/benozo/conduit/mcp"
)

// Custom web server using MCP components as pure library
func main() {
	// Initialize MCP components
	memory := mcp.NewMemory()
	tools := mcp.NewToolRegistry()

	// Register tools
	tools.Register("uppercase", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := params["text"].(string)
		return map[string]string{"result": strings.ToUpper(text)}, nil
	})

	tools.Register("store", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := params["key"].(string)
		value := params["value"]
		memory.Set(key, value)
		return map[string]string{"status": "stored"}, nil
	})

	tools.Register("retrieve", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		key := params["key"].(string)
		value := memory.Get(key)
		if value == nil {
			return map[string]string{"error": "not found"}, nil
		}
		return map[string]interface{}{"value": value}, nil
	})

	// HTTP handlers using MCP components
	http.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		// List available tools (hardcoded for this example since we know what we registered)
		toolsList := []string{"uppercase", "store", "retrieve"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"tools": toolsList})
	})

	http.HandleFunc("/tool/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		// Extract tool name from URL
		toolName := strings.TrimPrefix(r.URL.Path, "/tool/")

		// Parse request body
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		// Call tool using MCP
		result, err := tools.Call(toolName, req, memory)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/memory", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			// Return all memory contents (for demo purposes)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Memory contents are private, use specific keys",
			})
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
Custom MCP Web Server

Available endpoints:
- GET  /tools       - List available tools
- POST /tool/{name} - Execute a tool
- GET  /memory      - Memory info

Example usage:
curl -X POST http://localhost:8080/tool/uppercase -d '{"text":"hello world"}'
curl -X POST http://localhost:8080/tool/store -d '{"key":"name","value":"John"}'
curl -X POST http://localhost:8080/tool/retrieve -d '{"key":"name"}'
`)
	})

	fmt.Println("Custom MCP web server starting on :8080")
	fmt.Println("Using MCP components as pure library (no built-in server)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
