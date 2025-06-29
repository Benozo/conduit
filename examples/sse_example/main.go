package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	conduit "github.com/benozo/conduit/lib"
	"github.com/benozo/conduit/lib/tools"
	"github.com/benozo/conduit/mcp"
)

func main() {
	log.Println("=== Conduit HTTP/SSE Example ===")
	log.Println("This server demonstrates HTTP API with Server-Sent Events")
	log.Println("Perfect for web applications and real-time integrations")

	// Create configuration for HTTP mode
	config := conduit.DefaultConfig()
	config.Mode = mcp.ModeHTTP
	config.Port = 8090
	config.EnableCORS = true
	config.EnableLogging = true

	// Create server
	server := conduit.NewServer(config)

	// Register all available tool packages
	tools.RegisterTextTools(server)
	tools.RegisterMemoryTools(server)
	tools.RegisterUtilityTools(server)

	// Register custom tools for SSE demonstration
	server.RegisterTool("sse_demo", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		return map[string]interface{}{
			"result":      "Hello from HTTP/SSE server!",
			"mode":        "HTTP with SSE",
			"description": "This tool demonstrates HTTP/SSE integration",
			"port":        8090,
			"endpoints":   []string{"/mcp", "/schema", "/health"},
		}, nil
	})

	server.RegisterTool("streaming_demo", func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
		text := "Streaming demonstration with chunked responses"
		chunks := []string{}

		// Split text into chunks for demonstration
		words := []string{"Streaming", "demonstration", "with", "chunked", "responses"}
		for i, word := range words {
			chunks = append(chunks, fmt.Sprintf("Chunk %d: %s", i+1, word))
		}

		return map[string]interface{}{
			"result":      text,
			"chunks":      chunks,
			"stream_info": "This demonstrates how tools can prepare data for streaming",
			"timestamp":   time.Now().Format(time.RFC3339),
		}, nil
	})

	// Start server in a goroutine so we can add our custom endpoints
	go func() {
		log.Printf("Starting HTTP/SSE server on port %d...", config.Port)
		log.Printf("Available endpoints:")
		log.Printf("  GET  http://localhost:%d/schema - Tool schemas", config.Port)
		log.Printf("  POST http://localhost:%d/mcp - MCP endpoint with SSE", config.Port)
		log.Printf("  GET  http://localhost:%d/health - Health check", config.Port)
		log.Printf("  GET  http://localhost:%d/demo - Custom demo page", config.Port)
		log.Printf("  GET  http://localhost:%d/sse-test - SSE test endpoint", config.Port)

		// Add custom demo endpoints before starting the server
		http.HandleFunc("/demo", demoPageHandler)
		http.HandleFunc("/sse-test", sseTestHandler)

		if err := server.Start(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Keep the main goroutine alive
	select {}
}

// Custom demo page handler
func demoPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Conduit HTTP/SSE Demo</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { margin: 20px 0; padding: 10px; border: 1px solid #ddd; }
        button { padding: 10px 20px; margin: 5px; }
        #output { margin-top: 20px; padding: 10px; background: #f5f5f5; }
        pre { white-space: pre-wrap; }
    </style>
</head>
<body>
    <h1>Conduit HTTP/SSE Demo</h1>
    <p>This page demonstrates the HTTP API and Server-Sent Events capabilities.</p>
    
    <div class="endpoint">
        <h3>Available Endpoints:</h3>
        <ul>
            <li><strong>GET /schema</strong> - List all available tools</li>
            <li><strong>POST /mcp</strong> - MCP protocol endpoint with SSE support</li>
            <li><strong>GET /health</strong> - Health check</li>
            <li><strong>GET /sse-test</strong> - SSE test endpoint</li>
        </ul>
    </div>

    <div class="endpoint">
        <h3>Quick Tests:</h3>
        <button onclick="testSchema()">Get Tool Schema</button>
        <button onclick="testHealth()">Health Check</button>
        <button onclick="testSSE()">Test SSE Stream</button>
        <button onclick="testTool()">Test Tool Call</button>
    </div>

    <div id="output"></div>

    <script>
        function log(message) {
            document.getElementById('output').innerHTML += '<pre>' + message + '</pre>';
        }

        function testSchema() {
            fetch('/schema')
                .then(r => r.json())
                .then(data => log('Schema: ' + JSON.stringify(data, null, 2)))
                .catch(e => log('Error: ' + e));
        }

        function testHealth() {
            fetch('/health')
                .then(r => r.text())
                .then(data => log('Health: ' + data))
                .catch(e => log('Error: ' + e));
        }

        function testSSE() {
            const eventSource = new EventSource('/sse-test');
            eventSource.onmessage = function(event) {
                log('SSE: ' + event.data);
            };
            eventSource.onerror = function(event) {
                log('SSE Error: Connection failed');
                eventSource.close();
            };
            setTimeout(() => eventSource.close(), 10000); // Close after 10 seconds
        }

        function testTool() {
            fetch('/mcp', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    jsonrpc: '2.0',
                    id: 1,
                    method: 'tools/call',
                    params: {name: 'sse_demo', arguments: {}}
                })
            })
            .then(r => r.text())
            .then(data => log('Tool Result: ' + data))
            .catch(e => log('Error: ' + e));
        }
    </script>
</body>
</html>`
	w.Write([]byte(html))
}

// SSE test handler
func sseTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send initial message
	fmt.Fprintf(w, "data: %s\n\n", `{"type": "start", "message": "SSE stream started"}`)
	flusher.Flush()

	// Send periodic messages
	for i := 1; i <= 5; i++ {
		time.Sleep(1 * time.Second)

		message := map[string]interface{}{
			"type":      "data",
			"message":   fmt.Sprintf("Streaming message %d", i),
			"timestamp": time.Now().Format(time.RFC3339),
			"count":     i,
		}

		jsonData, _ := json.Marshal(message)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}

	// Send completion message
	fmt.Fprintf(w, "data: %s\n\n", `{"type": "complete", "message": "SSE stream completed"}`)
	flusher.Flush()
}
