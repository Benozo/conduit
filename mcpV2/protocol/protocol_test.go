package protocol

import (
	"testing"
)

func TestJSONRPCMessage(t *testing.T) {
	// Test request message
	req := NewJSONRPCRequest("test_method", map[string]interface{}{"param": "value"})

	if req.Version != JSONRPCVersion {
		t.Errorf("Expected version %s, got %s", JSONRPCVersion, req.Version)
	}

	if req.Method != "test_method" {
		t.Errorf("Expected method test_method, got %s", req.Method)
	}

	if !req.IsRequest() {
		t.Error("Message should be identified as request")
	}

	if req.IsNotification() {
		t.Error("Message should not be identified as notification")
	}

	if req.IsResponse() {
		t.Error("Message should not be identified as response")
	}
}

func TestJSONRPCNotification(t *testing.T) {
	// Test notification message
	notif := NewJSONRPCNotification("test_notification", nil)

	if notif.Version != JSONRPCVersion {
		t.Errorf("Expected version %s, got %s", JSONRPCVersion, notif.Version)
	}

	if notif.Method != "test_notification" {
		t.Errorf("Expected method test_notification, got %s", notif.Method)
	}

	if notif.IsRequest() {
		t.Error("Message should not be identified as request")
	}

	if !notif.IsNotification() {
		t.Error("Message should be identified as notification")
	}

	if notif.IsResponse() {
		t.Error("Message should not be identified as response")
	}
}

func TestJSONRPCResponse(t *testing.T) {
	// Test response message
	resp := NewJSONRPCResponse("test_id", map[string]interface{}{"result": "success"})

	if resp.Version != JSONRPCVersion {
		t.Errorf("Expected version %s, got %s", JSONRPCVersion, resp.Version)
	}

	if resp.ID != "test_id" {
		t.Errorf("Expected ID test_id, got %v", resp.ID)
	}

	if resp.IsRequest() {
		t.Error("Message should not be identified as request")
	}

	if resp.IsNotification() {
		t.Error("Message should not be identified as notification")
	}

	if !resp.IsResponse() {
		t.Error("Message should be identified as response")
	}
}

func TestRPCError(t *testing.T) {
	err := NewRPCError(InvalidParams, "Invalid parameters")

	if err.Code != InvalidParams {
		t.Errorf("Expected code %d, got %d", InvalidParams, err.Code)
	}

	if err.Message != "Invalid parameters" {
		t.Errorf("Expected message 'Invalid parameters', got %s", err.Message)
	}

	if err.Error() != "RPC error -32602: Invalid parameters" {
		t.Errorf("Unexpected error string: %s", err.Error())
	}
}

func TestToolValidation(t *testing.T) {
	// Valid tool
	tool := &Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: JSONSchema{Type: "object"},
	}

	if err := tool.Validate(); err != nil {
		t.Errorf("Valid tool should not return error: %v", err)
	}

	// Invalid tool (no name)
	invalidTool := &Tool{
		Description: "A test tool without name",
		InputSchema: JSONSchema{Type: "object"},
	}

	if err := invalidTool.Validate(); err == nil {
		t.Error("Tool without name should return error")
	}
}
