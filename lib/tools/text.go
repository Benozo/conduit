package tools

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/benozo/conduit/mcp"
)

// ToolRegistrar defines the interface for registering tools
type ToolRegistrar interface {
	RegisterTool(string, mcp.ToolFunc)
}

// RegisterTextTools adds comprehensive text manipulation tools
func RegisterTextTools(server ToolRegistrar) {
	server.RegisterTool("uppercase", UppercaseFunc)
	server.RegisterTool("lowercase", LowercaseFunc)
	server.RegisterTool("reverse", ReverseFunc)
	server.RegisterTool("word_count", WordCountFunc)
	server.RegisterTool("trim", TrimFunc)
	server.RegisterTool("title_case", TitleCaseFunc)
	server.RegisterTool("snake_case", SnakeCaseFunc)
	server.RegisterTool("camel_case", CamelCaseFunc)
	server.RegisterTool("replace", ReplaceFunc)
	server.RegisterTool("extract_words", ExtractWordsFunc)
	server.RegisterTool("sort_words", SortWordsFunc)
	server.RegisterTool("char_count", CharCountFunc)
	server.RegisterTool("remove_whitespace", RemoveWhitespaceFunc)
}

var UppercaseFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	return map[string]string{"result": strings.ToUpper(text)}, nil
}

var LowercaseFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	return map[string]string{"result": strings.ToLower(text)}, nil
}

var ReverseFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return map[string]string{"result": string(runes)}, nil
}

var WordCountFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	words := strings.Fields(text)
	count := len(words)
	return map[string]interface{}{
		"result":     count,
		"text":       text,
		"word_count": count,
		"words":      words,
	}, nil
}

var TrimFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	return map[string]string{"result": strings.TrimSpace(text)}, nil
}

var TitleCaseFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	return map[string]string{"result": strings.Title(strings.ToLower(text))}, nil
}

var SnakeCaseFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])

	// Replace spaces and special characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	snake := re.ReplaceAllString(text, "_")

	// Convert to lowercase and remove leading/trailing underscores
	snake = strings.ToLower(strings.Trim(snake, "_"))

	return map[string]string{"result": snake}, nil
}

var CamelCaseFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])

	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	var camel strings.Builder
	for i, word := range words {
		if word == "" {
			continue
		}
		if i == 0 {
			camel.WriteString(strings.ToLower(word))
		} else {
			camel.WriteString(strings.Title(strings.ToLower(word)))
		}
	}

	return map[string]string{"result": camel.String()}, nil
}

var ReplaceFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	old := fmt.Sprintf("%v", params["old"])
	new := fmt.Sprintf("%v", params["new"])

	result := strings.ReplaceAll(text, old, new)

	return map[string]interface{}{
		"result": result,
		"old":    old,
		"new":    new,
		"count":  strings.Count(text, old),
	}, nil
}

var ExtractWordsFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])

	re := regexp.MustCompile(`\b\w+\b`)
	words := re.FindAllString(text, -1)

	return map[string]interface{}{
		"result": words,
		"count":  len(words),
	}, nil
}

var SortWordsFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])
	words := strings.Fields(text)

	// Get sort order (default: ascending)
	order := "asc"
	if orderParam, ok := params["order"]; ok {
		order = fmt.Sprintf("%v", orderParam)
	}

	sort.Strings(words)
	if order == "desc" {
		for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
			words[i], words[j] = words[j], words[i]
		}
	}

	result := strings.Join(words, " ")

	return map[string]interface{}{
		"result": result,
		"words":  words,
		"order":  order,
	}, nil
}

var CharCountFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])

	charCount := len([]rune(text))
	byteCount := len(text)
	wordCount := len(strings.Fields(text))
	lineCount := len(strings.Split(text, "\n"))

	return map[string]interface{}{
		"result":     charCount,
		"characters": charCount,
		"bytes":      byteCount,
		"words":      wordCount,
		"lines":      lineCount,
		"text":       text,
	}, nil
}

var RemoveWhitespaceFunc = func(params map[string]interface{}, memory *mcp.Memory) (interface{}, error) {
	text := fmt.Sprintf("%v", params["text"])

	// Remove all whitespace characters
	re := regexp.MustCompile(`\s+`)
	result := re.ReplaceAllString(text, "")

	return map[string]interface{}{
		"result":          result,
		"original_length": len(text),
		"new_length":      len(result),
		"removed_chars":   len(text) - len(result),
	}, nil
}
