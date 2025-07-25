// Package middleware provides a flexible middleware system for MCP requests and responses.
//
// Middleware allows for request/response interception, modification, and handling
// of cross-cutting concerns like logging, authentication, rate limiting, and metrics.
package middleware

import (
	"context"
	"time"

	"github.com/benozo/neuron-mcp/protocol"
)

// Handler represents the final handler in the middleware chain
type Handler func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error)

// Middleware represents a middleware function that can intercept and modify requests/responses
type Middleware func(next Handler) Handler

// Chain represents a chain of middleware
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new middleware chain
func NewChain(middlewares ...Middleware) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

// Add adds middleware to the chain
func (c *Chain) Add(middleware Middleware) *Chain {
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// Build builds the middleware chain with the final handler
func (c *Chain) Build(handler Handler) Handler {
	// Build the chain in reverse order
	final := handler
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		final = c.middlewares[i](final)
	}
	return final
}

// Logger interface for middleware logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// Metrics interface for middleware metrics
type Metrics interface {
	IncrementCounter(name string, tags map[string]string)
	RecordDuration(name string, duration time.Duration, tags map[string]string)
	RecordGauge(name string, value float64, tags map[string]string)
}

// RequestContext represents the context for middleware
type RequestContext struct {
	RequestID string
	Method    string
	UserID    string
	ClientIP  string
	StartTime time.Time
	Metadata  map[string]interface{}
}

// GetRequestContext extracts request context from context.Context
func GetRequestContext(ctx context.Context) *RequestContext {
	if reqCtx, ok := ctx.Value("requestContext").(*RequestContext); ok {
		return reqCtx
	}
	return &RequestContext{
		Metadata: make(map[string]interface{}),
	}
}

// WithRequestContext adds request context to context.Context
func WithRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	return context.WithValue(ctx, "requestContext", reqCtx)
}

// LoggingMiddleware logs all requests and responses
func LoggingMiddleware(logger Logger) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			reqCtx := GetRequestContext(ctx)
			start := time.Now()

			logger.Info("Request started",
				"method", req.Method,
				"id", req.ID,
				"requestId", reqCtx.RequestID,
			)

			resp, err := next(ctx, req)

			duration := time.Since(start)

			if err != nil {
				logger.Error("Request failed",
					"method", req.Method,
					"id", req.ID,
					"error", err.Error(),
					"duration", duration,
					"requestId", reqCtx.RequestID,
				)
			} else {
				logger.Info("Request completed",
					"method", req.Method,
					"id", req.ID,
					"duration", duration,
					"requestId", reqCtx.RequestID,
				)
			}

			return resp, err
		}
	}
}

// MetricsMiddleware collects metrics for requests
func MetricsMiddleware(metrics Metrics) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			start := time.Now()
			reqCtx := GetRequestContext(ctx)

			tags := map[string]string{
				"method": req.Method,
			}
			if reqCtx.UserID != "" {
				tags["user_id"] = reqCtx.UserID
			}

			metrics.IncrementCounter("mcp.requests.total", tags)

			resp, err := next(ctx, req)

			duration := time.Since(start)
			metrics.RecordDuration("mcp.requests.duration", duration, tags)

			if err != nil {
				errorTags := make(map[string]string)
				for k, v := range tags {
					errorTags[k] = v
				}
				errorTags["status"] = "error"
				metrics.IncrementCounter("mcp.requests.errors", errorTags)
			} else {
				successTags := make(map[string]string)
				for k, v := range tags {
					successTags[k] = v
				}
				successTags["status"] = "success"
				metrics.IncrementCounter("mcp.requests.success", successTags)
			}

			return resp, err
		}
	}
}

// AuthenticationMiddleware validates authentication tokens
func AuthenticationMiddleware(validator TokenValidator) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			// Skip authentication for initialize and ping methods
			if req.Method == "initialize" || req.Method == "ping" {
				return next(ctx, req)
			}

			token := extractAuthToken(req)
			if token == "" {
				return nil, protocol.NewRPCError(protocol.AuthenticationFailed, "authentication required")
			}

			user, err := validator.ValidateToken(ctx, token)
			if err != nil {
				return nil, protocol.NewRPCError(protocol.AuthenticationFailed, "invalid authentication token")
			}

			// Add user to request context
			reqCtx := GetRequestContext(ctx)
			reqCtx.UserID = user.ID
			ctx = WithRequestContext(ctx, reqCtx)

			return next(ctx, req)
		}
	}
}

// TokenValidator validates authentication tokens
type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (*User, error)
}

// User represents an authenticated user
type User struct {
	ID       string
	Username string
	Roles    []string
}

// RateLimitMiddleware limits request rates per user/IP
func RateLimitMiddleware(limiter RateLimiter) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			reqCtx := GetRequestContext(ctx)

			// Use user ID if available, otherwise use client IP
			key := reqCtx.UserID
			if key == "" {
				key = reqCtx.ClientIP
			}

			if !limiter.Allow(key) {
				return nil, protocol.NewRPCError(protocol.RateLimitExceeded, "rate limit exceeded")
			}

			return next(ctx, req)
		}
	}
}

// RateLimiter interface for rate limiting
type RateLimiter interface {
	Allow(key string) bool
}

// ValidationMiddleware validates request parameters
func ValidationMiddleware() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			// Basic JSON-RPC validation
			if req.Version != "2.0" {
				return nil, protocol.NewRPCError(protocol.InvalidRequest, "invalid JSON-RPC version")
			}

			if req.Method == "" && req.ID != nil {
				return nil, protocol.NewRPCError(protocol.InvalidRequest, "missing method")
			}

			return next(ctx, req)
		}
	}
}

// ErrorHandlingMiddleware provides centralized error handling
func ErrorHandlingMiddleware(logger Logger) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Panic in request handler",
						"method", req.Method,
						"panic", r,
					)
				}
			}()

			resp, err := next(ctx, req)

			// Log unexpected errors
			if err != nil {
				if rpcErr, ok := err.(*protocol.RPCError); ok {
					if rpcErr.Code == protocol.InternalError {
						logger.Error("Internal error",
							"method", req.Method,
							"error", err.Error(),
						)
					}
				} else {
					logger.Error("Unexpected error",
						"method", req.Method,
						"error", err.Error(),
					)
				}
			}

			return resp, err
		}
	}
}

// extractAuthToken extracts authentication token from request
func extractAuthToken(req *protocol.JSONRPCMessage) string {
	// Check meta for auth token
	if req.Meta != nil && req.Meta.Extra != nil {
		if token, ok := req.Meta.Extra["authToken"].(string); ok {
			return token
		}
	}

	// Could also check params or other locations
	return ""
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			done := make(chan struct{})
			var resp *protocol.JSONRPCMessage
			var err error

			go func() {
				resp, err = next(timeoutCtx, req)
				close(done)
			}()

			select {
			case <-done:
				return resp, err
			case <-timeoutCtx.Done():
				return nil, protocol.NewRPCError(protocol.RequestFailed, "request timeout")
			}
		}
	}
}

// CompressionMiddleware handles request/response compression
func CompressionMiddleware() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *protocol.JSONRPCMessage) (*protocol.JSONRPCMessage, error) {
			// Implementation would handle compression/decompression
			// This is a placeholder for the interface
			return next(ctx, req)
		}
	}
}
