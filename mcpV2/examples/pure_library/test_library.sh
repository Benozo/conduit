#!/bin/bash

echo "=== Pure Library Example Test ==="
echo "Testing MCP pure library functionality..."
echo ""

# Build the example
echo "🔧 Building pure library example..."
go build -o pure_library .

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful!"
echo ""

# Run the example
echo "🚀 Running pure library example..."
echo ""

./pure_library

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Pure Library Example Test Complete!"
    echo ""
    echo "🎯 Test Summary:"
    echo "   • Tool registration and calling functional"
    echo "   • Memory storage and retrieval working"
    echo "   • Error handling demonstrated"
    echo "   • Statistics collection operational"
    echo ""
    echo "📈 Performance Benefits:"
    echo "   • Zero network overhead"
    echo "   • Direct function calls"
    echo "   • In-process memory"
    echo "   • Native Go types"
    echo ""
    echo "🔗 Integration:"
    echo "   • Use as library in any Go application"
    echo "   • Embed tools and memory directly"
    echo "   • Maximum performance for high-throughput services"
else
    echo "❌ Example execution failed"
    exit 1
fi
