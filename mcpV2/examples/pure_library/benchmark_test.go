package main

import (
	"context"
	"testing"

	"github.com/benozo/neuron-mcp/library"
	"github.com/benozo/neuron-mcp/protocol"
)

// BenchmarkToolCall benchmarks direct tool calls in library mode
func BenchmarkToolCall(b *testing.B) {
	registry := library.NewComponentRegistry()

	// Register a simple tool
	registry.Tools().Register("benchmark_tool", func(ctx context.Context, params map[string]interface{}) (*protocol.ToolResult, error) {
		return &protocol.ToolResult{
			Content: []protocol.Content{{
				Type: "text",
				Text: "benchmark result",
			}},
		}, nil
	})

	ctx := context.Background()
	params := map[string]interface{}{"test": "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Tools().Call(ctx, "benchmark_tool", params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMemoryOperations benchmarks memory operations
func BenchmarkMemoryOperations(b *testing.B) {
	registry := library.NewComponentRegistry()
	memory := registry.Memory()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := memory.Set("key", map[string]interface{}{
				"value": i,
				"data":  "test data",
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Get", func(b *testing.B) {
		memory.Set("test_key", "test_value")
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := memory.Get("test_key")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkTextTransform benchmarks the text transform tool
func BenchmarkTextTransform(b *testing.B) {
	registry := library.NewComponentRegistry()
	registerLibraryTools(registry)

	ctx := context.Background()
	params := map[string]interface{}{
		"text":      "Hello, World! This is a longer text for benchmarking purposes.",
		"operation": "uppercase",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Tools().Call(ctx, "text_transform", params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCalculator benchmarks the calculator tool
func BenchmarkCalculator(b *testing.B) {
	registry := library.NewComponentRegistry()
	registerLibraryTools(registry)

	ctx := context.Background()
	params := map[string]interface{}{
		"operation": "multiply",
		"a":         123.456,
		"b":         789.012,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Tools().Call(ctx, "calculator", params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
