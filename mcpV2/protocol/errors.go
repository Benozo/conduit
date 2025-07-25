package protocol

import (
	"errors"
	"fmt"
)

// Common protocol errors
var (
	ErrNotInitialized       = errors.New("protocol not initialized")
	ErrInvalidMessage       = errors.New("invalid message format")
	ErrInvalidMethod        = errors.New("invalid method")
	ErrMethodNotFound       = errors.New("method not found")
	ErrInvalidParams        = errors.New("invalid parameters")
	ErrRequestFailed        = errors.New("request failed")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
)

// NewRPCError creates a new RPC error with the given code and message
func NewRPCError(code int, message string) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
	}
}

// NewRPCErrorWithData creates a new RPC error with data
func NewRPCErrorWithData(code int, message string, data interface{}) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Standard error constructors for common cases
func NewParseError(message string) *RPCError {
	if message == "" {
		message = "Parse error"
	}
	return NewRPCError(ParseError, message)
}

func NewInvalidRequestError(message string) *RPCError {
	if message == "" {
		message = "Invalid Request"
	}
	return NewRPCError(InvalidRequest, message)
}

func NewMethodNotFoundError(method string) *RPCError {
	message := "Method not found"
	if method != "" {
		message = fmt.Sprintf("Method not found: %s", method)
	}
	return NewRPCError(MethodNotFound, message)
}

func NewInvalidParamsError(message string) *RPCError {
	if message == "" {
		message = "Invalid params"
	}
	return NewRPCError(InvalidParams, message)
}

func NewInternalError(message string) *RPCError {
	if message == "" {
		message = "Internal error"
	}
	return NewRPCError(InternalError, message)
}

// MCP-specific error constructors
func NewNotInitializedError() *RPCError {
	return NewRPCError(NotInitialized, "Not initialized")
}

func NewRequestFailedError(message string) *RPCError {
	if message == "" {
		message = "Request failed"
	}
	return NewRPCError(RequestFailed, message)
}

func NewInvalidToolError(toolName string) *RPCError {
	message := "Invalid tool"
	if toolName != "" {
		message = fmt.Sprintf("Invalid tool: %s", toolName)
	}
	return NewRPCError(InvalidTool, message)
}

func NewInvalidResourceError(uri string) *RPCError {
	message := "Invalid resource"
	if uri != "" {
		message = fmt.Sprintf("Invalid resource: %s", uri)
	}
	return NewRPCError(InvalidResource, message)
}

func NewMethodDisabledError(method string) *RPCError {
	message := "Method disabled"
	if method != "" {
		message = fmt.Sprintf("Method disabled: %s", method)
	}
	return NewRPCError(MethodDisabled, message)
}

func NewInvalidPromptError(promptName string) *RPCError {
	message := "Invalid prompt"
	if promptName != "" {
		message = fmt.Sprintf("Invalid prompt: %s", promptName)
	}
	return NewRPCError(InvalidPrompt, message)
}

func NewAuthenticationFailedError(message string) *RPCError {
	if message == "" {
		message = "Authentication failed"
	}
	return NewRPCError(AuthenticationFailed, message)
}

func NewPermissionDeniedError(message string) *RPCError {
	if message == "" {
		message = "Permission denied"
	}
	return NewRPCError(PermissionDenied, message)
}

func NewRateLimitExceededError(message string) *RPCError {
	if message == "" {
		message = "Rate limit exceeded"
	}
	return NewRPCError(RateLimitExceeded, message)
}

// IsStandardError checks if an error code is a standard JSON-RPC error
func IsStandardError(code int) bool {
	return code >= -32768 && code <= -32000
}

// IsMCPError checks if an error code is an MCP-specific error
func IsMCPError(code int) bool {
	return code >= -32009 && code <= -32001
}

// GetErrorName returns a human-readable name for an error code
func GetErrorName(code int) string {
	switch code {
	case ParseError:
		return "ParseError"
	case InvalidRequest:
		return "InvalidRequest"
	case MethodNotFound:
		return "MethodNotFound"
	case InvalidParams:
		return "InvalidParams"
	case InternalError:
		return "InternalError"
	case NotInitialized:
		return "NotInitialized"
	case RequestFailed:
		return "RequestFailed"
	case InvalidTool:
		return "InvalidTool"
	case InvalidResource:
		return "InvalidResource"
	case MethodDisabled:
		return "MethodDisabled"
	case InvalidPrompt:
		return "InvalidPrompt"
	case AuthenticationFailed:
		return "AuthenticationFailed"
	case PermissionDenied:
		return "PermissionDenied"
	case RateLimitExceeded:
		return "RateLimitExceeded"
	default:
		if IsStandardError(code) {
			return "StandardError"
		}
		return "UnknownError"
	}
}
