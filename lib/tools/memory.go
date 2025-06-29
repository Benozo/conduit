package tools

import (
	"fmt"
	"time"

	"github.com/benozo/conduit/mcp"
)

// RegisterMemoryTools adds memory management tools
func RegisterMemoryTools(server ToolRegistrar) {
	server.RegisterTool("remember", RememberFunc)
	server.RegisterTool("recall", RecallFunc)
	server.RegisterTool("forget", ForgetFunc)
	server.RegisterTool("list_memories", ListMemoriesFunc)
	server.RegisterTool("clear_memory", ClearMemoryFunc)
	server.RegisterTool("memory_stats", MemoryStatsFunc)
}

var RememberFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	key := fmt.Sprintf("%v", params["key"])
	value := params["value"]

	// Add timestamp to the value
	valueWithMeta := map[string]interface{}{
		"value":     value,
		"timestamp": time.Now().Unix(),
		"type":      fmt.Sprintf("%T", value),
	}

	memory.Set(key, valueWithMeta)

	return map[string]interface{}{
		"result":    fmt.Sprintf("Remembered %s", key),
		"key":       key,
		"timestamp": time.Now().Unix(),
	}, nil
}

var RecallFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	key := fmt.Sprintf("%v", params["key"])
	stored := memory.Get(key)

	if stored == nil {
		return map[string]interface{}{
			"result": "Not found",
			"key":    key,
			"found":  false,
		}, nil
	}

	// Handle both new format (with metadata) and old format (direct value)
	if storedMap, ok := stored.(map[string]interface{}); ok {
		if value, exists := storedMap["value"]; exists {
			return map[string]interface{}{
				"result":    value,
				"key":       key,
				"found":     true,
				"timestamp": storedMap["timestamp"],
				"type":      storedMap["type"],
			}, nil
		}
	}

	// Fallback for old format
	return map[string]interface{}{
		"result": stored,
		"key":    key,
		"found":  true,
	}, nil
}

var ForgetFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	key := fmt.Sprintf("%v", params["key"])

	// Check if key exists before forgetting
	exists := memory.Get(key) != nil

	memory.Set(key, nil)

	return map[string]interface{}{
		"result":  fmt.Sprintf("Forgot %s", key),
		"key":     key,
		"existed": exists,
	}, nil
}

var ListMemoriesFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	all := memory.All()
	keys := make([]string, 0, len(all))
	memories := make(map[string]interface{})

	for key, value := range all {
		if value != nil {
			keys = append(keys, key)

			// Extract just the value for listing, not the full metadata
			if valueMap, ok := value.(map[string]interface{}); ok {
				if actualValue, exists := valueMap["value"]; exists {
					memories[key] = map[string]interface{}{
						"value":     actualValue,
						"timestamp": valueMap["timestamp"],
						"type":      valueMap["type"],
					}
					continue
				}
			}

			// Fallback for old format
			memories[key] = value
		}
	}

	return map[string]interface{}{
		"result":   keys,
		"count":    len(keys),
		"memories": memories,
	}, nil
}

var ClearMemoryFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	all := memory.All()
	count := 0

	// Count non-nil entries
	for _, value := range all {
		if value != nil {
			count++
		}
	}

	// Clear all memories
	for key := range all {
		memory.Set(key, nil)
	}

	return map[string]interface{}{
		"result":  "Memory cleared",
		"cleared": count,
	}, nil
}

var MemoryStatsFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	all := memory.All()

	totalKeys := len(all)
	activeKeys := 0
	typeCount := make(map[string]int)
	oldestTimestamp := int64(0)
	newestTimestamp := int64(0)

	for _, value := range all {
		if value != nil {
			activeKeys++

			// Analyze metadata if available
			if valueMap, ok := value.(map[string]interface{}); ok {
				if typeStr, exists := valueMap["type"]; exists {
					if typeString, ok := typeStr.(string); ok {
						typeCount[typeString]++
					}
				}

				if timestamp, exists := valueMap["timestamp"]; exists {
					if ts, ok := timestamp.(int64); ok {
						if oldestTimestamp == 0 || ts < oldestTimestamp {
							oldestTimestamp = ts
						}
						if ts > newestTimestamp {
							newestTimestamp = ts
						}
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"result":            fmt.Sprintf("%d active memories out of %d total keys", activeKeys, totalKeys),
		"total_keys":        totalKeys,
		"active_keys":       activeKeys,
		"inactive_keys":     totalKeys - activeKeys,
		"type_distribution": typeCount,
		"oldest_timestamp":  oldestTimestamp,
		"newest_timestamp":  newestTimestamp,
	}, nil
}
