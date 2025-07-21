#!/bin/bash

echo "🧠 Testing LLM-Powered Agent Swarm"
echo "================================="
echo ""

# Check if project builds
echo "📦 Building LLM Agent Swarm..."
cd /home/engineone/Projects/AI/ConduitMCP
if go build -o bin/agent_swarm_llm examples/agent_swarm_llm/main.go; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "🎯 Testing LLM-powered agent swarm framework:"
echo "   📋 Coordinator - Routes tasks using LLM reasoning"
echo "   ✍️  ContentCreator - Handles content with LLM intelligence"
echo "   📊 DataAnalyst - Performs analysis with LLM insights"
echo "   🧠 MemoryManager - Manages information with LLM understanding"
echo ""

echo "🔧 Prerequisites for full LLM functionality:"
echo "   1. Install Ollama: curl -fsSL https://ollama.ai/install.sh | sh"
echo "   2. Start Ollama: ollama serve"
echo "   3. Pull a model: ollama pull llama3.2"
echo "   4. Run: OLLAMA_URL=http://localhost:11434 OLLAMA_MODEL=llama3.2 ./bin/agent_swarm_llm"
echo ""

echo "💡 The example will fallback to rule-based logic if Ollama is not available."
echo ""

echo "🚀 Quick Test (without Ollama - rule-based fallback):"
echo "./bin/agent_swarm_llm"
echo ""

echo "🧠 Full LLM Test (with Ollama):"
echo "export OLLAMA_URL=http://localhost:11434"
echo "export OLLAMA_MODEL=llama3.2"
echo "./bin/agent_swarm_llm"
echo ""

echo "✅ LLM-Powered Agent Swarm ready!"
echo ""
echo "🔗 Related Examples:"
echo "   • examples/agent_swarm/ - Rule-based agent swarm"
echo "   • examples/agent_swarm_simple/ - Basic swarm concepts"
echo "   • examples/agent_swarm_workflows/ - Advanced workflow patterns"
echo "   • examples/ollama/ - Basic Ollama integration"
echo "   • examples/agents_ollama/ - LLM agents with tools"
