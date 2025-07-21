# builtin_tools_test

## 🧠 What It Does

Tests all 31 built-in MCP tools to verify they work correctly. This example runs comprehensive validation tests on text manipulation, memory management, utility functions, and encoding tools without requiring external services.

## ⚙️ Requirements

- **Go 1.21+** - For building and running
- **No external services** - Runs completely offline

## 🚀 How to Run

```bash
# Run the comprehensive tool test suite
go run main.go
```

## 🧪 Test Categories

The test suite validates these tool categories:

- **Text Tools** - `uppercase`, `lowercase`, `reverse`, `word_count`, `trim`
- **Memory Tools** - `remember`, `recall`, `forget`, `list_memories`, `clear_memory`
- **Encoding Tools** - `base64_encode`, `base64_decode`, `url_encode`, `url_decode`
- **Hash Tools** - `hash_md5`, `hash_sha256`
- **Utility Tools** - `timestamp`, `uuid`, `random_number`, `random_string`
- **JSON Tools** - `json_format`, `json_minify`

## ✅ Sample Output

```bash
🔧 Conduit Built-in Tools Test Suite
===================================

📝 Testing Text Tools...
✅ uppercase: "hello world" → "HELLO WORLD"
✅ lowercase: "HELLO WORLD" → "hello world" 
✅ reverse: "hello" → "olleh"
✅ word_count: "hello world test" → 3 words
✅ trim: "  hello  " → "hello"
✅ title_case: "hello world" → "Hello World"
✅ snake_case: "Hello World" → "hello_world"
✅ camel_case: "hello world" → "helloWorld"

💾 Testing Memory Tools...
✅ remember: Stored "test_key" = "test_value"
✅ recall: Retrieved "test_key" → "test_value"
✅ list_memories: Found 1 stored item
✅ forget: Removed "test_key"
✅ clear_memory: Cleared all memories

🔐 Testing Encoding Tools...
✅ base64_encode: "Hello" → "SGVsbG8="
✅ base64_decode: "SGVsbG8=" → "Hello"
✅ url_encode: "Hello World!" → "Hello%20World%21"
✅ url_decode: "Hello%20World%21" → "Hello World!"

🔒 Testing Hash Tools...
✅ hash_md5: "hello" → "5d41402abc4b2a76b9719d911017c592"
✅ hash_sha256: "hello" → "2cf24dba4f..."

🛠️ Testing Utility Tools...
✅ uuid: Generated valid UUID → "550e8400-e29b-41d4-a716-446655440000"
✅ timestamp: Generated ISO timestamp → "2025-07-22T10:30:00Z"
✅ random_number: Generated random in range → 42
✅ random_string: Generated 10-char string → "aBc123XyZ9"

📄 Testing JSON Tools...
✅ json_format: Formatted JSON with proper indentation
✅ json_minify: Minified JSON removing whitespace

🎉 Test Results: 31/31 tools passed
✅ All built-in tools are working correctly!
```

## 🔧 Tools Tested

### Text Processing (8 tools)
- `uppercase`, `lowercase`, `reverse`, `word_count`
- `trim`, `title_case`, `snake_case`, `camel_case`

### Memory Management (5 tools)  
- `remember`, `recall`, `forget`, `list_memories`, `clear_memory`

### Encoding & Decoding (4 tools)
- `base64_encode`, `base64_decode`, `url_encode`, `url_decode`

### Cryptographic Hashing (2 tools)
- `hash_md5`, `hash_sha256`

### Utility Functions (6 tools)
- `uuid`, `timestamp`, `random_number`, `random_string`, `char_count`, `extract_words`

### JSON Processing (2 tools)
- `json_format`, `json_minify`

### Advanced Text (4 tools)
- `replace`, `sort_words`, `remove_whitespace`, `memory_stats`

## 🎯 Key Features

- ✅ **Comprehensive Coverage**: Tests all 31 built-in tools
- ✅ **No Dependencies**: Runs offline without external services
- ✅ **Detailed Output**: Shows input, output, and validation for each tool
- ✅ **Error Detection**: Catches and reports any tool failures
- ✅ **Performance Timing**: Measures execution time for each category
- ✅ **Memory Validation**: Verifies memory persistence and cleanup

## 🔍 How It Works

1. **Initialization** → Creates MCP server and registers all tools
2. **Category Testing** → Runs tests grouped by tool type
3. **Individual Validation** → Tests each tool with known inputs
4. **Result Verification** → Validates outputs match expected results
5. **Summary Report** → Shows pass/fail count and performance metrics

## ⚠️ Troubleshooting

**Build Errors:**
```bash
# Ensure dependencies are current
go mod tidy

# Check Go version
go version  # Should be 1.21+
```

**Test Failures:**
- If any tools fail, check the detailed error output
- Memory tools require read/write permissions in current directory
- JSON tools require valid JSON input format

## 📚 Related Examples

- [`custom_tools/`](../custom_tools) - Create your own tools with validation
- [`pure_library/`](../pure_library) - Use tools in library mode
- [`stdio_example/`](../stdio_example) - Tools via MCP protocol
- [`agents_test/`](../agents_test) - Test agent functionality

## 🚀 Next Steps

After verifying tools work:

1. **Create Custom Tools**: Use [`custom_tools/`](../custom_tools) to build your own
2. **Try with LLMs**: Use tools with [`ollama/`](../ollama) or [`openai/`](../openai)
3. **Build Agents**: Create AI agents that use these tools
4. **Production Deploy**: Use in your applications via MCP protocol

## 🧪 Development Usage

This test suite is also useful for:

- **CI/CD Validation** - Ensure tools work in different environments
- **Regression Testing** - Verify updates don't break existing functionality  
- **Performance Benchmarking** - Measure tool execution times
- **Documentation Validation** - Confirm examples in docs are accurate
