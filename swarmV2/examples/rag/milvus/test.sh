#!/bin/bash

# Test script for Milvus RAG example
echo "=== Milvus RAG Example Test Script ==="

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "❌ Error: main.go not found. Please run this script from the milvus example directory."
    exit 1
fi

echo "🔍 Checking prerequisites..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed or not in PATH"
    exit 1
fi

echo "✅ Go is available: $(go version)"

# Check Go modules
echo "🔧 Checking Go modules..."
if ! go mod tidy; then
    echo "❌ Failed to tidy Go modules"
    exit 1
fi

echo "✅ Go modules are valid"

# Check Milvus connectivity (optional)
echo "🔍 Checking Milvus connectivity..."
if curl -f -s http://localhost:19530/healthz > /dev/null 2>&1; then
    echo "✅ Milvus is accessible at localhost:19530"
else
    echo "⚠️  Milvus is not accessible at localhost:19530"
    echo "   Make sure to start Milvus with: docker-compose up -d etcd minio milvus"
fi

# Check Ollama connectivity (optional)
echo "🔍 Checking Ollama connectivity..."
OLLAMA_URL="${OLLAMA_URL:-http://192.168.10.10:11434}"
if curl -f -s "$OLLAMA_URL/api/tags" > /dev/null 2>&1; then
    echo "✅ Ollama is accessible at $OLLAMA_URL"
    
    # Check if llama3.2 model is available
    if curl -s "$OLLAMA_URL/api/tags" | grep -q "llama3.2"; then
        echo "✅ llama3.2 model is available"
    else
        echo "⚠️  llama3.2 model not found. You can pull it with: ollama pull llama3.2"
    fi
else
    echo "⚠️  Ollama is not accessible at $OLLAMA_URL"
    echo "   Make sure Ollama is running with: ollama serve"
fi

# Test compilation
echo "🔨 Testing compilation..."
if go build -o /tmp/milvus-rag-test main.go; then
    echo "✅ Code compiles successfully"
    rm -f /tmp/milvus-rag-test
else
    echo "❌ Compilation failed"
    exit 1
fi

echo ""
echo "🎉 All checks passed! You can now run the example with:"
echo "   go run main.go"
echo ""
echo "💡 Optional environment variables:"
echo "   export MILVUS_URL=http://localhost:19530"
echo "   export OLLAMA_URL=http://192.168.10.10:11434"
echo "   export OLLAMA_MODEL=llama3.2"
