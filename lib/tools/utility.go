package tools

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/benozo/conduit/mcp"
)

// RegisterUtilityTools adds general utility tools
func RegisterUtilityTools(server ToolRegistrar) {
	server.RegisterTool("timestamp", TimestampFunc)
	server.RegisterTool("uuid", UUIDFunc)
	server.RegisterTool("base64_encode", Base64EncodeFunc)
	server.RegisterTool("base64_decode", Base64DecodeFunc)
	server.RegisterTool("url_encode", URLEncodeFunc)
	server.RegisterTool("url_decode", URLDecodeFunc)
	server.RegisterTool("hash_md5", HashMD5Func)
	server.RegisterTool("hash_sha256", HashSHA256Func)
	server.RegisterTool("json_format", JSONFormatFunc)
	server.RegisterTool("json_minify", JSONMinifyFunc)
	server.RegisterTool("random_number", RandomNumberFunc)
	server.RegisterTool("random_string", RandomStringFunc)
}

var TimestampFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	now := time.Now()

	format := "unix"
	if formatParam, ok := params["format"]; ok {
		format = fmt.Sprintf("%v", formatParam)
	}

	var result interface{}
	switch format {
	case "unix":
		result = now.Unix()
	case "iso":
		result = now.Format(time.RFC3339)
	case "rfc":
		result = now.Format(time.RFC1123)
	case "custom":
		layout := "2006-01-02 15:04:05"
		if layoutParam, ok := params["layout"]; ok {
			layout = fmt.Sprintf("%v", layoutParam)
		}
		result = now.Format(layout)
	default:
		result = now.Unix()
	}

	return map[string]interface{}{
		"result": result,
		"format": format,
		"unix":   now.Unix(),
		"iso":    now.Format(time.RFC3339),
		"rfc":    now.Format(time.RFC1123),
	}, nil
}

var UUIDFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	// Simple UUID v4 generation
	uuid := make([]byte, 16)
	rand.Read(uuid)

	// Set version (4) and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	result := fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])

	return map[string]string{"result": result}, nil
}

var Base64EncodeFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	encoded := base64.StdEncoding.EncodeToString([]byte(text))

	return map[string]interface{}{
		"result":        encoded,
		"original":      text,
		"original_size": len(text),
		"encoded_size":  len(encoded),
	}, nil
}

var Base64DecodeFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	encoded := fmt.Sprintf("%v", params["text"])

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return map[string]interface{}{
			"result": "",
			"error":  fmt.Sprintf("Invalid base64: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"result":       string(decoded),
		"encoded":      encoded,
		"decoded_size": len(decoded),
		"encoded_size": len(encoded),
	}, nil
}

var URLEncodeFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	encoded := url.QueryEscape(text)

	return map[string]interface{}{
		"result":   encoded,
		"original": text,
	}, nil
}

var URLDecodeFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	encoded := fmt.Sprintf("%v", params["text"])

	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		return map[string]interface{}{
			"result": "",
			"error":  fmt.Sprintf("Invalid URL encoding: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"result":  decoded,
		"encoded": encoded,
	}, nil
}

var HashMD5Func = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	hash := md5.Sum([]byte(text))
	result := hex.EncodeToString(hash[:])

	return map[string]interface{}{
		"result":    result,
		"original":  text,
		"algorithm": "MD5",
	}, nil
}

var HashSHA256Func = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	hash := sha256.Sum256([]byte(text))
	result := hex.EncodeToString(hash[:])

	return map[string]interface{}{
		"result":    result,
		"original":  text,
		"algorithm": "SHA256",
	}, nil
}

var JSONFormatFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	jsonStr := fmt.Sprintf("%v", params["text"])

	var parsed interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return map[string]interface{}{
			"result": jsonStr,
			"error":  fmt.Sprintf("Invalid JSON: %v", err),
		}, nil
	}

	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return map[string]interface{}{
			"result": jsonStr,
			"error":  fmt.Sprintf("Failed to format: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"result":         string(formatted),
		"original":       jsonStr,
		"original_size":  len(jsonStr),
		"formatted_size": len(formatted),
	}, nil
}

var JSONMinifyFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	jsonStr := fmt.Sprintf("%v", params["text"])

	var parsed interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return map[string]interface{}{
			"result": jsonStr,
			"error":  fmt.Sprintf("Invalid JSON: %v", err),
		}, nil
	}

	minified, err := json.Marshal(parsed)
	if err != nil {
		return map[string]interface{}{
			"result": jsonStr,
			"error":  fmt.Sprintf("Failed to minify: %v", err),
		}, nil
	}

	return map[string]interface{}{
		"result":        string(minified),
		"original":      jsonStr,
		"original_size": len(jsonStr),
		"minified_size": len(minified),
		"space_saved":   len(jsonStr) - len(minified),
	}, nil
}

var RandomNumberFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	min := 0
	max := 100

	if minParam, ok := params["min"]; ok {
		if minInt, err := strconv.Atoi(fmt.Sprintf("%v", minParam)); err == nil {
			min = minInt
		}
	}

	if maxParam, ok := params["max"]; ok {
		if maxInt, err := strconv.Atoi(fmt.Sprintf("%v", maxParam)); err == nil {
			max = maxInt
		}
	}

	if min >= max {
		return map[string]interface{}{
			"result": 0,
			"error":  "min must be less than max",
		}, nil
	}

	result := mathrand.Intn(max-min) + min

	return map[string]interface{}{
		"result": result,
		"min":    min,
		"max":    max,
	}, nil
}

var RandomStringFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	length := 10
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	if lengthParam, ok := params["length"]; ok {
		if lengthInt, err := strconv.Atoi(fmt.Sprintf("%v", lengthParam)); err == nil && lengthInt > 0 {
			length = lengthInt
		}
	}

	if charsetParam, ok := params["charset"]; ok {
		charsetStr := fmt.Sprintf("%v", charsetParam)
		if charsetStr != "" {
			charset = charsetStr
		}
	}

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[mathrand.Intn(len(charset))]
	}

	return map[string]interface{}{
		"result":  string(result),
		"length":  length,
		"charset": charset,
	}, nil
}
